package dao

import (
	"github.com/jmoiron/sqlx"
	"strings"
)

const THE_SETTINGS_ID = 1

var db *sqlx.DB

func GetPlaceHolders(n int) string {
	placeHolders := make([]string, n)
	for i := 0; i < n; i++ {
		placeHolders[i] = "?"
	}
	return strings.Join(placeHolders, ", ")
}
