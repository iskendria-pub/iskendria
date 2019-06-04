package util

import (
	"fmt"
	"html/template"
	"log"
)

func ParseTemplates(name string, parsedItems ...string) *template.Template {
	log.Printf("Entering ParseTemplates for name: %s\n", name)
	tmpl := template.New(name)
	var err error
	for _, parsedItem := range parsedItems {
		log.Printf("Parsing called template: %s\n", parsedItem)
		tmpl, err = tmpl.Parse(parsedItem)
		if err != nil {
			fmt.Println("Error parsing " + parsedItem)
			panic(err)
		}
	}
	return tmpl
}
