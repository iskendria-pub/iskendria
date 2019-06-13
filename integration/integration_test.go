package integration

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"gitlab.bbinfra.net/3estack/alexandria/cliAlexandria"
	"gitlab.bbinfra.net/3estack/alexandria/command"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"gitlab.bbinfra.net/3estack/alexandria/util"
	"log"
	"os"
	"testing"
)

func TestBootstrap(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestBootstrap", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	withLoggedInWithNewKey(doTestBootstrap, t)
}

func TestPersonCreate(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestPersonCreate", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(personCreate *command.PersonCreate, t *testing.T) {
		doTestPersonCreate(personCreate, t)
		checkStateBalanceOfKey(
			SUFFICIENT_BALANCE-priceMajorCreatePerson,
			cliAlexandria.LoggedIn().PublicKeyStr,
			t)
		checkDaoBalanceOfKey(
			SUFFICIENT_BALANCE-priceMajorCreatePerson,
			cliAlexandria.LoggedIn().PublicKeyStr,
			t)
	}
	withNewPersonCreate(f, t)
}

func TestPersonUpdatePropertiesAsSelf(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestPersonUpdatePropertiesAsSelf", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(originalPersonCreate *command.PersonCreate, t *testing.T) {
		doTestPersonCreate(originalPersonCreate, t)
		if err := cliAlexandria.Login(personPublicKeyFile, personPrivateKeyFile); err != nil {
			t.Error("Could not login as newly created person")
		}
		doTestPersonUpdate(originalPersonCreate, t)
		checkStateBalanceOfKey(
			SUFFICIENT_BALANCE-pricePersonEdit,
			"Fake key",
			t)
		checkDaoBalanceOfKey(
			SUFFICIENT_BALANCE-pricePersonEdit,
			"Fake key",
			t)
	}
	withNewPersonCreate(f, t)
}

func doTestPersonUpdate(originalPersonCreate *command.PersonCreate, t *testing.T) {
	newPublicKey := "Fake key"
	cmd, originalPersonId := getPersonUpdatePropertiesCommand(originalPersonCreate, newPublicKey, t)
	err := command.RunCommandForTest(cmd, "transactionPersonUpdate", blockchainAccess)
	if err != nil {
		t.Error("Could not run person update command: " + err.Error())
	}
	checkModifiedStatePerson(getStatePerson(originalPersonId, t), originalPersonId, newPublicKey, t)
	checkModifiedDaoPerson(getPersonByKey(newPublicKey, t), originalPersonId, newPublicKey, t)
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
	result.Organization = "Peter's Toko"
	result.Telephone = "088-3456789"
	result.Address = "Insulindeweg"
	result.PostalCode = "1234 AB"
	result.Country = "Australia"
	result.ExtraInfo = "Some fake data"
	return result
}

