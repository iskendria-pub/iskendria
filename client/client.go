package main

import (
	"fmt"
	"github.com/iskendria-pub/iskendria/blockchain"
	"github.com/iskendria-pub/iskendria/cli"
	"github.com/iskendria-pub/iskendria/cliIskendria"
	"github.com/iskendria-pub/iskendria/command"
	"github.com/iskendria-pub/iskendria/dao"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var description = strings.TrimSpace(`
Welcome to the Iskendria Client Tool. Use this tool to
register and to manage manuscripts, reviews and journals.
`)

var makeGreen = "\033[32m"

func main() {
	context := &cli.Cli{
		FullDescription:    description,
		OneLineDescription: "Iskendria Client Tool",
		Name:               "iskendria-client",
		FormatEscape:       makeGreen,
		EventPager:         cliIskendria.PageEventStreamMessages,
		Handlers: append(cliIskendria.CommonRootHandlers,
			cliIskendria.CommonDiagnosticsGroup,
			&cli.Cli{
				FullDescription:    "Welcome to the settings commands",
				OneLineDescription: "Settings",
				Name:               "settings",
				Handlers:           cliIskendria.CommonSettingsHandlers,
			},
			&cli.Cli{
				FullDescription:    "Welcome to the person commands",
				OneLineDescription: "Person",
				Name:               "person",
				Handlers: append(cliIskendria.CommonPersonHandlers,
					&cli.StructRunnerHandler{
						FullDescription:              "Welcome to the person update dialog.",
						OneLineDescription:           "Update person",
						Name:                         "updatePerson",
						ReferenceValueGetter:         personUpdateReference,
						ReferenceValueGetterArgNames: []string{},
						Action:                       cliIskendria.PersonUpdate,
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
				Handlers: append(cliIskendria.CommonJournalHandlers,
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
			&cli.Cli{
				FullDescription:    "Welcome to the manuscript commands",
				OneLineDescription: "Manuscript",
				Name:               "manuscript",
				Handlers: []cli.Handler{
					&cli.StructRunnerHandler{
						FullDescription:    "Create new manuscript",
						OneLineDescription: "Create new manuscript",
						Name:               "create",
						Action:             manuscriptCreate,
					},
					&cli.StructRunnerHandler{
						FullDescription:    "Create a new manuscript version",
						OneLineDescription: "Create new manuscript version",
						Name:               "createNewVersion",
						Action:             manuscriptCreateNewVersion,
					},
					&cli.SingleLineHandler{
						Name:     "acceptAuthorship",
						Handler:  manuscriptAcceptAuthorship,
						ArgNames: []string{"manuscript id"},
					},
					&cli.SingleLineHandler{
						Name:     "allowReview",
						Handler:  manuscriptAllowReview,
						ArgNames: []string{"manuscript id"},
					},
					&cli.StructRunnerHandler{
						FullDescription:    "Add positive review about manuscript",
						OneLineDescription: "Add positive review about manuscript",
						Name:               "addReviewPositive",
						Action:             addPositiveReview,
					},
					&cli.StructRunnerHandler{
						FullDescription:    "Add negative review about manuscript",
						OneLineDescription: "Add negative review about manuscript",
						Name:               "addReviewNegative",
						Action:             addNegativeReview,
					},
					&cli.StructRunnerHandler{
						FullDescription:    "Publish manuscript",
						OneLineDescription: "Publish manuscript",
						Name:               "publish",
						Action:             manuscriptPublish,
					},
					&cli.StructRunnerHandler{
						FullDescription:    "Reject manuscript",
						OneLineDescription: "Reject manuscript",
						Name:               "reject",
						Action:             manuscriptReject,
					},
					&cli.StructRunnerHandler{
						FullDescription:    "Assign manuscript to volume",
						OneLineDescription: "Assign manuscript to volume",
						Name:               "assign",
						Action:             manuscriptAssign,
					},
				},
			},
		),
	}
	fmt.Print(makeGreen)
	dbLogger := log.New(os.Stdout, "db", log.Flags())
	dao.Init("client.db", dbLogger)
	defer dao.Shutdown(dbLogger)
	cliIskendria.InitEventStream("./client-events.log", "client")
	context.Run()
}

func personUpdateReference(outputter cli.Outputter) *dao.PersonUpdate {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return nil
	}
	cliIskendria.OriginalPersonId = cliIskendria.LoggedInPerson.Id
	cliIskendria.OriginalPerson = dao.PersonToPersonUpdate(cliIskendria.LoggedInPerson)
	return cliIskendria.OriginalPerson
}

func personUpdateBiography(outputter cli.Outputter, biographyFileName string) {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	data, err := ioutil.ReadFile(biographyFileName)
	if err != nil {
		outputter(cliIskendria.ToIoError(err) + "\n")
		return
	}
	theCommand := command.GetCommandPersonUpdateBiography(
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedInPerson.BiographyHash,
		data,
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PricePersonEdit)
	err = blockchain.SendCommand(theCommand, outputter)
	if err != nil {
		outputter("Error sending command to blockchain: " + err.Error() + "\n")
	}
}

func personRemoveBiography(outputter cli.Outputter) {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	theCommand := command.GetCommandPersonOmitBiography(
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedInPerson.BiographyHash,
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PricePersonEdit)
	err := blockchain.SendCommand(theCommand, outputter)
	if err != nil {
		outputter("Error sending command to blockchain: " + err.Error() + "\n")
	}
}

func personVerifyBiography(outputter cli.Outputter, biographyFileName string) {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	data, err := ioutil.ReadFile(biographyFileName)
	if err != nil {
		outputter(cliIskendria.ToIoError(err))
		return
	}
	err = dao.VerifyPersonBiography(cliIskendria.LoggedInPerson.Id, data)
	if err != nil {
		outputter("Verification failed: " + err.Error() + "\n")
		return
	}
	outputter("Verified\n")
}

func personVerifyBiographyOmitted(outputter cli.Outputter) {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	err := dao.VerifyPersonBiography(cliIskendria.LoggedInPerson.Id, []byte{})
	if err != nil {
		outputter("Verification failed: " + err.Error() + "\n")
		return
	}
	outputter("Verified\n")
}

func journalCreate(outputter cli.Outputter, journal *command.Journal) {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	cmd, journalId := command.GetCommandJournalCreate(
		journal,
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PriceEditorCreateJournal)
	if err := blockchain.SendCommand(cmd, outputter); err != nil {
		outputter(cliIskendria.ToIoError(err) + "\n")
		return
	}
	outputter("The journalId of the created journal is: " + journalId + "\n")
}

func journalUpdatePropertiesReference(outputter cli.Outputter, journalId string) *command.Journal {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
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
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PriceEditorEditJournal)
	if err := blockchain.SendCommand(theCommand, outputter); err != nil {
		outputter(cliIskendria.ToIoError(err))
	}
}

func journalUpdateDescription(outputter cli.Outputter, journalId, descriptionFileName string) {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	data, err := ioutil.ReadFile(descriptionFileName)
	if err != nil {
		outputter(cliIskendria.ToIoError(err) + "\n")
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
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PriceEditorEditJournal)
	err = blockchain.SendCommand(theCommand, outputter)
	if err != nil {
		outputter("Error sending command to blockchain: " + err.Error() + "\n")
	}
}

func journalRemoveDescription(outputter cli.Outputter, journalId string) {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
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
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PriceEditorEditJournal)
	err = blockchain.SendCommand(theCommand, outputter)
	if err != nil {
		outputter("Error sending command to blockchain: " + err.Error() + "\n")
	}
}

func journalVerifyDescription(outputter cli.Outputter, journalId, descriptionFileName string) {
	data, err := ioutil.ReadFile(descriptionFileName)
	if err != nil {
		outputter(cliIskendria.ToIoError(err))
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
	cliIskendria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetCommandEditorInvite(
			journalId,
			editorId,
			cliIskendria.LoggedInPerson.Id,
			cliIskendria.LoggedIn(),
			cliIskendria.Settings.PriceEditorAddColleague)
	})
}

func acceptEditorship(outputter cli.Outputter, journalId string) {
	cliIskendria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetCommandEditorAcceptDuty(
			journalId,
			cliIskendria.LoggedInPerson.Id,
			cliIskendria.LoggedIn(),
			cliIskendria.Settings.PriceEditorAcceptDuty)
	})
}

