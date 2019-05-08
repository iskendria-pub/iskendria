package command

import (
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"gitlab.bbinfra.net/3estack/alexandria/model"
)

type Command struct {
	InputAddresses  []string
	OutputAddresses []string
	CryptoIdentity  *CryptoIdentity
	Command         *model.Command
}

type CryptoIdentity struct {
	PublicKeyStr string
	PublicKey    signing.PublicKey
	PrivateKey   signing.PrivateKey
}

func ApplyModelCommand(
	command *model.Command,
	signerKey string,
	transactionId string,
	ba BlockchainAccess) error {
	// TODO: Implement.
	return nil
}

type EventHandler func(*events_pb2.Event) error

var _ EventHandler = dao.HandleEvent
