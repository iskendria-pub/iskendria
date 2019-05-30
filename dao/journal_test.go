package dao

import (
	"fmt"
	"testing"
)

func TestSortJournals(t *testing.T) {
	jts := getJournalTestItems()
	cases := getCases(jts)
	for i, currentCase := range cases {
		sortJournals(currentCase)
		if !(currentCase[0].JournalId == "03" &&
			currentCase[1].JournalId == "04" &&
			currentCase[2].JournalId == "01" &&
			currentCase[3].JournalId == "02") {
			t.Error(fmt.Sprintf("Case %d not sorted correctly", i))
		}
	}
}

func getJournalTestItems() *journalTestItems {
	firstTitle := "A"
	secondTitle := "Z"
	first := &Journal{
		JournalId: "01",
		Title:     secondTitle,
	}
	second := &Journal{
		JournalId: "02",
		Title:     secondTitle,
	}
	third := &Journal{
		JournalId: "03",
		Title:     firstTitle,
	}
	fourth := &Journal{
		JournalId: "04",
		Title:     firstTitle,
	}
	return &journalTestItems{
		first:  first,
		second: second,
		third:  third,
		fourth: fourth,
	}
}

type journalTestItems struct {
	first,
	second,
	third,
	fourth *Journal
}

func getCases(jts *journalTestItems) [][]*Journal {
	return [][]*Journal{
		{
			jts.first, jts.second, jts.third, jts.fourth,
		},
		{
			jts.first, jts.second, jts.fourth, jts.third,
		},
		{
			jts.first, jts.third, jts.second, jts.fourth,
		},
		{
			jts.first, jts.third, jts.fourth, jts.second,
		},
		{
			jts.first, jts.fourth, jts.second, jts.third,
		},
		{
			jts.first, jts.fourth, jts.third, jts.second,
		},
		{
			jts.second, jts.first, jts.third, jts.fourth,
		},
		{
			jts.second, jts.first, jts.fourth, jts.third,
		},
		{
			jts.second, jts.third, jts.first, jts.fourth,
		},
		{
			jts.second, jts.third, jts.fourth, jts.first,
		},
		{
			jts.second, jts.fourth, jts.first, jts.third,
		},
		{
			jts.second, jts.fourth, jts.third, jts.first,
		},
		{
			jts.third, jts.first, jts.second, jts.fourth,
		},
		{
			jts.third, jts.first, jts.fourth, jts.second,
		},
		{
			jts.third, jts.second, jts.first, jts.fourth,
		},
		{
			jts.third, jts.second, jts.fourth, jts.first,
		},
		{
			jts.third, jts.fourth, jts.first, jts.second,
		},
		{
			jts.third, jts.fourth, jts.second, jts.first,
		},
		{
			jts.fourth, jts.first, jts.second, jts.third,
		},
		{
			jts.fourth, jts.first, jts.third, jts.second,
		},
		{
			jts.fourth, jts.second, jts.first, jts.third,
		},
		{
			jts.fourth, jts.second, jts.third, jts.first,
		},
		{
			jts.fourth, jts.third, jts.first, jts.second,
		},
		{
			jts.fourth, jts.third, jts.second, jts.first,
		},
	}
}
