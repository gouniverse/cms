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

func pagePagesPageCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(utils.Req(r, "name", ""), " ")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	entity := EntityCreateWithAttributes("page", map[string]interface{}{
		"name":   name,
		"status": "inactive",
		"title":  name,
		"alias":  "/" + utils.Slugify(name+"-"+utils.RandStr(16), '-'),
	})

	log.Println(entity)

	if entity == nil {
		api.Respond(w, r, api.Error("Page failed to be created"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Page saved successfully", map[string]interface{}{"page_id": entity.ID}))
	return
}

func pagePagesPageUpdateAjax(w http.ResponseWriter, r *http.Request) {
	pageID := strings.Trim(utils.Req(r, "page_id", ""), " ")
	alias := strings.Trim(utils.Req(r, "alias", ""), " ")
	content := strings.Trim(utils.Req(r, "content", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")
	title := strings.Trim(utils.Req(r, "title", ""), " ")
	templateID := strings.Trim(utils.Req(r, "template_id", ""), " ")

	if pageID == "" {
		api.Respond(w, r, api.Error("Page ID is required"))
		return
	}

	page := EntityFindByID(pageID)

	if page == nil {
		api.Respond(w, r, api.Error("Page NOT FOUND with ID "+pageID))
		return
	}

	if alias == "" {
		api.Respond(w, r, api.Error("alias is required field"))
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

	if title == "" {
		api.Respond(w, r, api.Error("title is required field"))
		return
	}

	isOk := EntityAttributesUpsert(pageID, map[string]interface{}{
		"alias":       alias,
		"content":     content,
		"name":        name,
		"status":      status,
		"template_id": templateID,
		"title":       title,
	})

	if isOk == false {
		api.Respond(w, r, api.Error("Page failed to be updated"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Page saved successfully", map[string]interface{}{"page_id": page.ID}))
	return
}

func pagePagesPageUpdate(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	log.Println(endpoint)

	pageID := utils.Req(r, "page_id", "")
	if pageID == "" {
		api.Respond(w, r, api.Error("Page ID is required"))
		return
	}

	page := EntityFindByID(pageID)

	if page == nil {
		api.Respond(w, r, api.Error("Page NOT FOUND with ID "+pageID))
		return
	}

	header := cmsHeader(r.Context().Value(keyEndpoint).(string))
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathPagesPageManager):                       "Pages",
		(endpoint + "?path=" + PathPagesPageUpdate + "&page_id=" + pageID): "Edit page",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "page-update")
	heading := hb.NewHeading1().HTML("Edit Page")
	button := hb.NewButton().HTML("Save").Attr("class", "btn btn-success float-end").Attr("v-on:click", "pageSave")
	heading.AddChild(button)

	tabNavigation := hb.NewUL().Attr("class", "nav nav-tabs").Attr("style", "margin-bottom: 3px;").Attr("role", "tablist")
	tabNavigationContent := hb.NewLI().Attr("class", "nav-item").AddChild(hb.NewHyperlink().Attr("id", "TabContent-tab").Attr("class", "nav-link active").Attr("href", "#TabContent").Attr("v-on:click", "tab('TabContent')").HTML("Content"))
	tabNavigationSeo := hb.NewLI().Attr("class", "nav-item").AddChild(hb.NewHyperlink().Attr("id", "TabSeo-tab").Attr("class", "nav-link").Attr("href", "#TabSeo").Attr("v-on:click", "tab('TabSeo')").HTML("SEO"))
	tabNavigationSettings := hb.NewLI().Attr("class", "nav-item").AddChild(hb.NewHyperlink().Attr("id", "TabSettings-tab").Attr("class", "nav-link").Attr("href", "#TabSettings").Attr("v-on:click", "tab('TabSettings')").HTML("Settings"))
	tabNavigation.AddChild(tabNavigationContent).AddChild(tabNavigationSeo).AddChild(tabNavigationSettings)

	tabContent := hb.NewDiv().Attr("class", "tab-content")
	tabContentContent := hb.NewDiv().Attr("id", "TabContent").Attr("class", "tab-pane fade show active").Attr("data-bs-toggle", "tab")
	tabContentSeo := hb.NewDiv().Attr("id", "TabSeo").Attr("class", "tab-pane fade").Attr("data-bs-toggle", "tab")
	tabContentSettings := hb.NewDiv().Attr("id", "TabSettings").Attr("class", "tab-pane fade").Attr("data-bs-toggle", "tab")
	tabContent.AddChild(tabContentContent).AddChild(tabContentSeo).AddChild(tabContentSettings)

	// <div class="form-group well" style="display:table;width:100%;margin-top:10px;padding:5px 10px;">
	//                     Page address: <a href="<?php echo $page->url(); ?>" target="_blank"><?php echo $page->url(); ?></a> &nbsp;&nbsp;&nbsp; (to change click on Settings tab)
	// 				</div>

	// divPath := hb.NewDiv().Attr("class","form-group well").Attr("style","display:table;width:100%;margin-top:10px;padding:5px 10px;").HTML("Page address: ").HTML(" &nbsp;&nbsp;&nbsp; (to change click on Settings tab)")
	formGroupStatus := hb.NewDiv().Attr("class", "form-group")
	formGroupStatusLabel := hb.NewLabel().HTML("Status").Attr("class", "form-label")
	formGroupStatusSelect := hb.NewSelect().Attr("class", "form-control").Attr("v-model", "pageModel.status")
	formGroupOptionsActive := hb.NewOption().Attr("value", "active").HTML("Active")
	formGroupOptionsInactive := hb.NewOption().Attr("value", "inactive").HTML("Inactive")
	formGroupOptionsTrash := hb.NewOption().Attr("value", "trash").HTML("Trash")
	formGroupStatus.AddChild(formGroupStatusLabel)
	formGroupStatus.AddChild(formGroupStatusSelect.AddChild(formGroupOptionsActive).AddChild(formGroupOptionsInactive).AddChild(formGroupOptionsTrash))

	formGroupName := hb.NewDiv().Attr("class", "form-group")
	formGroupNameLabel := hb.NewLabel().HTML("Name").Attr("class", "form-label")
	formGroupNameInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "pageModel.name")
	formGroupName.AddChild(formGroupNameLabel)
	formGroupName.AddChild(formGroupNameInput)

	formGroupAlias := hb.NewDiv().Attr("class", "form-group")
	formGroupAliasLabel := hb.NewLabel().HTML("Alias / Path").Attr("class", "form-label")
	formGroupAliasInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "pageModel.alias")
	formGroupAlias.AddChild(formGroupAliasLabel)
	formGroupAlias.AddChild(formGroupAliasInput)

	templateList := EntityList("template", 0, 100, "", "id", "asc")
	formGroupTemplate := hb.NewDiv().Attr("class", "form-group")
	formGroupTemplateLabel := hb.NewLabel().HTML("Template").Attr("class", "form-label")
	formGroupTemplateSelect := hb.NewSelect().Attr("class", "form-control").Attr("v-model", "pageModel.templateId")
	formGroupTemplateOptionsEmpty := hb.NewOption().Attr("value", "").HTML("- none -")
	formGroupTemplateSelect.AddChild(formGroupTemplateOptionsEmpty)
	for _, template := range templateList {
		formGroupTemplateOptionsTemplate := hb.NewOption().Attr("value", template.ID).HTML(template.GetAttributeValue("name", "n/a").(string))
		formGroupTemplateSelect.AddChild(formGroupTemplateOptionsTemplate)
	}
	formGroupTemplate.AddChild(formGroupTemplateLabel).AddChild(formGroupTemplateSelect)

	formGroupTitle := hb.NewDiv().Attr("class", "form-group")
	formGroupTitleLabel := hb.NewLabel().HTML("Title").Attr("class", "form-label")
	formGroupTitleInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "pageModel.title")
	formGroupTitle.AddChild(formGroupTitleLabel)
	formGroupTitle.AddChild(formGroupTitleInput)

	formGroupContent := hb.NewDiv().Attr("class", "form-group")
	formGroupContentLabel := hb.NewLabel().HTML("Content").Attr("class", "form-label")
	formGroupContentInput := hb.NewTextArea().Attr("class", "form-control CodeMirror").Attr("v-model", "pageModel.content")
	formGroupContent.AddChild(formGroupContentLabel)
	formGroupContent.AddChild(formGroupContentInput)

	tabContentContent.AddChild(formGroupTitle).AddChild(formGroupContent)
	tabContentSeo.AddChild(formGroupAlias)
	tabContentSettings.AddChild(formGroupStatus).AddChild(formGroupTemplate).AddChild(formGroupName)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))
	container.AddChild(tabNavigation)
	container.AddChild(tabContent)

	h := container.ToHTML()

	alias := page.GetAttributeValue("alias", "").(string)
	content := page.GetAttributeValue("content", "").(string)
	name := page.GetAttributeValue("name", "").(string)
	status := page.GetAttributeValue("status", "").(string)
	templateID := page.GetAttributeValue("template_id", "").(string)
	title := page.GetAttributeValue("title", "").(string)

	contentJSON, _ := json.Marshal(content)
	templateIDJSON, _ := json.Marshal(templateID)

	inlineScript := `
var pageUpdateUrl = "` + endpoint + `?path=pages/page-update-ajax";
var pageId = "` + pageID + `";
var alias = "` + alias + `";
var name = "` + name + `";
var status = "` + status + `";
var title = "` + title + `";
var content = ` + string(contentJSON) + `;
var templateId = ` + string(templateIDJSON) + `;
const PageUpdate = {
	data() {
		return {
			pageModel:{
				pageId: pageId,
				alias: alias,
				content: content,
				name: name,
				status: status,
				title: title,
				templateId: templateId,
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
					self.pageModel.content = editor.getValue();
				});
				$(document).on('change', '.CodeMirror', function() {
					self.pageModel.content = editor.getValue();
				});
				setInterval(()=>{
					self.pageModel.content = editor.getValue();
				}, 1000)
			}
		}, 500);
	},
	methods: {
		tab(id){
			$(".nav-link").removeClass("show active")
			$(".tab-pane").removeClass("show active")
			$("#"+id).addClass("show active")
			$("#"+id+"-tab").addClass("active")
		},
		pageSave(){
			var alias = this.pageModel.alias;
			var content = this.pageModel.content;
			var name = this.pageModel.name;
			var pageId = this.pageModel.pageId;
			var status = this.pageModel.status;
			var templateId = this.pageModel.templateId;
			var title = this.pageModel.title;
			
			$.post(pageUpdateUrl, {
				page_id:pageId,
				alias: alias,
				content: content,
				name: name,
				status: status,
				title: title,
				template_id: templateId,
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
					title: 'Page saved',
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
Vue.createApp(PageUpdate).mount('#page-update')
	`

	webpage := Webpage("Edit Page", h)
	webpage.AddStyleURLs([]string{
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/codemirror.min.css",
	})
	webpage.AddScriptURLs([]string{
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
	webpage.AddStyle(`	
.CodeMirror {
	border: 1px solid #eee;
	height: auto;
}
	`)
	webpage.AddScript(inlineScript)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

func pagePagesPageManager(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	log.Println(endpoint)

	header := cmsHeader(endpoint)
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathPagesPageManager): "Pages",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "page-manager")
	heading := hb.NewHeading1().HTML("Page Manager")
	button := hb.NewButton().HTML("New page").Attr("class", "btn btn-success float-end").Attr("v-on:click", "showPageCreateModal")
	heading.AddChild(button)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))

	modal := hb.NewDiv().Attr("id", "ModalPageCreate").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("New Page"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewDiv().Attr("class", "form-group").AddChild(hb.NewLabel().HTML("Name")).AddChild(hb.NewInput().Attr("class", "form-control").Attr("v-model", "pageCreateModel.name")))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-prsecondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Create & Continue").Attr("class", "btn btn-primary").Attr("v-on:click", "pageCreate"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	container.AddChild(modal)

	pages := EntityList("page", 0, 200, "", "id", "asc")

	table := hb.NewTable().Attr("class", "table table-responsive table-striped mt-3")
	thead := hb.NewThead()
	tbody := hb.NewTbody()
	table.AddChild(thead).AddChild(tbody)

	tr := hb.NewTR()
	th1 := hb.NewTD().HTML("Name")
	th2 := hb.NewTD().HTML("Status")
	th3 := hb.NewTD().HTML("Actions").Attr("style", "width:1px;")
	thead.AddChild(tr.AddChild(th1).AddChild(th2).AddChild(th3))

	for _, page := range pages {
		name := page.GetAttributeValue("name", "n/a").(string)
		alias := page.GetAttributeValue("alias", "n/a").(string)
		status := page.GetAttributeValue("status", "n/a").(string)
		buttonEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Attr("class", "btn btn-primary").Attr("v-on:click", "pageEdit('"+page.ID+"')")

		tr := hb.NewTR()
		td1 := hb.NewTD().HTML(name).AddChild(hb.NewDiv().HTML(alias).Attr("style", "font-size:11px;"))
		td2 := hb.NewTD().HTML(status)
		td3 := hb.NewTD().AddChild(buttonEdit)

		tbody.AddChild(tr.AddChild(td1).AddChild(td2).AddChild(td3))
	}
	container.AddChild(table)

	h := container.ToHTML()

	inlineScript := `
var pageCreateUrl = "` + endpoint + `?path=pages/page-create-ajax"
var pageUpdateUrl = "` + endpoint + `?path=pages/page-update"
const PageManager = {
	data() {
		return {
		  pageCreateModel:{
			  name:""
		  }
		}
	},
	methods: {
        showPageCreateModal(){
			//alert("Create page");
			var modalPageCreate = new bootstrap.Modal(document.getElementById('ModalPageCreate'));
			modalPageCreate.show();
		},
		pageCreate(){
			var name = this.pageCreateModel.name;
			$.post(pageCreateUrl, {name: name}).done((result)=>{
				if (result.status==="success"){
					var modalPageCreate = new bootstrap.Modal(document.getElementById('ModalPageCreate'));
					modalPageCreate.hide();

					return location.href = pageUpdateUrl+ "&page_id=" + result.data.page_id;
				}
				alert("Failed. " + result.message)
			}).fail((result)=>{
				alert("Failed" + result)
			});
		},
		pageEdit(pageId){
			return location.href = pageUpdateUrl+ "&page_id=" + pageId;
		}
	}
};
Vue.createApp(PageManager).mount('#page-manager')
	`

	webpage := Webpage("Page Manager", h)
	webpage.AddScript(inlineScript)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}
