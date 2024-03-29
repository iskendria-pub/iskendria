package command

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/iskendria-pub/iskendria/model"
	"log"
	"sort"
	"strconv"
)

func GetCommandJournalCreate(
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
					JournalId: journalId,
					Title:     jc.Title,
				},
			},
		},
	}, journalId
}

type Journal struct {
	Title string
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
	return result
}

func GetCommandJournalUpdateDescription(
	journalId string,
	origDescriptionHash string,
	updatedDescription []byte,
	signer string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	updatedDescriptionHash := model.HashBytes(updatedDescription)
	return &Command{
		InputAddresses:  []string{journalId, signer, model.GetSettingsAddress()},
		OutputAddresses: []string{journalId, signer},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signer,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandJournalUpdateProperties{
				CommandJournalUpdateProperties: &model.CommandJournalUpdateProperties{
					JournalId: journalId,
					DescriptionHashUpdate: &model.StringUpdate{
						OldValue: origDescriptionHash,
						NewValue: updatedDescriptionHash,
					},
				},
			},
		},
	}
}

func GetCommandJournalOmitDescription(
	journalId string,
	origDescriptionHash string,
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
				CommandJournalUpdateProperties: &model.CommandJournalUpdateProperties{
					JournalId: journalId,
					DescriptionHashUpdate: &model.StringUpdate{
						OldValue: origDescriptionHash,
						NewValue: "",
					},
				},
			},
		},
	}
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

func GetCommandEditorResign(
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

func GetCommandEditorInvite(
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

func GetCommandEditorAcceptDuty(
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
	expectedPrice := nbce.unmarshalledState.settings.PriceList.PriceEditorCreateJournal
	if nbce.price != expectedPrice {
		return nil, formatPriceError("PriceEditorCreateJournal", expectedPrice)
	}
	if err := nbce.readAndCheckJournal(c.JournalId, ADDRESS_EMPTY); err != nil {
		return nil, err
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
				journalId:   c.JournalId,
				editorId:    nbce.verifiedSignerId,
				editorState: model.EditorState_editorAccepted,
				timestamp:   nbce.timestamp,
			},
		},
	}, nil
}

func (nbce *nonBootstrapCommandExecution) readAndCheckJournal(
	journalId string, expectedAddressState addressState) error {
	if !model.IsJournalAddress(journalId) {
		return errors.New("Journal id is not a journal address: " + journalId)
	}
	data, err := nbce.blockchainAccess.GetState([]string{journalId})
	if err != nil {
		return errors.New("Could not read journal address: " + journalId)
	}
	err = nbce.unmarshalledState.add(data, []string{journalId})
	if err != nil {
		return err
	}
	if nbce.unmarshalledState.getAddressState(journalId) != expectedAddressState {
		return reportUnexpectedJournalAddressState(expectedAddressState, journalId)
	}
	return nil
}

func reportUnexpectedJournalAddressState(expectedAddressState addressState, journalId string) error {
	var msg string
	switch expectedAddressState {
	case ADDRESS_EMPTY:
		msg = "Journal already exist: "
	case ADDRESS_FILLED:
		msg = "Journal does not exist: "
	case ADDRESS_UNKNOWN:
		msg = "Internal error, journal was not read: "
	}
	err := errors.New(msg + journalId)
	return err
}

type singleUpdateJournalCreate struct {
	timestamp     int64
	journalCreate *model.CommandJournalCreate
}

var _ singleUpdate = new(singleUpdateJournalCreate)

func (u *singleUpdateJournalCreate) updateState(state *unmarshalledState) []string {
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
	return []string{journalId}
}

func (u *singleUpdateJournalCreate) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	eventType := model.AlexandriaPrefix + model.EV_TYPE_JOURNAL_CREATE
	log.Println("Sending event of type: " + eventType)
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
				Key:   model.EV_KEY_ID,
				Value: u.journalCreate.JournalId,
			},
			{
				Key:   model.EV_KEY_JOURNAL_DESCRIPTION_HASH,
				Value: u.journalCreate.DescriptionHash,
			},
			{
				Key:   model.EV_KEY_JOURNAL_TITLE,
				Value: u.journalCreate.Title,
			},
		}, []byte{})
}

type singleUpdateEditorCreate struct {
	journalId   string
	editorId    string
	editorState model.EditorState
	timestamp   int64
}

var _ singleUpdate = new(singleUpdateEditorCreate)

