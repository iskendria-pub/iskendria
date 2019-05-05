package dao

import (
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/model"
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
	_, err = db.Exec(fmt.Sprintf("INSERT INTO person VALUES (%s)", GetPlaceHolders(17)),
		firstId, 10000, 1100000000000, firstKey, "Martijn",
		"xxx@gmail.com", false, true, 100, "abcdef01",
		"PDF", "Free University", "020", "Boelelaan", "1000 AA",
		"Netherlands", "Some extra info")
	if err != nil {
		t.Error("Could not insert first person")
		return
	}
	_, err = db.Exec(fmt.Sprintf("INSERT INTO person VALUES (%s)", GetPlaceHolders(17)),
		secondId, 12000, 13000, secondKey, "Maurice",
		"yyy.gmail.com", true, false, 200, "01234567",
		"PDF", "Sorbonne", "0033", "75005 Parijs", "75005",
		"France", "Other extra info")
	if err != nil {
		t.Error("Could not insert second person")
		return
	}
	firstPerson, err := GetPersonById(firstId)
	if err != nil {
		t.Error("With filled database, GetPersonById failed: " + err.Error())
		return
	}
	if firstPerson.Id != firstId {
		t.Error("For first person, checking id failed")
	}
	if firstPerson.CreatedOn != int64(10000) {
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
	if firstPerson.IsSigned != true {
		t.Error(" For first person, checking signed failed")
	}
	if firstPerson.Balance != int32(100) {
		t.Error("For first person, checking balance failed")
	}
	if firstPerson.BiographyHash != "abcdef01" {
		t.Error("For first person, checking bibliography hash failed")
	}
	if firstPerson.BiographyFormat != "PDF" {
		t.Error("For first person, checking bibliography format failed")
	}
	if firstPerson.Organization != "Free University" {
		t.Error("For first person, checking organization failed")
	}
	if firstPerson.Telephone != "020" {
		t.Error("For first person, checking telephone failed")
	}
	if firstPerson.Address != "Boelelaan" {
		t.Error("For first person, checking address failed")
	}
	if firstPerson.PostalCode != "1000 AA" {
		t.Error("For first person, checking postal code failed")
	}
	if firstPerson.Country != "Netherlands" {
		t.Error("For first person, checking country failed")
	}
	if firstPerson.ExtraInfo != "Some extra info" {
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
	if secondPerson.ExtraInfo != "Other extra info" {
		t.Error("For second person, checking extraInfo failed")
	}
}
