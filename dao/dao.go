package dao

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/iskendria-pub/iskendria/model"
	"github.com/iskendria-pub/iskendria/util"
	"log"
	"strings"
)

var DbFileName string

var AllEventTypes = []string{
	model.SAWTOOTH_BLOCK_COMMIT,
	model.AlexandriaPrefix + model.EV_TYPE_TRANSACTION_CONTROL,
	model.AlexandriaPrefix + model.EV_TYPE_SETTINGS_CREATE,
	model.AlexandriaPrefix + model.EV_TYPE_SETTINGS_UPDATE,
	model.AlexandriaPrefix + model.EV_TYPE_SETTINGS_MODIFICATION_TIME,
	model.AlexandriaPrefix + model.EV_TYPE_JOURNAL_CREATE,
	model.AlexandriaPrefix + model.EV_TYPE_JOURNAL_UPDATE,
	model.AlexandriaPrefix + model.EV_TYPE_JOURNAL_MODIFICATION_TIME,
	model.AlexandriaPrefix + model.EV_TYPE_EDITOR_CREATE,
	model.AlexandriaPrefix + model.EV_TYPE_EDITOR_UPDATE,
	model.AlexandriaPrefix + model.EV_TYPE_EDITOR_DELETE,
	model.AlexandriaPrefix + model.EV_TYPE_VOLUME_CREATE,
	model.AlexandriaPrefix + model.EV_TYPE_PERSON_CREATE,
	model.AlexandriaPrefix + model.EV_TYPE_PERSON_UPDATE,
	model.AlexandriaPrefix + model.EV_TYPE_PERSON_MODIFICATION_TIME,
	model.AlexandriaPrefix + model.EV_TYPE_MANUSCRIPT_CREATE,
	model.AlexandriaPrefix + model.EV_TYPE_AUTHOR_CREATE,
	model.AlexandriaPrefix + model.EV_TYPE_MANUSCRIPT_UPDATE,
	model.AlexandriaPrefix + model.EV_TYPE_AUTHOR_UPDATE,
	model.AlexandriaPrefix + model.EV_TYPE_MANUSCRIPT_MODIFICATION_TIME,
	model.AlexandriaPrefix + model.EV_TYPE_MANUSCRIPT_THREAD_UPDATE,
	model.AlexandriaPrefix + model.EV_TYPE_REVIEW_CREATE,
	model.AlexandriaPrefix + model.EV_TYPE_REVIEW_USE_BY_EDITOR,
}

func Init(fname string, logger *log.Logger) {
	var err error
	DbFileName = fname
	util.RemoveFileIfExists(fname, logger)
	db, err = sqlx.Open("sqlite3", fname)
	if err != nil {
		logger.Fatal("Could not open sqlite3 database: " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		logger.Fatal("Could not ping database: " + err.Error())
	}
	createTables(logger)
	createContext()
}

func Shutdown(logger *log.Logger) {
	util.CloseDb(db, logger)
}

func ShutdownAndDelete(logger *log.Logger) {
	util.CloseDb(db, logger)
	util.RemoveExistingFile(DbFileName, logger)
	util.RemoveFileIfExists(DbFileName+"-journal", logger)
}

func createTables(logger *log.Logger) {
	tableCreateStatements := []string{
		model.TableCreateSettings,
		model.TableCreatePerson,
		model.TableCreateJournal,
		model.TableCreateEditor,
		model.TableCreateVolume,
		model.TableCreateManuscript,
		model.IndexCreateManuscript,
		model.TableCreateAuthor,
		model.TableCreateReview,
	}
	for _, stmt := range tableCreateStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			logger.Fatal(fmt.Sprintf("Table create statement failed: %s, error: %s",
				stmt, err))
		}
	}
}

func HandleEvent(input *events_pb2.Event, logger *log.Logger) error {
	ev, err := parseEvent(input, logger)
	if err != nil {
		return err
	}
	return ev.accept(theContext)
}

type EventHandler func(*events_pb2.Event, *log.Logger) error

var _ EventHandler = HandleEvent

func parseEvent(input *events_pb2.Event, logger *log.Logger) (event, error) {
	logEvent(input, logger)
	defer logger.Printf("Done parsing event of type %s\n", input.EventType)
	switch simplifyEventType(input.EventType) {
	case model.SAWTOOTH_BLOCK_COMMIT:
		return createSawtoothBlockCommitEvent(input)
	case model.EV_TYPE_TRANSACTION_CONTROL:
		return createTransactionControlEvent(input)
	case model.EV_TYPE_SETTINGS_CREATE:
		return createSettingsCreateEvent(input)
	case model.EV_TYPE_SETTINGS_UPDATE:
		return createSettingsUpdateEvent(input)
	case model.EV_TYPE_SETTINGS_MODIFICATION_TIME:
		return createSettingsModificationTimeEvent(input)
	case model.EV_TYPE_JOURNAL_CREATE:
		return createJournalCreateEvent(input, logger)
	case model.EV_TYPE_JOURNAL_UPDATE:
		return createJournalUpdateEvent(input)
	case model.EV_TYPE_JOURNAL_MODIFICATION_TIME:
		return createJournalModificationTimeEvent(input)
	case model.EV_TYPE_EDITOR_CREATE:
		return createEditorCreateEvent(input, logger)
	case model.EV_TYPE_EDITOR_DELETE:
		return createEditorDeleteEvent(input, logger)
	case model.EV_TYPE_EDITOR_UPDATE:
		return createEditorUpdateEvent(input, logger)
	case model.EV_TYPE_PERSON_CREATE:
		return createPersonCreateEvent(input)
	case model.EV_TYPE_VOLUME_CREATE:
		return createVolumeCreateEvent(input)
	case model.EV_TYPE_PERSON_UPDATE:
		return createPersonUpdateEvent(input)
	case model.EV_TYPE_PERSON_MODIFICATION_TIME:
		return createPersonModificationTimeEvent(input)
	case model.EV_TYPE_MANUSCRIPT_CREATE:
		return createManuscriptCreateEvent(input)
	case model.EV_TYPE_AUTHOR_CREATE:
		return createAuthorCreateEvent(input)
	case model.EV_TYPE_MANUSCRIPT_UPDATE:
		return createManuscriptUpdateEvent(input)
	case model.EV_TYPE_AUTHOR_UPDATE:
		return createAuthorUpdateEvent(input)
	case model.EV_TYPE_MANUSCRIPT_MODIFICATION_TIME:
		return createManuscriptModificationTimeEvent(input)
	case model.EV_TYPE_MANUSCRIPT_THREAD_UPDATE:
		return createManuscriptThreadUpdateEvent(input)
	case model.EV_TYPE_REVIEW_CREATE:
		return createReviewCreateEvent(input)
	case model.EV_TYPE_REVIEW_USE_BY_EDITOR:
		return createReviewUseByEditorEvent(input)
	default:
		return nil, errors.New("Unknown event type: " + input.EventType)
	}
}

func logEvent(ev *events_pb2.Event, logger *log.Logger) {
	var transactionId string
	var eventSeq string
	for _, a := range ev.Attributes {
		switch a.Key {
		case model.EV_KEY_TRANSACTION_ID:
			transactionId = a.Value
		case model.EV_KEY_EVENT_SEQ:
			eventSeq = a.Value
		}
	}
	logger.Printf("Parsing event with type %s, transactionId = %s, eventSeq = %s\n",
		ev.EventType, transactionId, eventSeq)
}

func simplifyEventType(orig string) string {
	components := strings.Split(orig, "/")
	if components[0] == model.FamilyName {
		return components[1]
	}
	return orig
}
