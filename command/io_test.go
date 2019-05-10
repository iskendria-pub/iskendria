package command

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"testing"
)

func TestSettings(t *testing.T) {
	timestamp := model.GetCurrentTime()
	settings := &model.StateSettings{
		CreatedOnOn: timestamp,
		ModifiedOn:  timestamp,
		PriceList: &model.PriceList{
			PricePersonEdit: int32(1),
		},
	}
	settingsBytes, err := proto.Marshal(settings)
	if err != nil {
		t.Error("Could not marshal fake settings")
		return
	}
	us := newUnmarshalledState()
	err = us.add(map[string][]byte{model.GetSettingsAddress(): settingsBytes}, []string{model.GetSettingsAddress()})
	if err != nil {
		t.Error("Could not set person address in unmarshalledState")
	}
	regenerated, err := us.read([]string{model.GetSettingsAddress()})
	if err != nil {
		t.Error("Reading back from unmarshalledState produced an error: " + err.Error())
	}
	if len(regenerated) != 1 {
		t.Error(fmt.Sprintf("Expected to read back one address, got %d", len(regenerated)))
	}
	regeneratedBytes := regenerated[model.GetSettingsAddress()]
	if !bytes.Equal(regeneratedBytes, settingsBytes) {
		t.Error("Regenerating bytes from unmarshalledState does not produce the original data")
	}
	if us.getAddressState(model.GetSettingsAddress()) != ADDRESS_FILLED {
		t.Error("Settings address was expected to be filled")
	}
}

func TestPerson(t *testing.T) {
	timestamp := model.GetCurrentTime()
	personId := model.CreatePersonAddress()
	person := &model.StatePerson{
		Id:         personId,
		CreatedOn:  timestamp,
		ModifiedOn: timestamp,
		Name:       "Martijn",
	}
	personBytes, err := proto.Marshal(person)
	if err != nil {
		t.Error("Could not marshal example person: " + err.Error())
	}
	us := newUnmarshalledState()
	err = us.add(map[string][]byte{personId: personBytes}, []string{personId})
	if err != nil {
		t.Error("Could not add person to unmarshaledState: " + err.Error())
	}
	regenerated, err := us.read([]string{personId})
	if err != nil {
		t.Error("Could not read person from unmarshaledState")
	}
	if len(regenerated) != 1 {
		t.Error(fmt.Sprintf("Expected to read one address, but got %d", len(regenerated)))
	}
	regeneratedBytes := regenerated[personId]
	if !bytes.Equal(regeneratedBytes, personBytes) {
		t.Error(fmt.Sprintf("Read person bytes differ from original bytes"))
	}
	if us.getAddressState(personId) != ADDRESS_FILLED {
		t.Error("Person address was expected to be filled: " + personId)
	}
}

func TestEmpty(t *testing.T) {
	personId := model.CreatePersonAddress()
	us := newUnmarshalledState()
	err := us.add(make(map[string][]byte), []string{model.GetSettingsAddress(), personId})
	if err != nil {
		t.Error("Could not add empty address to unmarshalledState")
	}
	if us.getAddressState(model.GetSettingsAddress()) != ADDRESS_EMPTY {
		t.Error("unmarshalledState did not detect empty settings address")
	}
	if us.getAddressState(personId) != ADDRESS_EMPTY {
		t.Error("unmarshalledState did not detect empty person address")
	}
}

func TestMultiplePersons(t *testing.T) {
	timestamp := model.GetCurrentTime()
	firstPersonId := model.CreatePersonAddress()
	firstPerson := &model.StatePerson{
		Id:         firstPersonId,
		CreatedOn:  timestamp,
		ModifiedOn: timestamp,
		Name:       "Martijn",
	}
	secondPersonId := model.CreatePersonAddress()
	secondPerson := &model.StatePerson{
		Id:         secondPersonId,
		CreatedOn:  timestamp,
		ModifiedOn: timestamp,
		Name:       "Arri",
	}
	firstPersonBytes, err := proto.Marshal(firstPerson)
	if err != nil {
		t.Error("Could not marshal example person: " + err.Error())
	}
	secondPersonBytes, err := proto.Marshal(secondPerson)
	if err != nil {
		t.Error("Could not marshal example person: " + err.Error())
	}
	us := newUnmarshalledState()
	err = us.add(map[string][]byte{
		firstPersonId:  firstPersonBytes,
		secondPersonId: secondPersonBytes}, []string{firstPersonId, secondPersonId})
	if err != nil {
		t.Error("Could not add person to unmarshaledState: " + err.Error())
	}
	regenerated, err := us.read([]string{firstPersonId, secondPersonId})
	if err != nil {
		t.Error("Could not read person from unmarshaledState")
	}
	if len(regenerated) != 2 {
		t.Error(fmt.Sprintf("Expected to read two addresses, but got %d", len(regenerated)))
	}
	firstRegeneratedBytes := regenerated[firstPersonId]
	if !bytes.Equal(firstRegeneratedBytes, firstPersonBytes) {
		t.Error("Read person bytes differ from original bytes")
	}
	secondRegeneratedBytes := regenerated[secondPersonId]
	if !bytes.Equal(secondRegeneratedBytes, secondPersonBytes) {
		t.Error("Read person bytes differ from original bytes")
	}
	if us.getAddressState(firstPersonId) != ADDRESS_FILLED {
		t.Error("Person address was expected to be filled: " + firstPersonId)
	}
	if us.getAddressState(secondPersonId) != ADDRESS_FILLED {
		t.Error("Person address was expected to be filled: " + secondPersonId)
	}
}
