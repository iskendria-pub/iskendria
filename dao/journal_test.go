package dao

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"log"
	"os"
	"testing"
)

func TestSortJournals(t *testing.T) {
	jts := getJournalTestItems()
	cases := getCases(jts)
	for i, currentCase := range cases {
		sortJournals(currentCase)
		if !(currentCase[0].JournalId == "03" &&
			currentCase[1].JournalId == "04" &&
			currentCase[2].JournalId == "01" &&
			currentCase[3].JournalId == "02") {
			t.Error(fmt.Sprintf("Case %d not sorted correctly", i))
		}
	}
}

func getJournalTestItems() *journalTestItems {
	firstTitle := "A"
	secondTitle := "Z"
	first := &Journal{
		JournalId: "01",
		Title:     secondTitle,
	}
	second := &Journal{
		JournalId: "02",
		Title:     secondTitle,
	}
	third := &Journal{
		JournalId: "03",
		Title:     firstTitle,
	}
	fourth := &Journal{
		JournalId: "04",
		Title:     firstTitle,
	}
	return &journalTestItems{
		first:  first,
		second: second,
		third:  third,
		fourth: fourth,
	}
}

type journalTestItems struct {
	first,
	second,
	third,
	fourth *Journal
}

func getCases(jts *journalTestItems) [][]*Journal {
	return [][]*Journal{
		{
			jts.first, jts.second, jts.third, jts.fourth,
		},
		{
			jts.first, jts.second, jts.fourth, jts.third,
		},
		{
			jts.first, jts.third, jts.second, jts.fourth,
		},
		{
			jts.first, jts.third, jts.fourth, jts.second,
		},
		{
			jts.first, jts.fourth, jts.second, jts.third,
		},
		{
			jts.first, jts.fourth, jts.third, jts.second,
		},
		{
			jts.second, jts.first, jts.third, jts.fourth,
		},
		{
			jts.second, jts.first, jts.fourth, jts.third,
		},
		{
			jts.second, jts.third, jts.first, jts.fourth,
		},
		{
			jts.second, jts.third, jts.fourth, jts.first,
		},
		{
			jts.second, jts.fourth, jts.first, jts.third,
		},
		{
			jts.second, jts.fourth, jts.third, jts.first,
		},
		{
			jts.third, jts.first, jts.second, jts.fourth,
		},
		{
			jts.third, jts.first, jts.fourth, jts.second,
		},
		{
			jts.third, jts.second, jts.first, jts.fourth,
		},
		{
			jts.third, jts.second, jts.fourth, jts.first,
		},
		{
			jts.third, jts.fourth, jts.first, jts.second,
		},
		{
			jts.third, jts.fourth, jts.second, jts.first,
		},
		{
			jts.fourth, jts.first, jts.second, jts.third,
		},
		{
			jts.fourth, jts.first, jts.third, jts.second,
		},
		{
			jts.fourth, jts.second, jts.first, jts.third,
		},
		{
			jts.fourth, jts.second, jts.third, jts.first,
		},
		{
			jts.fourth, jts.third, jts.first, jts.second,
		},
		{
			jts.fourth, jts.third, jts.second, jts.first,
		},
	}
}

/*
Tests function GetAllJournals. We test that the id of the journal,
the id of an editor, the isSigned property of a journal and the
isSigned property of a person are not confused. The code fetches
journal-editor combinations from the database and then branches
on whether the first record of a new journal is being processed.
Both branches are covered.

We use that the code sorts everything on id, provided no titles
and names are provided. We create four journals each having two
editors. We consider the booleans journal-is-signed, first-editor-signed
and second-editor-signed. For journal-is-signed and first-editor-signed
we need all combinations. For journal-is-signed and second-editor-signed
we also need all combinations. But for first-editor-signed and
second-editor-signed we do not need all combinations.

We add a fifth and sixth journal to check that journals without editors
are also covered. The fifth has no editors at all while the sixth has
only a proposed editor. Finally, we give the first journal a proposed
editor to check that proposed editors are omitted.
*/
func TestGetAllJournals(t *testing.T) {
	logger := log.New(os.Stdout, "testGetAllJournals", log.Flags())
	Init("testGetAllJournals.db", logger)
	defer ShutdownAndDelete(logger)
	doTestGetAllJournals(logger, t)
}

