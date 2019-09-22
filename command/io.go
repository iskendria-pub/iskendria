package command

import (
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
	"github.com/iskendria-pub/iskendria/dao"
	"log"
)

var _ processor_pb2.TpProcessRequest

type BlockchainAccess interface {
	GetState(addresses []string) (map[string][]byte, error)
	SetState(pairs map[string][]byte) ([]string, error)
	AddEvent(eventType string, attributes []processor.Attribute, eventData []byte) error
}

var _ BlockchainAccess = new(processor.Context)

type addressState int

const (
	ADDRESS_UNKNOWN = addressState(0)
	ADDRESS_EMPTY   = addressState(1)
	ADDRESS_FILLED  = addressState(2)
)

type blockchainStub struct {
	data         map[string][]byte
	eventHandler EventHandler
	logger       *log.Logger
}

func NewBlockchainStub(eventHandler EventHandler, logger *log.Logger) BlockchainAccess {
	return &blockchainStub{
		data:         make(map[string][]byte),
		eventHandler: eventHandler,
		logger:       logger,
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
	}, bs.logger)
}

type EventHandler func(*events_pb2.Event, *log.Logger) error

var _ EventHandler = dao.HandleEvent
