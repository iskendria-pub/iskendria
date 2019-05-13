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
		OutputAddresses: []string{personId, signerId},
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

func GetPersonUpdatePropertiesCommand(
	personId string,
	orig,
	updated *dao.PersonUpdate,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	return &Command{
		InputAddresses:  []string{model.GetSettingsAddress(), personId, signerId},
		OutputAddresses: []string{personId, signerId},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signerId,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandPersonUpdateProperties{
				CommandPersonUpdateProperties: createModelCommandPersonUpdateProperties(personId, orig, updated),
			},
		},
	}
}

func GetPersonUpdateSetMajorCommand(
	personId,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	update := &authorizationUpdate{
		majorUpdate:  model.BoolUpdate_MAKE_TRUE,
		signedUpdate: model.BoolUpdate_UNMODIFIED,
	}
	return getGenericPersonUpdateAuthorizationCommand(personId, signerId, cryptoIdentity, price, update)
}

func GetPersonUpdateUnsetMajorCommand(
	personId,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	update := &authorizationUpdate{
		majorUpdate:  model.BoolUpdate_MAKE_FALSE,
		signedUpdate: model.BoolUpdate_UNMODIFIED,
	}
	return getGenericPersonUpdateAuthorizationCommand(personId, signerId, cryptoIdentity, price, update)
}

func GetPersonUpdateSetSignedCommand(
	personId,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	update := &authorizationUpdate{
		majorUpdate:  model.BoolUpdate_UNMODIFIED,
		signedUpdate: model.BoolUpdate_MAKE_TRUE,
	}
	return getGenericPersonUpdateAuthorizationCommand(personId, signerId, cryptoIdentity, price, update)
}

func GetPersonUpdateUnsetSignedCommand(
	personId,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	update := &authorizationUpdate{
		majorUpdate:  model.BoolUpdate_UNMODIFIED,
		signedUpdate: model.BoolUpdate_MAKE_FALSE,
	}
	return getGenericPersonUpdateAuthorizationCommand(personId, signerId, cryptoIdentity, price, update)
}

