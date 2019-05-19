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
	"testing"
)

const TIME_DIFF_THRESHOLD_SECONDS = 10

const personPublicKeyFile = "person.pub"
const personPrivateKeyFile = "person.priv"

const SUFFICIENT_BALANCE = int32(1000)

const priceMajorEditSettings int32 = 101
const priceMajorCreatePerson int32 = 102
const priceMajorChangePersonAuthorization int32 = 103
const pricePersonEdit int32 = 105

var logger *log.Logger
var blockchainAccess command.BlockchainAccess

func withInitializedDao(testFunc func(t *testing.T), t *testing.T) {
	dao.Init("testBootstrap.db", logger)
	defer dao.ShutdownAndDelete(logger)
	err := dao.StartFakeBlock("blockId", "")
	if err != nil {
		t.Error("Error starting fake block: " + err.Error())
	}
	testFunc(t)
}

func withLoggedInWithNewKey(testFunc func(t *testing.T), t *testing.T) {
	withLogin := func(t *testing.T) {
		publicKeyFile := "testBootstrap.pub"
		privateKeyFile := "testBootstrap.priv"
		err := cliAlexandria.CreateKeyPair(publicKeyFile, privateKeyFile)
		if err != nil {
			t.Error("Could not create keypair: " + err.Error())
		}
		defer cliAlexandria.RemoveKeyFiles(publicKeyFile, privateKeyFile, logger)
		err = cliAlexandria.Login(publicKeyFile, privateKeyFile)
		if err != nil {
			t.Error("Could not login: " + err.Error())
		}
		defer func() { _ = cliAlexandria.Logout() }()
		testFunc(t)
	}
	withInitializedDao(withLogin, t)
}

func withNewPersonCreate(testFunc func(personCreate *command.PersonCreate, t *testing.T), t *testing.T) {
	withNewPersonCreate := func(t *testing.T) {
		doTestBootstrap(t)
		err := cliAlexandria.CreateKeyPair(personPublicKeyFile, personPrivateKeyFile)
		if err != nil {
			t.Error("Could not create key pair for new person")
		}
		defer cliAlexandria.RemoveKeyFiles(personPublicKeyFile, personPrivateKeyFile, logger)
		_, personPublicKey, err := cliAlexandria.ReadPublicKeyFile(personPublicKeyFile)
		if err != nil {
			t.Error("Could not read public key file of new person")
		}
		personCreate := &command.PersonCreate{
			PublicKey: personPublicKey,
			Name:      "Rens",
			Email:     "rens@xxx.nl",
		}
		testFunc(personCreate, t)
	}
	withLoggedInWithNewKey(withNewPersonCreate, t)
}

func doTestBootstrap(t *testing.T) {
	bootstrap := getBootstrap()
	transactionId := "transactionId"
	commandBootstrap := command.GetBootstrapCommand(bootstrap, cliAlexandria.LoggedIn())
	err := command.RunCommandForTest(commandBootstrap, transactionId, blockchainAccess)
	if err != nil {
		t.Error("Error executing command: " + err.Error())
	}
	readSettings, err := dao.GetSettings()
	if err != nil {
		t.Error("After doing bootstrap, could not read settings from database: " + err.Error())
	}
	persons, err := dao.SearchPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr)
	if err != nil {
		t.Error("After doing bootstrap, could not read person")
	}
	if len(persons) != 1 {
		t.Error("Expected that bootstrapping would produce exactly one person")
	}
	person := persons[0]
	checkBootstrapStateSettings(getStateSettings(t), t)
	checkBootstrapStatePerson(getStatePerson(person.Id, t), t)
	checkBootstrapDaoSettings(readSettings, t)
	checkBootstrapDaoPerson(person, t)
	addSufficientBalance(person.Id, t)
}

