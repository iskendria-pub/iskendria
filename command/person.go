package command

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"gitlab.bbinfra.net/3estack/alexandria/model"
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
		InputAddresses:  []string{personId, signerId},
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
				CommandPersonUpdateProperties: createModelCommandPersonUpdate(orig, updated),
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

type singleUpdatePersonCreate struct {
	timestamp    int64
	personCreate *model.CommandPersonCreate
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
				Value: "xxx@gmail.com",
			},
		},
		[]byte{})
}
