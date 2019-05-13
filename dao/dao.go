package dao

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"gitlab.bbinfra.net/3estack/alexandria/util"
	"log"
	"strings"
)

var DbFileName string

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
	}
	for _, stmt := range tableCreateStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			logger.Fatal(fmt.Sprintf("Table create statement failed: %s, error: %s",
				stmt, err))
		}
	}
}

func HandleEvent(input *events_pb2.Event) error {
	ev, err := parseEvent(input)
	if err != nil {
		return err
	}
	return ev.accept(theContext)
}

func parseEvent(input *events_pb2.Event) (event, error) {
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
	case model.EV_TYPE_PERSON_CREATE:
		return createPersonCreateEvent(input)
	case model.EV_TYPE_PERSON_UPDATE:
		return createPersonUpdateEvent(input)
	case model.EV_TYPE_PERSON_MODIFICATION_TIME:
		return createPersonModificationTimeEvent(input)
	default:
		return nil, errors.New("Unknown event type: " + input.EventType)
	}
}

func simplifyEventType(orig string) string {
	components := strings.Split(orig, "/")
	if components[0] == model.FamilyName {
		return components[1]
	}
	return orig
}
