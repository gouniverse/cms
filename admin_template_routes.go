package cms

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	hb "github.com/gouniverse/html"
	"github.com/gouniverse/utils"
)

func pageTemplatesTemplateCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(utils.Req(r, "name", ""), " ")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	entity := GetEntityStore().EntityCreateWithAttributes("template", map[string]interface{}{
		"name":   name,
		"status": "inactive",
	})

	log.Println(entity)

	if entity == nil {
		api.Respond(w, r, api.Error("Template failed to be created"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Template saved successfully", map[string]interface{}{"template_id": entity.ID}))
	return
}

func pageTemplateTemplateManager(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	log.Println(endpoint)

	header := cmsHeader(endpoint)
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathTemplatesTemplateManager): "Templates",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "template-manager")
	heading := hb.NewHeading1().HTML("Template Manager")
	button := hb.NewButton().HTML("New template").Attr("class", "btn btn-success float-end").Attr("v-on:click", "showTemplateCreateModal")
	heading.AddChild(button)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))

	modal := hb.NewDiv().Attr("id", "ModalTemplateCreate").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("New Template"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewDiv().Attr("class", "form-group").AddChild(hb.NewLabel().HTML("Name")).AddChild(hb.NewInput().Attr("class", "form-control").Attr("v-model", "templateCreateModel.name")))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-prsecondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Create & Continue").Attr("class", "btn btn-primary").Attr("v-on:click", "templateCreate"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	container.AddChild(modal)

	templates := GetEntityStore().EntityList("template", 0, 200, "", "id", "asc")

	table := hb.NewTable().Attr("class", "table table-responsive table-striped mt-3")
	thead := hb.NewThead()
	tbody := hb.NewTbody()
	table.AddChild(thead).AddChild(tbody)

	tr := hb.NewTR()
	th1 := hb.NewTD().HTML("Name")
	th2 := hb.NewTD().HTML("Status")
	th3 := hb.NewTD().HTML("Actions").Attr("style", "width:1px;")
	thead.AddChild(tr.AddChild(th1).AddChild(th2).AddChild(th3))

	for _, template := range templates {
		name := template.GetString("name", "n/a")
		status := template.GetString("status", "n/a")
		buttonEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Attr("class", "btn btn-primary").Attr("v-on:click", "templateEdit('"+template.ID+"')")

		tr := hb.NewTR()
		td1 := hb.NewTD().HTML(name)
		td2 := hb.NewTD().HTML(status)
		td3 := hb.NewTD().AddChild(buttonEdit)

		tbody.AddChild(tr.AddChild(td1).AddChild(td2).AddChild(td3))
	}
	container.AddChild(table)

	h := container.ToHTML()

	inlineScript := `
var templateCreateUrl = "` + endpoint + `?path=templates/template-create-ajax"
var templateUpdateUrl = "` + endpoint + `?path=templates/template-update"
const TemplateManager = {
	data() {
		return {
		  templateCreateModel:{
			  name:""
		  }
		}
	},
	methods: {
        showTemplateCreateModal(){
			//alert("Create template");
			var modalTemplateCreate = new bootstrap.Modal(document.getElementById('ModalTemplateCreate'));
			modalTemplateCreate.show();
		},
		templateCreate(){
			var name = this.templateCreateModel.name;
			$.post(templateCreateUrl, {name: name}).done((result)=>{
				if (result.status==="success"){
					var modalTemplateCreate = new bootstrap.Modal(document.getElementById('ModalTemplateCreate'));
				        modalTemplateCreate.hide();
					return location.href = templateUpdateUrl+ "&template_id=" + result.data.template_id;
				}
				alert("Failed. " + result.message)
			}).fail((result)=>{
				alert("Failed" + result)
			});
		},
		templateEdit(templateId){
			return location.href = templateUpdateUrl+ "&template_id=" + templateId;
		}
	}
};
Vue.createApp(TemplateManager).mount('#template-manager')
	`

	webtemplate := Webpage("Template Manager", h)
	webtemplate.AddScript(inlineScript)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webtemplate.ToHTML()))
}