func checkModifiedStatePerson(person *model.StatePerson, expectedPersonId, expectedPublicKey string, t *testing.T) {
	if person.Id != expectedPersonId {
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

func checkModifiedDaoPerson(person *dao.Person, expectedId, expectedPublicKey string, t *testing.T) {
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

func TestPersonUpdatePropertiesAsSMajor(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestPersonUpdatePropertiesAsMajor", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(originalPersonCreate *command.PersonCreate, t *testing.T) {
		doTestPersonCreate(originalPersonCreate, t)
		doTestPersonUpdate(originalPersonCreate, t)
		checkStateBalanceOfKey(
			SUFFICIENT_BALANCE-priceMajorCreatePerson-pricePersonEdit,
			cliAlexandria.LoggedIn().PublicKeyStr,
			t)
		checkDaoBalanceOfKey(
			SUFFICIENT_BALANCE-priceMajorCreatePerson-pricePersonEdit,
			cliAlexandria.LoggedIn().PublicKeyStr,
			t)
	}
	withNewPersonCreate(f, t)
}

func TestPersonBiography(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestPersonUpdatePropertiesAsMajor", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(t *testing.T) {
		doTestBootstrap(t)
		personId := getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id
		biography := []byte("This is my biography")
		cmd := command.GetCommandPersonUpdateBiography(
			personId,
			"",
			biography,
			personId,
			cliAlexandria.LoggedIn(),
			pricePersonEdit)
		err := command.RunCommandForTest(cmd, "transactionIdPersonSetBiography", blockchainAccess)
		if err != nil {
			t.Error(err)
		}
		err = dao.VerifyPersonBiography(personId, biography)
		if err != nil {
			t.Error(err)
		}
		cmd = command.GetCommandPersonOmitBiography(
			personId,
			model.HashBytes(biography),
			personId,
			cliAlexandria.LoggedIn(),
			pricePersonEdit)
		err = command.RunCommandForTest(cmd, "transactionIdPersonRemoveBiography", blockchainAccess)
		if err != nil {
			t.Error(err)
		}
		err = dao.VerifyPersonBiography(personId, []byte{})
		if err != nil {
			t.Error(err)
		}
	}
	withLoggedInWithNewKey(f, t)
}

func TestPersonUpdateSetMajor(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestPersonUpdateSetMajor", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(originalPersonCreate *command.PersonCreate, t *testing.T) {
		doTestPersonUpdateSetMajor(originalPersonCreate, t)
		expectedBalance := SUFFICIENT_BALANCE - priceMajorCreatePerson - priceMajorChangePersonAuthorization
		checkStateBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
		checkDaoBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	}
	withNewPersonCreate(f, t)
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
	updatedDaoPerson := getPersonByKey(originalPersonCreate.PublicKey, t)
	if updatedDaoPerson.IsMajor != true {
		t.Error("Person was not updated")
	}
	if updatedDaoPerson.IsSigned != false {
		t.Error("Person should not have been signed")
	}
	updatedStatePerson := getStatePerson(updatedDaoPerson.Id, t)
	if updatedStatePerson.IsMajor != true {
		t.Error("Person was not updated")
	}
	if updatedStatePerson.IsSigned != false {
		t.Error("Person should not have been signed")
	}
}

func TestPersonUpdateSetSigned(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestPersonUpdateSetSigned", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(originalPersonCreate *command.PersonCreate, t *testing.T) {
		doTestPersonUpdateSetSigned(originalPersonCreate, t)
		expectedBalance := SUFFICIENT_BALANCE - priceMajorCreatePerson - priceMajorChangePersonAuthorization
		checkStateBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
		checkDaoBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	}
	withNewPersonCreate(f, t)
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
	updatedDaoPerson := getPersonByKey(originalPersonCreate.PublicKey, t)
	if updatedDaoPerson.IsMajor != false {
		t.Error("Person should not have become major")
	}
	if updatedDaoPerson.IsSigned != true {
		t.Error("Person was not signed")
	}
	updatedStatePerson := getStatePerson(updatedDaoPerson.Id, t)
	if updatedStatePerson.IsMajor != false {
		t.Error("Person should not have become major")
	}
	if updatedStatePerson.IsSigned != true {
		t.Error("Person was not signed")
	}
}

func TestPersonUpdateUnsetMajor(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestPersonUpdateUnsetMajor", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(t *testing.T) {
		doTestPersonUpdateUnsetMajor(t)
		expectedBalance := SUFFICIENT_BALANCE - priceMajorChangePersonAuthorization
		checkStateBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
		checkDaoBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	}
	withLoggedInWithNewKey(f, t)
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
	updatedDaoPerson := getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t)
	if updatedDaoPerson.IsMajor != false {
		t.Error("Majorship was not unset")
	}
	if updatedDaoPerson.IsSigned != true {
		t.Error("Signed should not have been changed")
	}
	updatedStatePerson := getStatePerson(updatedDaoPerson.Id, t)
	if updatedStatePerson.IsMajor != false {
		t.Error("Majorship was not unset")
	}
	if updatedStatePerson.IsSigned != true {
		t.Error("Signed should not have been changed")
	}
}

