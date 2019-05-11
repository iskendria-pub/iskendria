package integration

import (
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

func withInitializedDao(testFunc func(logger *log.Logger, t *testing.T), logger *log.Logger, t *testing.T) {
	dao.Init("testBootstrap.db", logger)
	defer dao.ShutdownAndDelete(logger)
	err := dao.StartFakeBlock("blockId", "")
	if err != nil {
		t.Error("Error starting fake block: " + err.Error())
	}
	testFunc(logger, t)
}

func withLoggedInWithNewKey(
	testFunc func(blockchainAccess command.BlockchainAccess, _ *log.Logger, t *testing.T),
	blockchainAccess command.BlockchainAccess,
	logger *log.Logger,
	t *testing.T) {
	withLogin := func(logger *log.Logger, t *testing.T) {
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
		testFunc(blockchainAccess, logger, t)
	}
	withInitializedDao(withLogin, logger, t)
}

func withNewPersonCreate(
	testFunc func(personCreate *command.PersonCreate, _ command.BlockchainAccess, _ *log.Logger, t *testing.T),
	blockchainAccess command.BlockchainAccess,
	logger *log.Logger,
	t *testing.T) {
	withNewPersonCreate := func(blockchainAccess command.BlockchainAccess, logger *log.Logger, t *testing.T) {
		doTestBootstrap(blockchainAccess, logger, t)
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
		testFunc(personCreate, blockchainAccess, logger, t)
	}
	withLoggedInWithNewKey(withNewPersonCreate, blockchainAccess, logger, t)
}

func doTestBootstrap(blockchainAccess command.BlockchainAccess, _ *log.Logger, t *testing.T) {
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
	checkSettings(readSettings, t)
	checkPerson(person, t)
}

func getBootstrap() *command.Bootstrap {
	return &command.Bootstrap{
		PriceMajorEditSettings:               101,
		PriceMajorCreatePerson:               102,
		PriceMajorChangePersonAuthorization:  103,
		PriceMajorChangeJournalAuthorization: 104,
		PricePersonEdit:                      105,
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

func checkSettings(settings *dao.Settings, t *testing.T) {
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

func checkPerson(person *dao.Person, t *testing.T) {
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

func doTestPersonCreate(
	personCreate *command.PersonCreate,
	blockchainAccess command.BlockchainAccess,
	_ *log.Logger,
	t *testing.T) {
	newPersonKey := personCreate.PublicKey
	personCreateCommand := command.GetPersonCreateCommand(
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
	checkCreatedPerson(createdPersons[0], newPersonKey, t)
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

func checkCreatedPerson(person *dao.Person, expectedPublicKey string, t *testing.T) {
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