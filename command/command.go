package command

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"github.com/mhdirkse/alexandria/util"
	"gitlab.bbinfra.net/3estack/alexandria/model"
)

type Command struct {
	InputAddresses  []string
	OutputAddresses []string
	CryptoIdentity  *CryptoIdentity
	Command         *model.Command
}

type CryptoIdentity struct {
	PublicKeyStr string
	PublicKey    signing.PublicKey
	PrivateKey   signing.PrivateKey
}

func ApplyModelCommand(
	command *model.Command,
	signerKey string,
	transactionId string,
	ba BlockchainAccess) error {
	ce := &commandExecution{
		command:          command,
		signerKey:        signerKey,
		transactionId:    transactionId,
		blockchainAccess: ba,
	}
	return ce.run()
}

type commandExecution struct {
	command          *model.Command
	signerKey        string
	transactionId    string
	blockchainAccess BlockchainAccess
}

func (ce *commandExecution) run() error {
	if ce.command == nil {
		return nil
	}
	if ce.command.Timestamp == int64(0) {
		return errors.New("Timestamp field of command was not set")
	}
	if ce.signerKey == "" {
		return errors.New("Cannot run command with empty signerKey")
	}
	if ce.transactionId == "" {
		return errors.New("Cannot run command with empty transactionId")
	}
	if ce.blockchainAccess == nil {
		return errors.New("Cannot run command with nil blockchainAccess")
	}
	if !model.IsPersonAddress(ce.command.Signer) {
		return errors.New("Bootstrap: signer id is not a person address: " + ce.command.Signer)
	}
	return ce.doRun()
}

func (ce *commandExecution) doRun() error {
	u, err := ce.check()
	if err != nil {
		return err
	}
	if u == nil {
		return nil
	}
	return ce.runUpdater(u)
}

func (ce *commandExecution) check() (*updater, error) {
	var u *updater
	var err error
	switch ce.command.Body.(type) {
	case *model.Command_Bootstrap:
		u, err = ce.checkBootstrap(ce.command.GetBootstrap())
	case nil:
		return nil, nil
	default:
		u, err = ce.checkNonBootstrap()
	}
	return u, err
}

func (ce *commandExecution) checkBootstrap(bootstrap *model.CommandBootstrap) (*updater, error) {
	if bootstrap.PriceList == nil {
		return nil, errors.New("Bootstrap: missing price list")
	}
	if bootstrap.FirstMajor == nil {
		return nil, errors.New("Bootstrap: missing first major")
	}
	us := newUnmarshalledState()
	requested := []string{model.GetSettingsAddress(), ce.command.Signer}
	data, err := ce.blockchainAccess.GetState(requested)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read settings address or signer address: %v", requested))
	}
	err = us.add(data, requested)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not unmarshal"))
	}
	if us.getAddressState(model.GetSettingsAddress()) != ADDRESS_EMPTY {
		return nil, errors.New("Blockchain was bootstrapped already")
	}
	if ce.signerKey != bootstrap.FirstMajor.PublicKey {
		return nil, errors.New("Bootstrap command is not signed with the key of the first major")
	}
	if ce.command.Price != int32(0) {
		return nil, errors.New("Bootstrapping should have price zero")
	}
	if bootstrap.FirstMajor.NewPersonId != ce.command.Signer {
		return nil, errors.New("Bootstrap: signer should be id of first major")
	}
	if err := checkSanityPersonCreate(bootstrap.FirstMajor); err != nil {
		return nil, err
	}
	return &updater{
		unmarshalledState: us,
		updates: []singleUpdate{
			&singleUpdateCreateSettings{
				timestamp: ce.command.Timestamp,
				priceList: bootstrap.PriceList,
			},
			&singleUpdatePersonCreate{
				timestamp:    ce.command.Timestamp,
				personCreate: bootstrap.FirstMajor,
			},
		},
	}, nil
}

func checkSanityPersonCreate(personCreate *model.CommandPersonCreate) error {
	if personCreate.NewPersonId == "" {
		return errors.New("personCreate.NewPersonId should be filled")
	}
	if personCreate.PublicKey == "" {
		return errors.New("personCreate.PublicKey should be filled")
	}
	if personCreate.Email == "" {
		return errors.New("personCreate.Email should be filled")
	}
	if personCreate.Name == "" {
		return errors.New("personCreate.Name should be filled")
	}
	return nil
}

func (ce *commandExecution) checkNonBootstrap() (*updater, error) {
	return nil, nil
}

func (ce *commandExecution) runUpdater(u *updater) error {
	if err := ce.updateState(u); err != nil {
		return err
	}
	return ce.writeEvents(u)
}

func (ce *commandExecution) updateState(u *updater) error {
	addressesToWrite := make(map[string]bool)
	for _, su := range u.updates {
		addressToWrite := su.updateState(u.unmarshalledState)
		addressesToWrite[addressToWrite] = true
	}
	dataToWrite, err := u.unmarshalledState.read(util.SetToSlice(addressesToWrite))
	if err != nil {
		return err
	}
	writtenAddresses, err := ce.blockchainAccess.SetState(dataToWrite)
	if err != nil {
		return err
	}
	if len(writtenAddresses) < len(addressesToWrite) {
		return errors.New(fmt.Sprintf("Could not write all addresses to be written. Tried: %v, succeeded: %v",
			addressesToWrite, writtenAddresses))
	}
	return nil
}

func (ce *commandExecution) writeEvents(u *updater) error {
	numUpdates := len(u.updates)
	for eventSeq, su := range u.updates {
		if err := su.issueEvent(int32(eventSeq), ce.transactionId, ce.blockchainAccess); err != nil {
			return err
		}
	}
	return ce.writeTransactionControlEvent(int32(numUpdates))
}

func (ce *commandExecution) writeTransactionControlEvent(numNonControlEvents int32) error {
	eventSeqString := fmt.Sprintf("%d", numNonControlEvents)
	numEventsString := fmt.Sprintf("%d", numNonControlEvents+1)
	timestampString := fmt.Sprintf("%d", ce.command.Timestamp)
	err := ce.blockchainAccess.AddEvent(
		model.EV_TRANSACTION_CONTROL,
		[]processor.Attribute{
			{
				Key:   model.TRANSACTION_ID,
				Value: ce.transactionId,
			},
			{
				Key:   model.EVENT_SEQ,
				Value: eventSeqString,
			},
			{
				Key:   model.TIMESTAMP,
				Value: timestampString,
			},
			{
				Key:   model.NUM_EVENTS,
				Value: numEventsString,
			},
		},
		[]byte{})
	return err
}
