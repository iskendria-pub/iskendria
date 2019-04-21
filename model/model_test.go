package model

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.bbinfra.net/3estack/alexandria/util"
	"log"
	"os"
	"testing"
)

func TestCreatePersonTable(t *testing.T) {
	logger := getLogger()
	dbFname := "./testPerson.db"
	db, err := sql.Open("sqlite3", dbFname)
	if err != nil {
		t.Error(err)
	}
	defer removeTestDatabase(db, dbFname, logger)
	_, err = db.Exec(CreatePersonTableSql)
	if err != nil {
		t.Error(err)
	}
}

func getLogger() *log.Logger {
	return log.New(os.Stdout, "model-test", log.Lshortfile)
}

func removeTestDatabase(db *sql.DB, dbFname string, logger *log.Logger) {
	util.CloseDb(db, logger)
	removeFile(dbFname, logger)
}

func removeFile(fname string, logger *log.Logger) {
	err := os.Remove(fname)
	if err != nil {
		logger.Println(err)
	}
}