func TestPersonUpdateUnsetSigned(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestPersonUpdateUnsetSigned", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(t *testing.T) {
		doTestPersonUpdateUnsetSigned(t)
		expectedBalance := SUFFICIENT_BALANCE - priceMajorChangePersonAuthorization
		checkStateBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
		checkDaoBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	}
	withLoggedInWithNewKey(f, t)
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
	updatedDaoPerson := getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t)
	if updatedDaoPerson.IsMajor != true {
		t.Error("Majorship should not have been changed")
	}
	if updatedDaoPerson.IsSigned != false {
		t.Error("Signed was not unset")
	}
	updatedStatePerson := getStatePerson(updatedDaoPerson.Id, t)
	if updatedStatePerson.IsMajor != true {
		t.Error("Majorship should not have been changed")
	}
	if updatedStatePerson.IsSigned != false {
		t.Error("Signed was not unset")
	}
}

func TestPersonUpdateIncBalance(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestPersonUpdateIncBalance", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	withLoggedInWithNewKey(doTestPersonUpdateIncBalance, t)
}

func doTestPersonUpdateIncBalance(t *testing.T) {
	doTestBootstrap(t)
	theBalanceIncrement := int32(50)
	originalPerson := getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t)
	cmd := command.GetPersonUpdateIncBalanceCommand(
		originalPerson.Id,
		theBalanceIncrement,
		originalPerson.Id,
		cliAlexandria.LoggedIn(),
		int32(0))
	err := command.RunCommandForTest(cmd, "transactionIdIncBalance", blockchainAccess)
	if err != nil {
		t.Error("Could not run person update inc balance command: " + err.Error())
	}
	updatedDaoPerson := getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t)
	expectedBalance := SUFFICIENT_BALANCE + theBalanceIncrement
	if updatedDaoPerson.Balance != expectedBalance {
		t.Error("Balance has not been incremented")
	}
	updatedStatePerson := getStatePerson(updatedDaoPerson.Id, t)
	if updatedStatePerson.Balance != expectedBalance {
		t.Error("Balance has not been incremented")
	}
}

func TestSettingsUpdate(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestSettingsUpdate", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(t *testing.T) {
		doTestSettingsUpdate(t)
		expectedBalance := SUFFICIENT_BALANCE - priceMajorEditSettings
		checkStateBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
		checkDaoBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	}
	withLoggedInWithNewKey(f, t)
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
	checkUpdatedStateSettings(getStateSettings(t), t)
	checkUpdatedDaoSettings(updated, t)
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

