package command

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"gitlab.bbinfra.net/3estack/alexandria/util"
	"sort"
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

func GetCommandManuscriptAcceptAuthorship(
	manuscript *dao.Manuscript,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	authorIds := make([]string, len(manuscript.Authors))
	for i, a := range manuscript.Authors {
		authorIds[i] = a.PersonId
	}
	return &Command{
		InputAddresses: append(
			authorIds,
			model.GetSettingsAddress(),
			manuscript.Id,
			manuscript.ThreadId,
			signerId),
		OutputAddresses: []string{manuscript.Id, manuscript.ThreadId, signerId},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signerId,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandManuscriptAcceptAuthorship{
				CommandManuscriptAcceptAuthorship: &model.CommandManuscriptAcceptAuthorship{
					ManuscriptId: manuscript.Id,
					Author:       daoManudcriptToCommandAuthors(manuscript),
				},
			},
		},
	}
}

func daoManudcriptToCommandAuthors(manuscript *dao.Manuscript) []*model.Author {
	result := make([]*model.Author, len(manuscript.Authors))
	for i, a := range manuscript.Authors {
		result[i] = &model.Author{
			AuthorId:     a.PersonId,
			DidSign:      a.DidSign,
			AuthorNumber: a.AuthorNumber,
		}
	}
	return result
}

func GetCommandManuscriptAllowReview(
	threadId string,
	daoThreadReference []dao.ReferenceThreadItem,
	journalId string,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	manuscriptIds := make([]string, len(daoThreadReference))
	for i, r := range daoThreadReference {
		manuscriptIds[i] = r.Id
	}
	return &Command{
		InputAddresses:  append(manuscriptIds, threadId, model.GetSettingsAddress(), signerId, journalId),
		OutputAddresses: append(manuscriptIds, threadId, signerId),
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signerId,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandManuscriptAllowReview{
				CommandManuscriptAllowReview: &model.CommandManuscriptAllowReview{
					ThreadId:        threadId,
					ThreadReference: daoThreadReferenceToCommandReferenceThread(daoThreadReference),
				},
			},
		},
	}
}

func daoThreadReferenceToCommandReferenceThread(
	daoReferenceThread []dao.ReferenceThreadItem) []*model.ThreadReferenceItem {
	result := make([]*model.ThreadReferenceItem, len(daoReferenceThread))
	for i, r := range daoReferenceThread {
		result[i] = &model.ThreadReferenceItem{
			ManuscriptId:     r.Id,
			ManuscriptStatus: model.GetManuscriptStatusCode(r.Status),
		}
	}
	return result
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
	status := getNewManuscriptStatus(len(c.AuthorId) == 1, false)
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

func getNewManuscriptStatus(allAuthorsWillHaveSigned, isThreadReviewable bool) model.ManuscriptStatus {
	status := model.ManuscriptStatus_init
	if allAuthorsWillHaveSigned {
		status = model.ManuscriptStatus_new
	}
	if status == model.ManuscriptStatus_new && isThreadReviewable {
		status = model.ManuscriptStatus_reviewable
	}
	return status
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
	status := getNewManuscriptStatus(len(c.AuthorId) == 1, manuscriptThread.IsReviewable)
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
	thread := state.manuscriptThreads[u.manuscriptThreadId]
	thread.ManuscriptId = util.EconomicStringSliceAppend(
		thread.ManuscriptId, u.manuscriptId)
	sort.Slice(thread.ManuscriptId, func(i, j int) bool {
		return thread.ManuscriptId[i] < thread.ManuscriptId[j]
	})
	return []string{u.manuscriptId, u.manuscriptThreadId}
}

func (nbce *nonBootstrapCommandExecution) checkManuscriptAcceptAuthorship(
	c *model.CommandManuscriptAcceptAuthorship) (*updater, error) {
	expectedPrice := nbce.unmarshalledState.settings.PriceList.PriceAuthorAcceptAuthorship
	if nbce.price != expectedPrice {
		return nil, formatPriceError("PriceAuthorAcceptAuthorship", expectedPrice)
	}
	if err := checkSanityManuscriptAcceptAuthorship(c); err != nil {
		return nil, err
	}
	err := nbce.readAndCheckAddresses(
		append(getAuthorIds(c.Author), c.ManuscriptId),
		[]string{})
	if err != nil {
		return nil, err
	}
	manuscript := nbce.unmarshalledState.manuscripts[c.ManuscriptId]
	err = nbce.readAndCheckAddresses([]string{manuscript.ThreadId}, []string{})
	if err != nil {
		return nil, err
	}
	if manuscript.Status != model.ManuscriptStatus_init {
		return nil, errors.New("All authors already accepted authorship")
	}
	err = checkAuthors(manuscript.Author, c.Author)
	if err != nil {
		return nil, err
	}
	err = nbce.checkSignerIsAuthor(c.ManuscriptId)
	if err != nil {
		return nil, err
	}
	doesAuthorUpdate, allAuthorsWillHaveSigned := getCommandManuscriptAcceptAuthorshipWork(
		c, nbce.verifiedSignerId)
	updates := []singleUpdate{}
	if doesAuthorUpdate {
		updates = append(updates, &singleUpdateAuthorUpdate{
			manuscriptId: c.ManuscriptId,
			authorId:     nbce.verifiedSignerId,
			timestamp:    nbce.timestamp,
		})
	}
	if allAuthorsWillHaveSigned {
		status := getNewManuscriptStatus(
			allAuthorsWillHaveSigned,
			nbce.unmarshalledState.manuscriptThreads[manuscript.ThreadId].IsReviewable)
		updates = append(updates, &singleUpdateManuscriptUpdateStatus{
			manuscriptId: c.ManuscriptId,
			newStatus:    status,
			timestamp:    nbce.timestamp,
		})
	}
	updates = nbce.addSingleUpdateManuscriptModificationTimeIfNeeded(updates, c.ManuscriptId)
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates:           updates,
	}, nil
}

