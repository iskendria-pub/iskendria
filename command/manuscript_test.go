package command

import (
	"fmt"
	"github.com/iskendria-pub/iskendria/model"
	"testing"
)

func TestGetCommandManuscriptAcceptAuthorshipWork(t *testing.T) {
	contexts := []acceptAuthorshipContext{
		{
			numAuthors:                       1,
			numSigner:                        0,
			numbersAlreadySigned:             []int{},
			isThreadReviewable:               false,
			expectedDoesAuthorUpdate:         true,
			expectedAllAuthorsWillHaveSigned: true,
			expectedNewStatus:                model.ManuscriptStatus_new,
		},
		{
			numAuthors:                       1,
			numSigner:                        0,
			numbersAlreadySigned:             []int{},
			isThreadReviewable:               true,
			expectedDoesAuthorUpdate:         true,
			expectedAllAuthorsWillHaveSigned: true,
			expectedNewStatus:                model.ManuscriptStatus_reviewable,
		},
		{
			numAuthors:                       2,
			numSigner:                        1,
			numbersAlreadySigned:             []int{0},
			isThreadReviewable:               false,
			expectedDoesAuthorUpdate:         true,
			expectedAllAuthorsWillHaveSigned: true,
			expectedNewStatus:                model.ManuscriptStatus_new,
		},
		{
			numAuthors:                       2,
			numSigner:                        1,
			numbersAlreadySigned:             []int{},
			isThreadReviewable:               false,
			expectedDoesAuthorUpdate:         true,
			expectedAllAuthorsWillHaveSigned: false,
			expectedNewStatus:                model.ManuscriptStatus_init,
		},
	}
	for _, c := range contexts {
		authorIds := make([]string, c.numAuthors)
		for i := range authorIds {
			authorIds[i] = model.CreatePersonAddress()
		}
		authors := make([]*model.Author, c.numAuthors)
		for authorNumber := range authors {
			didSign := false
			for _, authorNumberThatSigned := range c.numbersAlreadySigned {
				if authorNumber == authorNumberThatSigned {
					didSign = true
				}
			}
			authors[authorNumber] = &model.Author{
				AuthorId:     authorIds[authorNumber],
				DidSign:      didSign,
				AuthorNumber: int32(authorNumber),
			}
		}
		cmd := &model.CommandManuscriptAcceptAuthorship{
			Author: authors,
		}
		actualDoesAuthorUpdate, actualAllAuthorsWillHaveSigned := getCommandManuscriptAcceptAuthorshipWork(
			cmd, authorIds[c.numSigner])
		if actualDoesAuthorUpdate != c.expectedDoesAuthorUpdate {
			t.Error("DoesAuthorUpdate mismatch")
		}
		if actualAllAuthorsWillHaveSigned != c.expectedAllAuthorsWillHaveSigned {
			t.Error("AllAuthorsWillHaveSigned mismatch")
		}
		actualNewManuscriptStatus := getNewManuscriptStatus(actualAllAuthorsWillHaveSigned, c.isThreadReviewable)
		if actualNewManuscriptStatus != c.expectedNewStatus {
			t.Error(fmt.Sprintf("NewManuscriptStatus mismatch, expected %s got %s",
				model.GetManuscriptStatusString(c.expectedNewStatus),
				model.GetManuscriptStatusString(actualNewManuscriptStatus)))
		}
	}
}

type acceptAuthorshipContext struct {
	numAuthors                       int
	numSigner                        int
	numbersAlreadySigned             []int
	isThreadReviewable               bool
	expectedDoesAuthorUpdate         bool
	expectedAllAuthorsWillHaveSigned bool
	expectedNewStatus                model.ManuscriptStatus
}

func TestGetBlockchainSignedHistoricAuthors(t *testing.T) {
	nbce := &nonBootstrapCommandExecution{
		unmarshalledState: &unmarshalledState{
			manuscriptThreads: map[string]*model.StateManuscriptThread{
				"thread": {
					Id:           "thread",
					ManuscriptId: []string{"m1", "m2"},
				},
			},
			manuscripts: map[string]*model.StateManuscript{
				"m1": {
					Id:            "m1",
					ThreadId:      "thread",
					VersionNumber: 0,
					Author: []*model.Author{
						{
							AuthorId:     "a1",
							DidSign:      true,
							AuthorNumber: 0,
						},
						{
							AuthorId:     "a2",
							DidSign:      false,
							AuthorNumber: 1,
						},
					},
				},
				"m2": {
					Id:            "m2",
					ThreadId:      "thread",
					VersionNumber: 1,
					Author: []*model.Author{
						{
							AuthorId:     "a3",
							DidSign:      true,
							AuthorNumber: 2,
						},
						{
							AuthorId:     "a4",
							DidSign:      true,
							AuthorNumber: 3,
						},
					},
				},
			},
		},
	}
	actual := nbce.getBlockchainSignedHistoricAuthors("thread")
	expected := []string{"a1", "a3", "a4"}
	if len(actual) != len(expected) {
		t.Error(fmt.Sprintf("Invalid number of historic authors. Expected %d, got %d",
			len(expected), len(actual)))
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Error(fmt.Sprintf("Historic signer author mismatch for index %d. Expected %s, got %s",
				i, expected[i], actual[i]))
		}
	}
}
