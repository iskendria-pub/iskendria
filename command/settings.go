package command

import (
	"errors"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"gitlab.bbinfra.net/3estack/alexandria/model"
)

type Bootstrap struct {
	PriceMajorEditSettings               int32
	PriceMajorCreatePerson               int32
	PriceMajorChangePersonAuthorization  int32
	PriceMajorChangeJournalAuthorization int32
	PricePersonEdit                      int32
	PriceAuthorSubmitNewManuscript       int32
	PriceAuthorSubmitNewVersion          int32
	PriceAuthorAcceptAuthorship          int32
	PriceReviewerSubmit                  int32
	PriceEditorAllowManuscriptReview     int32
	PriceEditorRejectManuscript          int32
	PriceEditorPublishManuscript         int32
	PriceEditorAssignManuscript          int32
	PriceEditorCreateJournal             int32
	PriceEditorCreateVolume              int32
	PriceEditorEditJournal               int32
	PriceEditorAddColleague              int32
	PriceEditorAcceptDuty                int32
	Name                                 string
	Email                                string
}

func GetBootstrapCommand(bootstrap *Bootstrap, cryptoIdentity *CryptoIdentity) *Command {
	personId := model.CreatePersonAddress()
	return &Command{
		InputAddresses:  []string{model.GetSettingsAddress(), personId},
		OutputAddresses: []string{model.GetSettingsAddress(), personId},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Price:     int32(0),
			Signer:    personId,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_Bootstrap{
				Bootstrap: &model.CommandBootstrap{
					PriceList: &model.PriceList{
						PriceMajorEditSettings:               bootstrap.PriceMajorEditSettings,
						PriceMajorCreatePerson:               bootstrap.PriceMajorCreatePerson,
						PriceMajorChangePersonAuthorization:  bootstrap.PriceMajorChangePersonAuthorization,
						PriceMajorChangeJournalAuthorization: bootstrap.PriceMajorChangeJournalAuthorization,
						PricePersonEdit:                      bootstrap.PricePersonEdit,
						PriceAuthorSubmitNewManuscript:       bootstrap.PriceAuthorSubmitNewManuscript,
						PriceAuthorSubmitNewVersion:          bootstrap.PriceAuthorSubmitNewVersion,
						PriceAuthorAcceptAuthorship:          bootstrap.PriceAuthorAcceptAuthorship,
						PriceReviewerSubmit:                  bootstrap.PriceReviewerSubmit,
						PriceEditorAllowManuscriptReview:     bootstrap.PriceEditorAllowManuscriptReview,
						PriceEditorRejectManuscript:          bootstrap.PriceEditorRejectManuscript,
						PriceEditorPublishManuscript:         bootstrap.PriceEditorPublishManuscript,
						PriceEditorAssignManuscript:          bootstrap.PriceEditorAssignManuscript,
						PriceEditorCreateJournal:             bootstrap.PriceEditorCreateJournal,
						PriceEditorCreateVolume:              bootstrap.PriceEditorCreateVolume,
						PriceEditorEditJournal:               bootstrap.PriceEditorEditJournal,
						PriceEditorAddColleague:              bootstrap.PriceEditorAddColleague,
						PriceEditorAcceptDuty:                bootstrap.PriceEditorAcceptDuty,
					},
					FirstMajor: &model.CommandPersonCreate{
						NewPersonId: personId,
						PublicKey:   cryptoIdentity.PublicKeyStr,
						Name:        bootstrap.Name,
						Email:       bootstrap.Email,
					},
				},
			},
		},
	}
}

func GetSettingsUpdateCommand(
	orig,
	updated *dao.Settings,
	signerId string,
	cryptoIdentity *CryptoIdentity,
	price int32) *Command {
	return &Command{
		InputAddresses:  []string{model.GetSettingsAddress(), signerId},
		OutputAddresses: []string{model.GetSettingsAddress(), signerId},
		CryptoIdentity:  cryptoIdentity,
		Command: &model.Command{
			Price:     price,
			Signer:    signerId,
			Timestamp: model.GetCurrentTime(),
			Body: &model.Command_CommandSettingsUpdate{
				CommandSettingsUpdate: createModelCommandSettingsUpdate(orig, updated),
			},
		},
	}
}

type singleUpdateSettingsCreate struct {
	timestamp int64
	priceList *model.PriceList
}

var _ singleUpdate = new(singleUpdateSettingsCreate)

