package command

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/iskendria-pub/iskendria/model"
	"testing"
)

func TestSettings(t *testing.T) {
	timestamp := model.GetCurrentTime()
	settings := &model.StateSettings{
		CreatedOn:  timestamp,
		ModifiedOn: timestamp,
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

func TestJournal(t *testing.T) {
	timestamp := model.GetCurrentTime()
	journalId := model.CreateJournalAddress()
	journal := &model.StateJournal{
		Id:              journalId,
		CreatedOn:       timestamp,
		ModifiedOn:      timestamp,
		Title:           "Some journal",
		IsSigned:        false,
		DescriptionHash: "",
	}
	journalBytes, err := proto.Marshal(journal)
	if err != nil {
		t.Error("Could not marshal example journal " + err.Error())
	}
	us := newUnmarshalledState()
	err = us.add(map[string][]byte{journalId: journalBytes}, []string{journalId})
	if err != nil {
		t.Error("Could not add journal to unmarshalled state: " + err.Error())
	}
	regenerated, err := us.read([]string{journalId})
	if err != nil {
		t.Error("Could not read journal from unmarshalledState")
	}
	if len(regenerated) != 1 {
		t.Error(fmt.Sprintf("Expected to read one journal, but got %d", len(regenerated)))
	}
	if !bytes.Equal(regenerated[journalId], journalBytes) {
		t.Error("Regenerated bytes differ from original bytes")
	}
	if us.getAddressState(journalId) != ADDRESS_FILLED {
		t.Error("Expected journal address to be filled")
	}
}

func TestVolume(t *testing.T) {
	timestamp := model.GetCurrentTime()
	volumeId := model.CreateVolumeAddress()
	journalId := model.CreateJournalAddress()
	volume := &model.StateVolume{
		Id:        volumeId,
		CreatedOn: timestamp,
		JournalId: journalId,
		Issue:     "Some issue",
	}
	volumeBytes, err := proto.Marshal(volume)
	if err != nil {
		t.Error("Could not marshal example volume " + err.Error())
	}
	us := newUnmarshalledState()
	err = us.add(map[string][]byte{volumeId: volumeBytes}, []string{volumeId})
	if err != nil {
		t.Error("Could not add volume to unmarshalled state: " + err.Error())
	}
	regenerated, err := us.read([]string{volumeId})
	if err != nil {
		t.Error("Could not read volume from unmarshalledState")
	}
	if len(regenerated) != 1 {
		t.Error(fmt.Sprintf("Expected to read one volume, but got %d", len(regenerated)))
	}
	if !bytes.Equal(regenerated[volumeId], volumeBytes) {
		t.Error("Regenerated bytes differ from original bytes")
	}
	if us.getAddressState(volumeId) != ADDRESS_FILLED {
		t.Error("Expected volume address to be filled")
	}
}

func TestEmpty(t *testing.T) {
	personId := model.CreatePersonAddress()
	journalId := model.CreateJournalAddress()
	volumeId := model.CreateVolumeAddress()
	us := newUnmarshalledState()
	err := us.add(make(map[string][]byte), []string{model.GetSettingsAddress(), personId, journalId, volumeId})
	if err != nil {
		t.Error("Could not add empty addresses to unmarshalledState")
	}
	if us.getAddressState(model.GetSettingsAddress()) != ADDRESS_EMPTY {
		t.Error("unmarshalledState did not detect empty settings address")
	}
	if us.getAddressState(personId) != ADDRESS_EMPTY {
		t.Error("unmarshalledState did not detect empty person address")
	}
	if us.getAddressState(journalId) != ADDRESS_EMPTY {
		t.Error("unmarshalledState did not detect empty journal address")
	}
	if us.getAddressState(volumeId) != ADDRESS_EMPTY {
		t.Error("unmarshalledState did not detect empty volume address")
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
