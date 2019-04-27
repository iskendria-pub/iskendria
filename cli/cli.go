package cli

import (
	"errors"
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/util"
	"reflect"
	"sort"
	"strings"
)

type Handler interface {
	build() runnableHandler
}

type runnableHandler interface {
	getName() string
	handleLine(words []string) error
}

var _ Handler = new(SingleLineHandler)
var _ runnableHandler = new(SingleLineHandler)

type SingleLineHandler struct {
	Name     string
	Handler  interface{}
	ArgNames []string
}

func (slh *SingleLineHandler) build() runnableHandler {
	reflectHandlerType := reflect.TypeOf(slh.Handler)
	if reflectHandlerType.Kind() != reflect.Func {
		panic(fmt.Sprintf("Handler is not a function: %+v", slh))
	}
	expectedNumFunctionArgs := len(slh.ArgNames) + 1
	if reflectHandlerType.NumIn() != expectedNumFunctionArgs {
		panic(fmt.Sprintf(
			"Number of handler arguments does not match number of argument names or outputter function is missing: %+v",
			slh))
	}
	if reflectHandlerType.NumOut() != 0 {
		panic(fmt.Sprintf("Handler should not return anything: %+v", slh))
	}
	reflectFirstArgumentType := reflectHandlerType.In(0)
	if reflectFirstArgumentType.Kind() != reflect.Func {
		panic(fmt.Sprintf("The first argument of a handler should be of type func(string): %+v", slh))
	}
	if reflectFirstArgumentType.NumIn() != 1 {
		panic(fmt.Sprintf("The first argument of a handler should be a function with one argument: %+v", slh))
	}
	if reflectFirstArgumentType.NumOut() != 0 {
		panic(fmt.Sprintf("The first argument of a handler should be a function without outputs: %+v", slh))
	}
	if reflectFirstArgumentType.In(0).Kind() != reflect.String {
		panic(fmt.Sprintf(
			"The first argument of a handler should be a function with a string argument: %+v", slh))
	}
	return slh
}

func (slh *SingleLineHandler) getName() string {
	return slh.Name
}

func (command *SingleLineHandler) handleLine(words []string) error {
	expectedNumArgs := len(command.ArgNames)
	actualNumArgs := len(words) - 1
	if expectedNumArgs != actualNumArgs {
		return errors.New(fmt.Sprintf("Wrong number of arguments, expected %d, got %d",
			expectedNumArgs, actualNumArgs))
	}
	values, err := getValues(command, words[1:])
	if err != nil {
		return err
	}
	callResultValue := reflect.ValueOf(command.Handler).Call(values)
	if len(callResultValue) != 0 {
		panic("Expected exactly one result")
	}
	return nil
}

func getValues(handler *SingleLineHandler, argWords []string) ([]reflect.Value, error) {
	numArgWords := len(handler.ArgNames)
	values := make([]reflect.Value, numArgWords+1)
	values[0] = reflect.ValueOf(outputToStdout)
	for argWordsIndex, word := range argWords {
		allArgsIndex := argWordsIndex + 1
		argumentType := reflect.TypeOf(handler.Handler).In(allArgsIndex)
		value, err := getValue(word, argumentType)
		if err != nil {
			return values, err
		}
		values[allArgsIndex] = value
	}
	return values, nil
}

var _ Handler = new(Cli)

type Cli struct {
	FullDescription    string
	OneLineDescription string
	Name               string
	FormatEscape       string
	Handlers           []Handler
}

func (c *Cli) Run() {
	c.buildMain().run()
}

func (c *Cli) build() runnableHandler {
	return runnableHandler(c.buildMain())
}

