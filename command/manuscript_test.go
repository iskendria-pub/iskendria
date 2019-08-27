package command

import (
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/model"
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
