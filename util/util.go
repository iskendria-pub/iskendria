package util

import (
	"database/sql"
	"log"
	"strings"
)

func UnTitle(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

func CloseDb(db *sql.DB, logger *log.Logger) {
	err := db.Close()
	if err != nil {
		logger.Println("Could not close Sqlite 3 database: " + err.Error())
	}
}
