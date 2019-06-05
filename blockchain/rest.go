package blockchain

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"gitlab.bbinfra.net/3estack/alexandria/cli"
	"gitlab.bbinfra.net/3estack/alexandria/command"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

const restPort = "8008"

func getFirstBlock(logger *log.Logger) (string, error) {
	url := fmt.Sprintf("http://%s:%s/blocks?reverse", getIp(), restPort)
	logger.Printf("Reading %s...\n", url)
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func() { _ = response.Body.Close }()
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
	payloadSha512 := model.HashBytes(payloadBytes)
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