func checkSanityManuscriptAcceptAuthorship(c *model.CommandManuscriptAcceptAuthorship) error {
	if !model.IsManuscriptAddress(c.ManuscriptId) {
		return errors.New("Not a manuscript: " + c.ManuscriptId)
	}
	for _, a := range c.Author {
		if !model.IsPersonAddress(a.AuthorId) {
			return errors.New("AuthorId is not a person: " + a.AuthorId)
		}
	}
	return nil
}

func getAuthorIds(authors []*model.Author) []string {
	result := make([]string, len(authors))
	for i, a := range authors {
		result[i] = a.AuthorId
	}
	return result
}

func checkAuthors(expectedAuthors, actualAuthors []*model.Author) error {
	if len(expectedAuthors) != len(actualAuthors) {
		return errors.New(fmt.Sprintf("Number of authors mismatch, expected %d but got %d",
			len(expectedAuthors), len(actualAuthors)))
	}
	for i := range expectedAuthors {
		expected := expectedAuthors[i]
		actual := actualAuthors[i]
		if expected.AuthorId != actual.AuthorId {
			return errors.New(fmt.Sprintf("Author id mismatch for author #%d, expected %s but got %s",
				i+1, expected.AuthorId, actual.AuthorId))
		}
		if expected.AuthorNumber != int32(i) {
			return errors.New("Blockchain not consistent, author number does not agree with array index")
		}
		if actual.AuthorNumber != int32(i) {
			return errors.New("AuthorNumber mismatch")
		}
		if expected.DidSign != actual.DidSign {
			return errors.New(fmt.Sprintf("DidSign mismatch for author #%d, expected %v but got %v",
				i+1, expected.DidSign, actual.DidSign))
		}
	}
	return nil
}

func getCommandManuscriptAcceptAuthorshipWork(
	c *model.CommandManuscriptAcceptAuthorship,
	signerId string) (doesAuthorUpdate, allAuthorsWillHaveSigned bool) {
	allAuthorsWillHaveSigned = true
	doesAuthorUpdate = false
	for _, a := range c.Author {
		if a.AuthorId == signerId {
			if !a.DidSign {
				doesAuthorUpdate = true
			}
		} else {
			if !a.DidSign {
				allAuthorsWillHaveSigned = false
			}
		}
	}
	return
}

func (nbce *nonBootstrapCommandExecution) checkSignerIsAuthor(manuscriptId string) error {
	for _, a := range nbce.unmarshalledState.manuscripts[manuscriptId].Author {
		if a.AuthorId == nbce.verifiedSignerId {
			return nil
		}
	}
	return errors.New("You are not an author of manuscript " + manuscriptId)
}