func getStateSettings(t *testing.T) *model.StateSettings {
	settingsData, err := blockchainAccess.GetState([]string{model.GetSettingsAddress()})
	if err != nil {
		t.Error("integration.getStateSettings: Cannot read settings address: " + err.Error())
	}
	if len(settingsData) != 1 {
		t.Error("integration.getStateSettings: No blockchain state for settings")
	}
	settings := &model.StateSettings{}
	err = proto.Unmarshal(settingsData[model.GetSettingsAddress()], settings)
	if err != nil {
		t.Error("integration.getStateSettings: Could not unmarshall")
	}
	return settings
}

func getStatePerson(personId string, t *testing.T) *model.StatePerson {
	personData, err := blockchainAccess.GetState([]string{personId})
	if err != nil {
		t.Error("integration.getStatePerson: Cannot read person address: " + err.Error())
	}
	if len(personData) != 1 {
		t.Error("integration.getStatePerson: No blockchain state for person: " + personId)
	}
	person := &model.StatePerson{}
	err = proto.Unmarshal(personData[personId], person)
	if err != nil {
		t.Error("integration.getStatePerson: Could not unmarshall")
	}
	return person
}

func getBootstrap() *command.Bootstrap {
	return &command.Bootstrap{
		PriceMajorEditSettings:               priceMajorEditSettings,
		PriceMajorCreatePerson:               priceMajorCreatePerson,
		PriceMajorChangePersonAuthorization:  priceMajorChangePersonAuthorization,
		PriceMajorChangeJournalAuthorization: 104,
		PricePersonEdit:                      pricePersonEdit,
		PriceAuthorSubmitNewManuscript:       106,
		PriceAuthorSubmitNewVersion:          107,
		PriceAuthorAcceptAuthorship:          108,
		PriceReviewerSubmit:                  109,
		PriceEditorAllowManuscriptReview:     110,
		PriceEditorRejectManuscript:          111,
		PriceEditorPublishManuscript:         112,
		PriceEditorAssignManuscript:          113,
		PriceEditorCreateJournal:             114,
		PriceEditorCreateVolume:              115,
		PriceEditorEditJournal:               116,
		PriceEditorAddColleague:              117,
		PriceEditorAcceptDuty:                118,
		Name:                                 "Brita",
		Email:                                "brita@xxx.nl",
	}
}

func checkBootstrapStateSettings(settings *model.StateSettings, t *testing.T) {
	if util.Abs(settings.CreatedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("CreatedOn mismatch")
	}
	if util.Abs(settings.ModifiedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("ModifiedOn mismatch")
	}
	if settings.PriceList.PriceMajorEditSettings != 101 {
		t.Error("PriceMajorEditSettings mismatch")
	}
	if settings.PriceList.PriceMajorCreatePerson != 102 {
		t.Error("PriceMajorCreatePerson mismatch")
	}
	if settings.PriceList.PriceMajorChangePersonAuthorization != 103 {
		t.Error("PriceMajorChangePersonAuthorization mismatch")
	}
	if settings.PriceList.PriceMajorChangeJournalAuthorization != 104 {
		t.Error("PriceMajorChangeJournalAuthorization mismatch")
	}
	if settings.PriceList.PricePersonEdit != 105 {
		t.Error("PricePersonEdit mismatch")
	}
	if settings.PriceList.PriceAuthorSubmitNewManuscript != 106 {
		t.Error("PriceAuthorSubmitNewManuscript mismatch")
	}
	if settings.PriceList.PriceAuthorSubmitNewVersion != 107 {
		t.Error("PriceAuthorSubmitNewVersion mismatch")
	}
	if settings.PriceList.PriceAuthorAcceptAuthorship != 108 {
		t.Error("PriceAuthorAcceptAuthorship mismatch")
	}
	if settings.PriceList.PriceReviewerSubmit != 109 {
		t.Error("PriceReviewerSubmit mismatch")
	}
	if settings.PriceList.PriceEditorAllowManuscriptReview != 110 {
		t.Error("PriceEditorAllowManuscriptReview mismatch")
	}
	if settings.PriceList.PriceEditorRejectManuscript != 111 {
		t.Error("PriceEditorRejectManuscript mismatch")
	}
	if settings.PriceList.PriceEditorPublishManuscript != 112 {
		t.Error("PriceEditorPublishManuscript mismatch")
	}
	if settings.PriceList.PriceEditorAssignManuscript != 113 {
		t.Error("PriceEditorAssignManuscript mismatch")
	}
	if settings.PriceList.PriceEditorCreateJournal != 114 {
		t.Error("PriceEditorCreateJournal mismatch")
	}
	if settings.PriceList.PriceEditorCreateVolume != 115 {
		t.Error("PriceEditorCreateVolume mismatch")
	}
	if settings.PriceList.PriceEditorEditJournal != 116 {
		t.Error("PriceEditorEditJournal mismatch")
	}
	if settings.PriceList.PriceEditorAddColleague != 117 {
		t.Error("PriceEditorAddColleague mismatch")
	}
	if settings.PriceList.PriceEditorAcceptDuty != 118 {
		t.Error("PriceEditorAcceptDuty mismatch")
	}
}