func checkUpdatedStateSettings(settings *model.StateSettings, t *testing.T) {
	if util.Abs(settings.CreatedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("CreatedOn mismatch")
	}
	if util.Abs(settings.ModifiedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("ModifiedOn mismatch")
	}
	if settings.PriceList.PriceMajorEditSettings != 201 {
		t.Error("PriceMajorEditSettings mismatch")
	}
	if settings.PriceList.PriceMajorCreatePerson != 202 {
		t.Error("PriceMajorCreatePerson mismatch")
	}
	if settings.PriceList.PriceMajorChangePersonAuthorization != 203 {
		t.Error("PriceMajorChangePersonAuthorization mismatch")
	}
	if settings.PriceList.PriceMajorChangeJournalAuthorization != 204 {
		t.Error("PriceMajorChangeJournalAuthorization mismatch")
	}
	if settings.PriceList.PricePersonEdit != 205 {
		t.Error("PricePersonEdit mismatch")
	}
	if settings.PriceList.PriceAuthorSubmitNewManuscript != 206 {
		t.Error("PriceAuthorSubmitNewManuscript mismatch")
	}
	if settings.PriceList.PriceAuthorSubmitNewVersion != 207 {
		t.Error("PriceAuthorSubmitNewVersion mismatch")
	}
	if settings.PriceList.PriceAuthorAcceptAuthorship != 208 {
		t.Error("PriceAuthorAcceptAuthorship mismatch")
	}
	if settings.PriceList.PriceReviewerSubmit != 209 {
		t.Error("PriceReviewerSubmit mismatch")
	}
	if settings.PriceList.PriceEditorAllowManuscriptReview != 210 {
		t.Error("PriceEditorAllowManuscriptReview mismatch")
	}
	if settings.PriceList.PriceEditorRejectManuscript != 211 {
		t.Error("PriceEditorRejectManuscript mismatch")
	}
	if settings.PriceList.PriceEditorPublishManuscript != 212 {
		t.Error("PriceEditorPublishManuscript mismatch")
	}
	if settings.PriceList.PriceEditorAssignManuscript != 213 {
		t.Error("PriceEditorAssignManuscript mismatch")
	}
	if settings.PriceList.PriceEditorCreateJournal != 214 {
		t.Error("PriceEditorCreateJournal mismatch")
	}
	if settings.PriceList.PriceEditorCreateVolume != 215 {
		t.Error("PriceEditorCreateVolume mismatch")
	}
	if settings.PriceList.PriceEditorEditJournal != 216 {
		t.Error("PriceEditorEditJournal mismatch")
	}
	if settings.PriceList.PriceEditorAddColleague != 217 {
		t.Error("PriceEditorAddColleague mismatch")
	}
	if settings.PriceList.PriceEditorAcceptDuty != 218 {
		t.Error("PriceEditorAcceptDuty mismatch")
	}

}
func checkUpdatedDaoSettings(updated *dao.Settings, t *testing.T) {
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

func TestJournalCreate(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestJournalCreate", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	withNewJournalCreate(doTestJournalCreate, t)
}

func TestJournalUpdateProperties(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestJournalUpdateProperties", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(journal *command.Journal, personCreate *command.PersonCreate, initialBalance int32, t *testing.T) {
		doTestJournalCreate(journal, personCreate, initialBalance, t)
		updated := &command.Journal{
			Title: "Changed title",
		}
		journalId := getTheOnlyDaoJournal(t).JournalId
		cmd := command.GetCommandJournalUpdateProperties(
			journalId,
			getOriginalCommandJournal(),
			updated,
			getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
			cliAlexandria.LoggedIn(),
			priceEditorEditJournal)
		err := command.RunCommandForTest(cmd, "transactionIdJournalUpdateProperties", blockchainAccess)
		if err != nil {
			t.Error(err)
		}
		checkDaoJournalUpdatedProperties(getTheOnlyDaoJournal(t), t)
		checkStateJournalUpdatedProperties(getStateJournal(journalId, t), t)
		expectedBalance := initialBalance - priceEditorCreateJournal - priceEditorEditJournal
		checkDaoBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
		checkStateBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	}
	withNewJournalCreate(f, t)
}

func checkDaoJournalUpdatedProperties(journal *dao.Journal, t *testing.T) {
	if journal.Title != "Changed title" {
		t.Error("Title mismatch")
	}
}

func checkStateJournalUpdatedProperties(journal *model.StateJournal, t *testing.T) {
	if journal.Title != "Changed title" {
		t.Error("Title mismatch")
	}
}

func TestJournalUpdateAuthorization(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestJournalUpdateAuthorization", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(journal *command.Journal, personCreate *command.PersonCreate, initialBalance int32, t *testing.T) {
		doTestJournalCreate(journal, personCreate, initialBalance, t)
		journalId := getTheOnlyDaoJournal(t).JournalId
		cmd := command.GetCommandJournalUpdateAuthorization(
			journalId,
			true,
			getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
			cliAlexandria.LoggedIn(),
			priceMajorChangeJournalAuthorization)
		err := command.RunCommandForTest(cmd, "transactionIdJournalUpdateProperties", blockchainAccess)
		if err != nil {
			t.Error(err)
		}
		checkDaoJournalUpdatedAuthorization(getTheOnlyDaoJournal(t), t)
		checkStateJournalUpdatedAuthorization(getStateJournal(journalId, t), t)
		expectedBalance := initialBalance - priceEditorCreateJournal - priceMajorChangeJournalAuthorization
		checkDaoBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
		checkStateBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	}
	withNewJournalCreate(f, t)
}

func checkDaoJournalUpdatedAuthorization(journal *dao.Journal, t *testing.T) {
	if journal.IsSigned != true {
		t.Error("Setting IsSigned of journal was not done")
	}
}

