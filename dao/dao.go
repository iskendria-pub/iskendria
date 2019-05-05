package dao

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"gitlab.bbinfra.net/3estack/alexandria/util"
	"log"
	"os"
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
}

func Shutdown(logger *log.Logger) {
	util.CloseDb(db, logger)
}

func ShutdownAndDelete(logger *log.Logger) {
	util.CloseDb(db, logger)
	err := os.Remove(DbFileName)
	if err != nil {
		logger.Fatal("Could not remove database file: " + err.Error())
	}
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
