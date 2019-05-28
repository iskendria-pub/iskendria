package main

import (
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/blockchain"
	"gitlab.bbinfra.net/3estack/alexandria/cli"
	"gitlab.bbinfra.net/3estack/alexandria/cliAlexandria"
	"gitlab.bbinfra.net/3estack/alexandria/command"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"log"
	"os"
	"strings"
)

var description = strings.TrimSpace(`
Welcome to the Alexandria Client Tool. Use this tool to
register and to manage manuscripts, reviews and journals.
`)

var makeGreen = "\033[32m"

func main() {
	context := &cli.Cli{
		FullDescription:    description,
		OneLineDescription: "Alexandria Client Tool",
		Name:               "alexandria-client",
		FormatEscape:       makeGreen,
		EventPager:         cliAlexandria.PageEventStreamMessages,
		Handlers: append(cliAlexandria.CommonRootHandlers,
			cliAlexandria.CommonDiagnosticsGroup,
			&cli.Cli{
				FullDescription:    "Welcome to the settings commands",
				OneLineDescription: "Settings",
				Name:               "settings",
				Handlers:           cliAlexandria.CommonSettingsHandlers,
			},
			&cli.Cli{
				FullDescription:    "Welcome to the person commands",
				OneLineDescription: "Person",
				Name:               "person",
				Handlers: append(cliAlexandria.CommonPersonHandlers,
					&cli.StructRunnerHandler{
						FullDescription:              "Welcome to the person update dialog.",
						OneLineDescription:           "Update person",
						Name:                         "updatePerson",
						ReferenceValueGetter:         personUpdateReference,
						ReferenceValueGetterArgNames: []string{},
						Action:                       cliAlexandria.PersonUpdate,
					},
				),
			},
			&cli.Cli{
				FullDescription:    "Welcome to the journal commands.",
				OneLineDescription: "Journal",
				Name:               "journal",
				Handlers: []cli.Handler{
					&cli.StructRunnerHandler{
						FullDescription:    "Welcome to the journal create dialog.",
						OneLineDescription: "Create journal",
						Name:               "createJournal",
						Action:             journalCreate,
					},
				},
			},
		),
	}
	fmt.Print(makeGreen)
	dbLogger := log.New(os.Stdout, "db", log.Flags())
	dao.Init("client.db", dbLogger)
	defer dao.Shutdown(dbLogger)
	cliAlexandria.InitEventStream("./client-events.log", "client")
	context.Run()
}

func personUpdateReference(outputter cli.Outputter) *dao.PersonUpdate {
	if !cliAlexandria.CheckBootstrappedAndKnownPerson(outputter) {
		return nil
	}
	cliAlexandria.OriginalPersonId = cliAlexandria.LoggedInPerson.Id
	cliAlexandria.OriginalPerson = dao.PersonToPersonUpdate(cliAlexandria.LoggedInPerson)
	return cliAlexandria.OriginalPerson
}

func journalCreate(outputter cli.Outputter, journal *command.Journal) {
	if !cliAlexandria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	cmd, journalId := command.GetCommandJournalCreate(
		journal,
		cliAlexandria.LoggedInPerson.Id,
		cliAlexandria.LoggedIn(),
		cliAlexandria.Settings.PriceEditorCreateJournal)
	if err := blockchain.SendCommand(cmd, outputter); err != nil {
		outputter(cliAlexandria.ToIoError(err) + "\n")
		return
	}
	outputter("The journalId of the created journal is: " + journalId + "\n")
}
