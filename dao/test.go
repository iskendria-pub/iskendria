package dao

import (
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"gitlab.bbinfra.net/3estack/alexandria/model"
)

func StartFakeBlock(currentBlockId, previousBlockId string) error {
	return HandleEvent(&events_pb2.Event{
		EventType: model.EV_SAWTOOTH_BLOCK_COMMIT,
		Attributes: []*events_pb2.Event_Attribute{
			{
				Key:   model.SAWTOOTH_PREVIOUS_BLOCK_ID,
				Value: previousBlockId,
			},
			{
				Key:   model.SAWTOOTH_CURRENT_BLOCK_ID,
				Value: currentBlockId,
			},
		},
	})
}
