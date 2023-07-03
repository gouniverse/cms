package cms

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/bs"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/responses"
	"github.com/gouniverse/utils"
)

func (cms Cms) pageTemplatesTemplateCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(utils.Req(r, "name", ""), " ")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	template, err := cms.EntityStore.EntityCreate("template")

	if err != nil {
		api.Respond(w, r, api.Error("Template failed to be created: "+err.Error()))
		return
	}

	if template == nil {
		api.Respond(w, r, api.Error("Template failed to be created"))
		return
	}

	template.SetString("name", name)
	template.SetString("status", "inactive")
	cms.EntityStore.EntityUpdate(*template)

	api.Respond(w, r, api.SuccessWithData("Template saved successfully", map[string]interface{}{"template_id": template.ID()}))
}

func (cms Cms) pageTemplatesTemplateManager(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	// log.Println(endpoint)

	templates, err := cms.EntityStore.EntityList(entitystore.EntityQueryOptions{
		EntityType: "template",
		Offset:     0,
		Limit:      200,
		SortBy:     "id",
		SortOrder:  "asc",
	})

	if err != nil {
		api.Respond(w, r, api.Error("Templates failed to be fetched: "+err.Error()))
		return
	}

	header := cms.cmsHeader(endpoint)
	breadcrumbs := cms.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  endpoint,
			Name: "Home",
		},
		{
			URL:  (endpoint + "?path=" + PathTemplatesTemplateManager),
			Name: "Templates",
		},
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "template-manager")
	heading := hb.NewHeading1().HTML("Template Manager")
	button := hb.NewButton().HTML("New template").Attr("class", "btn btn-success float-end").Attr("v-on:click", "showTemplateCreateModal")
	heading.AddChild(button)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrumbs))

	container.AddChild(pageTemplatesTemplateTrashModal())
	container.AddChild(pageTemplatesTemplateCreateModal())

	table := hb.NewTable().Attr("id", "TableTemplates").Attr("class", "table table-responsive table-striped mt-3")
	thead := hb.NewThead()
	tbody := hb.NewTbody()
	table.AddChild(thead).AddChild(tbody)

	tr := hb.NewTR()
	th1 := hb.NewTD().HTML("Name")
	th2 := hb.NewTD().HTML("Status")
	th3 := hb.NewTD().HTML("Actions").Attr("style", "width:120px;")
	thead.AddChild(tr.AddChild(th1).AddChild(th2).AddChild(th3))

	for _, template := range templates {
		name, err := template.GetString("name", "n/a")
		if err != nil {
			api.Respond(w, r, api.Error("Attribute 'name' failed to be fetched: "+err.Error()))
			return
		}
		status, err := template.GetString("status", "n/a")
		if err != nil {
			api.Respond(w, r, api.Error("Attribute 'status' failed to be fetched: "+err.Error()))
			return
		}
		buttonEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Attr("class", "btn btn-primary btn-sm").Attr("v-on:click", "templateEdit('"+template.ID()+"')").Attr("style", "margin-right:5px")
		buttonTrash := hb.NewButton().HTML("Trash").Attr("type", "button").Attr("class", "btn btn-danger btn-sm").Attr("v-on:click", "showTemplateTrashModal('"+template.ID()+"')")

		tr := hb.NewTR()
		td1 := hb.NewTD().HTML(name)
		td2 := hb.NewTD().HTML(status)
		td3 := hb.NewTD().AddChild(buttonEdit).AddChild(buttonTrash)

		tbody.AddChild(tr.AddChild(td1).AddChild(td2).AddChild(td3))
	}
	container.AddChild(table)

	h := container.ToHTML()

	inlineScript := `
var templateCreateUrl = "` + endpoint + `?path=templates/template-create-ajax";
var templateTrashUrl = "` + endpoint + `?path=templates/template-trash-ajax";
var templateUpdateUrl = "` + endpoint + `?path=templates/template-update";
const TemplateManager = {
	data() {
		return {
		  templateCreateModel:{
			name:""
		  },
		  templateTrashModel:{
			id:""
		  }
		}
	},
	created(){
		this.initDataTable();
	},
	methods: {
		initDataTable(){
			$(() => {
				$('#TableTemplates').DataTable({
					"order": [[ 0, "asc" ]] // 1st column
				});
			});
		},
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
		},
		showTemplateTrashModal(templateId){
			this.templateTrashModel.id = templateId;
			var modalTemplateTrash = new bootstrap.Modal(document.getElementById('ModalTemplateTrash'));
			modalTemplateTrash.show();
		},
		templateTrash(){
            let templateId = this.templateTrashModel.id;
			$.post(templateTrashUrl, {template_id: templateId}).done((result)=>{
				if (result.status==="success"){
					var ModalTemplateTrash = new bootstrap.Modal(document.getElementById('ModalTemplateTrash'));
				    ModalTemplateTrash.hide();
					location.href = location.href;
					return;
				}
				alert("Failed. " + result.message)
			}).fail((result)=>{
				alert("Failed" + result)
			});
		}
	}
};
Vue.createApp(TemplateManager).mount('#template-manager')
	`

	if cms.funcLayout("") != "" {
		out := hb.NewWrap().Children([]*hb.Tag{
			hb.NewStyleURL(cdn.JqueryDataTablesCss_1_13_4()),
			hb.NewHTML(h),
			hb.NewScriptURL(cdn.Jquery_3_6_4()),
			hb.NewScriptURL(cdn.VueJs_3()),
			hb.NewScriptURL(cdn.Sweetalert2_10()),
			hb.NewScriptURL(cdn.JqueryDataTablesJs_1_13_4()),
			hb.NewScript(inlineScript),
		}).ToHTML()
		responses.HTMLResponse(w, r, cms.funcLayout(out))
		return
	}

	webpage := Webpage("Template Manager", h)
	webpage.AddStyleURL(cdn.JqueryDataTablesCss_1_13_4())
	webpage.AddScriptURL(cdn.JqueryDataTablesJs_1_13_4())
	webpage.AddScript(inlineScript)
	responses.HTMLResponse(w, r, webpage.ToHTML())
}

