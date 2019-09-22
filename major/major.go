package main

import (
	"fmt"
	"github.com/iskendria-pub/iskendria/blockchain"
	"github.com/iskendria-pub/iskendria/cli"
	"github.com/iskendria-pub/iskendria/cliIskendria"
	"github.com/iskendria-pub/iskendria/command"
	"github.com/iskendria-pub/iskendria/dao"
	"log"
	"os"
	"strings"
)

const (
	makeRed = "\033[31m"
)

var description = strings.TrimSpace(`
Iskendria Major Tool. Iskendria manages scientific publications
on a Hyperledger Sawtooth blockchain. This is the user interface
for majors, members of the Iskendria team. This tool allows them
to manage the blockchain.
`)

func main() {
	context := &cli.Cli{
		FullDescription:    description,
		OneLineDescription: "Iskendria Major Tool",
		Name:               "iskendria-major",
		FormatEscape:       makeRed,
		EventPager:         cliIskendria.PageEventStreamMessages,
		Handlers: append(cliIskendria.CommonRootHandlers,
			cliIskendria.CommonDiagnosticsGroup,
			&cli.Cli{
				FullDescription:    "Welcome to the Bootstrap and Settings Update commands",
				OneLineDescription: "Settings",
				Name:               "settings",
				Handlers: append(cliIskendria.CommonSettingsHandlers,
					&cli.StructRunnerHandler{
						FullDescription:    "Welcome to the dialog to bootstrap Iskendria",
						OneLineDescription: "Bootstrap",
						Name:               "bootstrap",
						Action:             bootstrap,
					},
					&cli.StructRunnerHandler{
						FullDescription:              "Welcome to the settings update dialog",
						OneLineDescription:           "Update settings",
						Name:                         "updateSettings",
						ReferenceValueGetter:         settingsUpdateReference,
						ReferenceValueGetterArgNames: []string{},
						Action:                       settingsUpdate,
					},
				),
			},
			&cli.Cli{
				FullDescription:    "Welcome to the person commands",
				OneLineDescription: "Person",
				Name:               "person",
				Handlers: append(cliIskendria.CommonPersonHandlers,
					&cli.StructRunnerHandler{
						FullDescription:    "Welcome to the person create dialog.",
						OneLineDescription: "Create Person",
						Name:               "createPerson",
						Action:             personCreate,
					},
					&cli.StructRunnerHandler{
						FullDescription:              "Welcome to the person update dialog.",
						OneLineDescription:           "Update person",
						Name:                         "updatePerson",
						ReferenceValueGetter:         personUpdateReference,
						ReferenceValueGetterArgNames: []string{"person id"},
						Action:                       cliIskendria.PersonUpdate,
					},
					&cli.SingleLineHandler{
						Name:     "setMajor",
						Handler:  setMajor,
						ArgNames: []string{"person id"},
					},
					&cli.SingleLineHandler{
						Name:     "unsetMajor",
						Handler:  unsetMajor,
						ArgNames: []string{"person id"},
					},
					&cli.SingleLineHandler{
						Name:     "setSigned",
						Handler:  setSigned,
						ArgNames: []string{"person id"},
					},
					&cli.SingleLineHandler{
						Name:     "unsetSigned",
						Handler:  unsetSigned,
						ArgNames: []string{"person id"},
					},
					&cli.SingleLineHandler{
						Name:     "incBalance",
						Handler:  incBalance,
						ArgNames: []string{"person id", "amount"},
					},
				),
			},
			&cli.Cli{
				FullDescription:    "Welcome to the journal commands.",
				OneLineDescription: "Journal",
				Name:               "journal",
				Handlers: append(cliIskendria.CommonJournalHandlers,
					&cli.SingleLineHandler{
						Name:     "setSigned",
						Handler:  journalSetSigned,
						ArgNames: []string{"journal id"},
					},
					&cli.SingleLineHandler{
						Name:     "unsetSigned",
						Handler:  journalUnsetSigned,
						ArgNames: []string{"journal id"},
					},
				),
			},
		),
	}
	if len(os.Args) >= 2 {
		cli.InputScript = os.Args[1]
	}
	fmt.Print(makeRed)
	dbLogger := log.New(os.Stdout, "db", log.Flags())
	dao.Init("major.db", dbLogger)
	defer dao.Shutdown(dbLogger)
	cliIskendria.InitEventStream("./major-events.log", "major")
	context.Run()
}

