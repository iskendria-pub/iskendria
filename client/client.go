package main

import (
	"gitlab.bbinfra.net/3estack/alexandria/blockchain"
	"gitlab.bbinfra.net/3estack/alexandria/cli"
	"gitlab.bbinfra.net/3estack/alexandria/cliAlexandria"
	"gitlab.bbinfra.net/3estack/alexandria/command"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
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
		Handlers: []cli.Handler{
			&cli.StructRunnerHandler{
				FullDescription:              "Welcome to the person update dialog.",
				OneLineDescription:           "Person Update",
				Name:                         "person-update",
				ReferenceValueGetter:         personUpdateReference,
				ReferenceValueGetterArgNames: []string{},
				Action:                       personUpdate,
			},
		},
	}
	context.Run()
}

func personUpdateReference(outputter cli.Outputter) *dao.PersonUpdate {
	if !cliAlexandria.CheckBootstrappedAndKnownPerson(outputter) {
		return nil
	}
	originalPerson = dao.PersonToPersonUpdate(cliAlexandria.LoggedInPerson)
	return originalPerson
}

var originalPerson *dao.PersonUpdate

func personUpdate(outputter cli.Outputter, newPerson *dao.PersonUpdate) {
	theCommand := command.GetPersonUpdatePropertiesCommand(
		cliAlexandria.LoggedInPerson.Id,
		originalPerson,
		newPerson,
		cliAlexandria.LoggedInPerson.Id,
		cliAlexandria.LoggedIn(),
		cliAlexandria.Settings.PricePersonEdit)
	if err := blockchain.SendCommand(theCommand); err != nil {
		outputter(cliAlexandria.ToIoError(err))
	}
}
