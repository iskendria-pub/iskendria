package main

import (
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"os"
	"reflect"
	"text/template"
)

var templateCommand = `
package command

import (
	"errors"
	"fmt"
    "gitlab.bbinfra.net/3estack/alexandria/dao"
    "gitlab.bbinfra.net/3estack/alexandria/model"
)

{{define "Update"}}
	if updated.{{.Field}} != orig.{{.Field}} {
		oldValue := orig.{{.Field}}
		newValue := updated.{{.Field}}
		theUpdate := &model.{{.Kind}}Update{
			OldValue: oldValue,
			NewValue: newValue,
		}
		result.{{.Field}}Update = theUpdate
	}
{{end}}

func createModelCommandSettingsUpdate(orig, updated *dao.Settings) *model.CommandSettingsUpdate {
	result := &model.CommandSettingsUpdate{}
{{range .DaoSettingsUpdate.UpdatesWithKinds}}
	{{template "Update" .}}
{{end}}
    return result
}

func checkModelCommandSettingsUpdate(
c *model.CommandSettingsUpdate, oldSettings *model.StateSettings) error {
	{{range .CommandSettingsUpdateProperties}}
	if c.{{.CommandField}}Update != nil && c.{{.CommandField}}Update.OldValue != oldSettings.PriceList.{{.CommandField}} {
		return errors.New(fmt.Sprintf("{{.CommandField}} mismatch. Expected %d, got %d",
			c.{{.CommandField}}Update.OldValue, oldSettings.PriceList.{{.CommandField}}))
	}
	{{end}}
    return nil
}

func createSingleUpdatesSettingsUpdate(
c *model.CommandSettingsUpdate, oldSettings *model.StateSettings, timestamp int64) []singleUpdate {
	result := []singleUpdate{}
{{range .CommandSettingsUpdateProperties}}
	if c.{{.CommandField}}Update != nil {
        var toAppend singleUpdate = &singleUpdateSettingsUpdate{
			newValue: c.{{.CommandField}}Update.NewValue,
			stateField: &oldSettings.PriceList.{{.CommandField}},
			eventKey: model.{{.EventKey}},
			timestamp: timestamp,
		}
		result = append(result, toAppend)
	}
{{end}}
    return result
}

func createModelCommandPersonUpdateProperties(
personId string,
orig, updated *dao.PersonUpdate) *model.CommandPersonUpdateProperties {
    result := &model.CommandPersonUpdateProperties{}
    result.PersonId = personId
{{range .DaoPersonUpdateProperties.UpdatesWithKinds}}
    {{template "Update" .}}
{{end}}
    return result
}

func checkModelCommandPersonUpdateProperties(
c *model.CommandPersonUpdateProperties, oldPerson *model.StatePerson) error {
	{{range .CommandPersonUpdateProperties}}
	if c.{{.CommandField}}Update != nil && c.{{.CommandField}}Update.OldValue != oldPerson.{{.CommandField}} {
		return errors.New(fmt.Sprintf("Person update properties value mismatch. Expected %s, got %s",
			c.{{.CommandField}}Update.OldValue, oldPerson.{{.CommandField}}))
	}
	{{end}}
    return nil
}

func createSingleUpdatesPersonUpdateProperties(
c *model.CommandPersonUpdateProperties, oldPerson *model.StatePerson, timestamp int64) []singleUpdate {
	result := []singleUpdate{}
{{range .CommandPersonUpdateProperties}}
	if c.{{.CommandField}}Update != nil {
        var toAppend singleUpdate = &singleUpdatePersonPropertyUpdate{
			newValue: c.{{.CommandField}}Update.NewValue,
			stateField: &oldPerson.{{.CommandField}},
			eventKey: model.{{.EventKey}},
			personId: c.PersonId,
			timestamp: timestamp,
		}
		result = append(result, toAppend)
	}
{{end}}
    return result
}
`

type Config struct {
	DaoSettingsUpdate               *Update
	DaoPersonUpdateProperties       *Update
	CommandSettingsUpdateProperties []*modelCommandCheck
	CommandPersonUpdateProperties   []*modelCommandCheck
}

type Update struct {
	Kind   string
	Fields []string
}

func (u *Update) UpdatesWithKinds() []*UpdateWithKind {
	result := make([]*UpdateWithKind, len(u.Fields))
	for i, f := range u.Fields {
		result[i] = &UpdateWithKind{
			Kind:  u.Kind,
			Field: f,
		}
	}
	return result
}

type UpdateWithKind struct {
	Kind  string
	Field string
}

type modelCommandCheck struct {
	CommandField string
	EventKey     string
}

func main() {
	c := &Config{
		DaoSettingsUpdate: &Update{
			Kind:   "Int",
			Fields: getDaoSettingsUpdatePropertiesFields(),
		},
		DaoPersonUpdateProperties: &Update{
			Kind:   "String",
			Fields: getDaoPersonUpdatePropertiesFields(),
		},
		CommandSettingsUpdateProperties: getCommandSettingsUpdatePropertiesFields(),
		CommandPersonUpdateProperties:   getCommandPersonUpdatePropertiesFields(),
	}
	tmpl, err := template.New("templateCommand").Parse(templateCommand)
	if err != nil {
		fmt.Println("Error parsing template")
		fmt.Println(err)
		return
	}
	err = tmpl.Execute(os.Stdout, c)
	if err != nil {
		fmt.Println(err)
	}
}

