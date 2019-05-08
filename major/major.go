package main

import (
	"gitlab.bbinfra.net/3estack/alexandria/blockchain"
	"gitlab.bbinfra.net/3estack/alexandria/cli"
	"gitlab.bbinfra.net/3estack/alexandria/cliAlexandria"
	"gitlab.bbinfra.net/3estack/alexandria/command"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
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
		Handlers: append(cliAlexandria.CommonHandlers,
			&cli.Cli{
				FullDescription:    "Welcome to the Bootstrap and Settings Update commands",
				OneLineDescription: "Settings",
				Name:               "settings",
				Handlers: []cli.Handler{
					&cli.StructRunnerHandler{
						FullDescription:    "Welcome to the dialog to bootstrap Alexandria",
						OneLineDescription: "Bootstrap",
						Name:               "bootstrap",
						Action:             command.GetBootstrapCommand,
					},
					&cli.StructRunnerHandler{
						FullDescription:              "Welcome to the settings update dialog",
						OneLineDescription:           "Settings Update",
						Name:                         "settingsUpdate",
						ReferenceValueGetter:         settingsUpdateReference,
						ReferenceValueGetterArgNames: []string{},
						Action:                       settingsUpdate,
					},
				},
			},
			&cli.Cli{
				FullDescription:    "Welcome to the person commands",
				OneLineDescription: "Person",
				Name:               "person",
				Handlers: []cli.Handler{
					&cli.StructRunnerHandler{
						FullDescription:    "Welcome to the person create dialog.",
						OneLineDescription: "Create Person",
						Name:               "createPerson",
						Action:             createPerson,
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
				},
			},
		),
	}
	context.Run()
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
	if err := blockchain.SendCommand(theCommand); err != nil {
		outputter(cliAlexandria.ToIoError(err))
		return
	}
}

func createPerson(outputter cli.Outputter, personInput *command.PersonCreate) {
	cliAlexandria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetPersonCreateCommand(
			personInput,
			cliAlexandria.LoggedInPerson.Id,
			cliAlexandria.LoggedIn(),
			cliAlexandria.Settings.PriceMajorCreatePerson)
	})
}

func setMajor(outputter cli.Outputter, personId string) {
	cliAlexandria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetPersonSetMajorCommand(
			personId,
			cliAlexandria.LoggedInPerson.Id,
			cliAlexandria.LoggedIn(),
			cliAlexandria.Settings.PriceMajorChangePersonAuthorization)
	})
}

func unsetMajor(outputter cli.Outputter, personId string) {
	cliAlexandria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetPersonUnsetMajorCommand(
			personId,
			cliAlexandria.LoggedInPerson.Id,
			cliAlexandria.LoggedIn(),
			cliAlexandria.Settings.PriceMajorChangePersonAuthorization)
	})
}

func setSigned(outputter cli.Outputter, personId string) {
	cliAlexandria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetPersonSetSignedCommand(
			personId,
			cliAlexandria.LoggedInPerson.Id,
			cliAlexandria.LoggedIn(),
			cliAlexandria.Settings.PriceMajorChangePersonAuthorization)
	})
}

func unsetSigned(outputter cli.Outputter, personId string) {
	cliAlexandria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetPersonUnsetSignedCommand(
			personId,
			cliAlexandria.LoggedInPerson.Id,
			cliAlexandria.LoggedIn(),
			cliAlexandria.Settings.PriceMajorChangePersonAuthorization)
	})
}

func incBalance(outputter cli.Outputter, personId string, amount int32) {
	cliAlexandria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetPersonIncBalanceCommand(
			personId,
			amount,
			cliAlexandria.LoggedInPerson.Id,
			cliAlexandria.LoggedIn(),
			int32(0))
	})
}
