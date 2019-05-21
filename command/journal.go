package command

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"sort"
)

func GetJournalCreateCommand(
	jc *Journal,
	signer string,
	cryptoIdentity *CryptoIdentity,
	price int32) (*Command, string) {
	journalId := model.CreateJournalAddress()
	return &Command{
		InputAddresses:  []string{journalId, signer, model.GetSettingsAddress()},
		OutputAddresses: []string{journalId, signer},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signer,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandJournalCreate{
				CommandJournalCreate: &model.CommandJournalCreate{
					JournalId:       journalId,
					Title:           jc.Title,
					DescriptionHash: jc.DescriptionHash,
				},
			},
		},
	}, journalId
}

type Journal struct {
	Title           string
	DescriptionHash string
}

func GetCommandJournalUpdateProperties(
	journalId string,
	orig,
	updated *Journal,
	signer string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	return &Command{
		InputAddresses:  []string{journalId, signer, model.GetSettingsAddress()},
		OutputAddresses: []string{journalId, signer},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signer,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandJournalUpdateProperties{
				CommandJournalUpdateProperties: getModelCommandJournalUpdateProperties(
					journalId, orig, updated),
			},
		},
	}
}

func getModelCommandJournalUpdateProperties(journalId string, orig, updated *Journal) *model.CommandJournalUpdateProperties {
	result := &model.CommandJournalUpdateProperties{}
	result.JournalId = journalId
	if orig.Title != updated.Title {
		theUpdate := &model.StringUpdate{
			OldValue: orig.Title,
			NewValue: updated.Title,
		}
		result.TitleUpdate = theUpdate
	}
	if orig.DescriptionHash != updated.DescriptionHash {
		theUpdate := &model.StringUpdate{
			OldValue: orig.DescriptionHash,
			NewValue: updated.DescriptionHash,
		}
		result.DescriptionHashUpdate = theUpdate
	}
	return result
}

func GetCommandJournalUpdateAuthorization(
	journalId string,
	makeSigned bool,
	signer string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	return &Command{
		InputAddresses:  []string{journalId, signer, model.GetSettingsAddress()},
		OutputAddresses: []string{journalId, signer},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signer,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandJournalUpdateAuthorization{
				CommandJournalUpdateAuthorization: &model.CommandJournalUpdateAuthorization{
					JournalId:  journalId,
					MakeSigned: makeSigned,
				},
			},
		},
	}
}

func GetCommandJournalEditorResign(
	journalId string,
	signer string,
	cryptoIdentity *CryptoIdentity) *Command {
	return &Command{
		InputAddresses:  []string{journalId, signer, model.GetSettingsAddress()},
		OutputAddresses: []string{journalId, signer},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signer,
			Price:     int32(0),
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandJournalEditorResign{
				CommandJournalEditorResign: &model.CommandJournalEditorResign{
					JournalId: journalId,
				},
			},
		},
	}
}

func GetCommandJournalEditorInvite(
	journalId,
	editorId,
	signer string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	return &Command{
		InputAddresses:  []string{journalId, signer, model.GetSettingsAddress()},
		OutputAddresses: []string{journalId, signer},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signer,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandJournalEditorInvite{
				CommandJournalEditorInvite: &model.CommandJournalEditorInvite{
					JournalId:       journalId,
					InvitedEditorId: editorId,
				},
			},
		},
	}
}

func GetCommandJournalEditorAcceptDuty(
	journalId string,
	signer string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	return &Command{
		InputAddresses:  []string{journalId, signer, model.GetSettingsAddress()},
		OutputAddresses: []string{journalId, signer},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signer,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandJournalEditorAcceptDuty{
				CommandJournalEditorAcceptDuty: &model.CommandJournalEditorAcceptDuty{
					JournalId: journalId,
				},
			},
		},
	}
}

func (nbce *nonBootstrapCommandExecution) checkJournalCreate(c *model.CommandJournalCreate) (*updater, error) {
	if !model.IsJournalAddress(c.JournalId) {
		return nil, errors.New("Journal id is not a journal address: " + c.JournalId)
	}
	data, err := nbce.blockchainAccess.GetState([]string{c.JournalId})
	if err != nil {
		return nil, errors.New("Could not read journal address: " + c.JournalId)
	}
	err = nbce.unmarshalledState.add(data, []string{c.JournalId})
	if err != nil {
		return nil, err
	}
	if nbce.unmarshalledState.getAddressState(c.JournalId) != ADDRESS_EMPTY {
		return nil, errors.New("Journal already exists: " + c.JournalId)
	}
	if c.Title == "" {
		return nil, errors.New("When creating a journal, the title is mandatory")
	}
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates: []singleUpdate{
			&singleUpdateJournalCreate{
				timestamp:     nbce.timestamp,
				journalCreate: c,
			},
			&singleUpdateEditorCreate{
				journalId: c.JournalId,
				editorId:  nbce.verifiedSignerId,
				timestamp: nbce.timestamp,
			},
		},
	}, nil
}

