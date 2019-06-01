package cliAlexandria

import (
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/cli"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"gitlab.bbinfra.net/3estack/alexandria/model"
)

var CommonJournalHandlers = []cli.Handler{
	&cli.SingleLineHandler{
		Name:     "showJournal",
		Handler:  showJournal,
		ArgNames: []string{"journal id"},
	},
}

func showJournal(outputter cli.Outputter, journalId string) {
	journal, err := dao.GetJournalIncludingProposedEditors(journalId)
	if err != nil {
		outputter(fmt.Sprintf("Journal does not exist: %s, detailed message: %s\n",
			journalId, err.Error()))
		return
	}
	tableAcceptedEditors := getEditorsWithState(journal, model.EditorState_editorAccepted)
	tableProposedEditors := getEditorsWithState(journal, model.EditorState_editorProposed)
	tableJournal := cli.StructToTable(JournalToJournalWithoutEditorsView(journal))
	outputter("Journal properties:\n\n" + tableJournal.String() +
		"\nAccepted editors\n\n" + tableAcceptedEditors.String() +
		"\nProposed editors\n\n" + tableProposedEditors.String() + "\n")
}

func getEditorsWithState(
	journal *dao.JournalIncludingProposedEditors, editorStateEnum model.EditorState) *cli.TableType {
	editorState := model.GetEditorStateString(editorStateEnum)
	editors := []*dao.Editor{}
	for _, e := range journal.AllEditors {
		if e.EditorState == editorState {
			editors = append(editors, &dao.Editor{
				PersonId:   e.PersonId,
				PersonName: e.PersonName,
			})
		}
	}
	result := cli.NewTable(len(editors), 2)
	for i := range editors {
		result.Set(i, 0, editors[i].PersonName)
		result.Set(i, 1, editors[i].PersonId)
	}
	return result
}

func JournalToJournalWithoutEditorsView(journal *dao.JournalIncludingProposedEditors) *JournalWithoutEditorsView {
	return &JournalWithoutEditorsView{
		JournalId:       journal.JournalId,
		CreatedOn:       formatTime(journal.CreatedOn),
		ModifiedOn:      formatTime(journal.ModifiedOn),
		Title:           journal.Title,
		IsSigned:        journal.IsSigned,
		Descriptionhash: journal.Descriptionhash,
	}
}

type JournalWithoutEditorsView struct {
	JournalId       string
	CreatedOn       string
	ModifiedOn      string
	Title           string
	IsSigned        bool
	Descriptionhash string
}
