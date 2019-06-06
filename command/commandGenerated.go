package command

import (
	"errors"
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"gitlab.bbinfra.net/3estack/alexandria/model"
)

func createModelCommandSettingsUpdate(orig, updated *dao.Settings) *model.CommandSettingsUpdate {
	result := &model.CommandSettingsUpdate{}

	if updated.PriceMajorEditSettings != orig.PriceMajorEditSettings {
		oldValue := orig.PriceMajorEditSettings
		newValue := updated.PriceMajorEditSettings
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceMajorEditSettingsUpdate = theUpdate
	}

	if updated.PriceMajorCreatePerson != orig.PriceMajorCreatePerson {
		oldValue := orig.PriceMajorCreatePerson
		newValue := updated.PriceMajorCreatePerson
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceMajorCreatePersonUpdate = theUpdate
	}

	if updated.PriceMajorChangePersonAuthorization != orig.PriceMajorChangePersonAuthorization {
		oldValue := orig.PriceMajorChangePersonAuthorization
		newValue := updated.PriceMajorChangePersonAuthorization
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceMajorChangePersonAuthorizationUpdate = theUpdate
	}

	if updated.PriceMajorChangeJournalAuthorization != orig.PriceMajorChangeJournalAuthorization {
		oldValue := orig.PriceMajorChangeJournalAuthorization
		newValue := updated.PriceMajorChangeJournalAuthorization
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceMajorChangeJournalAuthorizationUpdate = theUpdate
	}

	if updated.PricePersonEdit != orig.PricePersonEdit {
		oldValue := orig.PricePersonEdit
		newValue := updated.PricePersonEdit
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PricePersonEditUpdate = theUpdate
	}

	if updated.PriceAuthorSubmitNewManuscript != orig.PriceAuthorSubmitNewManuscript {
		oldValue := orig.PriceAuthorSubmitNewManuscript
		newValue := updated.PriceAuthorSubmitNewManuscript
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceAuthorSubmitNewManuscriptUpdate = theUpdate
	}

	if updated.PriceAuthorSubmitNewVersion != orig.PriceAuthorSubmitNewVersion {
		oldValue := orig.PriceAuthorSubmitNewVersion
		newValue := updated.PriceAuthorSubmitNewVersion
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceAuthorSubmitNewVersionUpdate = theUpdate
	}

	if updated.PriceAuthorAcceptAuthorship != orig.PriceAuthorAcceptAuthorship {
		oldValue := orig.PriceAuthorAcceptAuthorship
		newValue := updated.PriceAuthorAcceptAuthorship
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceAuthorAcceptAuthorshipUpdate = theUpdate
	}

	if updated.PriceReviewerSubmit != orig.PriceReviewerSubmit {
		oldValue := orig.PriceReviewerSubmit
		newValue := updated.PriceReviewerSubmit
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceReviewerSubmitUpdate = theUpdate
	}

	if updated.PriceEditorAllowManuscriptReview != orig.PriceEditorAllowManuscriptReview {
		oldValue := orig.PriceEditorAllowManuscriptReview
		newValue := updated.PriceEditorAllowManuscriptReview
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceEditorAllowManuscriptReviewUpdate = theUpdate
	}

	if updated.PriceEditorRejectManuscript != orig.PriceEditorRejectManuscript {
		oldValue := orig.PriceEditorRejectManuscript
		newValue := updated.PriceEditorRejectManuscript
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceEditorRejectManuscriptUpdate = theUpdate
	}

	if updated.PriceEditorPublishManuscript != orig.PriceEditorPublishManuscript {
		oldValue := orig.PriceEditorPublishManuscript
		newValue := updated.PriceEditorPublishManuscript
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceEditorPublishManuscriptUpdate = theUpdate
	}

	if updated.PriceEditorAssignManuscript != orig.PriceEditorAssignManuscript {
		oldValue := orig.PriceEditorAssignManuscript
		newValue := updated.PriceEditorAssignManuscript
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceEditorAssignManuscriptUpdate = theUpdate
	}

	if updated.PriceEditorCreateJournal != orig.PriceEditorCreateJournal {
		oldValue := orig.PriceEditorCreateJournal
		newValue := updated.PriceEditorCreateJournal
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceEditorCreateJournalUpdate = theUpdate
	}

	if updated.PriceEditorCreateVolume != orig.PriceEditorCreateVolume {
		oldValue := orig.PriceEditorCreateVolume
		newValue := updated.PriceEditorCreateVolume
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceEditorCreateVolumeUpdate = theUpdate
	}

	if updated.PriceEditorEditJournal != orig.PriceEditorEditJournal {
		oldValue := orig.PriceEditorEditJournal
		newValue := updated.PriceEditorEditJournal
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceEditorEditJournalUpdate = theUpdate
	}

	if updated.PriceEditorAddColleague != orig.PriceEditorAddColleague {
		oldValue := orig.PriceEditorAddColleague
		newValue := updated.PriceEditorAddColleague
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceEditorAddColleagueUpdate = theUpdate
	}

	if updated.PriceEditorAcceptDuty != orig.PriceEditorAcceptDuty {
		oldValue := orig.PriceEditorAcceptDuty
		newValue := updated.PriceEditorAcceptDuty
		theUpdate := &model.IntUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PriceEditorAcceptDutyUpdate = theUpdate
	}

	return result
}

