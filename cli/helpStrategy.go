package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type lineHandlerType func(string) error

type helpStrategy interface {
	run(lineHandler lineHandlerType)
	getFormattedHelpScreenTitle() string
	getOneLineDescription() string
	getName() string
	setParent(parent helpStrategy)
}

type helpStrategyImpl struct {
	parent             *helpStrategyImpl
	fullDescription    string
	oneLineDescription string
	name               string
	formatEscape       string
	stopWords          map[string]bool
}

func (hsi *helpStrategyImpl) run(lineHandler lineHandlerType) {
	outputToStdout(hsi.getFormatEscape())
	defer outputToStdout(UNDO_FORMAT)
	outputToStdout(hsi.fullDescription + "\n\n")
	reader := bufio.NewReader(os.Stdin)
	stop := hsi.nextLine(reader, lineHandler)
	for !stop {
		stop = hsi.nextLine(reader, lineHandler)
	}
}

func (hsi *helpStrategyImpl) nextLine(reader *bufio.Reader, lineHandler lineHandlerType) bool {
	hsi.prompt()
	outputToStdout(UNDO_FORMAT)
	defer outputToStdout(hsi.getFormatEscape())
	line, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return false
	}
	outputToStdout(hsi.getFormatEscape())
	if err := lineHandler(line); err != nil {
		fmt.Println(err)
	}
	_, stop := hsi.stopWords[line]
	return stop
}

func (hsi *helpStrategyImpl) getFormattedHelpScreenTitle() string {
	var sb strings.Builder
	title := hsi.getHelpScreenTitle()
	sb.WriteString("\n" + title + "\n")
	sb.WriteString(strings.Repeat("-", len(title)) + "\n\n")
	return sb.String()
}

func (hsi *helpStrategyImpl) getHelpScreenTitle() string {
	if hsi.parent != nil {
		return hsi.parent.getHelpScreenTitle() + " > " + hsi.oneLineDescription
	}
	return hsi.oneLineDescription
}

func (hsi *helpStrategyImpl) getOneLineDescription() string {
	return hsi.oneLineDescription
}

func (hsi *helpStrategyImpl) prompt() {
	outputToStdout(hsi.getPath() + " |> ")
}

func (hsi *helpStrategyImpl) getPath() string {
	if hsi.parent != nil {
		return hsi.parent.getPath() + "/" + hsi.name
	}
	return hsi.name
}

func (hsi *helpStrategyImpl) getName() string {
	return hsi.name
}

func (hsi *helpStrategyImpl) setParent(parentContextStrategy helpStrategy) {
	switch specificParent := parentContextStrategy.(type) {
	case *helpStrategyImpl:
		hsi.parent = specificParent
	default:
		panic("Unexpected implementation of helpStrategy")
	}
}

func (hsi *helpStrategyImpl) getFormatEscape() string {
	if hsi.formatEscape == "" && hsi.parent != nil {
		return hsi.parent.getFormatEscape()
	}
	return hsi.formatEscape
}
