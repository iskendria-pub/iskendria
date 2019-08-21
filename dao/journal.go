package dao

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/jmoiron/sqlx"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"log"
	"sort"
	"strconv"
	"strings"
)

func createJournalCreateEvent(ev *events_pb2.Event, logger *log.Logger) (event, error) {
	dm := &dataManipulationJournalCreate{}
	result := &dataManipulationEvent{
		dataManipulation: dm,
	}
	var err error
	var i64 int64
	for _, a := range ev.Attributes {
		logAttribute(a, logger)
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

func createEditorCreateEvent(ev *events_pb2.Event, logger *log.Logger) (event, error) {
	dm := &dataManipulationEditorCreate{}
	result := &dataManipulationEvent{
		dataManipulation: dm,
	}
	var err error
	var i64 int64
	for _, a := range ev.Attributes {
		logAttribute(a, logger)
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

func createEditorDeleteEvent(ev *events_pb2.Event, logger *log.Logger) (event, error) {
	dm := &dataManipulationEditorDelete{}
	result := &dataManipulationEvent{
		dataManipulation: dm,
	}
	var err error
	var i64 int64
	for _, a := range ev.Attributes {
		logAttribute(a, logger)
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

func createEditorUpdateEvent(ev *events_pb2.Event, logger *log.Logger) (event, error) {
	dm := &dataManipulationEditorUpdate{}
	result := &dataManipulationEvent{
		dataManipulation: dm,
	}
	var err error
	var i64 int64
	for _, a := range ev.Attributes {
		logAttribute(a, logger)
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

/*
Get all journals ordered by Title, journals with the same
title being sorted by journal id. Only editors with the
accepted state are included. This function is for use
by the portal tool, which should not show editors with
the proposed state.
*/
func GetAllJournals() ([]*Journal, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Commit() }()
	result, err := getAllJournalsHavingEditors(tx)
	if err != nil {
		return nil, err
	}
	journalsWithoutEditors, err := getAllJournalsNotHavingEditors(tx)
	if err != nil {
		return nil, err
	}
	result = append(result, journalsWithoutEditors...)
	sortJournals(result)
	return result, nil
}

func getAllJournalsHavingEditors(tx *sqlx.Tx) ([]*Journal, error) {
	return getJournalsWithEditorsForQuery(tx, getJournalEditorCombinationsQuery())
}

func getJournalsHavingSpecificEditor(tx *sqlx.Tx, editorId string) ([]*Journal, error) {
	journalIds := []JournalId{}
	err := tx.Select(&journalIds, getJournalIdsWithSpecificEditorQuery(editorId))
	if err != nil {
		return nil, err
	}
	result := make([]*Journal, 0)
	for _, journalId := range journalIds {
		journal, err := getJournalFromTransaction(tx, journalId.JournalId)
		if err != nil {
			return nil, err
		}
		result = append(result, journal)
	}
	sortJournals(result)
	return result, nil
}

type JournalId struct {
	JournalId string
}

func getJournalsWithEditorsForQuery(tx *sqlx.Tx, query string) ([]*Journal, error) {
	journalEditorCombinations := []JournalEditorCombination{}
	err := tx.Select(&journalEditorCombinations, query)
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
				PersonId:       jec.PersonId,
				PersonName:     jec.PersonName,
				PersonIsSigned: jec.PersonIsSigned,
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
	PersonId       string
	PersonName     string
	PersonIsSigned bool
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
	PersonIsSigned  bool
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
  person.name AS personname,
  person.issigned AS personissigned
FROM journal, editor, person
WHERE editor.journalid = journal.journalid
  AND editor.editorState = "%s"
  AND person.id = editor.personid
ORDER BY journal.title, journal.journalId, person.name, editor.personId
`, model.GetEditorStateString(model.EditorState_editorAccepted)))
}

func getJournalIdsWithSpecificEditorQuery(personId string) string {
	return strings.TrimSpace(fmt.Sprintf(`
SELECT DISTINCT
  journal.journalid
FROM journal, editor, person
WHERE editor.journalid = journal.journalid
  AND editor.editorState = "%s"
  AND person.id = editor.personid
  AND person.id = "%s"
ORDER BY journal.title, journal.journalId, person.name, editor.personId
`, model.GetEditorStateString(model.EditorState_editorAccepted), personId))
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
			PersonId:       jec.PersonId,
			PersonName:     jec.PersonName,
			PersonIsSigned: jec.PersonIsSigned,
		},
	}
	return journal
}

func getAllJournalsNotHavingEditors(tx *sqlx.Tx) ([]*Journal, error) {
	journalsWithoutEditors := []JournalExcludingEditors{}
	err := tx.Select(&journalsWithoutEditors, getJournalsNotHavingEditorsQuery())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not get journals without editors: %s", err.Error()))
	}
	result := make([]*Journal, 0, len(journalsWithoutEditors))
	for _, j := range journalsWithoutEditors {
		result = append(result, journalExcludingEditorsToJournal(&j))
	}
	return result, nil
}

func getJournalsNotHavingEditorsQuery() string {
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

type JournalExcludingEditors struct {
	JournalId       string
	CreatedOn       int64
	ModifiedOn      int64
	Title           string
	IsSigned        bool
	Descriptionhash string
}

func journalExcludingEditorsToJournal(jwe *JournalExcludingEditors) *Journal {
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

/*
Get a specific journal with its accepted editors. This function
is for use by the portal tool, which should not show proposed
editors.
*/
func GetJournal(journalId string) (*Journal, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, errors.New("Could not start database transaction")
	}
	defer func() { _ = tx.Commit() }()
	return getJournalFromTransaction(tx, journalId)
}

func getJournalFromTransaction(tx *sqlx.Tx, journalId string) (*Journal, error) {
	jwe, err := getJournalExcludingEditors(journalId, tx)
	if err != nil {
		return nil, err
	}
	journal := journalExcludingEditorsToJournal(jwe)
	editors := []Editor{}
	err = tx.Select(&editors, getAcceptedEditorsOfSpecificJournalQuery(journalId))
	if err != nil {
		return nil, err
	}
	journal.AcceptedEditors = make([]*Editor, 0, len(editors))
	for _, e := range editors {
		journal.AcceptedEditors = append(journal.AcceptedEditors, &Editor{
			PersonId:       e.PersonId,
			PersonName:     e.PersonName,
			PersonIsSigned: e.PersonIsSigned,
		})
	}
	return journal, nil
}

func getJournalExcludingEditors(journalId string, tx *sqlx.Tx) (*JournalExcludingEditors, error) {
	result := &JournalExcludingEditors{}
	err := tx.Get(result, "SELECT * FROM journal WHERE journalid = ?", journalId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getAcceptedEditorsOfSpecificJournalQuery(journalId string) string {
	return strings.TrimSpace(fmt.Sprintf(`
SELECT 
  editor.personid AS personid,
  person.name AS personname,
  person.issigned AS personissigned
FROM editor, person
WHERE
  editor.journalid = "%s"
  AND person.id = editor.personId
  AND editor.editorstate = "%s"
ORDER BY editor.personid
`, journalId, model.GetEditorStateString(model.EditorState_editorAccepted)))
}

func VerifyJournalDescription(journalId string, data []byte) error {
	tx, err := db.Beginx()
	if err != nil {
		return errors.New("Could not start database transaction")
	}
	defer func() { _ = tx.Commit() }()
	jwe, err := getJournalExcludingEditors(journalId, tx)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not get journal with journalId %s, error is %s",
			journalId, err.Error()))
	}
	if jwe.Descriptionhash == "" {
		if len(data) == 0 {
			return nil
		}
		return errors.New("Verification failed. There is no description on the blockchain")
	}
	hashOfData := model.HashBytes(data)
	if jwe.Descriptionhash != hashOfData {
		return errors.New("Verification failed")
	}
	return nil
}

type JournalIncludingProposedEditors struct {
	JournalId       string
	CreatedOn       int64
	ModifiedOn      int64
	Title           string
	IsSigned        bool
	Descriptionhash string
	AllEditors      []*EditorWithState
}

type EditorWithState struct {
	PersonId       string
	PersonName     string
	PersonIsSigned bool
	EditorState    string
}

/*
Get journal and include both accepted and proposed editors. This
function is for use by the client tool and the major tool.
*/
func GetJournalIncludingProposedEditors(journalId string) (*JournalIncludingProposedEditors, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, errors.New("Could not start database transaction")
	}
	defer func() { _ = tx.Commit() }()
	jwe, err := getJournalExcludingEditors(journalId, tx)
	if err != nil {
		return nil, err
	}
	journal := journalExcludingEditorsToJournalIncludingProposedEditors(jwe)
	editors := []EditorWithState{}
	err = tx.Select(&editors, getEditorsOfSpecificJournalIncludingEditorStateQuery(journalId))
	if err != nil {
		return nil, err
	}
	journal.AllEditors = make([]*EditorWithState, 0, len(editors))
	for _, e := range editors {
		journal.AllEditors = append(journal.AllEditors, &EditorWithState{
			PersonId:       e.PersonId,
			PersonName:     e.PersonName,
			PersonIsSigned: e.PersonIsSigned,
			EditorState:    e.EditorState,
		})
	}
	return journal, nil
}

func journalExcludingEditorsToJournalIncludingProposedEditors(
	jwe *JournalExcludingEditors) *JournalIncludingProposedEditors {
	return &JournalIncludingProposedEditors{
		JournalId:       jwe.JournalId,
		CreatedOn:       jwe.CreatedOn,
		ModifiedOn:      jwe.ModifiedOn,
		Title:           jwe.Title,
		IsSigned:        jwe.IsSigned,
		Descriptionhash: jwe.Descriptionhash,
	}
}

func getEditorsOfSpecificJournalIncludingEditorStateQuery(journalId string) string {
	return strings.TrimSpace(fmt.Sprintf(`
SELECT 
  editor.personid AS personid,
  editor.editorstate AS editorstate,
  person.name AS personname,
  person.issigned AS personissigned
FROM editor, person
WHERE
  editor.journalid = "%s"
  AND person.id = editor.personId
ORDER BY editor.personid
`, journalId))
}

func createVolumeCreateEvent(ev *events_pb2.Event) (event, error) {
	dm := &dataManipulationVolumeCreate{}
	result := &dataManipulationEvent{dataManipulation: dm}
	var i64 int64
	var err error
	for _, a := range ev.Attributes {
		switch a.Key {
		case model.EV_KEY_TRANSACTION_ID:
			result.transactionId = a.Value
		case model.EV_KEY_EVENT_SEQ:
			i64, err = strconv.ParseInt(a.Value, 10, 32)
			result.eventSeq = int32(i64)
		case model.EV_KEY_TIMESTAMP:
			i64, err = strconv.ParseInt(a.Value, 10, 64)
			dm.createdOn = i64
		case model.EV_KEY_ID:
			dm.volumeId = a.Value
		case model.EV_KEY_JOURNAL_ID:
			dm.journalId = a.Value
		case model.EV_KEY_VOLUME_ISSUE:
			dm.issue = a.Value
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationVolumeCreate struct {
	volumeId  string
	journalId string
	createdOn int64
	issue     string
}

var _ dataManipulation = new(dataManipulationVolumeCreate)

func (dm *dataManipulationVolumeCreate) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec("INSERT INTO volume VALUES(?, ?, ?, ?)",
		dm.volumeId, dm.createdOn, dm.journalId, dm.issue)
	return err
}

func GetVolumesOfJournal(journalId string) ([]Volume, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Commit() }()
	result := []Volume{}
	err = tx.Select(&result, "SELECT * FROM volume WHERE journalId = ? ORDER BY issue DESC",
		journalId)
	return result, err
}

type Volume struct {
	VolumeId  string
	CreatedOn int64
	JournalId string
	Issue     string
}

func GetVolume(volumeId string) (*Volume, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Commit() }()
	return getVolumeFromTransaction(tx, volumeId)
}

func getVolumeFromTransaction(tx *sqlx.Tx, volumeId string) (*Volume, error) {
	result := &Volume{}
	err := tx.Get(result, "SELECT * FROM volume WHERE volumeId = ? ORDER BY issue DESC",
		volumeId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetVolumeView(volumeId string) (*VolumeView, error) {
	result := new(VolumeView)
	var err error
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Commit() }()
	result.Volume, err = getVolumeFromTransaction(tx, volumeId)
	if err != nil {
		return nil, err
	}
	result.Journal, err = getJournalFromTransaction(tx, result.Volume.JournalId)
	if err != nil {
		return nil, err
	}
	manuscriptIds, err := getManuscriptsOfVolumeFromTransaction(volumeId, tx)
	if err != nil {
		return nil, err
	}
	manuscripts := make([]*Manuscript, len(manuscriptIds))
	for index, id := range manuscriptIds {
		manuscripts[index], err = getManuscriptFromTransaction(tx, id)
		if err != nil {
			return nil, err
		}
	}
	result.Manuscripts = manuscripts
	return result, nil
}

type VolumeView struct {
	Volume      *Volume
	Journal     *Journal
	Manuscripts []*Manuscript
}
