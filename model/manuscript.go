package model

import "github.com/google/uuid"

var TableCreateManuscript = `
CREATE TABLE manuscript (
    id VARCHAR primary key not null,
    createdon integer not null,
    modifiedon integer not null,
    hash VARCHAR not null,
    threadid VARCHAR not null,
    versionnumber integer not null,
    commitmsg VARCHAR not null,
    title VARCHAR not null,
    status integer not null,
    journalid VARCHAR not null,
    volumeid VARCHAR not null,
    firstpage VARCHAR not null,
    lastpage VARCHAR not null,
    isreviewable bool not null
)
`

var TableCreateAuthor = `
CREATE TABLE author (
    manuscriptid VARCHAR not null,
    personid VARCHAR not null,
    didsign bool not null,
    authornumber integer not null,
    PRIMARY KEY (manuscriptid, personid),
    FOREIGN KEY (manuscriptid) REFERENCES manuscript(id),
    FOREIGN KEY (personid) REFERENCES person(id)
)
`

// TODO: review

const (
	EV_TYPE_MANUSCRIPT_CREATE            = "evManuscriptCreate"
	EV_TYPE_MANUSCRIPT_UPDATE            = "evManuscriptUpdate"
	EV_TYPE_MANUSCRIPT_MODIFICATION_TIME = "evManuscriptModificationTime"
	EV_TYPE_AUTHOR_CREATE                = "evAuthorCreate"
	EV_TYPE_AUTHOR_UPDATE                = "evAuthorUpdate"
)

const (
	EV_KEY_MANUSCRIPT_HASH           = "hash"
	EV_KEY_MANUSCRIPT_THREAD_ID      = "threadId"
	EV_KEY_MANUSCRIPT_VERSION_NUMBER = "versionNumber"
	EV_KEY_MANUSCRIPT_COMMIT_MSG     = "commitMsg"
	EV_KEY_MANUSCRIPT_TITLE          = "title"
	EV_KEY_MANUSCRIPT_STATUS         = "status"
	EV_KEY_VOLUME_ID                 = "volumeId"
	EV_KEY_MANUSCRIPT_FIRST_PAGE     = "firstPage"
	EV_KEY_MANUSCRIPT_LAST_PAGE      = "lastPage"
	EV_KEY_MANUSCRIPT_IS_REVIEWABLE  = "isReviewable"
)

const (
	EV_KEY_MANUSCRIPT_ID   = "manuscriptId"
	EV_KEY_PERSON_ID       = "personId"
	EV_KEY_AUTHOR_DID_SIGN = "didSign"
	EV_KEY_AUTHOR_NUMBER   = "authorNumber"
)

const manuscriptAddressPrefix = "10"

func CreateManuscriptAddress() string {
	var theUuid uuid.UUID = uuid.New()
	uuidDigest := hexdigestOfUuid(theUuid)
	return Namespace + manuscriptAddressPrefix + uuidDigest[:62]
}

func IsManuscriptAddress(address string) bool {
	return getAddressPrefixFromAddress(address) == manuscriptAddressPrefix
}

func GetManuscriptStatusString(status ManuscriptStatus) string {
	switch status {
	case ManuscriptStatus_init:
		return "INIT"
	case ManuscriptStatus_new:
		return "NEW"
	case ManuscriptStatus_reviewable:
		return "REVIEWABLE"
	case ManuscriptStatus_rejected:
		return "REJECTED"
	case ManuscriptStatus_published:
		return "PUBLISHED"
	case ManuscriptStatus_assigned:
		return "ASSIGNED"
	default:
		panic("Invalud manuscript status")
	}
}

func GetManuscriptJudgementString(judgement ManuscriptJudgement) string {
	switch judgement {
	case ManuscriptJudgement_judgementAccepted:
		return "ACCEPTED"
	case ManuscriptJudgement_judgementRejected:
		return "REJECTED"
	default:
		panic("Invalid manuscript judgement")
	}
}