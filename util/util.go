package util

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"strings"
)

func UnTitle(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

func CloseDb(db *sqlx.DB, logger *log.Logger) {
	err := db.Close()
	if err != nil {
		logger.Println("Could not close Sqlite 3 database: " + err.Error())
	}
}

func RemoveFileIfExists(fname string, logger *log.Logger) {
	_, err := os.Stat(fname)
	if err == nil {
		RemoveExistingFile(fname, logger)
		return
	}
	if !os.IsNotExist(err) {
		logger.Fatal(fmt.Sprintf("Something is wrong with database file: %s, error: %s",
			fname, err.Error()))
	}
}

func RemoveExistingFile(fname string, logger *log.Logger) {
	err := os.Remove(fname)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Could not remove existing database file: %s, error: %s",
			fname, err.Error()))
	}
}

func CheckPanicked(panicer func()) (didPanic bool) {
	didPanic = false
	defer func() {
		if recover() != nil {
			didPanic = true
		}
	}()
	panicer()
	return
}

// Does not work for -1 << 63
func Abs(i int64) int64 {
	if i >= 0 {
		return i
	}
	if i == (-1 << 63) {
		panic(fmt.Sprintf("Inverse is out of range: %d", i))
	}
	return -i
}