func resignAsEditor(outputter cli.Outputter, journalId string) {
	cliIskendria.SendCommandAsPerson(outputter, func() *command.Command {
		return command.GetCommandEditorResign(
			journalId,
			cliIskendria.LoggedInPerson.Id,
			cliIskendria.LoggedIn())
	})
}

func volumeCreate(outputter cli.Outputter, volume *command.Volume) {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	cmd, volumeId := command.GetCommandVolumeCreate(
		volume,
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PriceEditorCreateVolume)
	if err := blockchain.SendCommand(cmd, outputter); err != nil {
		outputter(cliIskendria.ToIoError(err) + "\n")
		return
	}
	outputter("The volumeId of the created volume is: " + volumeId + "\n")
}

func manuscriptCreate(outputter cli.Outputter, manuscriptCreate *ManuscriptCreate) {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	manuscriptData, err := ioutil.ReadFile(manuscriptCreate.ManuscriptFileName)
	if err != nil {
		outputter(cliIskendria.ToIoError(err))
	}
	cmd, manuscriptId := command.GetCommandManuscriptCreate(
		&command.ManuscriptCreate{
			TheManuscript: manuscriptData,
			CommitMsg:     manuscriptCreate.CommitMsg,
			Title:         manuscriptCreate.Title,
			AuthorId:      manuscriptCreate.AuthorId,
			JournalId:     manuscriptCreate.JournalId,
		},
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PriceAuthorSubmitNewManuscript)
	if err := blockchain.SendCommand(cmd, outputter); err != nil {
		outputter(cliIskendria.ToIoError(err))
	}
	outputter("The manuscriptId of the created manuscript is: " + manuscriptId + "\n")
}

