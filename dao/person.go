package dao

import (
	"database/sql"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/jmoiron/sqlx"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"strconv"
)

type Person struct {
	Id              string
	CreatedOn       int64  `db:"createdon"`
	ModifiedOn      int64  `db:"modifiedon"`
	PublicKey       string `db:"publickey"`
	Name            string
	Email           string
	IsMajor         bool `db:"ismajor"`
	IsSigned        bool `db:"issigned"`
	Balance         int32
	BiographyHash   string `db:"biographyhash"`
	BiographyFormat string `db:"biographyformat"`
	Organization    string
	Telephone       string
	Address         string
	PostalCode      string `db:"postalcode"`
	Country         string
	ExtraInfo       string `db:"extrainfo"`
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
	transactionId := ""
	eventSeq := int32(0)
	dataManipulation := &dataManipulationPersonCreate{}
	for _, attribute := range event.Attributes {
		switch attribute.Key {
		case model.TRANSACTION_ID:
			transactionId = attribute.Value
		case model.TIMESTAMP:
			i64, err = strconv.ParseInt(attribute.Value, 10, 64)
			dataManipulation.timestamp = i64
		case model.EVENT_SEQ:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			eventSeq = int32(i64)
		case model.ID:
			dataManipulation.id = attribute.Value
		case model.PERSON_PUBLIC_KEY:
			dataManipulation.publicKey = attribute.Value
		case model.PERSON_NAME:
			dataManipulation.name = attribute.Value
		case model.PERSON_EMAIL:
			dataManipulation.email = attribute.Value
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
}

var _ dataManipulation = new(dataManipulationPersonCreate)

func (dmpc *dataManipulationPersonCreate) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec(fmt.Sprintf("INSERT INTO person VALUES (%s)", GetPlaceHolders(17)),
		dmpc.id, dmpc.timestamp, dmpc.timestamp, dmpc.publicKey, dmpc.name,
		dmpc.email, false, false, int32(0), "",
		"", "", "", "", "",
		"", "")
	return err
}
