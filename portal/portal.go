package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.bbinfra.net/3estack/alexandria/cliAlexandria"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"html/template"
	"log"
	"net/http"
	"os"
)

var editorsTemplate = `
{{define "editors"}}
{{range .}}<a href="/person/{{.PersonId}}">{{.PersonName}}</a>{{end}}
{{end}}
`

var journalsTemplate = `
<head>
  <title>Alexandria</title>
  <link rel="stylesheet" href="/public/alexandria.css"/>
</head>
<body>
  <h1>Alexandria</h1>
  {{range .}}
  <div class="journal">
    <div class="title"><a href="/journal/{{.JournalId}}">{{.Title}}</a></div>
    <div class="editors">{{template "editors" .AcceptedEditors}}</div>
  </div>
  {{end}}
</body>
`

var journalTemplate = `
<head>
  <title>Alexandria</title>
  <link rel="stylesheet" href="/public/alexandria.css"/>
</head>
<body>
  <h1>Alexandria</h1>
  <h2>{{.Title}}</h2>
  <div>{{template "editors" .AcceptedEditors}}</div>
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
	r := mux.NewRouter()
	r.HandleFunc("/index.html", handleJournals)
	r.HandleFunc("/journal/{id}", handleJournal)
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func handleJournals(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.New("journalsTemplate").Parse(editorsTemplate)
	if err != nil {
		fmt.Println("Error parsing editorsTemplate")
		fmt.Println(err)
		return
	}
	tmpl, err = tmpl.Parse(journalsTemplate)
	if err != nil {
		fmt.Println("Error parsing journalsTemplate")
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

func handleJournal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	journalId := vars["id"]
	tmpl, err := template.New("journalsTemplate").Parse(editorsTemplate)
	if err != nil {
		fmt.Println("Error parsing editorsTemplate")
		fmt.Println(err)
		return
	}
	tmpl, err = tmpl.Parse(journalTemplate)
	if err != nil {
		fmt.Println("Error parsing journalTemplate")
		fmt.Println(err)
		return
	}
	journal, err := dao.GetJournal(journalId)
	if err != nil {
		fmt.Println("Error reading journal from database: " + err.Error())
		return
	}
	err = tmpl.Execute(w, journal)
	if err != nil {
		fmt.Println(err)
	}
}