type ManuscriptCreate struct {
	ManuscriptFileName string
	CommitMsg          string
	Title              string
	AuthorId           []string
	JournalId          string
}

func manuscriptCreateNewVersion(outputter cli.Outputter, manuscriptCreateNewVersion *ManuscriptCreateNewVersion) {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	manuscriptData, err := ioutil.ReadFile(manuscriptCreateNewVersion.ManuscriptFileName)
	if err != nil {
		outputter(cliIskendria.ToIoError(err))
		return
	}
	previousManuscript, err := dao.GetManuscript(manuscriptCreateNewVersion.PreviousManuscriptId)
	if err != nil {
		outputter(fmt.Sprintf("Invalid previous manuscript id: %s, error msg: %s",
			manuscriptCreateNewVersion.PreviousManuscriptId, err))
		return
	}
	threadReference, err := dao.GetReferenceThread(previousManuscript.ThreadId)
	if err != nil {
		outputter(fmt.Sprintf("Could not get list of manuscripts in thread %s: %s",
			previousManuscript.ThreadId, err.Error()))
		return
	}
	historicAuthors, err := dao.GetHistoricSignedAuthors(previousManuscript.ThreadId)
	if err != nil {
		outputter(fmt.Sprintf("Could not get list of historic signed authors for thead %s: %s",
			previousManuscript.ThreadId, err.Error()))
		return
	}
	cmd, manuscriptId := command.GetCommandManuscriptCreateNewVersion(
		&command.ManuscriptCreateNewVersion{
			TheManuscript:        manuscriptData,
			CommitMsg:            manuscriptCreateNewVersion.CommitMsg,
			Title:                manuscriptCreateNewVersion.Title,
			AuthorId:             manuscriptCreateNewVersion.AuthorId,
			PreviousManuscriptId: manuscriptCreateNewVersion.PreviousManuscriptId,
			ThreadId:             previousManuscript.ThreadId,
			JournalId:            previousManuscript.JournalId,
		},
		threadReference,
		historicAuthors,
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PriceAuthorSubmitNewVersion)
	if err := blockchain.SendCommand(cmd, outputter); err != nil {
		outputter(cliIskendria.ToIoError(err))
		return
	}
	outputter("The manuscriptId of the created manuscript is: " + manuscriptId + "\n")
}

type ManuscriptCreateNewVersion struct {
	ManuscriptFileName   string
	CommitMsg            string
	Title                string
	AuthorId             []string
	PreviousManuscriptId string
}

func manuscriptAcceptAuthorship(outputter cli.Outputter, manuscriptId string) {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	manuscript, err := dao.GetManuscript(manuscriptId)
	if err != nil {
		outputter(fmt.Sprintf("Unknown manuscript id: %s, error message: %s",
			manuscriptId, err.Error()))
		return
	}
	cmd := command.GetCommandManuscriptAcceptAuthorship(
		manuscript,
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PriceAuthorAcceptAuthorship)
	if err := blockchain.SendCommand(cmd, outputter); err != nil {
		outputter(cliIskendria.ToIoError(err))
		return
	}
}

