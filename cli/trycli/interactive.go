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

Please enter "help" to start or "exit" to quit.`)

var childDescription = strings.TrimSpace(`
Welcome to the side track of "White against black". You can return
to the game using "exit". Execute "help" to see what you
can do here.
`)

var makeGreen = "\033[32m"

func main() {
	childGroup := cli.NewGroupContext(childDescription, "Side Track", "sideTrack", makeGreen,
		[]interface{}{
			&cli.SingleLineHandler{
				Name:     "and",
				Handler:  and,
				ArgNames: []string{"first value", "second value"},
			},
		},
	)
	context := cli.NewGroupContext(
		description,
		"Black and White",
		"black-white",
		makeGreen,
		[]interface{}{
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
			childGroup,
		})
	childGroup.SetParent(context)
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
