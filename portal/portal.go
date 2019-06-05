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
	r.HandleFunc("/update/{journalId}", uploadFile)
	r.HandleFunc("/verifyAndRefresh/{journalId}", verifyAndRefresh)
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
			JournalId:            journal.JournalId,
			JsUrl:                manageDocumentsJsUrl,
			DescriptionControlId: "descriptionId",
			UpdateUrlComponent:   "update",
			VerifyUrlComponent:   "verifyAndRefresh",
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
	result.JournalView.InitialDescription = string(description)
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
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/index.html", http.StatusSeeOther)
		return
	}
	vars := mux.Vars(r)
	journalId := vars["journalId"]
	log.Printf("Uploading file for journal id " + journalId)
	journal, err := dao.GetJournal(journalId)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, fmt.Sprintf("Journal not found: %s", journalId))
		return
	}
	theHash := journal.Descriptionhash
	if !checkNoCorrectDescriptionOverwritten(journalId, theHash, w) {
		return
	}
	file, handle, err := r.FormFile("file")
	defer func() { _ = file.Close() }()
	saveFile(theHash, w, file, handle)
}

func checkNoCorrectDescriptionOverwritten(
	journalId,
	theHash string,
	w http.ResponseWriter) bool {
	log.Printf("The hash of the description is: " + theHash)
	var hasOldDescription = true
	var oldDescription = []byte{}
	var err error
	if theHash != "" {
		oldDescription, hasOldDescription, err = theDocuments.searchDescription(theHash)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, err.Error())
			return false
		}
	}
	if hasOldDescription && dao.VerifyJournalDescription(journalId, oldDescription) == nil {
		jsonResponse(w, http.StatusForbidden, "The correct description was present already, do not overwrite")
		return false
	}
	return true
}

func saveFile(theHash string, w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader) {
	log.Printf("Entering saveFile...\n")
	defer log.Printf("Leaving saveFile\n")
	log.Printf("File is %s", handle.Filename)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = theDocuments.save(theHash, data)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonSuccessResponse(w, &manageDocument.PortalResponse{
		Description: string(data),
		Message:     "File uploaded successfully!",
	})
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = fmt.Fprint(w, message)
}

func jsonSuccessResponse(w http.ResponseWriter, jsonMessage *manageDocument.PortalResponse) {
	body, err := json.Marshal(jsonMessage)
	if err != nil {
		panic(err)
	}
	jsonResponse(w, http.StatusOK, string(body))
}

func verifyAndRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/index.html", http.StatusSeeOther)
		return
	}
	vars := mux.Vars(r)
	journalId := vars["journalId"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, fmt.Sprintf(
			"Could not read body of POST request: %s", err.Error()))
		return
	}
	defer func() { _ = r.Body.Close() }()
	request := &manageDocument.PortalRequest{}
	err = json.Unmarshal(body, request)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, fmt.Sprintf(
			"Could not parse body of POST request as PortalRequest: %s", err.Error()))
		return
	}
	if err = dao.VerifyJournalDescription(journalId, []byte(request.Description)); err == nil {
		jsonSuccessResponse(w, &manageDocument.PortalResponse{
			Description: request.Description,
			Message:     "Verification successful, description was correct",
		})
		return
	}
	journal, err := dao.GetJournal(journalId)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, fmt.Sprintf("Journal not found: %s, detailed message is: %s",
			journalId, err))
		return
	}
	theHash := journal.Descriptionhash
	log.Printf("The hash is: " + theHash)
	updatedDescription, hasUpdatedDescription, err := theDocuments.searchDescription(theHash)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if hasUpdatedDescription {
		jsonSuccessResponse(w, &manageDocument.PortalResponse{
			Description: string(updatedDescription),
			Message:     "Updated the description",
		})
		return
	}
	jsonSuccessResponse(w, &manageDocument.PortalResponse{
		Message:      "Please upload the description",
		UploadNeeded: true,
	})
}
