package cli

import (
	"errors"
	"fmt"
	"github.com/iskendria-pub/iskendria/util"
	"reflect"
	"sort"
	"strings"
)

type runnableHandler interface {
	getName() string
	getHelpIndexLine() string
	handleLine(words []string) error
	isArgumentsAfterEqualSign() bool
}

type singleLineRunnableHandler struct {
	name     string
	handler  interface{}
	argNames []string
}

var _ runnableHandler = new(singleLineRunnableHandler)

func (slrh *singleLineRunnableHandler) getName() string {
	return slrh.name
}

func (slrh *singleLineRunnableHandler) getHelpIndexLine() string {
	return combineNameAndArgs(slrh.name, slrh.argNames)
}

func combineNameAndArgs(name string, argNames []string) string {
	items := []string{name}
	if argNames != nil {
		for _, arg := range argNames {
			items = append(items, formatArgument(arg))
		}
	}
	return strings.Join(items, " ")
}

func formatArgument(arg string) string {
	return "<" + arg + ">"
}

func (slrh *singleLineRunnableHandler) handleLine(words []string) error {
	expectedNumArgs := len(slrh.argNames)
	actualNumArgs := len(words) - 1
	if expectedNumArgs != actualNumArgs {
		return invalidNumberOfArguments(expectedNumArgs, actualNumArgs)
	}
	switch reflect.TypeOf(slrh.handler).NumOut() {
	case 0:
		return slrh.handleLineUsingOutputter(words[1:])
	case 1:
		return slrh.handleLineUsingReturnedError(words[1:])
	}
	return nil
}

func invalidNumberOfArguments(expectedNumArgs int, actualNumArgs int) error {
	return errors.New(fmt.Sprintf(
		"Invalid number of arguments: expected %d, got %d", expectedNumArgs, actualNumArgs))
}

func (slrh *singleLineRunnableHandler) handleLineUsingOutputter(argWords []string) error {
	callResultValue, err := callIncludingOutputter(
		slrh.handler,
		slrh.argNames,
		argWords)
	if err != nil {
		return err
	}
	if len(callResultValue) != 0 {
		panic("Did not expect a result")
	}
	return nil
}

func callIncludingOutputter(handler interface{}, argNames, argWords []string) ([]reflect.Value, error) {
	values, err := getValues(
		reflect.TypeOf(handler),
		argNames,
		argWords,
		func(i int) int { return i + 1 })
	if err != nil {
		return nil, err
	}
	values = append([]reflect.Value{reflect.ValueOf(outputToStdout)}, values...)
	callResultValue := reflect.ValueOf(handler).Call(values)
	return callResultValue, nil
}

