package dao

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/jmoiron/sqlx"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"sort"
	"strconv"
	"strings"
)

func createJournalCreateEvent(ev *events_pb2.Event) (event, error) {
	dm := &dataManipulationJournalCreate{}
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
			dm.journalId = a.Value
		case model.EV_KEY_TIMESTAMP:
			i64, err = strconv.ParseInt(a.Value, 10, 64)
			dm.timestamp = i64
		case model.EV_KEY_JOURNAL_TITLE:
			dm.title = a.Value
		case model.EV_KEY_JOURNAL_DESCRIPTION_HASH:
			dm.descriptionHash = a.Value
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationJournalCreate struct {
	journalId       string
	timestamp       int64
	title           string
	descriptionHash string
}

var _ dataManipulation = new(dataManipulationJournalCreate)

func (dm *dataManipulationJournalCreate) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec(fmt.Sprintf("INSERT INTO journal VALUES (%s)", GetPlaceHolders(6)),
		dm.journalId, dm.timestamp, dm.timestamp, dm.title, false, dm.descriptionHash)
	return err
}

func createEditorCreateEvent(ev *events_pb2.Event) (event, error) {
	dm := &dataManipulationEditorCreate{}
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
		case model.EV_KEY_JOURNAL_ID:
			dm.journalId = a.Value
		case model.EV_KEY_EDITOR_ID:
			dm.personId = a.Value
		case model.EV_KEY_EDITOR_STATE:
			dm.editorState = a.Value
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationEditorCreate struct {
	journalId   string
	personId    string
	editorState string
}

var _ dataManipulation = new(dataManipulationEditorCreate)

func (dm *dataManipulationEditorCreate) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec(fmt.Sprintf("INSERT INTO editor VALUES (%s)", GetPlaceHolders(3)),
		dm.journalId, dm.personId, dm.editorState)
	return err
}

func createJournalUpdateEvent(ev *events_pb2.Event) (event, error) {
	dmProperties := &dataManipulationJournalUpdateProperties{}
	dmAuthorization := &dataManipulationJournalUpdateAuthorization{}
	result := &dataManipulationEvent{}
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
		case model.EV_KEY_ID:
			dmProperties.id = a.Value
			dmAuthorization.id = a.Value
		case model.EV_KEY_JOURNAL_TITLE, model.EV_KEY_JOURNAL_DESCRIPTION_HASH:
			result.dataManipulation = dmProperties
			dmProperties.field = a.Key
			dmProperties.newValue = a.Value
		case model.EV_KEY_JOURNAL_IS_SIGNED:
			b, err = strconv.ParseBool(a.Value)
			result.dataManipulation = dmAuthorization
			dmAuthorization.newIsSigned = b
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationJournalUpdateProperties struct {
	id       string
	field    string
	newValue string
}

var _ dataManipulation = new(dataManipulationJournalUpdateProperties)

func (dm *dataManipulationJournalUpdateProperties) apply(tx *sqlx.Tx) error {
	query := fmt.Sprintf("UPDATE journal SET %s = \"%s\" WHERE journalId = \"%s\"",
		dm.field, dm.newValue, dm.id)
	_, err := tx.Exec(query)
	return err
}

type dataManipulationJournalUpdateAuthorization struct {
	id          string
	newIsSigned bool
}

var _ dataManipulation = new(dataManipulationJournalUpdateAuthorization)

func (dm *dataManipulationJournalUpdateAuthorization) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec("UPDATE journal SET issigned = ? WHERE journalId = ?",
		dm.newIsSigned, dm.id)
	return err
}

func createEditorDeleteEvent(ev *events_pb2.Event) (event, error) {
	dm := &dataManipulationEditorDelete{}
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
		case model.EV_KEY_JOURNAL_ID:
			dm.journalId = a.Value
		case model.EV_KEY_EDITOR_ID:
			dm.personId = a.Value
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationEditorDelete struct {
	journalId string
	personId  string
}

var _ dataManipulation = new(dataManipulationEditorDelete)

func (dm *dataManipulationEditorDelete) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec("DELETE FROM editor WHERE journalid = ? AND personid = ?",
		dm.journalId, dm.personId)
	return err
}

func createEditorUpdateEvent(ev *events_pb2.Event) (event, error) {
	dm := &dataManipulationEditorUpdate{}
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
		case model.EV_KEY_JOURNAL_ID:
			dm.journalId = a.Value
		case model.EV_KEY_EDITOR_ID:
			dm.personId = a.Value
		case model.EV_KEY_EDITOR_STATE:
			dm.newEditorState = a.Value
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationEditorUpdate struct {
	journalId      string
	personId       string
	newEditorState string
}

var _ dataManipulation = new(dataManipulationEditorUpdate)

func (dm *dataManipulationEditorUpdate) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec("UPDATE editor SET editorstate = ? WHERE journalid = ? AND personid = ?",
		dm.newEditorState, dm.journalId, dm.personId)
	return err
}

