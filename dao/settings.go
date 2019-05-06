package dao

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

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

type dataManipulationSettingsCreate struct {
	timestamp int64
	priceMajorEditSettings int32
	priceMajorCreatePerson int32
	priceMajorChangePersonAuthorization int32
	priceMajorChangeJournalAuthorization int32
	pricePersonEdit int32
	priceAuthorSubmitNewManuscript int32
	priceAuthorSubmitNewVersion int32
	priceAuthorAcceptAuthorship int32
	priceReviewerSubmit int32
	priceEditorAllowManuscriptReview int32
	priceEditorRejectManuscript int32
	priceEditorPublishManuscript int32
	priceEditorAssignManuscript int32
	priceEditorCreateJournal int32
	priceEditorCreateVolume int32
	priceEditorEditJournal int32
	priceEditorAddColleague int32
	priceEditorAcceptDuty int32
}

var _ dataManipulation = new(dataManipulationSettingsCreate)

func (dmsc *dataManipulationSettingsCreate) apply(tx *sqlx.Tx) error {
	_, err := tx.Exec(fmt.Sprintf("INSERT INTO settings VALUES (%s)", GetPlaceHolders(21)),
		// id, createdOn, modifiedOn
		THE_SETTINGS_ID, dmsc.timestamp, dmsc.timestamp,
		// prices
		dmsc.priceMajorEditSettings,
		dmsc.priceMajorCreatePerson,
		dmsc.priceMajorChangePersonAuthorization,
		dmsc.priceMajorChangeJournalAuthorization,
		dmsc.pricePersonEdit,
		dmsc.priceAuthorSubmitNewManuscript,
		dmsc.priceAuthorSubmitNewVersion,
		dmsc.priceAuthorAcceptAuthorship,
		dmsc.priceReviewerSubmit,
		dmsc.priceEditorAllowManuscriptReview,
		dmsc.priceEditorRejectManuscript,
		dmsc.priceEditorPublishManuscript,
		dmsc.priceEditorAssignManuscript,
		dmsc.priceEditorCreateJournal,
		dmsc.priceEditorCreateVolume,
		dmsc.priceEditorEditJournal,
		dmsc.priceEditorAddColleague,
		dmsc.priceEditorAcceptDuty)
	return err
}