func checkBootstrapDaoSettings(settings *dao.Settings, t *testing.T) {
	if settings.PriceMajorEditSettings != 101 {
		t.Error("PriceMajorEditSettings mismatch")
	}
	if settings.PriceMajorCreatePerson != 102 {
		t.Error("PriceMajorCreatePerson mismatch")
	}
	if settings.PriceMajorChangePersonAuthorization != 103 {
		t.Error("PriceMajorChangePersonAuthorization mismatch")
	}
	if settings.PriceMajorChangeJournalAuthorization != 104 {
		t.Error("PriceMajorChangeJournalAuthorization mismatch")
	}
	if settings.PricePersonEdit != 105 {
		t.Error("PricePersonEdit mismatch")
	}
	if settings.PriceAuthorSubmitNewManuscript != 106 {
		t.Error("PriceAuthorSubmitNewManuscript mismatch")
	}
	if settings.PriceAuthorSubmitNewVersion != 107 {
		t.Error("PriceAuthorSubmitNewVersion mismatch")
	}
	if settings.PriceAuthorAcceptAuthorship != 108 {
		t.Error("PriceAuthorAcceptAuthorship mismatch")
	}
	if settings.PriceReviewerSubmit != 109 {
		t.Error("PriceReviewerSubmit mismatch")
	}
	if settings.PriceEditorAllowManuscriptReview != 110 {
		t.Error("PriceEditorAllowManuscriptReview mismatch")
	}
	if settings.PriceEditorRejectManuscript != 111 {
		t.Error("PriceEditorRejectManuscript mismatch")
	}
	if settings.PriceEditorPublishManuscript != 112 {
		t.Error("PriceEditorPublishManuscript mismatch")
	}
	if settings.PriceEditorAssignManuscript != 113 {
		t.Error("PriceEditorAssignManuscript mismatch")
	}
	if settings.PriceEditorCreateJournal != 114 {
		t.Error("PriceEditorCreateJournal mismatch")
	}
	if settings.PriceEditorCreateVolume != 115 {
		t.Error("PriceEditorCreateVolume mismatch")
	}
	if settings.PriceEditorEditJournal != 116 {
		t.Error("PriceEditorEditJournal mismatch")
	}
	if settings.PriceEditorAddColleague != 117 {
		t.Error("PriceEditorAddColleague mismatch")
	}
	if settings.PriceEditorAcceptDuty != 118 {
		t.Error("PriceEditorAcceptDuty mismatch")
	}
}

