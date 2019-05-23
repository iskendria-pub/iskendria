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

const (
	makeRed = "\033[31m"
)

var description = strings.TrimSpace(`
Alexandria Major Tool. Alexandria manages scientific publications
on a Hyperledger Sawtooth blockchain. This is the user interface
for majors, members of the Alexandria team. This tool allows them
to manage the blockchain.
`)

func main() {
	context := &cli.Cli{
		FullDescription:    description,
		OneLineDescription: "Alexandria Major Tool",
		Name:               "alexandria-major",
		FormatEscape:       makeRed,
		EventPager:         cliAlexandria.PageEventStreamMessages,
		Handlers: append(cliAlexandria.CommonRootHandlers,
			cliAlexandria.CommonDiagnosticsGroup,
			&cli.Cli{
				FullDescription:    "Welcome to the Bootstrap and Settings Update commands",
				OneLineDescription: "Settings",
				Name:               "settings",
				Handlers: append(cliAlexandria.CommonSettingsHandlers,
					&cli.StructRunnerHandler{
						FullDescription:    "Welcome to the dialog to bootstrap Alexandria",
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
				Handlers: append(cliAlexandria.CommonPersonHandlers,
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
						Action:                       cliAlexandria.PersonUpdate,
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
		),
	}
	if len(os.Args) >= 2 {
		cli.InputScript = os.Args[1]
	}
	fmt.Print(makeRed)
	dbLogger := log.New(os.Stdout, "db", log.Flags())
	dao.Init("major.db", dbLogger)
	defer dao.Shutdown(dbLogger)
	cliAlexandria.InitEventStream("./major-events.log", "major")
	context.Run()
}

func bootstrap(outputter cli.Outputter, bootstrap *command.Bootstrap) {
	if !cliAlexandria.IsLoggedIn() {
		outputter("You should login before you can bootstrap\n")
		return
	}
	cmd := command.GetBootstrapCommand(bootstrap, cliAlexandria.LoggedIn())
	if err := blockchain.SendCommand(cmd, outputter); err != nil {
		outputter(cliAlexandria.ToIoError(err))
	}
}

func settingsUpdateReference(outputter cli.Outputter) *dao.Settings {
	if !cliAlexandria.CheckBootstrappedAndKnownPerson(outputter) {
		return nil
	}
	return cliAlexandria.Settings
}

func settingsUpdate(outputter cli.Outputter, updated *dao.Settings) {
	theCommand := command.GetSettingsUpdateCommand(
		cliAlexandria.Settings,
		updated,
		cliAlexandria.LoggedInPerson.Id,
		cliAlexandria.LoggedIn(),
		cliAlexandria.Settings.PriceMajorEditSettings)
	if err := blockchain.SendCommand(theCommand, outputter); err != nil {
		outputter(cliAlexandria.ToIoError(err))
	}
}

func personCreate(outputter cli.Outputter, personInput *command.PersonCreate) {
	if !cliAlexandria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	cmd, personId := command.GetPersonCreateCommand(
		personInput,
		cliAlexandria.LoggedInPerson.Id,
		cliAlexandria.LoggedIn(),
		cliAlexandria.Settings.PriceMajorCreatePerson)
	if err := blockchain.SendCommand(cmd, outputter); err != nil {
		outputter(cliAlexandria.ToIoError(err) + "\n")
		return
	}
	outputter("The personId of the created person is: " + personId + "\n")
}

func setMajor(outputter cli.Outputter, personId string) {
	cliAlexandria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetPersonUpdateSetMajorCommand(
			personId,
			cliAlexandria.LoggedInPerson.Id,
			cliAlexandria.LoggedIn(),
			cliAlexandria.Settings.PriceMajorChangePersonAuthorization)
	})
}

func unsetMajor(outputter cli.Outputter, personId string) {
	cliAlexandria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetPersonUpdateUnsetMajorCommand(
			personId,
			cliAlexandria.LoggedInPerson.Id,
			cliAlexandria.LoggedIn(),
			cliAlexandria.Settings.PriceMajorChangePersonAuthorization)
	})
}

func setSigned(outputter cli.Outputter, personId string) {
	cliAlexandria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetPersonUpdateSetSignedCommand(
			personId,
			cliAlexandria.LoggedInPerson.Id,
			cliAlexandria.LoggedIn(),
			cliAlexandria.Settings.PriceMajorChangePersonAuthorization)
	})
}

func unsetSigned(outputter cli.Outputter, personId string) {
	cliAlexandria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetPersonUpdateUnsetSignedCommand(
			personId,
			cliAlexandria.LoggedInPerson.Id,
			cliAlexandria.LoggedIn(),
			cliAlexandria.Settings.PriceMajorChangePersonAuthorization)
	})
}

func incBalance(outputter cli.Outputter, personId string, amount int32) {
	cliAlexandria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetPersonUpdateIncBalanceCommand(
			personId,
			amount,
			cliAlexandria.LoggedInPerson.Id,
			cliAlexandria.LoggedIn(),
			int32(0))
	})
}

func personUpdateReference(outputter cli.Outputter, personId string) *dao.PersonUpdate {
	if !cliAlexandria.CheckBootstrappedAndKnownPerson(outputter) {
		return nil
	}
	person, err := dao.GetPersonById(personId)
	if err != nil {
		outputter(fmt.Sprintf("Could not find person %s, error: %s\n", personId, err.Error()))
		return nil
	}
	cliAlexandria.OriginalPersonId = personId
	cliAlexandria.OriginalPerson = dao.PersonToPersonUpdate(person)
	return cliAlexandria.OriginalPerson
}