func (u *singleUpdateEditorCreate) updateState(state *unmarshalledState) []string {
	journal := state.journals[u.journalId]
	journal.EditorInfo = append(journal.EditorInfo, &model.EditorInfo{
		EditorId:    u.editorId,
		EditorState: u.editorState,
	})
	sort.Slice(journal.EditorInfo, func(i, j int) bool {
		return journal.EditorInfo[i].EditorId < journal.EditorInfo[j].EditorId
	})
	return []string{u.journalId}
}

func (u *singleUpdateEditorCreate) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	eventType := model.AlexandriaPrefix + model.EV_TYPE_EDITOR_CREATE
	log.Println("Sending event of type: " + eventType)
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
				Key:   model.EV_KEY_EDITOR_ID,
				Value: u.editorId,
			},
			{
				Key:   model.EV_KEY_EDITOR_STATE,
				Value: model.GetEditorStateString(u.editorState),
			},
		}, []byte{})
}

func (nbce *nonBootstrapCommandExecution) checkJournalUpdateProperties(c *model.CommandJournalUpdateProperties) (
	*updater, error) {
	expectedPrice := nbce.unmarshalledState.settings.PriceList.PriceEditorEditJournal
	if nbce.price != expectedPrice {
		return nil, formatPriceError("PriceEditorEditJournal", expectedPrice)
	}
	if err := nbce.readAndCheckJournal(c.JournalId, ADDRESS_FILLED); err != nil {
		return nil, err
	}
	if !nbce.signerIsEditor(c.JournalId, []model.EditorState{model.EditorState_editorAccepted}) {
		return nil, errors.New(fmt.Sprintf(
			"You are not editor of journal %s, or you still have to accept editorship", c.JournalId))
	}
	oldJournal := nbce.unmarshalledState.journals[c.JournalId]
	if c.TitleUpdate != nil && c.TitleUpdate.OldValue != oldJournal.Title {
		return nil, errors.New(fmt.Sprintf("Title mismatch: expected %s, got %s",
			c.TitleUpdate.OldValue, oldJournal.Title))
	}
	if c.DescriptionHashUpdate != nil && c.DescriptionHashUpdate.OldValue != oldJournal.DescriptionHash {
		return nil, errors.New(fmt.Sprintf("DescriptionHash mismatch: expected %s, got %s",
			c.DescriptionHashUpdate.OldValue, oldJournal.DescriptionHash))
	}
	singleUpdates := createSingleUpdatesJournalUpdateProperties(c, oldJournal, nbce.timestamp)
	singleUpdates = nbce.addSingleUpdateJournalModificationTimeIfNeeded(singleUpdates, c.JournalId)
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates:           singleUpdates,
	}, nil
}

func createSingleUpdatesJournalUpdateProperties(
	c *model.CommandJournalUpdateProperties, oldJournal *model.StateJournal, timestamp int64) []singleUpdate {
	result := []singleUpdate{}
	if c.TitleUpdate != nil {
		result = append(result, &singleUpdateJournalUpdateProperties{
			newValue:   c.TitleUpdate.NewValue,
			stateField: &oldJournal.Title,
			eventKey:   model.EV_KEY_JOURNAL_TITLE,
			journalId:  c.JournalId,
			timestamp:  timestamp,
		})
	}
	if c.DescriptionHashUpdate != nil {
		result = append(result, &singleUpdateJournalUpdateProperties{
			newValue:   c.DescriptionHashUpdate.NewValue,
			stateField: &oldJournal.DescriptionHash,
			eventKey:   model.EV_KEY_JOURNAL_DESCRIPTION_HASH,
			journalId:  c.JournalId,
			timestamp:  timestamp,
		})
	}
	return result
}

type singleUpdateJournalUpdateProperties struct {
	newValue   string
	stateField *string
	eventKey   string
	journalId  string
	timestamp  int64
}

var _ singleUpdate = new(singleUpdateJournalUpdateProperties)

func (u *singleUpdateJournalUpdateProperties) updateState(*unmarshalledState) (writtenAddress []string) {
	*u.stateField = u.newValue
	return []string{u.journalId}
}

func (u *singleUpdateJournalUpdateProperties) issueEvent(
	eventSeq int32, transactionId string, ba BlockchainAccess) error {
	eventType := model.AlexandriaPrefix + model.EV_TYPE_JOURNAL_UPDATE
	log.Println("Sending event of type: " + eventType)
	return ba.AddEvent(
		eventType,
		[]processor.Attribute{
			{
				Key:   model.EV_KEY_TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.EV_KEY_EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.EV_KEY_TIMESTAMP,
				Value: fmt.Sprintf("%d", u.timestamp),
			},
			{
				Key:   model.EV_KEY_ID,
				Value: u.journalId,
			},
			{
				Key:   u.eventKey,
				Value: u.newValue,
			},
		}, []byte{})
}