// pageTemplatesTemplateUpdate shows the template edit page
func (cms Cms) pageTemplatesTemplateUpdate(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	// log.Println(endpoint)

	templateID := utils.Req(r, "template_id", "")
	if templateID == "" {
		api.Respond(w, r, api.Error("Template ID is required"))
		return
	}

	template, _ := cms.EntityStore.EntityFindByID(templateID)

	if template == nil {
		api.Respond(w, r, api.Error("Template NOT FOUND with ID "+templateID))
		return
	}

	header := cms.cmsHeader(r.Context().Value(keyEndpoint).(string))
	breadcrumbs := cms.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  endpoint,
			Name: "Home",
		},
		{
			URL:  (endpoint + "?path=" + PathTemplatesTemplateManager),
			Name: "Templates",
		},
		{
			URL:  (endpoint + "?path=" + PathTemplatesTemplateUpdate + "&template_id=" + templateID),
			Name: "Edit template",
		},
	})

	container := hb.NewDiv().Class("container").ID("template-update")
	heading := hb.NewHeading1().HTML("Edit Template")
	button := hb.NewButton().HTML("Save").Class("btn btn-success float-end").Attr("v-on:click", "templateSave")
	heading.AddChild(button)

	formGroupStatus := hb.NewDiv().Class("form-group mb-3")
	formGroupStatusLabel := hb.NewLabel().HTML("Status").Class("form-label")
	formGroupStatusSelect := hb.NewSelect().Class("form-select").Attr("v-model", "templateModel.status")
	formGroupOptionsActive := hb.NewOption().Value("active").HTML("Active")
	formGroupOptionsInactive := hb.NewOption().Value("inactive").HTML("Inactive")
	formGroupOptionsTrash := hb.NewOption().Value("trash").HTML("Trash")
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
	container.AddChild(hb.NewHTML(breadcrumbs))
	container.AddChild(formGroupStatus).AddChild(formGroupName).AddChild(formGroupContent)

	h := container.ToHTML()

	name, _ := template.GetString("name", "")
	status, _ := template.GetString("status", "")
	content, _ := template.GetString("content", "")
	templateJSON, _ := json.Marshal(templateID)
	nameJSON, _ := json.Marshal(name)
	statusJSON, _ := json.Marshal(status)
	contentJSON, _ := json.Marshal(content)

	inlineScript := `
var templateTrashUrl = "` + endpoint + `?path=templates/template-trash-ajax";
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

	if cms.funcLayout("") != "" {
		out := hb.NewWrap().Children([]*hb.Tag{
			hb.NewStyleURL(codemirrorCss),
			hb.NewStyle(`.CodeMirror {
				border: 1px solid #eee;
				height: auto;
			}`),
			hb.NewHTML(h),
			hb.NewScriptURL(cdn.Jquery_3_6_4()),
			hb.NewScriptURL(cdn.VueJs_3()),
			hb.NewScriptURL(cdn.Sweetalert2_10()),
			hb.NewScriptURL(codemirrorJs),
			hb.NewScriptURL(codemirrorHtmlmixedJs),
			hb.NewScriptURL(codemirrorJavascriptJs),
			hb.NewScriptURL(codemirrorCssJs),
			hb.NewScriptURL(codemirrorClikeJs),
			hb.NewScriptURL(codemirrorPhpJs),
			hb.NewScriptURL(codemirrorFormattingJs),
			hb.NewScriptURL(codemirrorMatchBracketsJs),
			hb.NewScript(inlineScript),
		}).ToHTML()
		responses.HTMLResponse(w, r, cms.funcLayout(out))
		return
	}

	webpage := Webpage("Edit Template", h).
		AddStyleURLs([]string{
			codemirrorCss,
		}).
		AddScriptURLs([]string{
			codemirrorJs,
			codemirrorXmlJs,
			codemirrorHtmlmixedJs,
			codemirrorJavascriptJs,
			codemirrorCssJs,
			codemirrorClikeJs,
			codemirrorPhpJs,
			codemirrorFormattingJs,
			codemirrorMatchBracketsJs,
		}).
		AddStyle(`	
