package cli

import (
	"errors"
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/util"
	"reflect"
	"sort"
	"strings"
)

/*
This package implements a generic interactive command-line interface.
To use it, fill a cli.Cli struct with help text and instances of
Handler. Then call the Run method. See trycli/interactive.go to see
how the cli package is used. When the cli starts, the end user gets
a welcome message and a prompt. She can enter commands and sees
a prompt after a command is completed. The user defines the
user interface by filling Cli with help texts and Handler instances.

Cli itself implements Handler because groups of commands can be nested.
Other Handler implementations are SingleLineHandler and
StructRunnerHandler.

SingleLineHandler wraps a Golang-function with
arbitrary arguments of type bool, int32, int64 or string.
There are two possibilities for a SingleLineHandler:

* It can take an extra argument of type Outputter that is called
to report outputs.
* It can take only the named arguments and return a value
that implements error.

StructRunnerHandler wraps a function taking a struct pointer,
the Action. This handler starts a dialog allowing the end user
to set all fields of the struct. A command "continue" is added
automatically that allows the end user to call the Action.
There is also "cancel" to go back without executing the
Action. The end user can do "review" to get an overview of
the entered properties.

StructRunnerHandler has an optional ReferenceValueGetter, a
function that returns a reference value. The reference value
should produce a value of the type needed by the Action.
When the dialog is started, the read value is initialized
with the reference value. There is a "clear" function that
fills the read value with the reference value, or the zero
values of the struct if there is no ReferenceValueGetter.

The cli is implemented as follows. The top-level Handler's
"build" method is executed to get a runnableHandler instance.
This is done recursively resulting in a nested structure of
runnableHandler instances. There are "var _ Handler = ... "
and "var _ runnableHandler = ..." to enforce that every
type implements the intended interfaces.

The runnableHandler implementations groupContextRunnable and
dialogContextRunnable have an interactionStrategy that implements
reading and parsing user input and that provides helper
functions for formatting. When the user has entered a line,
the parsed line is fed to groupContextRunnable.executeLine()
or dialogContextRunnable.executeLine(). The
groupContextRunnable or dialogContextRunnable selects the
appropriate runnableHandler and executes its handleLine
method. groupContextRunnable and dialogContextRunnable are
themselves runnableHandler implementations because a
dialog is called from a command group and a command group
can be the child of another command group. To summarize:
A command wrapped by a runnableHandler is executed using its
handleLine() method. When the runnableHandler is a group,
it has an executeLine() to select and run a command within
the group.

The interaction strategy does not provide help information.
groupContextRunnable and dialogContextRunnable have help
methods. The build() methods of Cli and StructRunnerHandler
wrap these help methods in runnableHandler instances
and register these handlers. This registration ensures
that a help screen's list of options shows the help
function itself as one of the options. dialogRunnableHandler
applies the same idea to provide "continue", "cancel",
"review" and "clear" commands.
*/
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
	switch reflectHandlerType.NumOut() {
	case 0:
		slh.checkSingleLineHandlerThatTakesOutputter(reflectHandlerType)
	case 1:
		slh.checkSingleLineHandlerThatReturnsError(reflectHandlerType)
	default:
		panic(fmt.Sprintf("Handler should return zero or one value: %+v", slh))
	}
	return slh
}

