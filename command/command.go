package command

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"gitlab.bbinfra.net/3estack/alexandria/util"
	"log"
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
	log.Println("Entering commandExecution.doRun...")
	defer log.Println("Left commandExecution.doRun")
	log.Printf("Have model command: %s\n", ce.command.String())
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
	log.Println("Entering commandExecution.checkBootstrap")
	defer log.Println("Left commandExecution.checkBootstrap")
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
			&singleUpdateSettingsCreate{
				timestamp: ce.command.Timestamp,
				priceList: bootstrap.PriceList,
			},
			&singleUpdatePersonCreate{
				timestamp:    ce.command.Timestamp,
				personCreate: bootstrap.FirstMajor,
				isSigned:     true,
				isMajor:      true,
			},
		},
	}, nil
}

func (ce *commandExecution) checkNonBootstrap() (*updater, error) {
	u := newUnmarshalledState()
	addressesToRead := []string{model.GetSettingsAddress(), ce.command.Signer}
	addressData, err := ce.blockchainAccess.GetState(addressesToRead)
	if err != nil {
		return nil, err
	}
	err = u.add(addressData, addressesToRead)
	if err != nil {
		return nil, err
	}
	if u.getAddressState(model.GetSettingsAddress()) != ADDRESS_FILLED {
		return nil, errors.New("Blockchain has not been bootstrapped")
	}
	if u.getAddressState(ce.command.Signer) != ADDRESS_FILLED {
		return nil, errors.New("Signer does not exist: " + ce.command.Signer)
	}
	if u.persons[ce.command.Signer].PublicKey != ce.signerKey {
		return nil, errors.New(fmt.Sprintf("Signer id does not match signing key, id = %s, key = %s",
			ce.command.Signer, ce.signerKey))
	}
	balance := u.persons[ce.command.Signer].Balance
	if balance < ce.command.Price {
		return nil, errors.New(fmt.Sprintf("Insufficient balance, got %d need %d",
			balance, ce.command.Price))
	}
	nbce := &nonBootstrapCommandExecution{
		verifiedSignerId:  ce.command.Signer,
		price:             ce.command.Price,
		timestamp:         ce.command.Timestamp,
		blockchainAccess:  ce.blockchainAccess,
		unmarshalledState: u,
	}
	return nbce.check(ce.command)
}

func (ce *commandExecution) runUpdater(u *updater) error {
	log.Println("Entering commandExecution.runUpdater...")
	defer log.Println("Left commandExecution.runUpdater")
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
	dataToWrite, err := u.unmarshalledState.read(util.MapStringBoolToSlice(addressesToWrite))
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
	eventType := model.AlexandriaPrefix + model.EV_TYPE_TRANSACTION_CONTROL
	log.Println("Sending event of type: " + eventType)
	err := ce.blockchainAccess.AddEvent(
		eventType,
		[]processor.Attribute{
			{
				Key:   model.EV_KEY_TRANSACTION_ID,
				Value: ce.transactionId,
			},
			{
				Key:   model.EV_KEY_EVENT_SEQ,
				Value: eventSeqString,
			},
			{
				Key:   model.EV_KEY_TIMESTAMP,
				Value: timestampString,
			},
			{
				Key:   model.EV_KEY_NUM_EVENTS,
				Value: numEventsString,
			},
		},
		[]byte{})
	return err
}

type nonBootstrapCommandExecution struct {
	verifiedSignerId  string
	price             int32
	timestamp         int64
	blockchainAccess  BlockchainAccess
	unmarshalledState *unmarshalledState
}

func (nbce *nonBootstrapCommandExecution) check(c *model.Command) (*updater, error) {
	result, err := nbce.checkSpecific(c)
	if err == nil {
		nbce.addBalanceDeduct(result)
	}
	return result, err
}

func (nbce *nonBootstrapCommandExecution) checkSpecific(c *model.Command) (*updater, error) {
	switch c.Body.(type) {
	case *model.Command_CommandSettingsUpdate:
		return nbce.checkSettingsUpdate(c.GetCommandSettingsUpdate())
	case *model.Command_CommandJournalCreate:
		return nbce.checkJournalCreate(c.GetCommandJournalCreate())
	case *model.Command_CommandJournalUpdateProperties:
		return nbce.checkJournalUpdateProperties(c.GetCommandJournalUpdateProperties())
	case *model.Command_PersonCreate:
		return nbce.checkPersonCreate(c.GetPersonCreate())
	case *model.Command_CommandPersonUpdateProperties:
		return nbce.checkPersonUpdateProperties(c.GetCommandPersonUpdateProperties())
	case *model.Command_CommandUpdateAuthorization:
		return nbce.checkPersonUpdateAuthorization(c.GetCommandUpdateAuthorization())
	case *model.Command_CommandPersonUpdateBalanceIncrement:
		return nbce.checkPersonUpdateIncBalance(c.GetCommandPersonUpdateBalanceIncrement())
	default:
		return nil, errors.New("Non-bootstrap command type not supported")
	}
}

func (nbce *nonBootstrapCommandExecution) addBalanceDeduct(u *updater) {
	// Balance deduction should happen first. If a balance increment update
	// had come before the deduction with the price, the updates would conflict.
	// Each update sets the balance to a pre-calculated value. Now this
	// conflict won't occur because incrementing the balance has price zero.
	var deductUpdate singleUpdate = &singleUpdatePersonIncBalance{
		personId:   nbce.verifiedSignerId,
		newBalance: nbce.unmarshalledState.persons[nbce.verifiedSignerId].Balance - nbce.price,
		timestamp:  nbce.timestamp,
	}
	// This one may be duplicate, but them both will be identical because
	// the timestamp comes from the client.
	var updateModificationTime singleUpdate = &singleUpdatePersonModificationTime{
		id: nbce.verifiedSignerId,
		timestamp: nbce.timestamp,
	}
	u.updates = append([]singleUpdate{deductUpdate, updateModificationTime}, u.updates...)
}
