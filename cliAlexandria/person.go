package cliAlexandria

import (
	"errors"
	"fmt"
	"github.com/iskendria-pub/iskendria/blockchain"
	"github.com/iskendria-pub/iskendria/cli"
	"github.com/iskendria-pub/iskendria/command"
	"github.com/iskendria-pub/iskendria/dao"
)

var CommonPersonHandlers = []cli.Handler{
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

func whoAmI(outputter cli.Outputter) {
	if !CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	person := daoPersonToPersonView(LoggedInPerson)
	outputter(cli.StructToTable(person).String())
}

func CheckBootstrappedAndKnownPerson(outputter cli.Outputter) bool {
	if !CheckBootstrappedAndLoggedIn(outputter) {
		return false
	}
	hasMultiplePersons, err := updateLoggedInPerson(loggedIn.PublicKeyStr)
	if err != nil {
		outputter("ERROR: The key you logged in with has no associated person\n")
		return false
	}
	if hasMultiplePersons {
		outputter("WARNING: The key you logged in with has multiple persons, chose arbitrarily id: " +
			LoggedInPerson.Id + "\n")
	}
	return true
}

var LoggedInPerson *dao.Person

func updateLoggedInPerson(loggedInPublicKey string) (hasMultiplePersons bool, err error) {
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
		hasMultiplePersons = true
	}
	LoggedInPerson = persons[0]
	return
}

func noPersonFoundForKey(key string) string {
	return fmt.Sprintf("No person found for key %s", key)
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
		return
	}
	if daoPerson == nil {
		outputter(personNotFound(personId) + "\n")
		return
	}
	person := daoPersonToPersonView(daoPerson)
	outputter(cli.StructToTable(person).String())
}

func personNotFound(personId string) string {
	return fmt.Sprintf("Person not found: %s", personId)
}

func SendCommandAsPerson(outputter cli.Outputter, commandFactory func() *command.Command) {
	if !CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	if err := blockchain.SendCommand(commandFactory(), outputter); err != nil {
		outputter(ToIoError(err))
	}
}

var OriginalPersonId string
var OriginalPerson *dao.PersonUpdate

func PersonUpdate(outputter cli.Outputter, newPerson *dao.PersonUpdate) {
	theCommand := command.GetPersonUpdatePropertiesCommand(
		OriginalPersonId,
		OriginalPerson,
		newPerson,
		LoggedInPerson.Id,
		LoggedIn(),
		Settings.PricePersonEdit)
	if err := blockchain.SendCommand(theCommand, outputter); err != nil {
		outputter(ToIoError(err))
	}
}
