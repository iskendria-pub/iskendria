package main

import (
	"fmt"
	"gitlab.bbinfra.net/3estack/alexandria/cliAlexandria"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"html/template"
	"log"
	"net/http"
	"os"
)

var journalsTemplate = `
<head>
  <title>Alexandria</title>
  <link rel="stylesheet" href="/alexandria.css"/>
</head>
<body>
  <h1>Alexandria</h1>
  {{range .}}
  <div class="journal">
    <div class="title">{{.Title}}</div>
    <div class="editors">{{range .AcceptedEditors}}{{.PersonName}} {{end}}</div>
  </div>
  {{end}}
</body>
`

func main() {
	dbLogger := log.New(os.Stdout, "db", log.Flags())
	initialize(dbLogger)
	defer dao.Shutdown(dbLogger)
	runHttpServer()
}

func initialize(dbLogger *log.Logger) {
	dao.Init("portal.db", dbLogger)
	cliAlexandria.InitEventStream("./portal-events.log", "portal")
	go func() {
		for {
			_ = cliAlexandria.ReadEventStreamStatus()
		}
	}()
}

func runHttpServer() {
	http.HandleFunc("/index.html", handleJournals)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func handleJournals(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.New("journalsTemplate").Parse(journalsTemplate)
	if err != nil {
		fmt.Println("Error parsing template")
		fmt.Println(err)
		return
	}
	journals, err := dao.GetAllJournals()
	if err != nil {
		fmt.Println("Error reading journals from database: " + err.Error())
		return
	}
	err = tmpl.Execute(w, journals)
	if err != nil {
		fmt.Println(err)
	}
}
