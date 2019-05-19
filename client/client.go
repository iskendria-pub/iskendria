package main

import (
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/cli"
	"gitlab.bbinfra.net/3estack/alexandria/cliAlexandria"
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
						OneLineDescription:           "Person Update",
						Name:                         "updatePerson",
						ReferenceValueGetter:         personUpdateReference,
						ReferenceValueGetterArgNames: []string{},
						Action:                       cliAlexandria.PersonUpdate,
					},
				),
			},
		),
	}
	fmt.Print(makeGreen)
	dao.Init("client.db", log.New(os.Stdout, "db", log.Flags()))
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
