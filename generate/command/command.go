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

func createModelCommandPersonUpdate(
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
	DaoSettingsUpdate             *Update
	DaoPersonUpdateProperties     *Update
	CommandPersonUpdateProperties []*modelCommandCheck
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
		CommandPersonUpdateProperties: getCommandPersonUpdatePropertiesFields(),
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

func getCommandPersonUpdatePropertiesFields() []*modelCommandCheck {
	return []*modelCommandCheck{
		{
			CommandField: "PublicKey",
			EventKey:     "PERSON_PUBLIC_KEY",
		},
		{
			CommandField: "Name",
			EventKey:     "PERSON_NAME",
		},
		{
			CommandField: "Email",
			EventKey:     "PERSON_EMAIL",
		},
		{
			CommandField: "BiographyHash",
			EventKey:     "PERSON_BIOGRAPHY_HASH",
		},
		{
			CommandField: "Organization",
			EventKey:     "PERSON_ORGANIZATION",
		},
		{
			CommandField: "Telephone",
			EventKey:     "PERSON_TELEPHONE",
		},
		{
			CommandField: "Address",
			EventKey:     "PERSON_ADDRESS",
		},
		{
			CommandField: "PostalCode",
			EventKey:     "PERSON_POSTAL_CODE",
		},
		{
			CommandField: "Country",
			EventKey:     "PERSON_COUNTRY",
		},
		{
			CommandField: "ExtraInfo",
			EventKey:     "PERSON_EXTRA_INFO",
		},
	}
}
