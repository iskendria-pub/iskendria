package model

import "github.com/google/uuid"

var TableCreatePerson = `
CREATE TABLE person (
	id varchar primary key not null,
	createdon integer not null,
	modifiedon integer not null,
	publickey varchar not null,
	name varchar not null,
	email varchar not null,
	ismajor bool not null,
	issigned bool not null,
	balance integer not null,
	biographyhash varchar not null,
	organization varchar not null,
	telephone varchar not null,
	address varchar not null,
	postalcode varchar not null,
	country varchar not null,
	extrainfo varchar not null
)`

const (
	EV_PERSON_CREATE            = "evPersonCreate"
	EV_PERSON_UPDATE            = "evPersonUpdate"
	EV_PERSON_MODIFICATION_TIME = "evPersonModificationTime"
)

const (
	PERSON_PUBLIC_KEY     = "publicKey"
	PERSON_NAME           = "name"
	PERSON_EMAIL          = "email"
	PERSON_IS_MAJOR       = "isMajor"
	PERSON_IS_SIGNED      = "isSigned"
	PERSON_BALANCE        = "balance"
	PERSON_BIOGRAPHY_HASH = "biographyHash"
	PERSON_ORGANIZATION   = "organization"
	PERSON_TELEPHONE      = "telephone"
	PERSON_ADDRESS        = "address"
	PERSON_POSTAL_CODE    = "postalCode"
	PERSON_COUNTRY        = "country"
	PERSON_EXTRA_INFO     = "extraInfo"
)

const personAddressPrefix = "01"

func CreatePersonAddress() string {
	var theUuid uuid.UUID = uuid.New()
	uuidDigest := hexdigestOfUuid(theUuid)
	return Namespace + personAddressPrefix + uuidDigest[:62]
}

func IsPersonAddress(address string) bool {
	return getAddressPrefixFromAddress(address) == personAddressPrefix
}
