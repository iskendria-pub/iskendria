package cliAlexandria

import (
	"gitlab.bbinfra.net/3estack/alexandria/util"
	"log"
	"os"
	"strings"
	"testing"
)

func TestKeys(t *testing.T) {
	logger := log.New(os.Stdout, "testKeys", log.Flags())
	publicKeyFile := "key.pub"
	privateKeyFile := "key.priv"
	err := createKeyPair(publicKeyFile, privateKeyFile)
	if err != nil {
		t.Error("Creating key pair failed")
	}
	defer removeKeyFileles(publicKeyFile, privateKeyFile, logger)
	err = login(publicKeyFile, privateKeyFile)
	if err != nil {
		t.Error("Could not log in with created keypair")
	}
	if loggedIn.PublicKeyStr == "" || loggedIn.PublicKeyStr != strings.TrimSpace(loggedIn.PublicKeyStr) {
		t.Error("Invalid public key string stored on login")
	}
	err = logout()
	if err != nil {
		t.Error("Could not logout")
	}
	if loggedIn.PublicKeyStr != "" {
		t.Error("Public key string not cleared on logout")
	}
}

func removeKeyFileles(f1, f2 string, logger *log.Logger) {
	util.RemoveExistingFile(f1, logger)
	util.RemoveExistingFile(f2, logger)
}
