package command

import (
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