type singleUpdateJournalCreate struct {
	timestamp     int64
	journalCreate *model.CommandJournalCreate
}

var _ singleUpdate = new(singleUpdateJournalCreate)

func (u *singleUpdateJournalCreate) updateState(state *unmarshalledState) string {
	journalId := u.journalCreate.JournalId
	journal := &model.StateJournal{
		Id:              journalId,
		CreatedOn:       u.timestamp,
		ModifiedOn:      u.timestamp,
		Title:           u.journalCreate.Title,
		IsSigned:        false,
		DescriptionHash: u.journalCreate.DescriptionHash,
		EditorInfo:      []*model.EditorInfo{},
	}
	state.journals[journalId] = journal
	return journalId
}

func (u *singleUpdateJournalCreate) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	eventType := model.AlexandriaPrefix + model.EV_TYPE_JOURNAL_CREATE
	return ba.AddEvent(
		eventType,
		[]processor.Attribute{
			{
				Key:   model.EV_KEY_TIMESTAMP,
				Value: fmt.Sprintf("%d", u.timestamp),
			},
			{
				Key:   model.EV_KEY_TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.EV_KEY_EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.EV_KEY_JOURNAL_ID,
				Value: u.journalCreate.JournalId,
			},
			{
				Key:   model.EV_KEY_DESCRIPTION_HASH,
				Value: u.journalCreate.DescriptionHash,
			},
			{
				Key:   model.EV_KEY_TITLE,
				Value: u.journalCreate.Title,
			},
		}, []byte{})
}

type singleUpdateEditorCreate struct {
	journalId string
	editorId  string
	timestamp int64
}

var _ singleUpdate = new(singleUpdateEditorCreate)

func (u *singleUpdateEditorCreate) updateState(state *unmarshalledState) string {
	journal := state.journals[u.journalId]
	journal.EditorInfo = append(journal.EditorInfo, &model.EditorInfo{
		EditorId:    u.editorId,
		EditorState: model.EditorState_editorAccepted,
	})
	sort.Slice(journal.EditorInfo, func(i, j int) bool {
		return journal.EditorInfo[i].EditorId < journal.EditorInfo[j].EditorId
	})
	return u.journalId
}

func (u *singleUpdateEditorCreate) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	eventType := model.AlexandriaPrefix + model.EV_TYPE_EDITOR_CREATE
	return ba.AddEvent(
		eventType,
		[]processor.Attribute{
			{
				Key:   model.EV_KEY_TIMESTAMP,
				Value: fmt.Sprintf("%d", u.timestamp),
			},
			{
				Key:   model.EV_KEY_TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.EV_KEY_EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.EV_KEY_JOURNAL_ID,
				Value: u.journalId,
			},
			{
				Key:   model.EV_KEY_PERSON_ID,
				Value: u.editorId,
			},
			{
				Key:   model.EV_KEY_EDITOR_STATE,
				Value: model.GetEditorStateString(model.EditorState_editorAccepted),
			},
		}, []byte{})
}

func (nbce *nonBootstrapCommandExecution) checkJournalUpdateProperties(c *model.CommandJournalUpdateProperties) (
	*updater, error) {
	return nil, nil
}

func (nbce *nonBootstrapCommandExecution) checkJournalUpdateAuthorization(c *model.CommandJournalUpdateAuthorization) (
	*updater, error) {
	return nil, nil
}

func (nbce *nonBootstrapCommandExecution) checkJournalEditorResign(c *model.CommandJournalEditorResign) (
	*updater, error) {
	return nil, nil
}

func (nbce *nonBootstrapCommandExecution) checkJournalEditorInvite(c *model.CommandJournalEditorInvite) (
	*updater, error) {
	return nil, nil
}

func (nbce *nonBootstrapCommandExecution) checkJournalEditorAcceptDuty(c *model.CommandJournalEditorAcceptDuty) (
	*updater, error) {
	return nil, nil
}
