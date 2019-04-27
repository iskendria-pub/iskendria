package main

import (
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

Also unrelated to this game, dialogs are tested. A dialog is a function
that requires a list of name/value pairs to be set interactively.
You can do "continue" to execute the function with the given
values.

Please enter "help" to start or "exit" to quit.`)

var childDescription = strings.TrimSpace(`
Welcome to the side track of "White against black". You can return
to the game using "exit". Execute "help" to see what you
can do here.
`)

var dialogDescription = strings.TrimSpace(`
Welcome to the dialog test. Please fill some properties and do
something trivial with them.
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
				FullDescription:    dialogDescription,
				OneLineDescription: "Dialog test",
				Name:               "dialog",
				Action:             executeDialogStruct,
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

func executeDialogStruct(outputter cli.Outputter, v *cli.DialogStruct) {
	outputter("Executing the dialog\n")
	outputter(fmt.Sprintf("We take f1 = %v\n", v.F1))
	outputter(fmt.Sprintf("We take f2 = %v\n", v.F2))
	outputter(fmt.Sprintf("We take f3 = %v\n", v.F3))
}
