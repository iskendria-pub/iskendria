package integration

import (
	"gitlab.bbinfra.net/3estack/alexandria/cliAlexandria"
	"gitlab.bbinfra.net/3estack/alexandria/command"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"log"
	"testing"
)

func withInitializedDao(testFunc func(logger *log.Logger, t *testing.T), logger *log.Logger, t *testing.T) {
	dao.Init("testBootstrap.db", logger)
	defer dao.ShutdownAndDelete(logger)
	err := dao.StartFakeBlock("blockId", "")
	if err != nil {
		t.Error("Error starting fake block: " + err.Error())
	}
	testFunc(logger, t)
}

func withLoggedInWithNewKey(
	testFunc func(
		blockchainAccess command.BlockchainAccess,
		logger *log.Logger,
		t *testing.T),
	blockchainAccess command.BlockchainAccess,
	logger *log.Logger,
	t *testing.T) {
	publicKeyFile := "testBootstrap.pub"
	privateKeyFile := "testBootstrap.priv"
	err := cliAlexandria.CreateKeyPair(publicKeyFile, privateKeyFile)
	if err != nil {
		t.Error("Could not create keypair: " + err.Error())
	}
	defer cliAlexandria.RemoveKeyFiles(publicKeyFile, privateKeyFile, logger)
	err = cliAlexandria.Login(publicKeyFile, privateKeyFile)
	if err != nil {
		t.Error("Could not login: " + err.Error())
	}
	defer func() { _ = cliAlexandria.Logout() }()
	testFunc(blockchainAccess, logger, t)
}
