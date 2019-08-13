package manageManuscript

import (
	"gitlab.bbinfra.net/3estack/alexandria/portal/util"
	"html/template"
)

var uploadTemplate = `
{{define "upload"}}
    <form enctype="multipart/form-data" class="uploadForm" id="uploadControl">
        <label class="uploadForm__label" for="uploadTrigger">Upload or verify manuscript:</label>
        <input class="uploadForm__input" type="file" name="file" id="uploadTrigger">
    </form>
{{end}}
`
var manageManuscriptTemplate = `
{{define "manageManuscript"}}
<div class="manageDocument">
  {{template "upload"}}
  <div class="notification" id="alert"></div>
  <script src="https://unpkg.com/axios/dist/axios.min.js"></script>
  <script src="{{.JsUrl}}"></script>
  <script>
    console.log("Going to initialize the context");
    var context = {
      subjectId: {{.SubjectId}},
      initialIsUploadNeeded: {{.InitialIsUploadNeeded}},
      updateUrlComponent: {{.UpdateUrlComponent}},
      downloadControlId: {{.DownloadControlId}},
    };
    console.log("Initialized the context");
    linkManageDocument(context);
  </script>
</div>
{{end}}
`

func ParseManageManuscriptTemplate(name string) *template.Template {
	return util.ParseTemplates(name,
		uploadTemplate, manageManuscriptTemplate)
}

type ManageManuscriptContext struct {
	SubjectId             string
	InitialIsUploadNeeded bool
	JsUrl                 string
	UpdateUrlComponent    string
	DownloadControlId     string
}

type PortalManuscriptResponse struct {
	Message      string
	UploadNeeded bool
	IsWarning    bool
}
