package dao

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/jmoiron/sqlx"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"strconv"
	"strings"
)

const THE_SETTINGS_ID = 1

type event interface {
	accept(context) error
}

type context interface {
	visitSawtoothBlockCommit(currentBlockId, previousBlockId string) error
	visitTransactionControl(transactionId string, eventSeq, numEvents int32) error
	visitDataManipulation(transactionId string, eventSeq int32, dataManipulation dataManipulation) error
}

type dataManipulation interface {
	apply(*sqlx.Tx) error
}

type dataManipulationEvent struct {
	transactionId    string
	eventSeq         int32
	dataManipulation dataManipulation
}

var _ event = new(dataManipulationEvent)

func (dme *dataManipulationEvent) accept(c context) error {
	return c.visitDataManipulation(dme.transactionId, dme.eventSeq, dme.dataManipulation)
}

func createSawtoothBlockCommitEvent(ev *events_pb2.Event) (event, error) {
	result := &sawtoothBlockCommitEvent{}
	for _, attribute := range ev.Attributes {
		switch attribute.Key {
		case model.SAWTOOTH_CURRENT_BLOCK_ID:
			result.currentBlock = attribute.Value
		case model.SAWTOOTH_PREVIOUS_BLOCK_ID:
			result.previousBlock = attribute.Value
		}
	}
	// Previous block id may be empty string for the first block.
	if result.currentBlock == "" {
		return nil, errors.New(fmt.Sprintf("Invalid block control event: %v", ev))
	}
	return result, nil
}

type sawtoothBlockCommitEvent struct {
	currentBlock  string
	previousBlock string
}

var _ event = new(sawtoothBlockCommitEvent)

func (sbe *sawtoothBlockCommitEvent) accept(c context) error {
	return c.visitSawtoothBlockCommit(sbe.currentBlock, sbe.previousBlock)
}

