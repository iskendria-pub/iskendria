package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.bbinfra.net/3estack/alexandria/cliAlexandria"
	"gitlab.bbinfra.net/3estack/alexandria/dao"
	"gitlab.bbinfra.net/3estack/alexandria/model"
	"gitlab.bbinfra.net/3estack/alexandria/portal/components/manageDocument"
	"gitlab.bbinfra.net/3estack/alexandria/portal/components/manageManuscript"
	"gitlab.bbinfra.net/3estack/alexandria/portal/util"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const CONTENT_DISPOSITION = "Content-Disposition"

var editorsTemplate = `
{{- define "editors" -}}
{{- range $index, $element := . -}}
{{- if $index -}}, {{end -}}
<a href="/person/{{.PersonId}}" {{if not .PersonIsSigned}}class="muted"{{end}}>{{.PersonName}}</a>
{{- end -}}
{{- end -}}
`

var authorsTemplate = `
{{- define "authors" -}}
{{- range $index, $element := . -}}
{{- if $index -}}, {{end -}}
<a href="/person/{{.PersonId}}" {{if not .DidSign}}class="muted"{{end}}>{{.PersonName}}</a>
{{- end -}}
{{- end -}}
`

var journalsTemplate = `
{{- define "journalsTemplate" -}}
  {{range .}}
  <div class="journal">
    <div class="title"><a href="/journal/{{.JournalId}}" {{if not .IsSigned}}class="muted"{{end}}>{{.Title}}</a></div>
    <div class="editors">{{template "editors" .AcceptedEditors}}</div>
  </div>
  {{end}}
{{- end -}}
`

var journalsPageTemplate = `
<head>
  <title>Alexandria</title>
  <link rel="stylesheet" href="/public/alexandria.css"/>
</head>
<body>
  <h1>Alexandria</h1>
  {{template "journalsTemplate" .}} 
</body>
`

var volumesTemplate = `
{{define "volumes"}}
<table>
{{range .}}
<tr>
  <td><a href="/volume/{{.JournalId}}/{{.VolumeId}}">{{.Issue}}</a></td>
</tr>
{{end}}
</table>
{{end}}
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
  </table>
  <h2>Description</h2>
  <div id="descriptionId">{{.InitialDescription}}</div>
  <p>
  {{end}}
  {{template "manageDocument" .ManageDocument}}
  <h2>Volumes</h2>
  {{template "volumes" .Volumes}}
</body>
`

var cvTemplate = `
<head>
  <title>Alexandria</title>
  <link rel="stylesheet" href="/public/alexandria.css"/>
</head>
<body>
  {{with .PersonView}}
  <h1>Alexandria</h1>
  <h2>{{.Name}}</h2>
   <table>
     <tr>
       <td>PersonId:</td>
       <td>{{.Id}}</td> 
     </tr>
     <tr>
       <td>Public key:</td>
       <td>{{.PublicKey}}</td> 
     </tr>
     <tr>
       <td>Email:</td>
       <td>{{.Email}}</td> 
     </tr>
     <tr>
       <td>Is major:</td>
       <td>{{.IsMajor}}</td> 
     </tr>
     <tr>
       <td>Is signed:</td>
       <td>{{.IsSigned}}</td> 
     </tr>
     <tr>
       <td>Balance:</td>
       <td>{{.Balance}}</td> 
     </tr>
     <tr>
       <td>Organization:</td>
       <td>{{.Organization}}</td> 
     </tr>
     <tr>
       <td>Telephone:</td>
       <td>{{.Telephone}}</td> 
     </tr>
     <tr>
       <td>Address:</td>
       <td>{{.Address}}</td> 
     </tr>
     <tr>
       <td>PostalCode:</td>
       <td>{{.PostalCode}}</td> 
     </tr>
     <tr>
       <td>Country:</td>
       <td>{{.Country}}</td> 
     </tr>
     <tr>
       <td>Extra info:</td>
       <td>{{.ExtraInfo}}</td> 
     </tr>
  </table>
  <h2>Biography</h2>
  <div id="biographyId">{{.InitialBiography}}</div>
  <p>
  {{end}}
  {{template "manageDocument" .ManageDocument}}
  <h2>Editor of:</h2>
  {{template "journalsTemplate" .Journals}}
</body>
`

