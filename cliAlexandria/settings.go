package cliAlexandria

import (
	"github.com/iskendria-pub/iskendria/cli"
	"github.com/iskendria-pub/iskendria/dao"
)

var CommonRootHandlers = []cli.Handler{
	&cli.SingleLineHandler{
		Name:     "login",
		Handler:  Login,
		ArgNames: []string{"public key file", "private key file"},
	},
	&cli.SingleLineHandler{
		Name:     "logout",
		Handler:  Logout,
		ArgNames: []string{},
	},
	&cli.SingleLineHandler{
		Name:     "createKeys",
		Handler:  CreateKeyPair,
		ArgNames: []string{"public key file", "private key file"},
	},
}

var CommonSettingsHandlers = []cli.Handler{
	&cli.SingleLineHandler{
		Name:     "showSettings",
		Handler:  showSettings,
		ArgNames: []string{},
	},
}

func CheckBootstrappedAndLoggedIn(outputter cli.Outputter) bool {
	var err error
	Settings, err = dao.GetSettings()
	if err != nil {
		outputter(err.Error() + "\n")
		return false
	}
	if Settings == nil {
		outputter("The Blockchain has not been bootstrapped yet, please do that first\n")
		return false
	}
	if !IsLoggedIn() {
		outputter("Pleas Login first\n")
		return false
	}
	return true
}

var Settings *dao.Settings

func showSettings(outputter cli.Outputter) {
	daoSettings, err := dao.GetSettings()
	if err != nil {
		outputter(err.Error() + "\n")
		return
	}
	if daoSettings == nil {
		outputter("Not bootstrapped\n")
		return
	}
	settings := daoSettingsToSettingsView(daoSettings)
	outputter(cli.StructToTable(settings).String())
}

func daoSettingsToSettingsView(settings *dao.Settings) *SettingsView {
	result := new(SettingsView)
	result.CreatedOn = formatTime(settings.CreatedOn)
	result.ModifiedOn = formatTime(settings.ModifiedOn)
	result.PriceMajorEditSettings = settings.PriceMajorEditSettings
	result.PriceMajorCreatePerson = settings.PriceMajorCreatePerson
	result.PriceMajorChangePersonAuthorization = settings.PriceMajorChangePersonAuthorization
	result.PriceMajorChangeJournalAuthorization = settings.PriceMajorChangeJournalAuthorization
	result.PricePersonEdit = settings.PricePersonEdit
	result.PriceAuthorSubmitNewManuscript = settings.PriceAuthorSubmitNewManuscript
	result.PriceAuthorSubmitNewVersion = settings.PriceAuthorSubmitNewVersion
	result.PriceAuthorAcceptAuthorship = settings.PriceAuthorAcceptAuthorship
	result.PriceReviewerSubmit = settings.PriceReviewerSubmit
	result.PriceEditorAllowManuscriptReview = settings.PriceEditorAllowManuscriptReview
	result.PriceEditorRejectManuscript = settings.PriceEditorRejectManuscript
	result.PriceEditorPublishManuscript = settings.PriceEditorPublishManuscript
	result.PriceEditorAssignManuscript = settings.PriceEditorAssignManuscript
	result.PriceEditorCreateJournal = settings.PriceEditorCreateJournal
	result.PriceEditorCreateVolume = settings.PriceEditorCreateVolume
	result.PriceEditorEditJournal = settings.PriceEditorEditJournal
	result.PriceEditorAddColleague = settings.PriceEditorAddColleague
	result.PriceEditorAcceptDuty = settings.PriceEditorAcceptDuty
	return result
}

type SettingsView struct {
	CreatedOn                            string
	ModifiedOn                           string
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
}