func bootstrap(outputter cli.Outputter, bootstrap *command.Bootstrap) {
	if !cliIskendria.IsLoggedIn() {
		outputter("You should login before you can bootstrap\n")
		return
	}
	cmd := command.GetBootstrapCommand(bootstrap, cliIskendria.LoggedIn())
	if err := blockchain.SendCommand(cmd, outputter); err != nil {
		outputter(cliIskendria.ToIoError(err))
	}
}

func settingsUpdateReference(outputter cli.Outputter) *dao.Settings {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return nil
	}
	return cliIskendria.Settings
}

func settingsUpdate(outputter cli.Outputter, updated *dao.Settings) {
	theCommand := command.GetSettingsUpdateCommand(
		cliIskendria.Settings,
		updated,
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PriceMajorEditSettings)
	if err := blockchain.SendCommand(theCommand, outputter); err != nil {
		outputter(cliIskendria.ToIoError(err))
	}
}

func personCreate(outputter cli.Outputter, personInput *command.PersonCreate) {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	cmd, personId := command.GetPersonCreateCommand(
		personInput,
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PriceMajorCreatePerson)
	if err := blockchain.SendCommand(cmd, outputter); err != nil {
		outputter(cliIskendria.ToIoError(err) + "\n")
		return
	}
	outputter("The personId of the created person is: " + personId + "\n")
}

func setMajor(outputter cli.Outputter, personId string) {
	cliIskendria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetPersonUpdateSetMajorCommand(
			personId,
			cliIskendria.LoggedInPerson.Id,
			cliIskendria.LoggedIn(),
			cliIskendria.Settings.PriceMajorChangePersonAuthorization)
	})
}

func unsetMajor(outputter cli.Outputter, personId string) {
	cliIskendria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetPersonUpdateUnsetMajorCommand(
			personId,
			cliIskendria.LoggedInPerson.Id,
			cliIskendria.LoggedIn(),
			cliIskendria.Settings.PriceMajorChangePersonAuthorization)
	})
}

func setSigned(outputter cli.Outputter, personId string) {
	cliIskendria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetPersonUpdateSetSignedCommand(
			personId,
			cliIskendria.LoggedInPerson.Id,
			cliIskendria.LoggedIn(),
			cliIskendria.Settings.PriceMajorChangePersonAuthorization)
	})
}

func unsetSigned(outputter cli.Outputter, personId string) {
	cliIskendria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetPersonUpdateUnsetSignedCommand(
			personId,
			cliIskendria.LoggedInPerson.Id,
			cliIskendria.LoggedIn(),
			cliIskendria.Settings.PriceMajorChangePersonAuthorization)
	})
}

func incBalance(outputter cli.Outputter, personId string, amount int32) {
	cliIskendria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetPersonUpdateIncBalanceCommand(
			personId,
			amount,
			cliIskendria.LoggedInPerson.Id,
			cliIskendria.LoggedIn(),
			int32(0))
	})
}

func personUpdateReference(outputter cli.Outputter, personId string) *dao.PersonUpdate {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return nil
	}
	person, err := dao.GetPersonById(personId)
	if err != nil {
		outputter(fmt.Sprintf("Could not find person %s, error: %s\n", personId, err.Error()))
		return nil
	}
	cliIskendria.OriginalPersonId = personId
	cliIskendria.OriginalPerson = dao.PersonToPersonUpdate(person)
	return cliIskendria.OriginalPerson
}

func journalSetSigned(outputter cli.Outputter, journalId string) {
	journalChangeAuthorization(outputter, journalId, true)
}

func journalUnsetSigned(outputter cli.Outputter, journalId string) {
	journalChangeAuthorization(outputter, journalId, false)
}

func journalChangeAuthorization(outputter cli.Outputter, journalId string, makeSigned bool) {
	cliIskendria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetCommandJournalUpdateAuthorization(
			journalId,
			makeSigned,
			cliIskendria.LoggedInPerson.Id,
			cliIskendria.LoggedIn(),
			cliIskendria.Settings.PriceMajorChangeJournalAuthorization)
	})
}
