package command

// THIS FILE HAS BEEN GENERATED BY .../generate/command/unmarshalledState/unmarshalledState
// DO NOT MODIFY!

import (
	proto "github.com/golang/protobuf/proto"
	"gitlab.bbinfra.net/3estack/alexandria/model"
)

type unmarshalledState struct {
	emptyAddresses    map[string]bool
	settings          *model.StateSettings
	persons           map[string]*model.StatePerson
	journals          map[string]*model.StateJournal
	volumes           map[string]*model.StateVolume
	manuscripts       map[string]*model.StateManuscript
	manuscriptThreads map[string]*model.StateManuscriptThread
	reviews           map[string]*model.StateReview
}

func newUnmarshalledState() *unmarshalledState {
	return &unmarshalledState{
		emptyAddresses:    make(map[string]bool),
		settings:          nil,
		persons:           make(map[string]*model.StatePerson),
		journals:          make(map[string]*model.StateJournal),
		volumes:           make(map[string]*model.StateVolume),
		manuscripts:       make(map[string]*model.StateManuscript),
		manuscriptThreads: make(map[string]*model.StateManuscriptThread),
		reviews:           make(map[string]*model.StateReview),
	}
}

func (us *unmarshalledState) getAddressState(address string) addressState {
	_, isEmpty := us.emptyAddresses[address]
	if isEmpty {
		return ADDRESS_EMPTY
	}
	switch {
	case address == model.GetSettingsAddress():
		if us.settings != nil {
			return ADDRESS_FILLED
		}
	case model.IsPersonAddress(address):
		_, found := us.persons[address]
		if found {
			return ADDRESS_FILLED
		}
	case model.IsJournalAddress(address):
		_, found := us.journals[address]
		if found {
			return ADDRESS_FILLED
		}
	case model.IsVolumeAddress(address):
		_, found := us.volumes[address]
		if found {
			return ADDRESS_FILLED
		}
	case model.IsManuscriptAddress(address):
		_, found := us.manuscripts[address]
		if found {
			return ADDRESS_FILLED
		}
	case model.IsManuscriptThreadAddress(address):
		_, found := us.manuscriptThreads[address]
		if found {
			return ADDRESS_FILLED
		}
	case model.IsReviewAddress(address):
		_, found := us.reviews[address]
		if found {
			return ADDRESS_FILLED
		}
	}
	return ADDRESS_UNKNOWN
}

