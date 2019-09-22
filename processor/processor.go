package main

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
	"github.com/iskendria-pub/iskendria/command"
	"github.com/iskendria-pub/iskendria/model"
	"log"
	"syscall"
)

func main() {
	endpoint := "tcp://127.0.0.1:4004"
	handler := &AlexandriaHandler{}
	processorObject := processor.NewTransactionProcessor(endpoint)
	processorObject.AddHandler(handler)
	processorObject.ShutdownOnSignal(syscall.SIGINT, syscall.SIGTERM)
	err := processorObject.Start()
	if err != nil {
		fmt.Println(err)
	}
}

type AlexandriaHandler struct{}

func (ch *AlexandriaHandler) FamilyName() string {
	return model.FamilyName
}

func (ch *AlexandriaHandler) FamilyVersions() []string {
	return []string{model.FamilyVersion}
}

func (ch *AlexandriaHandler) Namespaces() []string {
	return []string{model.Namespace}
}

func (ch *AlexandriaHandler) Apply(request *processor_pb2.TpProcessRequest, context *processor.Context) error {
	log.Println("Entering Apply...")
	defer log.Println("Left Apply")
	cmd, publicKey, transactionId, err := parseApply(request)
	if err != nil {
		return &processor.InvalidTransactionError{
			Msg: err.Error(),
		}
	}
	err = command.ApplyModelCommand(cmd, publicKey, transactionId, context)
	if err != nil {
		return &processor.InvalidTransactionError{
			Msg: err.Error(),
		}
	}
	return nil
}

func parseApply(request *processor_pb2.TpProcessRequest) (*model.Command, string, string, error) {
	// get public key.
	var header = request.Header
	var publicKey = header.SignerPublicKey
	var transactionId = request.Signature
	// Parse request.
	var payload = request.Payload
	var cmd = &model.Command{}
	if err := proto.Unmarshal(payload, cmd); err != nil {
		return nil, "", "", errors.New("Could not parse transaction: " + err.Error())
	}
	return cmd, publicKey, transactionId, nil
}
