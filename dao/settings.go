package dao

import "database/sql"

type Settings struct {
	Id                                   int32
	CreatedOn                            int64 `db:"createdon"`
	ModifiedOn                           int64 `db:"modifiedon"`
	PriceMajorEditSettings               int32 `db:"pricemajoreditsettings"`
	PriceMajorCreatePerson               int32 `db:"pricemajorcreateperson"`
	PriceMajorChangePersonAuthorization  int32 `db:"pricemajorchangepersonauthorization"`
	PriceMajorChangeJournalAuthorization int32 `db:"pricemajorchangejournalauthorization"`
	PricePersonEdit                      int32 `db:"pricepersonedit"`
	PriceAuthorSubmitNewManuscript       int32 `db:"priceauthorsubmitnewmanuscript"`
	PriceAuthorSubmitNewVersion          int32 `db:"priceauthorsubmitnewversion"`
	PriceAuthorAcceptAuthorship          int32 `db:"priceauthoracceptauthorship"`
	PriceReviewerSubmit                  int32 `db:"pricereviewersubmit"`
	PriceEditorAllowManuscriptReview     int32 `db:"priceeditorallowmanuscriptreview"`
	PriceEditorRejectManuscript          int32 `db:"priceeditorrejectmanuscript"`
	PriceEditorPublishManuscript         int32 `db:"priceeditorpublishmanuscript"`
	PriceEditorAssignManuscript          int32 `db:"priceeditorassignmanuscript"`
	PriceEditorCreateJournal             int32 `db:"priceeditorcreatejournal"`
	PriceEditorCreateVolume              int32 `db:"priceeditorcreatevolume"`
	PriceEditorEditJournal               int32 `db:"priceeditoreditjournal"`
	PriceEditorAddColleague              int32 `db:"priceeditoraddcolleague"`
	PriceEditorAcceptDuty                int32 `db:"priceeditoracceptduty"`
}

func GetSettings() (*Settings, error) {
	var settings = new(Settings)
	err := db.QueryRowx("SELECT * FROM settings WHERE id = ?", THE_SETTINGS_ID).StructScan(settings)
	if err == nil {
		return settings, nil
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return nil, err
}
