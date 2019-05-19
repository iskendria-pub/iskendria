package blockchain

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/messaging"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_event_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"github.com/pebbe/zmq4"
	"gitlab.bbinfra.net/3estack/alexandria/cli"
	"gitlab.bbinfra.net/3estack/alexandria/command"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

const envVarIp = "ALEXANDRIA_IP"
const restPort = "8008"
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

func getFirstBlock(logger *log.Logger) (string, error) {
	url := fmt.Sprintf("http://%s:%s/blocks?reverse", getIp(), restPort)
	logger.Printf("Reading %s...\n", url)
	response, err := http.Get(url)
	defer func() { _ = response.Body.Close }()
	if err != nil {
		return "", err
	}
	logger.Printf("Got status code %d\n", response.StatusCode)
	if response.StatusCode >= 400 {
		return "", errors.New(fmt.Sprintf("Request for blocks resulted in status code %d", response.StatusCode))
	}
	logger.Println("Headers are:")
	for _, h := range response.Header {
		logger.Printf("  %s\n", h)
	}
	jsonBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	logger.Printf("Read %d response bytes\n", len(jsonBody))
	logger.Printf("Response as string is:\n%s\n", string(jsonBody))
	parsedResponse := &sawtoothBlocksResponse{}
	err = json.Unmarshal(jsonBody, parsedResponse)
	if err != nil {
		return "", err
	}
	if len(parsedResponse.Data) == 0 {
		return "", errors.New("Request for blocks gave positive response, but no blocks were in")
	}
	return parsedResponse.Data[0].Header_signature, nil
}

func getIp() string {
	return os.Getenv(envVarIp)
}

type sawtoothBlocksResponse struct {
	Head string
	Link string
	Data []*blocksData
}

type blocksData struct {
	Header_signature string
}

func GetBatchStatus(batchId string) string {
	result := strings.Builder{}
	url := fmt.Sprintf("http://%s:%s/batch_statuses?id=%s", getIp(), restPort, batchId)
	result.WriteString("Sending to URL: " + url + "\n")
	response, err := http.Get(url)
	if err != nil {
		return err.Error()
	}
	msg := fmt.Sprintf("Got status code %d\n", response.StatusCode)
	result.WriteString(msg)
	writeMeaningOfstatusCode(response.StatusCode, &result)
	if response.Body != nil {
		defer func() { _ = response.Body.Close }()
	}
	if response.StatusCode >= 400 {
		return result.String()
	}
	jsonBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		result.WriteString(err.Error() + "\n")
		return result.String()
	}
	result.WriteString(fmt.Sprintf("Read %d response bytes\n", len(jsonBody)))
	result.WriteString(fmt.Sprintf("Response as string is:\n%s\n", string(jsonBody)))
	return result.String()
}

func writeMeaningOfstatusCode(statusCode int, result *strings.Builder) {
	switch statusCode {
	case 200:
		result.WriteString("Found batch\n")
	case 400:
		result.WriteString("Bad request\n")
	case 404:
		result.WriteString("Batch not found\n")
	case 500:
		result.WriteString("Internal server error\n")
	case 503:
		result.WriteString("Service unavailable\n")
	}
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
	EventStreamStatusChannel <- &EventStreamStatus{
		StatusCode: EVENT_STREAM_STATUS_RUNNING,
		Msg:        "Received events from blockchain",
	}
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

func SendCommand(c *command.Command, outputter cli.Outputter) error {
	ip := os.Getenv(envVarIp)
	if net.ParseIP(ip) == nil {
		return errors.New(fmt.Sprintf("Environment variable %s should hold an ip address, but is %s",
			envVarIp, ip))
	}
	ipAndPort := ip + ":" + restPort
	payloadBytes, err := proto.Marshal(c.Command)
	if err != nil {
		return err
	}
	context := signing.CreateContext(c.CryptoIdentity.PrivateKey.GetAlgorithmName())
	hasher := sha512.New()
	hasher.Write(payloadBytes)
	payloadSha512 := strings.ToLower(hex.EncodeToString(hasher.Sum(nil)))
	rawTransactionHeader := &transaction_pb2.TransactionHeader{
		SignerPublicKey:  c.CryptoIdentity.PublicKey.AsHex(),
		FamilyName:       model.FamilyName,
		FamilyVersion:    model.FamilyVersion,
		Dependencies:     []string{},
		BatcherPublicKey: c.CryptoIdentity.PublicKey.AsHex(),
		Inputs:           c.InputAddresses,
		Outputs:          c.OutputAddresses,
		PayloadSha512:    payloadSha512,
	}
	transactionHeaderBytes, err := proto.Marshal(rawTransactionHeader)
	if err != nil {
		return err
	}
	signature := hex.EncodeToString(context.Sign(transactionHeaderBytes, c.CryptoIdentity.PrivateKey))
	transaction := &transaction_pb2.Transaction{
		Header:          transactionHeaderBytes,
		HeaderSignature: signature,
		Payload:         payloadBytes,
	}
	transactionSignatures := []string{transaction.HeaderSignature}
	rawBatchHeader := &batch_pb2.BatchHeader{
		SignerPublicKey: c.CryptoIdentity.PublicKey.AsHex(),
		TransactionIds:  transactionSignatures,
	}
	batchHeaderBytes, err := proto.Marshal(rawBatchHeader)
	if err != nil {
		return err
	}
	batchSignature := hex.EncodeToString(context.Sign(batchHeaderBytes, c.CryptoIdentity.PrivateKey))
	batch := &batch_pb2.Batch{
		Header:          batchHeaderBytes,
		Transactions:    []*transaction_pb2.Transaction{transaction},
		HeaderSignature: batchSignature,
	}
	rawBatchList := &batch_pb2.BatchList{
		Batches: []*batch_pb2.Batch{batch},
	}
	batchListBytes, err := proto.Marshal(rawBatchList)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://%s/batches", ipAndPort)
	outputter(fmt.Sprintf("Sending command to %s\n", url))
	response, err := http.Post(
		url,
		"application/octet-stream",
		bytes.NewBuffer(batchListBytes))
	if err != nil {
		return err
	}
	outputter(fmt.Sprintf("Response status code: %d\n", response.StatusCode))
	batchSeq := TheBatchSequenceNumbers.add(batchSignature)
	outputter(fmt.Sprintf("You can request the status of batch %s using reference number %d\n",
		batchSignature, batchSeq))
	return nil
}

type BatchSequenceNumbers struct {
	nextSeq int32
	data    map[int32]string
}

func (b *BatchSequenceNumbers) add(batchId string) int32 {
	theSeq := b.nextSeq
	b.data[theSeq] = batchId
	b.nextSeq++
	return theSeq
}

func (b *BatchSequenceNumbers) Get(seq int32) (string, error) {
	batchId, found := b.data[seq]
	if !found {
		return "", errors.New(fmt.Sprintf("No batch for batch sequence number: %d", seq))
	}
	return batchId, nil
}

var TheBatchSequenceNumbers = &BatchSequenceNumbers{
	nextSeq: 0,
	data:    make(map[int32]string),
}
