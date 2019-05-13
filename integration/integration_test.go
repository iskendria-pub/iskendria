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

func TestSettingsUpdate(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestSettingsUpdate", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent)
	withLoggedInWithNewKey(doTestSettingsUpdate, t)
}

func doTestSettingsUpdate(t *testing.T) {
	doTestBootstrap(t)
	signer := getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t)
	origSettings := getSettings(t)
	var settingsUpdate = getSettingsUpdate()
	cmd := command.GetSettingsUpdateCommand(
		origSettings,
		settingsUpdate,
		signer.Id,
		cliAlexandria.LoggedIn(),
		origSettings.PriceMajorEditSettings)
	err := command.RunCommandForTest(cmd, "transactionIdSettingsUpdate", blockchainAccess)
	if err != nil {
		t.Error("Could not run settings update command: " + err.Error())
	}
	updated := getSettings(t)
	checkUpdatedSettings(updated, t)
}

func getSettingsUpdate() *dao.Settings {
	return &dao.Settings{
		PriceMajorEditSettings:               201,
		PriceMajorCreatePerson:               202,
		PriceMajorChangePersonAuthorization:  203,
		PriceMajorChangeJournalAuthorization: 204,
		PricePersonEdit:                      205,
		PriceAuthorSubmitNewManuscript:       206,
		PriceAuthorSubmitNewVersion:          207,
		PriceAuthorAcceptAuthorship:          208,
		PriceReviewerSubmit:                  209,
		PriceEditorAllowManuscriptReview:     210,
		PriceEditorRejectManuscript:          211,
		PriceEditorPublishManuscript:         212,
		PriceEditorAssignManuscript:          213,
		PriceEditorCreateJournal:             214,
		PriceEditorCreateVolume:              215,
		PriceEditorEditJournal:               216,
		PriceEditorAddColleague:              217,
		PriceEditorAcceptDuty:                218,
	}
}

func checkUpdatedSettings(updated *dao.Settings, t *testing.T) {
	if updated.PriceMajorEditSettings != int32(201) {
		t.Error("PriceMajorEditSettings mismatch")
	}
	if updated.PriceMajorCreatePerson != int32(202) {
		t.Error("PriceMajorCreatePerson mismatch")
	}
	if updated.PriceMajorChangePersonAuthorization != int32(203) {
		t.Error("PriceMajorChangePersonAuthorization mismatch")
	}
	if updated.PriceMajorChangeJournalAuthorization != int32(204) {
		t.Error("PriceMajorChangeJournalAuthorization mismatch")
	}
	if updated.PricePersonEdit != int32(205) {
		t.Error("PricePersonEdit mismatch")
	}
	if updated.PriceAuthorSubmitNewManuscript != int32(206) {
		t.Error("PriceAuthorSubmitNewManuscript mismatch")
	}
	if updated.PriceAuthorSubmitNewVersion != int32(207) {
		t.Error("PriceAuthorSubmitNewVersion mismatch")
	}
	if updated.PriceAuthorAcceptAuthorship != int32(208) {
		t.Error("PriceAuthorAcceptAuthorship mismatch")
	}
	if updated.PriceReviewerSubmit != int32(209) {
		t.Error("PriceReviewerSubmit mismatch")
	}
	if updated.PriceEditorAllowManuscriptReview != int32(210) {
		t.Error("PriceEditorAllowManuscriptReview mismatch")
	}
	if updated.PriceEditorRejectManuscript != int32(211) {
		t.Error("PriceEditorRejectManuscript mismatch")
	}
	if updated.PriceEditorPublishManuscript != int32(212) {
		t.Error("PriceEditorPublishManuscript mismatch")
	}
	if updated.PriceEditorAssignManuscript != int32(213) {
		t.Error("PriceEditorAssignManuscript mismatch")
	}
	if updated.PriceEditorCreateJournal != int32(214) {
		t.Error("PriceEditorCreateJournal mismatch")
	}
	if updated.PriceEditorCreateVolume != int32(215) {
		t.Error("PriceEditorCreateVolume mismatch")
	}
	if updated.PriceEditorEditJournal != int32(216) {
		t.Error("PriceEditorEditJournal mismatch")
	}
	if updated.PriceEditorAddColleague != int32(217) {
		t.Error("PriceEditorAddColleague mismatch")
	}
	if updated.PriceEditorAcceptDuty != int32(218) {
		t.Error("PriceEditorAcceptDuty mismatch")
	}
}