.CodeMirror {
	border: 1px solid #eee;
	height: auto;
}`).
		AddScript(inlineScript)

	responses.HTMLResponse(w, r, webpage.ToHTML())
}

// pageTemplatesTemplateTrashAjax - moves the template to the trash
func (cms Cms) pageTemplatesTemplateTrashAjax(w http.ResponseWriter, r *http.Request) {
	templateID := strings.Trim(utils.Req(r, "template_id", ""), " ")

	if templateID == "" {
		api.Respond(w, r, api.Error("Template ID is required"))
		return
	}

	template, _ := cms.EntityStore.EntityFindByID(templateID)

	if template == nil {
		api.Respond(w, r, api.Error("Template NOT FOUND with ID "+templateID))
		return
	}

	isOk, err := cms.EntityStore.EntityTrash(templateID)

	if err != nil {
		api.Respond(w, r, api.Error("Template failed to be moved to trash "+err.Error()))
		return
	}

	if !isOk {
		api.Respond(w, r, api.Error("Template failed to be moved to trash"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Template moved to trash successfully", map[string]interface{}{"template_id": template.ID()}))
}

// pageTemplatesTemplateUpdateAjax - saves the template via Ajax
func (cms Cms) pageTemplatesTemplateUpdateAjax(w http.ResponseWriter, r *http.Request) {
	templateID := strings.Trim(utils.Req(r, "template_id", ""), " ")
	content := strings.Trim(utils.Req(r, "content", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")
	handle := strings.Trim(utils.Req(r, "handle", ""), " ")

	if templateID == "" {
		api.Respond(w, r, api.Error("Template ID is required"))
		return
	}

	template, _ := cms.EntityStore.EntityFindByID(templateID)

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

	err := cms.EntityStore.AttributeSetString(template.ID(), "content", content)
	if err != nil {
		api.Respond(w, r, api.Error("Content failed to be updated: "+err.Error()))
		return
	}

	err = cms.EntityStore.AttributeSetString(template.ID(), "name", name)
	if err != nil {
		api.Respond(w, r, api.Error("Name failed to be updated: "+err.Error()))
		return
	}

	err = cms.EntityStore.AttributeSetString(template.ID(), "status", status)
	if err != nil {
		api.Respond(w, r, api.Error("Status failed to be updated: "+err.Error()))
		return
	}

	template.SetHandle(handle)
	errUpdate := cms.EntityStore.EntityUpdate(*template)

	if errUpdate != nil {
		api.Respond(w, r, api.Error("Template failed to be updated: "+errUpdate.Error()))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Template saved successfully", map[string]interface{}{"template_id": template.ID()}))
}

func pageTemplatesTemplateTrashModal() *hb.Tag {
	modal := hb.NewDiv().Attr("id", "ModalTemplateTrash").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("Trash Template"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewParagraph().HTML("Are you sure you want to move this template to trash bin?"))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-secondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Move to trash bin").Attr("class", "btn btn-danger").Attr("v-on:click", "templateTrash"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	return modal
}

func pageTemplatesTemplateCreateModal() *hb.Tag {
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
	return modal
}
