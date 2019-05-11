package command

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"gitlab.bbinfra.net/3estack/alexandria/model"
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

type unmarshalledState struct {
	emptyAddresses map[string]bool
	settings       *model.StateSettings
	persons        map[string]*model.StatePerson
}

func newUnmarshalledState() *unmarshalledState {
	return &unmarshalledState{
		emptyAddresses: make(map[string]bool),
		settings:       nil,
		persons:        make(map[string]*model.StatePerson),
	}
}

func (us *unmarshalledState) getAddressState(address string) addressState {
	_, isEmpty := us.emptyAddresses[address]
	if isEmpty {
		return ADDRESS_EMPTY
	}
	switch {
	case address == model.GetSettingsAddress():
		if us.settings != nil {
			return ADDRESS_FILLED
		}
	case model.IsPersonAddress(address):
		_, found := us.persons[address]
		if found {
			return ADDRESS_FILLED
		}
	}
	return ADDRESS_UNKNOWN
}

func (us *unmarshalledState) add(readData map[string][]byte, requestedAddresses []string) error {
	for _, ra := range requestedAddresses {
		var err error
		contents, isAvailable := readData[ra]
		if isAvailable {
			err = us.addAvailable(ra, contents)
		} else {
			us.emptyAddresses[ra] = true
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (us *unmarshalledState) addAvailable(address string, contents []byte) error {
	var err error
	switch {
	case address == model.GetSettingsAddress():
		err = us.addSettings(contents)
	case model.IsPersonAddress(address):
		err = us.addPerson(address, contents)
	}
	return err
}

func (us *unmarshalledState) addSettings(contents []byte) error {
	settings := &model.StateSettings{}
	err := proto.Unmarshal(contents, settings)
	if err != nil {
		return err
	}
	us.settings = settings
	return nil
}

func (us *unmarshalledState) addPerson(personId string, contents []byte) error {
	person := &model.StatePerson{}
	err := proto.Unmarshal(contents, person)
	if err != nil {
		return err
	}
	us.persons[personId] = person
	return nil
}

func (us *unmarshalledState) read(addresses []string) (map[string][]byte, error) {
	result := make(map[string][]byte)
	var err error
	for _, address := range addresses {
		switch {
		case address == model.GetSettingsAddress():
			err = us.readSettings(result)
		case model.IsPersonAddress(address):
			err = us.readPerson(address, result)
		}
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func (us *unmarshalledState) readSettings(result map[string][]byte) error {
	marshalled, err := proto.Marshal(us.settings)
	if err != nil {
		return err
	}
	result[model.GetSettingsAddress()] = marshalled
	return nil
}

func (us *unmarshalledState) readPerson(personId string, result map[string][]byte) error {
	marshalled, err := proto.Marshal(us.persons[personId])
	if err != nil {
		return err
	}
	result[personId] = marshalled
	return nil
}