func (slh *SingleLineHandler) checkSingleLineHandlerThatTakesOutputter(reflectHandlerType reflect.Type) {
	expectedNumFunctionArgs := len(slh.ArgNames) + 1
	if reflectHandlerType.NumIn() != expectedNumFunctionArgs {
		panic(fmt.Sprintf(
			"Number of handler arguments does not match number of argument names or outputter function is missing: %+v",
			slh))
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
}

func (slh *SingleLineHandler) checkSingleLineHandlerThatReturnsError(reflectHandlerType reflect.Type) {
	if reflectHandlerType.NumIn() != len(slh.ArgNames) {
		panic(fmt.Sprintf(
			"Number of handler arguments does not match number of argument names: %+v",
			slh))
	}
	reflectReturnType := reflectHandlerType.Out(0)
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if reflectReturnType == errorType {
		return
	}
	if !reflectReturnType.Implements(errorType) {
		panic(fmt.Sprintf(
			"Return type does not implement error: %s", reflectReturnType))
	}
}

func (slh *SingleLineHandler) getName() string {
	return slh.Name
}

func (slh *SingleLineHandler) handleLine(words []string) error {
	expectedNumArgs := len(slh.ArgNames)
	actualNumArgs := len(words) - 1
	if expectedNumArgs != actualNumArgs {
		return errors.New(fmt.Sprintf("Wrong number of arguments, expected %d, got %d",
			expectedNumArgs, actualNumArgs))
	}
	switch reflect.TypeOf(slh.Handler).NumOut() {
	case 0:
		return slh.handleLineUsingOutputter(words[1:])
	case 1:
		return slh.handleLineUsingReturnedError(words[1:])
	}
	return nil
}

func (slh *SingleLineHandler) handleLineUsingOutputter(argWords []string) error {
	values, err := getValues(slh, argWords, func(i int) int { return i + 1 })
	if err != nil {
		return err
	}
	values = append([]reflect.Value{reflect.ValueOf(outputToStdout)}, values...)
	callResultValue := reflect.ValueOf(slh.Handler).Call(values)
	if len(callResultValue) != 0 {
		panic("Did not expect a result")
	}
	return nil
}

func (slh *SingleLineHandler) handleLineUsingReturnedError(argWords []string) error {
	values, err := getValues(slh, argWords, func(i int) int { return i })
	if err != nil {
		return err
	}
	callResultValue := reflect.ValueOf(slh.Handler).Call(values)
	callResult := callResultValue[0].Interface()
	if callResult == nil {
		outputToStdout(OK + "\n")
		return nil
	}
	err, isConvertedSuccessfully := callResult.(error)
	if !isConvertedSuccessfully {
		panic(fmt.Sprintf("Did not get proper error object: %v", callResultValue))
	}
	return err
}

func getValues(handler *SingleLineHandler, argWords []string, namedArgIdxToArgIdx func(int) int) ([]reflect.Value, error) {
	numArgWords := len(handler.ArgNames)
	values := make([]reflect.Value, numArgWords)
	for index, word := range argWords {
		argumentType := reflect.TypeOf(handler.Handler).In(namedArgIdxToArgIdx(index))
		value, err := getValue(word, argumentType)
		if err != nil {
			return values, err
		}
		values[index] = value
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
		interactionStrategy: &interactionStrategyImpl{
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
	interactionStrategy interactionStrategy
	handlers            []runnableHandler
}

var _ runnableHandler = new(groupContextRunnable)

func (gcr *groupContextRunnable) getName() string {
	return gcr.interactionStrategy.getName()
}

func (gcr *groupContextRunnable) init() {
	for _, handler := range gcr.handlers {
		switch specificHandler := handler.(type) {
		case *groupContextRunnable:
			specificHandler.interactionStrategy.setParent(gcr.interactionStrategy)
		case *dialogContextRunnable:
			specificHandler.interactionStrategy.setParent(gcr.interactionStrategy)
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
	outputter(gcr.interactionStrategy.getFormattedHelpScreenTitle())
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
	gcr.interactionStrategy.run(func(line string) error {
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
	FullDescription      string
	OneLineDescription   string
	Name                 string
	Action               interface{}
	ReferenceValueGetter interface{}
}

var _ Handler = new(StructRunnerHandler)

func (srh *StructRunnerHandler) build() runnableHandler {
	actionInputType := srh.getAndCheckActionType()
	srh.checkReferenceValueGetter(actionInputType)
	result := &dialogContextRunnable{
		interactionStrategy: &interactionStrategyImpl{
			fullDescription:    srh.FullDescription,
			oneLineDescription: srh.OneLineDescription,
			name:               srh.Name,
			stopWords:          map[string]bool{CANCEL: true, CONTINUE: true},
		},
		action:               srh.Action,
		actionInpuType:       actionInputType,
		referenceValueGetter: srh.ReferenceValueGetter,
	}
	srh.addCommandHandlers(result)
	srh.addPropertyHandlers(actionInputType, result)
	return result
}

func (srh *StructRunnerHandler) getAndCheckActionType() reflect.Type {
	actionInputType := reflect.TypeOf(srh.Action).In(1).Elem()
	if actionInputType.Kind() != reflect.Struct {
		panic("The second argument of Action is expected to be a pointer to struct")
	}
	return actionInputType
}

func (srh *StructRunnerHandler) checkReferenceValueGetter(actionInputType reflect.Type) {
	if srh.ReferenceValueGetter != nil {
		referenceGetterType := reflect.TypeOf(srh.ReferenceValueGetter)
		if referenceGetterType.NumIn() != 0 || referenceGetterType.NumOut() != 1 {
			panic("Reference value getter should have no inputs and one output")
		}
		referenceType := referenceGetterType.Out(0).Elem()
		if actionInputType != referenceType {
			panic("The ReferenceValueGetter must produce the type needed by Action")
		}
	}
}

func (srh *StructRunnerHandler) addCommandHandlers(dcr *dialogContextRunnable) {
	helpHandler := &SingleLineHandler{
		Name:     HELP,
		Handler:  dcr.help,
		ArgNames: []string{},
	}
	reviewHandler := &SingleLineHandler{
		Name:     REVIEW,
		Handler:  dcr.review,
		ArgNames: []string{},
	}
	clearHandler := &SingleLineHandler{
		Name:     CLEAR,
		Handler:  dcr.clear,
		ArgNames: []string{},
	}
	continueHandler := &SingleLineHandler{
		Name:     CONTINUE,
		Handler:  dcr.doContinue,
		ArgNames: []string{},
	}
	cancelHandler := &SingleLineHandler{
		Name:     CANCEL,
		Handler:  dcr.cancel,
		ArgNames: []string{},
	}
	dcr.handlers = []runnableHandler{
		helpHandler.build(),
		reviewHandler.build(),
		clearHandler.build(),
		continueHandler.build(),
		cancelHandler.build(),
	}
}

func (srh *StructRunnerHandler) addPropertyHandlers(actionInputType reflect.Type, dcr *dialogContextRunnable) {
	dialogPropertyHandlers := make([]runnableHandler, actionInputType.NumField())
	for i := 0; i < actionInputType.NumField(); i++ {
		f := actionInputType.Field(i)
		if f.Name != strings.Title(f.Name) {
			panic("Field is not exported: " + f.Name)
		}
		dph := &dialogPropertyHandler{
			name:         util.UnTitle(f.Name),
			propertyType: f.Type,
			fieldNumber:  i,
		}
		dialogPropertyHandlers[i] = dph
	}
	dcr.handlers = append(dcr.handlers, dialogPropertyHandlers...)
}

type dialogContextRunnable struct {
	interactionStrategy  interactionStrategy
	handlers             []runnableHandler
	readValue            reflect.Value
	referenceValue       reflect.Value
	referenceValueGetter interface{}
	actionInpuType       reflect.Type
	action               interface{}
}

var _ runnableHandler = new(dialogContextRunnable)

func (dcr *dialogContextRunnable) getName() string {
	return dcr.interactionStrategy.getName()
}

func (dcr *dialogContextRunnable) help(outputter Outputter) {
	outputter(dcr.interactionStrategy.getFormattedHelpScreenTitle())
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
	table := StructToTable(dcr.readValue.Interface())
	outputter(table.String())
}

func (dcr *dialogContextRunnable) clear(outputter Outputter) {
	dcr.initReadValue()
}

func (dcr *dialogContextRunnable) doContinue(outputter Outputter) {
	actionValue := reflect.ValueOf(dcr.action)
	actionValue.Call([]reflect.Value{
		reflect.ValueOf(outputter),
		dcr.readValue,
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
	dcr.initReferenceValue()
	dcr.initReadValue()
	dcr.interactionStrategy.run(func(line string) error {
		return dcr.executeLine(line)
	})
}

func (dcr *dialogContextRunnable) initReferenceValue() {
	dcr.referenceValue = reflect.New(dcr.actionInpuType)
	if dcr.referenceValueGetter != nil {
		dcr.referenceValue = reflect.ValueOf(dcr.referenceValueGetter).Call([]reflect.Value{})[0]
	}
}

func (dcr *dialogContextRunnable) initReadValue() {
	dcr.readValue = reflect.New(dcr.actionInpuType)
	dcr.readValue.Elem().Set(dcr.referenceValue.Elem())
	for _, handler := range dcr.handlers {
		switch specificHandler := handler.(type) {
		case *dialogPropertyHandler:
			specificHandler.readValue = dcr.readValue
		}
	}
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
	name         string
	propertyType reflect.Type
	fieldNumber  int
	readValue    reflect.Value
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
			return errors.New(fmt.Sprintf("Type mismatch or value out of range for property %s: %s",
				dph.name, words[1]))
		}
	}
	dph.readValue.Elem().Field(dph.fieldNumber).Set(value)
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
			context.interactionStrategy.getName(),
			context.interactionStrategy.getOneLineDescription()))
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
			dialog.interactionStrategy.getName(),
			dialog.interactionStrategy.getOneLineDescription()))
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
		name:  "Properties that can be set",
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