func (u *singleUpdateSettingsCreate) updateState(state *unmarshalledState) (writtenAddress string) {
	state.settings = &model.StateSettings{
		CreatedOn:  u.timestamp,
		ModifiedOn: u.timestamp,
		PriceList: &model.PriceList{
			PriceMajorEditSettings:               u.priceList.PriceMajorEditSettings,
			PriceMajorCreatePerson:               u.priceList.PriceMajorCreatePerson,
			PriceMajorChangePersonAuthorization:  u.priceList.PriceMajorChangePersonAuthorization,
			PriceMajorChangeJournalAuthorization: u.priceList.PriceMajorChangeJournalAuthorization,
			PricePersonEdit:                      u.priceList.PricePersonEdit,
			PriceAuthorSubmitNewManuscript:       u.priceList.PriceAuthorSubmitNewManuscript,
			PriceAuthorSubmitNewVersion:          u.priceList.PriceAuthorSubmitNewVersion,
			PriceAuthorAcceptAuthorship:          u.priceList.PriceAuthorAcceptAuthorship,
			PriceReviewerSubmit:                  u.priceList.PriceReviewerSubmit,
			PriceEditorAllowManuscriptReview:     u.priceList.PriceEditorAllowManuscriptReview,
			PriceEditorRejectManuscript:          u.priceList.PriceEditorRejectManuscript,
			PriceEditorPublishManuscript:         u.priceList.PriceEditorPublishManuscript,
			PriceEditorAssignManuscript:          u.priceList.PriceEditorAssignManuscript,
			PriceEditorCreateJournal:             u.priceList.PriceEditorCreateJournal,
			PriceEditorCreateVolume:              u.priceList.PriceEditorCreateVolume,
			PriceEditorEditJournal:               u.priceList.PriceEditorEditJournal,
			PriceEditorAddColleague:              u.priceList.PriceEditorAddColleague,
			PriceEditorAcceptDuty:                u.priceList.PriceEditorAcceptDuty,
		},
	}
	return model.GetSettingsAddress()
}

func (u *singleUpdateSettingsCreate) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	return ba.AddEvent(
		model.EV_TYPE_SETTINGS_CREATE,
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
				Key:   model.EV_KEY_PRICE_MAJOR_EDIT_SETTINGS,
				Value: fmt.Sprintf("%d", u.priceList.PriceMajorEditSettings),
			},
			{
				Key:   model.EV_KEY_PRICE_MAJOR_CREATE_PERSON,
				Value: fmt.Sprintf("%d", u.priceList.PriceMajorCreatePerson),
			},
			{
				Key:   model.EV_KEY_PRICE_MAJOR_CHANGE_PERSON_AUTHORIZATION,
				Value: fmt.Sprintf("%d", u.priceList.PriceMajorChangePersonAuthorization),
			},
			{
				Key:   model.EV_KEY_PRICE_MAJOR_CHANGE_JOURNAL_AUTHORIZATION,
				Value: fmt.Sprintf("%d", u.priceList.PriceMajorChangeJournalAuthorization),
			},
			{
				Key:   model.EV_KEY_PRICE_PERSON_EDIT,
				Value: fmt.Sprintf("%d", u.priceList.PricePersonEdit),
			},
			{
				Key:   model.EV_KEY_PRICE_AUTHOR_SUBMIT_NEW_MANUSCRIPT,
				Value: fmt.Sprintf("%d", u.priceList.PriceAuthorSubmitNewManuscript),
			},
			{
				Key:   model.EV_KEY_PRICE_AUTHOR_SUBMIT_NEW_VERSION,
				Value: fmt.Sprintf("%d", u.priceList.PriceAuthorSubmitNewVersion),
			},
			{
				Key:   model.EV_KEY_PRICE_AUTHOR_ACCEPT_AUTHORSHIP,
				Value: fmt.Sprintf("%d", u.priceList.PriceAuthorAcceptAuthorship),
			},
			{
				Key:   model.EV_KEY_PRICE_REVIEWER_SUBMIT,
				Value: fmt.Sprintf("%d", u.priceList.PriceReviewerSubmit),
			},
			{
				Key:   model.EV_KEY_PRICE_EDITOR_ALLOW_MANUSCRIPT_REVIEW,
				Value: fmt.Sprintf("%d", u.priceList.PriceEditorAllowManuscriptReview),
			},
			{
				Key:   model.EV_KEY_PRICE_EDITOR_REJECT_MANUSCRIPT,
				Value: fmt.Sprintf("%d", u.priceList.PriceEditorRejectManuscript),
			},
			{
				Key:   model.EV_KEY_PRICE_EDITOR_PUBLISH_MANUSCRIPT,
				Value: fmt.Sprintf("%d", u.priceList.PriceEditorPublishManuscript),
			},
			{
				Key:   model.EV_KEY_PRICE_EDITOR_ASSIGN_MANUSCRIPT,
				Value: fmt.Sprintf("%d", u.priceList.PriceEditorAssignManuscript),
			},
			{
				Key:   model.EV_KEY_PRICE_EDITOR_CREATE_JOURNAL,
				Value: fmt.Sprintf("%d", u.priceList.PriceEditorCreateJournal),
			},
			{
				Key:   model.EV_KEY_PRICE_EDITOR_CREATE_VOLUME,
				Value: fmt.Sprintf("%d", u.priceList.PriceEditorCreateVolume),
			},
			{
				Key:   model.EV_KEY_PRICE_EDITOR_EDIT_JOURNAL,
				Value: fmt.Sprintf("%d", u.priceList.PriceEditorEditJournal),
			},
			{
				Key:   model.EV_KEY_PRICE_EDITOR_ADD_COLLEAGUE,
				Value: fmt.Sprintf("%d", u.priceList.PriceEditorAddColleague),
			},
			{
				Key:   model.EV_KEY_PRICE_EDITOR_ACCEPT_DUTY,
				Value: fmt.Sprintf("%d", u.priceList.PriceEditorAcceptDuty),
			},
		},
		[]byte{})
}

