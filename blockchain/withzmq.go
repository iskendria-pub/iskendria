package blockchain

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/messaging"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_event_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"
	"github.com/pebbe/zmq4"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"log"
	"os"
)

const zmqEventPort = "4004"

var EventStreamStatusChannel = make(chan *EventStreamStatus, eventStreamStatusChannelCapacity)

const eventStreamStatusChannelCapacity = 20

type EventStreamStatus struct {
	StatusCode EventStreamStatusCode
	Msg        string
}

type EventStreamStatusCode int32

const (
	EVENT_STREAM_STATUS_STOPPED      = EventStreamStatusCode(0)
	EVENT_STREAM_STATUS_INITIALIZING = EventStreamStatusCode(1)
	EVENT_STREAM_STATUS_RUNNING      = EventStreamStatusCode(2)
)

func RequestEventStream(eventHandler dao.EventHandler, fname, tag string) {
	file, err := os.Create(fname)
	if err != nil {
		EventStreamStatusChannel <- &EventStreamStatus{
			StatusCode: EVENT_STREAM_STATUS_STOPPED,
			Msg: fmt.Sprintf("Could not create logfile %s, error %s",
				fname, err.Error()),
		}
		return
	}
	defer func() { _ = file.Close() }()
	logger := log.New(file, tag, log.Flags())
	firstBlock, err := getFirstBlock(logger)
	if err != nil {
		EventStreamStatusChannel <- &EventStreamStatus{
			StatusCode: EVENT_STREAM_STATUS_STOPPED,
			Msg:        "Could not get first block: " + err.Error(),
		}
		return
	}
	eventRequester := func(c *messaging.ZmqConnection) error {
		return handleEvents(eventHandler, c, logger)
	}
	withSawtoothZmqConnection(eventRequester, firstBlock)
}

func handleEvents(eventHandler dao.EventHandler, connection *messaging.ZmqConnection, _ *log.Logger) error {
	EventStreamStatusChannel <- &EventStreamStatus{
		StatusCode: EVENT_STREAM_STATUS_RUNNING,
		Msg:        "Requesting events from blockchain...",
	}
	defer func() {
		EventStreamStatusChannel <- &EventStreamStatus{
			StatusCode: EVENT_STREAM_STATUS_STOPPED,
			Msg:        "Stopped.",
		}
	}()
	for {
		_, message, err := connection.RecvMsg()
		if err != nil {
			return err
		}
		err = handleEventsMessage(message, eventHandler)
		if err != nil {
			return err
		}
	}
}

func withSawtoothZmqConnection(theEventListener eventListener, firstBlock string) {
	err := doWithSawtoothZmqConnection(theEventListener, firstBlock)
	if err != nil {
		EventStreamStatusChannel <- &EventStreamStatus{
			StatusCode: EVENT_STREAM_STATUS_STOPPED,
			Msg:        err.Error(),
		}
	}
}

func doWithSawtoothZmqConnection(theEventListner eventListener, firstBlock string) error {
	ctx, err := zmq4.NewContext()
	if err != nil {
		return err
	}
	url := "tcp://" + getIp() + ":" + zmqEventPort
	connection, err := messaging.NewConnection(ctx, zmq4.DEALER, url, false)
	if err != nil {
		return err
	}
	defer connection.Close()
	request := &client_event_pb2.ClientEventsSubscribeRequest{
		Subscriptions:     getEventSubscriptions(),
		LastKnownBlockIds: []string{firstBlock},
	}
	serializedRequest, err := proto.Marshal(request)
	if err != nil {
		return err
	}
	corrId, err := connection.SendNewMsg(
		validator_pb2.Message_CLIENT_EVENTS_SUBSCRIBE_REQUEST,
		serializedRequest)
	if err != nil {
		return err
	}
	_, response, err := connection.RecvMsgWithId(corrId)
	if err != nil {
		return err
	}
	responseContent := &client_event_pb2.ClientEventsSubscribeResponse{}
	err = proto.Unmarshal(response.Content, responseContent)
	if err != nil {
		EventStreamStatusChannel <- &EventStreamStatus{
			StatusCode: EVENT_STREAM_STATUS_STOPPED,
			Msg: fmt.Sprintf("Could not unmarshal ClientEventsSubscribeResponse, error: %s",
				err.Error()),
		}
		return err
	}
	if responseContent.Status != client_event_pb2.ClientEventsSubscribeResponse_OK {
		return errors.New("ClientEventsSubscribeResponse is negative")
	}
	defer sendCloseSubscriptionRequest(connection)
	return theEventListner(connection)
}

type eventListener func(*messaging.ZmqConnection) error

func getEventSubscriptions() []*events_pb2.EventSubscription {
	result := make([]*events_pb2.EventSubscription, len(dao.AllEventTypes))
	for i, et := range dao.AllEventTypes {
		result[i] = &events_pb2.EventSubscription{
			EventType: et,
			Filters:   []*events_pb2.EventFilter{},
		}
	}
	return result
}

func sendCloseSubscriptionRequest(connection *messaging.ZmqConnection) {
	eventsUnsubscribeRequest := &client_event_pb2.ClientEventsUnsubscribeRequest{}
	serializedUnsubscibeRequest, err := proto.Marshal(eventsUnsubscribeRequest)
	if err != nil {
		panic(err)
	}
	corrId, err := connection.SendNewMsg(
		validator_pb2.Message_CLIENT_EVENTS_UNSUBSCRIBE_REQUEST, serializedUnsubscibeRequest)
	if err != nil {
		panic(err)
	}
	_, eventsUnsubscribeResponse, err := connection.RecvMsgWithId(corrId)
	if err != nil {
		panic(err)
	}
	responseContent := &client_event_pb2.ClientEventsUnsubscribeResponse{}
	err = proto.Unmarshal(eventsUnsubscribeResponse.Content, responseContent)
	if err != nil {
		panic(err)
	}
	if responseContent.Status != client_event_pb2.ClientEventsUnsubscribeResponse_OK {
		panic("Negative response on client events unsubscribe request")
	}
}

func handleEventsMessage(message *validator_pb2.Message, eventHandler dao.EventHandler) error {
	if message.MessageType != validator_pb2.Message_CLIENT_EVENTS {
		return errors.New("Received message is not requested for")
	}
	event_list := &events_pb2.EventList{}
	err := proto.Unmarshal(message.Content, event_list)
	if err != nil {
		return err
	}
	for _, event := range event_list.Events {
		err = eventHandler(event)
		if err != nil {
			return err
		}
	}
	return nil
}