func getGenericPersonUpdateAuthorizationCommand(
	personId,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32,
	update *authorizationUpdate) *Command {
	return &Command{
		InputAddresses:  []string{model.GetSettingsAddress(), signerId, personId},
		OutputAddresses: []string{personId, signerId},
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

func GetPersonUpdateIncBalanceCommand(
	personId string,
	amount int32,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	return &Command{
		InputAddresses:  []string{model.GetSettingsAddress(), signerId, personId},
		OutputAddresses: []string{personId, signerId},
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
	expectedPrice := nbce.unmarshalledState.settings.PriceList.PriceMajorCreatePerson
	if nbce.price != expectedPrice {
		return nil, formatPriceError("PriceMajorCreatePerson", expectedPrice)
	}
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

func checkSanityPersonCreate(personCreate *model.CommandPersonCreate) error {
	if personCreate.NewPersonId == "" {
		return errors.New("personCreate.NewPersonId should be filled")
	}
	if !model.IsPersonAddress(personCreate.NewPersonId) {
		return errors.New("personCreate.NewPersonId should be a person address")
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
		model.EV_TYPE_PERSON_CREATE,
		[]processor.Attribute{
			{
				Key:   model.EV_KEY_TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.EV_KEY_TIMESTAMP,
				Value: fmt.Sprintf("%d", u.timestamp),
			},
			{
				Key:   model.EV_KEY_EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.EV_KEY_ID,
				Value: u.personCreate.NewPersonId,
			},
			{
				Key:   model.EV_KEY_PERSON_NAME,
				Value: u.personCreate.Name,
			},
			{
				Key:   model.EV_KEY_PERSON_PUBLIC_KEY,
				Value: u.personCreate.PublicKey,
			},
			{
				Key:   model.EV_KEY_PERSON_EMAIL,
				Value: u.personCreate.Email,
			},
			{
				Key:   model.EV_KEY_PERSON_IS_MAJOR,
				Value: strconv.FormatBool(u.isMajor),
			},
			{
				Key:   model.EV_KEY_PERSON_IS_SIGNED,
				Value: strconv.FormatBool(u.isSigned),
			},
		},
		[]byte{})
}

func (nbce *nonBootstrapCommandExecution) checkPersonUpdateProperties(c *model.CommandPersonUpdateProperties) (
	*updater, error) {
	expectedPrice := nbce.unmarshalledState.settings.PriceList.PricePersonEdit
	if nbce.price != expectedPrice {
		return nil, formatPriceError("PricePersonEdit", expectedPrice)
	}
	if c.PersonId != nbce.verifiedSignerId {
		return nil, errors.New("Person update properties not authorized. Properties can not be updated by someone else")
	}
	oldPerson := nbce.unmarshalledState.persons[nbce.verifiedSignerId]
	if err := checkModelCommandPersonUpdateProperties(c, oldPerson); err != nil {
		return nil, err
	}
	singleUpdates := createSingleUpdatesPersonUpdateProperties(c, oldPerson, nbce.timestamp)
	singleUpdates = nbce.addSingleUpdatePersonModificationTimeIfNeeded(singleUpdates, oldPerson.Id)
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates:           singleUpdates,
	}, nil
}

func (nbce *nonBootstrapCommandExecution) addSingleUpdatePersonModificationTimeIfNeeded(
	singleUpdates []singleUpdate, subjectId string) []singleUpdate {
	if len(singleUpdates) >= 1 {
		singleUpdates = append(singleUpdates, &singleUpdatePersonModificationTime{
			timestamp: nbce.timestamp,
			personId:  subjectId,
		})
	}
	return singleUpdates
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
		model.EV_TYPE_PERSON_UPDATE,
		[]processor.Attribute{
			{
				Key:   model.EV_KEY_TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.EV_KEY_TIMESTAMP,
				Value: fmt.Sprintf("%d", su.timestamp),
			},
			{
				Key:   model.EV_KEY_EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.EV_KEY_ID,
				Value: su.personId,
			},
			{
				Key:   su.eventKey,
				Value: su.newValue,
			},
		}, []byte{})
}

type singleUpdatePersonModificationTime struct {
	timestamp int64
	personId  string
}

var _ singleUpdate = new(singleUpdatePersonModificationTime)

func (u *singleUpdatePersonModificationTime) updateState(state *unmarshalledState) (writtenAddress string) {
	state.persons[u.personId].ModifiedOn = u.timestamp
	return u.personId
}

func (u *singleUpdatePersonModificationTime) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	return ba.AddEvent(
		model.EV_TYPE_PERSON_MODIFICATION_TIME,
		[]processor.Attribute{
			{
				Key:   model.EV_KEY_TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.EV_KEY_TIMESTAMP,
				Value: fmt.Sprintf("%d", u.timestamp),
			},
			{
				Key:   model.EV_KEY_EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.EV_KEY_ID,
				Value: u.personId,
			},
		},
		[]byte{})
}

func (nbce *nonBootstrapCommandExecution) checkPersonUpdateAuthorization(
	c *model.CommandPersonUpdateAuthorization) (*updater, error) {
	expectedPrice := nbce.unmarshalledState.settings.PriceList.PriceMajorChangePersonAuthorization
	if nbce.price != expectedPrice {
		return nil, formatPriceError("PriceMajorChangePersonAuthorization", expectedPrice)
	}
	if !nbce.unmarshalledState.persons[nbce.verifiedSignerId].IsMajor {
		return nil, errors.New("Only majors can do this")
	}
	data, err := nbce.blockchainAccess.GetState([]string{c.PersonId})
	if err != nil {
		return nil, errors.New("Could not read person address: " + c.PersonId)
	}
	err = nbce.unmarshalledState.add(data, []string{c.PersonId})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not unmarshall person %s, error %s",
			c.PersonId, err))
	}
	if nbce.unmarshalledState.getAddressState(c.PersonId) != ADDRESS_FILLED {
		return nil, errors.New(fmt.Sprintf("Person to modify does not exist: " + c.PersonId))
	}
	oldPerson := nbce.unmarshalledState.persons[c.PersonId]
	singleUpdates := createSingleUpdatesPersonUpdateAuthorization(c, oldPerson, nbce.timestamp)
	singleUpdates = nbce.addSingleUpdatePersonModificationTimeIfNeeded(singleUpdates, oldPerson.Id)
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates:           singleUpdates,
	}, nil
}

