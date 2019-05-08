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

func TestStringToSliceSet(t *testing.T) {
	resultNil := StringSliceToSet(nil)
	if resultNil == nil || len(resultNil) >= 1 {
		t.Error("Expect map of zero length for input nil")
	}
	resultEmpty := StringSliceToSet([]string{})
	if resultEmpty == nil || len(resultEmpty) >= 1 {
		t.Error("Expect map of zero length for empty slice")
	}
	resultOne := StringSliceToSet([]string{"one"})
	if len(resultOne) != 1 {
		t.Error("Expect map of length one for slice of length one")
	}
	_, found := resultOne["one"]
	if !found {
		t.Error("Word \"one\" was in slice, but not in resulting map")
	}
}

func TestStringSetToSlice(t *testing.T) {
	resultNil := StringSetToSlice(nil)
	if resultNil == nil || len(resultNil) >= 1 {
		t.Error("For string set nil, expected slice of length zero")
	}
	resultEmpty := StringSetToSlice(map[string][]byte{})
	if resultEmpty == nil || len(resultEmpty) >= 1 {
		t.Error("For empty string set, expected slice of length zero")
	}
	resultOne := StringSetToSlice(map[string][]byte{"one": {byte(1)}})
	if len(resultOne) != 1 {
		t.Error("For set of length one, expected slice of length one")
	}
	if resultOne[0] != "one" {
		t.Error("Word \"one\" from set was not put into slice")
	}
}

func TestStringHasAll(t *testing.T) {
	if StringSetHasAll(map[string]bool{"one": true}, []string{"one"}) != true {
		t.Error("When all items of slice are in map, expect true")
	}
	if StringSetHasAll(map[string]bool{"two": true}, []string{"one"}) != false {
		t.Error("When some items of slice are not in map, expect false")
	}
	if StringSetHasAll(map[string]bool{"one": true}, []string{}) != true {
		t.Error("When slice empty, the result should be true")
	}
	if StringSetHasAll(map[string]bool{"one": true}, []string{"one", "two"}) != false {
		t.Error("When only some of slice are in map, result should be false")
	}
}
