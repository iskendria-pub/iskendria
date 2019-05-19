package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	HELP        = "help"
	EXIT        = "exit"
	CONTINUE    = "continue"
	CANCEL      = "cancel"
	REVIEW      = "review"
	CLEAR       = "clear"
	UNDO_FORMAT = "\033[0m"
	OK          = "Ok"
)

type Outputter func(string)

func outputToStdout(value string) {
	fmt.Print(value)
}

func getValue(word string, expectedType reflect.Type) (reflect.Value, error) {
	switch expectedType.Kind() {
	case reflect.String:
		return reflect.ValueOf(word), nil
	case reflect.Int32:
		return getValueInt(word, 32)
	case reflect.Int64:
		return getValueInt(word, 64)
	case reflect.Bool:
		value, err := strconv.ParseBool(word)
		if err != nil {
			return reflect.ValueOf(false), errors.New("Invalid boolean value")
		}
		return reflect.ValueOf(value), nil
	default:
		panic("Unsupported type for word: " + word)
	}
}

func getValueInt(word string, numBits int) (reflect.Value, error) {
	value, err := strconv.ParseInt(word, 10, numBits)
	if err != nil {
		return reflect.ValueOf(0), errors.New("Invalid integer value, possibly value out of range")
	}
	rawValue := reflect.ValueOf(value)
	switch numBits {
	case 32:
		typeInt32 := reflect.TypeOf(int32(0))
		return rawValue.Convert(typeInt32), nil
	case 64:
		return rawValue, err
	}
	panic("Unsupported number of bits")
}

type lineGroups []*lineGroup

func (lgs *lineGroups) String() string {
	filled := make([]string, 0)
	for _, lg := range []*lineGroup(*lgs) {
		if !lg.isEmpty() {
			filled = append(filled, lg.String())
		}
	}
	return strings.Join(filled, "\n") + "\n"
}

type lineGroup struct {
	name  string
	lines []string
}

const lineGroupIndent = 2

func (lg *lineGroup) String() string {
	var sb strings.Builder
	sb.WriteString(lg.name + ":\n")
	for _, line := range lg.lines {
		sb.WriteString(strings.Repeat(" ", lineGroupIndent) + line + "\n")
	}
	return sb.String()
}

func (lg *lineGroup) isEmpty() bool {
	return lg.lines == nil || len(lg.lines) == 0
}

type inputSource interface {
	readLine() (string, bool, error)
	open()
	close()
}

type inputSourceConsole struct {
	reader *bufio.Reader
}

var _ inputSource = new(inputSourceConsole)

func (inp *inputSourceConsole) readLine() (string, bool, error) {
	line, err := inp.reader.ReadString('\n')
	return line, false, err
}

func (inp *inputSourceConsole) open() {
	inp.reader = bufio.NewReader(os.Stdin)
}

func (inp *inputSourceConsole) close() {
}

type inputSourceFile struct {
	fname  string
	f      *os.File
	reader *bufio.Reader
}

var _ inputSource = new(inputSourceFile)

func (fi *inputSourceFile) readLine() (string, bool, error) {
	line, err := fi.reader.ReadString('\n')
	fmt.Print(line)
	if err == io.EOF {
		return line, true, err
	}
	return line, false, err
}

func (fi *inputSourceFile) open() {
	var err error
	fi.f, err = os.Open(fi.fname)
	if err != nil {
		panic(err)
	}
	fi.reader = bufio.NewReader(fi.f)
}

func (fi *inputSourceFile) close() {
	_ = fi.f.Close()
}