func doTestGetAllJournals(_ *log.Logger, t *testing.T) {
	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() { _ = tx.Commit() }()
	insertJournal("j1", false, tx, t)
	insertJournal("j2", false, tx, t)
	insertJournal("j3", true, tx, t)
	insertJournal("j4", true, tx, t)
	insertJournal("j5", true, tx, t)
	insertJournal("j6", true, tx, t)
	insertPerson("p1", false, tx, t)
	insertPerson("p2", true, tx, t)
	insertPerson("p3", false, tx, t)
	insertPerson("notAccepted", true, tx, t)

	// j1 and j2 cover journalIsSigned = true

	// j1 covers first-editor-signed = false and second-editor-signed = true
	insertAcceptedEditor("j1", "p1", tx, t)
	insertAcceptedEditor("j1", "p2", tx, t)
	insertEditor(
		"j1",
		"notAccepted",
		model.GetEditorStateString(model.EditorState_editorProposed),
		tx,
		t)

	// j2 covers first-editor-signed = true and second-editor-signed = false
	insertAcceptedEditor("j2", "p2", tx, t)
	insertAcceptedEditor("j2", "p3", tx, t)

	// j3 and j4 cover journalIsSigned = true

	insertAcceptedEditor("j3", "p1", tx, t)
	insertAcceptedEditor("j3", "p2", tx, t)
	insertAcceptedEditor("j4", "p2", tx, t)
	insertAcceptedEditor("j4", "p3", tx, t)

	err = tx.Commit()
	if err != nil {
		t.Error(err)
	}

	actual, err := GetAllJournals()
	if err != nil {
		t.Error(err)
	}

	if len(actual) != 6 {
		t.Error("Expected six journals")
	}
	for i := 0; i < 4; i++ {
		if len(actual[i].AcceptedEditors) != 2 {
			t.Error(fmt.Sprintf("Expected that journal #%d has two editors", i))
		}
	}
	for i := 4; i < 6; i++ {
		if len(actual[i].AcceptedEditors) != 0 {
			t.Error("Journal without editors mismatch")
		}
	}
	if actual[0].JournalId != "j1" || actual[0].IsSigned != false {
		t.Error("Journal 0 not as expected")
	}
	if actual[1].JournalId != "j2" || actual[1].IsSigned != false {
		t.Error("Journal 1 not as expected")
	}
	if actual[2].JournalId != "j3" || actual[2].IsSigned != true {
		t.Error("Journal 2 not as expected")
	}
	if actual[3].JournalId != "j4" || actual[3].IsSigned != true {
		t.Error("Journal 3 not as expected")
	}

	checkActualEditor(actual[0].AcceptedEditors[0], "p1", false, t)
	checkActualEditor(actual[0].AcceptedEditors[1], "p2", true, t)
	checkActualEditor(actual[1].AcceptedEditors[0], "p2", true, t)
	checkActualEditor(actual[1].AcceptedEditors[1], "p3", false, t)
	checkActualEditor(actual[2].AcceptedEditors[0], "p1", false, t)
	checkActualEditor(actual[2].AcceptedEditors[1], "p2", true, t)
	checkActualEditor(actual[3].AcceptedEditors[0], "p2", true, t)
	checkActualEditor(actual[3].AcceptedEditors[1], "p3", false, t)
}

func insertJournal(id string, isSigned bool, tx *sqlx.Tx, t *testing.T) {
	_, err := tx.Exec(fmt.Sprintf("INSERT INTO journal VALUES (%s)", GetPlaceHolders(6)),
		id, 0, 0, "", isSigned, "")
	if err != nil {
		t.Error(err)
	}
}

func insertPerson(id string, isSigned bool, tx *sqlx.Tx, t *testing.T) {
	_, err := tx.Exec(fmt.Sprintf("INSERT INTO person VALUES (%s)", GetPlaceHolders(16)),
		id, 0, 0, "", "",
		"", false, isSigned, int32(0), "",
		"", "", "", "", "",
		"")
	if err != nil {
		t.Error(err)
	}
}