func (c *Cli) buildMain() *groupContextRunnable {
	runnableHandlers := make([]runnableHandler, len(c.Handlers))
	for i := range c.Handlers {
		runnableHandlers[i] = c.Handlers[i].build()
	}
	result := &groupContextRunnable{
		helpStrategy: &helpStrategyImpl{
			fullDescription:    c.FullDescription,
			oneLineDescription: c.OneLineDescription,
			name:               c.Name,
			formatEscape:       c.FormatEscape,
			stopWords:          map[string]bool{EXIT: true},
		},
		handlers: runnableHandlers,
	}
	var helpHandler = &SingleLineHandler{
		Name:     HELP,
		Handler:  func(outputter Outputter) { result.help(outputter) },
		ArgNames: []string{},
	}
	var exitHandler = &SingleLineHandler{
		Name:     EXIT,
		Handler:  func(Outputter) {},
		ArgNames: []string{},
	}
	result.addRunnableHandler(helpHandler.build())
	result.addRunnableHandler(exitHandler.build())
	result.init()
	return result
}

type groupContextRunnable struct {
	helpStrategy helpStrategy
	handlers     []runnableHandler
}

func (gcr *groupContextRunnable) getName() string {
	return gcr.helpStrategy.getName()
}

func (gcr *groupContextRunnable) init() {
	for _, handler := range gcr.handlers {
		switch specificHandler := handler.(type) {
		case *groupContextRunnable:
			specificHandler.helpStrategy.setParent(gcr.helpStrategy)
		case *dialogContextRunnable:
			specificHandler.helpStrategy.setParent(gcr.helpStrategy)
		}
	}
	sort.Slice(gcr.handlers, func(i, j int) bool {
		return lessRunnableHandler(gcr.handlers[i], gcr.handlers[j])
	})
}

func lessRunnableHandler(first, second runnableHandler) bool {
	return first.getName() < second.getName()
}

func (gcr *groupContextRunnable) addRunnableHandler(handler runnableHandler) {
	gcr.handlers = append(gcr.handlers, handler)
}

func (gcr *groupContextRunnable) help(outputter Outputter) {
	outputter(gcr.helpStrategy.getFormattedHelpScreenTitle())
	gcr.showHandlers(outputter)
}

func (gcr *groupContextRunnable) showHandlers(outputter Outputter) {
	groups, commands, dialogs := gcr.splitGroupsCommandsAndDialogs()
	parts := lineGroups{
		listSingleLineHandlers(commands),
		listDialogContextRunnables(dialogs),
		listGroupContextRunnables(groups),
	}
	outputter(parts.String())
}

func (gcr *groupContextRunnable) splitGroupsCommandsAndDialogs() (
	[]*groupContextRunnable, []*SingleLineHandler, []*dialogContextRunnable) {
	groups := make([]*groupContextRunnable, 0)
	commands := make([]*SingleLineHandler, 0)
	dialogs := make([]*dialogContextRunnable, 0)
	for _, handler := range gcr.handlers {
		switch specificHandler := handler.(type) {
		case *groupContextRunnable:
			groups = append(groups, specificHandler)
		case *SingleLineHandler:
			commands = append(commands, specificHandler)
		case *dialogContextRunnable:
			dialogs = append(dialogs, specificHandler)
		default:
			panic("Unknown handler type")
		}
	}
	return groups, commands, dialogs
}

func (gcr *groupContextRunnable) run() {
	gcr.helpStrategy.run(func(line string) error {
		return gcr.executeLine(line)
	})
}

func (gcr *groupContextRunnable) executeLine(line string) error {
	words := strings.Split(strings.TrimSpace(line), " ")
	handler := getHandler(words[0], gcr.handlers)
	if handler == nil {
		return errors.New("Name does not match a group or a command: " + words[0])
	}
	return handler.handleLine(words)
}

func (gcr *groupContextRunnable) handleLine(words []string) error {
	if len(words) != 1 {
		return errors.New("Entering a group does not require arguments")
	}
	gcr.run()
	return nil
}

type StructRunnerHandler struct {
	FullDescription    string
	OneLineDescription string
	Name               string
	Action             interface{}
}

