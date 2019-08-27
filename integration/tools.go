package integration

import (
	"errors"
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

var majorName = "Brita"

const SUFFICIENT_BALANCE = int32(1000)

const THE_LOGICAL_PUBLICATION_TIME = int64(481)

const priceMajorEditSettings int32 = 101
const priceMajorCreatePerson int32 = 102
const priceMajorChangePersonAuthorization int32 = 103
const priceMajorChangeJournalAuthorization = 104
const pricePersonEdit int32 = 105
const priceAuthorSubmitNewManuscript int32 = 106
const priceAuthorSubmitNewVersion int32 = 107
const priceAuthorAcceptAuthorship int32 = 108
const priceReviewerSubmit int32 = 109
const priceEditorAllowManuscriptReview int32 = 110
const priceEditorRejectManuscript int32 = 111
const priceEditorPublishManuscript int32 = 112
const priceEditorAssignManuscript int32 = 113
const priceEditorCreateJournal int32 = 114
const priceEditorCreateVolume int32 = 115
const priceEditorEditJournal int32 = 116
const priceEditorAddColleague int32 = 117
const priceEditorAcceptDuty int32 = 118

var logger *log.Logger
var blockchainAccess command.BlockchainAccess

func withInitializedDao(testFunc func(t *testing.T), t *testing.T) {
	dao.Init("testBootstrap.db", logger)
	defer dao.ShutdownAndDelete(logger)
	err := dao.StartFakeBlock("blockId", "", logger)
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

func withNewJournalCreate(testFunc func(*command.Journal, *command.PersonCreate, int32, *testing.T), t *testing.T) {
	journal := getOriginalCommandJournal()
	f := func(personCreate *command.PersonCreate, t *testing.T) {
		doTestPersonCreate(personCreate, t)
		initialBalance := SUFFICIENT_BALANCE - priceMajorCreatePerson
		checkStateBalanceOfKey(initialBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
		checkDaoBalanceOfKey(initialBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
		testFunc(journal, personCreate, initialBalance, t)
	}
	withNewPersonCreate(f, t)
}

func withNewManuscriptCreate(
	testFunc func(*command.ManuscriptCreate, *command.Journal, *command.PersonCreate, int32, *testing.T),
	numAuthors int,
	t *testing.T) {
	f := func(
		journal *command.Journal,
		personCreate *command.PersonCreate,
		initialBalance int32,
		t *testing.T) {
		doTestJournalCreate(journal, personCreate, initialBalance, t)
		initialBalance -= priceEditorCreateJournal
		manuscriptCreate := &command.ManuscriptCreate{
			TheManuscript: []byte("Lorem ipsum"),
			CommitMsg:     "Initial version",
			Title:         "My Manuscript",
			AuthorId:      getAuthorsForWithNewManuscriptId(numAuthors, personCreate, t),
			JournalId:     getTheOnlyDaoJournal(t).JournalId,
		}
		testFunc(manuscriptCreate, journal, personCreate, initialBalance, t)
	}
	withNewJournalCreate(f, t)
}

func withReviewCreated(
	testFunc func(review *dao.Review, manuscript *dao.Manuscript, initialBalance int32, t *testing.T),
	t *testing.T) {
	f := func(
		manuscriptCreate *command.ManuscriptCreate,
		journal *command.Journal,
		personCreate *command.PersonCreate,
		initialBalance int32,
		t *testing.T) {
		review, manuscript, initialBalance := doTestWritePositiveReview(manuscriptCreate, initialBalance, t)
		testFunc(review, manuscript, initialBalance, t)
	}
	withNewManuscriptCreate(f, 1, t)
}

func getAuthorsForWithNewManuscriptId(
	numAuthors int, personCreate *command.PersonCreate, t *testing.T) []string {
	result := make([]string, numAuthors)
	for i := 0; i < numAuthors; i++ {
		switch i {
		case 0:
			result[0] = getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id
		case 1:
			result[1] = getPersonByKey(personCreate.PublicKey, t).Id
		}
	}
	return result
}

func getOriginalCommandJournal() *command.Journal {
	return &command.Journal{
		Title: "The Journal",
	}
}

func checkCreatedDaoJournal(journal *dao.Journal, journalId string, editorId string, t *testing.T) {
	if journal.JournalId != journalId {
		t.Error("JournalId mismatch")
	}
	if util.Abs(journal.CreatedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("CreatedOn mismatch")
	}
	if util.Abs(journal.ModifiedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("ModifiedOn mismatch")
	}
	if journal.Title != "The Journal" {
		t.Error("Title mismatch")
	}
	if journal.IsSigned != false {
		t.Error("IsSigned mismatch")
	}
	if journal.Descriptionhash != "" {
		t.Error("DescriptionHash mismatch")
	}
	if len(journal.AcceptedEditors) != 1 {
		t.Error("Length mismatch of accepted editors")
	}
	if journal.AcceptedEditors[0].PersonId != editorId {
		t.Error("Editor PersonId mismatch")
	}
	if journal.AcceptedEditors[0].PersonName != majorName {
		t.Error(fmt.Sprintf("Editor PersonName mismatch, expected %s, got %s",
			majorName, journal.AcceptedEditors[0].PersonName))
	}
}

func checkCreatedStateJournal(journal *model.StateJournal, journalId, editorId string, t *testing.T) {
	if journal.Id != journalId {
		t.Error("Id mismatch")
	}
	if journal.Title != "The Journal" {
		t.Error("Title mismatch")
	}
	if journal.IsSigned != false {
		t.Error("IsSigned mismatch")
	}
	if util.Abs(journal.ModifiedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("ModifiedOn mismatch")
	}
	if util.Abs(journal.CreatedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("CreatedOn mismatch")
	}
	if journal.DescriptionHash != "" {
		t.Error("DescriptionHash mismatch")
	}
	if len(journal.EditorInfo) != 1 {
		t.Error("EditorInfo length mismatch")
	}
	if journal.EditorInfo[0].EditorId != editorId {
		t.Error("EditorInfo.EditorId mismatch")
	}
	if journal.EditorInfo[0].EditorState != model.EditorState_editorAccepted {
		t.Error("EditorInfo.EditorState mismatch")
	}
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
		PriceMajorChangeJournalAuthorization: priceMajorChangeJournalAuthorization,
		PricePersonEdit:                      pricePersonEdit,
		PriceAuthorSubmitNewManuscript:       priceAuthorSubmitNewManuscript,
		PriceAuthorSubmitNewVersion:          priceAuthorSubmitNewVersion,
		PriceAuthorAcceptAuthorship:          priceAuthorAcceptAuthorship,
		PriceReviewerSubmit:                  priceReviewerSubmit,
		PriceEditorAllowManuscriptReview:     priceEditorAllowManuscriptReview,
		PriceEditorRejectManuscript:          priceEditorRejectManuscript,
		PriceEditorPublishManuscript:         priceEditorPublishManuscript,
		PriceEditorAssignManuscript:          priceEditorAssignManuscript,
		PriceEditorCreateJournal:             priceEditorCreateJournal,
		PriceEditorCreateVolume:              priceEditorCreateVolume,
		PriceEditorEditJournal:               priceEditorEditJournal,
		PriceEditorAddColleague:              priceEditorAddColleague,
		PriceEditorAcceptDuty:                priceEditorAcceptDuty,
		Name:                                 majorName,
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
	if settings.PriceList.PriceMajorEditSettings != priceMajorEditSettings {
		t.Error("PriceMajorEditSettings mismatch")
	}
	if settings.PriceList.PriceMajorCreatePerson != priceMajorCreatePerson {
		t.Error("PriceMajorCreatePerson mismatch")
	}
	if settings.PriceList.PriceMajorChangePersonAuthorization != priceMajorChangePersonAuthorization {
		t.Error("PriceMajorChangePersonAuthorization mismatch")
	}
	if settings.PriceList.PriceMajorChangeJournalAuthorization != priceMajorChangeJournalAuthorization {
		t.Error("PriceMajorChangeJournalAuthorization mismatch")
	}
	if settings.PriceList.PricePersonEdit != pricePersonEdit {
		t.Error("PricePersonEdit mismatch")
	}
	if settings.PriceList.PriceAuthorSubmitNewManuscript != priceAuthorSubmitNewManuscript {
		t.Error("PriceAuthorSubmitNewManuscript mismatch")
	}
	if settings.PriceList.PriceAuthorSubmitNewVersion != priceAuthorSubmitNewVersion {
		t.Error("PriceAuthorSubmitNewVersion mismatch")
	}
	if settings.PriceList.PriceAuthorAcceptAuthorship != priceAuthorAcceptAuthorship {
		t.Error("PriceAuthorAcceptAuthorship mismatch")
	}
	if settings.PriceList.PriceReviewerSubmit != priceReviewerSubmit {
		t.Error("PriceReviewerSubmit mismatch")
	}
	if settings.PriceList.PriceEditorAllowManuscriptReview != priceEditorAllowManuscriptReview {
		t.Error("PriceEditorAllowManuscriptReview mismatch")
	}
	if settings.PriceList.PriceEditorRejectManuscript != priceEditorRejectManuscript {
		t.Error("PriceEditorRejectManuscript mismatch")
	}
	if settings.PriceList.PriceEditorPublishManuscript != priceEditorPublishManuscript {
		t.Error("PriceEditorPublishManuscript mismatch")
	}
	if settings.PriceList.PriceEditorAssignManuscript != priceEditorAssignManuscript {
		t.Error("PriceEditorAssignManuscript mismatch")
	}
	if settings.PriceList.PriceEditorCreateJournal != priceEditorCreateJournal {
		t.Error("PriceEditorCreateJournal mismatch")
	}
	if settings.PriceList.PriceEditorCreateVolume != priceEditorCreateVolume {
		t.Error("PriceEditorCreateVolume mismatch")
	}
	if settings.PriceList.PriceEditorEditJournal != priceEditorEditJournal {
		t.Error("PriceEditorEditJournal mismatch")
	}
	if settings.PriceList.PriceEditorAddColleague != priceEditorAddColleague {
		t.Error("PriceEditorAddColleague mismatch")
	}
	if settings.PriceList.PriceEditorAcceptDuty != priceEditorAcceptDuty {
		t.Error("PriceEditorAcceptDuty mismatch")
	}
}

func checkBootstrapDaoSettings(settings *dao.Settings, t *testing.T) {
	if settings.PriceMajorEditSettings != priceMajorEditSettings {
		t.Error("PriceMajorEditSettings mismatch")
	}
	if settings.PriceMajorCreatePerson != priceMajorCreatePerson {
		t.Error("PriceMajorCreatePerson mismatch")
	}
	if settings.PriceMajorChangePersonAuthorization != priceMajorChangePersonAuthorization {
		t.Error("PriceMajorChangePersonAuthorization mismatch")
	}
	if settings.PriceMajorChangeJournalAuthorization != priceMajorChangeJournalAuthorization {
		t.Error("PriceMajorChangeJournalAuthorization mismatch")
	}
	if settings.PricePersonEdit != pricePersonEdit {
		t.Error("PricePersonEdit mismatch")
	}
	if settings.PriceAuthorSubmitNewManuscript != priceAuthorSubmitNewManuscript {
		t.Error("PriceAuthorSubmitNewManuscript mismatch")
	}
	if settings.PriceAuthorSubmitNewVersion != priceAuthorSubmitNewVersion {
		t.Error("PriceAuthorSubmitNewVersion mismatch")
	}
	if settings.PriceAuthorAcceptAuthorship != priceAuthorAcceptAuthorship {
		t.Error("PriceAuthorAcceptAuthorship mismatch")
	}
	if settings.PriceReviewerSubmit != priceReviewerSubmit {
		t.Error("PriceReviewerSubmit mismatch")
	}
	if settings.PriceEditorAllowManuscriptReview != priceEditorAllowManuscriptReview {
		t.Error("PriceEditorAllowManuscriptReview mismatch")
	}
	if settings.PriceEditorRejectManuscript != priceEditorRejectManuscript {
		t.Error("PriceEditorRejectManuscript mismatch")
	}
	if settings.PriceEditorPublishManuscript != priceEditorPublishManuscript {
		t.Error("PriceEditorPublishManuscript mismatch")
	}
	if settings.PriceEditorAssignManuscript != priceEditorAssignManuscript {
		t.Error("PriceEditorAssignManuscript mismatch")
	}
	if settings.PriceEditorCreateJournal != priceEditorCreateJournal {
		t.Error("PriceEditorCreateJournal mismatch")
	}
	if settings.PriceEditorCreateVolume != priceEditorCreateVolume {
		t.Error("PriceEditorCreateVolume mismatch")
	}
	if settings.PriceEditorEditJournal != priceEditorEditJournal {
		t.Error("PriceEditorEditJournal mismatch")
	}
	if settings.PriceEditorAddColleague != priceEditorAddColleague {
		t.Error("PriceEditorAddColleague mismatch")
	}
	if settings.PriceEditorAcceptDuty != priceEditorAcceptDuty {
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
	if person.Name != majorName {
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
	if person.Name != majorName {
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
		int32(priceMajorCreatePerson))
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

func doTestJournalCreate(journal *command.Journal, _ *command.PersonCreate, initialBalance int32, t *testing.T) {
	editorId := getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id
	cmd, journalId := command.GetCommandJournalCreate(
		journal,
		editorId,
		cliAlexandria.LoggedIn(),
		priceEditorCreateJournal)
	err := command.RunCommandForTest(cmd, "transactionJournalCreate", blockchainAccess)
	if err != nil {
		t.Error()
	}
	checkCreatedDaoJournal(getTheOnlyDaoJournal(t), journalId, editorId, t)
	checkCreatedStateJournal(getStateJournal(journalId, t), journalId, editorId, t)
	expectedBalance := initialBalance - priceEditorCreateJournal
	checkDaoBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	checkStateBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
}

func getTheOnlyDaoJournal(t *testing.T) *dao.Journal {
	journals, err := dao.GetAllJournals()
	if err != nil {
		t.Error(err)
	}
	if len(journals) != 1 {
		t.Error(fmt.Sprintf("Expected to have exactly one journal, but got %d", len(journals)))
	}
	actualJournal := journals[0]
	return actualJournal
}

func getStateJournal(journalId string, t *testing.T) *model.StateJournal {
	data, err := blockchainAccess.GetState([]string{journalId})
	if err != nil {
		t.Error(err)
	}
	if len(data) != 1 {
		t.Error("Expected to read one address")
	}
	journalBytes := data[journalId]
	journal := &model.StateJournal{}
	err = proto.Unmarshal(journalBytes, journal)
	if err != nil {
		t.Error(err)
	}
	return journal
}

func doTestManuscriptCreate(
	manuscriptCreate *command.ManuscriptCreate,
	personCreate *command.PersonCreate,
	initialBalance int32,
	t *testing.T) (string, string) {
	cmd, manuscriptId := command.GetCommandManuscriptCreate(
		manuscriptCreate,
		getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
		cliAlexandria.LoggedIn(),
		priceAuthorSubmitNewManuscript)
	err := command.RunCommandForTest(cmd, "transactionIdAuthorSubmitNewJournal", blockchainAccess)
	if err != nil {
		t.Error(err)
	}
	threadId := checkCreatedStateManuscript(
		getStateManuscript(manuscriptId),
		manuscriptId,
		getTheOnlyDaoJournal(t).JournalId,
		getPersonByKey(personCreate.PublicKey, t).Id,
		t)
	daoManuscript, err := dao.GetManuscript(manuscriptId)
	if err != nil {
		t.Error(err)
	}
	checkCreatedDaoManuscript(
		daoManuscript,
		manuscriptId,
		threadId,
		getTheOnlyDaoJournal(t).JournalId,
		getPersonByKey(personCreate.PublicKey, t).Id,
		t)
	expectedBalance := initialBalance - priceAuthorSubmitNewManuscript
	checkDaoBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	checkStateBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	return manuscriptId, threadId
}

func getStateManuscript(manuscriptId string) *model.StateManuscript {
	resultMap, err := blockchainAccess.GetState([]string{manuscriptId})
	if err != nil {
		panic(err)
	}
	if len(resultMap) != 1 {
		panic("Did not find manuscriptId: " + manuscriptId)
	}
	result := &model.StateManuscript{}
	err = proto.Unmarshal(resultMap[manuscriptId], result)
	if err != nil {
		panic(err)
	}
	return result
}

func checkCreatedStateManuscript(
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
	if manuscript.Hash != model.HashBytes([]byte("Lorem ipsum")) {
		t.Error("Hash mismatch")
	}
	if !model.IsManuscriptThreadAddress(manuscript.ThreadId) {
		t.Error("ThreadId mismatch")
	}
	if manuscript.VersionNumber != int32(0) {
		t.Error("VersionNumber mismatch")
	}
	if manuscript.CommitMsg != "Initial version" {
		t.Error("CommitMsg mismatch")
	}
	if manuscript.Title != "My Manuscript" {
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

func checkCreatedStateFirstAuthor(author *model.Author, t *testing.T) {
	if author.AuthorId != getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id {
		t.Error("First author id mismatch")
	}
	if author.AuthorNumber != int32(0) {
		t.Error("First author number should be zero")
	}
	if author.DidSign != true {
		t.Error("First author was expected to be signed")
	}
}

func checkCreatedStateSecondAuthor(author *model.Author, expectedAuthorId string, t *testing.T) {
	if author.AuthorId != expectedAuthorId {
		t.Error("Second author id mismatch")
	}
	if author.AuthorNumber != 1 {
		t.Error("Second author number should be one")
	}
	if author.DidSign != false {
		t.Error("Second author was expected not to be signed")
	}
}

func checkCreatedDaoManuscript(
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
	if manuscript.Hash != model.HashBytes([]byte("Lorem ipsum")) {
		t.Error("Hash mismatch")
	}
	if manuscript.ThreadId != threadId {
		t.Error("ThreadId mismatch")
	}
	if manuscript.VersionNumber != int32(0) {
		t.Error("VersionNumber mismatch")
	}
	if manuscript.CommitMsg != "Initial version" {
		t.Error("CommitMsg mismatch")
	}
	if manuscript.Title != "My Manuscript" {
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

func checkCreatedDaoFirstAuthor(author *dao.Author, manuscriptId string, t *testing.T) {
	if author.ManuscriptId != manuscriptId {
		t.Error("First author manuscriptId mismatch")
	}
	if author.PersonId != getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id {
		t.Error(("First author personId mismatch"))
	}
	if author.DidSign != true {
		t.Error("First author didSign mismatch")
	}
	if author.AuthorNumber != 0 {
		t.Error("First author authorNumber mismatch")
	}
	if author.PersonName != majorName {
		t.Error("First author PersonName mismatch")
	}
}

func checkCreatedDaoSecondAuthor(author *dao.Author, manuscriptId, secondAuthorId string, t *testing.T) {
	if author.ManuscriptId != manuscriptId {
		t.Error("Second author manuscriptId mismatch")
	}
	if author.PersonId != secondAuthorId {
		t.Error("Second author personId mismatch")
	}
	if author.AuthorNumber != 1 {
		t.Error("Second author authorNumber mismatch")
	}
	if author.DidSign != false {
		t.Error("Second author didSign mismatch")
	}
	if author.PersonName != "Rens" {
		t.Error("Second author PersonName mismatch")
	}
}

func getStateThread(threadId string, t *testing.T) *model.StateManuscriptThread {
	data, err := blockchainAccess.GetState([]string{threadId})
	if err != nil {
		t.Error(err)
	}
	stateBytes, ok := data[threadId]
	if !ok {
		t.Error(errors.New("Thread address was not filled: " + threadId))
	}
	state := &model.StateManuscriptThread{}
	err = proto.Unmarshal(stateBytes, state)
	if err != nil {
		t.Error(err)
	}
	return state
}

func doTestWritePositiveReview(manuscriptCreate *command.ManuscriptCreate, initialBalance int32, t *testing.T) (
	*dao.Review, *dao.Manuscript, int32) {
	cmdManuscriptCreate, manuscriptId := command.GetCommandManuscriptCreate(
		manuscriptCreate,
		getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
		cliAlexandria.LoggedIn(),
		priceAuthorSubmitNewManuscript)
	err := command.RunCommandForTest(
		cmdManuscriptCreate, "transactionIdManuscriptCreateOneAuthor", blockchainAccess)
	if err != nil {
		t.Error(err)
	}
	_ = runEditorAllowReview(manuscriptId, t)
	reviewCreate := &command.ReviewCreate{
		ManuscriptId: manuscriptId,
		TheReview:    []byte("My review"),
	}
	cmdWriteReview, reviewId := command.GetCommandWritePositiveReview(
		reviewCreate,
		getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
		cliAlexandria.LoggedIn(),
		priceReviewerSubmit)
	err = command.RunCommandForTest(cmdWriteReview, "transactionIdWriteReview", blockchainAccess)
	if err != nil {
		t.Error(err)
	}
	daoReview, err := dao.GetReview(reviewId)
	if err != nil {
		t.Error(err)
	}
	checkStateReview(getStateReview(reviewId, t), reviewId, manuscriptId, model.Judgement_POSITIVE, t)
	checkDaoReview(daoReview, reviewId, manuscriptId, model.Judgement_POSITIVE, t)
	expectedBalance := initialBalance -
		priceAuthorSubmitNewManuscript -
		priceEditorAllowManuscriptReview -
		priceReviewerSubmit
	checkStateBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	checkDaoBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	finalManuscript, err := dao.GetManuscript(manuscriptId)
	if err != nil {
		t.Error(err)
	}
	return daoReview, finalManuscript, expectedBalance
}

func runEditorAllowReview(manuscriptId string, t *testing.T) *dao.Manuscript {
	manuscript, err := dao.GetManuscript(manuscriptId)
	if err != nil {
		t.Error(err)
	}
	threadReference, err := dao.GetReferenceThread(manuscript.ThreadId)
	if err != nil {
		t.Error(err)
	}
	cmdAllowReview := command.GetCommandManuscriptAllowReview(
		manuscript.ThreadId,
		threadReference,
		getTheOnlyDaoJournal(t).JournalId,
		getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
		cliAlexandria.LoggedIn(),
		priceEditorAllowManuscriptReview)
	err = command.RunCommandForTest(
		cmdAllowReview, "transactionIdManuscriptAllowReview", blockchainAccess)
	if err != nil {
		t.Error(err)
	}
	return manuscript
}

func checkStateReview(
	r *model.StateReview, reviewId, manuscriptId string, expectedJudgement model.Judgement, t *testing.T) {
	if r.Id != reviewId {
		t.Error("ReviewId mismatch")
	}
	if util.Abs(r.CreatedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("CreatedOn mismatch")
	}
	if r.ManuscriptId != manuscriptId {
		t.Error("ManuscriptId mismatch")
	}
	if r.ReviewAuthorId != getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id {
		t.Error("ReviewAuthorId mismatch")
	}
	if r.Hash != model.HashBytes([]byte("My review")) {
		t.Error("Hash mismatch")
	}
	if r.Judgement != expectedJudgement {
		t.Error("Judgement mismatch")
	}
	if r.IsUsedByEditor != false {
		t.Error("IsUsedByEditor mismatch")
	}
}

func checkDaoReview(r *dao.Review, reviewId, manuscriptId string, expectedJudgement model.Judgement, t *testing.T) {
	if r.Id != reviewId {
		t.Error("ReviewId mismatch")
	}
	if util.Abs(r.CreatedOn-model.GetCurrentTime()) >= TIME_DIFF_THRESHOLD_SECONDS {
		t.Error("CreatedOn mismatch")
	}
	if r.ManuscriptId != manuscriptId {
		t.Error("ManuscriptId mismatch")
	}
	if r.ReviewAuthorId != getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id {
		t.Error("ReviewAuthorId mismatch")
	}
	if r.Hash != model.HashBytes([]byte("My review")) {
		t.Error("Hash mismatch")
	}
	if r.Judgement != model.GetJudgementString(expectedJudgement) {
		t.Error("Judgement mismatch")
	}
	if r.IsUsedByEditor != false {
		t.Error("IsUsedByEditor mismatch")
	}
}

func doTestWriteNegativeReview(manuscriptCreate *command.ManuscriptCreate, initialBalance int32, t *testing.T) (
	*dao.Review, *dao.Manuscript, int32) {
	cmdManuscriptCreate, manuscriptId := command.GetCommandManuscriptCreate(
		manuscriptCreate,
		getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
		cliAlexandria.LoggedIn(),
		priceAuthorSubmitNewManuscript)
	err := command.RunCommandForTest(
		cmdManuscriptCreate, "transactionIdManuscriptCreateOneAuthor", blockchainAccess)
	if err != nil {
		t.Error(err)
	}
	_ = runEditorAllowReview(manuscriptId, t)
	reviewCreate := &command.ReviewCreate{
		ManuscriptId: manuscriptId,
		TheReview:    []byte("My review"),
	}
	cmdWriteReview, reviewId := command.GetCommandWriteNegativeReview(
		reviewCreate,
		getPersonByKey(cliAlexandria.LoggedIn().PublicKeyStr, t).Id,
		cliAlexandria.LoggedIn(),
		priceReviewerSubmit)
	err = command.RunCommandForTest(cmdWriteReview, "transactionIdWriteReview", blockchainAccess)
	if err != nil {
		t.Error(err)
	}
	daoReview, err := dao.GetReview(reviewId)
	if err != nil {
		t.Error(err)
	}
	checkStateReview(getStateReview(reviewId, t), reviewId, manuscriptId, model.Judgement_NEGATIVE, t)
	checkDaoReview(daoReview, reviewId, manuscriptId, model.Judgement_NEGATIVE, t)
	expectedBalance := initialBalance -
		priceAuthorSubmitNewManuscript -
		priceEditorAllowManuscriptReview -
		priceReviewerSubmit
	checkStateBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	checkDaoBalanceOfKey(expectedBalance, cliAlexandria.LoggedIn().PublicKeyStr, t)
	finalManuscript, err := dao.GetManuscript(manuscriptId)
	if err != nil {
		t.Error(err)
	}
	return daoReview, finalManuscript, expectedBalance
}

func getStateReview(reviewId string, t *testing.T) *model.StateReview {
	data, err := blockchainAccess.GetState([]string{reviewId})
	if err != nil {
		t.Error(err)
	}
	stateBytes, ok := data[reviewId]
	if !ok {
		t.Error(errors.New("Review address was not filled: " + reviewId))
	}
	state := &model.StateReview{}
	err = proto.Unmarshal(stateBytes, state)
	if err != nil {
		t.Error(err)
	}
	return state
}
