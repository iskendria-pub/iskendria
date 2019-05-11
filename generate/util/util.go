package main

import (
	"fmt"
	"os"
	"text/template"
)

var templateUtil = `
package util

{{range .}}
func {{.FunctionName}}(sm map[string]{{.ValueType}}) []string {
	result := make([]string, 0)
	if sm == nil {
		return result
	}
	for s := range sm {
		result = append(result, s)
	}
	return result
}
{{end}}

`

type mapToSlice struct {
	FunctionName string
	ValueType    string
}

func main() {
	c := []mapToSlice{
		{
			FunctionName: "MapStringByteArrayToSlice",
			ValueType:    "[]byte",
		},
		{
			FunctionName: "MapStringBoolToSlice",
			ValueType:    "bool",
		},
	}
	tmpl, err := template.New("templateUtil").Parse(templateUtil)
	if err != nil {
		fmt.Println("Error parsing template")
		fmt.Println(err)
		return
	}
	err = tmpl.Execute(os.Stdout, c)
	if err != nil {
		fmt.Println(err)
	}
}