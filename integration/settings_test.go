package integration

import (
	"gitlab.bbinfra.net/3estack/alexandria/cliAlexandria"
	"gitlab.bbinfra.net/3estack/alexandria/command"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"log"
	"os"
	"testing"
)

func TestBootstrap(t *testing.T) {
	logger := log.New(os.Stdout, "integration.TestBootstrap", log.Flags())
	blockchainAccess := command.NewBlockchainStub(dao.HandleEvent)
	withLoggedInWithNewKey(doTestBootstrap, blockchainAccess, logger, t)
}

func TestPersonCreate(t *testing.T) {
	logger := log.New(os.Stdout, "integration.TestBootstrap", log.Flags())
	blockchainAccess := command.NewBlockchainStub(dao.HandleEvent)
	withNewPersonCreate(doTestPersonCreate, blockchainAccess, logger, t)
}

func TestPersonUpdate(t *testing.T) {
	logger := log.New(os.Stdout, "integration.TestBootstrap", log.Flags())
	blockchainAccess := command.NewBlockchainStub(dao.HandleEvent)
	withNewPersonCreate(doTestPersonUpdate, blockchainAccess, logger, t)
}

func doTestPersonUpdate(
	originalPersonCreate *command.PersonCreate,
	blockchainAccess command.BlockchainAccess,
	logger *log.Logger,
	t *testing.T) {
	doTestPersonCreate(originalPersonCreate, blockchainAccess, logger, t)
	if err := cliAlexandria.Login(personPublicKeyFile, personPrivateKeyFile); err != nil {
		t.Error("Could not login as newly created person")
	}
	newPublicKey := "Fake key"
	cmd, originalPersonId := getPersonUpdateCommand(originalPersonCreate, newPublicKey, t)
	err := command.RunCommandForTest(cmd, "transactionPersonUpdate", blockchainAccess)
	if err != nil {
		t.Error("Could not run person update command: " + err.Error())
	}
	checkModifiedPerson(getPersonByKey(newPublicKey, t), originalPersonId, newPublicKey, t)
}

func getPersonUpdateCommand(
	originalPersonCreate *command.PersonCreate,
	newPublicKey string,
	t *testing.T) (*command.Command, string) {
	originalPerson := getPersonByKey(originalPersonCreate.PublicKey, t)
	originalPersonUpdate := dao.PersonToPersonUpdate(originalPerson)
	newPersonUpdate := getNewPersonUpdate(newPublicKey)
	settings, err := dao.GetSettings()
	if err != nil {
		t.Error(err)
	}
	cmd := command.GetPersonUpdatePropertiesCommand(
		originalPerson.Id,
		originalPersonUpdate,
		newPersonUpdate,
		getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
		cliAlexandria.LoggedIn(),
		settings.PricePersonEdit)
	return cmd, originalPerson.Id
}

func getNewPersonUpdate(newPublicKey string) *dao.PersonUpdate {
	result := new(dao.PersonUpdate)
	result.PublicKey = newPublicKey
	result.Name = "Peter"
	result.Email = "peter@gmail.com"
	result.BiographyHash = "01234567"
	result.Organization = "Peter's Toko"
	result.Telephone = "088-3456789"
	result.Address = "Insulindeweg"
	result.PostalCode = "1234 AB"
	result.Country = "Australia"
	result.ExtraInfo = "Some fake data"
	return result
}

func checkModifiedPerson(person *dao.Person, expectedId, expectedPublicKey string, t *testing.T) {
	if person.Id != expectedId {
		t.Error("Id mismatch")
	}
	// Please check createdOn and modifiedOn manually, difficult to test
	if person.PublicKey != expectedPublicKey {
		t.Error("PublicKey mismatch")
	}
	if person.Name != "Peter" {
		t.Error("Name mismatch")
	}
	if person.Email != "peter@gmail.com" {
		t.Error("Email mismatch")
	}
	if person.IsMajor != false {
		t.Error("IsMajor mismatch")
	}
	if person.IsSigned != false {
		t.Error("IsSigned mismatch")
	}
	if person.Balance != int32(0) {
		t.Error("Balance mismatch")
	}
	if person.BiographyHash != "01234567" {
		t.Error("BiographyHash mismatch")
	}
	if person.Organization != "Peter's Toko" {
		t.Error("Organization mismatch")
	}
	if person.Telephone != "088-3456789" {
		t.Error("Telephone mismatch")
	}
	if person.Address != "Insulindeweg" {
		t.Error("Address mismatch")
	}
	if person.PostalCode != "1234 AB" {
		t.Error("PostalCode mismatch")
	}
	if person.Country != "Australia" {
		t.Error("Country mismatch")
	}
	if person.ExtraInfo != "Some fake data" {
		t.Error("ExtraInfo mismatch")
	}
}
