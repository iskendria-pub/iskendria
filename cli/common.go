package cli

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

const (
	HELP        = "help"
	EXIT        = "exit"
	UNDO_FORMAT = "\033[0m"
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
		panic("Unsupported type")
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
