package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Context interface {
	Run()
}

type Outputter func(string)

func outputToStdout(value string) {
	fmt.Print(value)
}

const (
	HELP = "help"
	EXIT = "exit"
	UNDO_FORMAT = "\033[0m"
)

type SingleLineHandler struct {
	Name     string
	Handler  interface{}
	ArgNames []string
}

type groupContext struct {
	description string
	title    string
	pathTag string
	formatEscape string
	handlers []*SingleLineHandler
}

func NewGroupContext(description, title, pathTag, formatEscape string, inputHandlers []*SingleLineHandler) Context {
	handlers := make([]*SingleLineHandler, len(inputHandlers))
	for i, handler := range inputHandlers {
		checkHandler(handler)
		handlers[i] = copyHandler(inputHandlers[i])
	}
	result := &groupContext{
		description: description,
		title:    title,
		pathTag: pathTag,
		formatEscape: formatEscape,
		handlers: handlers,
	}
	var helpHandler *SingleLineHandler = &SingleLineHandler{
		Name:     HELP,
		Handler:  func(outputter Outputter) { result.help(outputter) },
		ArgNames: []string{},
	}
	var exitHandler *SingleLineHandler = &SingleLineHandler{
		Name:     EXIT,
		Handler:  func(Outputter) {},
		ArgNames: []string{},
	}
	checkHandler(helpHandler)
	checkHandler(exitHandler)
	result.addHandler(helpHandler)
	result.addHandler(exitHandler)
	result.init()
	return result
}

func checkHandler(handler *SingleLineHandler) {
	reflectHandlerType := reflect.TypeOf(handler.Handler)
	if reflectHandlerType.Kind() != reflect.Func {
		panic("Handler is not a function")
	}
	expectedNumFunctionArgs := len(handler.ArgNames) + 1
	if reflectHandlerType.NumIn() != expectedNumFunctionArgs {
		panic("Number of handler arguments does not match number of argument names or outputter function is missing")
	}
	if reflectHandlerType.NumOut() != 0 {
		panic("Handler should not return anything")
	}
	reflectFirstArgumentType := reflectHandlerType.In(0)
	if reflectFirstArgumentType.Kind() != reflect.Func {
		panic("The first argument of a handler should be of type func(string)")
	}
	if reflectFirstArgumentType.NumIn() != 1 {
		panic("The first argument of a handler should be a function with one argument")
	}
	if reflectFirstArgumentType.NumOut() != 0 {
		panic("The first argument of a handler should be a function without outputs")
	}
	if reflectFirstArgumentType.In(0).Kind() != reflect.String {
		panic("The first argument of a handler should be a function with a string argument")
	}
}

func copyHandler(orig *SingleLineHandler) *SingleLineHandler {
	copyArgNames := append(make([]string, 0), orig.ArgNames...)
	return &SingleLineHandler{
		Name:     orig.Name,
		Handler:  orig.Handler,
		ArgNames: copyArgNames,
	}
}

func (cg *groupContext) help(outputter Outputter) {
	outputter("\n" + cg.title + "\n")
	outputter(strings.Repeat("-", len(cg.title)) + "\n\n")
	cg.listCommands(outputter)
}

func (cg *groupContext) listCommands(outputter Outputter) {
	var sb strings.Builder
	for _, handler := range cg.handlers {
		items := []string{handler.Name}
		for _, arg := range handler.ArgNames {
			items = append(items, "<"+arg+">")
		}
		sb.WriteString(strings.Join(items, " ") + "\n")
	}
	outputter(sb.String())
}

func (cg *groupContext) addHandler(handler *SingleLineHandler) {
	cg.handlers = append(cg.handlers, copyHandler(handler))
}

func (cg *groupContext) init() {
	sort.Slice(cg.handlers, func(i, j int) bool {
		return cg.handlers[i].Name < cg.handlers[j].Name
	})
}

func (cg *groupContext) Run() {
	outputToStdout(cg.formatEscape)
	defer outputToStdout(UNDO_FORMAT)
	outputToStdout(cg.description + "\n\n")
	reader := bufio.NewReader(os.Stdin)
	stop := cg.nextLine(reader)
	for !stop {
		stop = cg.nextLine(reader)
	}
}

