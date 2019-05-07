package dao

import (
	"database/sql"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/jmoiron/sqlx"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"strconv"
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

func createSettingsCreateEvent(event *events_pb2.Event) (event, error) {
	var err error
	var i64 int64
	transactionId := ""
	eventSeq := int32(0)
	dataManipulation := &dataManipulationSettingsCreate{}
	for _, attribute := range event.Attributes {
		switch attribute.Key {
		case model.TRANSACTION_ID:
			transactionId = attribute.Value
		case model.TIMESTAMP:
			i64, err = strconv.ParseInt(attribute.Value, 10, 64)
			dataManipulation.timestamp = i64
		case model.EVENT_SEQ:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			eventSeq = int32(i64)
		case model.PRICE_MAJOR_EDIT_SETTINGS:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceMajorEditSettings = int32(i64)
		case model.PRICE_MAJOR_CREATE_PERSON:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceMajorCreatePerson = int32(i64)
		case model.PRICE_MAJOR_CHANGE_PERSON_AUTHORIZATION:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceMajorChangePersonAuthorization = int32(i64)
		case model.PRICE_MAJOR_CHANGE_JOURNAL_AUTHORIZATION:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceMajorChangeJournalAuthorization = int32(i64)
		case model.PRICE_PERSON_EDIT:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.pricePersonEdit = int32(i64)
		case model.PRICE_AUTHOR_SUBMIT_NEW_MANUSCRIPT:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceAuthorSubmitNewManuscript = int32(i64)
		case model.PRICE_AUTHOR_SUBMIT_NEW_VERSION:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceAuthorSubmitNewVersion = int32(i64)
		case model.PRICE_AUTHOR_ACCEPT_AUTHORSHIP:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceAuthorAcceptAuthorship = int32(i64)
		case model.PRICE_REVIEWER_SUBMIT:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceReviewerSubmit = int32(i64)
		case model.PRICE_EDITOR_ALLOW_MANUSCRIPT_REVIEW:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceEditorAllowManuscriptReview = int32(i64)
		case model.PRICE_EDITOR_REJECT_MANUSCRIPT:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceEditorRejectManuscript = int32(i64)
		case model.PRICE_EDITOR_PUBLISH_MANUSCRIPT:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceEditorPublishManuscript = int32(i64)
		case model.PRICE_EDITOR_ASSIGN_MANUSCRIPT:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceEditorAssignManuscript = int32(i64)
		case model.PRICE_EDITOR_CREATE_JOURNAL:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceEditorCreateJournal = int32(i64)
		case model.PRICE_EDITOR_CREATE_VOLUME:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceEditorCreateVolume = int32(i64)
		case model.PRICE_EDITOR_EDIT_JOURNAL:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceEditorEditJournal = int32(i64)
		case model.PRICE_EDITOR_ADD_COLLEAGUE:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceEditorAddColleague = int32(i64)
		case model.PRICE_EDITOR_ACCEPT_DUTY:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			dataManipulation.priceEditorAcceptDuty = int32(i64)
		}
		if err != nil {
			return nil, err
		}
	}
	return &dataManipulationEvent{
		transactionId:    transactionId,
		eventSeq:         eventSeq,
		dataManipulation: dataManipulation,
	}, nil
}

type dataManipulationSettingsCreate struct {
	timestamp                            int64
	priceMajorEditSettings               int32
	priceMajorCreatePerson               int32
	priceMajorChangePersonAuthorization  int32
	priceMajorChangeJournalAuthorization int32
	pricePersonEdit                      int32
	priceAuthorSubmitNewManuscript       int32
	priceAuthorSubmitNewVersion          int32
	priceAuthorAcceptAuthorship          int32
	priceReviewerSubmit                  int32
	priceEditorAllowManuscriptReview     int32
	priceEditorRejectManuscript          int32
	priceEditorPublishManuscript         int32
	priceEditorAssignManuscript          int32
	priceEditorCreateJournal             int32
	priceEditorCreateVolume              int32
	priceEditorEditJournal               int32
	priceEditorAddColleague              int32
	priceEditorAcceptDuty                int32
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
