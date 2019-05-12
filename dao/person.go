package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/jmoiron/sqlx"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"strconv"
	"strings"
)

type Person struct {
	Id            string
	CreatedOn     int64  `db:"createdon"`
	ModifiedOn    int64  `db:"modifiedon"`
	PublicKey     string `db:"publickey"`
	Name          string
	Email         string
	IsMajor       bool `db:"ismajor"`
	IsSigned      bool `db:"issigned"`
	Balance       int32
	BiographyHash string `db:"biographyhash"`
	Organization  string
	Telephone     string
	Address       string
	PostalCode    string `db:"postalcode"`
	Country       string
	ExtraInfo     string `db:"extrainfo"`
}

func SearchPersonByKey(key string) ([]*Person, error) {
	persons := make([]Person, 0)
	err := db.Select(&persons, "SELECT * FROM person WHERE publickey = ?", key)
	if err != nil {
		return nil, err
	}
	result := make([]*Person, len(persons))
	for i := 0; i < len(persons); i++ {
		result[i] = &persons[i]
	}
	return result, nil
}

func GetPersonById(id string) (*Person, error) {
	var person = new(Person)
	err := db.QueryRowx("SELECT * FROM person WHERE id = ?", id).StructScan(person)
	if err == nil {
		return person, nil
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nil, err
}

type PersonUpdate struct {
	PublicKey     string
	Name          string
	Email         string
	BiographyHash string
	Organization  string
	Telephone     string
	Address       string
	PostalCode    string
	Country       string
	ExtraInfo     string
}

func PersonToPersonUpdate(p *Person) *PersonUpdate {
	return &PersonUpdate{
		PublicKey:     p.PublicKey,
		Name:          p.Name,
		Email:         p.Email,
		BiographyHash: p.BiographyHash,
		Organization:  p.Organization,
		Telephone:     p.Telephone,
		Address:       p.Address,
		PostalCode:    p.PostalCode,
		Country:       p.Country,
		ExtraInfo:     p.ExtraInfo,
	}
}

func createPersonCreateEvent(event *events_pb2.Event) (event, error) {
	var err error
	var i64 int64
	var b bool
	transactionId := ""
	eventSeq := int32(0)
	dataManipulation := &dataManipulationPersonCreate{}
	for _, attribute := range event.Attributes {
		switch attribute.Key {
		case model.EV_KEY_TRANSACTION_ID:
			transactionId = attribute.Value
		case model.EV_KEY_TIMESTAMP:
			i64, err = strconv.ParseInt(attribute.Value, 10, 64)
			dataManipulation.timestamp = i64
		case model.EV_KEY_EVENT_SEQ:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			eventSeq = int32(i64)
		case model.EV_KEY_ID:
			dataManipulation.id = attribute.Value
		case model.EV_KEY_PERSON_PUBLIC_KEY:
			dataManipulation.publicKey = attribute.Value
		case model.EV_KEY_PERSON_NAME:
			dataManipulation.name = attribute.Value
		case model.EV_KEY_PERSON_EMAIL:
			dataManipulation.email = attribute.Value
		case model.EV_KEY_PERSON_IS_MAJOR:
			b, err = strconv.ParseBool(attribute.Value)
			dataManipulation.isMajor = b
		case model.EV_KEY_PERSON_IS_SIGNED:
			b, err = strconv.ParseBool(attribute.Value)
			dataManipulation.isSigned = b
		}
		if err != nil {
			return nil, err
		}
	}
	return &dataManipulationEvent{
		transactionId:    transactionId,
		eventSeq:         eventSeq,
		dataManipulation: dataManipulation,
	}, nil
}

type dataManipulationPersonCreate struct {
	timestamp int64
	id        string
	publicKey string
	name      string
	email     string
	isMajor   bool
	isSigned  bool
}

var _ dataManipulation = new(dataManipulationPersonCreate)

func (dmpc *dataManipulationPersonCreate) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec(fmt.Sprintf("INSERT INTO person VALUES (%s)", GetPlaceHolders(16)),
		dmpc.id, dmpc.timestamp, dmpc.timestamp, dmpc.publicKey, dmpc.name,
		dmpc.email, dmpc.isMajor, dmpc.isSigned, int32(0), "",
		"", "", "", "", "",
		"")
	return err
}

