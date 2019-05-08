package command

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"gitlab.bbinfra.net/3estack/alexandria/util"
)

// Test running a command and check that only its input and output addresses
// are accessed.
func RunCommandForTest(c *Command, transactionId string, ba BlockchainAccess) error {
	return ApplyModelCommand(
		c.Command,
		c.CryptoIdentity.PublicKeyStr,
		transactionId,
		newAccessCheckingDecorator(ba, c.InputAddresses, c.OutputAddresses))
}

type accessCheckingDecorator struct {
	readableAddresses map[string]bool
	writableAddresses map[string]bool
	delegate          BlockchainAccess
}

func newAccessCheckingDecorator(ba BlockchainAccess, inputAddresses, outputAddresses []string) BlockchainAccess {
	return &accessCheckingDecorator{
		readableAddresses: util.StringSliceToSet(inputAddresses),
		writableAddresses: util.StringSliceToSet(outputAddresses),
		delegate:          ba,
	}
}

func (acd *accessCheckingDecorator) GetState(addresses []string) (map[string][]byte, error) {
	if !util.StringSetHasAll(acd.readableAddresses, addresses) {
		toReport := getAddressesNotIn(addresses, acd.readableAddresses)
		return nil, errors.New(fmt.Sprintf("Some addresses were not readable: %v", toReport))
	}
	return acd.delegate.GetState(addresses)
}

func getAddressesNotIn(addresses []string, theSet map[string]bool) []string {
	toReport := make([]string, 0, len(addresses))
	for _, a := range addresses {
		_, isIn := theSet[a]
		if !isIn {
			toReport = append(toReport, a)
		}
	}
	return toReport
}

func (acd *accessCheckingDecorator) SetState(pairs map[string][]byte) ([]string, error) {
	writtenAddresses := util.StringSetToSlice(pairs)
	if !util.StringSetHasAll(acd.writableAddresses, writtenAddresses) {
		toReport := getAddressesNotIn(writtenAddresses, acd.writableAddresses)
		return nil, errors.New(fmt.Sprintf("Some addresses were not writable: %v", toReport))
	}
	return acd.delegate.SetState(pairs)
}

func (acd *accessCheckingDecorator) AddEvent(
	eventType string, attributes []processor.Attribute, eventData []byte) error {
	return acd.delegate.AddEvent(eventType, attributes, eventData)
}

type blockchainStub struct {
	data         map[string][]byte
	eventHandler EventHandler
}

func NewBlockchainStub(eventHandler EventHandler) BlockchainAccess {
	return &blockchainStub{
		data:         make(map[string][]byte),
		eventHandler: eventHandler,
	}
}

func (bs *blockchainStub) GetState(addresses []string) (map[string][]byte, error) {
	result := make(map[string][]byte)
	for _, a := range addresses {
		contents, found := bs.data[a]
		if found {
			result[a] = contents
		}
	}
	return result, nil
}

func (bs *blockchainStub) SetState(pairs map[string][]byte) ([]string, error) {
	result := make([]string, 0)
	for address, contents := range pairs {
		bs.data[address] = contents
		result = append(result, address)
	}
	return result, nil
}

func (bs *blockchainStub) AddEvent(eventType string, attributes []processor.Attribute, eventData []byte) error {
	eventAddtributes := make([]*events_pb2.Event_Attribute, 0)
	for _, a := range attributes {
		eventAddtributes = append(eventAddtributes, &events_pb2.Event_Attribute{
			Key:   a.Key,
			Value: a.Value,
		})
	}
	return bs.eventHandler(&events_pb2.Event{
		EventType:  eventType,
		Attributes: eventAddtributes,
		Data:       eventData,
	})
}

type EventHandler func(*events_pb2.Event) error

var _ EventHandler = dao.HandleEvent
