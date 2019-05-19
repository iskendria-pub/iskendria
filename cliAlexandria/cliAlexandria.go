package cliAlexandria

import (
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/blockchain"
	"gitlab.bbinfra.net/3estack/alexandria/cli"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
)

func InitEventStream(fname, tag string) {
	go blockchain.RequestEventStream(dao.HandleEvent, fname, tag)
	awaitEventStreamRunning(func(s string) { fmt.Print(s) })
}

func awaitEventStreamRunning(outputter cli.Outputter) {
	for {
		eventStreamStatus := readEventStreamStatus()
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
	outputter(LastEventStatus + "\n")
}

var LastEventStatus string

func readEventStreamStatus() *blockchain.EventStreamStatus {
	eventStatus := <-blockchain.EventStreamStatusChannel
	updateLastEventStatus(eventStatus)
	return eventStatus
}

func updateLastEventStatus(status *blockchain.EventStreamStatus) {
	switch status.StatusCode {
	case blockchain.EVENT_STREAM_STATUS_STOPPED:
		LastEventStatus = "STOPPED"
	case blockchain.EVENT_STREAM_STATUS_RUNNING:
		LastEventStatus = "RUNNING"
	case blockchain.EVENT_STREAM_STATUS_INITIALIZING:
		LastEventStatus = "INITIALIZING"
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
