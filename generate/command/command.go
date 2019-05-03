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

func createModelCommandPersonUpdate(orig, updated *dao.PersonUpdate) *model.CommandPersonUpdateProperties {
    result := &model.CommandPersonUpdateProperties{}
{{range .DaoPersonUpdateProperties.UpdatesWithKinds}}
    {{template "Update" .}}
{{end}}
    return result
}
`

type Config struct {
	DaoSettingsUpdate         *Update
	DaoPersonUpdateProperties *Update
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

func main() {
	c := &Config{
		DaoSettingsUpdate: &Update{
			Kind: "Int",
			Fields: []string{
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
			},
		},
		DaoPersonUpdateProperties: &Update{
			Kind:   "String",
			Fields: getDaoPersonUpdatePropertiesFields(),
		},
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

func getDaoPersonUpdatePropertiesFields() []string {
	t := reflect.TypeOf(dao.PersonUpdate{})
	result := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		result[i] = t.Field(i).Name
	}
	return result
}