var volumeTemplate = `
<head>
  <title>Alexandria</title>
  <link rel="stylesheet" href="/public/alexandria.css"/>
</head>
<body>
  <h1>Alexandria</h1>
  <h2>{{.JournalTitle}}</h2>
  <h2>{{.Issue}}</h2>
</body>
`

var manuscriptTemplate = `
<head>
  <title>Alexandria</title>
  <link rel="stylesheet" href="/public/alexandria.css"/>
</head>
<body>
  {{with .Manuscript}}
  <h1>Alexandria</h1>
  <h2>{{.Title}}</h2>
  <div class="authors">{{template "authors" .Authors}}</div>
  <p/>
  <table>
    <tr>
      <td>Status:</td>
      <td>{{.Status}}</td>
    </tr>
    <tr>
      <td>Version number:</td>
      <td>{{.VersionNumber}}</td>
    </tr>
    <tr>
      <td>Commit message:</td>
      <td>{{.CommitMsg}}</td>
    </tr>
    <tr>
      <td>Id:</td>
      <td>{{.Id}}</td>
    </tr>
    <tr>
      <td>Thread id:</td>
      <td>{{.ThreadId}}</td>
    </tr>
  </table>
  <p>
  <form>
    <input type="button" value="Download" id="manuscriptDownload" onclick="window.location.href='/manuscriptDownload/{{.Id}}.pdf'" disabled/>
  </form> 
  <p>
  {{end}}
  <h2>Journal</h2>
  {{template "journalsTemplate" .Journals}}
  <h2>Manage</h2>
  {{template "manageManuscript" .ManageManuscript}}
</body>
`

const manageDocumentsJsUrl = "/manageDocument/manageDocument.js"
const manageManuscriptsJsUrl = "/manageManuscript/manageManuscript.js"

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
	r.HandleFunc("/journalUpdate/{journalId}", journalUpdate)
	r.HandleFunc("/journalVerifyAndRefresh/{journalId}", journalVerifyAndRefresh)
	r.HandleFunc("/person/{id}", handlePerson)
	r.HandleFunc("/personUpdate/{id}", personUpdate)
	r.HandleFunc("/personVerifyAndRefresh/{id}", personVerifyAndRefresh)
	r.HandleFunc("/volume/{journalId}/{volumeId}", handleVolume)
	r.HandleFunc("/manuscript/{manuscriptId}", handleManuscript)
	r.HandleFunc("/manuscriptUpdate/{id}", manuscriptUpdate)
	r.HandleFunc("/manuscriptDownload/{id}.pdf", handleManuscriptDownload)
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	r.PathPrefix("/manageDocument/").Handler(
		http.StripPrefix("/manageDocument/", http.FileServer(http.Dir("./components/manageDocument"))))
	r.PathPrefix("/manageManuscript/").Handler(
		http.StripPrefix("/manageManuscript/", http.FileServer(http.Dir("./components/manageManuscript"))))
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

var parsedJournalsPageTemplate = util.ParseTemplates("journalsPageTemplate",
	editorsTemplate, journalsTemplate, journalsPageTemplate)

func handleJournals(w http.ResponseWriter, _ *http.Request) {
	journals, err := dao.GetAllJournals()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(("Error reading journals from database: " + err.Error())))
		return
	}
	err = parsedJournalsPageTemplate.Execute(w, journals)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error executing template: " + err.Error()))
	}
}

var parsedJournalTemplate = parseTemplatesWithManageDocument(
	"journalTemplate", editorsTemplate, volumesTemplate, journalTemplate)

func parseTemplatesWithManageDocument(name string, templatesToAdd ...string) *template.Template {
	result := manageDocument.ParseManageDocumentTemplate(name)
	for _, toAdd := range templatesToAdd {
		var err error
		result, err = result.Parse(toAdd)
		if err != nil {
			panic(err)
		}
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
	volumes, err := dao.GetVolumesOfJournal(journalId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error reading volumes from database: " + err.Error()))
		return
	}
	err = parsedJournalTemplate.Execute(w, journalToJournalContext(journal, volumes))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error executing template: " + err.Error()))
	}
}

