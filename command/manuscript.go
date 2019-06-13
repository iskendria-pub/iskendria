package command

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"gitlab.bbinfra.net/3estack/alexandria/util"
	"strconv"
)

type ManuscriptCreate struct {
	TheManuscript []byte
	CommitMsg     string
	Title         string
	AuthorId      []string
	JournalId     string
}

func GetCommandManuscriptCreate(
	manuscriptCreate *ManuscriptCreate,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) (*Command, string) {
	manuscriptId := model.CreateManuscriptAddress()
	threadId := model.CreateManuscriptThreadAddress()
	theHash := model.HashBytes(manuscriptCreate.TheManuscript)
	return &Command{
		InputAddresses: append(
			manuscriptCreate.AuthorId,
			model.GetSettingsAddress(),
			signerId,
			manuscriptCreate.JournalId,
			manuscriptId,
			threadId),
		OutputAddresses: []string{signerId, manuscriptId, threadId},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signerId,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandManuscriptCreate{
				CommandManuscriptCreate: &model.CommandManuscriptCreate{
					ManuscriptId:       manuscriptId,
					ManuscriptThreadId: threadId,
					Hash:               theHash,
					CommitMsg:          manuscriptCreate.CommitMsg,
					Title:              manuscriptCreate.Title,
					AuthorId:           manuscriptCreate.AuthorId,
					JournalId:          manuscriptCreate.JournalId,
				},
			},
		},
	}, manuscriptId
}

type ManuscriptCreateNewVersion struct {
	TheManuscript        []byte
	CommitMsg            string
	Title                string
	AuthorId             []string
	PreviousManuscriptId string
	ThreadId             string
	JournalId            string
}

func GetCommandManuscriptCreateNewVersion(
	manuscriptCreateNewVersion *ManuscriptCreateNewVersion,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) (*Command, string) {
	manuscriptId := model.CreateManuscriptAddress()
	return &Command{
		InputAddresses: append(manuscriptCreateNewVersion.AuthorId,
			model.GetSettingsAddress(),
			signerId,
			manuscriptCreateNewVersion.JournalId,
			manuscriptId,
			manuscriptCreateNewVersion.PreviousManuscriptId,
			manuscriptCreateNewVersion.ThreadId),
		OutputAddresses: []string{signerId, manuscriptId, manuscriptCreateNewVersion.ThreadId},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signerId,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandManuscriptCreateNewVersion{
				CommandManuscriptCreateNewVersion: &model.CommandManuscriptCreateNewVersion{
					ManuscriptId:         manuscriptId,
					PreviousManuscriptId: manuscriptCreateNewVersion.PreviousManuscriptId,
					Hash:                 model.HashBytes(manuscriptCreateNewVersion.TheManuscript),
					CommitMsg:            manuscriptCreateNewVersion.CommitMsg,
					Title:                manuscriptCreateNewVersion.Title,
					AuthorId:             manuscriptCreateNewVersion.AuthorId,
				},
			},
		},
	}, manuscriptId
}

func (nbce *nonBootstrapCommandExecution) checkManuscriptCreate(c *model.CommandManuscriptCreate) (*updater, error) {
	expectedPrice := nbce.unmarshalledState.settings.PriceList.PriceAuthorSubmitNewManuscript
	if nbce.price != expectedPrice {
		return nil, formatPriceError("PriceAuthorSubmitNewManuscript", expectedPrice)
	}
	if err := checkSanityManuscriptCreate(c); err != nil {
		return nil, err
	}
	err := nbce.readAndCheckAddresses(
		append(c.AuthorId, c.JournalId),
		[]string{c.ManuscriptId, c.ManuscriptThreadId},
	)
	if err != nil {
		return nil, err
	}
	isSignerAuthor := false
	for _, a := range c.AuthorId {
		if a == nbce.verifiedSignerId {
			isSignerAuthor = true
			break
		}
	}
	if !isSignerAuthor {
		return nil, errors.New("A manuscript should be submitted by one of its authors")
	}
	status := model.ManuscriptStatus_init
	if len(c.AuthorId) == 1 {
		status = model.ManuscriptStatus_new
	}
	updates := []singleUpdate{
		&singleUpdateManuscriptCreate{
			singleUpdateManuscriptCreateBase: singleUpdateManuscriptCreateBase{
				manuscriptId:       c.ManuscriptId,
				manuscriptThreadId: c.ManuscriptThreadId,
				timestamp:          nbce.timestamp,
				hash:               c.Hash,
				versionNumber:      int32(0),
				commitMsg:          c.CommitMsg,
				title:              c.Title,
				status:             status,
				journalId:          c.JournalId,
			},
		},
	}
	updates = nbce.addAuthorUpdates(c.AuthorId, updates, c.ManuscriptId)
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates:           updates,
	}, nil
}

