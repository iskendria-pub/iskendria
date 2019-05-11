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

func createModelCommandPersonUpdate(
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

	if updated.BiographyHash != orig.BiographyHash {
		oldValue := orig.BiographyHash
		newValue := updated.BiographyHash
		theUpdate := &model.StringUpdate{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.BiographyHashUpdate = theUpdate
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
			eventKey:   model.PERSON_PUBLIC_KEY,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.NameUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.NameUpdate.NewValue,
			stateField: &oldPerson.Name,
			eventKey:   model.PERSON_NAME,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.EmailUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.EmailUpdate.NewValue,
			stateField: &oldPerson.Email,
			eventKey:   model.PERSON_EMAIL,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.BiographyHashUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.BiographyHashUpdate.NewValue,
			stateField: &oldPerson.BiographyHash,
			eventKey:   model.PERSON_BIOGRAPHY_HASH,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.OrganizationUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.OrganizationUpdate.NewValue,
			stateField: &oldPerson.Organization,
			eventKey:   model.PERSON_ORGANIZATION,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.TelephoneUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.TelephoneUpdate.NewValue,
			stateField: &oldPerson.Telephone,
			eventKey:   model.PERSON_TELEPHONE,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.AddressUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.AddressUpdate.NewValue,
			stateField: &oldPerson.Address,
			eventKey:   model.PERSON_ADDRESS,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.PostalCodeUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.PostalCodeUpdate.NewValue,
			stateField: &oldPerson.PostalCode,
			eventKey:   model.PERSON_POSTAL_CODE,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.CountryUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.CountryUpdate.NewValue,
			stateField: &oldPerson.Country,
			eventKey:   model.PERSON_COUNTRY,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	if c.ExtraInfoUpdate != nil {
		var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue:   c.ExtraInfoUpdate.NewValue,
			stateField: &oldPerson.ExtraInfo,
			eventKey:   model.PERSON_EXTRA_INFO,
			personId:   c.PersonId,
			timestamp:  timestamp,
		}
		result = append(result, toAppend)
	}

	return result
}
