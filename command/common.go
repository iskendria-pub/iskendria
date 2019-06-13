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

func (nbce *nonBootstrapCommandExecution) readAndCheckAddresses(
	addressesExpectedFilled []string,
	addressesExpectedEmpty []string) error {
	toRead := append(addressesExpectedFilled, addressesExpectedEmpty...)
	readState, err := nbce.blockchainAccess.GetState(toRead)
	if err != nil {
		return err
	}
	err = nbce.unmarshalledState.add(readState, toRead)
	if err != nil {
		return err
	}
	for _, a := range addressesExpectedFilled {
		if nbce.unmarshalledState.getAddressState(a) != ADDRESS_FILLED {
			return errors.New("Address was not filled: " + a)
		}
	}
	for _, a := range addressesExpectedEmpty {
		if nbce.unmarshalledState.getAddressState(a) != ADDRESS_EMPTY {
			return errors.New("Manuscript id or manuscript thread id already in use: " + a)
		}
	}
	return nil
}

func (nbce *nonBootstrapCommandExecution) addAuthorUpdates(
	authorIds []string,
	formalUpdates []singleUpdate,
	manuscriptId string) []singleUpdate {
	for i, a := range authorIds {
		didSign := (a == nbce.verifiedSignerId)
		formalUpdates = append(formalUpdates, &singleUpdateAuthorCreate{
			manuscriptId: manuscriptId,
			authorId:     a,
			didSign:      didSign,
			authorNumber: int32(i),
			timestamp:    nbce.timestamp,
		})
	}
	return formalUpdates
}
