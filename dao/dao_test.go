package dao

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"gitlab.bbinfra.net/3estack/alexandria/util"
	"log"
	"os"
	"strconv"
	"testing"
)

const (
	MAX_TIME_DIFF_SEC = 10
)

const (
	firstBlock       = "first block"
	secondBlock      = "second block"
	thirdBlock       = "third block"
	theTransactionId = "transactionId"
)

func TestBootstrapHappy(t *testing.T) {
	personId := model.CreatePersonAddress()
	scenarios := [][]*events_pb2.Event{
		{
			getBlockCommitEvent(firstBlock, ""),
			getTransactionControlEvent(theTransactionId, 0, 3),
			getCreateSettingsEvent(theTransactionId, 1),
			getCreatePersonEvent(theTransactionId, 2, personId),
		},
		{
			getBlockCommitEvent(firstBlock, ""),
			getTransactionControlEvent(theTransactionId, 0, 3),
			getCreateSettingsEvent(theTransactionId, 1),
			getCreatePersonEvent(theTransactionId, 2, personId),
			getBlockCommitEvent(secondBlock, firstBlock),
			getBlockCommitEvent(thirdBlock, secondBlock),
		},
		{
			getBlockCommitEvent(firstBlock, ""),
			getCreateSettingsEvent(theTransactionId, 0),
			getTransactionControlEvent(theTransactionId, 1, 3),
			getCreatePersonEvent(theTransactionId, 2, personId),
		},
		{
			getBlockCommitEvent(firstBlock, ""),
			getCreateSettingsEvent(theTransactionId, 0),
			getCreatePersonEvent(theTransactionId, 1, personId),
			getTransactionControlEvent(theTransactionId, 2, 3),
		},
	}
	for i, events := range scenarios {
		fmt.Printf("TestBootstrapHappy scenario %d\n", i)
		applyEventsHappy(events, personId, t)
	}
}

func applyEventsHappy(events []*events_pb2.Event, personId string, t *testing.T) {
	logger := log.New(os.Stdout, "testBootstrap", log.Flags())
	dbFile := "testBootstrap.db"
	util.RemoveFileIfExists(dbFile, logger)
	Init(dbFile, logger)
	defer ShutdownAndDelete(logger)
	for _, ev := range events {
		if err := HandleEvent(ev); err != nil {
			t.Error(fmt.Sprintf("Error handling event: error %s, event %v", err, ev))
		}
	}
	actualSettings, err := GetSettings()
	if err != nil {
		t.Error("Could not get settings: " + err.Error())
		return
	}
	actualPerson, err := GetPersonById(personId)
	if err != nil {
		t.Error("Could not get person: " + err.Error())
		return
	}
	if actualSettings == nil {
		t.Error("No settings were present in db")
		return
	}
	if actualPerson == nil {
		t.Error(fmt.Sprintf("Person was not present in db: %s", personId))
		return
	}
	if util.Abs(actualSettings.CreatedOn-model.GetCurrentTime()) >= MAX_TIME_DIFF_SEC {
		t.Error("Settings creation time mismatch")
	}
	if util.Abs(actualSettings.ModifiedOn-model.GetCurrentTime()) >= MAX_TIME_DIFF_SEC {
		t.Error("Settings modification time mismatch")
	}
	if actualSettings.Id != THE_SETTINGS_ID {
		t.Error("Unexpected id value for settings")
	}
	if actualSettings.PriceMajorEditSettings != int32(1) ||
		actualSettings.PriceMajorCreatePerson != int32(2) ||
		actualSettings.PriceMajorChangePersonAuthorization != int32(3) ||
		actualSettings.PriceMajorChangeJournalAuthorization != int32(4) ||
		actualSettings.PricePersonEdit != int32(5) ||
		actualSettings.PriceAuthorSubmitNewManuscript != int32(6) ||
		actualSettings.PriceAuthorSubmitNewVersion != int32(7) ||
		actualSettings.PriceAuthorAcceptAuthorship != int32(8) ||
		actualSettings.PriceReviewerSubmit != int32(9) ||
		actualSettings.PriceEditorAllowManuscriptReview != int32(10) ||
		actualSettings.PriceEditorRejectManuscript != int32(11) ||
		actualSettings.PriceEditorPublishManuscript != int32(12) ||
		actualSettings.PriceEditorAssignManuscript != int32(13) ||
		actualSettings.PriceEditorCreateJournal != int32(14) ||
		actualSettings.PriceEditorCreateVolume != int32(15) ||
		actualSettings.PriceEditorEditJournal != int32(16) ||
		actualSettings.PriceEditorAddColleague != int32(17) ||
		actualSettings.PriceEditorAcceptDuty != int32(18) {
		t.Error("Price mismatch")
	}
	if actualPerson.Id != personId {
		t.Error(fmt.Sprintf("Person id mismatch, expected %s, got %s",
			personId, actualPerson.Id))
	}
	if util.Abs(actualPerson.CreatedOn-model.GetCurrentTime()) >= int64(MAX_TIME_DIFF_SEC) {
		t.Error("Creation time mismatch")
	}
	if util.Abs(actualPerson.ModifiedOn-model.GetCurrentTime()) >= int64(MAX_TIME_DIFF_SEC) {
		t.Error("Modification time mismatch")
	}
	if actualPerson.PublicKey != "Key Martijn" {
		t.Error("Publick key mismatch")
	}
	if actualPerson.Name != "Martijn" {
		t.Error("Name mismatch")
	}
	if actualPerson.Email != "xxx@gmail.com" {
		t.Error("Email mismatch")
	}
	if actualPerson.IsMajor != false {
		t.Error("Should not be major")
	}
	if actualPerson.IsSigned != false {
		t.Error("Should not be signed")
	}
	if actualPerson.Balance != int32(0) {
		t.Error("Should not have balance")
	}
	if actualPerson.BiographyHash != "" {
		t.Error("Should not have bibliography hash")
	}
	if actualPerson.BiographyFormat != "" {
		t.Error("Should not have bibliography format")
	}
	if actualPerson.Organization != "" {
		t.Error("Should not have organization")
	}
	if actualPerson.Telephone != "" {
		t.Error("Should not have telephone")
	}
	if actualPerson.Address != "" {
		t.Error("Should not have address")
	}
	if actualPerson.PostalCode != "" {
		t.Error("Should not have postal code")
	}
	if actualPerson.Country != "" {
		t.Error("Should not have country")
	}
	if actualPerson.ExtraInfo != "" {
		t.Error("Should not have extra info")
	}
}

