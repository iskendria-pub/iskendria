package command

import (
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"log"
	"os"
	"testing"
)

const priceAuthorSubmitNewManuscript = int32(5)
const priceAuthorSubmitNewVersion = int32(6)

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
	checkAuthorCreateSigner(actualSignerCreate, context.signer.Id, context.manuscriptId, t)
	checkAuthorCreateOther(actualOtherCreate, context.other.Id, context.manuscriptId, t)
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
	signerId string,
	manuscriptId string,
	t *testing.T) {
	if actualSignerCreate.manuscriptId != manuscriptId {
		t.Error("actualSignerCreate manuscriptId mismatch")
	}
	if actualSignerCreate.authorNumber != int32(0) {
		t.Error("actualSignerCreate authorNumber mismatch")
	}
	if actualSignerCreate.didSign != true {
		t.Error("actualSignerCreate didSign mismatch")
	}
	if actualSignerCreate.authorId != signerId {
		t.Error("actualSignerCreate authorId mismatch")
	}
}

func checkAuthorCreateOther(
	actualOtherCreate *singleUpdateAuthorCreate,
	otherId string,
	manuscriptId string,
	t *testing.T) {
	if actualOtherCreate.manuscriptId != manuscriptId {
		t.Error("actualOtherCreate manuscriptId mismatch")
	}
	if actualOtherCreate.authorNumber != int32(1) {
		t.Error("actualOtherCreate authorNumber mismatch")
	}
	if actualOtherCreate.didSign != false {
		t.Error("actualOtherCreate didSign mismatch")
	}
	if actualOtherCreate.authorId != otherId {
		t.Error("actualOtherCreate authorId mismatch")
	}
}

func TestCheckManuscriptCreateWithOneAuthor(t *testing.T) {
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
	checkAuthorCreateSigner(actualSignerCreate, context.signer.Id, context.manuscriptId, t)
}