func journalToJournalContext(journal *dao.Journal, volumes []dao.Volume) *JournalContext {
	result := &JournalContext{
		JournalView: JournalView{
			JournalId:       journal.JournalId,
			Title:           journal.Title,
			AcceptedEditors: journal.AcceptedEditors,
		},
		ManageDocument: manageDocument.ManageDocumentContext{
			SubjectId:            journal.JournalId,
			JsUrl:                manageDocumentsJsUrl,
			DescriptionControlId: "descriptionId",
			UpdateUrlComponent:   "journalUpdate",
			VerifyUrlComponent:   "journalVerifyAndRefresh",
			SubjectWord:          "description",
		},
		Volumes: volumesToVolumeViews(volumes, journal.JournalId, journal.Title),
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
	Volumes        []*VolumeView
}

type JournalView struct {
	JournalId          string
	Title              string
	AcceptedEditors    []*dao.Editor
	InitialDescription string
}

type VolumeView struct {
	VolumeId     string
	Issue        string
	JournalId    string
	JournalTitle string
}

func volumesToVolumeViews(volumes []dao.Volume, journalId, journalTitle string) []*VolumeView {
	result := make([]*VolumeView, 0, len(volumes))
	for _, volume := range volumes {
		result = append(result, &VolumeView{
			VolumeId:     volume.VolumeId,
			Issue:        volume.Issue,
			JournalId:    journalId,
			JournalTitle: journalTitle,
		})
	}
	return result
}

func journalUpdate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering journalUpdate...\n")
	defer log.Printf("Leaving journalUpdate\n")
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
	if !checkNoJournalCorrectDescriptionOverwritten(journalId, theHash, w) {
		return
	}
	file, handle, err := r.FormFile("file")
	defer func() { _ = file.Close() }()
	saveFile(theHash, w, file, handle)
}

func checkNoJournalCorrectDescriptionOverwritten(
	journalId,
	theHash string,
	w http.ResponseWriter) bool {
	return checkNoDescriptionOverwritten(theHash, w, func(oldDescription []byte) error {
		return dao.VerifyJournalDescription(journalId, oldDescription)
	})
}

func checkNoDescriptionOverwritten(
	theHash string,
	w http.ResponseWriter,
	verifyDescriptionFunc func([]byte) error) bool {
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
	if hasOldDescription && verifyDescriptionFunc(oldDescription) == nil {
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

func journalVerifyAndRefresh(w http.ResponseWriter, r *http.Request) {
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
	log.Printf("Verification failed, setting up upload\n")
	journal, err := dao.GetJournal(journalId)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, fmt.Sprintf(
			"Journal not found: %s, detailed message is: %s", journalId, err))
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
		IsWarning:    true,
	})
}

func handlePerson(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering handlePerson...")
	defer log.Println("Left handlePerson")
	vars := mux.Vars(r)
	personId := vars["id"]
	cv, err := dao.GetCV(personId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error reading person cv from database: " + err.Error()))
		return
	}
	err = parsedPersonTemplate.Execute(w, personToCVContext(cv))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error executing template: " + err.Error()))
	}
}

var parsedPersonTemplate = parseTemplatesWithManageDocument("cvTemplate",
	editorsTemplate, journalsTemplate, cvTemplate)

func personToCVContext(cv *dao.CV) *CVContext {
	result := &CVContext{
		PersonView: PersonView{
			Id:           cv.Person.Id,
			PublicKey:    cv.Person.PublicKey,
			Name:         cv.Person.Name,
			Email:        cv.Person.Email,
			IsMajor:      cv.Person.IsMajor,
			IsSigned:     cv.Person.IsSigned,
			Balance:      cv.Person.Balance,
			Organization: cv.Person.Organization,
			Telephone:    cv.Person.Telephone,
			Address:      cv.Person.Address,
			PostalCode:   cv.Person.PostalCode,
			Country:      cv.Person.Country,
			ExtraInfo:    cv.Person.ExtraInfo,
		},
		ManageDocument: manageDocument.ManageDocumentContext{
			SubjectId:            cv.Person.Id,
			JsUrl:                manageDocumentsJsUrl,
			DescriptionControlId: "biographyId",
			UpdateUrlComponent:   "personUpdate",
			VerifyUrlComponent:   "personVerifyAndRefresh",
			SubjectWord:          "biography",
		},
		Journals: cv.Journals,
	}
	if cv.Person.BiographyHash == "" {
		return result
	}
	description, isAvailable, err := theDocuments.searchDescription(cv.Person.BiographyHash)
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
	result.PersonView.InitialBiography = string(description)
	return result
}

