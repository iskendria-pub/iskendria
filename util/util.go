package util

import (
	"database/sql"
	"log"
)

func CloseDb(db *sql.DB, logger *log.Logger) {
	err := db.Close()
	if err != nil {
		logger.Println("Could not close Sqlite 3 database: " + err.Error())
	}
}
