package model

import "strings"

var TableCreateSettings = `
	createdOn integer not null,
	modifiedOn integer not null,
	priceMajorEditSettings integer not null,
	priceMajorCreatePerson integer not null,
	priceMajorChangePersonAuthorization integer not null,
	priceMajorChangeJournalAuthorization integer not null,
	pricePersonEdit integer not null,
	priceAuthorSubmitNewManuscript integer not null,
	priceAuthorSubmitNewVersion integer not null,
	priceAuthorAcceptAuthorship integer not null,
	priceReviewerSubmit integer not null,
	priceEditorAllowManuscriptReview integer not null,
	priceEditorRejectManuscript integer not null,
	priceEditorPublishManuscript integer not null,
	priceEditorAssignManuscript integer not null,
	priceEditorCreateJournal integer not null,
	priceEditorCreateVolume integer not null,
	priceEditorEditJournal integer not null,
	priceEditorAddColleague integer not null,
	priceEditorAcceptDuty integer not null
`

const (
	PRICE_MAJOR_EDIT_SETTINGS                = "priceMajorEditSettings"
	PRICE_MAJOR_CREATE_PERSON                = "priceMajorCreatePerson"
	PRICE_MAJOR_CHANGE_PERSON_AUTHORIZATION  = "priceMajorChangePersonAuthorization"
	PRICE_MAJOR_CHANGE_JOURNAL_AUTHORIZATION = "priceMajorChangeJournalAuthorization"
	PRICE_PERSON_EDIT                        = "pricePersonEdit"
	PRICE_AUTHOR_SUBMIT_NEW_MANUSCRIPT       = "priceAuthorSubmitNewManuscript"
	PRICE_AUTHOR_SUBMIT_NEW_VERSION          = "priceAuthorSubmitNewVersion"
	PRICE_AUTHOR_ACCEPT_AUTHORSHIP           = "priceAuthorAcceptAuthorship"
	PRICE_REVIEWER_SUBMIT                    = "priceReviewerSubmit"
	PRICE_EDITOR_ALLOW_MANUSCRIPT_REVIEW     = "priceEditorAllowManuscriptReview"
	PRICE_EDITOR_REJECT_MANUSCRIPT           = "priceEditorRejectManuscript"
	PRICE_EDITOR_PUBLISH_MANUSCRIPT          = "priceEditorPublishManuscript"
	PRICE_EDITOR_ASSIGN_MANUSCRIPT           = "priceEditorAssignManuscript"
	PRICE_EDITOR_CREATE_JOURNAL              = "priceEditorCreateJournal"
	PRICE_EDITOR_CREATE_VOLUME               = "priceEditorCreateVolume"
	PRICE_EDITOR_EDIT_JOURNAL                = "priceEditorEditJournal"
	PRICE_EDITOR_ADD_COLLEAGUE               = "priceEditorAddColleague"
	PRICE_EDITOR_ACCEPT_DUTY                 = "priceEditorAcceptDuty"
)

func GetSettingsAddress() string {
	return Namespace + strings.Repeat("0", 64)
}