func checkSanityManuscriptCreate(c *model.CommandManuscriptCreate) error {
	if !model.IsManuscriptAddress(c.ManuscriptId) {
		return errors.New("Not a manuscript address: " + c.ManuscriptId)
	}
	if !model.IsManuscriptThreadAddress(c.ManuscriptThreadId) {
		return errors.New("Not a manuscript thread address: " + c.ManuscriptThreadId)
	}
	if c.Hash == "" {
		return errors.New("Hash is empty")
	}
	// CommitMsg is allowed to be empty.
	if c.Title == "" {
		return errors.New("Title is empty")
	}
	for _, authorId := range c.AuthorId {
		if !model.IsPersonAddress(authorId) {
			return errors.New("Author is is not a person id: " + authorId)
		}
	}
	if !model.IsJournalAddress(c.JournalId) {
		return errors.New("JournalId is not a journal: " + c.JournalId)
	}
	return nil
}

type singleUpdateManuscriptCreateBase struct {
	manuscriptId       string
	manuscriptThreadId string
	timestamp          int64
	hash               string
	versionNumber      int32
	commitMsg          string
	title              string
	status             model.ManuscriptStatus
	journalId          string
}

func (u *singleUpdateManuscriptCreateBase) updateStateManuscript(state *unmarshalledState) {
	state.manuscripts[u.manuscriptId] = &model.StateManuscript{
		Id:            u.manuscriptId,
		CreatedOn:     u.timestamp,
		ModifiedOn:    u.timestamp,
		Hash:          u.hash,
		ThreadId:      u.manuscriptThreadId,
		VersionNumber: u.versionNumber,
		CommitMsg:     u.commitMsg,
		Title:         u.title,
		Author:        []*model.Author{},
		Status:        u.status,
		JournalId:     u.journalId,
	}
}

func (u *singleUpdateManuscriptCreateBase) issueEvent(
	eventSeq int32, transactionId string, ba BlockchainAccess) error {
	return ba.AddEvent(
		model.AlexandriaPrefix+model.EV_TYPE_MANUSCRIPT_CREATE,
		[]processor.Attribute{
			{
				Key:   model.EV_KEY_TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.EV_KEY_EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.EV_KEY_TIMESTAMP,
				Value: fmt.Sprintf("%d", u.timestamp),
			},
			{
				Key:   model.EV_KEY_MANUSCRIPT_ID,
				Value: u.manuscriptId,
			},
			{
				Key:   model.EV_KEY_MANUSCRIPT_THREAD_ID,
				Value: u.manuscriptThreadId,
			},
			{
				Key:   model.EV_KEY_MANUSCRIPT_HASH,
				Value: u.hash,
			},
			{
				Key:   model.EV_KEY_MANUSCRIPT_VERSION_NUMBER,
				Value: fmt.Sprintf("%d", u.versionNumber),
			},
			{
				Key:   model.EV_KEY_MANUSCRIPT_COMMIT_MSG,
				Value: u.commitMsg,
			},
			{
				Key:   model.EV_KEY_MANUSCRIPT_TITLE,
				Value: u.title,
			},
			{
				Key:   model.EV_KEY_MANUSCRIPT_STATUS,
				Value: model.GetManuscriptStatusString(u.status),
			},
			{
				Key:   model.EV_KEY_JOURNAL_ID,
				Value: u.journalId,
			},
		}, []byte{})
}

type singleUpdateManuscriptCreate struct {
	singleUpdateManuscriptCreateBase
}

var _ singleUpdate = new(singleUpdateManuscriptCreate)

func (u *singleUpdateManuscriptCreate) updateState(state *unmarshalledState) []string {
	u.updateStateManuscript(state)
	state.manuscriptThreads[u.manuscriptThreadId] = &model.StateManuscriptThread{
		Id:           u.manuscriptThreadId,
		ManuscriptId: []string{u.manuscriptId},
		IsReviewable: false,
	}
	return []string{u.manuscriptId, u.manuscriptThreadId}
}

type singleUpdateAuthorCreate struct {
	manuscriptId string
	authorId     string
	didSign      bool
	authorNumber int32
	timestamp    int64
}

var _ singleUpdate = new(singleUpdateAuthorCreate)

func (u *singleUpdateAuthorCreate) updateState(state *unmarshalledState) []string {
	theManuscript := state.manuscripts[u.manuscriptId]
	numExistingAuthors := int32(len(theManuscript.Author))
	if u.authorNumber != numExistingAuthors {
		panic("Internal error: sequence of author create events has been messed up")
	}
	authors := make([]*model.Author, numExistingAuthors+1)
	for i := int32(0); i < numExistingAuthors; i++ {
		authors[i] = theManuscript.Author[i]
	}
	authors[numExistingAuthors] = &model.Author{
		AuthorId:     u.authorId,
		DidSign:      u.didSign,
		AuthorNumber: numExistingAuthors,
	}
	theManuscript.Author = authors
	return []string{u.manuscriptId}
}