func checkStateJournalUpdatedAuthorization(journal *model.StateJournal, t *testing.T) {
	if journal.IsSigned != true {
		t.Error("Setting IsSigned of journal was not done")
	}
}

func TestJournalEditorResign(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestJournalEditorResign", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(journal *command.Journal, personCreate *command.PersonCreate, initialBalance int32, t *testing.T) {
		doTestJournalCreate(journal, personCreate, initialBalance, t)
		journalId := getTheOnlyDaoJournal(t).JournalId
		cmd := command.GetCommandEditorResign(
			journalId,
			getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
			cliAlexandria.LoggedIn())
		err := command.RunCommandForTest(cmd, "transactionIdJournalEditorResign", blockchainAccess)
		if err != nil {
			t.Error(err)
		}
		checkDaoJournalEditorResigned(getTheOnlyDaoJournal(t), t)
		checkStateJournalEditorResigned(getStateJournal(journalId, t), t)
		expectedBalance := initialBalance - priceEditorCreateJournal - 0
		checkDaoBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
		checkStateBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	}
	withNewJournalCreate(f, t)
}

func checkDaoJournalEditorResigned(journal *dao.Journal, t *testing.T) {
	if len(journal.AcceptedEditors) >= 1 {
		t.Error("Last editor was not removed")
	}
}

func checkStateJournalEditorResigned(journal *model.StateJournal, t *testing.T) {
	if len(journal.EditorInfo) >= 1 {
		t.Error("Last editor was not removed")
	}
}

func TestJournalNewEditor(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestJournalNewEditor", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(journal *command.Journal, personCreate *command.PersonCreate, initialBalance int32, t *testing.T) {
		doTestJournalCreate(journal, personCreate, initialBalance, t)
		journalId := getTheOnlyDaoJournal(t).JournalId
		doTestJournalEditorInvite(journalId, personCreate, t, initialBalance-priceEditorCreateJournal)
		err := cliAlexandria.Login(personPublicKeyFile, personPrivateKeyFile)
		if err != nil {
			t.Error("Could not login as newly proposed editor")
		}
		doTestJournalEditorAcceptDuty(journalId, t)
	}
	withNewJournalCreate(f, t)
}

func doTestJournalEditorAcceptDuty(journalId string, t *testing.T) {
	cmd := command.GetCommandEditorAcceptDuty(
		journalId,
		getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
		cliAlexandria.LoggedIn(),
		priceEditorAcceptDuty)
	err := command.RunCommandForTest(cmd, "transactionIdEditorAcceptDuty", blockchainAccess)
	if err != nil {
		t.Error(err)
	}
	checkDaoJournalEditorAcceptedDuty(getTheOnlyDaoJournal(t), t)
	checkStateJournalEditorStates(getStateJournal(journalId, t), 0, 2, t)
	expectedBalance := SUFFICIENT_BALANCE - priceEditorAcceptDuty
	checkDaoBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	checkStateBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
}

func checkDaoJournalEditorAcceptedDuty(journal *dao.Journal, t *testing.T) {
	if len(journal.AcceptedEditors) != 2 {
		t.Error("Expected exactly two accepted editors")
	}
}

func doTestJournalEditorInvite(journalId string, personCreate *command.PersonCreate, t *testing.T, initialBalance int32) {
	cmd := command.GetCommandEditorInvite(
		journalId,
		getPersonByKey(personCreate.PublicKey, t).Id,
		getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
		cliAlexandria.LoggedIn(),
		priceEditorAddColleague)
	err := command.RunCommandForTest(cmd, "transactionIdJournalEditorResign", blockchainAccess)
	if err != nil {
		t.Error(err)
	}
	checkDaoJournalEditorInvited(getTheOnlyDaoJournal(t), t)
	checkStateJournalEditorStates(getStateJournal(journalId, t), 1, 1, t)
	expectedBalance := initialBalance - priceEditorAddColleague
	checkDaoBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	checkStateBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
}