func (slrh *singleLineRunnableHandler) handleLineUsingReturnedError(argWords []string) error {
	values, err := getValues(reflect.TypeOf(slrh.handler), slrh.argNames, argWords, func(i int) int { return i })
	if err != nil {
		return err
	}
	callResultValue := reflect.ValueOf(slrh.handler).Call(values)
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

func getValues(handlerType reflect.Type, argNames []string, argWords []string, namedArgIdxToArgIdx func(int) int) (
	[]reflect.Value, error) {
	numArgWords := len(argNames)
	values := make([]reflect.Value, numArgWords)
	for index, word := range argWords {
		argumentType := handlerType.In(namedArgIdxToArgIdx(index))
		value, err := getValue(word, argumentType)
		if err != nil {
			return values, err
		}
		values[index] = value
	}
	return values, nil
}

func (slrh *singleLineRunnableHandler) isArgumentsAfterEqualSign() bool {
	return false
}

type groupContextRunnable struct {
	*handlersForGroup
	interactionStrategy interactionStrategy
}

var _ runnableHandler = new(groupContextRunnable)

func (gcr *groupContextRunnable) getName() string {
	return gcr.interactionStrategy.getName()
}

func (gcr *groupContextRunnable) getHelpIndexLine() string {
	return fmt.Sprintf("%s - %s",
		gcr.interactionStrategy.getName(),
		gcr.interactionStrategy.getOneLineDescription())
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

func (gcr *groupContextRunnable) help(outputter Outputter) {
	outputter(gcr.interactionStrategy.getFormattedHelpScreenTitle())
	gcr.showHandlers(outputter)
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
		return errors.New("entering a group does not require arguments")
	}
	gcr.run()
	return nil
}

func (grc *groupContextRunnable) isArgumentsAfterEqualSign() bool {
	return false
}

type dialogContextRunnable struct {
	*handlersForDialog
	interactionStrategy          interactionStrategy
	action                       interface{}
	actionInpuType               reflect.Type
	referenceValueGetter         interface{}
	referenceValueGetterArgNames []string
	readValue                    reflect.Value
	referenceValue               reflect.Value
}

var _ runnableHandler = new(dialogContextRunnable)

func (dcr *dialogContextRunnable) getName() string {
	return dcr.interactionStrategy.getName()
}

func (dcr *dialogContextRunnable) getHelpIndexLine() string {
	return fmt.Sprintf("%s - %s",
		combineNameAndArgs(dcr.interactionStrategy.getName(), dcr.referenceValueGetterArgNames),
		dcr.interactionStrategy.getOneLineDescription())
}

func (dcr *dialogContextRunnable) help(outputter Outputter) {
	outputter(dcr.interactionStrategy.getFormattedHelpScreenTitle())
	dcr.showHandlers(outputter)
}

func (dcr *dialogContextRunnable) review(outputter Outputter) {
	numFields := dcr.readValue.Elem().NumField()
	reviewValues := make([]*reviewValue, numFields)
	for i := 0; i < numFields; i++ {
		reviewValues[i] = newReviewValue(
			util.UnTitle(dcr.readValue.Elem().Type().Field(i).Name),
			dcr.readValue.Elem().Field(i))
	}
	numLines := 0
	for i := 0; i < len(reviewValues); i++ {
		numLines += reviewValues[i].numLines
	}
	table := NewTable(numLines, 2)
	lineNum := 0
	for i := 0; i < len(reviewValues); i++ {
		if len(reviewValues[i].formattedValues) == 0 {
			table.Set(lineNum, 0, reviewValues[i].fieldName)
			lineNum++
		} else {
			for j := range reviewValues[i].formattedValues {
				if j == 0 {
					table.Set(lineNum, 0, reviewValues[i].fieldName)
				}
				table.Set(lineNum, 1, reviewValues[i].formattedValues[j])
				lineNum++
			}
		}
	}
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

func (dcr *dialogContextRunnable) add(outputter Outputter, variableName, addedItem string) {
	value, err := dcr.initListManipulation(variableName, addedItem)
	if err != nil {
		outputter(err.Error() + "\n")
		return
	}
	theSlice := make([]string, value.Len()+1)
	for i := 0; i < value.Len(); i++ {
		theSlice[i] = value.Index(i).Interface().(string)
	}
	theSlice[value.Len()] = addedItem
	value.Set(reflect.ValueOf(theSlice))
}

func (dcr *dialogContextRunnable) initListManipulation(variableName, item string) (reflect.Value, error) {
	if item == "" {
		return reflect.ValueOf(nil), errors.New("Cannot add empty value")
	}
	fieldName := strings.Title(variableName)
	value, found := dcr.getReflectValueFromStruct(fieldName)
	if !found {
		return reflect.ValueOf(nil), errors.New("Field does not exist: " + variableName)
	}
	return value, nil
}

func (dcr *dialogContextRunnable) getReflectValueFromStruct(fieldName string) (reflect.Value, bool) {
	for i := 0; i < dcr.readValue.Elem().NumField(); i++ {
		if dcr.readValue.Elem().Type().Field(i).Name == fieldName {
			return dcr.readValue.Elem().Field(i), true
		}
	}
	return reflect.ValueOf(nil), false
}

func (dcr *dialogContextRunnable) removeItem(outputter Outputter, variableName, removedItem string) {
	value, err := dcr.initListManipulation(variableName, removedItem)
	if err != nil {
		outputter(err.Error() + "\n")
		return
	}
	didRemove := false
	theSlice := make([]string, 0, value.Len())
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i).Interface().(string)
		if item == removedItem {
			didRemove = true
		} else {
			theSlice = append(theSlice, item)
		}
	}
	if !didRemove {
		outputter("Item to be removed was not present\n")
		return
	}
	value.Set(reflect.ValueOf(theSlice))
}

func (dcr *dialogContextRunnable) insert(outputter Outputter, variableName, insertedItem string, oneBasedIndex int32) {
	value, err := dcr.initListManipulation(variableName, insertedItem)
	if err != nil {
		outputter(err.Error() + "\n")
		return
	}
	index := oneBasedIndex - 1
	if index < 0 || index >= int32((value.Len()+1)) {
		outputter(fmt.Sprintf("Index out of range: %d\n", oneBasedIndex))
		return
	}
	offset := int32(0)
	theSlice := make([]string, value.Len()+1)
	for i := int32(0); i < int32(value.Len()); i++ {
		if i == index {
			theSlice[i] = insertedItem
			offset = 1
		}
		theSlice[i+offset] = value.Index(int(i)).Interface().(string)
	}
	if offset == 0 {
		theSlice[value.Len()] = insertedItem
	}
	value.Set(reflect.ValueOf(theSlice))
}

