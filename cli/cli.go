package cli

import (
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/util"
	"reflect"
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

There is one more feature to describe. The user may have
background processes. These may have messages to report
to the user. This package provides a hook that is called
each time the end user presses Enter. This way, the user
interface provides frequent opportunities to show messages,
but messages do not interfere when the end user is
typing a command. Please set the EventPager property
of struct Cli.

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

var _ Handler = new(SingleLineHandler)

type SingleLineHandler struct {
	Name     string
	Handler  interface{}
	ArgNames []string
}

func (slh *SingleLineHandler) build() runnableHandler {
	checkFunctionWithArgNames(slh.Handler, slh.ArgNames)
	return &singleLineRunnableHandler{
		name:     slh.Name,
		handler:  slh.Handler,
		argNames: slh.ArgNames,
	}
}

func checkFunctionWithArgNames(f interface{}, argNames []string) {
	reflectHandlerType := reflect.TypeOf(f)
	if reflectHandlerType.Kind() != reflect.Func {
		panic("Handler is not a function")
	}
	switch reflectHandlerType.NumOut() {
	case 0:
		checkInputTypesIncludingOutputter(reflectHandlerType, argNames)
	case 1:
		checkFunctionThatReturnsError(reflectHandlerType, argNames)
	default:
		panic("Handler should return zero or one value")
	}
}

func checkInputTypesIncludingOutputter(reflectHandlerType reflect.Type, argNames []string) {
	expectedNumFunctionArgs := len(argNames) + 1
	if reflectHandlerType.NumIn() != expectedNumFunctionArgs {
		panic(fmt.Sprintf(
			"Number of handler arguments does not match number of argument names or outputter function is missing: %v",
			reflectHandlerType))
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

func checkFunctionThatReturnsError(reflectHandlerType reflect.Type, argNames []string) {
	if reflectHandlerType.NumIn() != len(argNames) {
		panic("Number of handler arguments does not match number of argument names")
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

var _ Handler = new(Cli)

type Cli struct {
	FullDescription    string
	OneLineDescription string
	Name               string
	FormatEscape       string
	Handlers           []Handler
	EventPager         func(Outputter)
}

func (c *Cli) Run() {
	if InputScript == "" {
		inp = new(inputSourceConsole)
	} else {
		inp = &inputSourceFile{
			fname: InputScript,
		}
	}
	inp.open()
	defer inp.close()
	c.buildMain().run()
}

var InputScript string

func (c *Cli) build() runnableHandler {
	return runnableHandler(c.buildMain())
}

func (c *Cli) buildMain() *groupContextRunnable {
	runnableHandlers := make([]runnableHandler, len(c.Handlers))
	for i := range c.Handlers {
		runnableHandlers[i] = c.Handlers[i].build()
	}
	result := &groupContextRunnable{
		handlersForGroup: &handlersForGroup{
			runnableHandlers,
		},
		interactionStrategy: &interactionStrategyImpl{
			fullDescription:    c.FullDescription,
			oneLineDescription: c.OneLineDescription,
			name:               c.Name,
			formatEscape:       c.FormatEscape,
			stopWords:          map[string]bool{EXIT: true},
			eventPager:         c.EventPager,
		},
	}
	c.addGeneratedCommandHandlers(result)
	result.init()
	return result
}

func (c *Cli) addGeneratedCommandHandlers(gcr *groupContextRunnable) {
	var helpHandler = &SingleLineHandler{
		Name:     HELP,
		Handler:  func(outputter Outputter) { gcr.help(outputter) },
		ArgNames: []string{},
	}
	var exitHandler = &SingleLineHandler{
		Name:     EXIT,
		Handler:  func(Outputter) {},
		ArgNames: []string{},
	}
	gcr.handlers = append(gcr.handlers, helpHandler.build())
	gcr.handlers = append(gcr.handlers, exitHandler.build())
}

type StructRunnerHandler struct {
	FullDescription              string
	OneLineDescription           string
	Name                         string
	Action                       interface{}
	ReferenceValueGetter         interface{}
	ReferenceValueGetterArgNames []string
}

var _ Handler = new(StructRunnerHandler)

func (srh *StructRunnerHandler) build() runnableHandler {
	actionInputType := srh.getAndCheckActionType()
	srh.checkReferenceValueGetter(actionInputType)
	result := &dialogContextRunnable{
		handlersForDialog: new(handlersForDialog),
		interactionStrategy: &interactionStrategyImpl{
			fullDescription:    srh.FullDescription,
			oneLineDescription: srh.OneLineDescription,
			name:               srh.Name,
			stopWords:          map[string]bool{CANCEL: true, CONTINUE: true},
		},
		action:                       srh.Action,
		actionInpuType:               actionInputType,
		referenceValueGetter:         srh.ReferenceValueGetter,
		referenceValueGetterArgNames: srh.ReferenceValueGetterArgNames,
	}
	srh.addGeneratedCommandHandlers(result)
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
		if referenceGetterType.NumIn() == 0 || referenceGetterType.NumOut() != 1 {
			panic(fmt.Sprintf("Reference value getter should have at least one input and one output: %v",
				referenceGetterType))
		}
		checkInputTypesIncludingOutputter(referenceGetterType, srh.ReferenceValueGetterArgNames)
		referenceType := referenceGetterType.Out(0).Elem()
		if actionInputType != referenceType {
			panic("The ReferenceValueGetter must produce the type needed by Action")
		}
	}
}

func (srh *StructRunnerHandler) addGeneratedCommandHandlers(dcr *dialogContextRunnable) {
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
