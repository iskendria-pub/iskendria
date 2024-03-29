package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/iskendria-pub/iskendria/cliIskendria"
	"github.com/iskendria-pub/iskendria/dao"
	"github.com/iskendria-pub/iskendria/model"
	"github.com/iskendria-pub/iskendria/portal/components/manageDocument"
	"github.com/iskendria-pub/iskendria/portal/components/manageManuscript"
	"github.com/iskendria-pub/iskendria/portal/util"
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
  <title>Iskendria</title>
  <link rel="stylesheet" href="/public/alexandria.css"/>
</head>
<body>
  <h1>Iskendria</h1>
  {{template "journalsTemplate" .}} 
</body>
`

var volumesTemplate = `
{{define "volumes"}}
<table>
{{range .}}
<tr>
  <td><a href="/volume/{{.VolumeId}}">{{.Issue}}</a></td>
</tr>
{{end}}
</table>
{{end}}
`

var journalTemplate = `
<head>
  <title>Iskendria</title>
  <link rel="stylesheet" href="/public/alexandria.css"/>
</head>
<body>
  {{with .JournalView}}
  <h1>Iskendria</h1>
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
  <table>
  <tr>
    <td><a href="/published/{{.JournalView.JournalId}}">Published, not yet assigned to volume</a></td>
  </tr>
  {{range .Volumes}}
  <tr>
    <td><a href="/volume/{{.VolumeId}}">{{.Issue}}</a></td>
  </tr>
  {{end}}
  </table>
</body>
`

var cvTemplate = `
<head>
  <title>Iskendria</title>
  <link rel="stylesheet" href="/public/alexandria.css"/>
</head>
<body>
  {{with .PersonView}}
  <h1>Iskendria</h1>
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
  <h2>Publications</h2>
  {{template "manuscriptsTemplate" .Manuscripts}}
</body>
`

var volumeTemplate = `
<head>
  <title>Iskendria</title>
  <link rel="stylesheet" href="/public/alexandria.css"/>
</head>
<body>
  <h1>Iskendria</h1>
  {{template "journalsTemplate" .Journal}}
  <h2>{{.Volume.Issue}}</h2>
  <table>
    <tr>
      <td>Volume id:</td>
      <td>{{.Volume.VolumeId}}</td>
    </tr>
  </table>
  <br />
  {{with .Manuscripts}}
  <table>
    {{- range . -}}
    <tr>
      <td>{{.FirstPage}} &hyphen; {{.LastPage}}</td>
      <td>&#x2005;</td>
      <td><a href="/manuscript/{{.Id}}">{{.Title}}</a>
        <div class="authors">{{template "authors" .Authors}}</div>
      </td>
    </tr>
    {{- end -}}
  </table>
  {{end}}
</body>
`

var publishedManuscriptsPageTemplate = `
<head>
  <title>Iskendria</title>
  <link rel="stylesheet" href="/public/alexandria.css"/>
</head>
<body>
  <h1>Iskendria</h1>
  {{template "journalsTemplate" .Journal}}
  <h2>Published, not yet assigned to volume</h2>
  {{with .Manuscripts}}
  <table>
    {{- range . -}}
    <tr>
      <td><a href="/manuscript/{{.Id}}">{{.Title}}</a>
        <div class="authors">{{template "authors" .Authors}}</div>
      </td>
    </tr>
    {{- end -}}
  </table>
  {{end}}
</body>
`

var manuscriptTemplate = `
<head>
  <title>Iskendria</title>
  <link rel="stylesheet" href="/public/alexandria.css"/>
</head>
<body>
  {{with .Manuscript}}
  <h1>Iskendria</h1>
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
    <tr>
      <td>First page</td>
      <td>{{.FirstPage}}</td>
    </tr>
    <tr>
      <td>Last page</td>
      <td>{{.LastPage}}</td>
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
  {{if .Volumes}}
  {{template "volumes" .Volumes}}
  {{else}}
  Not assigned to volume.
  {{end}}
  <h2>Reviews</h2>
  {{template "reviewList" .Reviews}}
  <h2>Manage</h2>
  {{template "manageManuscript" .ManageManuscript}}
