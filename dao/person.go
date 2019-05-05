package dao

import (
	"database/sql"
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
	persons := []Person{}
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
