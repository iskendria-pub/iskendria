package command

import (
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"log"
	"os"
	"testing"
)

const priceAuthorSubmitNewManuscript = int32(5)

func getNonBootstrapCommandExecution(context *testContext) *nonBootstrapCommandExecution {
	logger := log.New(os.Stdout, "testAuthor", log.Flags())
	ba := NewBlockchainStub(nil, logger)
	totalUnmarshalledState := &unmarshalledState{
		emptyAddresses: make(map[string]bool),
		settings: &model.StateSettings{
			PriceList: &model.PriceList{
				PriceAuthorSubmitNewManuscript: priceAuthorSubmitNewManuscript,
			},
		},
		persons: map[string]*model.StatePerson{
			context.signer.Id: context.signer,
			context.other.Id:  context.other,
		},
		journals: map[string]*model.StateJournal{
			context.journal.Id: context.journal,
		},
		manuscripts:       make(map[string]*model.StateManuscript),
		manuscriptThreads: make(map[string]*model.StateManuscriptThread),
	}
	readUnmarshalledState := &unmarshalledState{
		emptyAddresses: make(map[string]bool),
		settings: &model.StateSettings{
			PriceList: &model.PriceList{
				PriceAuthorSubmitNewManuscript: priceAuthorSubmitNewManuscript,
			},
		},
		persons: map[string]*model.StatePerson{
			context.signer.Id: context.signer,
		},
		journals:          make(map[string]*model.StateJournal),
		manuscripts:       make(map[string]*model.StateManuscript),
		manuscriptThreads: make(map[string]*model.StateManuscriptThread),
	}
	data, err := totalUnmarshalledState.read([]string{
		context.signer.Id,
		context.other.Id,
		model.GetSettingsAddress(),
		context.journal.Id})
	if err != nil {
		panic(err)
	}
	writtenAddresses, err := ba.SetState(data)
	if err != nil {
		panic(err)
	}
	if len(writtenAddresses) != 4 {
		panic("Not all addresses were written")
	}
	return &nonBootstrapCommandExecution{
		verifiedSignerId:  context.signer.Id,
		price:             priceAuthorSubmitNewManuscript,
		timestamp:         model.GetCurrentTime(),
		blockchainAccess:  ba,
		unmarshalledState: readUnmarshalledState,
	}
}

type testContext struct {
	signer       *model.StatePerson
	other        *model.StatePerson
	journal      *model.StateJournal
	manuscriptId string
	threadId     string
}

func getTestContext() *testContext {
	return &testContext{
		signer:       getStatePerson("John"),
		other:        getStatePerson("Jane"),
		journal:      getStateJournal("My journal"),
		manuscriptId: model.CreateManuscriptAddress(),
		threadId:     model.CreateManuscriptThreadAddress(),
	}
}

func getStatePerson(name string) *model.StatePerson {
	personId := model.CreatePersonAddress()
	return &model.StatePerson{
		Id:   personId,
		Name: name,
	}
}

func getStateJournal(title string) *model.StateJournal {
	journalId := model.CreateJournalAddress()
	return &model.StateJournal{
		Id:    journalId,
		Title: title,
	}
}

func TestCheckManuscriptCreateWithTwoAuthors(t *testing.T) {
	context := getTestContext()
	nbce := getNonBootstrapCommandExecution(context)
	authorIds := []string{context.signer.Id, context.other.Id}
	c := getManuscriptCreateTestCommand(context, authorIds)
	u, err := nbce.checkManuscriptCreate(c)
	if err != nil {
		t.Error(err)
	}
	if len(u.updates) != 3 {
		t.Error("Expected three singleUpdate objects")
		return
	}
	actualCreateManuscript := u.updates[0].(*singleUpdateManuscriptCreate)
	actualSignerCreate := u.updates[1].(*singleUpdateAuthorCreate)
	actualOtherCreate := u.updates[2].(*singleUpdateAuthorCreate)
	if actualCreateManuscript == nil || actualSignerCreate == nil || actualOtherCreate == nil {
		t.Error("Unexpected singleUpdate types")
		return
	}
	checkManuscriptCreate(actualCreateManuscript, context, model.ManuscriptStatus_init, t)
	checkAuthorCreateSigner(actualSignerCreate, context, t)
	checkAuthorCreateOther(actualOtherCreate, context, t)
}

