package model

import "github.com/google/uuid"

var TableCreatePerson = `
create table Person (
	id string primary key not null,
	createdOn integer not null,
	modifiedOn integer not null,
	publicKey string not null,
	name string not null,
	email string not null,
	isMajor bool not null,
	isSigned bool not null,
	balance integer not null,
	biographyHash string not null,
	biographyFormat string not null,
	organization string not null,
	telephone string not null,
	address string not null,
	postalCode string not null,
	country string not null,
	extraInfo string not null,
)`

const (
	PERSON_PUBLIC_KEY       = "publicKey"
	PERSON_NAME             = "name"
	PERSON_EMAIL            = "email"
	PERSON_IS_MAJOR         = "isMajor"
	PERSON_IS_SIGNED        = "isSigned"
	PERSON_BALANCE          = "balance"
	PERSON_BIOGRAPHY_HASH   = "biographyHash"
	PERSON_BIOGRAPHY_FORMAT = "biographyFormat"
	PERSON_ORGANIZATION     = "organization"
	PERSON_TELEPHONE        = "telephone"
	PERSON_ADDRESS          = "address"
	PERSON_POSTAL_CODE      = "postalCode"
	PERSON_COUNTRY          = "country"
	PERSON_EXTRA_INFO       = "extraInfo"
)

const personAddressPrefix = "01"

func CreatePersonAddress() string {
	var uuid uuid.UUID = uuid.New()
	uuidDigest := hexdigestOfUuid(uuid)
	return Namespace + personAddressPrefix + uuidDigest[:62]
}

func IsPersonAddress(address string) bool {
	return getAddressPrefixFromAddress(address) == personAddressPrefix
}