func (us *unmarshalledState) add(readData map[string][]byte, requestedAddresses []string) error {
	for _, ra := range requestedAddresses {
		var err error
		contents, isAvailable := readData[ra]
		if isAvailable {
			err = us.addAvailable(ra, contents)
		} else {
			us.emptyAddresses[ra] = true
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (us *unmarshalledState) addAvailable(address string, contents []byte) error {
	var err error
	switch {
	case address == model.GetSettingsAddress():
		err = us.addSettings(contents)
	case model.IsPersonAddress(address):
		err = us.addPerson(address, contents)
	case model.IsJournalAddress(address):
		err = us.addJournal(address, contents)
	case model.IsVolumeAddress(address):
		err = us.addVolume(address, contents)
	case model.IsManuscriptAddress(address):
		err = us.addManuscript(address, contents)
	case model.IsManuscriptThreadAddress(address):
		err = us.addManuscriptThread(address, contents)
	case model.IsReviewAddress(address):
		err = us.addReview(address, contents)
	}
	return err
}

func (us *unmarshalledState) addSettings(contents []byte) error {
	settings := &model.StateSettings{}
	err := proto.Unmarshal(contents, settings)
	if err != nil {
		return err
	}
	us.settings = settings
	return nil
}

func (us *unmarshalledState) addPerson(theId string, contents []byte) error {
	modelContainer := &model.StatePerson{}
	err := proto.Unmarshal(contents, modelContainer)
	if err != nil {
		return err
	}
	us.persons[theId] = modelContainer
	return nil
}
func (us *unmarshalledState) addJournal(theId string, contents []byte) error {
	modelContainer := &model.StateJournal{}
	err := proto.Unmarshal(contents, modelContainer)
	if err != nil {
		return err
	}
	us.journals[theId] = modelContainer
	return nil
}
func (us *unmarshalledState) addVolume(theId string, contents []byte) error {
	modelContainer := &model.StateVolume{}
	err := proto.Unmarshal(contents, modelContainer)
	if err != nil {
		return err
	}
	us.volumes[theId] = modelContainer
	return nil
}
func (us *unmarshalledState) addManuscript(theId string, contents []byte) error {
	modelContainer := &model.StateManuscript{}
	err := proto.Unmarshal(contents, modelContainer)
	if err != nil {
		return err
	}
	us.manuscripts[theId] = modelContainer
	return nil
}
func (us *unmarshalledState) addManuscriptThread(theId string, contents []byte) error {
	modelContainer := &model.StateManuscriptThread{}
	err := proto.Unmarshal(contents, modelContainer)
	if err != nil {
		return err
	}
	us.manuscriptThreads[theId] = modelContainer
	return nil
}
func (us *unmarshalledState) addReview(theId string, contents []byte) error {
	modelContainer := &model.StateReview{}
	err := proto.Unmarshal(contents, modelContainer)
	if err != nil {
		return err
	}
	us.reviews[theId] = modelContainer
	return nil
}
func (us *unmarshalledState) read(addresses []string) (map[string][]byte, error) {
	result := make(map[string][]byte)
	var err error
	for _, address := range addresses {
		switch {
		case address == model.GetSettingsAddress():
			err = us.readSettings(result)
		case model.IsPersonAddress(address):
			err = us.readPerson(address, result)
		case model.IsJournalAddress(address):
			err = us.readJournal(address, result)
		case model.IsVolumeAddress(address):
			err = us.readVolume(address, result)
		case model.IsManuscriptAddress(address):
			err = us.readManuscript(address, result)
		case model.IsManuscriptThreadAddress(address):
			err = us.readManuscriptThread(address, result)
		case model.IsReviewAddress(address):
			err = us.readReview(address, result)
		}
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func (us *unmarshalledState) readSettings(result map[string][]byte) error {
	marshalled, err := proto.Marshal(us.settings)
	if err != nil {
		return err
	}
	result[model.GetSettingsAddress()] = marshalled
	return nil
}

func (us *unmarshalledState) readPerson(theId string, result map[string][]byte) error {
	marshalled, err := proto.Marshal(us.persons[theId])
	if err != nil {
		return err
	}
	result[theId] = marshalled
	return nil
}
func (us *unmarshalledState) readJournal(theId string, result map[string][]byte) error {
	marshalled, err := proto.Marshal(us.journals[theId])
	if err != nil {
		return err
	}
	result[theId] = marshalled
	return nil
}
func (us *unmarshalledState) readVolume(theId string, result map[string][]byte) error {
	marshalled, err := proto.Marshal(us.volumes[theId])
	if err != nil {
		return err
	}
	result[theId] = marshalled
	return nil
}
func (us *unmarshalledState) readManuscript(theId string, result map[string][]byte) error {
	marshalled, err := proto.Marshal(us.manuscripts[theId])
	if err != nil {
		return err
	}
	result[theId] = marshalled
	return nil
}
func (us *unmarshalledState) readManuscriptThread(theId string, result map[string][]byte) error {
	marshalled, err := proto.Marshal(us.manuscriptThreads[theId])
	if err != nil {
		return err
	}
	result[theId] = marshalled
	return nil
}
func (us *unmarshalledState) readReview(theId string, result map[string][]byte) error {
	marshalled, err := proto.Marshal(us.reviews[theId])
	if err != nil {
		return err
	}
	result[theId] = marshalled
	return nil
}
