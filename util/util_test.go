package util

import (
	"fmt"
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

func TestAbs(t *testing.T) {
	if Abs(int64(1)) != int64(1) {
		t.Error("Abs of positive value fails")
	}
	if Abs(int64(-1)) != int64(1) {
		t.Error("Abs of negative value fails")
	}
	if Abs(int64(0)) != int64(0) {
		t.Error("Abs of zero fails")
	}
	tooNegative := int64(-1) << 63
	if !CheckPanicked(func() { _ = Abs(tooNegative) }) {
		t.Error("Abs did not detect that the smallest negative value has no inverse")
	}
	if Abs(tooNegative+1) != -(tooNegative + 1) {
		t.Error(fmt.Sprintf("Could not handle value: %d", tooNegative+1))
	}
	maxPositive := -(tooNegative + 1)
	if Abs(maxPositive) != maxPositive {
		t.Error(fmt.Sprintf("Could not handle value: %d", maxPositive))
	}
}
