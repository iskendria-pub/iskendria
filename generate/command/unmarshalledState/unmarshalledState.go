package main

import (
	"fmt"
	"os"
	"text/template"
)

var templateUnmarshalledState = `
package command

import (
    "gitlab.bbinfra.net/3estack/alexandria/model"
  	proto "github.com/golang/protobuf/proto"
)

type unmarshalledState struct {
	emptyAddresses map[string]bool
	settings       *model.StateSettings
    {{range . -}}
    {{.UnmarshalledContainerField}} map[string]*model.{{.ModelStateField}}
    {{ end -}}
}

func newUnmarshalledState() *unmarshalledState {
	return &unmarshalledState{
		emptyAddresses: make(map[string]bool),
		settings:       nil,
        {{range . -}}
		{{.UnmarshalledContainerField}}: make(map[string]*model.{{.ModelStateField}}),
        {{end -}}
	}
}

func (us *unmarshalledState) getAddressState(address string) addressState {
	_, isEmpty := us.emptyAddresses[address]
	if isEmpty {
		return ADDRESS_EMPTY
	}
	switch {
	case address == model.GetSettingsAddress():
		if us.settings != nil {
			return ADDRESS_FILLED
		}
    {{range . -}}
	case model.{{.ModelAddressTypeChecker}}(address):
		_, found := us.{{.UnmarshalledContainerField}}[address]
		if found {
			return ADDRESS_FILLED
		}
    {{end -}}
	}
	return ADDRESS_UNKNOWN
}

func (us *unmarshalledState) add(readData map[string][]byte, requestedAddresses []string) error {
	for _, ra := range requestedAddresses {
		var err error
		contents, isAvailable := readData[ra]
		if isAvailable {
			err = us.addAvailable(ra, contents)
		} else {
			us.emptyAddresses[ra] = true
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (us *unmarshalledState) addAvailable(address string, contents []byte) error {
	var err error
	switch {
	case address == model.GetSettingsAddress():
		err = us.addSettings(contents)
    {{range . -}}	
    case model.{{.ModelAddressTypeChecker}}(address):
		err = us.add{{.Tag}}(address, contents)
    {{end -}}
    }
	return err
}

func (us *unmarshalledState) addSettings(contents []byte) error {
	settings := &model.StateSettings{}
	err := proto.Unmarshal(contents, settings)
	if err != nil {
		return err
	}
	us.settings = settings
	return nil
}

{{range . -}}
func (us *unmarshalledState) add{{.Tag}}(theId string, contents []byte) error {
	modelContainer := &model.{{.ModelStateField}}{}
	err := proto.Unmarshal(contents, modelContainer)
	if err != nil {
		return err
	}
	us.{{.UnmarshalledContainerField}}[theId] = modelContainer
	return nil
}
{{end -}}

func (us *unmarshalledState) read(addresses []string) (map[string][]byte, error) {
	result := make(map[string][]byte)
	var err error
	for _, address := range addresses {
		switch {
		case address == model.GetSettingsAddress():
			err = us.readSettings(result)
        {{range . -}}
		case model.{{.ModelAddressTypeChecker}}(address):
			err = us.read{{.Tag}}(address, result)
        {{end -}}
		}
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func (us *unmarshalledState) readSettings(result map[string][]byte) error {
	marshalled, err := proto.Marshal(us.settings)
	if err != nil {
		return err
	}
	result[model.GetSettingsAddress()] = marshalled
	return nil
}

{{range . -}}
func (us *unmarshalledState) read{{.Tag}}(theId string, result map[string][]byte) error {
	marshalled, err := proto.Marshal(us.{{.UnmarshalledContainerField}}[theId])
	if err != nil {
		return err
	}
	result[theId] = marshalled
	return nil
}
{{end -}}
`

type Config struct {
	Tag                        string
	UnmarshalledContainerField string
	ModelStateField            string
	ModelAddressTypeChecker    string
}

func main() {
	c := []Config{
		{
			Tag:                        "Person",
			UnmarshalledContainerField: "persons",
			ModelStateField:            "StatePerson",
			ModelAddressTypeChecker:    "IsPersonAddress",
		},
		{
			Tag:                        "Journal",
			UnmarshalledContainerField: "journals",
			ModelStateField:            "StateJournal",
			ModelAddressTypeChecker:    "IsJournalAddress",
		},
		{
			Tag:                        "Volume",
			UnmarshalledContainerField: "volumes",
			ModelStateField:            "StateVolume",
			ModelAddressTypeChecker:    "IsVolumeAddress",
		},
		{
			Tag:                        "Manuscript",
			UnmarshalledContainerField: "manuscripts",
			ModelStateField:            "StateManuscript",
			ModelAddressTypeChecker:    "IsManuscriptAddress",
		},
		{
			Tag:                        "ManuscriptThread",
			UnmarshalledContainerField: "manuscriptThreads",
			ModelStateField:            "StateManuscriptThread",
			ModelAddressTypeChecker:    "IsManuscriptThreadAddress",
		},
	}
	tmpl, err := template.New("templateUnmarshalledState").Parse(templateUnmarshalledState)
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
