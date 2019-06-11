package main

import (
	"errors"
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/cli"
	"strings"
)

var description = strings.TrimSpace(`
Welcome to "White against black". This is a test of the cli package
which is part of Alexandria. We have a fake dual against black and white.
You can set their strength and let them play. The one you gave the
highest strength wins.

Unrelated to this game, there are test functions to check whether
package cli can detect malformat numbers and boolean values. Out-of-range
values for 32-bit integers and 64-bit integers should be detected.
There is also a test that all types of SingleLineHandler work. Such
a handler can take an Outputter as the first argument, or it can
only take the arguments named in ArgNames and return error. This is
tested with the "expectTrue" command.

Also unrelated to this game, dialogs are tested. A dialog is a function
that requires a list of name/value pairs to be set interactively.
You can do "continue" to execute the function with the given
values.

There are four dialog tests, two without a reference value and two with
a reference value. One of the dialogs with reference values has a fixed
reference value; this dialog does not require an argument to start.
The other uses a dynamic reference value that depends on the input
value supplied when entering the dialog.

The first dialog without reference value and the two dialogs with
reference value have a []string member called "list". When there
is a list member, four additional handlers are generated: add,
insert, removeItem and removeIndex. The last dialog does not have
a list member. Please verify that these four methods are not
present here.

Finally, there is a group to test event paging. After every ten
enter-presses within this group, a message is displayed to simulate
an application event that has to be paged. There is a subgroup
to test that a parent event pager is applied recursively in
subgroups.

Note that your options are sorted alphabetically. Please enter "help" to
start or "exit" to quit.`)

var childDescription = strings.TrimSpace(`
Welcome to the side track of "White against black". You can return
to the game using "exit". Execute "help" to see what you
can do here.
`)

var dialogDescriptionNoRef = strings.TrimSpace(`
Welcome to the dialog test without reference. Please fill some 
properties. We will do something trivial with them on "continue".

Please note that your options are not sorted alphabetically, but
according to the field order of the Golang struct being filled.

The "list" field refers to a []string field in the handling
Golang code. Please edit this field like "list = first second".

This test is without a reference value.
`)

var dialogDescriptionRef = strings.TrimSpace(`
Welcome to the dialog test with reference. Please fill some 
properties. We will do something trivial with them on "continue".
Note that when you enter the dialog, the struct is already
filled with the reference values.

The "list" field refers to a []string field in the handling
Golang code. Please edit this field like "list = first second".

Please note that your options are not sorted alphabetically, but
according to the field order of the Golang struct being filled.
`)

var dialogDescriptionRefArg = strings.TrimSpace(`
Welcome to the dialog test that uses a calculated reference
value. This dialog maintaines two reference values "one"
and "two". To enter the dialog, you have to enter the
number. If the number is not 1 or 2, a message is printed
and the dialog is not entered.
`)

var makeGreen = "\033[32m"

func main() {
	context := &cli.Cli{
		FullDescription:    description,
		OneLineDescription: "Black and White",
		Name:               "black-white",
		FormatEscape:       makeGreen,
		Handlers: []cli.Handler{
			&cli.SingleLineHandler{
				Name:     "strengthWhite",
				Handler:  strengthWhite,
				ArgNames: []string{"Strength of white"},
			},
			&cli.SingleLineHandler{
				Name:     "strengthBlack",
				Handler:  strengthBlack,
				ArgNames: []string{"Strength of black"},
			},
			&cli.SingleLineHandler{
				Name:     "play",
				Handler:  play,
				ArgNames: []string{},
			},
			&cli.SingleLineHandler{
				Name:     "inc32bit",
				Handler:  incInt32,
				ArgNames: []string{"32-bit value to increment"},
			},
			&cli.SingleLineHandler{
				Name:     "inc64bit",
				Handler:  incInt64,
				ArgNames: []string{"64-bit value to increment"},
			},
			&cli.SingleLineHandler{
				Name:     "or",
				Handler:  or,
				ArgNames: []string{"first value", "second value"},
			},
			&cli.SingleLineHandler{
				Name:     "expectTrue",
				Handler:  errorReturnerExpectingTrue,
				ArgNames: []string{"test value"},
			},
			cli.Handler(&cli.Cli{
				FullDescription:    childDescription,
				OneLineDescription: "Side Track",
				Name:               "sideTrack",
				Handlers: []cli.Handler{
					&cli.SingleLineHandler{
						Name:     "and",
						Handler:  and,
						ArgNames: []string{"first value", "second value"},
					},
				},
			}),
			cli.Handler(&cli.StructRunnerHandler{
				FullDescription:    dialogDescriptionNoRef,
				OneLineDescription: "Dialog test without reference",
				Name:               "dialog-noref",
				Action:             executeDialogStruct,
			}),
			cli.Handler(&cli.StructRunnerHandler{
				FullDescription:      dialogDescriptionRef,
				OneLineDescription:   "Dialog test with reference",
				Name:                 "dialog-ref",
				Action:               executeDialogStruct,
				ReferenceValueGetter: createReference,
			}),
			cli.Handler(&cli.StructRunnerHandler{
				FullDescription:              dialogDescriptionRefArg,
				OneLineDescription:           "Dialog test with chosen reference",
				Name:                         "dialog-ref-chosen",
				Action:                       executeDialogStruct,
				ReferenceValueGetter:         chooseReference,
				ReferenceValueGetterArgNames: []string{"chosen number"},
			}),
			cli.Handler(&cli.StructRunnerHandler{
				FullDescription: strings.TrimSpace(`
Dialog without reference and without list.
Methods add, insert, removeItem and removeIndex should not be present`),
				OneLineDescription: "Dialog test without lists",
				Name:               "dialog-nolist",
				Action:             executeDialogStructNoList,
			}),
			cli.Handler(&cli.Cli{
				FullDescription:    "With event paging, parent. You should see a message after each 10 enter presses",
				OneLineDescription: "With event paging",
				Name:               "with-paging",
				EventPager:         countNumberOfCalls,
				Handlers: []cli.Handler{
					&cli.SingleLineHandler{
						Name:     "and",
						Handler:  and,
						ArgNames: []string{"first value", "second value"},
					},
					&cli.Cli{
						FullDescription:    "Sub group - continues counting enters",
						OneLineDescription: "Sub group that continues counting enters",
						Name:               "subgroup",
						Handlers:           []cli.Handler{},
					},
				},
			}),
		},
	}
	context.Run()
}

