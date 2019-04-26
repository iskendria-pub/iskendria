package cli

import (
	"errors"
	"fmt"
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
		panic("Handler is not a function")
	}
	expectedNumFunctionArgs := len(slh.ArgNames) + 1
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

func lessRunnableHandler(first, second runnableHandler) bool {
	return first.getName() < second.getName()
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
		default:
			panic("Unexpected handler type")
		}
	}
}

func (gcr *groupContextRunnable) addRunnableHandler(handler runnableHandler) {
	gcr.handlers = append(gcr.handlers, handler)
}

func (gcr *groupContextRunnable) help(outputter Outputter) {
	outputter(gcr.helpStrategy.getFormattedHelpScreenTitle())
	gcr.showHandlers(outputter)
}

func (gcr *groupContextRunnable) showHandlers(outputter Outputter) {
	groups, commands := gcr.splitGroupsAndCommands()
	sort.Slice(groups, func(i, j int) bool {
		return lessRunnableHandler(groups[i], groups[j])
	})
	sort.Slice(commands, func(i, j int) bool {
		return lessRunnableHandler(groups[i], groups[j])
	})
	outputter(gcr.combineGroupsAndCommands(commands, groups))
}

func (gcr *groupContextRunnable) splitGroupsAndCommands() ([]*groupContextRunnable, []*SingleLineHandler) {
	groups := make([]*groupContextRunnable, 0)
	commands := make([]*SingleLineHandler, 0)
	for _, handler := range gcr.handlers {
		switch specificHandler := handler.(type) {
		case *groupContextRunnable:
			groups = append(groups, specificHandler)
		case *SingleLineHandler:
			commands = append(commands, specificHandler)
		default:
			panic("Unknown handler type")
		}
	}
	return groups, commands
}

func (gcr *groupContextRunnable) combineGroupsAndCommands(
	commands []*SingleLineHandler, groups []*groupContextRunnable) string {
	parts := make([]string, 2)
	parts[0] = gcr.listCommands(commands)
	parts[1] = gcr.listGroups(groups)
	filledParts := make([]string, 0)
	for _, part := range parts {
		if part != "" {
			filledParts = append(filledParts, part)
		}
	}
	return strings.Join(filledParts, "\n") + "\n"
}

func (gcr *groupContextRunnable) listCommands(handlers []*SingleLineHandler) string {
	if len(handlers) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("Commands:\n")
	for _, handler := range handlers {
		items := []string{handler.Name}
		for _, arg := range handler.ArgNames {
			items = append(items, "<"+arg+">")
		}
		sb.WriteString("  " + strings.Join(items, " ") + "\n")
	}
	return sb.String()
}

func (gcr *groupContextRunnable) listGroups(contexts []*groupContextRunnable) string {
	if len(contexts) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("Groups\n")
	for _, context := range contexts {
		sb.WriteString(fmt.Sprintf("  %s - %s\n",
			context.helpStrategy.getName(),
			context.helpStrategy.getOneLineDescription()))
	}
	return sb.String()
}

func (gcr *groupContextRunnable) run() {
	gcr.helpStrategy.run(func(line string) error {
		return gcr.executeLine(line)
	})
}

func (gcr *groupContextRunnable) executeLine(line string) error {
	words := strings.Split(strings.TrimSpace(line), " ")
	handler := gcr.getHandler(words[0])
	if handler == nil {
		return errors.New("Name does not match a group or a command: " + words[0])
	}
	return handler.handleLine(words)
}

func (gcr *groupContextRunnable) getHandler(name string) runnableHandler {
	for _, handler := range gcr.handlers {
		if handler.getName() == name {
			return handler
		}
	}
	return nil
}

func (group *groupContextRunnable) handleLine(words []string) error {
	if len(words) != 1 {
		return errors.New("Entering a group does not require arguments")
	}
	group.run()
	return nil
}