func checkDaoJournalEditorInvited(journal *dao.Journal, t *testing.T) {
	if len(journal.AcceptedEditors) != 1 {
		t.Error("New editor with proposed state should not be shown yet")
	}
}

func checkStateJournalEditorStates(
	journal *model.StateJournal, expectedNumProposed, expectedNumAccepted int, t *testing.T) {
	numProposed := 0
	numAccepted := 0
	for _, e := range journal.EditorInfo {
		switch e.EditorState {
		case model.EditorState_editorProposed:
			numProposed++
		case model.EditorState_editorAccepted:
			numAccepted++
		}
	}
	if numProposed != expectedNumProposed {
		t.Error(fmt.Sprintf("Expected %d proposed editors, got %d",
			expectedNumProposed, numProposed))
	}
	if numAccepted != expectedNumAccepted {
		t.Error(fmt.Sprintf("Expected %d accepted editors, got %d",
			expectedNumAccepted, numAccepted))
	}
}

func TestJournalDescription(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestJournalUpdateDescription", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(journal *command.Journal, personCreate *command.PersonCreate, initialBalance int32, t *testing.T) {
		doTestJournalCreate(journal, personCreate, initialBalance, t)
		journalId := getTheOnlyDaoJournal(t).JournalId
		description := []byte("This is the description")
		cmd := command.GetCommandJournalUpdateDescription(
			journalId,
			"",
			description,
			getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
			cliAlexandria.LoggedIn(),
			priceEditorEditJournal)
		err := command.RunCommandForTest(cmd, "transactionIdUpdateDescription", blockchainAccess)
		if err != nil {
			t.Error(err)
		}
		err = dao.VerifyJournalDescription(journalId, description)
		if err != nil {
			t.Error(err)
		}
		cmd = command.GetCommandJournalOmitDescription(
			journalId,
			getTheOnlyDaoJournal(t).Descriptionhash,
			getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
			cliAlexandria.LoggedIn(),
			priceEditorEditJournal)
		err = command.RunCommandForTest(cmd, "transactionIdRemoveDescription", blockchainAccess)
		if err != nil {
			t.Error(err)
		}
		err = dao.VerifyJournalDescription(journalId, []byte{})
		if err != nil {
			t.Error(err)
		}
	}
	withNewJournalCreate(f, t)
}

func TestCreateVolume(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestJournalUpdateDescription", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(journal *command.Journal, personCreate *command.PersonCreate, initialBalance int32, t *testing.T) {
		doTestJournalCreate(journal, personCreate, initialBalance, t)
		journalId := getTheOnlyDaoJournal(t).JournalId
		vol := &command.Volume{
			JournalId: journalId,
			Issue:     "My issue",
		}
		cmd, volumeId := command.GetCommandVolumeCreate(
			vol,
			getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
			cliAlexandria.LoggedIn(),
			priceEditorCreateVolume)
		err := command.RunCommandForTest(cmd, "transactionIdVolumeCreate", blockchainAccess)
		if err != nil {
			t.Error(err)
		}
		checkCreatedDaoVolume(volumeId, journalId, t)
		checkDaoGetSingleVolume(volumeId, journalId, t)
		checkCreatedStateVolume(getStateVolume(volumeId, t), volumeId, journalId, t)
		expectedBalance := initialBalance - priceEditorCreateJournal - priceEditorCreateVolume
		checkDaoBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
		checkStateBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	}
	withNewJournalCreate(f, t)
}

func checkCreatedDaoVolume(volumeId, journalId string, t *testing.T) {
	volumes, err := dao.GetVolumesOfJournal(journalId)
	if err != nil {
		t.Error(err)
	}
	if len(volumes) != 1 {
		t.Error("Expected to have exactly one volume")
	}
	theVolume := &volumes[0]
	checkDaoVolumeContents(theVolume, volumeId, journalId, t)
}

