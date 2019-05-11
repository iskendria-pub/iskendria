package dao

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var theExpectations = []expectation{
	{
		g:        func(s *Settings) int32 { return s.PriceMajorEditSettings },
		expected: 100,
	},
	{
		g:        func(s *Settings) int32 { return s.PriceMajorCreatePerson },
		expected: 200,
	},
	{
		g:        func(s *Settings) int32 { return s.PriceMajorChangePersonAuthorization },
		expected: 300,
	},
	{
		g:        func(s *Settings) int32 { return s.PriceMajorChangeJournalAuthorization },
		expected: 400,
	},
	{
		g:        func(s *Settings) int32 { return s.PricePersonEdit },
		expected: 500,
	},
	{
		g:        func(s *Settings) int32 { return s.PriceAuthorSubmitNewManuscript },
		expected: 600,
	},
	{
		g:        func(s *Settings) int32 { return s.PriceAuthorSubmitNewVersion },
		expected: 700,
	},
	{
		g:        func(s *Settings) int32 { return s.PriceAuthorAcceptAuthorship },
		expected: 800,
	},
	{
		g:        func(s *Settings) int32 { return s.PriceReviewerSubmit },
		expected: 900,
	},
	{
		g:        func(s *Settings) int32 { return s.PriceEditorAllowManuscriptReview },
		expected: 1000,
	},
	{
		g:        func(s *Settings) int32 { return s.PriceEditorRejectManuscript },
		expected: 1100,
	},
	{
		g:        func(s *Settings) int32 { return s.PriceEditorPublishManuscript },
		expected: 1200,
	},
	{
		g:        func(s *Settings) int32 { return s.PriceEditorAssignManuscript },
		expected: 1300,
	},
	{
		g:        func(s *Settings) int32 { return s.PriceEditorCreateJournal },
		expected: 1400,
	},
	{
		g:        func(s *Settings) int32 { return s.PriceEditorCreateVolume },
		expected: 1500,
	},
	{
		g:        func(s *Settings) int32 { return s.PriceEditorEditJournal },
		expected: 1600,
	},
	{
		g:        func(s *Settings) int32 { return s.PriceEditorAddColleague },
		expected: 1700,
	},
	{
		g:        func(s *Settings) int32 { return s.PriceEditorAcceptDuty },
		expected: 1800,
	},
}

type expectation struct {
	g        settingsGetter
	expected int32
}

type settingsGetter func(*Settings) int32

var settingsCreate = &dataManipulationSettingsCreate{
	timestamp:                            10000,
	priceMajorEditSettings:               100,
	priceMajorCreatePerson:               200,
	priceMajorChangePersonAuthorization:  300,
	priceMajorChangeJournalAuthorization: 400,
	pricePersonEdit:                      500,
	priceAuthorSubmitNewManuscript:       600,
	priceAuthorSubmitNewVersion:          700,
	priceAuthorAcceptAuthorship:          800,
	priceReviewerSubmit:                  900,
	priceEditorAllowManuscriptReview:     1000,
	priceEditorRejectManuscript:          1100,
	priceEditorPublishManuscript:         1200,
	priceEditorAssignManuscript:          1300,
	priceEditorCreateJournal:             1400,
	priceEditorCreateVolume:              1500,
	priceEditorEditJournal:               1600,
	priceEditorAddColleague:              1700,
	priceEditorAcceptDuty:                1800,
}

func TestGetSettings(t *testing.T) {
	logger := log.New(os.Stdout, "testGetSettings", log.Flags())
	Init("testGetSettings.db", logger)
	defer ShutdownAndDelete(logger)
	settings, err := GetSettings()
	if err != nil {
		t.Error("With empty settings table, GetSettings() gave an error: " + err.Error())
	}
	if settings != nil {
		t.Error("With empty settings table, GetSettings() gave a non-nil value")
	}
	tx, err := db.Beginx()
	if err != nil {
		t.Error("Could not start transaction: " + err.Error())
		return
	}
	err = settingsCreate.apply(tx)
	if err != nil {
		t.Error("Applying settings create gave an error: " + err.Error())
		return
	}
	err = tx.Commit()
	if err != nil {
		t.Error("Could not commit transaction: " + err.Error())
		return
	}
	settings, err = GetSettings()
	if err != nil {
		t.Error("With filled settings table, GetSettings() gave an error: " + err.Error())
		return
	}
	if settings.Id != int32(THE_SETTINGS_ID) {
		t.Error("With filled settings table, id error")
	}
	if settings.CreatedOn != int64(10000) {
		t.Error("With filled settings table, createdOn error")
	}
	if settings.ModifiedOn != int64(10000) {
		t.Error("With filled settings table, modifiedOn error")
	}
	for i, e := range theExpectations {
		if e.g(settings) != e.expected {
			t.Error(fmt.Sprintf("From read settings, expectation %d failed", i))
		}
	}
}