type singleUpdateManuscriptUpdateStatus struct {
	manuscriptId string
	newStatus    model.ManuscriptStatus
	timestamp    int64
}

var _ singleUpdate = new(singleUpdateManuscriptUpdateStatus)

func (u *singleUpdateManuscriptUpdateStatus) updateState(state *unmarshalledState) (writtenAddresses []string) {
	state.manuscripts[u.manuscriptId].Status = u.newStatus
	return []string{u.manuscriptId}
}

func (u *singleUpdateManuscriptUpdateStatus) issueEvent(
	eventSeq int32, transactionId string, ba BlockchainAccess) error {
	return ba.AddEvent(
		model.AlexandriaPrefix+model.EV_TYPE_MANUSCRIPT_UPDATE,
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
				Key:   model.EV_KEY_ID,
				Value: u.manuscriptId,
			},
			{
				Key:   model.EV_KEY_MANUSCRIPT_STATUS,
				Value: model.GetManuscriptStatusString(u.newStatus),
			},
		}, []byte{})
}

type singleUpdateAuthorUpdate struct {
	manuscriptId string
	authorId     string
	timestamp    int64
}

var _ singleUpdate = new(singleUpdateAuthorUpdate)

func (u *singleUpdateAuthorUpdate) updateState(state *unmarshalledState) (writtenAddresses []string) {
	manuscript := state.manuscripts[u.manuscriptId]
	for _, a := range manuscript.Author {
		if a.AuthorId == u.authorId {
			a.DidSign = true
		}
	}
	return []string{u.manuscriptId}
}

func (u *singleUpdateAuthorUpdate) issueEvent(
	eventSeq int32, transactionId string, ba BlockchainAccess) error {
	return ba.AddEvent(
		model.AlexandriaPrefix+model.EV_TYPE_AUTHOR_UPDATE,
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
				Value: strconv.FormatBool(true),
			},
		}, []byte{})
}

func (nbce *nonBootstrapCommandExecution) checkManuscriptAllowReview(c *model.CommandManuscriptAllowReview) (
	*updater, error) {
	expectedPrice := nbce.unmarshalledState.settings.PriceList.PriceEditorAllowManuscriptReview
	if nbce.price != expectedPrice {
		return nil, formatPriceError("PriceEditorAllowManuscriptReview", expectedPrice)
	}
	if err := checkSanityManuscriptAllowReview(c); err != nil {
		return nil, err
	}
	err := nbce.readAndCheckAddresses(
		append(getManuscriptIds(c.ThreadReference), c.ThreadId),
		[]string{})
	if err != nil {
		return nil, err
	}
	err = nbce.checkThreadReference(c.ThreadId, c.ThreadReference)
	if err != nil {
		return nil, err
	}
	err = nbce.checkAllowReviewSignerIsJournalEditor(c)
	if err != nil {
		return nil, err
	}
	updates := []singleUpdate{
		&singleUpdateManuscriptThreadAllowReview{
			threadId:  c.ThreadId,
			timestamp: nbce.timestamp,
		},
	}
	for _, threadReferenceItem := range c.ThreadReference {
		allAuthorsSigned := nbce.IsAllAuthorsOfThreadReferenceItemSigned(threadReferenceItem)
		newStatus := getNewManuscriptStatus(allAuthorsSigned, true)
		if newStatus != threadReferenceItem.ManuscriptStatus {
			updates = append(updates, &singleUpdateManuscriptUpdateStatus{
				manuscriptId: threadReferenceItem.ManuscriptId,
				newStatus:    newStatus,
				timestamp:    nbce.timestamp,
			},
				&singleUpdateManuscriptModificationTime{
					id:        threadReferenceItem.ManuscriptId,
					timestamp: nbce.timestamp,
				})
		}
	}
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates:           updates,
	}, nil
}