func checkDaoVolumeContents(actual *dao.Volume, expectedVolumeId, expectedJournalId string, t *testing.T) {
	if actual.VolumeId != expectedVolumeId {
		t.Error("volumeId mismatch")
	}
	if actual.JournalId != expectedJournalId {
		t.Error("journalId mismatch")
	}
	if actual.Issue != "My issue" {
		t.Error("issue mismatch")
	}
	if util.Abs(actual.CreatedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("CreatedOn mismatch")
	}
}

func checkCreatedStateVolume(theVolume *model.StateVolume, volumeId, journalId string, t *testing.T) {
	if theVolume.Id != volumeId {
		t.Error("volumeId mismatch")
	}
	if theVolume.JournalId != journalId {
		t.Error("journalId mismatch")
	}
	if theVolume.Issue != "My issue" {
		t.Error("issue mismatch")
	}
	if util.Abs(theVolume.CreatedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("CreatedOn mismatch")
	}
}

func checkDaoGetSingleVolume(volumeId, journalId string, t *testing.T) {
	volume, err := dao.GetVolume(volumeId)
	if err != nil {
		t.Error(err)
	}
	checkDaoVolumeContents(volume, volumeId, journalId, t)
}

func getStateVolume(volumeId string, t *testing.T) *model.StateVolume {
	data, err := blockchainAccess.GetState([]string{volumeId})
	if err != nil {
		t.Error(err)
	}
	if len(data) != 1 {
		t.Error("Expected to read one address")
	}
	volumeBytes := data[volumeId]
	volume := &model.StateVolume{}
	err = proto.Unmarshal(volumeBytes, volume)
	if err != nil {
		t.Error(err)
	}
	return volume
}

func TestManuscriptCreate(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestJournalUpdateDescription", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(
		manuscriptCreate *command.ManuscriptCreate,
		journal *command.Journal,
		personCreate *command.PersonCreate,
		initialBalance int32,
		t *testing.T) {
		doTestManuscriptCreate(manuscriptCreate, personCreate, initialBalance, t)
	}
	withNewManuscriptCreate(f, t)
}

func TestManuscriptCreateNewVersion(t *testing.T) {
	logger = log.New(os.Stdout, "integration.TestJournalUpdateDescription", log.Flags())
	blockchainAccess = command.NewBlockchainStub(dao.HandleEvent, logger)
	f := func(
		manuscriptCreate *command.ManuscriptCreate,
		journal *command.Journal,
		personCreate *command.PersonCreate,
		initialBalance int32,
		t *testing.T) {
		previousManuscriptId, threadId := doTestManuscriptCreate(manuscriptCreate, personCreate, initialBalance, t)
		manuscriptCreateNewVersion := &command.ManuscriptCreateNewVersion{
			TheManuscript: []byte("New version text"),
			CommitMsg:     "Next version",
			Title:         "My manuscript",
			AuthorId: []string{
				getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
				getPersonByKey(personCreate.PublicKey, t).Id,
			},
			PreviousManuscriptId: previousManuscriptId,
			ThreadId:             threadId,
			JournalId:            getTheOnlyDaoJournal(t).JournalId,
		}
		cmd, newManuscriptId := command.GetCommandManuscriptCreateNewVersion(
			manuscriptCreateNewVersion,
			getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
			cliAlexandria.LoggedIn(),
			priceAuthorSubmitNewVersion)
		err := command.RunCommandForTest(cmd, "transactionIdManuscriptCreateNewVersion", blockchainAccess)
		if err != nil {
			t.Error(err)
		}
		actualThreadId := checkCreatedStateManuscriptNewVersion(
			getStateManuscript(newManuscriptId),
			newManuscriptId,
			getTheOnlyDaoJournal(t).JournalId,
			getPersonByKey(personCreate.PublicKey, t).Id,
			t)
		if actualThreadId != threadId {
			t.Error("ThreadId mismatch")
		}
		daoManuscriptNewVersion, err := dao.GetManuscript(newManuscriptId)
		if err != nil {
			t.Error(err)
		}
		checkCreatedDaoManuscriptNewVersion(
			daoManuscriptNewVersion,
			newManuscriptId,
			threadId,
			getTheOnlyDaoJournal(t).JournalId,
			getPersonByKey(personCreate.PublicKey, t).Id,
			t)
		expectedBalance := initialBalance - priceAuthorSubmitNewManuscript - priceAuthorSubmitNewVersion
		checkDaoBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
		checkStateBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	}
	withNewManuscriptCreate(f, t)
}