type CVContext struct {
	PersonView     PersonView
	Journals       []*dao.Journal
	ManageDocument manageDocument.ManageDocumentContext
}

type PersonView struct {
	Id               string
	PublicKey        string
	Name             string
	Email            string
	IsMajor          bool
	IsSigned         bool
	Balance          int32
	InitialBiography string
	Organization     string
	Telephone        string
	Address          string
	PostalCode       string
	Country          string
	ExtraInfo        string
}

func personUpdate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering personUpdate...\n")
	defer log.Printf("Leaving personUpdate\n")
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/index.html", http.StatusSeeOther)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	log.Printf("Uploading file for person id " + id)
	person, err := dao.GetPersonById(id)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, fmt.Sprintf("Person not found: %s", id))
		return
	}
	theHash := person.BiographyHash
	if !checkNoPersonCorrectDescriptionOverwritten(id, theHash, w) {
		return
	}
	file, handle, err := r.FormFile("file")
	defer func() { _ = file.Close() }()
	saveFile(theHash, w, file, handle)
}

func checkNoPersonCorrectDescriptionOverwritten(
	personId,
	theHash string,
	w http.ResponseWriter) bool {
	return checkNoDescriptionOverwritten(theHash, w, func(oldDescription []byte) error {
		return dao.VerifyPersonBiography(personId, oldDescription)
	})
}

func personVerifyAndRefresh(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering personVerifyAndRefresh...\n")
	defer log.Printf("Left personVerifyAndRefresh\n")
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/index.html", http.StatusSeeOther)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	log.Printf("Have id: %s\n", id)
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
	if err = dao.VerifyPersonBiography(id, []byte(request.Description)); err == nil {
		jsonSuccessResponse(w, &manageDocument.PortalResponse{
			Description: request.Description,
			Message:     "Verification successful, biography was correct",
		})
		return
	}
	log.Printf("Verification failed, setting up upload\n")
	person, err := dao.GetPersonById(id)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, fmt.Sprintf("Person not found: %s, detailed message is: %s",
			id, err))
		return
	}
	theHash := person.BiographyHash
	log.Printf("The hash is: " + theHash)
	updatedBiography, hasUpdatedBiography, err := theDocuments.searchDescription(theHash)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if hasUpdatedBiography {
		jsonSuccessResponse(w, &manageDocument.PortalResponse{
			Description: string(updatedBiography),
			Message:     "Updated the biography",
		})
		return
	}
	jsonSuccessResponse(w, &manageDocument.PortalResponse{
		Message:      "Please upload the biography",
		UploadNeeded: true,
		IsWarning:    true,
	})
}

