package model

// The field names are derived from the event keys.
// When an event key is taken to lower case, the
// corresponding field name is obtained.
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
	EV_TYPE_PERSON_CREATE            = "evPersonCreate"
	EV_TYPE_PERSON_UPDATE            = "evPersonUpdate"
	EV_TYPE_PERSON_MODIFICATION_TIME = "evPersonModificationTime"
)

const (
	EV_KEY_PERSON_PUBLIC_KEY     = "publicKey"
	EV_KEY_PERSON_NAME           = "name"
	EV_KEY_PERSON_EMAIL          = "email"
	EV_KEY_PERSON_IS_MAJOR       = "isMajor"
	EV_KEY_PERSON_IS_SIGNED      = "isSigned"
	EV_KEY_PERSON_BALANCE        = "balance"
	EV_KEY_PERSON_BIOGRAPHY_HASH = "biographyHash"
	EV_KEY_PERSON_ORGANIZATION   = "organization"
	EV_KEY_PERSON_TELEPHONE      = "telephone"
	EV_KEY_PERSON_ADDRESS        = "address"
	EV_KEY_PERSON_POSTAL_CODE    = "postalCode"
	EV_KEY_PERSON_COUNTRY        = "country"
	EV_KEY_PERSON_EXTRA_INFO     = "extraInfo"
)
