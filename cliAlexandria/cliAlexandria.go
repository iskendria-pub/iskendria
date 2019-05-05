package cliAlexandria

import (
	"errors"
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/blockchain"
	"gitlab.bbinfra.net/3estack/alexandria/cli"
	"gitlab.bbinfra.net/3estack/alexandria/command"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
)

var LoggedInPerson *dao.Person
var Settings *dao.Settings

var CommonHandlers = []cli.Handler{
	&cli.SingleLineHandler{
		Name:     "login",
		Handler:  login,
		ArgNames: []string{"public key file", "private key file"},
	},
	&cli.SingleLineHandler{
		Name:     "logout",
		Handler:  logout,
		ArgNames: []string{},
	},
	&cli.SingleLineHandler{
		Name:     "createKeys",
		Handler:  createKeyPair,
		ArgNames: []string{"public key file", "private key file"},
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
	isDuplicate, err := updateLoggedInPerson(LoggedInPublicKeyStr)
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
		outputter("Pleas login first\n")
		return false
	}
	return true
}

func ToIoError(err error) string {
	return "IO error while sending command: " + err.Error() + "\n"
}

func whoAmI(outputter cli.Outputter) {
	if !CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	outputter(cli.StructToTable(LoggedInPerson).String())
}

func whoIs(outputter cli.Outputter, personId string) {
	person, err := dao.GetPersonById(personId)
	if err != nil {
		outputter(personNotFound(personId) + ", error: " + err.Error() + "\n")
	}
	if person == nil {
		outputter(personNotFound(personId) + "\n")
	}
	outputter(cli.StructToTable(person).String())
}

func personNotFound(personId string) string {
	return fmt.Sprintf("Person not found: %s", personId)
}