func getDaoSettingsUpdatePropertiesFields() []string {
	return []string{
		"PriceMajorEditSettings",
		"PriceMajorCreatePerson",
		"PriceMajorChangePersonAuthorization",
		"PriceMajorChangeJournalAuthorization",
		"PricePersonEdit",
		"PriceAuthorSubmitNewManuscript",
		"PriceAuthorSubmitNewVersion",
		"PriceAuthorAcceptAuthorship",
		"PriceReviewerSubmit",
		"PriceEditorAllowManuscriptReview",
		"PriceEditorRejectManuscript",
		"PriceEditorPublishManuscript",
		"PriceEditorAssignManuscript",
		"PriceEditorCreateJournal",
		"PriceEditorCreateVolume",
		"PriceEditorEditJournal",
		"PriceEditorAddColleague",
		"PriceEditorAcceptDuty",
	}
}

func getDaoPersonUpdatePropertiesFields() []string {
	t := reflect.TypeOf(dao.PersonUpdate{})
	result := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		result[i] = t.Field(i).Name
	}
	return result
}

func getCommandSettingsUpdatePropertiesFields() []*modelCommandCheck {
	return []*modelCommandCheck{
		{
			CommandField: "PriceMajorEditSettings",
			EventKey:     "EV_KEY_PRICE_MAJOR_EDIT_SETTINGS",
		},
		{
			CommandField: "PriceMajorCreatePerson",
			EventKey:     "EV_KEY_PRICE_MAJOR_CREATE_PERSON",
		},
		{
			CommandField: "PriceMajorChangePersonAuthorization",
			EventKey:     "EV_KEY_PRICE_MAJOR_CHANGE_PERSON_AUTHORIZATION",
		},
		{
			CommandField: "PriceMajorChangeJournalAuthorization",
			EventKey:     "EV_KEY_PRICE_MAJOR_CHANGE_JOURNAL_AUTHORIZATION",
		},
		{
			CommandField: "PricePersonEdit",
			EventKey:     "EV_KEY_PRICE_PERSON_EDIT",
		},
		{
			CommandField: "PriceAuthorSubmitNewManuscript",
			EventKey:     "EV_KEY_PRICE_AUTHOR_SUBMIT_NEW_MANUSCRIPT",
		},
		{
			CommandField: "PriceAuthorSubmitNewVersion",
			EventKey:     "EV_KEY_PRICE_AUTHOR_SUBMIT_NEW_VERSION",
		},
		{
			CommandField: "PriceAuthorAcceptAuthorship",
			EventKey:     "EV_KEY_PRICE_AUTHOR_ACCEPT_AUTHORSHIP",
		},
		{
			CommandField: "PriceReviewerSubmit",
			EventKey:     "EV_KEY_PRICE_REVIEWER_SUBMIT",
		},
		{
			CommandField: "PriceEditorAllowManuscriptReview",
			EventKey:     "EV_KEY_PRICE_EDITOR_ALLOW_MANUSCRIPT_REVIEW",
		},
		{
			CommandField: "PriceEditorRejectManuscript",
			EventKey:     "EV_KEY_PRICE_EDITOR_REJECT_MANUSCRIPT",
		},
		{
			CommandField: "PriceEditorPublishManuscript",
			EventKey:     "EV_KEY_PRICE_EDITOR_PUBLISH_MANUSCRIPT",
		},
		{
			CommandField: "PriceEditorAssignManuscript",
			EventKey:     "EV_KEY_PRICE_EDITOR_ASSIGN_MANUSCRIPT",
		},
		{
			CommandField: "PriceEditorCreateJournal",
			EventKey:     "EV_KEY_PRICE_EDITOR_CREATE_JOURNAL",
		},
		{
			CommandField: "PriceEditorCreateVolume",
			EventKey:     "EV_KEY_PRICE_EDITOR_CREATE_VOLUME",
		},
		{
			CommandField: "PriceEditorEditJournal",
			EventKey:     "EV_KEY_PRICE_EDITOR_EDIT_JOURNAL",
		},
		{
			CommandField: "PriceEditorAddColleague",
			EventKey:     "EV_KEY_PRICE_EDITOR_ADD_COLLEAGUE",
		},
		{
			CommandField: "PriceEditorAcceptDuty",
			EventKey:     "EV_KEY_PRICE_EDITOR_ACCEPT_DUTY",
		},
	}
}

func getCommandPersonUpdatePropertiesFields() []*modelCommandCheck {
	return []*modelCommandCheck{
		{
			CommandField: "PublicKey",
			EventKey:     "EV_KEY_PERSON_PUBLIC_KEY",
		},
		{
			CommandField: "Name",
			EventKey:     "EV_KEY_PERSON_NAME",
		},
		{
			CommandField: "Email",
			EventKey:     "EV_KEY_PERSON_EMAIL",
		},
		{
			CommandField: "BiographyHash",
			EventKey:     "EV_KEY_PERSON_BIOGRAPHY_HASH",
		},
		{
			CommandField: "Organization",
			EventKey:     "EV_KEY_PERSON_ORGANIZATION",
		},
		{
			CommandField: "Telephone",
			EventKey:     "EV_KEY_PERSON_TELEPHONE",
		},
		{
			CommandField: "Address",
			EventKey:     "EV_KEY_PERSON_ADDRESS",
		},
		{
			CommandField: "PostalCode",
			EventKey:     "EV_KEY_PERSON_POSTAL_CODE",
		},
		{
			CommandField: "Country",
			EventKey:     "EV_KEY_PERSON_COUNTRY",
		},
		{
			CommandField: "ExtraInfo",
			EventKey:     "EV_KEY_PERSON_EXTRA_INFO",
		},
	}
}
