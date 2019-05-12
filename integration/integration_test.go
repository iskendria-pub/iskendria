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
	logger = log.New(os.Stdout, "integration.TestBootstrap", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent)
	withLoggedInWithNewKey(doTestBootstrap, t)
}

func TestPersonCreate(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestPersonCreate", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent)
	withNewPersonCreate(doTestPersonCreate, t)
}

func TestPersonUpdateProperties(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestPersonUpdateProperties", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent)
	withNewPersonCreate(doTestPersonUpdate, t)
}

func doTestPersonUpdate(originalPersonCreate *command.PersonCreate, t *testing.T) {
	doTestPersonCreate(originalPersonCreate, t)
	if err := cliAlexandria.Login(personPublicKeyFile, personPrivateKeyFile); err != nil {
		t.Error("Could not login as newly created person")
	}
	newPublicKey := "Fake key"
	cmd, originalPersonId := getPersonUpdatePropertiesCommand(originalPersonCreate, newPublicKey, t)
	err := command.RunCommandForTest(cmd, "transactionPersonUpdate", blockchainAccess)
	if err != nil {
		t.Error("Could not run person update command: " + err.Error())
	}
	checkModifiedPerson(getPersonByKey(newPublicKey, t), originalPersonId, newPublicKey, t)
}

func getPersonUpdatePropertiesCommand(originalPersonCreate *command.PersonCreate, newPublicKey string, t *testing.T) (
	*command.Command, string) {
	originalPerson := getPersonByKey(originalPersonCreate.PublicKey, t)
	originalPersonUpdate := dao.PersonToPersonUpdate(originalPerson)
	newPersonUpdate := getNewDaoPersonUpdate(newPublicKey)
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

func getNewDaoPersonUpdate(newPublicKey string) *dao.PersonUpdate {
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

func TestPersonUpdateSetMajor(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestPersonUpdateSetMajor", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent)
	withNewPersonCreate(doTestPersonUpdateSetMajor, t)
}

func doTestPersonUpdateSetMajor(originalPersonCreate *command.PersonCreate, t *testing.T) {
	doTestPersonCreate(originalPersonCreate, t)
	cmd := command.GetPersonUpdateSetMajorCommand(
		getPersonByKey(originalPersonCreate.PublicKey, t).Id,
		getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
		cliAlexandria.LoggedIn(),
		getSettings(t).PriceMajorChangePersonAuthorization)
	err := command.RunCommandForTest(cmd, "transactionSetMajor", blockchainAccess)
	if err != nil {
		t.Error("Could not run person update set major command: " + err.Error())
	}
	updatedPerson := getPersonByKey(originalPersonCreate.PublicKey, t)
	if updatedPerson.IsMajor != true {
		t.Error("Person was not updated")
	}
	if updatedPerson.IsSigned != false {
		t.Error("Person should not have been signed")
	}
}

func TestPersonUpdateSetSigned(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestPersonUpdateSetSigned", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent)
	withNewPersonCreate(doTestPersonUpdateSetSigned, t)
}

func doTestPersonUpdateSetSigned(originalPersonCreate *command.PersonCreate, t *testing.T) {
	doTestPersonCreate(originalPersonCreate, t)
	cmd := command.GetPersonUpdateSetSignedCommand(
		getPersonByKey(originalPersonCreate.PublicKey, t).Id,
		getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
		cliAlexandria.LoggedIn(),
		getSettings(t).PriceMajorChangePersonAuthorization)
	err := command.RunCommandForTest(cmd, "transactionSetSigned", blockchainAccess)
	if err != nil {
		t.Error("Could not run person update set signed command: " + err.Error())
	}
	updatedPerson := getPersonByKey(originalPersonCreate.PublicKey, t)
	if updatedPerson.IsMajor != false {
		t.Error("Person should not have become major")
	}
	if updatedPerson.IsSigned != true {
		t.Error("Person was not signed")
	}
}

func TestPersonUpdateUnsetMajor(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestPersonUpdateUnsetMajor", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent)
	withLoggedInWithNewKey(doTestPersonUpdateUnsetMajor, t)
}

func doTestPersonUpdateUnsetMajor(t *testing.T) {
	doTestBootstrap(t)
	settings := getSettings(t)
	originalPerson := getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t)
	cmd := command.GetPersonUpdateUnsetMajorCommand(
		originalPerson.Id,
		originalPerson.Id,
		cliAlexandria.LoggedIn(),
		settings.PriceMajorChangePersonAuthorization)
	err := command.RunCommandForTest(cmd, "transactionIdUnsetMajor", blockchainAccess)
	if err != nil {
		t.Error("Could not run person unset major command: " + err.Error())
	}
	updated := getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t)
	if updated.IsMajor != false {
		t.Error("Majorship was not unset")
	}
	if updated.IsSigned != true {
		t.Error("Signed should not have been changed")
	}
}

func TestPersonUpdateUnsetSigned(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestPersonUpdateUnsetSigned", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent)
	withLoggedInWithNewKey(doTestPersonUpdateUnsetSigned, t)
}

func doTestPersonUpdateUnsetSigned(t *testing.T) {
	doTestBootstrap(t)
	settings := getSettings(t)
	originalPerson := getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t)
	cmd := command.GetPersonUpdateUnsetSignedCommand(
		originalPerson.Id,
		originalPerson.Id,
		cliAlexandria.LoggedIn(),
		settings.PriceMajorChangePersonAuthorization)
	err := command.RunCommandForTest(cmd, "transactionIdUnsetMajor", blockchainAccess)
	if err != nil {
		t.Error("Could not run person unset signed command: " + err.Error())
	}
	updated := getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t)
	if updated.IsMajor != true {
		t.Error("Majorship should not have been changed")
	}
	if updated.IsSigned != false {
		t.Error("Signed was not unset")
	}
}

func TestPersonUpdateIncBalance(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestPersonUpdateIncBalance", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent)
	withLoggedInWithNewKey(doTestPersonUpdateIncBalance, t)
}

func doTestPersonUpdateIncBalance(t *testing.T) {
	doTestBootstrap(t)
	theBalanceIncrement := 50
	originalPerson := getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t)
	cmd := command.GetPersonUpdateIncBalanceCommand(
		originalPerson.Id,
		int32(theBalanceIncrement),
		originalPerson.Id,
		cliAlexandria.LoggedIn(),
		int32(0))
	err := command.RunCommandForTest(cmd, "transactionIdIncBalance", blockchainAccess)
	if err != nil {
		t.Error("Could not run person update inc balance command: " + err.Error())
	}
	updated := getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t)
	if updated.Balance != int32(theBalanceIncrement) {
		t.Error("Balance has not been incremented")
	}
}
