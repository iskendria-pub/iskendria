package command

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"strconv"
)

type PersonCreate struct {
	PublicKey string
	Name      string
	Email     string
}

func GetPersonCreateCommand(
	pc *PersonCreate,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	personId := model.CreatePersonAddress()
	return &Command{
		InputAddresses:  []string{personId, signerId, model.GetSettingsAddress()},
		OutputAddresses: []string{personId},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signerId,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_PersonCreate{
				PersonCreate: &model.CommandPersonCreate{
					NewPersonId: personId,
					PublicKey:   pc.PublicKey,
					Name:        pc.Name,
					Email:       pc.Email,
				},
			},
		},
	}
}

func GetPersonUpdateCommand(
	personId string,
	orig,
	updated *dao.PersonUpdate,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	return &Command{
		InputAddresses:  []string{model.GetSettingsAddress(), personId, signerId},
		OutputAddresses: []string{personId},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signerId,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandPersonUpdateProperties{
				CommandPersonUpdateProperties: createModelCommandPersonUpdate(personId, orig, updated),
			},
		},
	}
}

func GetPersonSetMajorCommand(
	personId,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	update := &authorizationUpdate{
		majorUpdate:  model.BoolUpdate_MAKE_TRUE,
		signedUpdate: model.BoolUpdate_UNMODIFIED,
	}
	return getGenericPersonAuthorizationUpdateCommand(personId, signerId, cryptoIdentity, price, update)
}

func GetPersonUnsetMajorCommand(
	personId,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	update := &authorizationUpdate{
		majorUpdate:  model.BoolUpdate_MAKE_FALSE,
		signedUpdate: model.BoolUpdate_UNMODIFIED,
	}
	return getGenericPersonAuthorizationUpdateCommand(personId, signerId, cryptoIdentity, price, update)
}

func GetPersonSetSignedCommand(
	personId,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	update := &authorizationUpdate{
		majorUpdate:  model.BoolUpdate_UNMODIFIED,
		signedUpdate: model.BoolUpdate_MAKE_TRUE,
	}
	return getGenericPersonAuthorizationUpdateCommand(personId, signerId, cryptoIdentity, price, update)
}

func GetPersonUnsetSignedCommand(
	personId,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	update := &authorizationUpdate{
		majorUpdate:  model.BoolUpdate_UNMODIFIED,
		signedUpdate: model.BoolUpdate_MAKE_FALSE,
	}
	return getGenericPersonAuthorizationUpdateCommand(personId, signerId, cryptoIdentity, price, update)
}

func getGenericPersonAuthorizationUpdateCommand(
	personId,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32,
	update *authorizationUpdate) *Command {
	return &Command{
		InputAddresses:  []string{model.GetSettingsAddress(), signerId, personId},
		OutputAddresses: []string{personId},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signerId,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandUpdateAuthorization{
				CommandUpdateAuthorization: &model.CommandPersonUpdateAuthorization{
					PersonId:   personId,
					MakeMajor:  update.majorUpdate,
					MakeSigned: update.signedUpdate,
				},
			},
		},
	}
}

type authorizationUpdate struct {
	majorUpdate  model.BoolUpdate
	signedUpdate model.BoolUpdate
}

func GetPersonIncBalanceCommand(
	personId string,
	amount int32,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	return &Command{
		InputAddresses:  []string{model.GetSettingsAddress(), signerId, personId},
		OutputAddresses: []string{personId},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signerId,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandPersonUpdateBalanceIncrement{
				CommandPersonUpdateBalanceIncrement: &model.CommandPersonUpdateBalanceIncrement{
					PersonId:         personId,
					BalanceIncrement: amount,
				},
			},
		},
	}
}

func (nbce *nonBootstrapCommandExecution) checkPersonCreate(c *model.CommandPersonCreate) (*updater, error) {
	if err := checkSanityPersonCreate(c); err != nil {
		return nil, err
	}
	readData, err := nbce.blockchainAccess.GetState([]string{c.NewPersonId})
	if err != nil {
		return nil, err
	}
	err = nbce.unmarshalledState.add(readData, []string{c.NewPersonId})
	if err != nil {
		return nil, err
	}
	if nbce.unmarshalledState.getAddressState(c.NewPersonId) != ADDRESS_EMPTY {
		return nil, errors.New("Address collision when creating person: " + c.NewPersonId)
	}
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates: []singleUpdate{
			&singleUpdatePersonCreate{
				timestamp:    nbce.timestamp,
				personCreate: c,
				isSigned:     false,
				isMajor:      false,
			},
		},
	}, nil
}