func getNonBootstrapCommandExecutionWithInitialManuscript(
	context *testContextWithInitialManuscript) *nonBootstrapCommandExecution {
	logger := log.New(os.Stdout, "testAuthor", log.Flags())
	ba := NewBlockchainStub(nil, logger)
	totalUnmarshalledState := &unmarshalledState{
		emptyAddresses: make(map[string]bool),
		settings: &model.StateSettings{
			PriceList: &model.PriceList{
				PriceAuthorSubmitNewVersion: priceAuthorSubmitNewVersion,
			},
		},
		persons: map[string]*model.StatePerson{
			context.signer.Id: context.signer,
			context.other.Id:  context.other,
		},
		journals: map[string]*model.StateJournal{
			context.journal.Id: context.journal,
		},
		manuscripts: map[string]*model.StateManuscript{
			context.initialManuscriptId: {
				Id:            context.initialManuscriptId,
				Hash:          "2468ace0",
				ThreadId:      context.threadId,
				VersionNumber: 0,
				Title:         "My Test Manuscript",
				Author: []*model.Author{
					{
						AuthorId:     context.signer.Id,
						DidSign:      false,
						AuthorNumber: 0,
					},
					{
						AuthorId:     context.other.Id,
						DidSign:      false,
						AuthorNumber: 1,
					},
				},
				Status:    model.ManuscriptStatus_init,
				JournalId: context.journal.Id,
			},
		},
		manuscriptThreads: map[string]*model.StateManuscriptThread{
			context.threadId: {
				Id:           context.threadId,
				ManuscriptId: []string{context.initialManuscriptId},
				IsReviewable: context.isThreadReviewable,
			},
		},
	}
	readUnmarshalledState := &unmarshalledState{
		emptyAddresses: make(map[string]bool),
		settings: &model.StateSettings{
			PriceList: &model.PriceList{
				PriceAuthorSubmitNewVersion: priceAuthorSubmitNewVersion,
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
		context.journal.Id,
		context.initialManuscriptId,
		context.threadId,
	})
	if err != nil {
		panic(err)
	}
	writtenAddresses, err := ba.SetState(data)
	if err != nil {
		panic(err)
	}
	if len(writtenAddresses) != 6 {
		panic("Not all addresses were written")
	}
	return &nonBootstrapCommandExecution{
		verifiedSignerId:  context.signer.Id,
		price:             priceAuthorSubmitNewVersion,
		timestamp:         model.GetCurrentTime(),
		blockchainAccess:  ba,
		unmarshalledState: readUnmarshalledState,
	}
}

type testContextWithInitialManuscript struct {
	signer              *model.StatePerson
	other               *model.StatePerson
	journal             *model.StateJournal
	initialManuscriptId string
	threadId            string
	isThreadReviewable  bool
	manuscriptId        string
	manuscriptStatus    model.ManuscriptStatus
}

func TestCheckManuscriptCreateNewVersionWithTwoAuthors(t *testing.T) {
	contexts := []testContextWithInitialManuscript{
		*getTestContextWithInitialManuscript(false, model.ManuscriptStatus_init),
		*getTestContextWithInitialManuscript(true, model.ManuscriptStatus_init),
	}
	for _, context := range contexts {
		nbce := getNonBootstrapCommandExecutionWithInitialManuscript(&context)
		authorIds := []string{context.signer.Id, context.other.Id}
		c := getManuscriptCreateNewVersionTestCommand(&context, authorIds)
		u, err := nbce.checkManuscriptCreateNewVersion(c)
		if err != nil {
			t.Error(err)
		}
		if len(u.updates) != 3 {
			t.Error("Expected three singleUpdate objects")
			return
		}
		actualCreateManuscript := u.updates[0].(*singleUpdateManuscriptCreateNewVersion)
		actualSignerCreate := u.updates[1].(*singleUpdateAuthorCreate)
		actualOtherCreate := u.updates[2].(*singleUpdateAuthorCreate)
		if actualCreateManuscript == nil || actualSignerCreate == nil || actualOtherCreate == nil {
			t.Error("Unexpected singleUpdate types")
			return
		}
		checkManuscriptCreateNewVersion(actualCreateManuscript, &context, context.manuscriptStatus, t)
		checkAuthorCreateSigner(actualSignerCreate, context.signer.Id, context.manuscriptId, t)
		checkAuthorCreateOther(actualOtherCreate, context.other.Id, context.manuscriptId, t)
	}
}

func getTestContextWithInitialManuscript(
	isThreadReviewable bool,
	manuscriptStatus model.ManuscriptStatus) *testContextWithInitialManuscript {
	return &testContextWithInitialManuscript{
		signer:              getStatePerson("John"),
		other:               getStatePerson("Jane"),
		journal:             getStateJournal("My journal"),
		initialManuscriptId: model.CreateManuscriptAddress(),
		threadId:            model.CreateManuscriptThreadAddress(),
		isThreadReviewable:  isThreadReviewable,
		manuscriptId:        model.CreateManuscriptAddress(),
		manuscriptStatus:    manuscriptStatus,
	}
}

func getManuscriptCreateNewVersionTestCommand(
	context *testContextWithInitialManuscript, authorIds []string) *model.CommandManuscriptCreateNewVersion {
	c := &model.CommandManuscriptCreateNewVersion{
		ManuscriptId:         context.manuscriptId,
		PreviousManuscriptId: context.initialManuscriptId,
		Hash:                 "someOtherHash",
		CommitMsg:            "Next version",
		Title:                "My manuscript",
		AuthorId:             authorIds,
	}
	return c
}

func checkManuscriptCreateNewVersion(
	actualCreateManuscript *singleUpdateManuscriptCreateNewVersion,
	context *testContextWithInitialManuscript,
	expectedManuscriptState model.ManuscriptStatus,
	t *testing.T) {
	if actualCreateManuscript.manuscriptId != context.manuscriptId {
		t.Error("ManuscriptId mismatch")
	}
	if actualCreateManuscript.versionNumber != int32(1) {
		t.Error("Expected version number one")
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

func TestCheckManuscriptCreateNewVersionWithOneAuthor(t *testing.T) {
	contexts := []testContextWithInitialManuscript{
		*getTestContextWithInitialManuscript(false, model.ManuscriptStatus_new),
		*getTestContextWithInitialManuscript(true, model.ManuscriptStatus_reviewable),
	}
	for _, context := range contexts {
		nbce := getNonBootstrapCommandExecutionWithInitialManuscript(&context)
		authorIds := []string{context.signer.Id}
		c := getManuscriptCreateNewVersionTestCommand(&context, authorIds)
		u, err := nbce.checkManuscriptCreateNewVersion(c)
		if err != nil {
			t.Error(err)
		}
		if len(u.updates) != 2 {
			t.Error("Expected two singleUpdate objects")
			return
		}
		actualCreateManuscript := u.updates[0].(*singleUpdateManuscriptCreateNewVersion)
		actualSignerCreate := u.updates[1].(*singleUpdateAuthorCreate)
		if actualCreateManuscript == nil || actualSignerCreate == nil {
			t.Error("Unexpected singleUpdate types")
			return
		}
		checkManuscriptCreateNewVersion(actualCreateManuscript, &context, context.manuscriptStatus, t)
		checkAuthorCreateSigner(actualSignerCreate, context.signer.Id, context.manuscriptId, t)
	}
}

func TestGetCommandManuscriptAcceptAuthorshipWork(t *testing.T) {
	contexts := []acceptAuthorshipContext{
		{
			numAuthors:                       1,
			numSigner:                        0,
			numbersAlreadySigned:             []int{},
			isThreadReviewable:               false,
			expectedDoesAuthorUpdate:         true,
			expectedAllAuthorsWillHaveSigned: true,
			expectedNewStatus:                model.ManuscriptStatus_new,
		},
		{
			numAuthors:                       1,
			numSigner:                        0,
			numbersAlreadySigned:             []int{},
			isThreadReviewable:               true,
			expectedDoesAuthorUpdate:         true,
			expectedAllAuthorsWillHaveSigned: true,
			expectedNewStatus:                model.ManuscriptStatus_reviewable,
		},
		{
			numAuthors:                       2,
			numSigner:                        1,
			numbersAlreadySigned:             []int{0},
			isThreadReviewable:               false,
			expectedDoesAuthorUpdate:         true,
			expectedAllAuthorsWillHaveSigned: true,
			expectedNewStatus:                model.ManuscriptStatus_new,
		},
		{
			numAuthors:                       2,
			numSigner:                        1,
			numbersAlreadySigned:             []int{},
			isThreadReviewable:               false,
			expectedDoesAuthorUpdate:         true,
			expectedAllAuthorsWillHaveSigned: false,
			expectedNewStatus:                model.ManuscriptStatus_init,
		},
	}
	for _, c := range contexts {
		authorIds := make([]string, c.numAuthors)
		for i := range authorIds {
			authorIds[i] = model.CreatePersonAddress()
		}
		authors := make([]*model.Author, c.numAuthors)
		for authorNumber := range authors {
			didSign := false
			for _, authorNumberThatSigned := range c.numbersAlreadySigned {
				if authorNumber == authorNumberThatSigned {
					didSign = true
				}
			}
			authors[authorNumber] = &model.Author{
				AuthorId:     authorIds[authorNumber],
				DidSign:      didSign,
				AuthorNumber: int32(authorNumber),
			}
		}
		cmd := &model.CommandManuscriptAcceptAuthorship{
			Author: authors,
		}
		actualDoesAuthorUpdate, actualAllAuthorsWillHaveSigned := getCommandManuscriptAcceptAuthorshipWork(
			cmd, authorIds[c.numSigner])
		if actualDoesAuthorUpdate != c.expectedDoesAuthorUpdate {
			t.Error("DoesAuthorUpdate mismatch")
		}
		if actualAllAuthorsWillHaveSigned != c.expectedAllAuthorsWillHaveSigned {
			t.Error("AllAuthorsWillHaveSigned mismatch")
		}
		actualNewManuscriptStatus := getNewManuscriptStatus(actualAllAuthorsWillHaveSigned, c.isThreadReviewable)
		if actualNewManuscriptStatus != c.expectedNewStatus {
			t.Error(fmt.Sprintf("NewManuscriptStatus mismatch, expected %s got %s",
				model.GetManuscriptStatusString(c.expectedNewStatus),
				model.GetManuscriptStatusString(actualNewManuscriptStatus)))
		}
	}
}

type acceptAuthorshipContext struct {
	numAuthors                       int
	numSigner                        int
	numbersAlreadySigned             []int
	isThreadReviewable               bool
	expectedDoesAuthorUpdate         bool
	expectedAllAuthorsWillHaveSigned bool
	expectedNewStatus                model.ManuscriptStatus
}
