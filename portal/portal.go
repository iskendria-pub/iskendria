package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.bbinfra.net/3estack/alexandria/cliAlexandria"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"gitlab.bbinfra.net/3estack/alexandria/portal/components/manageDocument"
	"gitlab.bbinfra.net/3estack/alexandria/portal/util"
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

var journalTemplate = `
<head>
  <title>Alexandria</title>
  <link rel="stylesheet" href="/public/alexandria.css"/>
</head>
<body>
  {{with .JournalView}}
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
      <td id="descriptionId">{{.InitialDescription}}</td>
    </tr>
  </table>
  {{end}}
  {{template "manageDocument" .ManageDocument}}
</body>
`

const manageDocumentsJsUrl = "/manageDocument/manageDocument.js"

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
	r.HandleFunc("/verify/{theHash}", verify)
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	r.PathPrefix("/manageDocument/").Handler(
		http.StripPrefix("/manageDocument/", http.FileServer(http.Dir("./components/manageDocument"))))
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

var parsedJournalsTemplate = util.ParseTemplates("journalsTemplate", editorsTemplate, journalsTemplate)

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

var parsedJournalTemplate = parseJournalTemplate()

func parseJournalTemplate() *template.Template {
	result := manageDocument.ParseManageDocumentTemplate("journalTemplate")
	result, err := result.Parse(editorsTemplate)
	if err != nil {
		panic(err)
	}
	result, err = result.Parse(journalTemplate)
	if err != nil {
		panic(err)
	}
	return result
}

func handleJournal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	journalId := vars["id"]
	journal, err := dao.GetJournal(journalId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error reading journal from database: " + err.Error()))
		return
	}
	err = parsedJournalTemplate.Execute(w, journalToJournalContext(journal))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error executing template: " + err.Error()))
	}
}

func journalToJournalContext(journal *dao.Journal) *JournalContext {
	result := &JournalContext{
		JournalView: JournalView{
			JournalId:       journal.JournalId,
			Title:           journal.Title,
			AcceptedEditors: journal.AcceptedEditors,
		},
		ManageDocument: manageDocument.ManageDocumentContext{
			DescriptionHash:      journal.Descriptionhash,
			JsUrl:                manageDocumentsJsUrl,
			DescriptionControlId: "descriptionId",
		},
	}
	if journal.Descriptionhash == "" {
		return result
	}
	description, isAvailable, err := theDocuments.searchDescription(journal.Descriptionhash)
	if err != nil {
		result.ManageDocument.InitialIsUploadNeeded = true
		result.ManageDocument.HasInitialDescriptionError = true
		result.ManageDocument.InitialDescriptionError = err.Error()
		return result
	}
	if !isAvailable {
		result.ManageDocument.InitialIsUploadNeeded = true
		return result
	}
	result.JournalView.InitialDescription = description
	return result
}

type JournalContext struct {
	JournalView    JournalView
	ManageDocument manageDocument.ManageDocumentContext
}

type JournalView struct {
	JournalId          string
	Title              string
	AcceptedEditors    []*dao.Editor
	InitialDescription string
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
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "%v", err)
		return
	}
	err = theDocuments.Save(theHash, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "%v", err)
		return
	}
	result, err := json.Marshal(&manageDocument.SaveFileSuccessResponse{
		Text:    string(data),
		Message: "File uploaded successfully!",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "%v", err)
		return
	}
	jsonResponse(w, http.StatusCreated, string(result))
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = fmt.Fprint(w, message)
}

func verify(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	theHash := vars["theHash"]
	log.Printf("The hash is: " + theHash)
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/index.html", http.StatusSeeOther)
		return
	}
	if theHash == "" {
		jsonResponse(w, http.StatusOK, "No description on blockchain, nothing to verify")
		return
	}
	_, isAvailable, err := theDocuments.searchDescription(theHash)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, "Verification failed")
		return
	}
	if !isAvailable {
		jsonResponse(w, http.StatusNotFound, "Description was not present, you have to upload it again")
		return
	}
	jsonResponse(w, http.StatusOK, "Verified!")
}