func insertEditor(journalId, personId string, editorState string, tx *sqlx.Tx, t *testing.T) {
	_, err := tx.Exec("INSERT INTO editor(journalid, personid, editorstate) VALUES(?, ?, ?)",
		journalId, personId, editorState)
	if err != nil {
		t.Error(err)
	}
}

func insertAcceptedEditor(journalId, personId string, tx *sqlx.Tx, t *testing.T) {
	insertEditor(journalId, personId, model.GetEditorStateString(model.EditorState_editorAccepted), tx, t)
}

func checkActualEditor(actualEditor *Editor, expectedPersonId string, expectedIsSigned bool, t *testing.T) {
	if actualEditor.PersonId != expectedPersonId {
		t.Error("Person id mismatch")
	}
	if actualEditor.PersonIsSigned != expectedIsSigned {
		t.Error("PersonIsSigned mismatch")
	}
}

/*
This test covers both GetJournal and GetJournal and GetJournalIncludingProposedEditors
GetJournal returns a journal and provides all accepted editors.
GetJournalIncludingProposedEditors returns a journal and provides all accepted and
proposed editors.
*/
func TestGettingOneJournal(t *testing.T) {
	logger := log.New(os.Stdout, "testGettingOneJournal", log.Flags())
	Init("testGettingOneJournal.db", logger)
	defer ShutdownAndDelete(logger)
	doTestGettingOneJournal(logger, t)
}

func doTestGettingOneJournal(_ *log.Logger, t *testing.T) {
	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() { _ = tx.Commit() }()

	insertJournal("j1", false, tx, t)
	insertJournal("j2", false, tx, t)
	insertPerson("p1", false, tx, t)
	insertPerson("p2", false, tx, t)
	insertEditor("j2", "p1", model.GetEditorStateString(model.EditorState_editorProposed), tx, t)
	insertEditor("j2", "p2", model.GetEditorStateString(model.EditorState_editorAccepted), tx, t)
	err = tx.Commit()
	if err != nil {
		t.Error(err)
	}

	actualOnlyAccepted, err := GetJournal("j1")
	if err != nil {
		t.Error(err)
	}
	if actualOnlyAccepted == nil {
		t.Error("Expected that GetJournal() returns a journal")
	}
	if len(actualOnlyAccepted.AcceptedEditors) != 0 {
		t.Error("Journal j1 is expected not to have accepted editors")
	}
	actualOnlyAccepted, err = GetJournal("j2")
	if err != nil {
		t.Error(err)
	}
	if actualOnlyAccepted == nil {
		t.Error("Expected that GetJournal() returns a journal")
	}
	if len(actualOnlyAccepted.AcceptedEditors) != 1 {
		t.Error("Expected that j2 has one accepted editor")
	}
	if actualOnlyAccepted.AcceptedEditors[0].PersonId != "p2" {
		t.Error("Got the wrong accepted editor")
	}

	actualAll, err := GetJournalIncludingProposedEditors("j1")
	if err != nil {
		t.Error(err)
	}
	if actualAll == nil {
		t.Error("Expected that GetJournalIncludingProposedEditors returns a journal")
	}
	if len(actualAll.AllEditors) != 0 {
		t.Error("Expected no editors")
	}
	actualAll, err = GetJournalIncludingProposedEditors("j2")
	if err != nil {
		t.Error(err)
	}
	if len(actualAll.AllEditors) != 2 {
		t.Error("Expected two editors")
	}
	checkAllEditorsObject(
		actualAll.AllEditors[0], "p1", model.GetEditorStateString(model.EditorState_editorProposed), t)
	checkAllEditorsObject(
		actualAll.AllEditors[1], "p2", model.GetEditorStateString(model.EditorState_editorAccepted), t)
}

func checkAllEditorsObject(e *EditorWithState, expectedPersonId, expectedEditorState string, t *testing.T) {
	if e.PersonId != expectedPersonId {
		t.Error("Person mismatch")
	}
	if e.EditorState != expectedEditorState {
		t.Error("EditorState mismatch")
	}
}