</body>
`

var manuscriptsTemplate = `
{{- define "manuscriptsTemplate" -}}
  {{range .}}
  <div class="manuscript">
    <div class="title"><a href="/manuscript/{{.Id}}">{{.Title}}</a></div>
    <div class="authors">{{template "authors" .Authors}}</div>
  </div>
  {{end}}
{{- end -}}
`

var reviewPageTemplate = `
<head>
  <title>Iskendria</title>
  <link rel="stylesheet" href="/public/alexandria.css"/>
</head>
<body>
  <h1>Iskendria</h1>
  <h2>Subject of review</h2>
  {{template "manuscriptsTemplate" .Manuscripts}}
  <h2>Review</h2>
  {{with .Review}}
  <table>
    <tr>
      <td>Id:</td>
      <td>{{.Id}}
    </tr>
    <tr>
      <td>Review author:</td>
      {{- end -}}
      {{- with .ReviewAuthor -}}
      <td><a href="/person/{{.PersonId}}" {{if not .PersonIsSigned}}class="muted"{{end}}>{{.PersonName}}</a></td>
      {{- end -}}
    </tr>
    <tr>
      {{- with .Review -}}
      <td>Judgement:</td>
      <td>{{.Judgement}}</td>
    </tr>
    <tr>
      <td>Used by editor:</td>
      <td>{{.IsUsedByEditor}}
    </tr>
  </table>
  {{- end -}}
  <h2>Review text</h2>
  <div id="reviewTextId">{{.ReviewText}}</div>
  <p>
  {{template "manageDocument" .ManageDocument}}
</body>
`

var reviewListTemplate = `
{{define "reviewList"}}
{{range .}}
<table>
<tr><td>
<a href="/person/{{.PersonId}}" {{if not .PersonIsSigned}}class="muted"{{end}}>{{.PersonName}}</a>
</td></tr>
<tr><td>
<div id="reviewTextId">{{.ReviewText}}</div>
</td></tr>
<tr><td>
{{.Judgement}}
</td></tr>
<tr><td>
{{if .IsUsedByEditor}}Used by editor to judge{{end}}
</td></tr>
<tr><td>
<a href="/review/{{.Id}}">Manage</a>
</tr></td>
</table>
<p>
{{end}}
{{end}}
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
	cliIskendria.InitEventStream("./portal-events.log", "portal")
	go func() {
		for {
			_ = cliIskendria.ReadEventStreamStatus()
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
	r.HandleFunc("/volume/{volumeId}", handleVolume)
	r.HandleFunc("/published/{journalId}", handlePublished)
	r.HandleFunc("/manuscript/{manuscriptId}", handleManuscript)
	r.HandleFunc("/manuscriptUpdate/{id}", manuscriptUpdate)
	r.HandleFunc("/manuscriptDownload/{id}.pdf", handleManuscriptDownload)
	r.HandleFunc("/review/{id}", handleReviewDetail)
	r.HandleFunc("/reviewUpdate/{id}", reviewUpdate)
	r.HandleFunc("/reviewVerifyAndRefresh/{id}", reviewVerifyAndRefresh)
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
		log.Printf("Parsing added template: %s\n", toAdd)
		var err error
		result, err = result.Parse(toAdd)
		if err != nil {
			fmt.Println("Error parsing " + toAdd)
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
	editorsTemplate, journalsTemplate, authorsTemplate, manuscriptsTemplate, cvTemplate)

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
		Journals:    cv.Journals,
		Manuscripts: cv.Manuscripts,
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
	Manuscripts    []*dao.Manuscript
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
	volumeId := vars["volumeId"]
	volume, err := dao.GetVolumeView(volumeId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Could not find volume for volumeId: "+volumeId)
		return
	}
	err = parsedVolumeTemplate.Execute(w, volumeToVolumeContext(volume))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Could not parse volume template: "+err.Error())
		return
	}
}

