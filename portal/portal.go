package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.bbinfra.net/3estack/alexandria/cliAlexandria"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"html/template"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

var editorsTemplate = `
{{- define "editors" -}}
{{- range $index, $element := . -}}
{{- if $index -}}, {{end -}}
<a href="/person/{{.PersonId}}">{{.PersonName}}</a>
{{- end -}}
{{- end -}}
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

var uploadTemplate = `
{{define "upload"}}
    <form action="/upload" method="post" enctype="multipart/form-data" class="uploadForm">
        <input class="uploadForm__input" type="file" name="file" id="inputFile">
        <label class="uploadForm__label" for="inputFile">
            Select a file
        </label>
    </form>
    <div class="notification" id="alert"></div>
    <script src="https://unpkg.com/axios/dist/axios.min.js"></script>
    <script src="/public/api.js"></script>
    <script src="/public/app.js"></script>
{{end}}
`

var journalTemplate = `
<head>
  <title>Alexandria</title>
  <link rel="stylesheet" href="/public/alexandria.css"/>
</head>
<body>
  <h1>Alexandria</h1>
  <h2>{{.Title}}</h2>
  <table>
    <tr>
      <td>Id:</td>
      <td>{{.JournalId}}</td>
    </tr>
    <tr>
      <td>Editors:</td>
      <td><div>{{template "editors" .AcceptedEditors}}</div></td>
    </tr>
  </table>
  {{template "upload"}}
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
	r.HandleFunc("/upload", UploadFile)
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func parseTemplates(name string, parsedItems ...string) *template.Template {
	tmpl := template.New("journalsTemplate")
	var err error
	for _, parsedItem := range parsedItems {
		tmpl, err = tmpl.Parse(parsedItem)
		if err != nil {
			fmt.Println("Error parsing " + parsedItem)
			panic(err)
		}
	}
	return tmpl
}

var parsedJournalsTemplate = parseTemplates("journalsTemplate", editorsTemplate, journalsTemplate)

func handleJournals(w http.ResponseWriter, _ *http.Request) {
	journals, err := dao.GetAllJournals()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(("Error reading journals from database: " + err.Error())))
		return
	}
	err = parsedJournalsTemplate.Execute(w, journals)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error executing template: " + err.Error()))
	}
}

var parsedJournalTemplate = parseTemplates("journalTemplate", editorsTemplate, uploadTemplate, journalTemplate)

func handleJournal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	journalId := vars["id"]
	journal, err := dao.GetJournal(journalId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error reading journal from database: " + err.Error()))
		return
	}
	err = parsedJournalTemplate.Execute(w, journal)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error executing template: " + err.Error()))
	}
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering UploadFile...\n")
	defer log.Printf("Leaving UploadFile\n")
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/index.html", http.StatusSeeOther)
		return
	}

	file, handle, err := r.FormFile("file")
	if err != nil {
		_, _ = fmt.Fprintf(w, "%v", err)
		return
	}
	defer func() { _ = file.Close() }()

	/*
		mimeType := handle.Header.Get("Content-Type")
		switch mimeType {
		case "image/jpeg":
			saveFile(w, file, handle)
		case "image/png":
			saveFile(w, file, handle)
		default:
			jsonResponse(w, http.StatusBadRequest, "The format file is not valid.")
		}
	*/
	saveFile(w, file, handle)
}

func saveFile(w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader) {
	log.Printf("Entering saveFile...\n")
	defer log.Printf("Leaving saveFile\n")
	log.Printf("File is %s", handle.Filename)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		_, _ = fmt.Fprintf(w, "%v", err)
		return
	}

	err = ioutil.WriteFile("./uploaded", data, 0666)
	if err != nil {
		_, _ = fmt.Fprintf(w, "%v", err)
		return
	}
	jsonResponse(w, http.StatusCreated, "File uploaded successfully!.")
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = fmt.Fprint(w, message)
}
