package main

import (
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/blockchain"
	"gitlab.bbinfra.net/3estack/alexandria/cli"
	"gitlab.bbinfra.net/3estack/alexandria/cliAlexandria"
	"gitlab.bbinfra.net/3estack/alexandria/command"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"io/ioutil"
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
					&cli.Cli{
						FullDescription:    "Welcome to the person biography commands",
						OneLineDescription: "Update/Verify/Remove biography",
						Name:               "biography",
						Handlers: []cli.Handler{
							&cli.SingleLineHandler{
								Name:     "updateBiography",
								Handler:  personUpdateBiography,
								ArgNames: []string{"biography filename"},
							},
							&cli.SingleLineHandler{
								Name:     "removeBiography",
								Handler:  personRemoveBiography,
								ArgNames: []string{},
							},
							&cli.SingleLineHandler{
								Name:     "verifyBiography",
								Handler:  personVerifyBiography,
								ArgNames: []string{"biography filename"},
							},
							&cli.SingleLineHandler{
								Name:     "verifyBiographyOmitted",
								Handler:  personVerifyBiographyOmitted,
								ArgNames: []string{},
							},
						},
					},
				),
			},
			&cli.Cli{
				FullDescription:    "Welcome to the journal commands.",
				OneLineDescription: "Journal",
				Name:               "journal",
				Handlers: append(cliAlexandria.CommonJournalHandlers,
					&cli.StructRunnerHandler{
						FullDescription:    "Welcome to the journal create dialog.",
						OneLineDescription: "Create journal",
						Name:               "createJournal",
						Action:             journalCreate,
					},
					&cli.StructRunnerHandler{
						FullDescription:              "Welcome to the journal update properties dialog",
						OneLineDescription:           "Update journal properties",
						Name:                         "updateProperties",
						ReferenceValueGetter:         journalUpdatePropertiesReference,
						ReferenceValueGetterArgNames: []string{"journal id"},
						Action:                       journalUpdateProperties,
					},
					&cli.Cli{
						FullDescription:    "Welcome to the journal description commands",
						OneLineDescription: "Update/Verify/Remove description",
						Name:               "description",
						Handlers: []cli.Handler{
							&cli.SingleLineHandler{
								Name:     "updateDescription",
								Handler:  journalUpdateDescription,
								ArgNames: []string{"journal id", "description file"},
							},
							&cli.SingleLineHandler{
								Name:     "removeDescription",
								Handler:  journalRemoveDescription,
								ArgNames: []string{"journal id"},
							},
							&cli.SingleLineHandler{
								Name:     "verifyDescription",
								Handler:  journalVerifyDescription,
								ArgNames: []string{"journal id", "description file"},
							},
							&cli.SingleLineHandler{
								Name:     "verifyDescriptionOmitted",
								Handler:  journalVerifyDescriptionOmitted,
								ArgNames: []string{"journalId"},
							},
						},
					},
					&cli.SingleLineHandler{
						Name:     "proposeEditor",
						Handler:  proposeEditor,
						ArgNames: []string{"journal id", "editor person id"},
					},
					&cli.SingleLineHandler{
						Name:     "acceptEditorship",
						Handler:  acceptEditorship,
						ArgNames: []string{"journal id"},
					},
					&cli.SingleLineHandler{
						Name:     "resignAsEditor",
						Handler:  resignAsEditor,
						ArgNames: []string{"journal id"},
					},
					&cli.StructRunnerHandler{
						FullDescription:    "Welcome to the volume create dialog",
						OneLineDescription: "Create volume",
						Name:               "createVolume",
						Action:             volumeCreate,
					},
				),
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

func personUpdateBiography(outputter cli.Outputter, biographyFileName string) {
	if !cliAlexandria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	data, err := ioutil.ReadFile(biographyFileName)
	if err != nil {
		outputter(cliAlexandria.ToIoError(err) + "\n")
		return
	}
	theCommand := command.GetCommandPersonUpdateBiography(
		cliAlexandria.LoggedInPerson.Id,
		cliAlexandria.LoggedInPerson.BiographyHash,
		data,
		cliAlexandria.LoggedInPerson.Id,
		cliAlexandria.LoggedIn(),
		cliAlexandria.Settings.PricePersonEdit)
	err = blockchain.SendCommand(theCommand, outputter)
	if err != nil {
		outputter("Error sending command to blockchain: " + err.Error() + "\n")
	}
}

func personRemoveBiography(outputter cli.Outputter) {
	if !cliAlexandria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	theCommand := command.GetCommandPersonOmitBiography(
		cliAlexandria.LoggedInPerson.Id,
		cliAlexandria.LoggedInPerson.BiographyHash,
		cliAlexandria.LoggedInPerson.Id,
		cliAlexandria.LoggedIn(),
		cliAlexandria.Settings.PricePersonEdit)
	err := blockchain.SendCommand(theCommand, outputter)
	if err != nil {
		outputter("Error sending command to blockchain: " + err.Error() + "\n")
	}
}

func personVerifyBiography(outputter cli.Outputter, biographyFileName string) {
	if !cliAlexandria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	data, err := ioutil.ReadFile(biographyFileName)
	if err != nil {
		outputter(cliAlexandria.ToIoError(err))
		return
	}
	err = dao.VerifyPersonBiography(cliAlexandria.LoggedInPerson.Id, data)
	if err != nil {
		outputter("Verification failed: " + err.Error() + "\n")
		return
	}
	outputter("Verified\n")
}

func personVerifyBiographyOmitted(outputter cli.Outputter) {
	if !cliAlexandria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	err := dao.VerifyPersonBiography(cliAlexandria.LoggedInPerson.Id, []byte{})
	if err != nil {
		outputter("Verification failed: " + err.Error() + "\n")
		return
	}
	outputter("Verified\n")
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

func journalUpdatePropertiesReference(outputter cli.Outputter, journalId string) *command.Journal {
	if !cliAlexandria.CheckBootstrappedAndKnownPerson(outputter) {
		return nil
	}
	daoJournal, err := dao.GetJournal(journalId)
	if err != nil {
		outputter(fmt.Sprintf("Journal does not exist: %s, detailed error message: %s\n",
			journalId, err.Error()))
		return nil
	}
	originalJournalId = journalId
	originalJournal = &command.Journal{
		Title: daoJournal.Title,
	}
	return originalJournal
}

var originalJournalId string
var originalJournal *command.Journal

func journalUpdateProperties(outputter cli.Outputter, journal *command.Journal) {
	theCommand := command.GetCommandJournalUpdateProperties(
		originalJournalId,
		originalJournal,
		journal,
		cliAlexandria.LoggedInPerson.Id,
		cliAlexandria.LoggedIn(),
		cliAlexandria.Settings.PriceEditorEditJournal)
	if err := blockchain.SendCommand(theCommand, outputter); err != nil {
		outputter(cliAlexandria.ToIoError(err))
	}
}

func journalUpdateDescription(outputter cli.Outputter, journalId, descriptionFileName string) {
	if !cliAlexandria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	data, err := ioutil.ReadFile(descriptionFileName)
	if err != nil {
		outputter(cliAlexandria.ToIoError(err) + "\n")
		return
	}
	origJournal, err := dao.GetJournal(journalId)
	if err != nil {
		outputter(fmt.Sprintf("Journal does not exist, journalId = %s, detailed error = %s\n",
			journalId, err.Error()))
		return
	}
	theCommand := command.GetCommandJournalUpdateDescription(
		journalId,
		origJournal.Descriptionhash,
		data,
		cliAlexandria.LoggedInPerson.Id,
		cliAlexandria.LoggedIn(),
		cliAlexandria.Settings.PriceEditorEditJournal)
	err = blockchain.SendCommand(theCommand, outputter)
	if err != nil {
		outputter("Error sending command to blockchain: " + err.Error() + "\n")
	}
}

func journalRemoveDescription(outputter cli.Outputter, journalId string) {
	if !cliAlexandria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	origJournal, err := dao.GetJournal(journalId)
	if err != nil {
		outputter(fmt.Sprintf("Journal does not exist, journalId = %s, detailed error = %s\n",
			journalId, err.Error()))
		return
	}
	theCommand := command.GetCommandJournalOmitDescription(
		journalId,
		origJournal.Descriptionhash,
		cliAlexandria.LoggedInPerson.Id,
		cliAlexandria.LoggedIn(),
		cliAlexandria.Settings.PriceEditorEditJournal)
	err = blockchain.SendCommand(theCommand, outputter)
	if err != nil {
		outputter("Error sending command to blockchain: " + err.Error() + "\n")
	}
}

func journalVerifyDescription(outputter cli.Outputter, journalId, descriptionFileName string) {
	data, err := ioutil.ReadFile(descriptionFileName)
	if err != nil {
		outputter(cliAlexandria.ToIoError(err))
		return
	}
	err = dao.VerifyJournalDescription(journalId, data)
	if err != nil {
		outputter("Verification failed: " + err.Error() + "\n")
		return
	}
	outputter("Verified\n")
}

func journalVerifyDescriptionOmitted(outputter cli.Outputter, journalId string) {
	err := dao.VerifyJournalDescription(journalId, []byte{})
	if err != nil {
		outputter("Verification failed: " + err.Error() + "\n")
		return
	}
	outputter("Verified\n")
}

func proposeEditor(outputter cli.Outputter, journalId, editorId string) {
	cliAlexandria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetCommandEditorInvite(
			journalId,
			editorId,
			cliAlexandria.LoggedInPerson.Id,
			cliAlexandria.LoggedIn(),
			cliAlexandria.Settings.PriceEditorAddColleague)
	})
}

func acceptEditorship(outputter cli.Outputter, journalId string) {
	cliAlexandria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetCommandEditorAcceptDuty(
			journalId,
			cliAlexandria.LoggedInPerson.Id,
			cliAlexandria.LoggedIn(),
			cliAlexandria.Settings.PriceEditorAcceptDuty)
	})
}

func resignAsEditor(outputter cli.Outputter, journalId string) {
	cliAlexandria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetCommandEditorResign(
			journalId,
			cliAlexandria.LoggedInPerson.Id,
			cliAlexandria.LoggedIn())
	})
}

func volumeCreate(outputter cli.Outputter, volume *command.Volume) {
	if !cliAlexandria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	cmd, volumeId := command.GetCommandVolumeCreate(
		volume,
		cliAlexandria.LoggedInPerson.Id,
		cliAlexandria.LoggedIn(),
		cliAlexandria.Settings.PriceEditorCreateVolume)
	if err := blockchain.SendCommand(cmd, outputter); err != nil {
		outputter(cliAlexandria.ToIoError(err) + "\n")
		return
	}
	outputter("The volumeId of the created volume is: " + volumeId + "\n")
}