func getBlockCommitEvent(current, previous string) *events_pb2.Event {
	return &events_pb2.Event{
		EventType: model.EV_SAWTOOTH_BLOCK_COMMIT,
		Attributes: []*events_pb2.Event_Attribute{
			{
				Key:   model.SAWTOOTH_CURRENT_BLOCK_ID,
				Value: current,
			},
			{
				Key:   model.SAWTOOTH_PREVIOUS_BLOCK_ID,
				Value: previous,
			},
		},
	}
}

func getTransactionControlEvent(transactionId string, eventSeq, numEvents int32) *events_pb2.Event {
	return &events_pb2.Event{
		EventType: model.FamilyName + "/" + model.EV_TRANSACTION_CONTROL,
		Attributes: []*events_pb2.Event_Attribute{
			{
				Key:   model.TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.EVENT_SEQ,
				Value: strconv.FormatInt(int64(eventSeq), 10),
			},
			{
				Key:   model.NUM_EVENTS,
				Value: strconv.FormatInt(int64(numEvents), 10),
			},
		},
	}
}

func getCreateSettingsEvent(transactionId string, eventSeq int32) *events_pb2.Event {
	return &events_pb2.Event{
		EventType: model.EV_SETTINGS_CREATE,
		Attributes: []*events_pb2.Event_Attribute{
			{
				Key:   model.TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.TIMESTAMP,
				Value: strconv.FormatInt(model.GetCurrentTime(), 10),
			},
			{
				Key:   model.EVENT_SEQ,
				Value: strconv.FormatInt(int64(eventSeq), 10),
			},
			{
				Key:   model.PRICE_MAJOR_EDIT_SETTINGS,
				Value: "1",
			},
			{
				Key:   model.PRICE_MAJOR_CREATE_PERSON,
				Value: "2",
			},
			{
				Key:   model.PRICE_MAJOR_CHANGE_PERSON_AUTHORIZATION,
				Value: "3",
			},
			{
				Key:   model.PRICE_MAJOR_CHANGE_JOURNAL_AUTHORIZATION,
				Value: "4",
			},
			{
				Key:   model.PRICE_PERSON_EDIT,
				Value: "5",
			},
			{
				Key:   model.PRICE_AUTHOR_SUBMIT_NEW_MANUSCRIPT,
				Value: "6",
			},
			{
				Key:   model.PRICE_AUTHOR_SUBMIT_NEW_VERSION,
				Value: "7",
			},
			{
				Key:   model.PRICE_AUTHOR_ACCEPT_AUTHORSHIP,
				Value: "8",
			},
			{
				Key:   model.PRICE_REVIEWER_SUBMIT,
				Value: "9",
			},
			{
				Key:   model.PRICE_EDITOR_ALLOW_MANUSCRIPT_REVIEW,
				Value: "10",
			},
			{
				Key:   model.PRICE_EDITOR_REJECT_MANUSCRIPT,
				Value: "11",
			},
			{
				Key:   model.PRICE_EDITOR_PUBLISH_MANUSCRIPT,
				Value: "12",
			},
			{
				Key:   model.PRICE_EDITOR_ASSIGN_MANUSCRIPT,
				Value: "13",
			},
			{
				Key:   model.PRICE_EDITOR_CREATE_JOURNAL,
				Value: "14",
			},
			{
				Key:   model.PRICE_EDITOR_CREATE_VOLUME,
				Value: "15",
			},
			{
				Key:   model.PRICE_EDITOR_EDIT_JOURNAL,
				Value: "16",
			},
			{
				Key:   model.PRICE_EDITOR_ADD_COLLEAGUE,
				Value: "17",
			},
			{
				Key:   model.PRICE_EDITOR_ACCEPT_DUTY,
				Value: "18",
			},
		},
	}
}

func getCreatePersonEvent(transactionId string, eventSeq int32, personId string) *events_pb2.Event {
	return &events_pb2.Event{
		EventType: model.EV_PERSON_CREATE,
		Attributes: []*events_pb2.Event_Attribute{
			{
				Key:   model.TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.TIMESTAMP,
				Value: strconv.FormatInt(model.GetCurrentTime(), 10),
			},
			{
				Key:   model.EVENT_SEQ,
				Value: strconv.FormatInt(int64(eventSeq), 10),
			},
			{
				Key:   model.ID,
				Value: personId,
			},
			{
				Key:   model.PERSON_NAME,
				Value: "Martijn",
			},
			{
				Key:   model.PERSON_PUBLIC_KEY,
				Value: "Key Martijn",
			},
			{
				Key:   model.PERSON_EMAIL,
				Value: "xxx@gmail.com",
			},
		},
	}
}