func (nbce *nonBootstrapCommandExecution) checkSettingsUpdate(c *model.CommandSettingsUpdate) (*updater, error) {
	expectedPrice := nbce.unmarshalledState.settings.PriceList.PriceMajorEditSettings
	if nbce.price != expectedPrice {
		return nil, formatPriceError("PriceMajorEditSettings", expectedPrice)
	}
	if !nbce.unmarshalledState.persons[nbce.verifiedSignerId].IsMajor {
		return nil, errors.New("Only majors can update settings")
	}
	oldSettings := nbce.unmarshalledState.settings
	if err := checkModelCommandSettingsUpdate(c, oldSettings); err != nil {
		return nil, err
	}
	singleUpdates := createSingleUpdatesSettingsUpdate(c, oldSettings, nbce.timestamp)
	singleUpdates = nbce.addSingleUpdateSettingsModificationTimeIfNeeded(singleUpdates)
	return &updater{
		unmarshalledState: nbce.unmarshalledState,
		updates:           singleUpdates,
	}, nil
}

func (nbce *nonBootstrapCommandExecution) addSingleUpdateSettingsModificationTimeIfNeeded(
	singleUpdates []singleUpdate) []singleUpdate {
	if len(singleUpdates) >= 1 {
		singleUpdates = append(singleUpdates, &singleUpdateSettingsModificationTime{
			timestamp: nbce.timestamp,
		})
	}
	return singleUpdates
}

type singleUpdateSettingsUpdate struct {
	stateField *int32
	newValue   int32
	eventKey   string
	timestamp  int64
}

var _ singleUpdate = new(singleUpdateSettingsUpdate)

func (u *singleUpdateSettingsUpdate) updateState(_ *unmarshalledState) (writtenAddress string) {
	*u.stateField = u.newValue
	return model.GetSettingsAddress()
}

func (u *singleUpdateSettingsUpdate) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	return ba.AddEvent(model.EV_TYPE_SETTINGS_UPDATE,
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
				Key:   u.eventKey,
				Value: fmt.Sprintf("%d", u.newValue),
			},
		},
		[]byte{})
}

type singleUpdateSettingsModificationTime struct {
	timestamp int64
}

var _ singleUpdate = new(singleUpdateSettingsModificationTime)

func (u *singleUpdateSettingsModificationTime) updateState(state *unmarshalledState) (writtenAddress string) {
	state.settings.ModifiedOn = u.timestamp
	return model.GetSettingsAddress()
}

func (u *singleUpdateSettingsModificationTime) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	return ba.AddEvent(
		model.EV_TYPE_SETTINGS_MODIFICATION_TIME,
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
		},
		[]byte{})
}