func createTransactionControlEvent(ev *events_pb2.Event) (event, error) {
	result := &transactionControlEvent{}
	var err error
	var i64 int64
	for _, attribute := range ev.Attributes {
		switch attribute.Key {
		case model.EV_KEY_TRANSACTION_ID:
			result.transactionId = attribute.Value
		case model.EV_KEY_EVENT_SEQ:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			result.eventSeq = int32(i64)
		case model.EV_KEY_NUM_EVENTS:
			i64, err = strconv.ParseInt(attribute.Value, 10, 32)
			result.numEvents = int32(i64)
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type transactionControlEvent struct {
	transactionId string
	eventSeq      int32
	numEvents     int32
}

var _ event = new(transactionControlEvent)

func (tce *transactionControlEvent) accept(c context) error {
	return c.visitTransactionControl(tce.transactionId, tce.eventSeq, tce.numEvents)
}

type contextImpl struct {
	stateNoBlock        context
	stateNoTransaction  context
	stateInTransaction  context
	stateKnownNumEvents context
	state               context
	blockId             string
	transactionId       string
	transaction         *sqlx.Tx
	seenEvents          map[int32]bool
	expectedNumEvents   int32
}

func createContext() {
	contextImpl := &contextImpl{}
	contextImpl.stateNoBlock = &contextStateNoBlock{parent: contextImpl}
	contextImpl.stateNoTransaction = &contextStateNoTransaction{parent: contextImpl}
	contextImpl.stateInTransaction = &contextStateInTransaction{parent: contextImpl}
	contextImpl.stateKnownNumEvents = &contextStateKnownNumEvents{parent: contextImpl}
	contextImpl.state = contextImpl.stateNoBlock
	theContext = contextImpl
}

func (c *contextImpl) visitSawtoothBlockCommit(currentBlockId, previousBlockId string) error {
	return c.state.visitSawtoothBlockCommit(currentBlockId, previousBlockId)
}

func (c *contextImpl) visitTransactionControl(transactionId string, eventSeq, numEvents int32) error {
	return c.state.visitTransactionControl(transactionId, eventSeq, numEvents)
}

func (c *contextImpl) visitDataManipulation(
	transactionId string, eventSeq int32, dataManipulation dataManipulation) error {
	return c.state.visitDataManipulation(transactionId, eventSeq, dataManipulation)
}

func (c *contextImpl) changeStateNewBlock(blockId string) {
	c.blockId = blockId
	c.changeStateNoTransaction()
}

func (c *contextImpl) changeStateNoTransaction() {
	c.transactionId = ""
	if c.transaction != nil {
		err := c.transaction.Commit()
		if err != nil {
			panic("Could not commit transaction: " + err.Error())
		}
		c.transaction = nil
	}
	c.seenEvents = make(map[int32]bool)
	c.expectedNumEvents = 0
	c.state = c.stateNoTransaction
}

func (c *contextImpl) changeStateInTransaction(transactionId string) {
	c.transactionId = transactionId
	var err error
	c.transaction, err = db.Beginx()
	if err != nil {
		panic("Could not start transaction: " + err.Error())
	}
	c.seenEvents = make(map[int32]bool)
	c.expectedNumEvents = 0
	c.state = c.stateInTransaction
}

func (c *contextImpl) changeStateKnownNumEvents(numEvents int32) {
	c.expectedNumEvents = numEvents
	c.state = c.stateKnownNumEvents
}

func (c *contextImpl) checkAndApply(
	transactionId string, eventSeq int32, dm dataManipulation) error {
	if err := c.check(transactionId, eventSeq); err != nil {
		return err
	}
	if err := dm.apply(c.transaction); err != nil {
		return errors.New(fmt.Sprintf("Error executing event for transaction id %s, error %s, eventSeq = %d",
			err, transactionId, eventSeq))
	}
	c.seenEvents[eventSeq] = true
	return nil
}

func (c *contextImpl) check(transactionId string, eventSeq int32) error {
	if transactionId != c.transactionId {
		return errors.New(fmt.Sprintf("Unexpected transaction id change, from: %s, to: %s",
			c.transactionId, transactionId))
	}
	_, alreadySeen := c.seenEvents[eventSeq]
	if alreadySeen {
		return errors.New(fmt.Sprintf("Event already seen for transaction id %s, eventSeq = %d",
			transactionId, eventSeq))
	}
	return nil
}

func (c *contextImpl) allSeen() bool {
	for seq := int32(0); seq < c.expectedNumEvents; seq++ {
		_, seen := c.seenEvents[seq]
		if !seen {
			return false
		}
	}
	return true
}

type contextStateNoBlock struct {
	parent *contextImpl
}

func (csnb *contextStateNoBlock) visitSawtoothBlockCommit(currentBlockId, previousBlockId string) error {
	csnb.parent.changeStateNewBlock(currentBlockId)
	return nil
}

func (csnb *contextStateNoBlock) visitTransactionControl(transactionId string, eventSeq, numEvents int32) error {
	return errors.New("Expect a block control before a transaction control")
}

func (csnb *contextStateNoBlock) visitDataManipulation(
	transactionId string, eventSeq int32, dataManipulation dataManipulation) error {
	return errors.New("Expect a block control event before a data manipulation")
}

type contextStateNoTransaction struct {
	parent *contextImpl
}

func (csnt *contextStateNoTransaction) visitSawtoothBlockCommit(currentBlockId, previousBlockId string) error {
	if csnt.parent.blockId != previousBlockId {
		return errors.New(fmt.Sprintf("Fork detected, expected previous %s but was %s, new = %s",
			csnt.parent.blockId, previousBlockId, currentBlockId))
	}
	csnt.parent.changeStateNewBlock(currentBlockId)
	return nil
}

func (csnt *contextStateNoTransaction) visitTransactionControl(transactionId string, eventSeq, numEvents int32) error {
	csnt.parent.changeStateInTransaction(transactionId)
	return csnt.parent.visitTransactionControl(transactionId, eventSeq, numEvents)
}

func (csnt *contextStateNoTransaction) visitDataManipulation(
	transactionId string, eventSeq int32, dataManipulation dataManipulation) error {
	csnt.parent.changeStateInTransaction(transactionId)
	return csnt.parent.visitDataManipulation(transactionId, eventSeq, dataManipulation)
}

type contextStateInTransaction struct {
	parent *contextImpl
}

func (csit *contextStateInTransaction) visitSawtoothBlockCommit(currentBlockId, previousBlockId string) error {
	return errors.New(fmt.Sprintf(
		"Changing block while in transaction without known event count; transactionId: %s, previousBlock: %s, currentBlock: %s",
		csit.parent.transactionId, previousBlockId, currentBlockId))
}

func (csit *contextStateInTransaction) visitTransactionControl(transactionId string, eventSeq, numEvents int32) error {
	if err := csit.parent.check(transactionId, eventSeq); err != nil {
		return err
	}
	csit.parent.seenEvents[eventSeq] = true
	csit.parent.changeStateKnownNumEvents(numEvents)
	if csit.parent.allSeen() {
		csit.parent.changeStateNoTransaction()
	}
	return nil
}

func (csit *contextStateInTransaction) visitDataManipulation(
	transactionId string, eventSeq int32, dm dataManipulation) error {
	return csit.parent.checkAndApply(transactionId, eventSeq, dm)
}

type contextStateKnownNumEvents struct {
	parent *contextImpl
}

func (cskne *contextStateKnownNumEvents) visitSawtoothBlockCommit(currentBlockId, previousBlockId string) error {
	return errors.New(fmt.Sprintf(
		"Changing block while in transaction with known event count; transactionId: %s, previousBlock: %s, currentBlock: %s",
		cskne.parent.transactionId, previousBlockId, currentBlockId))
}

func (cskne *contextStateKnownNumEvents) visitTransactionControl(
	transactionId string, eventSeq, numEvents int32) error {
	return errors.New(fmt.Sprintf("Duplicate control event for transaction id %s, numEvents = %d",
		transactionId, numEvents))
}

func (cskne *contextStateKnownNumEvents) visitDataManipulation(
	transactionId string, eventSeq int32, dm dataManipulation) error {
	if err := cskne.parent.checkAndApply(transactionId, eventSeq, dm); err != nil {
		return err
	}
	if cskne.parent.allSeen() {
		cskne.parent.changeStateNoTransaction()
	}
	return nil
}

var theContext context = new(contextImpl)

var db *sqlx.DB

func GetPlaceHolders(n int) string {
	placeHolders := make([]string, n)
	for i := 0; i < n; i++ {
		placeHolders[i] = "?"
	}
	return strings.Join(placeHolders, ", ")
}
