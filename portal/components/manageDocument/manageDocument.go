package manageDocument

import (
	"gitlab.bbinfra.net/3estack/alexandria/portal/util"
	"html/template"
)

var uploadTemplate = `
{{define "upload"}}
    <form enctype="multipart/form-data" class="uploadForm" id="uploadControl" hidden>
        <label class="uploadForm__label" for="uploadTrigger">Upload description file:</label>
        <input class="uploadForm__input" type="file" name="file" id="uploadTrigger">
    </form>
{{end}}
`
var verifyTemplate = `
{{define "verify"}}
  <table id="verifyControl" hidden>
    <tr>
      <td>Verify description:</td>
      <td><button type="button" id="verifyTrigger">Verify and refresh</button></td>
    <tr>
  </table>
{{end}}
`

var manageDocumentTemplate = `
{{define "manageDocument"}}
<div class="manageDocument">
  {{template "upload"}}
  {{template "verify"}}
  <div class="notification" id="alert"></div>
  <script src="https://unpkg.com/axios/dist/axios.min.js"></script>
  <script src="{{.JsUrl}}"></script>
  <script>
    console.log("Going to initialize the context");
    var context = {
      journalId: {{.JournalId}},
      hasInitialDescriptionError: {{.HasInitialDescriptionError}},
      initialDescriptionError: {{.InitialDescriptionError}},
      initialIsUploadNeeded: {{.InitialIsUploadNeeded}},
      descriptionControlId: {{.DescriptionControlId}},
      updateUrlComponent: {{.UpdateUrlComponent}},
      verifyUrlComponent: {{.VerifyUrlComponent}}
    };
    console.log("Initialized the context");
    linkManageDocument(context);
  </script>
</div>
{{end}}
`

func ParseManageDocumentTemplate(name string) *template.Template {
	return util.ParseTemplates(name,
		uploadTemplate, verifyTemplate, manageDocumentTemplate)
}

type ManageDocumentContext struct {
	JournalId                  string
	HasInitialDescriptionError bool
	InitialDescriptionError    string
	InitialIsUploadNeeded      bool
	JsUrl                      string
	DescriptionControlId       string
	UpdateUrlComponent         string
	VerifyUrlComponent         string
}

type PortalResponse struct {
	Description  string
	Message      string
	UploadNeeded bool
}

type PortalRequest struct {
	Description string
}