func checkModelCommandSettingsUpdate(
	c *model.CommandSettingsUpdate, oldSettings *model.StateSettings) error {

	if c.PriceMajorEditSettingsUpdate != nil && c.PriceMajorEditSettingsUpdate.OldValue != oldSettings.PriceList.PriceMajorEditSettings {
		return errors.New(fmt.Sprintf("PriceMajorEditSettings mismatch. Expected %d, got %d",
			c.PriceMajorEditSettingsUpdate.OldValue, oldSettings.PriceList.PriceMajorEditSettings))
	}

	if c.PriceMajorCreatePersonUpdate != nil && c.PriceMajorCreatePersonUpdate.OldValue != oldSettings.PriceList.PriceMajorCreatePerson {
		return errors.New(fmt.Sprintf("PriceMajorCreatePerson mismatch. Expected %d, got %d",
			c.PriceMajorCreatePersonUpdate.OldValue, oldSettings.PriceList.PriceMajorCreatePerson))
	}

	if c.PriceMajorChangePersonAuthorizationUpdate != nil && c.PriceMajorChangePersonAuthorizationUpdate.OldValue != oldSettings.PriceList.PriceMajorChangePersonAuthorization {
		return errors.New(fmt.Sprintf("PriceMajorChangePersonAuthorization mismatch. Expected %d, got %d",
			c.PriceMajorChangePersonAuthorizationUpdate.OldValue, oldSettings.PriceList.PriceMajorChangePersonAuthorization))
	}

	if c.PriceMajorChangeJournalAuthorizationUpdate != nil && c.PriceMajorChangeJournalAuthorizationUpdate.OldValue != oldSettings.PriceList.PriceMajorChangeJournalAuthorization {
		return errors.New(fmt.Sprintf("PriceMajorChangeJournalAuthorization mismatch. Expected %d, got %d",
			c.PriceMajorChangeJournalAuthorizationUpdate.OldValue, oldSettings.PriceList.PriceMajorChangeJournalAuthorization))
	}

	if c.PricePersonEditUpdate != nil && c.PricePersonEditUpdate.OldValue != oldSettings.PriceList.PricePersonEdit {
		return errors.New(fmt.Sprintf("PricePersonEdit mismatch. Expected %d, got %d",
			c.PricePersonEditUpdate.OldValue, oldSettings.PriceList.PricePersonEdit))
	}

	if c.PriceAuthorSubmitNewManuscriptUpdate != nil && c.PriceAuthorSubmitNewManuscriptUpdate.OldValue != oldSettings.PriceList.PriceAuthorSubmitNewManuscript {
		return errors.New(fmt.Sprintf("PriceAuthorSubmitNewManuscript mismatch. Expected %d, got %d",
			c.PriceAuthorSubmitNewManuscriptUpdate.OldValue, oldSettings.PriceList.PriceAuthorSubmitNewManuscript))
	}

	if c.PriceAuthorSubmitNewVersionUpdate != nil && c.PriceAuthorSubmitNewVersionUpdate.OldValue != oldSettings.PriceList.PriceAuthorSubmitNewVersion {
		return errors.New(fmt.Sprintf("PriceAuthorSubmitNewVersion mismatch. Expected %d, got %d",
			c.PriceAuthorSubmitNewVersionUpdate.OldValue, oldSettings.PriceList.PriceAuthorSubmitNewVersion))
	}

	if c.PriceAuthorAcceptAuthorshipUpdate != nil && c.PriceAuthorAcceptAuthorshipUpdate.OldValue != oldSettings.PriceList.PriceAuthorAcceptAuthorship {
		return errors.New(fmt.Sprintf("PriceAuthorAcceptAuthorship mismatch. Expected %d, got %d",
			c.PriceAuthorAcceptAuthorshipUpdate.OldValue, oldSettings.PriceList.PriceAuthorAcceptAuthorship))
	}

	if c.PriceReviewerSubmitUpdate != nil && c.PriceReviewerSubmitUpdate.OldValue != oldSettings.PriceList.PriceReviewerSubmit {
		return errors.New(fmt.Sprintf("PriceReviewerSubmit mismatch. Expected %d, got %d",
			c.PriceReviewerSubmitUpdate.OldValue, oldSettings.PriceList.PriceReviewerSubmit))
	}

	if c.PriceEditorAllowManuscriptReviewUpdate != nil && c.PriceEditorAllowManuscriptReviewUpdate.OldValue != oldSettings.PriceList.PriceEditorAllowManuscriptReview {
		return errors.New(fmt.Sprintf("PriceEditorAllowManuscriptReview mismatch. Expected %d, got %d",
			c.PriceEditorAllowManuscriptReviewUpdate.OldValue, oldSettings.PriceList.PriceEditorAllowManuscriptReview))
	}

	if c.PriceEditorRejectManuscriptUpdate != nil && c.PriceEditorRejectManuscriptUpdate.OldValue != oldSettings.PriceList.PriceEditorRejectManuscript {
		return errors.New(fmt.Sprintf("PriceEditorRejectManuscript mismatch. Expected %d, got %d",
			c.PriceEditorRejectManuscriptUpdate.OldValue, oldSettings.PriceList.PriceEditorRejectManuscript))
	}

	if c.PriceEditorPublishManuscriptUpdate != nil && c.PriceEditorPublishManuscriptUpdate.OldValue != oldSettings.PriceList.PriceEditorPublishManuscript {
		return errors.New(fmt.Sprintf("PriceEditorPublishManuscript mismatch. Expected %d, got %d",
			c.PriceEditorPublishManuscriptUpdate.OldValue, oldSettings.PriceList.PriceEditorPublishManuscript))
	}

	if c.PriceEditorAssignManuscriptUpdate != nil && c.PriceEditorAssignManuscriptUpdate.OldValue != oldSettings.PriceList.PriceEditorAssignManuscript {
		return errors.New(fmt.Sprintf("PriceEditorAssignManuscript mismatch. Expected %d, got %d",
			c.PriceEditorAssignManuscriptUpdate.OldValue, oldSettings.PriceList.PriceEditorAssignManuscript))
	}

	if c.PriceEditorCreateJournalUpdate != nil && c.PriceEditorCreateJournalUpdate.OldValue != oldSettings.PriceList.PriceEditorCreateJournal {
		return errors.New(fmt.Sprintf("PriceEditorCreateJournal mismatch. Expected %d, got %d",
			c.PriceEditorCreateJournalUpdate.OldValue, oldSettings.PriceList.PriceEditorCreateJournal))
	}

	if c.PriceEditorCreateVolumeUpdate != nil && c.PriceEditorCreateVolumeUpdate.OldValue != oldSettings.PriceList.PriceEditorCreateVolume {
		return errors.New(fmt.Sprintf("PriceEditorCreateVolume mismatch. Expected %d, got %d",
			c.PriceEditorCreateVolumeUpdate.OldValue, oldSettings.PriceList.PriceEditorCreateVolume))
	}

	if c.PriceEditorEditJournalUpdate != nil && c.PriceEditorEditJournalUpdate.OldValue != oldSettings.PriceList.PriceEditorEditJournal {
		return errors.New(fmt.Sprintf("PriceEditorEditJournal mismatch. Expected %d, got %d",
			c.PriceEditorEditJournalUpdate.OldValue, oldSettings.PriceList.PriceEditorEditJournal))
	}

	if c.PriceEditorAddColleagueUpdate != nil && c.PriceEditorAddColleagueUpdate.OldValue != oldSettings.PriceList.PriceEditorAddColleague {
		return errors.New(fmt.Sprintf("PriceEditorAddColleague mismatch. Expected %d, got %d",
			c.PriceEditorAddColleagueUpdate.OldValue, oldSettings.PriceList.PriceEditorAddColleague))
	}

	if c.PriceEditorAcceptDutyUpdate != nil && c.PriceEditorAcceptDutyUpdate.OldValue != oldSettings.PriceList.PriceEditorAcceptDuty {
		return errors.New(fmt.Sprintf("PriceEditorAcceptDuty mismatch. Expected %d, got %d",
			c.PriceEditorAcceptDutyUpdate.OldValue, oldSettings.PriceList.PriceEditorAcceptDuty))
	}

	return nil
}