type strengthsType struct {
	white int32
	black int32
}

var strengths strengthsType

func strengthWhite(outputter cli.Outputter, strength int32) {
	strengths.white = strength
	outputter("Ok\n")
}

func strengthBlack(outputter cli.Outputter, strength int32) {
	strengths.black = strength
	outputter("Ok\n")
}

func play(outputter cli.Outputter) {
	if strengths.white > strengths.black {
		outputter("White wins\n")
	} else if strengths.black > strengths.white {
		outputter("Black wins\n")
	} else {
		outputter("Draw\n")
	}
}

func incInt32(outputter cli.Outputter, v int32) {
	outputter(fmt.Sprintf("%d\n", v+1))
}

func incInt64(outputter cli.Outputter, v int64) {
	outputter(fmt.Sprintf("%d\n", v+1))
}

func or(outputter cli.Outputter, v1, v2 bool) {
	outputter(fmt.Sprintf("%v\n", v1 || v2))
}

func and(outputter cli.Outputter, v1, v2 bool) {
	outputter(fmt.Sprintf("%v\n", v1 && v2))
}

// The sort order should be first, second, fourth. We
// test that the fields are not sorted alphabetically
// in the help text of the dialog.
type DialogStruct struct {
	First  bool
	Second int32
	Fourth string
	List   []string
}

func executeDialogStruct(outputter cli.Outputter, v *DialogStruct) {
	outputter("Executing the dialog\n")
	outputter(fmt.Sprintf("We take first = %v\n", v.First))
	outputter(fmt.Sprintf("We take second = %v\n", v.Second))
	outputter(fmt.Sprintf("We take fourth = %v\n", v.Fourth))
	outputter(fmt.Sprintf("The list field has length %d\n", len(v.List)))
	for i, item := range v.List {
		outputter(fmt.Sprintf("  List item #%d has value %s\n", i+1, item))
	}
}

func errorReturnerExpectingTrue(testValue bool) error {
	if !testValue {
		return errors.New("We expected to see true here")
	}
	return nil
}

func createReference(_ cli.Outputter) *DialogStruct {
	return &DialogStruct{
		First:  true,
		Second: 8,
		Fourth: "Martijn Dirkse",
		List:   []string{"firstListItem", "secondListItem"},
	}
}

func chooseReference(outputter cli.Outputter, number int32) *DialogStruct {
	references := map[int32]*DialogStruct{
		1: {
			First:  false,
			Second: 45,
			Fourth: "Martijn Dirkse",
		},
		2: {
			First:  true,
			Second: 43,
			Fourth: "Arri Dirkse",
			List:   []string{"firstListItem"},
		},
	}
	reference, found := references[number]
	if !found {
		outputter(fmt.Sprintf("Number not found: %d\n", number))
		return nil
	}
	return reference
}

var callCount int

func countNumberOfCalls(outputter cli.Outputter) {
	callCount++
	if callCount%10 == 0 {
		outputter(fmt.Sprintf("*** Number of enters: %d\n", callCount))
	}
}

type DialogStructNoList struct {
	The32Bit  int32
	The64Bit  int64
	TheBool   bool
	TheString string
}

func executeDialogStructNoList(outputter cli.Outputter, input *DialogStructNoList) {
	outputter(fmt.Sprintf("The32Bit:  %d\n", input.The32Bit))
	outputter(fmt.Sprintf("The64Bit:  %d\n", input.The64Bit))
	outputter(fmt.Sprintf("TheBool:   %v\n", input.TheBool))
	outputter(fmt.Sprintf("TheString: %s\n", input.TheString))
}