// pageTemplatesTemplateUpdate shows the template edit page
func pageTemplatesTemplateUpdate(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	log.Println(endpoint)

	templateID := utils.Req(r, "template_id", "")
	if templateID == "" {
		api.Respond(w, r, api.Error("Template ID is required"))
		return
	}

	template := GetEntityStore().EntityFindByID(templateID)

	if template == nil {
		api.Respond(w, r, api.Error("Template NOT FOUND with ID "+templateID))
		return
	}

	header := cmsHeader(r.Context().Value(keyEndpoint).(string))
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathTemplatesTemplateManager):                               "Templates",
		(endpoint + "?path=" + PathTemplatesTemplateUpdate + "&template_id=" + templateID): "Edit template",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "template-update")
	heading := hb.NewHeading1().HTML("Edit Template")
	button := hb.NewButton().HTML("Save").Attr("class", "btn btn-success float-end").Attr("v-on:click", "templateSave")
	heading.AddChild(button)

	formGroupStatus := hb.NewDiv().Attr("class", "form-group mb-3")
	formGroupStatusLabel := hb.NewLabel().HTML("Status").Attr("class", "form-label")
	formGroupStatusSelect := hb.NewSelect().Attr("class", "form-control").Attr("v-model", "templateModel.status")
	formGroupOptionsActive := hb.NewOption().Attr("value", "active").HTML("Active")
	formGroupOptionsInactive := hb.NewOption().Attr("value", "inactive").HTML("Inactive")
	formGroupOptionsTrash := hb.NewOption().Attr("value", "trash").HTML("Trash")
	formGroupStatus.AddChild(formGroupStatusLabel)
	formGroupStatus.AddChild(formGroupStatusSelect.AddChild(formGroupOptionsActive).AddChild(formGroupOptionsInactive).AddChild(formGroupOptionsTrash))

	formGroupName := hb.NewDiv().Attr("class", "form-group mb-3")
	formGroupNameLabel := hb.NewLabel().HTML("Name").Attr("class", "form-label")
	formGroupNameInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "templateModel.name")
	formGroupName.AddChild(formGroupNameLabel)
	formGroupName.AddChild(formGroupNameInput)

	formGroupContent := hb.NewDiv().Attr("class", "form-group mb-3")
	formGroupContentLabel := hb.NewLabel().HTML("Content").Attr("class", "form-label")
	formGroupContentInfo := hb.NewParagraph().HTML("Available variables: [[PageContent]], [[PageCanonicalUrl]], [[PageMetaDescription]], [[PageMetaKeywords]], [[PageMetaRobots]], [[PageTitle]]").Attr("class", "alert alert-info text-small")
	formGroupContentInput := hb.NewTextArea().Attr("class", "form-control CodeMirror").Attr("v-model", "templateModel.content")
	formGroupContent.AddChild(formGroupContentLabel)
	formGroupContent.AddChild(formGroupContentInfo)
	formGroupContent.AddChild(formGroupContentInput)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))
	container.AddChild(formGroupStatus).AddChild(formGroupName).AddChild(formGroupContent)

	h := container.ToHTML()

	name := template.GetString("name", "")
	status := template.GetString("status", "")
	content := template.GetString("content", "")
	templateJSON, _ := json.Marshal(templateID)
	nameJSON, _ := json.Marshal(name)
	statusJSON, _ := json.Marshal(status)
	contentJSON, _ := json.Marshal(content)

	inlineScript := `
var templateUpdateUrl = "` + endpoint + `?path=templates/template-update-ajax";
var templateId = ` + string(templateJSON) + `;
var name = ` + string(nameJSON) + `;
var status = ` + string(statusJSON) + `;
var content = ` + string(contentJSON) + `;
const TemplateUpdate = {
	data() {
		return {
			templateModel:{
				templateId: templateId,
				content: content,
				name: name,
				status: status,
		    }
		}
	},
	created(){
		var self = this;
		setTimeout(function () {
			if ($('.CodeMirror').length > 0) {
				var editor = CodeMirror.fromTextArea($('.CodeMirror').get(0), {
					lineNumbers: true,
					matchBrackets: true,
					mode: "application/x-httpd-php",
					indentUnit: 4,
					indentWithTabs: true,
					enterMode: "keep", tabMode: "shift"
				});
				$(document).on('mouseup', '.CodeMirror', function() {
					self.templateModel.content = editor.getValue();
				});
				$(document).on('change', '.CodeMirror', function() {
					self.templateModel.content = editor.getValue();
				});
				setInterval(()=>{
					self.templateModel.content = editor.getValue();
				}, 1000)
			}
		}, 500);
	},
	methods: {
		templateSave(){
			var content = this.templateModel.content;
			var name = this.templateModel.name;
			var templateId = this.templateModel.templateId;
			var status = this.templateModel.status;
			
			$.post(templateUpdateUrl, {
				template_id:templateId,
				content: content,
				name: name,
				status: status,
			}).done((response)=>{
				if (response.status !== "success") {
					return Swal.fire({
						icon: 'error',
						title: 'Oops...',
						text: response.message,
					});
				}

				return Swal.fire({
					icon: 'success',
					title: 'Template saved',
				});
			}).fail((result)=>{
				console.log(result);
				return Swal.fire({
					icon: 'error',
					title: 'Oops...',
					text: result,
				});
			});
		}
	}
};
Vue.createApp(TemplateUpdate).mount('#template-update')
	`

	webtemplate := Webpage("Edit Template", h)

	// <style>
	// .CodeMirror {
	// 	border: 1px solid #eee;
	// 	height: auto;
	// }
	// </style>

	webtemplate.AddStyleURLs([]string{
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/codemirror.min.css",
	})
	webtemplate.AddScriptURLs([]string{
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/codemirror.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/xml/xml.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/htmlmixed/htmlmixed.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/javascript/javascript.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/css/css.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/clike/clike.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/php/php.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/2.36.0/formatting.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.22.0/addon/edit/matchbrackets.min.js",
	})
	webtemplate.AddStyle(`	
.CodeMirror {
	border: 1px solid #eee;
	height: auto;
}
	`)
	webtemplate.AddScript(inlineScript)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webtemplate.ToHTML()))
}

// pageTemplatesTemplateUpdateAjax - saves the template via Ajax
func pageTemplatesTemplateUpdateAjax(w http.ResponseWriter, r *http.Request) {
	templateID := strings.Trim(utils.Req(r, "template_id", ""), " ")
	content := strings.Trim(utils.Req(r, "content", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")

	if templateID == "" {
		api.Respond(w, r, api.Error("Template ID is required"))
		return
	}

	template := GetEntityStore().EntityFindByID(templateID)

	if template == nil {
		api.Respond(w, r, api.Error("Template NOT FOUND with ID "+templateID))
		return
	}

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	if status == "" {
		api.Respond(w, r, api.Error("status is required field"))
		return
	}

	template.SetString("content", content)
	template.SetString("name", name)
	isOk := template.SetString("status", status)

	if isOk == false {
		api.Respond(w, r, api.Error("Template failed to be updated"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Template saved successfully", map[string]interface{}{"template_id": template.ID}))
	return
}