var _ Handler = new(StructRunnerHandler)

func (srh *StructRunnerHandler) build() runnableHandler {
	s := reflect.TypeOf(srh.Action).In(1).Elem()
	if s.Kind() != reflect.Struct {
		panic("The second argument of Action is expected to be a pointer to struct")
	}
	containerValue := reflect.New(s)
	result := &dialogContextRunnable{
		helpStrategy: &helpStrategyImpl{
			fullDescription:    srh.FullDescription,
			oneLineDescription: srh.OneLineDescription,
			name:               srh.Name,
			stopWords:          map[string]bool{CANCEL: true, CONTINUE: true},
		},
		containerValue: containerValue,
		action:         srh.Action,
	}
	helpHandler := &SingleLineHandler{
		Name:     HELP,
		Handler:  result.help,
		ArgNames: []string{},
	}
	reviewHandler := &SingleLineHandler{
		Name:     REVIEW,
		Handler:  result.review,
		ArgNames: []string{},
	}
	continueHandler := &SingleLineHandler{
		Name:     CONTINUE,
		Handler:  result.doContinue,
		ArgNames: []string{},
	}
	cancelHandler := &SingleLineHandler{
		Name:     CANCEL,
		Handler:  result.cancel,
		ArgNames: []string{},
	}
	result.handlers = []runnableHandler{
		helpHandler.build(),
		reviewHandler.build(),
		continueHandler.build(),
		cancelHandler.build(),
	}
	dialogPropertyHandlers := make([]runnableHandler, s.NumField())
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if f.Name != strings.Title(f.Name) {
			panic("Field is not exported: " + f.Name)
		}
		dph := &dialogPropertyHandler{
			name:           util.UnTitle(f.Name),
			propertyType:   f.Type,
			fieldNumber:    i,
			containerValue: containerValue,
		}
		dialogPropertyHandlers[i] = dph
	}
	result.handlers = append(result.handlers, dialogPropertyHandlers...)
	return result
}

type dialogContextRunnable struct {
	helpStrategy   helpStrategy
	handlers       []runnableHandler
	containerValue reflect.Value
	action         interface{}
}

var _ runnableHandler = new(dialogContextRunnable)

func (dcr *dialogContextRunnable) getName() string {
	return dcr.helpStrategy.getName()
}

func (dcr *dialogContextRunnable) help(outputter Outputter) {
	outputter(dcr.helpStrategy.getFormattedHelpScreenTitle())
	dcr.showHandlers(outputter)
}

func (dcr *dialogContextRunnable) showHandlers(outputter Outputter) {
	commands, properties := dcr.splitCommandsAndProperties()
	parts := lineGroups{
		listDialogPropertyHandlers(properties),
		listSingleLineHandlers(commands),
	}
	outputter(parts.String())
}

func (dcr *dialogContextRunnable) splitCommandsAndProperties() (
	[]*SingleLineHandler, []*dialogPropertyHandler) {
	commands := make([]*SingleLineHandler, 0)
	properties := make([]*dialogPropertyHandler, 0)
	for _, handler := range dcr.handlers {
		switch specificHandler := handler.(type) {
		case *SingleLineHandler:
			commands = append(commands, specificHandler)
		case *dialogPropertyHandler:
			properties = append(properties, specificHandler)
		default:
			panic("Unknown handler type")
		}
	}
	return commands, properties
}

func (dcr *dialogContextRunnable) review(outputter Outputter) {
	var propertyHandlers []*dialogPropertyHandler
	_, propertyHandlers = dcr.splitCommandsAndProperties()
	numFields := dcr.containerValue.Elem().NumField()
	result := NewTable(numFields, 2)
	for i := 0; i < numFields; i++ {
		result.Set(i, 0, propertyHandlers[i].name)
		result.Set(i, 1, fmt.Sprintf("%v", dcr.containerValue.Elem().Field(i).Interface()))
	}
	outputter(result.String())
}