func checkSanityManuscriptAllowReview(c *model.CommandManuscriptAllowReview) error {
	if !model.IsManuscriptThreadAddress(c.ThreadId) {
		return errors.New("Not a manuscript thread: " + c.ThreadId)
	}
	if len(c.ThreadReference) == 0 {
		return errors.New("Thread without manuscirpts: " + c.ThreadId)
	}
	for _, item := range c.ThreadReference {
		if !model.IsManuscriptAddress(item.ManuscriptId) {
			return errors.New("Not a manuscript: " + item.ManuscriptId)
		}
		if item.ManuscriptStatus < model.ManuscriptStatus(model.MinManuscriptStatus) ||
			item.ManuscriptStatus > model.ManuscriptStatus(model.MaxManuscriptStatus) {
			return errors.New(fmt.Sprintf("ManuscriptStatus out of range: %d", item.ManuscriptStatus))
		}
	}
	return nil
}

func getManuscriptIds(referenceThread []*model.ThreadReferenceItem) []string {
	result := make([]string, len(referenceThread))
	for i, r := range referenceThread {
		result[i] = r.ManuscriptId
	}
	return result
}

func (nbce *nonBootstrapCommandExecution) checkThreadReference(
	threadId string, r []*model.ThreadReferenceItem) error {
	blockchainManuscriptIds := nbce.unmarshalledState.manuscriptThreads[threadId].ManuscriptId
	if len(r) != len(blockchainManuscriptIds) {
		return errors.New("The number of reference manuscript in the command does not match. " +
			fmt.Sprintf("Expected %d, got %d",
				len(blockchainManuscriptIds), len(r)))
	}
	for i := 0; i < len(blockchainManuscriptIds); i++ {
		if r[i].ManuscriptId != blockchainManuscriptIds[i] {
			return errors.New(fmt.Sprintf("Manuscript id %d does not match. Expected %s, got %s",
				i+1, blockchainManuscriptIds[i], r[i].ManuscriptId))
		}
		expectedStatus := nbce.unmarshalledState.manuscripts[blockchainManuscriptIds[i]].Status
		if r[i].ManuscriptStatus != expectedStatus {
			return errors.New(fmt.Sprintf("Status of manuscript %d does not match. Expected %s, got %s",
				i+1,
				model.GetManuscriptStatusString(expectedStatus),
				model.GetManuscriptStatusString(r[i].ManuscriptStatus)))
		}
	}
	return nil
}

func (nbce *nonBootstrapCommandExecution) checkAllowReviewSignerIsJournalEditor(
	c *model.CommandManuscriptAllowReview) error {
	journalId := nbce.unmarshalledState.manuscripts[c.ThreadReference[0].ManuscriptId].JournalId
	err := nbce.readAndCheckAddresses([]string{journalId}, []string{})
	if err != nil {
		return err
	}
	journal := nbce.unmarshalledState.journals[journalId]
	isSignerJournalEditor := false
	for _, e := range journal.EditorInfo {
		if e.EditorId == nbce.verifiedSignerId {
			isSignerJournalEditor = true
		}
	}
	if !isSignerJournalEditor {
		return errors.New("You are not editor of journal " + journalId)
	}
	return nil
}

func (nbce *nonBootstrapCommandExecution) IsAllAuthorsOfThreadReferenceItemSigned(manuscript *model.ThreadReferenceItem) bool {
	allAuthorsSigned := true
	for _, a := range nbce.unmarshalledState.manuscripts[manuscript.ManuscriptId].Author {
		if !a.DidSign {
			allAuthorsSigned = false
			break
		}
	}
	return allAuthorsSigned
}

type singleUpdateManuscriptThreadAllowReview struct {
	threadId  string
	timestamp int64
}

var _ singleUpdate = new(singleUpdateManuscriptThreadAllowReview)

func (u *singleUpdateManuscriptThreadAllowReview) updateState(state *unmarshalledState) (writtenAddresses []string) {
	state.manuscriptThreads[u.threadId].IsReviewable = true
	return []string{u.threadId}
}

func (u *singleUpdateManuscriptThreadAllowReview) issueEvent(
	eventSeq int32, transactionId string, ba BlockchainAccess) error {
	return ba.AddEvent(
		model.AlexandriaPrefix+model.EV_TYPE_MANUSCRIPT_THREAD_UPDATE,
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
				Key:   model.EV_KEY_MANUSCRIPT_THREAD_ID,
				Value: u.threadId,
			},
		}, []byte{})
}