func volumeToVolumeContext(volume *dao.VolumeView) *VolumeContext {
	return &VolumeContext{
		Volume: &VolumeView{
			VolumeId:     volume.Volume.VolumeId,
			Issue:        volume.Volume.Issue,
			JournalId:    volume.Journal.JournalId,
			JournalTitle: volume.Journal.Title,
		},
		Journal:     []*dao.Journal{volume.Journal},
		Manuscripts: volume.Manuscripts,
	}
}

type VolumeContext struct {
	Volume      *VolumeView
	Journal     []*dao.Journal
	Manuscripts []*dao.Manuscript
}

var parsedVolumeTemplate = util.ParseTemplates("volume",
	editorsTemplate, journalsTemplate, authorsTemplate, volumeTemplate)

func handlePublished(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering handlePublished...\n")
	defer log.Printf("Left handlePublished\n")
	vars := mux.Vars(r)
	journalId := vars["journalId"]
	published, err := dao.GetPublishedManuscriptView(journalId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Could not find journal for journalId: "+journalId)
		return
	}
	err = parsedPublishedTemplate.Execute(w, publishedToPublishedContext(published))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Could not parse volume template: "+err.Error())
		return
	}
}

func publishedToPublishedContext(published *dao.PublishedManuscriptView) *PublishedContext {
	return &PublishedContext{
		Journal:     []*dao.Journal{published.Journal},
		Manuscripts: published.Manuscripts,
	}
}

type PublishedContext struct {
	Journal     []*dao.Journal
	Manuscripts []*dao.Manuscript
}

var parsedPublishedTemplate = util.ParseTemplates("published",
	editorsTemplate, journalsTemplate, authorsTemplate, publishedManuscriptsPageTemplate)

func handleManuscript(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering handleManuscript...\n")
	defer log.Printf("Left handleManuscript\n")
	vars := mux.Vars(r)
	manuscriptId := vars["manuscriptId"]
	manuscript, err := dao.GetManuscriptView(manuscriptId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(
			w,
			"Could not find manuscript for manuscriptId %s, error %s",
			manuscriptId,
			err.Error())
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
	extraTemplates := []string{
		authorsTemplate,
		editorsTemplate,
		journalsTemplate,
		reviewListTemplate,
		volumesTemplate,
		manuscriptTemplate}
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
	Reviews  []*ReviewListItem
	Volumes  []*VolumeView
}

type ReviewListItem struct {
	PersonId       string
	PersonIsSigned bool
	PersonName     string
	Id             string
	ReviewText     string
	Judgement      string
	IsUsedByEditor bool
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
		Reviews: extendedReviewsToReviewListItems(manuscript.Reviews),
		Volumes: getManuscriptVolumes(
			manuscript.Volume,
			manuscript.Journal.JournalId,
			manuscript.Journal.Title),
	}
}

func extendedReviewsToReviewListItems(source []*dao.ExtendedReview) []*ReviewListItem {
	result := make([]*ReviewListItem, len(source))
	for i, s := range source {
		reviewText, _, err := theDocuments.searchDescription(s.Hash)
		if err != nil {
			reviewText = []byte("ERROR getting review text: " + err.Error())
		}
		result[i] = &ReviewListItem{
			PersonId:       s.PersonId,
			PersonIsSigned: s.PersonIsSigned,
			PersonName:     s.PersonName,
			Id:             s.Id,
			ReviewText:     string(reviewText),
			Judgement:      s.Judgement,
			IsUsedByEditor: s.IsUsedByEditor,
		}
	}
	return result
}