func GetAllJournals() ([]*Journal, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Commit() }()
	result, err := getAllJournalsWithEditors(tx)
	if err != nil {
		return nil, err
	}
	journalsWithoutEditors, err := getAllJournalsWithoutEditors(tx)
	if err != nil {
		return nil, err
	}
	result = append(result, journalsWithoutEditors...)
	sortJournals(result)
	return result, nil
}

func getAllJournalsWithEditors(tx *sqlx.Tx) ([]*Journal, error) {
	journalEditorCombinations := []JournalEditorCombination{}
	err := tx.Select(&journalEditorCombinations, getJournalEditorCombinationsQuery())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not get journal editor combinations: %s", err.Error()))
	}
	result := make([]*Journal, 0)
	var currentJournal *Journal
	for _, jec := range journalEditorCombinations {
		if currentJournal == nil || currentJournal.JournalId != jec.JournalId {
			currentJournal = journalEditorCombinationToJournal(&jec)
			result = append(result, currentJournal)
		} else {
			currentJournal.AcceptedEditors = append(currentJournal.AcceptedEditors, &Editor{
				PersonId:   jec.PersonId,
				PersonName: jec.PersonName,
			})
		}
	}
	return result, nil
}

type Journal struct {
	JournalId       string
	CreatedOn       int64
	ModifiedOn      int64
	Title           string
	IsSigned        bool
	Descriptionhash string
	AcceptedEditors []*Editor
}

type Editor struct {
	PersonId   string
	PersonName string
}

type JournalEditorCombination struct {
	JournalId       string
	CreatedOn       int64
	ModifiedOn      int64
	Title           string
	IsSigned        bool
	Descriptionhash string
	PersonId        string
	PersonName      string
}

func getJournalEditorCombinationsQuery() string {
	return strings.TrimSpace(fmt.Sprintf(`
SELECT
  journal.journalid,
  journal.createdon,
  journal.modifiedon,
  journal.title,
  journal.issigned,
  journal.descriptionhash,
  editor.personid,
  person.name AS personname
FROM journal, editor, person
WHERE editor.journalid = journal.journalid
  AND editor.editorState = "%s"
  AND person.id = editor.personid
ORDER BY journal.title, journal.journalId, person.name, editor.personId
`, model.GetEditorStateString(model.EditorState_editorAccepted)))
}

func journalEditorCombinationToJournal(jec *JournalEditorCombination) *Journal {
	journal := new(Journal)
	journal.JournalId = jec.JournalId
	journal.CreatedOn = jec.CreatedOn
	journal.ModifiedOn = jec.ModifiedOn
	journal.Title = jec.Title
	journal.IsSigned = jec.IsSigned
	journal.Descriptionhash = jec.Descriptionhash
	journal.AcceptedEditors = []*Editor{
		{
			PersonId:   jec.PersonId,
			PersonName: jec.PersonName,
		},
	}
	return journal
}

func getAllJournalsWithoutEditors(tx *sqlx.Tx) ([]*Journal, error) {
	journalsWithoutEditors := []JournalWithoutEditors{}
	err := tx.Select(&journalsWithoutEditors, getJournalsWithoutEditorsQuery())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not get journals without editors: %s", err.Error()))
	}
	result := make([]*Journal, 0, len(journalsWithoutEditors))
	for _, j := range journalsWithoutEditors {
		result = append(result, journalWithoutEditorsToJournal(&j))
	}
	return result, nil
}

func getJournalsWithoutEditorsQuery() string {
	return strings.TrimSpace(fmt.Sprintf(`
SELECT
  journalid,
  createdon,
  modifiedon,
  title,
  issigned,
  descriptionhash
FROM journal
WHERE journalId NOT IN (
  SELECT journalId FROM editor
  WHERE editorState = "%s"
)`, model.GetEditorStateString(model.EditorState_editorAccepted)))
}

type JournalWithoutEditors struct {
	JournalId       string
	CreatedOn       int64
	ModifiedOn      int64
	Title           string
	IsSigned        bool
	Descriptionhash string
}

func journalWithoutEditorsToJournal(jwe *JournalWithoutEditors) *Journal {
	return &Journal{
		JournalId:       jwe.JournalId,
		CreatedOn:       jwe.CreatedOn,
		ModifiedOn:      jwe.ModifiedOn,
		Title:           jwe.Title,
		IsSigned:        jwe.IsSigned,
		Descriptionhash: jwe.Descriptionhash,
	}
}

func sortJournals(journals []*Journal) {
	sort.Slice(journals, func(i, j int) bool {
		if journals[i].Title < journals[j].Title {
			return true
		}
		if journals[i].Title == journals[j].Title && journals[i].JournalId < journals[j].JournalId {
			return true
		}
		return false
	})
}