func createPersonUpdateEvent(ev *events_pb2.Event) (event, error) {
	dmProperties := &dataManipulationPersonUpdateProperties{}
	dmAuthorization := &dataManipulationPersonUpdateAuthorization{}
	dmBalance := &dataManipulationPersonUpdateBalance{}
	result := &dataManipulationEvent{}
	for _, a := range ev.Attributes {
		var err error
		var i64 int64
		var b bool
		switch a.Key {
		case model.EV_KEY_TRANSACTION_ID:
			result.transactionId = a.Value
		case model.EV_KEY_EVENT_SEQ:
			i64, err = strconv.ParseInt(a.Value, 10, 32)
			result.eventSeq = int32(i64)
		case model.EV_KEY_TIMESTAMP:
			// Nothing to do.
		case model.EV_KEY_ID:
			dmProperties.id = a.Value
			dmAuthorization.id = a.Value
			dmBalance.id = a.Value
		case model.EV_KEY_PERSON_PUBLIC_KEY, model.EV_KEY_PERSON_NAME, model.EV_KEY_PERSON_EMAIL,
			model.EV_KEY_PERSON_BIOGRAPHY_HASH, model.EV_KEY_PERSON_ORGANIZATION, model.EV_KEY_PERSON_TELEPHONE,
			model.EV_KEY_PERSON_ADDRESS, model.EV_KEY_PERSON_POSTAL_CODE, model.EV_KEY_PERSON_COUNTRY,
			model.EV_KEY_PERSON_EXTRA_INFO:
			result.dataManipulation = dmProperties
			dmProperties.field = strings.ToLower(a.Key)
			dmProperties.newValue = a.Value
		case model.EV_KEY_PERSON_IS_MAJOR, model.EV_KEY_PERSON_IS_SIGNED:
			result.dataManipulation = dmAuthorization
			dmAuthorization.field = strings.ToLower(a.Key)
			b, err = strconv.ParseBool(a.Value)
			dmAuthorization.newValue = b
		case model.EV_KEY_PERSON_BALANCE:
			result.dataManipulation = dmBalance
			i64, err = strconv.ParseInt(a.Value, 10, 32)
			dmBalance.newValue = int32(i64)
		default:
			err = errors.New("createPersonUpdateEvent: unknown attribute " + a.Key)
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type dataManipulationPersonUpdateProperties struct {
	id       string
	field    string
	newValue string
}

var _ dataManipulation = new(dataManipulationPersonUpdateProperties)

func (dm *dataManipulationPersonUpdateProperties) apply(tx *sqlx.Tx) error {
	query := fmt.Sprintf("UPDATE person SET %s = \"%s\" WHERE id = \"%s\"",
		dm.field, dm.newValue, dm.id)
	_, err := tx.Exec(query)
	return err
}

type dataManipulationPersonUpdateAuthorization struct {
	id       string
	field    string
	newValue bool
}

var _ dataManipulation = new(dataManipulationPersonUpdateAuthorization)

func (dm *dataManipulationPersonUpdateAuthorization) apply(tx *sqlx.Tx) error {
	query := fmt.Sprintf("UPDATE person SET %s = %s WHERE id = \"%s\"",
		dm.field, strconv.FormatBool(dm.newValue), dm.id)
	_, err := tx.Exec(query)
	return err
}

type dataManipulationPersonUpdateBalance struct {
	id       string
	newValue int32
}

var _ dataManipulation = new(dataManipulationPersonUpdateBalance)

func (dm *dataManipulationPersonUpdateBalance) apply(tx *sqlx.Tx) error {
	query := fmt.Sprintf("UPDATE person SET %s = %d WHERE id = \"%s\"",
		strings.ToLower(model.EV_KEY_PERSON_BALANCE), dm.newValue, dm.id)
	_, err := tx.Exec(query)
	return err
}

func createPersonModificationTimeEvent(ev *events_pb2.Event) (event, error) {
	dm := &dataManipulationPersonModificationTime{}
	result := &dataManipulationEvent{
		dataManipulation: dm,
	}
	var err error
	var i64 int64
	for _, a := range ev.Attributes {
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