func checkBootstrapStatePerson(person *model.StatePerson, t *testing.T) {
	if !model.IsPersonAddress(person.Id) {
		t.Error("Id mismatch")
	}
	if util.Abs(person.CreatedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("CreatedOn mismatch")
	}
	if util.Abs(person.ModifiedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("ModifiedOn mismatch")
	}
	if person.PublicKey != cliAlexandria.LoggedIn().PublicKeyStr {
		t.Error("PublicKey mismatch")
	}
	if person.Name != "Brita" {
		t.Error("Name mismatch")
	}
	if person.Email != "brita@xxx.nl" {
		t.Error("Email mismatch")
	}
	if person.IsMajor != true {
		t.Error("IsMajor mismatch")
	}
	if person.IsSigned != true {
		t.Error("IsSigned mismatch")
	}
	if person.Balance != int32(0) {
		t.Error("Balance mismatch")
	}
	if person.BiographyHash != "" {
		t.Error("BiographyHash mismatch")
	}
	if person.Organization != "" {
		t.Error("Organization mismatch")
	}
	if person.Telephone != "" {
		t.Error("Telephone mismatch")
	}
	if person.Address != "" {
		t.Error("Address mismatch")
	}
	if person.PostalCode != "" {
		t.Error("PostalCode mismatch")
	}
	if person.Country != "" {
		t.Error("Country mismatch")
	}
	if person.ExtraInfo != "" {
		t.Error("ExtraInfo mismatch")
	}
}

func checkBootstrapDaoPerson(person *dao.Person, t *testing.T) {
	if !model.IsPersonAddress(person.Id) {
		t.Error("Id mismatch")
	}
	if util.Abs(person.CreatedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("CreatedOn mismatch")
	}
	if util.Abs(person.ModifiedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("ModifiedOn mismatch")
	}
	if person.PublicKey != cliAlexandria.LoggedIn().PublicKeyStr {
		t.Error("PublicKey mismatch")
	}
	if person.Name != "Brita" {
		t.Error("Name mismatch")
	}
	if person.Email != "brita@xxx.nl" {
		t.Error("Email mismatch")
	}
	if person.IsMajor != true {
		t.Error("IsMajor mismatch")
	}
	if person.IsSigned != true {
		t.Error("IsSigned mismatch")
	}
	if person.Balance != int32(0) {
		t.Error("Balance mismatch")
	}
	if person.BiographyHash != "" {
		t.Error("BiographyHash mismatch")
	}
	if person.Organization != "" {
		t.Error("Organization mismatch")
	}
	if person.Telephone != "" {
		t.Error("Telephone mismatch")
	}
	if person.Address != "" {
		t.Error("Address mismatch")
	}
	if person.PostalCode != "" {
		t.Error("PostalCode mismatch")
	}
	if person.Country != "" {
		t.Error("Country mismatch")
	}
	if person.ExtraInfo != "" {
		t.Error("ExtraInfo mismatch")
	}
}

func addSufficientBalance(personId string, t *testing.T) {
	cmd := command.GetPersonUpdateIncBalanceCommand(
		personId,
		SUFFICIENT_BALANCE,
		getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
		cliAlexandria.LoggedIn(),
		int32(0))
	err := command.RunCommandForTest(cmd, "transactionIdAddSufficientBalance", blockchainAccess)
	if err != nil {
		t.Error("integration.addSufficientBalance could not run command: " + err.Error())
	}
	person, err := dao.GetPersonById(personId)
	if err != nil {
		t.Error("integration.addSufficientBalance could not read person from db")
		return
	}
	if person.Balance != SUFFICIENT_BALANCE {
		t.Error("integration.addSufficientBalance: balance was not updated")
	}
}

func doTestPersonCreate(personCreate *command.PersonCreate, t *testing.T) {
	newPersonKey := personCreate.PublicKey
	personCreateCommand, newPersonId := command.GetPersonCreateCommand(
		personCreate,
		getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
		cliAlexandria.LoggedIn(),
		int32(102))
	err := command.RunCommandForTest(personCreateCommand, "secondTransaction", blockchainAccess)
	if err != nil {
		t.Error("Running person create command failed: " + err.Error())
	}
	createdPersons, err := dao.SearchPersonByKey(newPersonKey)
	if err != nil {
		t.Error("Could not read newly created person")
	}
	if len(createdPersons) != 1 {
		t.Error("Expected exactly one newly created person")
	}

	createdPerson := getPersonByKey(newPersonKey, t)
	if createdPerson.Id != newPersonId {
		t.Error("Returned created person id differs from person id found in database")
	}
	checkCreatedStatePerson(getStatePerson(createdPerson.Id, t), newPersonKey, t)
	checkCreatedDaoPerson(createdPersons[0], newPersonKey, t)
	addSufficientBalance(createdPersons[0].Id, t)
}