func handleVolume(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering handleVolume...\n")
	defer log.Printf("Left handleVolume\n")
	vars := mux.Vars(r)
	journalId := vars["journalId"]
	volumeId := vars["volumeId"]
	journal, err := dao.GetJournal(journalId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Could not find journal for journalId: "+journalId)
		return
	}
	volume, err := dao.GetVolume(volumeId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Could not find volume for volumeId: "+volumeId)
		return
	}
	err = parsedVolumeTemplate.Execute(w, &VolumeView{
		VolumeId:     volume.VolumeId,
		Issue:        volume.Issue,
		JournalId:    journal.JournalId,
		JournalTitle: journal.Title,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Could not parse volume template: "+err.Error())
		return
	}
}

var parsedVolumeTemplate = util.ParseTemplates("volume", volumeTemplate)

func handleManuscript(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering handleManuscript...\n")
	defer log.Printf("Left handleManuscript\n")
	vars := mux.Vars(r)
	manuscriptId := vars["manuscriptId"]
	manuscript, err := dao.GetManuscriptView(manuscriptId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Could not find manuscript for manuscriptId"+manuscriptId)
		return
	}
	_, hasExistingManuscript, err := theDocuments.searchDescription(manuscript.Manuscript.Hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Could not check whether manuscript file is present")
		return
	}
	err = parsedManuscriptTemplate.Execute(w, manuscriptToManuscriptContext(manuscript, hasExistingManuscript))
	if err != nil {
		log.Printf("Could not parse manuscript template: " + err.Error())
		return
	}
}

var parsedManuscriptTemplate = parseManuscriptTemplate()

func parseManuscriptTemplate() *template.Template {
	result := manageManuscript.ParseManageManuscriptTemplate("manuscript")
	var err error
	extraTemplates := []string{authorsTemplate, editorsTemplate, journalsTemplate, manuscriptTemplate}
	for _, t := range extraTemplates {
		result, err = result.Parse(t)
		if err != nil {
			log.Printf("Could not parse template: %s, error: %s\n", t, err.Error())
			return nil
		}
	}
	return result
}

type ManuscriptContext struct {
	Manuscript       *dao.Manuscript
	ManageManuscript *manageManuscript.ManageManuscriptContext
	// This is a list for technical reason. In fact a manuscript
	// is only published in one journal.
	Journals []*dao.Journal
}

func manuscriptToManuscriptContext(manuscript *dao.ManuscriptView, hasExistingManuscript bool) *ManuscriptContext {
	return &ManuscriptContext{
		Manuscript: manuscript.Manuscript,
		Journals:   []*dao.Journal{manuscript.Journal},
		ManageManuscript: &manageManuscript.ManageManuscriptContext{
			SubjectId:             manuscript.Manuscript.Id,
			JsUrl:                 manageManuscriptsJsUrl,
			UpdateUrlComponent:    "manuscriptUpdate",
			DownloadControlId:     "manuscriptDownload",
			InitialIsUploadNeeded: !hasExistingManuscript,
		},
	}
}

func manuscriptUpdate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering manuscriptUpdate...\n")
	defer log.Printf("Leaving manuscriptUpdate\n")
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/index.html", http.StatusSeeOther)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	log.Printf("Uploading file and verifying for manuscript id " + id)
	manuscript, err := dao.GetManuscript(id)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, fmt.Sprintf("Manuscript not found: %s", id))
		return
	}
	expectedHash := manuscript.Hash
	file, _, err := r.FormFile("file")
	defer func() { _ = file.Close() }()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "Could not read uploaded file: "+err.Error())
		return
	}
	actualHash := model.HashBytes(data)
	if actualHash != expectedHash {
		jsonManuscriptSuccessResponse(w, &manageManuscript.PortalManuscriptResponse{
			Message:      "Verification failed",
			UploadNeeded: true,
			IsWarning:    true,
		})
		return
	}
	existingManuscriptData, alreadyPresent, err := theDocuments.searchDescription(expectedHash)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError,
			"Could not search for existing manuscript: "+err.Error())
		return
	}
	if alreadyPresent {
		hashOfStoredManuscript := model.HashBytes(existingManuscriptData)
		if hashOfStoredManuscript == expectedHash {
			jsonManuscriptSuccessResponse(w, &manageManuscript.PortalManuscriptResponse{
				Message:      "Verification successful",
				UploadNeeded: false,
				IsWarning:    false,
			})
			return
		}
	}
	err = theDocuments.save(actualHash, data)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError,
			"Could not save uploaded manuscript: "+err.Error())
		return
	}
	jsonManuscriptSuccessResponse(w, &manageManuscript.PortalManuscriptResponse{
		Message: "Manuscript uploaded successfully!",
	})
}

func jsonManuscriptSuccessResponse(w http.ResponseWriter, jsonMessage *manageManuscript.PortalManuscriptResponse) {
	body, err := json.Marshal(jsonMessage)
	if err != nil {
		panic(err)
	}
	jsonResponse(w, http.StatusOK, string(body))
}

func handleManuscriptDownload(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering handleManuscriptDownload...\n")
	defer log.Printf("Left handleManuscriptDownload\n")
	vars := mux.Vars(r)
	id := vars["id"]
	manuscript, err := dao.GetManuscript(id)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, "Unknown manuscript id: "+id)
		return
	}
	w.Header().Set(CONTENT_DISPOSITION, "attachment")
	in, err := theDocuments.open(manuscript.Hash)
	if err != nil {
		log.Printf("Could not open file to download: %s\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer func() { _ = in.Close() }()
	_, err = io.Copy(w, in)
	if err != nil {
		log.Printf("Could not copy downloaded file to response stream: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