func createSingleUpdatesPersonUpdateAuthorization(
	c *model.CommandPersonUpdateAuthorization,
	oldPerson *model.StatePerson,
	timestamp int64) []singleUpdate {
	result := []singleUpdate{}
	switch c.MakeMajor {
	case model.BoolUpdate_UNMODIFIED:
		// Do nothing
	case model.BoolUpdate_MAKE_FALSE:
		toAppend := singleUpdatePersonAuthorizationUpdate{
			newValue:   false,
			stateField: &oldPerson.IsMajor,
			eventKey:   model.EV_KEY_PERSON_IS_MAJOR,
			personId:   oldPerson.Id,
			timestamp:  timestamp,
		}
		result = append(result, &toAppend)
	case model.BoolUpdate_MAKE_TRUE:
		toAppend := singleUpdatePersonAuthorizationUpdate{
			newValue:   true,
			stateField: &oldPerson.IsMajor,
			eventKey:   model.EV_KEY_PERSON_IS_MAJOR,
			personId:   oldPerson.Id,
			timestamp:  timestamp,
		}
		result = append(result, &toAppend)
	}
	switch c.MakeSigned {
	case model.BoolUpdate_UNMODIFIED:
		// Do nothing
	case model.BoolUpdate_MAKE_FALSE:
		toAppend := singleUpdatePersonAuthorizationUpdate{
			newValue:   false,
			stateField: &oldPerson.IsSigned,
			eventKey:   model.EV_KEY_PERSON_IS_SIGNED,
			personId:   oldPerson.Id,
			timestamp:  timestamp,
		}
		result = append(result, &toAppend)
	case model.BoolUpdate_MAKE_TRUE:
		toAppend := singleUpdatePersonAuthorizationUpdate{
			newValue:   true,
			stateField: &oldPerson.IsSigned,
			eventKey:   model.EV_KEY_PERSON_IS_SIGNED,
			personId:   oldPerson.Id,
			timestamp:  timestamp,
		}
		result = append(result, &toAppend)
	}
	return result
}

type singleUpdatePersonAuthorizationUpdate struct {
	newValue   bool
	stateField *bool
	eventKey   string
	personId   string
	timestamp  int64
}

var _ singleUpdate = new(singleUpdatePersonAuthorizationUpdate)

func (u *singleUpdatePersonAuthorizationUpdate) updateState(*unmarshalledState) (writtenAddress string) {
	*u.stateField = u.newValue
	return u.personId
}

func (u *singleUpdatePersonAuthorizationUpdate) issueEvent(
	eventSeq int32, transactionId string, ba BlockchainAccess) error {
	return ba.AddEvent(
		model.EV_TYPE_PERSON_UPDATE,
		[]processor.Attribute{
			{
				Key:   model.EV_KEY_TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.EV_KEY_TIMESTAMP,
				Value: fmt.Sprintf("%d", u.timestamp),
			},
			{
				Key:   model.EV_KEY_EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.EV_KEY_ID,
				Value: u.personId,
			},
			{
				Key:   u.eventKey,
				Value: strconv.FormatBool(u.newValue),
			},
		}, []byte{})
}

func (nbce *nonBootstrapCommandExecution) checkPersonUpdateIncBalance(c *model.CommandPersonUpdateBalanceIncrement) (
	*updater, error) {
	if nbce.price != int32(0) {
		return nil, formatPriceError("-", int32(0))
	}
	if !nbce.unmarshalledState.persons[nbce.verifiedSignerId].IsMajor {
		return nil, errors.New("Only majors can increment someone's balance")
	}
	data, err := nbce.blockchainAccess.GetState([]string{c.PersonId})
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Could not read person %s: %s", c.PersonId, err.Error()))
	}
	err = nbce.unmarshalledState.add(data, []string{c.PersonId})
	if err != nil {
		return nil, errors.New("Could not unmarshall person: " + c.PersonId)
	}
	if nbce.unmarshalledState.getAddressState(c.PersonId) != ADDRESS_FILLED {
		return nil, errors.New("The person whose balance is to be updated does not exist: " + c.PersonId)
	}
	oldPerson := nbce.unmarshalledState.persons[c.PersonId]
	singleUpdates := []singleUpdate{
		&singleUpdatePersonIncBalance{
			personId:   c.PersonId,
			newBalance: oldPerson.Balance + c.BalanceIncrement,
			timestamp:  nbce.timestamp,
		},
	}
	singleUpdates = nbce.addSingleUpdatePersonModificationTimeIfNeeded(singleUpdates, oldPerson.Id)
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates:           singleUpdates,
	}, nil
}

type singleUpdatePersonIncBalance struct {
	personId   string
	newBalance int32
	timestamp  int64
}

var _ singleUpdate = new(singleUpdatePersonIncBalance)

func (u *singleUpdatePersonIncBalance) updateState(state *unmarshalledState) (writtenAddress string) {
	state.persons[u.personId].Balance = u.newBalance
	return u.personId
}

func (u *singleUpdatePersonIncBalance) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	return ba.AddEvent(
		model.EV_TYPE_PERSON_UPDATE,
		[]processor.Attribute{
			{
				Key:   model.EV_KEY_TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.EV_KEY_EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.EV_KEY_TIMESTAMP,
				Value: fmt.Sprintf("%d", u.timestamp),
			},
			{
				Key:   model.EV_KEY_ID,
				Value: u.personId,
			},
			{
				Key:   model.EV_KEY_PERSON_BALANCE,
				Value: fmt.Sprintf("%d", u.newBalance),
			},
		},
		[]byte{})
}