func (cg *groupContext) nextLine(reader *bufio.Reader) bool {
	cg.prompt()
	outputToStdout(UNDO_FORMAT)
	defer outputToStdout(cg.formatEscape)
	line, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return false
	}
	if line == EXIT {
		return true
	}
	outputToStdout(cg.formatEscape)
	if err := cg.executeLine(line); err != nil {
		fmt.Println(err)
	}
	return false
}

func (cg *groupContext) prompt() {
	outputToStdout(cg.pathTag + " |> ")
}

func (cg *groupContext) executeLine(line string) error {
	words := strings.Split(strings.TrimSpace(line), " ")
	handler := cg.getHandler(words[0])
	if handler == nil {
		return errors.New("Command not found: " + words[0])
	}
	expectedNumArgs := len(handler.ArgNames)
	actualNumArgs := len(words) - 1
	if expectedNumArgs != actualNumArgs {
		return errors.New(fmt.Sprintf("Wrong number of arguments, expected %d, got %d",
			expectedNumArgs, actualNumArgs))
	}
	values, err := cg.getValues(handler, words[1:])
	if err != nil {
		return err
	}
	callResultValue := reflect.ValueOf(handler.Handler).Call(values)
	if len(callResultValue) != 0 {
		panic("Expected exactly one result")
	}
	return nil
}

func (cg *groupContext) getValues(handler *SingleLineHandler, argWords []string) ([]reflect.Value, error) {
	numArgWords := len(handler.ArgNames)
	values := make([]reflect.Value, numArgWords+1)
	values[0] = reflect.ValueOf(outputToStdout)
	for argWordsIndex, word := range argWords {
		allArgsIndex := argWordsIndex + 1
		argumentType := reflect.TypeOf(handler.Handler).In(allArgsIndex)
		value, err := cg.getValue(word, argumentType)
		if err != nil {
			return values, err
		}
		values[allArgsIndex] = value
	}
	return values, nil
}

func (cg *groupContext) getValue(word string, expectedType reflect.Type) (reflect.Value, error) {
	switch expectedType.Kind() {
	case reflect.String:
		return reflect.ValueOf(word), nil
	case reflect.Int32:
		return cg.getValueInt(word, 32)
	case reflect.Int64:
		return cg.getValueInt(word, 64)
	case reflect.Bool:
		value, err := strconv.ParseBool(word)
		if err != nil {
			return reflect.ValueOf(false), errors.New("Invalid boolean value")
		}
		return reflect.ValueOf(value), nil
	default:
		panic("Unsupported type")
	}
}

func (cg *groupContext) getValueInt(word string, numBits int) (reflect.Value, error) {
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

func (cg *groupContext) getHandler(name string) *SingleLineHandler {
	for _, handler := range cg.handlers {
		if handler.Name == name {
			return handler
		}
	}
	return nil
}

type tableType struct {
	numRows int
	numCols int
	data    []string
}

func newTable(numRows, numCols int) *tableType {
	return &tableType{
		numRows: numRows,
		numCols: numCols,
		data:    make([]string, numRows*numCols),
	}
}

func (table *tableType) format() string {
	fieldLengths := table.getFieldLengths()
	var sb strings.Builder
	for row := 0; row < table.numRows; row++ {
		rowOutput := make([]string, table.numCols)
		for col := 0; col < table.numCols; col++ {
			rowOutput[col] = table.formatField(table.get(row, col), fieldLengths[col])
		}
		sb.WriteString(strings.Join(rowOutput, " "))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (table *tableType) formatField(value string, fieldLength int) string {
	// Left align
	return value + strings.Repeat(" ", fieldLength-len(value))
}

func (table *tableType) getFieldLengths() []int {
	fieldLengths := make([]int, table.numCols)
	for col := 0; col < table.numCols; col++ {
		fieldLengths[col] = table.getFieldLength(col)
	}
	return fieldLengths
}

func (table *tableType) getFieldLength(col int) int {
	result := 0
	for row := 0; row < table.numRows; row++ {
		l := len(table.get(row, col))
		if l > result {
			result = l
		}
	}
	return result
}

func (table *tableType) get(row, col int) string {
	return table.data[table.getIndex(row, col)]
}

func (table *tableType) getIndex(row int, col int) int {
	return table.numCols*row + col
}

func (table *tableType) set(row, col int, value string) {
	table.data[table.getIndex(row, col)] = value
}
