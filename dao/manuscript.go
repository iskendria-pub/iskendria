package dao

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/jmoiron/sqlx"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"strconv"
	"strings"
)

func createManuscriptCreateEvent(ev *events_pb2.Event) (event, error) {
	dm := &dataManipulationManuscriptCreate{}
	result := &dataManipulationEvent{
		dataManipulation: dm,
	}
	var err error
	var i64 int64
	for _, a := range ev.Attributes {
		switch a.Key {
		case model.EV_KEY_TRANSACTION_ID:
			result.transactionId = a.Value
		case model.EV_KEY_EVENT_SEQ:
			i64, err = strconv.ParseInt(a.Value, 10, 32)
			result.eventSeq = int32(i64)
		case model.EV_KEY_TIMESTAMP:
			i64, err = strconv.ParseInt(a.Value, 10, 64)
			dm.timestamp = i64
		case model.EV_KEY_MANUSCRIPT_ID:
			dm.id = a.Value
		case model.EV_KEY_MANUSCRIPT_THREAD_ID:
			dm.threadId = a.Value
		case model.EV_KEY_MANUSCRIPT_HASH:
			dm.hash = a.Value
		case model.EV_KEY_MANUSCRIPT_VERSION_NUMBER:
			i64, err = strconv.ParseInt(a.Value, 10, 32)
			dm.versionNumber = int32(i64)
		case model.EV_KEY_MANUSCRIPT_COMMIT_MSG:
			dm.commitMsg = a.Value
		case model.EV_KEY_MANUSCRIPT_TITLE:
			dm.title = a.Value
		case model.EV_KEY_MANUSCRIPT_STATUS:
			dm.status = a.Value
		case model.EV_KEY_JOURNAL_ID:
			dm.journalid = a.Value
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationManuscriptCreate struct {
	id            string
	timestamp     int64
	hash          string
	threadId      string
	versionNumber int32
	commitMsg     string
	title         string
	status        string
	journalid     string
}

var _ dataManipulation = new(dataManipulationManuscriptCreate)

func (dm *dataManipulationManuscriptCreate) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec(fmt.Sprintf("INSERT INTO manuscript VALUES (%s)", GetPlaceHolders(14)),
		dm.id,
		dm.timestamp,
		dm.timestamp,
		dm.hash,
		dm.threadId,
		dm.versionNumber,
		dm.commitMsg,
		dm.title,
		dm.status,
		dm.journalid,
		"",
		"",
		"",
		false)
	return err
}

func createAuthorCreateEvent(ev *events_pb2.Event) (event, error) {
	dm := &dataManipulationAuthorCreate{}
	result := &dataManipulationEvent{
		dataManipulation: dm,
	}
	var err error
	var i64 int64
	var b bool
	for _, a := range ev.Attributes {
		switch a.Key {
		case model.EV_KEY_TRANSACTION_ID:
			result.transactionId = a.Value
		case model.EV_KEY_EVENT_SEQ:
			i64, err = strconv.ParseInt(a.Value, 10, 32)
			result.eventSeq = int32(i64)
		case model.EV_KEY_MANUSCRIPT_ID:
			dm.manuscriptId = a.Value
		case model.EV_KEY_PERSON_ID:
			dm.personId = a.Value
		case model.EV_KEY_AUTHOR_DID_SIGN:
			b, err = strconv.ParseBool(a.Value)
			dm.didSign = b
		case model.EV_KEY_AUTHOR_NUMBER:
			i64, err = strconv.ParseInt(a.Value, 10, 32)
			dm.authorNumber = int32(i64)
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationAuthorCreate struct {
	manuscriptId string
	personId     string
	didSign      bool
	authorNumber int32
}

var _ dataManipulation = new(dataManipulationAuthorCreate)

func (dm *dataManipulationAuthorCreate) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec(fmt.Sprintf("INSERT INTO author VALUES (%s)", GetPlaceHolders(4)),
		dm.manuscriptId,
		dm.personId,
		dm.didSign,
		dm.authorNumber)
	return err
}

func createManuscriptUpdateEvent(ev *events_pb2.Event) (event, error) {
	dm := &dataManipulationManuscriptUpdateString{}
	result := &dataManipulationEvent{
		dataManipulation: dm,
	}
	var err error
	var i64 int64
	for _, a := range ev.Attributes {
		switch a.Key {
		case model.EV_KEY_TRANSACTION_ID:
			result.transactionId = a.Value
		case model.EV_KEY_EVENT_SEQ:
			i64, err = strconv.ParseInt(a.Value, 10, 32)
			result.eventSeq = int32(i64)
		case model.EV_KEY_ID:
			dm.manuscriptId = a.Value
		case
			model.EV_KEY_MANUSCRIPT_STATUS,
			model.EV_KEY_VOLUME_ID,
			model.EV_KEY_MANUSCRIPT_FIRST_PAGE,
			model.EV_KEY_MANUSCRIPT_LAST_PAGE:
			dm.field = strings.ToLower(a.Key)
			dm.newValue = a.Value
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationManuscriptUpdateString struct {
	manuscriptId string
	field        string
	newValue     string
}

var _ dataManipulation = new(dataManipulationManuscriptUpdateString)

func (dm *dataManipulationManuscriptUpdateString) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec(fmt.Sprintf("UPDATE manuscript SET %s = ? WHERE id = ?",
		dm.field), dm.newValue, dm.manuscriptId)
	return err
}

func createAuthorUpdateEvent(ev *events_pb2.Event) (event, error) {
	dm := &dataManipulationAuthorSign{}
	result := &dataManipulationEvent{
		dataManipulation: dm,
	}
	var b bool
	var i64 int64
	var err error
	for _, a := range ev.Attributes {
		switch a.Key {
		case model.EV_KEY_TRANSACTION_ID:
			result.transactionId = a.Value
		case model.EV_KEY_EVENT_SEQ:
			i64, err = strconv.ParseInt(a.Value, 10, 32)
			result.eventSeq = int32(i64)
		case model.EV_KEY_MANUSCRIPT_ID:
			dm.manuscriptId = a.Value
		case model.EV_KEY_PERSON_ID:
			dm.personId = a.Value
		case model.EV_KEY_AUTHOR_DID_SIGN:
			b, err = strconv.ParseBool(a.Value)
			if err == nil && b == false {
				err = errors.New("author.didsign cannot be cleared")
			}
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationAuthorSign struct {
	manuscriptId string
	personId     string
}

var _ dataManipulation = new(dataManipulationAuthorSign)

func (dm *dataManipulationAuthorSign) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec("UPDATE author SET didsign = ? WHERE manuscriptid = ? AND personid = ?",
		true, dm.manuscriptId, dm.personId)
	return err
}

func createManuscriptThreadUpdateEvent(ev *events_pb2.Event) (event, error) {
	dm := &dataManipulationManuscriptThreadUpdate{}
	result := &dataManipulationEvent{
		dataManipulation: dm,
	}
	var i64 int64
	var err error
	for _, a := range ev.Attributes {
		switch a.Key {
		case model.EV_KEY_TRANSACTION_ID:
			result.transactionId = a.Value
		case model.EV_KEY_EVENT_SEQ:
			i64, err = strconv.ParseInt(a.Value, 10, 32)
			result.eventSeq = int32(i64)
		case model.EV_KEY_MANUSCRIPT_THREAD_ID:
			dm.threadId = a.Value
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationManuscriptThreadUpdate struct {
	threadId string
}

var _ dataManipulation = new(dataManipulationManuscriptThreadUpdate)

func (dm *dataManipulationManuscriptThreadUpdate) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec(
		"UPDATE manuscript SET isreviewable = ? WHERE threadid = ?",
		true, dm.threadId)
	return err
}

func createReviewCreateEvent(ev *events_pb2.Event) (event, error) {
	dm := &dataManipulationReviewCreate{}
	result := &dataManipulationEvent{
		dataManipulation: dm,
	}
	var err error
	var i64 int64
	for _, a := range ev.Attributes {
		switch a.Key {
		case model.EV_KEY_TRANSACTION_ID:
			result.transactionId = a.Value
		case model.EV_KEY_EVENT_SEQ:
			i64, err = strconv.ParseInt(a.Value, 10, 32)
			result.eventSeq = int32(i64)
		case model.EV_KEY_TIMESTAMP:
			i64, err = strconv.ParseInt(a.Value, 10, 64)
			dm.timestamp = i64
		case model.EV_KEY_ID:
			dm.id = a.Value
		case model.EV_KEY_MANUSCRIPT_ID:
			dm.manuscriptId = a.Value
		case model.EV_KEY_REVIEW_AUTHOR_ID:
			dm.reviewAuthorId = a.Value
		case model.EV_KEY_REVIEW_HASH:
			dm.hash = a.Value
		case model.EV_KEY_REVIEW_JUDGEMENT:
			dm.judgement = a.Value
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationReviewCreate struct {
	id             string
	timestamp      int64
	manuscriptId   string
	reviewAuthorId string
	hash           string
	judgement      string
}

var _ dataManipulation = new(dataManipulationReviewCreate)

func (dm *dataManipulationReviewCreate) apply(tx *sqlx.Tx) error {
	query := fmt.Sprintf("INSERT INTO review VALUES (%s)", GetPlaceHolders(7))
	_, err := tx.Exec(query,
		dm.id,
		dm.timestamp,
		dm.manuscriptId,
		dm.reviewAuthorId,
		dm.hash,
		dm.judgement,
		false)
	return err
}

func createReviewUseByEditorEvent(ev *events_pb2.Event) (event, error) {
	dm := &dataManipulationReviewUseByEditor{}
	result := &dataManipulationEvent{
		dataManipulation: dm,
	}
	var i64 int64
	var err error
	for _, a := range ev.Attributes {
		switch a.Key {
		case model.EV_KEY_TRANSACTION_ID:
			result.transactionId = a.Value
		case model.EV_KEY_EVENT_SEQ:
			i64, err = strconv.ParseInt(a.Value, 10, 32)
			result.eventSeq = int32(i64)
		case model.EV_KEY_ID:
			dm.reviewId = a.Value
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationReviewUseByEditor struct {
	reviewId string
}

var _ dataManipulation = new(dataManipulationReviewUseByEditor)

func (dm *dataManipulationReviewUseByEditor) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec("UPDATE review SET isusedbyeditor = ? WHERE id = ?",
		true, dm.reviewId)
	return err
}

func GetManuscript(manuscriptId string) (*Manuscript, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Commit() }()
	return getManuscriptFromTransaction(tx, manuscriptId)
}

func getManuscriptFromTransaction(tx *sqlx.Tx, manuscriptId string) (*Manuscript, error) {
	combinations := &[]ManuscriptAuthorCombination{}
	err := tx.Select(combinations, getGetManuscriptQuery(), manuscriptId)
	if err != nil {
		return nil, err
	}
	if len(*combinations) == 0 {
		return nil, errors.New("Manuscript not found: " + manuscriptId)
	}
	return combinationsToManuscript(combinations), nil
}

type Manuscript struct {
	Id            string
	CreatedOn     int64
	ModifiedOn    int64
	Hash          string
	ThreadId      string
	VersionNumber int32
	CommitMsg     string
	Title         string
	Status        string
	JournalId     string
	VolumeId      string
	FirstPage     string
	LastPage      string
	IsReviewable  bool
	Authors       []*Author
}

type Author struct {
	ManuscriptId string
	PersonId     string
	DidSign      bool
	AuthorNumber int32
	PersonName   string
}

type ManuscriptAuthorCombination struct {
	Id            string
	CreatedOn     int64
	ModifiedOn    int64
	Hash          string
	ThreadId      string
	VersionNumber int32
	CommitMsg     string
	Title         string
	Status        string
	JournalId     string
	VolumeId      string
	FirstPage     string
	LastPage      string
	IsReviewable  bool
	PersonId      string
	DidSign       bool
	AuthorNumber  int32
	PersonName    string
}

func getGetManuscriptQuery() string {
	return `
SELECT
	manuscript.id,
	manuscript.createdon,
	manuscript.modifiedon,
	manuscript.hash,
	manuscript.threadid,
	manuscript.versionnumber,
	manuscript.commitmsg,
	manuscript.title,
	manuscript.status,
	manuscript.journalid,
	manuscript.volumeid,
	manuscript.firstpage,
	manuscript.lastpage,
	manuscript.isreviewable,
	author.personid,
	author.didsign,
	author.authornumber,
    person.name AS personname
FROM manuscript, author, person
WHERE manuscript.id = author.manuscriptid
  AND person.id = author.personid
  AND manuscript.id = ?
ORDER BY author.authornumber
`
}

func combinationsToManuscript(combinations *[]ManuscriptAuthorCombination) *Manuscript {
	authors := make([]*Author, len(*combinations))
	result := &Manuscript{
		Authors: authors,
	}
	for i, c := range *combinations {
		result.Id = c.Id
		result.CreatedOn = c.CreatedOn
		result.ModifiedOn = c.ModifiedOn
		result.Hash = c.Hash
		result.ThreadId = c.ThreadId
		result.VersionNumber = c.VersionNumber
		result.CommitMsg = c.CommitMsg
		result.Title = c.Title
		result.Status = c.Status
		result.JournalId = c.JournalId
		result.VolumeId = c.VolumeId
		result.FirstPage = c.FirstPage
		result.LastPage = c.LastPage
		result.IsReviewable = c.IsReviewable
		result.Authors[i] = &Author{
			ManuscriptId: c.Id,
			PersonId:     c.PersonId,
			DidSign:      c.DidSign,
			AuthorNumber: c.AuthorNumber,
			PersonName:   c.PersonName,
		}
	}
	return result
}

func VerifyManuscript(manuscriptId string, data []byte) error {
	manuscript, err := GetManuscript(manuscriptId)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not get manuscript with manuscriptId %s, error is %s",
			manuscriptId, err.Error()))
	}
	hashOfData := model.HashBytes(data)
	if manuscript.Hash != hashOfData {
		return errors.New("Verification failed")
	}
	return nil
}

func GetReferenceThread(threadId string) ([]ReferenceThreadItem, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Commit() }()
	rawResult := make([]ReferenceThreadItem, 0)
	err = tx.Select(&rawResult, getQueryReferenceThread(), threadId)
	if err != nil {
		return nil, err
	}
	return rawResult, nil
}

type ReferenceThreadItem struct {
	Id     string
	Status string
}

func getQueryReferenceThread() string {
	return `
SELECT
  id,
  status
FROM manuscript
WHERE threadid = ?
ORDER BY versionnumber
`
}

func GetReview(reviewId string) (*Review, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Commit() }()
	return getReviewFromTransaction(tx, reviewId)
}

func getReviewFromTransaction(tx *sqlx.Tx, reviewId string) (*Review, error) {
	result := Review{}
	err := tx.Get(&result, "SELECT * FROM review WHERE id = ?", reviewId)
	return &result, err
}

type Review struct {
	Id             string
	CreatedOn      int64
	ManuscriptId   string
	ReviewAuthorId string
	Hash           string
	Judgement      string
	IsUsedByEditor bool
}

func GetManuscriptView(manuscriptId string) (*ManuscriptView, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Commit() }()
	result := new(ManuscriptView)
	result.Manuscript, err = getManuscriptFromTransaction(tx, manuscriptId)
	if err != nil {
		return nil, err
	}
	result.Journal, err = getJournalFromTransaction(tx, result.Manuscript.JournalId)
	if err != nil {
		return nil, err
	}
	result.Reviews, err = getManuscriptReviewsFromTransaction(tx, manuscriptId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getManuscriptReviewsFromTransaction(tx *sqlx.Tx, manuscriptId string) ([]*ExtendedReview, error) {
	reviews := &[]ExtendedReview{}
	err := tx.Select(reviews, getExtendedReviewQuery(), manuscriptId)
	if err != nil {
		return nil, err
	}
	result := []*ExtendedReview{}
	for _, r := range *reviews {
		result = append(result, &r)
	}
	return result, err
}

type ManuscriptView struct {
	Manuscript *Manuscript
	Journal    *Journal
	Reviews    []*ExtendedReview
}

type ExtendedReview struct {
	Id             string
	CreatedOn      int64
	ManuscriptId   string
	ReviewAuthorId string
	Hash           string
	Judgement      string
	IsUsedByEditor bool
	PersonId       string
	PersonIsSigned bool
	PersonName     string
}

func getExtendedReviewQuery() string {
	return `
SELECT
	review.id,
    review.createdon,
    review.manuscriptid,
    review.reviewauthorid,
    review.hash,
    review.judgement,
    review.isusedbyeditor,
    person.id AS personid,
    person.issigned AS personissigned,
    person.name AS personname
FROM review, person
WHERE review.reviewauthorid = person.id
AND review.manuscriptid = ?
ORDER BY review.createdon
`
}

func GetReviewDetailsView(reviewId string) (*ReviewView, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Commit() }()
	result := new(ReviewView)
	result.Review, err = getReviewFromTransaction(tx, reviewId)
	if err != nil {
		return nil, err
	}
	result.ReviewAuthor, err = getPersonByIdFromTransaction(result.Review.ReviewAuthorId, tx)
	if err != nil {
		return nil, err
	}
	result.Manuscript, err = getManuscriptFromTransaction(tx, result.Review.ManuscriptId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type ReviewView struct {
	Review       *Review
	ReviewAuthor *Person
	Manuscript   *Manuscript
}

type ReviewEditor struct {
	PersonId       string
	PersonName     string
	PersonIsSigned bool
}

func VerifyReview(reviewId string, data []byte) error {
	tx, err := db.Beginx()
	if err != nil {
		return errors.New("Could not start database transaction")
	}
	defer func() { _ = tx.Commit() }()
	review, err := getReviewFromTransaction(tx, reviewId)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not get review with reviewId %s, error is %s",
			reviewId, err.Error()))
	}
	hashOfData := model.HashBytes(data)
	if review.Hash != hashOfData {
		return errors.New("Verification failed")
	}
	return nil
}
