package cliAlexandria

import (
	"fmt"
	"github.com/iskendria-pub/iskendria/blockchain"
	"github.com/iskendria-pub/iskendria/cli"
	"github.com/iskendria-pub/iskendria/dao"
	"sync"
)

func InitEventStream(fname, tag string) {
	go blockchain.RequestEventStream(dao.HandleEvent, fname, tag)
	awaitEventStreamRunning(func(s string) { fmt.Print(s) })
}

func awaitEventStreamRunning(outputter cli.Outputter) {
	for {
		eventStreamStatus := ReadEventStreamStatus()
		outputter(formatEventStreamMessage(eventStreamStatus.Msg))
		switch eventStreamStatus.StatusCode {
		case blockchain.EVENT_STREAM_STATUS_INITIALIZING:
			// Continue waiting
		case blockchain.EVENT_STREAM_STATUS_STOPPED:
			outputter(formatEventStreamMessage("Could not initialize event stream from blockchain, stopping."))
			return
		case blockchain.EVENT_STREAM_STATUS_RUNNING:
			return
		}
	}
}

func eventStatus(outputter cli.Outputter) {
	outputter(LastEventStatus.Get() + "\n")
}

type ThreadSafeLastEventStatus struct {
	mux             sync.Mutex
	lastEventStatus string
}

func (es *ThreadSafeLastEventStatus) set(value string) {
	es.mux.Lock()
	defer es.mux.Unlock()
	es.lastEventStatus = value
}

func (es *ThreadSafeLastEventStatus) Get() string {
	es.mux.Lock()
	defer es.mux.Unlock()
	return es.lastEventStatus
}

var LastEventStatus = &ThreadSafeLastEventStatus{
	mux: sync.Mutex{},
}

func ReadEventStreamStatus() *blockchain.EventStreamStatus {
	eventStatus := <-blockchain.EventStreamStatusChannel
	updateLastEventStatus(eventStatus)
	return eventStatus
}

func updateLastEventStatus(status *blockchain.EventStreamStatus) {
	switch status.StatusCode {
	case blockchain.EVENT_STREAM_STATUS_STOPPED:
		LastEventStatus.set("STOPPED")
	case blockchain.EVENT_STREAM_STATUS_RUNNING:
		LastEventStatus.set("RUNNING")
	case blockchain.EVENT_STREAM_STATUS_INITIALIZING:
		LastEventStatus.set("INITIALIZING")
	}
}

func readEventStreamStatusNonblocking() *blockchain.EventStreamStatus {
	select {
	case eventStreamStatus := <-blockchain.EventStreamStatusChannel:
		updateLastEventStatus(eventStreamStatus)
		return eventStreamStatus
	default:
		return nil
	}
}

func formatEventStreamMessage(msg string) string {
	return fmt.Sprintf("*** %s\n", msg)
}

func PageEventStreamMessages(outputter cli.Outputter) {
	for {
		eventStreamStatus := readEventStreamStatusNonblocking()
		if eventStreamStatus == nil {
			break
		}
		outputter(formatEventStreamMessage(eventStreamStatus.Msg))
	}
}

var CommonDiagnosticsGroup = &cli.Cli{
	FullDescription:    "Welcome to the diagnostics commands",
	OneLineDescription: "Diagnostics",
	Name:               "diagnostics",
	Handlers: []cli.Handler{
		&cli.SingleLineHandler{
			Name:     "eventStatus",
			Handler:  eventStatus,
			ArgNames: []string{},
		},
		&cli.SingleLineHandler{
			Name:     "batchStatus",
			Handler:  batchStatus,
			ArgNames: []string{"batch seq"},
		},
	},
}

func batchStatus(outputter cli.Outputter, batchSeq int32) {
	batchId, err := blockchain.TheBatchSequenceNumbers.Get(batchSeq)
	if err != nil {
		outputter(fmt.Sprintf("No batch id known for batchSeq %d\n", batchSeq))
	}
	outputter(fmt.Sprintf("batchSeq %d corresponds to batch id %s\n", batchSeq, batchId))
	result := blockchain.GetBatchStatus(batchId)
	outputter(result)
}