func (dcr *dialogContextRunnable) doContinue(outputter Outputter) {
	actionValue := reflect.ValueOf(dcr.action)
	actionValue.Call([]reflect.Value{
		reflect.ValueOf(outputter),
		dcr.containerValue,
	})
}

func (dcr *dialogContextRunnable) cancel(_ Outputter) {
}

func (dcr *dialogContextRunnable) handleLine(words []string) error {
	if len(words) != 1 {
		return errors.New("Entering a group does not require arguments")
	}
	dcr.run()
	return nil
}

func (dcr *dialogContextRunnable) run() {
	dcr.helpStrategy.run(func(line string) error {
		return dcr.executeLine(line)
	})
}

func (dcr *dialogContextRunnable) executeLine(line string) error {
	rawWords := strings.Split(line, "=")
	words := make([]string, len(rawWords))
	for i, rawWord := range rawWords {
		words[i] = strings.TrimSpace(rawWord)
	}
	handler := getHandler(words[0], dcr.handlers)
	if handler == nil {
		return errors.New("Name does not match a command or a property: " + words[0])
	}
	return handler.handleLine(words)
}

type dialogPropertyHandler struct {
	name           string
	propertyType   reflect.Type
	fieldNumber    int
	containerValue reflect.Value
}

var _ runnableHandler = new(dialogPropertyHandler)

func (dph *dialogPropertyHandler) getName() string {
	return dph.name
}

func (dph *dialogPropertyHandler) handleLine(words []string) error {
	value := reflect.Zero(dph.propertyType)
	var err error
	if len(words) == 2 && words[1] != "" {
		if value, err = getValue(words[1], dph.propertyType); err != nil {
			return errors.New(fmt.Sprintf("Type mismatch for property %s: %s",
				dph.name, words[1]))
		}
	}
	setField(dph.containerValue, dph.fieldNumber, value)
	return nil
}

func listSingleLineHandlers(handlers []*SingleLineHandler) *lineGroup {
	if len(handlers) == 0 {
		return &lineGroup{}
	}
	var lines []string
	for _, handler := range handlers {
		items := []string{handler.Name}
		for _, arg := range handler.ArgNames {
			items = append(items, "<"+arg+">")
		}
		lines = append(lines, strings.Join(items, " "))
	}
	return &lineGroup{
		name:  "Commands",
		lines: lines,
	}
}

func listGroupContextRunnables(contexts []*groupContextRunnable) *lineGroup {
	if len(contexts) == 0 {
		return &lineGroup{}
	}
	var lines []string
	for _, context := range contexts {
		lines = append(lines, fmt.Sprintf("%s - %s",
			context.helpStrategy.getName(),
			context.helpStrategy.getOneLineDescription()))
	}
	return &lineGroup{
		name:  "Groups",
		lines: lines,
	}
}

func listDialogContextRunnables(dialogs []*dialogContextRunnable) *lineGroup {
	if len(dialogs) == 0 {
		return &lineGroup{}
	}
	var lines []string
	for _, dialog := range dialogs {
		lines = append(lines, fmt.Sprintf("%s - %s",
			dialog.helpStrategy.getName(),
			dialog.helpStrategy.getOneLineDescription()))
	}
	return &lineGroup{
		name:  "Dialogs",
		lines: lines,
	}
}

func listDialogPropertyHandlers(handlers []*dialogPropertyHandler) *lineGroup {
	if len(handlers) == 0 {
		return &lineGroup{}
	}
	var lines []string
	for _, handler := range handlers {
		lines = append(lines, handler.name)
	}
	return &lineGroup{
		name:  "Properties that can be set:",
		lines: lines,
	}
}

func getHandler(name string, handlers []runnableHandler) runnableHandler {
	for _, handler := range handlers {
		if handler.getName() == name {
			return handler
		}
	}
	return nil
}

type DialogStruct struct {
	F1 bool
	F2 int32
	F3 string
}
