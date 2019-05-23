package model

import (
	"fmt"
	"github.com/google/uuid"
)

var TableCreateJournal = `
CREATE TABLE journal (
	journalid VARCHAR primary key not null,
	createdon integer not null,
    modifiedon integer not null,
    title string not null,
    issigned bool not null,
    descriptionhash string not null
)
`

var TableCreateEditor = `
CREATE TABLE editor (
    journalid VARCHAR not null,
    personid VARCHAR not null,
    editorstate integer not null,
    PRIMARY KEY (journalid, personid),
    FOREIGN KEY (journalid) REFERENCES journal(journalid)
)
`

const (
	EV_TYPE_JOURNAL_CREATE            = "evJournalCreate"
	EV_TYPE_JOURNAL_UPDATE            = "evJournalUpdate"
	EV_TYPE_JOURNAL_MODIFICATION_TIME = "evJournalModification"
	EV_TYPE_EDITOR_CREATE             = "evEditorCreate"
	EV_TYPE_EDITOR_UPDATE             = "evEditorUpdate"
	EV_TYPE_EDITOR_DELETE             = "evEditorDelete"
)

const (
	EV_KEY_JOURNAL_ID       = "journalId"
	EV_KEY_TITLE            = "title"
	EV_KEY_IS_SIGNED        = "isSigned"
	EV_KEY_DESCRIPTION_HASH = "descriptionHash"
	EV_KEY_PERSON_ID        = "personId"
	EV_KEY_EDITOR_STATE     = "editorState"
)

const journalAddressPrefix = "20"

func CreateJournalAddress() string {
	var theUuid uuid.UUID = uuid.New()
	uuidDigest := hexdigestOfUuid(theUuid)
	return Namespace + journalAddressPrefix + uuidDigest[:62]
}

func IsJournalAddress(address string) bool {
	return getAddressPrefixFromAddress(address) == journalAddressPrefix
}

func GetEditorStateString(value EditorState) string {
	switch value {
	case EditorState_editorProposed:
		return "PROPOSED"
	case EditorState_editorAccepted:
		return "ACCEPTED"
	}
	panic(fmt.Sprintf("Unknown editor state: %d", value))
}