func createSingleUpdatesSettingsUpdate(
	c *model.CommandSettingsUpdate, oldSettings *model.StateSettings, timestamp int64) []singleUpdate {
	result := []singleUpdate{}

	if c.PriceMajorEditSettingsUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceMajorEditSettingsUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceMajorEditSettings,
			eventKey:   model.EV_KEY_PRICE_MAJOR_EDIT_SETTINGS,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PriceMajorCreatePersonUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceMajorCreatePersonUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceMajorCreatePerson,
			eventKey:   model.EV_KEY_PRICE_MAJOR_CREATE_PERSON,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PriceMajorChangePersonAuthorizationUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceMajorChangePersonAuthorizationUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceMajorChangePersonAuthorization,
			eventKey:   model.EV_KEY_PRICE_MAJOR_CHANGE_PERSON_AUTHORIZATION,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PriceMajorChangeJournalAuthorizationUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceMajorChangeJournalAuthorizationUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceMajorChangeJournalAuthorization,
			eventKey:   model.EV_KEY_PRICE_MAJOR_CHANGE_JOURNAL_AUTHORIZATION,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PricePersonEditUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PricePersonEditUpdate.NewValue,
			stateField: &oldSettings.PriceList.PricePersonEdit,
			eventKey:   model.EV_KEY_PRICE_PERSON_EDIT,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PriceAuthorSubmitNewManuscriptUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceAuthorSubmitNewManuscriptUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceAuthorSubmitNewManuscript,
			eventKey:   model.EV_KEY_PRICE_AUTHOR_SUBMIT_NEW_MANUSCRIPT,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PriceAuthorSubmitNewVersionUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceAuthorSubmitNewVersionUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceAuthorSubmitNewVersion,
			eventKey:   model.EV_KEY_PRICE_AUTHOR_SUBMIT_NEW_VERSION,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PriceAuthorAcceptAuthorshipUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceAuthorAcceptAuthorshipUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceAuthorAcceptAuthorship,
			eventKey:   model.EV_KEY_PRICE_AUTHOR_ACCEPT_AUTHORSHIP,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PriceReviewerSubmitUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceReviewerSubmitUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceReviewerSubmit,
			eventKey:   model.EV_KEY_PRICE_REVIEWER_SUBMIT,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PriceEditorAllowManuscriptReviewUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceEditorAllowManuscriptReviewUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceEditorAllowManuscriptReview,
			eventKey:   model.EV_KEY_PRICE_EDITOR_ALLOW_MANUSCRIPT_REVIEW,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PriceEditorRejectManuscriptUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceEditorRejectManuscriptUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceEditorRejectManuscript,
			eventKey:   model.EV_KEY_PRICE_EDITOR_REJECT_MANUSCRIPT,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PriceEditorPublishManuscriptUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceEditorPublishManuscriptUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceEditorPublishManuscript,
			eventKey:   model.EV_KEY_PRICE_EDITOR_PUBLISH_MANUSCRIPT,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PriceEditorAssignManuscriptUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceEditorAssignManuscriptUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceEditorAssignManuscript,
			eventKey:   model.EV_KEY_PRICE_EDITOR_ASSIGN_MANUSCRIPT,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PriceEditorCreateJournalUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceEditorCreateJournalUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceEditorCreateJournal,
			eventKey:   model.EV_KEY_PRICE_EDITOR_CREATE_JOURNAL,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PriceEditorCreateVolumeUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceEditorCreateVolumeUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceEditorCreateVolume,
			eventKey:   model.EV_KEY_PRICE_EDITOR_CREATE_VOLUME,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PriceEditorEditJournalUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceEditorEditJournalUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceEditorEditJournal,
			eventKey:   model.EV_KEY_PRICE_EDITOR_EDIT_JOURNAL,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PriceEditorAddColleagueUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceEditorAddColleagueUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceEditorAddColleague,
			eventKey:   model.EV_KEY_PRICE_EDITOR_ADD_COLLEAGUE,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PriceEditorAcceptDutyUpdate != nil {
		var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue:   c.PriceEditorAcceptDutyUpdate.NewValue,
			stateField: &oldSettings.PriceList.PriceEditorAcceptDuty,
			eventKey:   model.EV_KEY_PRICE_EDITOR_ACCEPT_DUTY,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	return result
}

func createModelCommandPersonUpdateProperties(
	personId string,
	orig, updated *dao.PersonUpdate) *model.CommandPersonUpdateProperties {
	result := &model.CommandPersonUpdateProperties{}
	result.PersonId = personId

	if updated.PublicKey != orig.PublicKey {
		oldValue := orig.PublicKey
		newValue := updated.PublicKey
		theUpdate := &model.StringUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PublicKeyUpdate = theUpdate
	}

	if updated.Name != orig.Name {
		oldValue := orig.Name
		newValue := updated.Name
		theUpdate := &model.StringUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.NameUpdate = theUpdate
	}

	if updated.Email != orig.Email {
		oldValue := orig.Email
		newValue := updated.Email
		theUpdate := &model.StringUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.EmailUpdate = theUpdate
	}

	if updated.Organization != orig.Organization {
		oldValue := orig.Organization
		newValue := updated.Organization
		theUpdate := &model.StringUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.OrganizationUpdate = theUpdate
	}

	if updated.Telephone != orig.Telephone {
		oldValue := orig.Telephone
		newValue := updated.Telephone
		theUpdate := &model.StringUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.TelephoneUpdate = theUpdate
	}

	if updated.Address != orig.Address {
		oldValue := orig.Address
		newValue := updated.Address
		theUpdate := &model.StringUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.AddressUpdate = theUpdate
	}

	if updated.PostalCode != orig.PostalCode {
		oldValue := orig.PostalCode
		newValue := updated.PostalCode
		theUpdate := &model.StringUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.PostalCodeUpdate = theUpdate
	}

	if updated.Country != orig.Country {
		oldValue := orig.Country
		newValue := updated.Country
		theUpdate := &model.StringUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.CountryUpdate = theUpdate
	}

	if updated.ExtraInfo != orig.ExtraInfo {
		oldValue := orig.ExtraInfo
		newValue := updated.ExtraInfo
		theUpdate := &model.StringUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.ExtraInfoUpdate = theUpdate
	}

	return result
}

func checkModelCommandPersonUpdateProperties(
	c *model.CommandPersonUpdateProperties, oldPerson *model.StatePerson) error {

	if c.PublicKeyUpdate != nil && c.PublicKeyUpdate.OldValue != oldPerson.PublicKey {
		return errors.New(fmt.Sprintf("Person update properties value mismatch. Expected %s, got %s",
			c.PublicKeyUpdate.OldValue, oldPerson.PublicKey))
	}

	if c.NameUpdate != nil && c.NameUpdate.OldValue != oldPerson.Name {
		return errors.New(fmt.Sprintf("Person update properties value mismatch. Expected %s, got %s",
			c.NameUpdate.OldValue, oldPerson.Name))
	}

	if c.EmailUpdate != nil && c.EmailUpdate.OldValue != oldPerson.Email {
		return errors.New(fmt.Sprintf("Person update properties value mismatch. Expected %s, got %s",
			c.EmailUpdate.OldValue, oldPerson.Email))
	}

	if c.BiographyHashUpdate != nil && c.BiographyHashUpdate.OldValue != oldPerson.BiographyHash {
		return errors.New(fmt.Sprintf("Person update properties value mismatch. Expected %s, got %s",
			c.BiographyHashUpdate.OldValue, oldPerson.BiographyHash))
	}

	if c.OrganizationUpdate != nil && c.OrganizationUpdate.OldValue != oldPerson.Organization {
		return errors.New(fmt.Sprintf("Person update properties value mismatch. Expected %s, got %s",
			c.OrganizationUpdate.OldValue, oldPerson.Organization))
	}

	if c.TelephoneUpdate != nil && c.TelephoneUpdate.OldValue != oldPerson.Telephone {
		return errors.New(fmt.Sprintf("Person update properties value mismatch. Expected %s, got %s",
			c.TelephoneUpdate.OldValue, oldPerson.Telephone))
	}

	if c.AddressUpdate != nil && c.AddressUpdate.OldValue != oldPerson.Address {
		return errors.New(fmt.Sprintf("Person update properties value mismatch. Expected %s, got %s",
			c.AddressUpdate.OldValue, oldPerson.Address))
	}

	if c.PostalCodeUpdate != nil && c.PostalCodeUpdate.OldValue != oldPerson.PostalCode {
		return errors.New(fmt.Sprintf("Person update properties value mismatch. Expected %s, got %s",
			c.PostalCodeUpdate.OldValue, oldPerson.PostalCode))
	}

	if c.CountryUpdate != nil && c.CountryUpdate.OldValue != oldPerson.Country {
		return errors.New(fmt.Sprintf("Person update properties value mismatch. Expected %s, got %s",
			c.CountryUpdate.OldValue, oldPerson.Country))
	}

	if c.ExtraInfoUpdate != nil && c.ExtraInfoUpdate.OldValue != oldPerson.ExtraInfo {
		return errors.New(fmt.Sprintf("Person update properties value mismatch. Expected %s, got %s",
			c.ExtraInfoUpdate.OldValue, oldPerson.ExtraInfo))
	}

	return nil
}

func createSingleUpdatesPersonUpdateProperties(
	c *model.CommandPersonUpdateProperties, oldPerson *model.StatePerson, timestamp int64) []singleUpdate {
	result := []singleUpdate{}

	if c.PublicKeyUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.PublicKeyUpdate.NewValue,
			stateField: &oldPerson.PublicKey,
			eventKey:   model.EV_KEY_PERSON_PUBLIC_KEY,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.NameUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.NameUpdate.NewValue,
			stateField: &oldPerson.Name,
			eventKey:   model.EV_KEY_PERSON_NAME,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.EmailUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.EmailUpdate.NewValue,
			stateField: &oldPerson.Email,
			eventKey:   model.EV_KEY_PERSON_EMAIL,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.BiographyHashUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.BiographyHashUpdate.NewValue,
			stateField: &oldPerson.BiographyHash,
			eventKey:   model.EV_KEY_PERSON_BIOGRAPHY_HASH,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.OrganizationUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.OrganizationUpdate.NewValue,
			stateField: &oldPerson.Organization,
			eventKey:   model.EV_KEY_PERSON_ORGANIZATION,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.TelephoneUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.TelephoneUpdate.NewValue,
			stateField: &oldPerson.Telephone,
			eventKey:   model.EV_KEY_PERSON_TELEPHONE,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.AddressUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.AddressUpdate.NewValue,
			stateField: &oldPerson.Address,
			eventKey:   model.EV_KEY_PERSON_ADDRESS,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PostalCodeUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.PostalCodeUpdate.NewValue,
			stateField: &oldPerson.PostalCode,
			eventKey:   model.EV_KEY_PERSON_POSTAL_CODE,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.CountryUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.CountryUpdate.NewValue,
			stateField: &oldPerson.Country,
			eventKey:   model.EV_KEY_PERSON_COUNTRY,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.ExtraInfoUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.ExtraInfoUpdate.NewValue,
			stateField: &oldPerson.ExtraInfo,
			eventKey:   model.EV_KEY_PERSON_EXTRA_INFO,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	return result
}
