package command

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"log"
)

func (nbce *nonBootstrapCommandExecution) addSingleUpdatePersonModificationTimeIfNeeded(
	singleUpdates []singleUpdate, subjectId string) []singleUpdate {
	if len(singleUpdates) >= 1 {
		singleUpdates = append(singleUpdates, &singleUpdatePersonModificationTime{
			timestamp: nbce.timestamp,
			id:        subjectId,
		})
	}
	return singleUpdates
}

type singleUpdatePersonModificationTime struct {
	timestamp int64
	id        string
}

var _ singleUpdate = new(singleUpdatePersonModificationTime)

func (u *singleUpdatePersonModificationTime) updateState(state *unmarshalledState) (writtenAddresses []string) {
	state.persons[u.id].ModifiedOn = u.timestamp
	return []string{u.id}
}

func (u *singleUpdatePersonModificationTime) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	eventType := model.AlexandriaPrefix + model.EV_TYPE_PERSON_MODIFICATION_TIME
	log.Println("Sending event of type: " + eventType)
	return ba.AddEvent(eventType,
		[]processor.Attribute{
			{
				Key:   model.EV_KEY_TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.EV_KEY_TIMESTAMP,
				Value: fmt.Sprintf("%d", u.timestamp),
			},
			{
				Key:   model.EV_KEY_EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.EV_KEY_ID,
				Value: u.id,
			},
		},
		[]byte{})
}

func (nbce *nonBootstrapCommandExecution) addSingleUpdateJournalModificationTimeIfNeeded(
	singleUpdates []singleUpdate, subjectId string) []singleUpdate {
	if len(singleUpdates) >= 1 {
		singleUpdates = append(singleUpdates, &singleUpdateJournalModificationTime{
			timestamp: nbce.timestamp,
			id:        subjectId,
		})
	}
	return singleUpdates
}

type singleUpdateJournalModificationTime struct {
	timestamp int64
	id        string
}

var _ singleUpdate = new(singleUpdateJournalModificationTime)

func (u *singleUpdateJournalModificationTime) updateState(state *unmarshalledState) (writtenAddresses []string) {
	state.journals[u.id].ModifiedOn = u.timestamp
	return []string{u.id}
}

func (u *singleUpdateJournalModificationTime) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	eventType := model.AlexandriaPrefix + model.EV_TYPE_JOURNAL_MODIFICATION_TIME
	log.Println("Sending event of type: " + eventType)
	return ba.AddEvent(eventType,
		[]processor.Attribute{
			{
				Key:   model.EV_KEY_TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.EV_KEY_TIMESTAMP,
				Value: fmt.Sprintf("%d", u.timestamp),
			},
			{
				Key:   model.EV_KEY_EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.EV_KEY_ID,
				Value: u.id,
			},
		},
		[]byte{})
}