func checkCreatedStateManuscriptNewVersion(
	manuscript *model.StateManuscript,
	manuscriptId string,
	journalId string,
	secondAuthorId string,
	t *testing.T) string {
	if manuscript.Id != manuscriptId {
		t.Error("ManuscriptId mismatch")
	}
	if util.Abs(model.GetCurrentTime()-manuscript.CreatedOn) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("CreatedOn mismatch")
	}
	if util.Abs(model.GetCurrentTime()-manuscript.ModifiedOn) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("ModifiedOn mismatch")
	}
	if manuscript.Hash != model.HashBytes([]byte("New version text")) {
		t.Error("Hash mismatch")
	}
	if !model.IsManuscriptThreadAddress(manuscript.ThreadId) {
		t.Error("ThreadId mismatch")
	}
	if manuscript.VersionNumber != int32(1) {
		t.Error("VersionNumber mismatch")
	}
	if manuscript.CommitMsg != "Next version" {
		t.Error("CommitMsg mismatch")
	}
	if manuscript.Title != "My manuscript" {
		t.Error("Title mismatch")
	}
	if len(manuscript.Author) != 2 {
		t.Error("Wrong number of authors")
	}
	checkCreatedStateFirstAuthor(manuscript.Author[0], t)
	checkCreatedStateSecondAuthor(manuscript.Author[1], secondAuthorId, t)
	if manuscript.Status != model.ManuscriptStatus_init {
		t.Error("Status mismatch")
	}
	if manuscript.JournalId != journalId {
		t.Error("JournalId mismatch")
	}
	if manuscript.VolumeId != "" {
		t.Error("VolumeId mismatch")
	}
	if manuscript.FirstPage != "" {
		t.Error("FirstPage mismatch")
	}
	if manuscript.LastPage != "" {
		t.Error("LastPage mismatch")
	}
	return manuscript.ThreadId
}

func checkCreatedDaoManuscriptNewVersion(
	manuscript *dao.Manuscript,
	manuscriptId,
	threadId string,
	journalId string,
	secondAuthorId string,
	t *testing.T) {
	if manuscriptId != manuscriptId {
		t.Error("ManuscriptId mismatch")
	}
	if util.Abs(model.GetCurrentTime()-manuscript.CreatedOn) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("CreatedOn mismatch")
	}
	if util.Abs(model.GetCurrentTime()-manuscript.ModifiedOn) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("ModifiedOn mismatch")
	}
	if manuscript.Hash != model.HashBytes([]byte("New version text")) {
		t.Error("Hash mismatch")
	}
	if manuscript.ThreadId != threadId {
		t.Error("ThreadId mismatch")
	}
	if manuscript.VersionNumber != int32(1) {
		t.Error("VersionNumber mismatch")
	}
	if manuscript.CommitMsg != "Next version" {
		t.Error("CommitMsg mismatch")
	}
	if manuscript.Title != "My manuscript" {
		t.Error("Title mismatch")
	}
	if manuscript.Status != model.GetManuscriptStatusString(model.ManuscriptStatus_init) {
		t.Error("Status mismatch")
	}
	if manuscript.JournalId != journalId {
		t.Error("JournalId mismatch")
	}
	if manuscript.VolumeId != "" {
		t.Error("VolumdId mismatch")
	}
	if manuscript.FirstPage != "" {
		t.Error("FirstPage mismatch")
	}
	if manuscript.LastPage != "" {
		t.Error("LastPage mismatch")
	}
	if manuscript.IsReviewable != false {
		t.Error("IsReviewable mismatch")
	}
	checkCreatedDaoFirstAuthor(manuscript.Authors[0], manuscriptId, t)
	checkCreatedDaoSecondAuthor(manuscript.Authors[1], manuscriptId, secondAuthorId, t)
}
