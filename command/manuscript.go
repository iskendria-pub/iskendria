package command

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"strconv"
)

type ManuscriptCreate struct {
	TheManuscript []byte
	CommitMsg     string
	Title         string
	AuthorId      []string
	JournalId     string
}

func GetManuscriptCreateCommand(
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

func (nbce *nonBootstrapCommandExecution) checkManuscriptCreate(c *model.CommandManuscriptCreate) (*updater, error) {
	expectedPrice := nbce.unmarshalledState.settings.PriceList.PriceAuthorSubmitNewManuscript
	if nbce.price != expectedPrice {
		return nil, formatPriceError("PriceAuthorSubmitNewManuscript", expectedPrice)
	}
	if err := checkSanityManuscriptCreate(c); err != nil {
		return nil, err
	}
	err := nbce.readAndCheckManuscriptCreateAddresses(c)
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
	}
	for i, a := range c.AuthorId {
		didSign := (a == nbce.verifiedSignerId)
		updates = append(updates, &singleUpdateAuthorCreate{
			manuscriptId: c.ManuscriptId,
			authorId:     a,
			didSign:      didSign,
			authorNumber: int32(i),
			timestamp:    nbce.timestamp,
		})
	}
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

func (nbce *nonBootstrapCommandExecution) readAndCheckManuscriptCreateAddresses(c *model.CommandManuscriptCreate) error {
	addressesExpectedFilled := append(c.AuthorId, c.JournalId)
	addressesExpectedEmpty := []string{c.ManuscriptId, c.ManuscriptThreadId}
	toRead := append(addressesExpectedFilled, addressesExpectedEmpty...)
	readState, err := nbce.blockchainAccess.GetState(toRead)
	if err != nil {
		return err
	}
	err = nbce.unmarshalledState.add(readState, toRead)
	if err != nil {
		return err
	}
	for _, a := range addressesExpectedFilled {
		if nbce.unmarshalledState.getAddressState(a) != ADDRESS_FILLED {
			return errors.New("Address was not filled: " + a)
		}
	}
	for _, a := range addressesExpectedEmpty {
		if nbce.unmarshalledState.getAddressState(a) != ADDRESS_EMPTY {
			return errors.New("Manuscript id or manuscript thread id already in use: " + a)
		}
	}
	return nil
}

type singleUpdateManuscriptCreate struct {
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

var _ singleUpdate = new(singleUpdateManuscriptCreate)

func (u *singleUpdateManuscriptCreate) updateState(state *unmarshalledState) []string {
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
	state.manuscriptThreads[u.manuscriptThreadId] = &model.StateManuscriptThread{
		Id:           u.manuscriptThreadId,
		ManuscriptId: []string{u.manuscriptId},
		IsReviewable: false,
	}
	return []string{u.manuscriptId, u.manuscriptThreadId}
}

func (u *singleUpdateManuscriptCreate) issueEvent(
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