func (dcr *dialogContextRunnable) removeIndex(outputter Outputter, variableName string, oneBasedIndex int32) {
	value, found := dcr.getReflectValueFromStruct(strings.Title(variableName))
	if !found {
		outputter(fmt.Sprintf("Field does not exist: %s\n", variableName))
		return
	}
	index := oneBasedIndex - 1
	if index < 0 || index >= int32(value.Len()) {
		outputter(fmt.Sprintf("Index out of bounds: %d\n", oneBasedIndex))
		return
	}
	theSlice := make([]string, value.Len()-1)
	offset := int32(0)
	for i := int32(0); i < int32(value.Len()); i++ {
		if i == index {
			offset = -1
		} else {
			theSlice[i+offset] = value.Index(int(i)).Interface().(string)
		}
	}
	value.Set(reflect.ValueOf(theSlice))
}

func (dcr *dialogContextRunnable) handleLine(words []string) error {
	err, stop := dcr.setReferenceValue(words)
	if err != nil {
		return err
	}
	if stop {
		return nil
	}
	dcr.initReadValue()
	dcr.interactionStrategy.run(func(line string) error {
		return dcr.executeLine(line)
	})
	return nil
}

func (dcr *dialogContextRunnable) setReferenceValue(words []string) (error, bool) {
	dcr.referenceValue = reflect.New(dcr.actionInpuType)
	if dcr.referenceValueGetter != nil {
		expectedNumArgs := len(dcr.referenceValueGetterArgNames)
		actualNumArgs := len(words) - 1
		if actualNumArgs != expectedNumArgs {
			return invalidNumberOfArguments(expectedNumArgs, actualNumArgs), true
		}
		callResult, err := callIncludingOutputter(
			dcr.referenceValueGetter, dcr.referenceValueGetterArgNames, words[1:])
		if err != nil {
			return err, true
		}
		dcr.referenceValue = callResult[0]
		if dcr.referenceValue.IsNil() {
			return nil, true
		}
	}
	return nil, false
}

func (dcr *dialogContextRunnable) initReadValue() {
	dcr.readValue = reflect.New(dcr.actionInpuType)
	dcr.readValue.Elem().Set(dcr.referenceValue.Elem())
	for _, handler := range dcr.handlers {
		switch specificHandler := handler.(type) {
		case *dialogPropertyHandler:
			specificHandler.readValue = dcr.readValue
		case *dialogListPropertyHandler:
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
		return dcr.searchAndExecuteHandlerWithoutEqualSign(line)
	}
	return handler.handleLine(words)
}

func (dcr *dialogContextRunnable) searchAndExecuteHandlerWithoutEqualSign(line string) error {
	words := strings.Split(strings.TrimSpace(line), " ")
	handler := getHandler(words[0], dcr.handlers)
	// Prevent property to be set without = sign
	if handler == nil {
		return errors.New("Command not found: " + words[0])
	}
	if handler.isArgumentsAfterEqualSign() {
		return errors.New("Property set requires \"=\" sign: " + words[0])
	}
	return handler.handleLine(words)
}

func (dcr *dialogContextRunnable) isArgumentsAfterEqualSign() bool {
	return false
}

type reviewValue struct {
	fieldName       string
	numLines        int
	formattedValues []string
}

func newReviewValue(fieldName string, value reflect.Value) *reviewValue {
	result := &reviewValue{
		fieldName: fieldName,
	}
	if value.Kind() == reflect.Slice {
		theSlice := make([]string, value.Len())
		for i := 0; i < value.Len(); i++ {
			theSlice[i] = value.Index(i).Interface().(string)
		}
		result.formattedValues = theSlice
		if len(result.formattedValues) == 0 {
			result.numLines = 1
		} else {
			result.numLines = len(result.formattedValues)
		}
		return result
	}
	result.numLines = 1
	result.formattedValues = []string{fmt.Sprintf("%v", value.Interface())}
	return result
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

func (dph *dialogPropertyHandler) getHelpIndexLine() string {
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

func (dph *dialogPropertyHandler) isArgumentsAfterEqualSign() bool {
	return true
}

type dialogListPropertyHandler struct {
	name        string
	fieldNumber int
	readValue   reflect.Value
}

var _ runnableHandler = new(dialogListPropertyHandler)

func (dlph *dialogListPropertyHandler) getName() string {
	return dlph.name
}

func (dlph *dialogListPropertyHandler) getHelpIndexLine() string {
	return dlph.name + " (list)"
}

func (dlph *dialogListPropertyHandler) handleLine(words []string) error {
	value := []string{}
	if len(words) == 2 && words[1] != "" {
		items := strings.Split(strings.TrimSpace(words[1]), " ")
		value = make([]string, 0, len(items))
		for _, item := range items {
			if item != "" {
				value = append(value, item)
			}
		}
	}
	dlph.readValue.Elem().Field(dlph.fieldNumber).Set(reflect.ValueOf(value))
	return nil
}

func getHandler(name string, handlers []runnableHandler) runnableHandler {
	for _, handler := range handlers {
		if handler.getName() == name {
			return handler
		}
	}
	return nil
}

func (dlph *dialogListPropertyHandler) isArgumentsAfterEqualSign() bool {
	return true
}
