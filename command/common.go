package command

type updater struct {
	unmarshalledState *unmarshalledState
	updates           []singleUpdate
}

type singleUpdate interface {
	updateState(*unmarshalledState) (writtenAddress string)
	issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error
}
