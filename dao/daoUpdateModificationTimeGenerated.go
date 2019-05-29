package dao

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/jmoiron/sqlx"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"strconv"
)

func createSettingsModificationTimeEvent(input *events_pb2.Event) (event, error) {
	dm := new(dataManipulationSettingsModificationTime)
	dm.id = THE_SETTINGS_ID
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
			err = errors.New("createSettingsModificationTimeEvent: Unknown event attribute: " + a.Key)
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationSettingsModificationTime struct {
	id        string
	timestamp int64
}

var _ dataManipulation = new(dataManipulationSettingsModificationTime)

func (dm *dataManipulationSettingsModificationTime) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec(fmt.Sprintf("UPDATE settings SET modifiedon = %d WHERE id = \"%s\"",
		dm.timestamp, dm.id))
	return err
}

func createPersonModificationTimeEvent(input *events_pb2.Event) (event, error) {
	dm := new(dataManipulationPersonModificationTime)

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
			err = errors.New("createPersonModificationTimeEvent: Unknown event attribute: " + a.Key)
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationPersonModificationTime struct {
	id        string
	timestamp int64
}

var _ dataManipulation = new(dataManipulationPersonModificationTime)

func (dm *dataManipulationPersonModificationTime) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec(fmt.Sprintf("UPDATE person SET modifiedon = %d WHERE id = \"%s\"",
		dm.timestamp, dm.id))
	return err
}

func createJournalModificationTimeEvent(input *events_pb2.Event) (event, error) {
	dm := new(dataManipulationJournalModificationTime)

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
			err = errors.New("createJournalModificationTimeEvent: Unknown event attribute: " + a.Key)
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationJournalModificationTime struct {
	id        string
	timestamp int64
}

var _ dataManipulation = new(dataManipulationJournalModificationTime)

func (dm *dataManipulationJournalModificationTime) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec(fmt.Sprintf("UPDATE journal SET modifiedon = %d WHERE journalId = \"%s\"",
		dm.timestamp, dm.id))
	return err
}
