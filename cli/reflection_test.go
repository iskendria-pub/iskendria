package cli

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReflection(t *testing.T) {
	d := &DialogStruct{}
	var i interface{} = d
	vp := reflect.ValueOf(i)
	fmt.Printf("Value vp = %v\n", vp)
	v := vp.Elem()
	fmt.Printf("Value v = %v\n", v)
	fmt.Printf("F1 = %v\n", v.Field(0).Interface())
}

func TestReflectionWithNew(t *testing.T) {
	d := DialogStruct{}
	dt := reflect.TypeOf(d)
	newVp := reflect.New(dt)
	newV := reflect.Indirect(newVp)
	fmt.Printf("newV.F0 = %v\n", newV.Field(0).Interface())
}

func TestReflectionWithFunc(t *testing.T) {
	var i interface{} = df
	st := reflect.TypeOf(i).In(1).Elem()
	fmt.Printf("Type st = %v\n", st)
	newV := reflect.New(st)
	var newVI interface{} = newV
	newVIV := reflect.ValueOf(newVI).Elem()
	// v := newV.Elem()
	fmt.Printf("newV.F0 = %v\n", newVIV.Field(0).Interface())
}

/*
type DialogStruct struct {
	F1 bool
	F2 int32
	F3 string
}
*/

func df(int, *DialogStruct) {}