func getManuscriptCreateTestCommand(context *testContext, authorIds []string) *model.CommandManuscriptCreate {
	c := &model.CommandManuscriptCreate{
		ManuscriptId:       context.manuscriptId,
		ManuscriptThreadId: context.threadId,
		Hash:               "someHash",
		Title:              "My manuscript",
		AuthorId:           authorIds,
		JournalId:          context.journal.Id,
	}
	return c
}

func checkManuscriptCreate(
	actualCreateManuscript *singleUpdateManuscriptCreate,
	context *testContext,
	expectedManuscriptState model.ManuscriptStatus,
	t *testing.T) {
	if actualCreateManuscript.manuscriptId != context.manuscriptId {
		t.Error("ManuscriptId mismatch")
	}
	if actualCreateManuscript.versionNumber != int32(0) {
		t.Error("Expected version number zero")
	}
	if actualCreateManuscript.manuscriptThreadId != context.threadId {
		t.Error("ThreadId mismatch")
	}
	if actualCreateManuscript.journalId != context.journal.Id {
		t.Error("JournalId mismatch")
	}
	if actualCreateManuscript.status != expectedManuscriptState {
		t.Error("Status mismatch")
	}
}

func checkAuthorCreateSigner(
	actualSignerCreate *singleUpdateAuthorCreate,
	context *testContext,
	t *testing.T) {
	if actualSignerCreate.manuscriptId != context.manuscriptId {
		t.Error("actualSignerCreate manuscriptId mismatch")
	}
	if actualSignerCreate.authorNumber != int32(0) {
		t.Error("actualSignerCreate authorNumber mismatch")
	}
	if actualSignerCreate.didSign != true {
		t.Error("actualSignerCreate didSign mismatch")
	}
	if actualSignerCreate.authorId != context.signer.Id {
		t.Error("actualSignerCreate authorId mismatch")
	}
}

func checkAuthorCreateOther(
	actualOtherCreate *singleUpdateAuthorCreate,
	context *testContext,
	t *testing.T) {
	if actualOtherCreate.manuscriptId != context.manuscriptId {
		t.Error("actualOtherCreate manuscriptId mismatch")
	}
	if actualOtherCreate.authorNumber != int32(1) {
		t.Error("actualOtherCreate authorNumber mismatch")
	}
	if actualOtherCreate.didSign != false {
		t.Error("actualOtherCreate didSign mismatch")
	}
	if actualOtherCreate.authorId != context.other.Id {
		t.Error("actualOtherCreate authorId mismatch")
	}
}

func TestCheckManuscriptCreateWithOneAuthors(t *testing.T) {
	context := getTestContext()
	nbce := getNonBootstrapCommandExecution(context)
	authorIds := []string{context.signer.Id}
	c := getManuscriptCreateTestCommand(context, authorIds)
	u, err := nbce.checkManuscriptCreate(c)
	if err != nil {
		t.Error(err)
	}
	if len(u.updates) != 2 {
		t.Error("Expected two singleUpdate objects")
		return
	}
	actualCreateManuscript := u.updates[0].(*singleUpdateManuscriptCreate)
	actualSignerCreate := u.updates[1].(*singleUpdateAuthorCreate)
	if actualCreateManuscript == nil || actualSignerCreate == nil {
		t.Error("Unexpected singleUpdate types")
		return
	}
	checkManuscriptCreate(actualCreateManuscript, context, model.ManuscriptStatus_new, t)
	checkAuthorCreateSigner(actualSignerCreate, context, t)
}