func (nbce *nonBootstrapCommandExecution) signerIsEditor(
	journalId string, allowedEditorStates []model.EditorState) bool {
	for _, e := range nbce.unmarshalledState.journals[journalId].EditorInfo {
		for _, allowedState := range allowedEditorStates {
			if e.EditorId == nbce.verifiedSignerId && e.EditorState == allowedState {
				return true
			}
		}
	}
	return false
}

func (nbce *nonBootstrapCommandExecution) checkJournalUpdateAuthorization(c *model.CommandJournalUpdateAuthorization) (
	*updater, error) {
	expectedPrice := nbce.unmarshalledState.settings.PriceList.PriceMajorChangeJournalAuthorization
	if nbce.price != expectedPrice {
		return nil, formatPriceError("PriceMajorChangeJournalAuthorization", expectedPrice)
	}
	if err := nbce.readAndCheckJournal(c.JournalId, ADDRESS_FILLED); err != nil {
		return nil, err
	}
	if !nbce.unmarshalledState.persons[nbce.verifiedSignerId].IsMajor {
		return nil, errors.New("Only majors can change the isSigned property of a journal")
	}
	oldJournal := nbce.unmarshalledState.journals[c.JournalId]
	updates := []singleUpdate{}
	if oldJournal.IsSigned != c.MakeSigned {
		updates = append(updates, &singleUpdateJournalUpdateAuthorization{
			journalId:  c.JournalId,
			makeSigned: c.MakeSigned,
			timestamp:  nbce.timestamp,
		})
	}
	updates = nbce.addSingleUpdateJournalModificationTimeIfNeeded(updates, c.JournalId)
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates:           updates,
	}, nil
}

type singleUpdateJournalUpdateAuthorization struct {
	journalId  string
	makeSigned bool
	timestamp  int64
}

var _ singleUpdate = new(singleUpdateJournalUpdateAuthorization)

func (u *singleUpdateJournalUpdateAuthorization) updateState(state *unmarshalledState) (writtenAddresses []string) {
	state.journals[u.journalId].IsSigned = u.makeSigned
	return []string{u.journalId}
}

func (u *singleUpdateJournalUpdateAuthorization) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	eventType := model.AlexandriaPrefix + model.EV_TYPE_JOURNAL_UPDATE
	log.Println("Sending event of type: " + eventType)
	return ba.AddEvent(eventType,
		[]processor.Attribute{
			{
				Key:   model.EV_KEY_TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.EV_KEY_EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.EV_KEY_TIMESTAMP,
				Value: fmt.Sprintf("%d", u.timestamp),
			},
			{
				Key:   model.EV_KEY_ID,
				Value: u.journalId,
			},
			{
				Key:   model.EV_KEY_JOURNAL_IS_SIGNED,
				Value: strconv.FormatBool(u.makeSigned),
			},
		}, []byte{})
}

func (nbce *nonBootstrapCommandExecution) checkJournalEditorResign(c *model.CommandJournalEditorResign) (
	*updater, error) {
	expectedPrice := int32(0)
	if nbce.price != expectedPrice {
		return nil, formatPriceError("<No price>", expectedPrice)
	}
	if err := nbce.readAndCheckJournal(c.JournalId, ADDRESS_FILLED); err != nil {
		return nil, err
	}
	if !nbce.signerIsEditor(c.JournalId, []model.EditorState{
		model.EditorState_editorProposed, model.EditorState_editorAccepted}) {
		return nil, errors.New("You are not the editor of the journal")
	}
	updates := []singleUpdate{
		&singleUpdateEditorDelete{
			journalId: c.JournalId,
			editorId:  nbce.verifiedSignerId,
			timestamp: nbce.timestamp,
		},
	}
	updates = nbce.addSingleUpdateJournalModificationTimeIfNeeded(updates, c.JournalId)
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates:           updates,
	}, nil
}

type singleUpdateEditorDelete struct {
	journalId string
	editorId  string
	timestamp int64
}

var _ singleUpdate = new(singleUpdateEditorDelete)

func (u *singleUpdateEditorDelete) updateState(
	state *unmarshalledState) (writtenAddresses []string) {
	existingEditors := state.journals[u.journalId].EditorInfo
	newEditors := make([]*model.EditorInfo, 0, len(existingEditors)-1)
	for _, e := range existingEditors {
		if e.EditorId != u.editorId {
			newEditors = append(newEditors, &model.EditorInfo{
				EditorId:    e.EditorId,
				EditorState: e.EditorState,
			})
		}
	}
	state.journals[u.journalId].EditorInfo = newEditors
	return []string{u.journalId}
}

