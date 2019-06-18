package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

var templateUpdateModificationTime string = strings.TrimSpace(`
package command

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"log"
)

{{range .}}
func (nbce *nonBootstrapCommandExecution) addSingleUpdate{{.Tag}}ModificationTimeIfNeeded(
	singleUpdates []singleUpdate, subjectId string) []singleUpdate {
	if len(singleUpdates) >= 1 {
		singleUpdates = append(singleUpdates, &singleUpdate{{.Tag}}ModificationTime{
			timestamp: nbce.timestamp,
			id:  subjectId,
		})
	}
	return singleUpdates
}

type singleUpdate{{.Tag}}ModificationTime struct {
	timestamp int64
	id  string
}

var _ singleUpdate = new(singleUpdate{{.Tag}}ModificationTime)

func (u *singleUpdate{{.Tag}}ModificationTime) updateState(state *unmarshalledState) (writtenAddresses []string) {
	state.{{.StateContainer}}[u.id].ModifiedOn = u.timestamp
	return []string{u.id}
}

func (u *singleUpdate{{.Tag}}ModificationTime) issueEvent(eventSeq int32, transactionId string, ba BlockchainAccess) error {
	eventType := model.AlexandriaPrefix + {{.EventType}}
	log.Println("Sending event of type: " + eventType)
	return ba.AddEvent(eventType,
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
			{
				Key:   model.EV_KEY_ID,
				Value: u.id,
			},
		},
		[]byte{})
}
{{end}}`)

type ModificationTimeUpdate struct {
	Tag            string
	EventType      string
	StateContainer string
}

func main() {
	c := []*ModificationTimeUpdate{
		{
			Tag:            "Person",
			EventType:      "model.EV_TYPE_PERSON_MODIFICATION_TIME",
			StateContainer: "persons",
		},
		{
			Tag:            "Journal",
			EventType:      "model.EV_TYPE_JOURNAL_MODIFICATION_TIME",
			StateContainer: "journals",
		},
		{
			Tag:            "Manuscript",
			EventType:      "model.EV_TYPE_MANUSCRIPT_MODIFICATION_TIME",
			StateContainer: "manuscripts",
		},
	}
	tmpl, err := template.New("templateUpdateModificationTime").Parse(templateUpdateModificationTime)
	if err != nil {
		fmt.Println("Error parsing template")
		fmt.Println(err)
		return
	}
	err = tmpl.Execute(os.Stdout, c)
	if err != nil {
		fmt.Println(err)
	}
}
