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

func GetPersonCreateCommand(pc *PersonCreate, signer string, price int32) *Command {
	personId := model.CreatePersonAddress()
	return &Command{
		inputAddresses:  []string{personId, signer},
		outputAddresses: []string{personId},
		command: &model.Command{
			Signer:    signer,
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

func GetPersonUpdateCommand(personId string, orig, updated *dao.PersonUpdate, signer string, price int32) *Command {
	return &Command{
		inputAddresses:  []string{model.GetSettingsAddress(), personId, signer},
		outputAddresses: []string{personId},
		command: &model.Command{
			Signer:    signer,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandPersonUpdateProperties{
				CommandPersonUpdateProperties: createModelCommandPersonUpdate(orig, updated),
			},
		},
	}
}

func GetPersonSetMajorCommand(personId, signer string, price int32) *Command {
	update := &authorizationUpdate{
		majorUpdate:  model.BoolUpdate_MAKE_TRUE,
		signedUpdate: model.BoolUpdate_UNMODIFIED,
	}
	return getGenericPersonAuthorizationUpdateCommand(personId, signer, price, update)
}

func GetPersonUnsetMajorCommand(personId, signer string, price int32) *Command {
	update := &authorizationUpdate{
		majorUpdate:  model.BoolUpdate_MAKE_FALSE,
		signedUpdate: model.BoolUpdate_UNMODIFIED,
	}
	return getGenericPersonAuthorizationUpdateCommand(personId, signer, price, update)
}

func GetPersonSetSignedCommand(personId, signer string, price int32) *Command {
	update := &authorizationUpdate{
		majorUpdate:  model.BoolUpdate_UNMODIFIED,
		signedUpdate: model.BoolUpdate_MAKE_TRUE,
	}
	return getGenericPersonAuthorizationUpdateCommand(personId, signer, price, update)
}

func GetPersonUnsetSignedCommand(personId, signer string, price int32) *Command {
	update := &authorizationUpdate{
		majorUpdate:  model.BoolUpdate_UNMODIFIED,
		signedUpdate: model.BoolUpdate_MAKE_FALSE,
	}
	return getGenericPersonAuthorizationUpdateCommand(personId, signer, price, update)
}

func getGenericPersonAuthorizationUpdateCommand(personId, signer string, price int32, update *authorizationUpdate) *Command {
	return &Command{
		inputAddresses:  []string{model.GetSettingsAddress(), signer, personId},
		outputAddresses: []string{personId},
		command: &model.Command{
			Signer:    signer,
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

func GetPersonIncBalanceCommand(personId string, amount int32, signer string, price int32) *Command {
	return &Command{
		inputAddresses:  []string{model.GetSettingsAddress(), signer, personId},
		outputAddresses: []string{personId},
		command: &model.Command{
			Signer:    signer,
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
