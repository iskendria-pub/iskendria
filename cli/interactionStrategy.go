package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type lineHandlerType func(string) error

type interactionStrategy interface {
	run(lineHandler lineHandlerType)
	getFormattedHelpScreenTitle() string
	getOneLineDescription() string
	getName() string
	setParent(parent interactionStrategy)
}

type interactionStrategyImpl struct {
	parent             *interactionStrategyImpl
	fullDescription    string
	oneLineDescription string
	name               string
	formatEscape       string
	stopWords          map[string]bool
	eventPager         func(Outputter)
}

func (isi *interactionStrategyImpl) run(lineHandler lineHandlerType) {
	outputToStdout(isi.getFormatEscape())
	defer outputToStdout(UNDO_FORMAT)
	outputToStdout(isi.fullDescription + "\n\n")
	reader := bufio.NewReader(os.Stdin)
	stop := isi.nextLine(reader, lineHandler)
	for !stop {
		stop = isi.nextLine(reader, lineHandler)
	}
}

func (isi *interactionStrategyImpl) nextLine(reader *bufio.Reader, lineHandler lineHandlerType) bool {
	isi.prompt()
	outputToStdout(UNDO_FORMAT)
	defer isi.afterEnter()
	line, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return false
	}
	outputToStdout(isi.getFormatEscape())
	if err := lineHandler(line); err != nil {
		fmt.Println(err)
	}
	_, stop := isi.stopWords[line]
	return stop
}

func (isi *interactionStrategyImpl) afterEnter() {
	outputToStdout(isi.getFormatEscape())
	isi.runEventPager()
}

func (isi *interactionStrategyImpl) runEventPager() {
	switch {
	case isi.eventPager != nil:
		isi.eventPager(outputToStdout)
	case isi.parent != nil:
		isi.parent.runEventPager()
	}
}

func (isi *interactionStrategyImpl) getFormattedHelpScreenTitle() string {
	var sb strings.Builder
	title := isi.getHelpScreenTitle()
	sb.WriteString("\n" + title + "\n")
	sb.WriteString(strings.Repeat("-", len(title)) + "\n\n")
	return sb.String()
}

func (isi *interactionStrategyImpl) getHelpScreenTitle() string {
	if isi.parent != nil {
		return isi.parent.getHelpScreenTitle() + " > " + isi.oneLineDescription
	}
	return isi.oneLineDescription
}

func (isi *interactionStrategyImpl) getOneLineDescription() string {
	return isi.oneLineDescription
}

func (isi *interactionStrategyImpl) prompt() {
	outputToStdout(isi.getPath() + " |> ")
}

func (isi *interactionStrategyImpl) getPath() string {
	if isi.parent != nil {
		return isi.parent.getPath() + "/" + isi.name
	}
	return isi.name
}

func (isi *interactionStrategyImpl) getName() string {
	return isi.name
}

func (isi *interactionStrategyImpl) setParent(parentContextStrategy interactionStrategy) {
	switch specificParent := parentContextStrategy.(type) {
	case *interactionStrategyImpl:
		isi.parent = specificParent
	default:
		panic("Unexpected implementation of interactionStrategy")
	}
}

func (isi *interactionStrategyImpl) getFormatEscape() string {
	if isi.formatEscape == "" && isi.parent != nil {
		return isi.parent.getFormatEscape()
	}
	return isi.formatEscape
}
