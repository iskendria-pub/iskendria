package cliAlexandria

import (
	"errors"
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/blockchain"
	"gitlab.bbinfra.net/3estack/alexandria/cli"
	"gitlab.bbinfra.net/3estack/alexandria/command"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"time"
)

var LoggedInPerson *dao.Person
var Settings *dao.Settings

var CommonHandlers = []cli.Handler{
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
	&cli.SingleLineHandler{
		Name:     "showSettings",
		Handler:  showSettings,
		ArgNames: []string{},
	},
	&cli.SingleLineHandler{
		Name:     "whoAmI",
		Handler:  whoAmI,
		ArgNames: []string{},
	},
	&cli.SingleLineHandler{
		Name:     "whoIs",
		Handler:  whoIs,
		ArgNames: []string{"person id"},
	},
}

func SendCommandAsPerson(outputter cli.Outputter, commandFactory func() *command.Command) {
	if !CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	if err := blockchain.SendCommand(commandFactory()); err != nil {
		outputter(ToIoError(err))
	}
}

func CheckBootstrappedAndKnownPerson(outputter cli.Outputter) bool {
	if !CheckBootstrappedAndLoggedIn(outputter) {
		return false
	}
	isDuplicate, err := updateLoggedInPerson(loggedIn.PublicKeyStr)
	if err != nil {
		outputter("ERROR: The key you logged in with has no associated person\n")
		return false
	}
	if isDuplicate {
		outputter("WARNING: The key you logged in with has multiple persons, chose arbitrarily id: " +
			LoggedInPerson.Id + "\n")
	}
	return true
}

func updateLoggedInPerson(loggedInPublicKey string) (isDuplicate bool, err error) {
	var persons []*dao.Person
	persons, err = dao.SearchPersonByKey(loggedInPublicKey)
	if err != nil {
		err = errors.New(noPersonFoundForKey(loggedInPublicKey) + ": " + err.Error())
		return
	}
	if len(persons) == 0 {
		err = errors.New(noPersonFoundForKey(loggedInPublicKey))
		return
	}
	if len(persons) >= 2 {
		isDuplicate = true
	}
	LoggedInPerson = persons[0]
	return
}

func noPersonFoundForKey(key string) string {
	return fmt.Sprintf("No person found for key %s", key)
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

func ToIoError(err error) string {
	return "IO error while sending command: " + err.Error() + "\n"
}

func showSettings(outputter cli.Outputter) {
	daoSettings, err := dao.GetSettings()
	if err != nil {
		outputter(err.Error() + "\n")
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

func formatTime(t int64) string {
	return time.Unix(t, 0).Format(time.UnixDate)
}

func whoAmI(outputter cli.Outputter) {
	if !CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	person := daoPersonToPersonView(LoggedInPerson)
	outputter(cli.StructToTable(person).String())
}

func daoPersonToPersonView(daoPerson *dao.Person) *PersonView {
	result := new(PersonView)
	result.Id = daoPerson.Id
	result.CreatedOn = formatTime(daoPerson.CreatedOn)
	result.ModifiedOn = formatTime(daoPerson.ModifiedOn)
	result.PublicKey = daoPerson.PublicKey
	result.Name = daoPerson.Name
	result.Email = daoPerson.Email
	result.IsMajor = daoPerson.IsMajor
	result.IsSigned = daoPerson.IsSigned
	result.Balance = daoPerson.Balance
	result.BiographyHash = daoPerson.BiographyHash
	result.Organization = daoPerson.Organization
	result.Telephone = daoPerson.Telephone
	result.Address = daoPerson.Address
	result.PostalCode = daoPerson.PostalCode
	result.Country = daoPerson.Country
	result.ExtraInfo = daoPerson.ExtraInfo
	return result
}

// This type represents a person as it has to be shown to
// end users. It is like dao.Person but the creation time
// and the modification time are formatted as strings.
type PersonView struct {
	Id            string
	CreatedOn     string
	ModifiedOn    string
	PublicKey     string
	Name          string
	Email         string
	IsMajor       bool
	IsSigned      bool
	Balance       int32
	BiographyHash string
	Organization  string
	Telephone     string
	Address       string
	PostalCode    string
	Country       string
	ExtraInfo     string
}

func whoIs(outputter cli.Outputter, personId string) {
	daoPerson, err := dao.GetPersonById(personId)
	if err != nil {
		outputter(personNotFound(personId) + ", error: " + err.Error() + "\n")
	}
	if daoPerson == nil {
		outputter(personNotFound(personId) + "\n")
	}
	person := daoPersonToPersonView(daoPerson)
	outputter(cli.StructToTable(person).String())
}

func personNotFound(personId string) string {
	return fmt.Sprintf("Person not found: %s", personId)
}
