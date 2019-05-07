package command

import (
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"gitlab.bbinfra.net/3estack/alexandria/model"
)

type Bootstrap struct {
	PriceMajorEditSettings               int32
	PriceMajorCreatePerson               int32
	PriceMajorChangePersonAuthorization  int32
	PriceMajorChangeJournalAuthorization int32
	PricePersonEdit                      int32
	PriceAuthorSubmitNewManuscript       int32
	PriceAuthorSubmitNewVersion          int32
	PriceAuthorAcceptAuthorship          int32
	PriceReviewerSubmit                  int32
	PriceEditorAllowManuscriptReview     int32
	PriceEditorRejectManuscript          int32
	PriceEditorPublishManuscript         int32
	PriceEditorAssignManuscript          int32
	PriceEditorCreateJournal             int32
	PriceEditorCreateVolume              int32
	PriceEditorEditJournal               int32
	PriceEditorAddColleague              int32
	PriceEditorAcceptDuty                int32
	Name                                 string
	Email                                string
}

func GetBootstrapCommand(bootstrap *Bootstrap, signerKey string) *Command {
	personId := model.CreatePersonAddress()
	return &Command{
		inputAddresses:  []string{model.GetSettingsAddress(), personId},
		outputAddresses: []string{model.GetSettingsAddress(), personId},
		command: &model.Command{
			Price:     int32(0),
			Signer:    personId,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_Bootstrap{
				Bootstrap: &model.CommandBootstrap{
					PriceList: &model.PriceList{
						PriceMajorEditSettings:               bootstrap.PriceMajorEditSettings,
						PriceMajorCreatePerson:               bootstrap.PriceMajorCreatePerson,
						PriceMajorChangePersonAuthorization:  bootstrap.PriceMajorChangePersonAuthorization,
						PriceMajorChangeJournalAuthorization: bootstrap.PriceMajorChangeJournalAuthorization,
						PricePersonEdit:                      bootstrap.PricePersonEdit,
						PriceAuthorSubmitNewManuscript:       bootstrap.PriceAuthorSubmitNewManuscript,
						PriceAuthorSubmitNewVersion:          bootstrap.PriceAuthorSubmitNewVersion,
						PriceAuthorAcceptAuthorship:          bootstrap.PriceAuthorAcceptAuthorship,
						PriceReviewerSubmit:                  bootstrap.PriceReviewerSubmit,
						PriceEditorAllowManuscriptReview:     bootstrap.PriceEditorAllowManuscriptReview,
						PriceEditorRejectManuscript:          bootstrap.PriceEditorRejectManuscript,
						PriceEditorPublishManuscript:         bootstrap.PriceEditorPublishManuscript,
						PriceEditorAssignManuscript:          bootstrap.PriceEditorAssignManuscript,
						PriceEditorCreateJournal:             bootstrap.PriceEditorCreateJournal,
						PriceEditorCreateVolume:              bootstrap.PriceEditorCreateVolume,
						PriceEditorEditJournal:               bootstrap.PriceEditorEditJournal,
						PriceEditorAddColleague:              bootstrap.PriceEditorAddColleague,
						PriceEditorAcceptDuty:                bootstrap.PriceEditorAcceptDuty,
					},
					FirstMajor: &model.CommandPersonCreate{
						NewPersonId: personId,
						PublicKey:   signerKey,
						Name:        bootstrap.Name,
						Email:       bootstrap.Email,
					},
				},
			},
		},
	}
}

func GetSettingsUpdateCommand(orig, updated *dao.Settings, signer string, price int32) *Command {
	return &Command{
		inputAddresses:  []string{model.GetSettingsAddress()},
		outputAddresses: []string{model.GetSettingsAddress()},
		command: &model.Command{
			Price:     price,
			Signer:    signer,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandSettingsUpdate{
				CommandSettingsUpdate: createModelCommandSettingsUpdate(orig, updated),
			},
		},
	}
}