func (u *singleUpdateAuthorCreate) issueEvent(
	eventSeq int32, transactionId string, ba BlockchainAccess) error {
	return ba.AddEvent(
		model.AlexandriaPrefix+model.EV_TYPE_AUTHOR_CREATE,
		[]processor.Attribute{
			{
				Key:   model.EV_KEY_TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.EV_KEY_EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.EV_KEY_TIMESTAMP,
				Value: fmt.Sprintf("%d", u.timestamp),
			},
			{
				Key:   model.EV_KEY_MANUSCRIPT_ID,
				Value: u.manuscriptId,
			},
			{
				Key:   model.EV_KEY_PERSON_ID,
				Value: u.authorId,
			},
			{
				Key:   model.EV_KEY_AUTHOR_DID_SIGN,
				Value: strconv.FormatBool(u.didSign),
			},
			{
				Key:   model.EV_KEY_AUTHOR_NUMBER,
				Value: fmt.Sprintf("%d", u.authorNumber),
			},
		}, []byte{})
}

func (nbce *nonBootstrapCommandExecution) checkManuscriptCreateNewVersion(
	c *model.CommandManuscriptCreateNewVersion) (*updater, error) {
	expectedPrice := nbce.unmarshalledState.settings.PriceList.PriceAuthorSubmitNewVersion
	if nbce.price != expectedPrice {
		return nil, formatPriceError("PriceAuthorSubmitNewVersion", expectedPrice)
	}
	if err := checkSanityManuscriptCreateNewVersion(c); err != nil {
		return nil, err
	}
	isSignerAuthor := false
	for _, a := range c.AuthorId {
		if a == nbce.verifiedSignerId {
			isSignerAuthor = true
			break
		}
	}
	if !isSignerAuthor {
		return nil, errors.New("A manuscript should be submitted by one of its authors")
	}
	err := nbce.readAndCheckAddresses(
		append(c.AuthorId, c.PreviousManuscriptId),
		[]string{c.ManuscriptId})
	if err != nil {
		return nil, err
	}
	previousManuscript := nbce.unmarshalledState.manuscripts[c.PreviousManuscriptId]
	err = nbce.readAndCheckAddresses(
		[]string{previousManuscript.ThreadId, previousManuscript.JournalId},
		[]string{})
	if err != nil {
		return nil, err
	}
	manuscriptThread := nbce.unmarshalledState.manuscriptThreads[previousManuscript.ThreadId]
	if manuscriptThread.ManuscriptId[len(manuscriptThread.ManuscriptId)-1] != c.PreviousManuscriptId {
		return nil, errors.New("You can only add a manuscript to the end of its thread")
	}
	status := model.ManuscriptStatus_init
	if len(c.AuthorId) == 1 {
		status = model.ManuscriptStatus_new
	}
	if status == model.ManuscriptStatus_new && manuscriptThread.IsReviewable {
		status = model.ManuscriptStatus_reviewable
	}
	versionNumber := int32(len(manuscriptThread.ManuscriptId))
	updates := []singleUpdate{
		&singleUpdateManuscriptCreateNewVersion{
			singleUpdateManuscriptCreateBase: singleUpdateManuscriptCreateBase{
				manuscriptId:       c.ManuscriptId,
				manuscriptThreadId: manuscriptThread.Id,
				timestamp:          nbce.timestamp,
				hash:               c.Hash,
				versionNumber:      versionNumber,
				commitMsg:          c.CommitMsg,
				title:              c.Title,
				status:             status,
				journalId:          previousManuscript.JournalId,
			},
		},
	}
	updates = nbce.addAuthorUpdates(c.AuthorId, updates, c.ManuscriptId)
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates:           updates,
	}, nil
}

func checkSanityManuscriptCreateNewVersion(c *model.CommandManuscriptCreateNewVersion) error {
	if !model.IsManuscriptAddress(c.ManuscriptId) {
		return errors.New("ManuscriptId is not a manuscript id: " + c.ManuscriptId)
	}
	if !model.IsManuscriptAddress(c.PreviousManuscriptId) {
		return errors.New("PreviousManuscriptId is not a manuscript id: " + c.PreviousManuscriptId)
	}
	if c.Hash == "" {
		return errors.New("Hash should not be empty")
	}
	if c.CommitMsg == "" {
		return errors.New("For versions after the first, the commit message is mandatory")
	}
	if c.Title == "" {
		return errors.New("Title is mandatory")
	}
	for _, authorId := range c.AuthorId {
		if !model.IsPersonAddress(authorId) {
			return errors.New("Author is not a person: " + authorId)
		}
	}
	return nil
}

type singleUpdateManuscriptCreateNewVersion struct {
	singleUpdateManuscriptCreateBase
}

var _ singleUpdate = new(singleUpdateManuscriptCreateNewVersion)

func (u *singleUpdateManuscriptCreateNewVersion) updateState(state *unmarshalledState) []string {
	u.updateStateManuscript(state)
	state.manuscriptThreads[u.manuscriptThreadId].ManuscriptId = util.EconomicStringSliceAppend(
		state.manuscriptThreads[u.manuscriptThreadId].ManuscriptId, u.manuscriptId)
	return []string{u.manuscriptId, u.manuscriptThreadId}
}