func getManuscriptVolumes(volume *dao.Volume, journalId, journalTitle string) []*VolumeView {
	if volume == nil {
		return []*VolumeView{}
	}
	return volumesToVolumeViews(
		[]dao.Volume{*volume},
		journalId,
		journalTitle)
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

func handleReviewDetail(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering handleReviewDetail...\n")
	defer log.Printf("Left handleReviewDetail\n")
	vars := mux.Vars(r)
	reviewId := vars["id"]
	review, err := dao.GetReviewDetailsView(reviewId)
	if err != nil {
		log.Printf("Could not get review for id: %s, error: %s\n",
			reviewId, err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	reviewText, hasReviewText, err := theDocuments.searchDescription(review.Review.Hash)
	if err != nil {
		log.Printf("Could not search review text: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = parsedReviewDetailsTemplate.Execute(w, reviewToReviewDetailsView(review, hasReviewText, string(reviewText)))
	if err != nil {
		log.Printf("Could not parse manuscript template: " + err.Error())
		return
	}
}

type ReviewDetailsView struct {
	Manuscripts    []*dao.Manuscript
	Review         *dao.Review
	ReviewAuthor   *dao.ReviewEditor
	ReviewText     string
	ManageDocument *manageDocument.ManageDocumentContext
}

func reviewToReviewDetailsView(reviewView *dao.ReviewView, hasReviewText bool, reviewText string) *ReviewDetailsView {
	return &ReviewDetailsView{
		Manuscripts: []*dao.Manuscript{
			reviewView.Manuscript,
		},
		Review: reviewView.Review,
		ReviewAuthor: &dao.ReviewEditor{
			PersonId:       reviewView.ReviewAuthor.Id,
			PersonName:     reviewView.ReviewAuthor.Name,
			PersonIsSigned: reviewView.ReviewAuthor.IsSigned,
		},
		ReviewText: reviewText,
		ManageDocument: &manageDocument.ManageDocumentContext{
			SubjectId:             reviewView.Review.Id,
			InitialIsUploadNeeded: !hasReviewText,
			JsUrl:                 manageDocumentsJsUrl,
			DescriptionControlId:  "reviewTextId",
			UpdateUrlComponent:    "reviewUpdate",
			VerifyUrlComponent:    "reviewVerifyAndRefresh",
			SubjectWord:           "review text",
		},
	}
}

func reviewUpdate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering reviewUpdate...\n")
	defer log.Printf("Leaving reviewUpdate\n")
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/index.html", http.StatusSeeOther)
		return
	}
	vars := mux.Vars(r)
	reviewId := vars["id"]
	log.Printf("Uploading file for review id " + reviewId)
	review, err := dao.GetReview(reviewId)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, fmt.Sprintf("Journal not found: %s", reviewId))
		return
	}
	theHash := review.Hash
	if !checkNoReviewCorrectDescriptionOverwritten(reviewId, theHash, w) {
		return
	}
	file, handle, err := r.FormFile("file")
	defer func() { _ = file.Close() }()
	saveFile(theHash, w, file, handle)
}

func checkNoReviewCorrectDescriptionOverwritten(
	reviewId,
	theHash string,
	w http.ResponseWriter) bool {
	return checkNoDescriptionOverwritten(theHash, w, func(oldDescription []byte) error {
		return dao.VerifyReview(reviewId, oldDescription)
	})
}

func reviewVerifyAndRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/index.html", http.StatusSeeOther)
		return
	}
	vars := mux.Vars(r)
	reviewId := vars["id"]
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
	if err = dao.VerifyReview(reviewId, []byte(request.Description)); err == nil {
		jsonSuccessResponse(w, &manageDocument.PortalResponse{
			Description: request.Description,
			Message:     "Verification successful, description was correct",
		})
		return
	}
	log.Printf("Verification failed, setting up upload\n")
	review, err := dao.GetReview(reviewId)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, fmt.Sprintf(
			"Review not found: %s, detailed message is: %s", reviewId, err))
		return
	}
	theHash := review.Hash
	log.Printf("The hash is: " + theHash)
	updatedDescription, hasUpdatedDescription, err := theDocuments.searchDescription(theHash)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if hasUpdatedDescription {
		jsonSuccessResponse(w, &manageDocument.PortalResponse{
			Description: string(updatedDescription),
			Message:     "Updated the review text",
		})
		return
	}
	jsonSuccessResponse(w, &manageDocument.PortalResponse{
		Message:      "Please upload the review",
		UploadNeeded: true,
		IsWarning:    true,
	})
}

var parsedReviewDetailsTemplate = parseTemplatesWithManageDocument(
	"reviewDetailsTemplate", authorsTemplate, manuscriptsTemplate, reviewPageTemplate)
