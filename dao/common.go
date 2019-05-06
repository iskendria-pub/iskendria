package dao

import (
	"github.com/jmoiron/sqlx"
	"strings"
)

const THE_SETTINGS_ID = 1

type event interface{
	accept(context) error
}

type context interface {
	visitDataManipulation(transactionId string, eventSeq int32, dataManipulation dataManipulation) error
	visitTransactionControl(transactionId string, numEvents int32) error
	visitBlockControl(currentBlockId, previousBlockId string) error
}

type dataManipulation interface {
	apply(*sqlx.Tx) error
}

var db *sqlx.DB

func GetPlaceHolders(n int) string {
	placeHolders := make([]string, n)
	for i := 0; i < n; i++ {
		placeHolders[i] = "?"
	}
	return strings.Join(placeHolders, ", ")
}
