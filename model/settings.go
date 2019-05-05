package model

import "strings"

var TableCreateSettings = `
	CREATE TABLE settings (
    id integer primary key not null,
	createdon integer not null,
	modifiedon integer not null,
	pricemajoreditsettings integer not null,
	pricemajorcreateperson integer not null,
	pricemajorchangepersonauthorization integer not null,
	pricemajorchangejournalauthorization integer not null,
	pricepersonedit integer not null,
	priceauthorsubmitnewmanuscript integer not null,
	priceauthorsubmitnewversion integer not null,
	priceauthoracceptauthorship integer not null,
	pricereviewersubmit integer not null,
	priceeditorallowmanuscriptreview integer not null,
	priceeditorrejectmanuscript integer not null,
	priceeditorpublishmanuscript integer not null,
	priceeditorassignmanuscript integer not null,
	priceeditorcreatejournal integer not null,
	priceeditorcreatevolume integer not null,
	priceeditoreditjournal integer not null,
	priceeditoraddcolleague integer not null,
	priceeditoracceptduty integer not null)
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