func (u *singleUpdateEditorDelete) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	eventType := model.AlexandriaPrefix + model.EV_TYPE_EDITOR_DELETE
	log.Println("Sending event of type: " + eventType)
	return ba.AddEvent(
		eventType,
		[]processor.Attribute{
			{
				Key:   model.EV_KEY_TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.EV_KEY_EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.EV_KEY_TIMESTAMP,
				Value: fmt.Sprintf("%d", u.timestamp),
			},
			{
				Key:   model.EV_KEY_JOURNAL_ID,
				Value: u.journalId,
			},
			{
				Key:   model.EV_KEY_EDITOR_ID,
				Value: u.editorId,
			},
		}, []byte{})
}

func (nbce *nonBootstrapCommandExecution) checkJournalEditorInvite(c *model.CommandJournalEditorInvite) (
	*updater, error) {
	expectedPrice := nbce.unmarshalledState.settings.PriceList.PriceEditorAddColleague
	if nbce.price != expectedPrice {
		return nil, formatPriceError("PriceEditorAddColleague", expectedPrice)
	}
	if err := nbce.readAndCheckJournal(c.JournalId, ADDRESS_FILLED); err != nil {
		return nil, err
	}
	if err := nbce.checkIsNotEditor(c.InvitedEditorId, c.JournalId); err != nil {
		return nil, err
	}
	updates := []singleUpdate{
		&singleUpdateEditorCreate{
			journalId:   c.JournalId,
			editorId:    c.InvitedEditorId,
			editorState: model.EditorState_editorProposed,
			timestamp:   nbce.timestamp,
		},
	}
	updates = nbce.addSingleUpdateJournalModificationTimeIfNeeded(updates, c.JournalId)
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates:           updates,
	}, nil
}

func (nbce *nonBootstrapCommandExecution) checkIsNotEditor(personId, journalId string) error {
	for _, e := range nbce.unmarshalledState.journals[journalId].EditorInfo {
		if e.EditorId == personId {
			return errors.New("You are editor already, or you were already proposed as editor")
		}
	}
	return nil
}

func (nbce *nonBootstrapCommandExecution) checkJournalEditorAcceptDuty(c *model.CommandJournalEditorAcceptDuty) (
	*updater, error) {
	expectedPrice := nbce.unmarshalledState.settings.PriceList.PriceEditorAcceptDuty
	if nbce.price != expectedPrice {
		return nil, formatPriceError("PriceEditorAcceptDuty", expectedPrice)
	}
	if err := nbce.readAndCheckJournal(c.JournalId, ADDRESS_FILLED); err != nil {
		return nil, err
	}

	if !nbce.signerIsEditor(c.JournalId, []model.EditorState{model.EditorState_editorProposed}) {
		return nil, errors.New("You already accepted editorship, or you are no editor at all")
	}
	updates := []singleUpdate{
		&singleUpdateEditorAcceptDuty{
			journalId: c.JournalId,
			editorId:  nbce.verifiedSignerId,
			timestamp: nbce.timestamp,
		},
	}
	updates = nbce.addSingleUpdateJournalModificationTimeIfNeeded(updates, c.JournalId)
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates:           updates,
	}, nil
}

type singleUpdateEditorAcceptDuty struct {
	journalId string
	editorId  string
	timestamp int64
}

var _ singleUpdate = new(singleUpdateEditorAcceptDuty)

func (u *singleUpdateEditorAcceptDuty) updateState(state *unmarshalledState) (writtenAddresses []string) {
	journal := state.journals[u.journalId]
	for _, e := range journal.EditorInfo {
		if e.EditorId == u.editorId {
			e.EditorState = model.EditorState_editorAccepted
		}
	}
	return []string{u.journalId}
}

func (u *singleUpdateEditorAcceptDuty) issueEvent(
	eventSeq int32, transactionId string, ba BlockchainAccess) error {
	eventType := model.AlexandriaPrefix + model.EV_TYPE_EDITOR_UPDATE
	attributes := []processor.Attribute{
		{
			Key:   model.EV_KEY_TRANSACTION_ID,
			Value: transactionId,
		},
		{
			Key:   model.EV_KEY_EVENT_SEQ,
			Value: fmt.Sprintf("%d", eventSeq),
		},
		{
			Key:   model.EV_KEY_TIMESTAMP,
			Value: fmt.Sprintf("%d", u.timestamp),
		},
		{
			Key:   model.EV_KEY_JOURNAL_ID,
			Value: u.journalId,
		},
		{
			Key:   model.EV_KEY_EDITOR_ID,
			Value: u.editorId,
		},
		{
			Key:   model.EV_KEY_EDITOR_STATE,
			Value: model.GetEditorStateString(model.EditorState_editorAccepted),
		},
	}
	log.Println("Sending event of type: " + eventType)
	return ba.AddEvent(eventType, attributes, []byte{})
}

