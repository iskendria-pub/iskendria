package dao

import (
	"github.com/iskendria-pub/iskendria/model"
	"log"
	"os"
	"testing"
)

func TestGetPerson(t *testing.T) {
	logger := log.New(os.Stdout, "testGetSettings", log.Flags())
	Init("testGetPerson.db", logger)
	defer ShutdownAndDelete(logger)
	firstId := model.CreatePersonAddress()
	secondId := model.CreatePersonAddress()
	firstKey := "someFirstKey"
	secondKey := "someSecondKey"
	expectedNil, err := GetPersonById(firstId)
	if err != nil {
		t.Error("When person table empty, error when querying person: " + err.Error())
	}
	if expectedNil != nil {
		t.Error("When person table empty, got non-nil person")
	}
	expectedLenZero, err := SearchPersonByKey(firstKey)
	if err != nil {
		t.Error("When person table empty, error when querying person by key: " + err.Error())
	}
	if len(expectedLenZero) != 0 {
		t.Error("When person table empty, got persons from key")
	}
	firstPersonCreate := &dataManipulationPersonCreate{
		id:        firstId,
		timestamp: 1100000000000,
		publicKey: firstKey,
		name:      "Martijn",
		email:     "xxx@gmail.com",
	}
	applyPersonCreate(firstPersonCreate, t)
	secondPersonCreate := &dataManipulationPersonCreate{
		id:        secondId,
		timestamp: 12000,
		publicKey: secondKey,
		name:      "Maurice",
		email:     "yyy.gmail.com",
	}
	applyPersonCreate(secondPersonCreate, t)
	firstPerson, err := GetPersonById(firstId)
	if err != nil {
		t.Error("With filled database, GetPersonById failed: " + err.Error())
		return
	}
	if firstPerson.Id != firstId {
		t.Error("For first person, checking id failed")
	}
	if firstPerson.CreatedOn != int64(1100000000000) {
		t.Error("For first person, checking createdOn failed")
	}
	if firstPerson.ModifiedOn != int64(1100000000000) {
		t.Error("For first person, checking modifiedOn failed")
	}
	if firstPerson.PublicKey != firstKey {
		t.Error("For first person, checking public key failed")
	}
	if firstPerson.Name != "Martijn" {
		t.Error("For first person, checking name failed")
	}
	if firstPerson.Email != "xxx@gmail.com" {
		t.Error("For first person, checking email failed")
	}
	if firstPerson.IsMajor != false {
		t.Error("For first person, checking majorship failed")
	}
	if firstPerson.IsSigned != false {
		t.Error(" For first person, checking signed failed")
	}
	if firstPerson.Balance != int32(0) {
		t.Error("For first person, checking balance failed")
	}
	if firstPerson.BiographyHash != "" {
		t.Error("For first person, checking bibliography hash failed")
	}
	if firstPerson.Organization != "" {
		t.Error("For first person, checking organization failed")
	}
	if firstPerson.Telephone != "" {
		t.Error("For first person, checking telephone failed")
	}
	if firstPerson.Address != "" {
		t.Error("For first person, checking address failed")
	}
	if firstPerson.PostalCode != "" {
		t.Error("For first person, checking postal code failed")
	}
	if firstPerson.Country != "" {
		t.Error("For first person, checking country failed")
	}
	if firstPerson.ExtraInfo != "" {
		t.Error("For first person, checking extraInfo failed")
	}
	persons, err := SearchPersonByKey(secondKey)
	if err != nil {
		t.Error("With filled database, SearchPersonByKey failed: " + err.Error())
		return
	}
	if len(persons) != 1 {
		t.Error("Expected exactly one person from SearchPersonByKey")
		return
	}
	secondPerson := persons[0]
	if secondPerson.Id != secondId {
		t.Error("For second person, checking id failed")
	}
	if secondPerson.PublicKey != secondKey {
		t.Error("For second person, checking key failed")
	}
	if secondPerson.ExtraInfo != "" {
		t.Error("For second person, checking extraInfo failed")
	}
}

func applyPersonCreate(pc *dataManipulationPersonCreate, t *testing.T) {
	tx, err := db.Beginx()
	if err != nil {
		t.Error("Error when starting transaction: " + err.Error())
	}
	err = pc.apply(tx)
	if err != nil {
		t.Error("Error applying data manipulation person create: " + err.Error())
	}
	err = tx.Commit()
	if err != nil {
		t.Error("Error committing transaction: " + err.Error())
	}
}
