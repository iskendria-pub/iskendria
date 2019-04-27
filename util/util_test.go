package util

import (
	"strings"
	"testing"
)

func TestTitle(t *testing.T) {
	title := strings.Title("someString")
	if title != "SomeString" {
		t.Error("strings.Title does not behave as expected")
	}
}

func TestUnTitle(t *testing.T) {
	title := "SomeString"
	untitled := UnTitle(title)
	if untitled != "someString" {
		t.Error("First character should have become lower case")
	}
	if UnTitle("") != "" {
		t.Error("UnTitle does not handle empty string properly")
	}
}
