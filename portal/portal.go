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
        <label class="uploadForm__label" for="inputFile">
            Not available, you can upload
        </label>
        <input class="uploadForm__input" type="file" name="file" id="inputFile">
    </form>
    <div class="notification" id="alert"></div>
    <script src="https://unpkg.com/axios/dist/axios.min.js"></script>
    <script src="/public/api.js"></script>
    <script src="/public/app.js"></script>
    <script>linkUploadForm({{.}}, document, axios)</script>
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
    <tr>
      <td>Description:</td>
      <td>{{if .IsUploadNeeded}}{{template "upload" .DescriptionHash}}{{else}}{{.Description}}{{end}}</td>
    </tr>
  </table>
  {{if .HasDescriptionError}}<div class="error">{{.DescriptionError}}<br/>{{end}}
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
	initDocuments()
}

func runHttpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/index.html", handleJournals)
	r.HandleFunc("/journal/{id}", handleJournal)
	r.HandleFunc("/upload/{theHash}", uploadFile)
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
	err = parsedJournalTemplate.Execute(w, journalToJournalView(journal))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error executing template: " + err.Error()))
	}
}

func journalToJournalView(journal *dao.Journal) *JournalView {
	result := &JournalView{
		JournalId:       journal.JournalId,
		Title:           journal.Title,
		DescriptionHash: journal.Descriptionhash,
		AcceptedEditors: journal.AcceptedEditors,
	}
	if journal.Descriptionhash == "" {
		return result
	}
	description, isAvailable, err := theDocuments.searchDescription(journal.Descriptionhash)
	if err != nil {
		result.IsUploadNeeded = true
		result.HasDescriptionError = true
		result.DescriptionError = err.Error()
		return result
	}
	if !isAvailable {
		result.IsUploadNeeded = true
		return result
	}
	result.Description = description
	return result
}

type JournalView struct {
	JournalId           string
	Title               string
	DescriptionHash     string
	AcceptedEditors     []*dao.Editor
	IsUploadNeeded      bool
	Description         string
	HasDescriptionError bool
	DescriptionError    string
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering uploadFile...\n")
	defer log.Printf("Leaving uploadFile\n")
	vars := mux.Vars(r)
	theHash := vars["theHash"]
	log.Printf("The hash is: " + theHash)
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
	saveFile(theHash, w, file, handle)
}

func saveFile(theHash string, w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader) {
	log.Printf("Entering saveFile...\n")
	defer log.Printf("Leaving saveFile\n")
	log.Printf("File is %s", handle.Filename)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		_, _ = fmt.Fprintf(w, "%v", err)
		return
	}
	err = theDocuments.Save(theHash, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "%v", err)
		return
	}
	jsonResponse(w, http.StatusCreated, "File uploaded successfully!")
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = fmt.Fprint(w, message)
}
