package manageDocument

import (
	"gitlab.bbinfra.net/3estack/alexandria/portal/util"
	"html/template"
)

var uploadTemplate = `
{{define "upload"}}
    <form enctype="multipart/form-data" class="uploadForm" id="uploadControl">
        <label class="uploadForm__label" for="uploadTrigger">Upload description file:</label>
        <input class="uploadForm__input" type="file" name="file" id="uploadTrigger">
    </form>
{{end}}
`
var verifyTemplate = `
{{define "verify"}}
  <table id="verifyControl">
    <tr>
      <td>Verify description:</td>
      <td><button type="button" id="verifyTrigger">Verify</button></td>
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
    var context = {
      descriptionHash: {{.DescriptionHash}},
      hasInitialDescriptionError: {{.HasInitialDescriptionError}},
      initialDescriptionError: {{.InitialDescriptionError}},
      initialIsUploadNeeded: {{.InitialIsUploadNeeded}},
      descriptionControlId: {{.DescriptionControlId}}
    }
    linkManageDocument(context)
  </script>
</div>
{{end}}
`

func ParseManageDocumentTemplate(name string) *template.Template {
	return util.ParseTemplates(name,
		uploadTemplate, verifyTemplate, manageDocumentTemplate)
}

type ManageDocumentContext struct {
	DescriptionHash            string
	HasInitialDescriptionError bool
	InitialDescriptionError    string
	InitialIsUploadNeeded      bool
	JsUrl                      string
	DescriptionControlId       string
}

type SaveFileSuccessResponse struct {
	Text    string
	Message string
}
