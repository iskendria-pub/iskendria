package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

var templateModificationTime = strings.TrimSpace(`
package dao

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/jmoiron/sqlx"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"strconv"
)

{{range .}}
func create{{.Tag}}ModificationTimeEvent(input *events_pb2.Event) (event, error) {
	dm := new(dataManipulation{{.Tag}}ModificationTime)
    {{if .IsSettings}} dm.id = THE_SETTINGS_ID {{end}}
	result := &dataManipulationEvent{
		dataManipulation: dm,
	}
	var err error
	var i64 int64
	for _, a := range input.Attributes {
		switch a.Key {
		case model.EV_KEY_TRANSACTION_ID:
			result.transactionId = a.Value
		case model.EV_KEY_EVENT_SEQ:
			i64, err = strconv.ParseInt(a.Value, 10, 32)
			result.eventSeq = int32(i64)
		case model.EV_KEY_TIMESTAMP:
			i64, err = strconv.ParseInt(a.Value, 10, 64)
			dm.timestamp = i64
        case model.EV_KEY_ID:
            dm.id = a.Value
		default:
			err = errors.New("create{{.Tag}}ModificationTimeEvent: Unknown event attribute: " + a.Key)
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulation{{.Tag}}ModificationTime struct {
	id string
	timestamp int64
}

var _ dataManipulation = new(dataManipulation{{.Tag}}ModificationTime)

func (dm *dataManipulation{{.Tag}}ModificationTime) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec(fmt.Sprintf("UPDATE {{.Table}} SET modifiedon = %d WHERE {{.IdColumn}} = \"%s\"",
		dm.timestamp, dm.id))
	return err
}
{{end}}
`)

type ModificationTimeUpdate struct {
	IdColumn   string
	IsSettings bool
	Tag        string
	Table      string
}

func main() {
	c := []*ModificationTimeUpdate{
		{
			IdColumn:   "id",
			IsSettings: true,
			Tag:        "Settings",
			Table:      "settings",
		},
		{
			IdColumn:   "id",
			IsSettings: false,
			Tag:        "Person",
			Table:      "person",
		},
		{
			IdColumn:   "journalId",
			IsSettings: false,
			Tag:        "Journal",
			Table:      "journal",
		},
	}
	tmpl, err := template.New("templateModificationTime").Parse(templateModificationTime)
	if err != nil {
		fmt.Println("Error parsing template")
		fmt.Println(err)
		return
	}
	err = tmpl.Execute(os.Stdout, c)
	if err != nil {
		fmt.Println(err)
	}
}