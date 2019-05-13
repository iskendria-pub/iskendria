package model

import (
	"strings"
)

// The field names are derived from the event keys.
// When an event key is taken to lower case, the
// corresponding field name is obtained.
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
	EV_TYPE_SETTINGS_CREATE            = "evSettingsCreate"
	EV_TYPE_SETTINGS_UPDATE            = "evSettingsUpdate"
	EV_TYPE_SETTINGS_MODIFICATION_TIME = "evSettingsModificationTime"
)

const (
	EV_KEY_PRICE_MAJOR_EDIT_SETTINGS                = "priceMajorEditSettings"
	EV_KEY_PRICE_MAJOR_CREATE_PERSON                = "priceMajorCreatePerson"
	EV_KEY_PRICE_MAJOR_CHANGE_PERSON_AUTHORIZATION  = "priceMajorChangePersonAuthorization"
	EV_KEY_PRICE_MAJOR_CHANGE_JOURNAL_AUTHORIZATION = "priceMajorChangeJournalAuthorization"
	EV_KEY_PRICE_PERSON_EDIT                        = "pricePersonEdit"
	EV_KEY_PRICE_AUTHOR_SUBMIT_NEW_MANUSCRIPT       = "priceAuthorSubmitNewManuscript"
	EV_KEY_PRICE_AUTHOR_SUBMIT_NEW_VERSION          = "priceAuthorSubmitNewVersion"
	EV_KEY_PRICE_AUTHOR_ACCEPT_AUTHORSHIP           = "priceAuthorAcceptAuthorship"
	EV_KEY_PRICE_REVIEWER_SUBMIT                    = "priceReviewerSubmit"
	EV_KEY_PRICE_EDITOR_ALLOW_MANUSCRIPT_REVIEW     = "priceEditorAllowManuscriptReview"
	EV_KEY_PRICE_EDITOR_REJECT_MANUSCRIPT           = "priceEditorRejectManuscript"
	EV_KEY_PRICE_EDITOR_PUBLISH_MANUSCRIPT          = "priceEditorPublishManuscript"
	EV_KEY_PRICE_EDITOR_ASSIGN_MANUSCRIPT           = "priceEditorAssignManuscript"
	EV_KEY_PRICE_EDITOR_CREATE_JOURNAL              = "priceEditorCreateJournal"
	EV_KEY_PRICE_EDITOR_CREATE_VOLUME               = "priceEditorCreateVolume"
	EV_KEY_PRICE_EDITOR_EDIT_JOURNAL                = "priceEditorEditJournal"
	EV_KEY_PRICE_EDITOR_ADD_COLLEAGUE               = "priceEditorAddColleague"
	EV_KEY_PRICE_EDITOR_ACCEPT_DUTY                 = "priceEditorAcceptDuty"
)

func GetSettingsAddress() string {
	return Namespace + strings.Repeat("0", 64)
}