type Volume struct {
	JournalId              string
	Issue                  string
	LogicalPublicationTime int64
}

func GetCommandVolumeCreate(
	v *Volume,
	signer string,
	cryptoIdentity *CryptoIdentity,
	price int32) (*Command, string) {
	volumeId := model.CreateVolumeAddress()
	return &Command{
		InputAddresses:  []string{model.GetSettingsAddress(), signer, volumeId, v.JournalId},
		OutputAddresses: []string{signer, volumeId},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Signer:    signer,
			Price:     price,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandVolumeCreate{
				CommandVolumeCreate: &model.CommandVolumeCreate{
					VolumeId:               volumeId,
					JournalId:              v.JournalId,
					Issue:                  v.Issue,
					LogicalPublicationTime: v.LogicalPublicationTime,
				},
			},
		},
	}, volumeId
}

func (nbce *nonBootstrapCommandExecution) checkVolumeCreate(c *model.CommandVolumeCreate) (*updater, error) {
	expectedPrice := nbce.unmarshalledState.settings.PriceList.PriceEditorCreateVolume
	if nbce.price != expectedPrice {
		return nil, formatPriceError("PriceEditorCreateVolume", expectedPrice)
	}
	if err := nbce.readAndCheckJournal(c.JournalId, ADDRESS_FILLED); err != nil {
		return nil, err
	}
	if !model.IsVolumeAddress(c.VolumeId) {
		return nil, errors.New("Volume id is not a volume address: " + c.VolumeId)
	}
	data, err := nbce.blockchainAccess.GetState([]string{c.VolumeId})
	if err != nil {
		return nil, errors.New("Could not read volume address: " + c.VolumeId)
	}
	err = nbce.unmarshalledState.add(data, []string{c.VolumeId})
	if err != nil {
		return nil, err
	}
	if nbce.unmarshalledState.getAddressState(c.VolumeId) != ADDRESS_EMPTY {
		return nil, errors.New("Volume already exists: " + c.VolumeId)
	}
	if c.Issue == "" {
		return nil, errors.New("Volume issue string should be filled but was empty")
	}
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates: []singleUpdate{
			&singleUpdateVolumeCreate{
				volumeId:               c.VolumeId,
				journalId:              c.JournalId,
				issue:                  c.Issue,
				timestamp:              nbce.timestamp,
				logicalPublicationTime: c.LogicalPublicationTime,
			},
		},
	}, nil
}

type singleUpdateVolumeCreate struct {
	volumeId               string
	journalId              string
	issue                  string
	timestamp              int64
	logicalPublicationTime int64
}

var _ singleUpdate = new(singleUpdateVolumeCreate)

func (u *singleUpdateVolumeCreate) updateState(state *unmarshalledState) (writtenAddresses []string) {
	state.volumes[u.volumeId] = &model.StateVolume{
		Id:                     u.volumeId,
		CreatedOn:              u.timestamp,
		JournalId:              u.journalId,
		Issue:                  u.issue,
		LogicalPublicationTime: u.logicalPublicationTime,
	}
	return []string{u.volumeId}
}

func (u *singleUpdateVolumeCreate) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	return ba.AddEvent(model.AlexandriaPrefix+model.EV_TYPE_VOLUME_CREATE,
		[]processor.Attribute{
			{
				Key:   model.EV_KEY_TRANSACTION_ID,
				Value: transactionId,
			},
			{
				Key:   model.EV_KEY_EVENT_SEQ,
				Value: fmt.Sprintf("%d", eventSeq),
			},
			{
				Key:   model.EV_KEY_TIMESTAMP,
				Value: fmt.Sprintf("%d", u.timestamp),
			},
			{
				Key:   model.EV_KEY_ID,
				Value: u.volumeId,
			},
			{
				Key:   model.EV_KEY_JOURNAL_ID,
				Value: u.journalId,
			},
			{
				Key:   model.EV_KEY_VOLUME_ISSUE,
				Value: u.issue,
			},
			{
				Key:   model.EV_KEY_VOLUME_LOGICAL_PUBLICATION_TIME,
				Value: fmt.Sprintf("%d", u.logicalPublicationTime),
			},
		}, []byte{})
}