func manuscriptAllowReview(outputter cli.Outputter, manuscriptId string) {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	manuscript, err := dao.GetManuscript(manuscriptId)
	if err != nil {
		outputter(fmt.Sprintf("Unknown manuscript id: %s, error message: %s",
			manuscriptId, err.Error()))
		return
	}
	referenceThread, err := dao.GetReferenceThread(manuscript.ThreadId)
	if err != nil {
		outputter(fmt.Sprintf("Error getting version history of manuscript: %s", err.Error()))
		return
	}
	cmd := command.GetCommandManuscriptAllowReview(
		manuscript.ThreadId,
		referenceThread,
		manuscript.JournalId,
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PriceEditorAllowManuscriptReview)
	if err := blockchain.SendCommand(cmd, outputter); err != nil {
		outputter(cliIskendria.ToIoError(err))
		return
	}
}

func addPositiveReview(outputter cli.Outputter, r *ReviewCreation) {
	addReview(outputter, r, getCommandReviewSubmitPositive)
}

func addNegativeReview(outputter cli.Outputter, r *ReviewCreation) {
	addReview(outputter, r, getCommandReviewSubmitNegative)
}

func addReview(outputter cli.Outputter, r *ReviewCreation, commandCreator reviewCreatorType) {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	cr, err := getCommandReviewCreate(r)
	if err != nil {
		outputter(err.Error() + "\n")
	}
	cmd, reviewId := commandCreator(cr)
	err = blockchain.SendCommand(cmd, outputter)
	if err != nil {
		outputter(cliIskendria.ToIoError(err))
		return
	}
	outputter(fmt.Sprintf("The id of the created review is: %s\n", reviewId))
}

type ReviewCreation struct {
	ManuscriptId string
	FileName     string
}

type reviewCreatorType func(*command.ReviewCreate) (*command.Command, string)

func getCommandReviewCreate(r *ReviewCreation) (*command.ReviewCreate, error) {
	reviewData, err := ioutil.ReadFile(r.FileName)
	if err != nil {
		return nil, err
	}
	return &command.ReviewCreate{
		ManuscriptId: r.ManuscriptId,
		TheReview:    reviewData,
	}, nil
}

func getCommandReviewSubmitPositive(cr *command.ReviewCreate) (*command.Command, string) {
	return command.GetCommandWritePositiveReview(
		cr,
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PriceReviewerSubmit)
}

func getCommandReviewSubmitNegative(cr *command.ReviewCreate) (*command.Command, string) {
	return command.GetCommandWriteNegativeReview(
		cr,
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PriceReviewerSubmit)
}

func manuscriptPublish(outputter cli.Outputter, judge *command.ManuscriptJudge) {
	manuscriptJudge(outputter, judge, getPositiveJudgeCommand)
}

func manuscriptReject(outputter cli.Outputter, judge *command.ManuscriptJudge) {
	manuscriptJudge(outputter, judge, getNegativeJudgeCommand)
}

func manuscriptJudge(outputter cli.Outputter, judge *command.ManuscriptJudge, commandGetter judgeCommandGetter) {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	manuscript, err := dao.GetManuscript(judge.ManuscriptId)
	if err != nil {
		outputter(err.Error() + "\n")
		return
	}
	cmd := commandGetter(judge, manuscript.JournalId)
	err = blockchain.SendCommand(cmd, outputter)
	if err != nil {
		outputter(cliIskendria.ToIoError(err))
		return
	}
}

type judgeCommandGetter func(*command.ManuscriptJudge, string) *command.Command

func getPositiveJudgeCommand(judge *command.ManuscriptJudge, journalId string) *command.Command {
	return command.GetCommandManuscriptPublish(
		judge,
		journalId,
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PriceEditorPublishManuscript)
}

func getNegativeJudgeCommand(judge *command.ManuscriptJudge, journalId string) *command.Command {
	return command.GetCommandManuscriptReject(
		judge,
		journalId,
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PriceEditorRejectManuscript)
}

func manuscriptAssign(outputter cli.Outputter, manuscriptAssign *command.ManuscriptAssign) {
	if !cliIskendria.CheckBootstrappedAndKnownPerson(outputter) {
		return
	}
	manuscript, err := dao.GetManuscript(manuscriptAssign.ManuscriptId)
	if err != nil {
		outputter(err.Error() + "\n")
		return
	}
	cmd := command.GetCommandManuscriptAssign(
		manuscriptAssign,
		manuscript.JournalId,
		cliIskendria.LoggedInPerson.Id,
		cliIskendria.LoggedIn(),
		cliIskendria.Settings.PriceEditorAssignManuscript)
	err = blockchain.SendCommand(cmd, outputter)
	if err != nil {
		outputter(cliIskendria.ToIoError(err))
	}
}