func (nbce *nonBootstrapCommandExecution) checkPersonUpdateProperties(c *model.CommandPersonUpdateProperties) (
	*updater, error) {
	if c.PersonId != nbce.verifiedSignerId {
		return nil, errors.New("Person update properties not authorized. Properties can not be updated by someone else")
	}
	oldPerson := nbce.unmarshalledState.persons[nbce.verifiedSignerId]
	if err := checkModelCommandPersonUpdateProperties(c, oldPerson); err != nil {
		return nil, err
	}
	singleUpdates := createSingleUpdatesPersonUpdateProperties(c, oldPerson, nbce.timestamp)
	singleUpdates = append(singleUpdates, &singleUpdatePersonTimestamp{
		timestamp: nbce.timestamp,
		personId:  nbce.verifiedSignerId,
	})
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates:           singleUpdates,
	}, nil
}

type singleUpdatePersonCreate struct {
	timestamp    int64
	personCreate *model.CommandPersonCreate
	isMajor      bool
	isSigned     bool
}

var _ singleUpdate = new(singleUpdatePersonCreate)

func (u *singleUpdatePersonCreate) updateState(state *unmarshalledState) (writtenAddress string) {
	personId := u.personCreate.NewPersonId
	person := &model.StatePerson{
		Id:         personId,
		CreatedOn:  u.timestamp,
		ModifiedOn: u.timestamp,
		PublicKey:  u.personCreate.PublicKey,
		Name:       u.personCreate.Name,
		Email:      u.personCreate.Email,
		IsMajor:    u.isMajor,
		IsSigned:   u.isSigned,
	}
	state.persons[personId] = person
	return personId
}

func (u *singleUpdatePersonCreate) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	return ba.AddEvent(
		model.EV_PERSON_CREATE,
		[]processor.Attribute{
			{
				Key:   model.TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.TIMESTAMP,
				Value: fmt.Sprintf("%d", u.timestamp),
			},
			{
				Key:   model.EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.ID,
				Value: u.personCreate.NewPersonId,
			},
			{
				Key:   model.PERSON_NAME,
				Value: u.personCreate.Name,
			},
			{
				Key:   model.PERSON_PUBLIC_KEY,
				Value: u.personCreate.PublicKey,
			},
			{
				Key:   model.PERSON_EMAIL,
				Value: u.personCreate.Email,
			},
			{
				Key:   model.PERSON_IS_MAJOR,
				Value: strconv.FormatBool(u.isMajor),
			},
			{
				Key:   model.PERSON_IS_SIGNED,
				Value: strconv.FormatBool(u.isSigned),
			},
		},
		[]byte{})
}

type singleUpdatePersonPropertyUpdate struct {
	newValue   string
	stateField *string
	eventKey   string
	personId   string
	timestamp  int64
}

var _ singleUpdate = new(singleUpdatePersonPropertyUpdate)

func (su singleUpdatePersonPropertyUpdate) updateState(_ *unmarshalledState) (writtenAddress string) {
	*su.stateField = su.newValue
	return su.personId
}

func (su singleUpdatePersonPropertyUpdate) issueEvent(
	eventSeq int32, transactionId string, ba BlockchainAccess) error {
	return ba.AddEvent(
		model.EV_PERSON_UPDATE,
		[]processor.Attribute{
			{
				Key:   model.TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.TIMESTAMP,
				Value: fmt.Sprintf("%d", su.timestamp),
			},
			{
				Key:   model.EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.ID,
				Value: su.personId,
			},
			{
				Key:   su.eventKey,
				Value: su.newValue,
			},
		}, []byte{})
}

type singleUpdatePersonTimestamp struct {
	timestamp int64
	personId  string
}

var _ singleUpdate = new(singleUpdatePersonTimestamp)

func (u *singleUpdatePersonTimestamp) updateState(state *unmarshalledState) (writtenAddress string) {
	state.persons[u.personId].ModifiedOn = u.timestamp
	return u.personId
}

func (u *singleUpdatePersonTimestamp) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	return ba.AddEvent(
		model.EV_PERSON_MODIFICATION_TIME,
		[]processor.Attribute{
			{
				Key:   model.TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.TIMESTAMP,
				Value: fmt.Sprintf("%d", u.timestamp),
			},
			{
				Key:   model.EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.ID,
				Value: u.personId,
			},
		},
		[]byte{})
}
