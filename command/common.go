package command

import (
	"errors"
	"fmt"
)

type updater struct {
	unmarshalledState *unmarshalledState
	updates           []singleUpdate
}

type singleUpdate interface {
	updateState(*unmarshalledState) (writtenAddresses []string)
	issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error
}

func formatPriceError(priceName string, expectedPrice int32) error {
	return errors.New(fmt.Sprintf("Price should be %s, which is %d",
		priceName, expectedPrice))
}
