package command

import (
	"github.com/hyperledger/sawtooth-sdk-go/signing"
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
	ba BlockchainAccess) {
}