func checkDaoBalanceOfKey(expectedSignerBalance int32, key string, t *testing.T) {
	signerPerson := getPersonByKey(key, t)
	if signerPerson.Balance != expectedSignerBalance {
		t.Error(fmt.Sprintf("Signer balance on client side not OK, expected %d, got %d",
			expectedSignerBalance, signerPerson.Balance))
	}
}

func checkStateBalanceOfKey(expectedSignerBalance int32, key string, t *testing.T) {
	daoSignerPerson := getPersonByKey(key, t)
	signerPerson := getStatePerson(daoSignerPerson.Id, t)
	if signerPerson.Balance != expectedSignerBalance {
		t.Error(fmt.Sprintf("Signer balance on blockchain not OK, expected %d, got %d",
			expectedSignerBalance, signerPerson.Balance))
	}
}

func getPersonByKey(key string, t *testing.T) *dao.Person {
	persons, err := dao.SearchPersonByKey(key)
	if err != nil {
		t.Error("Could not find person for logged in key")
	}
	if len(persons) != 1 {
		t.Error("Expected only one person id for signing key")
	}
	return persons[0]
}

func checkCreatedStatePerson(person *model.StatePerson, expectedPublicKey string, t *testing.T) {
	if !model.IsPersonAddress(person.Id) {
		t.Error("Id is not a person address")
	}
	if person.Id == getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id {
		t.Error("Id equals the id of the existing person")
	}
	if util.Abs(person.CreatedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("CreatedOn mismatch")
	}
	if util.Abs(person.ModifiedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("ModifiedOn mismatch")
	}
	if person.PublicKey != expectedPublicKey {
		t.Error("PublicKey mismatch")
	}
	if person.Name != "Rens" {
		t.Error("Name mismatch")
	}
	if person.Email != "rens@xxx.nl" {
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
	if person.BiographyHash != "" {
		t.Error("BiographyHash mismatch")
	}
	if person.Organization != "" {
		t.Error("Organization mismatch")
	}
	if person.Telephone != "" {
		t.Error("Telephone mismatch")
	}
	if person.Address != "" {
		t.Error("Address mismatch")
	}
	if person.PostalCode != "" {
		t.Error("PostalCode mismatch")
	}
	if person.Country != "" {
		t.Error("Country mismatch")
	}
	if person.ExtraInfo != "" {
		t.Error("ExtraInfo mismatch")
	}
}

func checkCreatedDaoPerson(person *dao.Person, expectedPublicKey string, t *testing.T) {
	if !model.IsPersonAddress(person.Id) {
		t.Error("Id is not a person address")
	}
	if person.Id == getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id {
		t.Error("Id equals the id of the existing person")
	}
	if util.Abs(person.CreatedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("CreatedOn mismatch")
	}
	if util.Abs(person.ModifiedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("ModifiedOn mismatch")
	}
	if person.PublicKey != expectedPublicKey {
		t.Error("PublicKey mismatch")
	}
	if person.Name != "Rens" {
		t.Error("Name mismatch")
	}
	if person.Email != "rens@xxx.nl" {
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
	if person.BiographyHash != "" {
		t.Error("BiographyHash mismatch")
	}
	if person.Organization != "" {
		t.Error("Organization mismatch")
	}
	if person.Telephone != "" {
		t.Error("Telephone mismatch")
	}
	if person.Address != "" {
		t.Error("Address mismatch")
	}
	if person.PostalCode != "" {
		t.Error("PostalCode mismatch")
	}
	if person.Country != "" {
		t.Error("Country mismatch")
	}
	if person.ExtraInfo != "" {
		t.Error("ExtraInfo mismatch")
	}
}

func getSettings(t *testing.T) *dao.Settings {
	settings, err := dao.GetSettings()
	if err != nil {
		t.Error(err)
	}
	return settings
}
