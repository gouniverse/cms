package cms

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/bs"
	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/responses"

	// "github.com/gouniverse/cms/ve"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/utils"
)

func (cms Cms) pagePagesPageCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(utils.Req(r, "name", ""), " ")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	page, err := cms.EntityStore.EntityCreate(ENTITY_TYPE_PAGE)

	if err != nil {
		api.Respond(w, r, api.Error("Page failed to be created: "+err.Error()))
		return
	}

	// log.Println(page)

	if page == nil {
		api.Respond(w, r, api.Error("Page failed to be created"))
		return
	}

	page.SetString("name", name)
	page.SetString("status", "inactive")
	page.SetString("title", name)
	page.SetString("alias", "/"+utils.StrSlugify(name+"-"+utils.StrRandom(16), '-'))

	api.Respond(w, r, api.SuccessWithData("Page saved successfully", map[string]interface{}{"page_id": page.ID()}))
}

func (cms Cms) pagePagesPageUpdateAjax(w http.ResponseWriter, r *http.Request) {
	pageID := strings.Trim(utils.Req(r, "page_id", ""), " ")
	alias := strings.Trim(utils.Req(r, "alias", ""), " ")
	canonicalURL := strings.Trim(utils.Req(r, "canonical_url", ""), " ")
	content := strings.Trim(utils.Req(r, "content", ""), " ")
	contentEditor := strings.Trim(utils.Req(r, "content_editor", ""), " ")
	metaDescription := strings.Trim(utils.Req(r, "meta_description", ""), " ")
	metaKeywords := strings.Trim(utils.Req(r, "meta_keywords", ""), " ")
	metaRobots := strings.Trim(utils.Req(r, "meta_robots", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")
	title := strings.Trim(utils.Req(r, "title", ""), " ")
	templateID := strings.Trim(utils.Req(r, "template_id", ""), " ")
	handle := strings.Trim(utils.Req(r, "handle", ""), " ")

	if pageID == "" {
		api.Respond(w, r, api.Error("Page ID is required"))
		return
	}

	page, _ := cms.EntityStore.EntityFindByID(pageID)

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

	page.SetString("alias", alias)
	page.SetString("canonical_url", canonicalURL)
	page.SetString("content", content)
	page.SetString("content_editor", contentEditor)
	page.SetString("meta_description", metaDescription)
	page.SetString("meta_keywords", metaKeywords)
	page.SetString("meta_robots", metaRobots)
	page.SetString("name", name)
	page.SetString("status", status)
	page.SetString("template_id", templateID)
	page.SetString("handle", handle)
	err := page.SetString("title", title)

	if err != nil {
		api.Respond(w, r, api.Error("Page failed to be updated: "+err.Error()))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Page saved successfully", map[string]interface{}{"page_id": page.ID()}))
}

func (cms Cms) pagePagesPageUpdate(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	// log.Println(endpoint)

	pageID := utils.Req(r, "page_id", "")
	if pageID == "" {
		api.Respond(w, r, api.Error("Page ID is required"))
		return
	}

	page, _ := cms.EntityStore.EntityFindByID(pageID)

	if page == nil {
		api.Respond(w, r, api.Error("Page NOT FOUND with ID "+pageID))
		return
	}

	header := cms.cmsHeader(r.Context().Value(keyEndpoint).(string))
	breadcrumbs := cms.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  endpoint,
			Name: "Home",
		},
		{
			URL:  (endpoint + "?path=" + PathPagesPageManager),
			Name: "Pages",
		},
		{
			URL:  (endpoint + "?path=" + PathPagesPageUpdate + "&page_id=" + pageID),
			Name: "Edit page",
		},
	})

	container := hb.NewDiv().ID("page-update").Class("container")
	heading := hb.NewHeading1().HTML("Edit Page")
	saveButton := hb.NewButton().HTML("Save").Class("btn btn-success float-end").Attr("v-on:click", "pageSave")
	heading.Child(saveButton)

	tabNavigation := bs.NavTabs().Style("margin-bottom: 3px;").Attr("role", "tablist")
	tabNavigationContent := bs.NavItem().Child(bs.NavLink().ID("TabContent-tab").Attr("href", "#TabContent").Attr("v-on:click", "tab('TabContent')").HTML("Content"))
	tabNavigationSeo := bs.NavItem().Child(bs.NavLink().ID("TabSeo-tab").Attr("href", "#TabSeo").Attr("v-on:click", "tab('TabSeo')").HTML("SEO"))
	tabNavigationSettings := bs.NavItem().Child(bs.NavLink().ID("TabSettings-tab").Attr("href", "#TabSettings").Attr("v-on:click", "tab('TabSettings')").HTML("Settings"))
	tabNavigation.AddChild(tabNavigationContent).AddChild(tabNavigationSeo).AddChild(tabNavigationSettings)

	tabContent := hb.NewDiv().Class("tab-content")
	tabContentContent := hb.NewDiv().ID("TabContent").Attr("class", "tab-pane fade show active").Attr("data-bs-toggle", "tab")
	tabContentSeo := hb.NewDiv().ID("TabSeo").Attr("class", "tab-pane fade").Attr("data-bs-toggle", "tab")
	tabContentSettings := hb.NewDiv().ID("TabSettings").Attr("class", "tab-pane fade").Attr("data-bs-toggle", "tab")
	tabContent.AddChild(tabContentContent).AddChild(tabContentSeo).AddChild(tabContentSettings)

	// <div class="form-group well" style="display:table;width:100%;margin-top:10px;padding:5px 10px;">
	//                     Page address: <a href="<?php echo $page->url(); ?>" target="_blank"><?php echo $page->url(); ?></a> &nbsp;&nbsp;&nbsp; (to change click on Settings tab)
	// 				</div>

	// Status
	formGroupStatus := bs.FormGroup().Children([]*hb.Tag{
		bs.FormLabel("Status"),
		bs.FormSelect().Attr("v-model", "pageModel.status").Children([]*hb.Tag{
			bs.FormSelectOption("active", "Active"),
			bs.FormSelectOption("inactive", "Inactive"),
			bs.FormSelectOption("trash", "Trash"),
		}),
	})

	// Content Editor
	formGroupEditor := bs.FormGroup().Children([]*hb.Tag{
		bs.FormLabel("Content Editor"),
		bs.FormSelect().Attr("v-model", "pageModel.contentEditor").Children([]*hb.Tag{
			bs.FormSelectOption("", "- none -"),
			bs.FormSelectOption("codemirror", "CodeMirror"),
			// bs.FormSelectOption("visual", "Visual Editor (Experimental)"),
		}),
		bs.FormText("The content editor allows you to select the mode for editing the content. Note you will need to save and refresh to activate"),
	})

	// Name
	formGroupName := bs.FormGroup().Children([]*hb.Tag{
		bs.FormLabel("Name"),
		bs.FormInput().Attr("v-model", "pageModel.name"),
		bs.FormText("The name of the page as visible in the admin panel. This is not vsible to the page vistors"),
	})

	// Alias
	formGroupAlias := bs.FormGroup().Children([]*hb.Tag{
		bs.FormLabel("Alias / Path"),
		bs.FormInput().Attr("v-model", "pageModel.alias"),
		bs.FormText("The relative path on the website where this page will be visible to the vistors"),
	})

	// Canonical Url
	formGroupCanonicalURL := bs.FormGroup().Children([]*hb.Tag{
		bs.FormLabel("Canonical Url"),
		bs.FormInput().Attr("v-model", "pageModel.canonicalUrl"),
	})

	// Meta Description
	formGroupMetaDescription := bs.FormGroup().Children([]*hb.Tag{
		bs.FormLabel("Meta Description"),
		bs.FormInput().Attr("v-model", "pageModel.metaDescription"),
	})

	// Meta Keywords
	formGroupMetaKeywords := hb.NewDiv().Class("form-group").Children([]*hb.Tag{
		bs.FormLabel("Meta Keywords"),
		bs.FormInput().Attr("v-model", "pageModel.metaKeywords"),
	})

	// Robots
	formGroupMetaRobots := hb.NewDiv().Class("form-group").Children([]*hb.Tag{
		bs.FormLabel("Meta Robots"),
		bs.FormInput().Attr("v-model", "pageModel.metaRobots"),
	})

	// Template
	templateList, err := cms.EntityStore.EntityList(entitystore.EntityQueryOptions{
		EntityType: "template",
		Offset:     0,
		Limit:      100,
		SortBy:     "id",
		SortOrder:  "asc",
	})
	if err != nil {
		api.Respond(w, r, api.Error("Entity list failed to be retrieved "+err.Error()))
		return
	}
	formGroupTemplateSelect := bs.FormSelect().Attr("v-model", "pageModel.templateId")
	formGroupTemplateOptionsEmpty := bs.FormSelectOption("", "- none -")
	formGroupTemplateSelect.Child(formGroupTemplateOptionsEmpty)
	for _, template := range templateList {
		templateName, _ := template.GetString("name", "n/a")
		formGroupTemplateOptionsTemplate := bs.FormSelectOption(template.ID(), templateName)
		formGroupTemplateSelect.Child(formGroupTemplateOptionsTemplate)
	}
	formGroupTemplate := bs.FormGroup().Children([]*hb.Tag{
		bs.FormLabel("Template").Class("form-label"),
		formGroupTemplateSelect,
		bs.FormText("Select the template that this page content will be displayed in. This feature is useful if you want to implement consistent layouts. Leaving the template field empty will display page as it is standalone"),
	})

	// Title
	formGroupTitle := hb.NewDiv().Class("form-group")
	formGroupTitleLabel := hb.NewLabel().HTML("Title").Class("form-label")
	formGroupTitleInput := hb.NewInput().Class("form-control").Attr("v-model", "pageModel.title")
	formGroupTitle.Child(formGroupTitleLabel).Child(formGroupTitleInput)

	// Content
	// editor, _ := page.GetString("content_editor", "")
	formGroupContent := bs.FormGroup()
	formGroupContentLabel := bs.FormLabel("Content")
	formGroupContentInput := hb.NewTextArea().Class("form-control CodeMirror").Attr("v-model", "pageModel.content")
	// formGroupContentInputVisualData := hb.NewTextArea().Class("form-control").Attr("v-model", "pageModel.contentVisual")
	// if editor == "visualeditor" {
	// formGroupContentInput.Style("display:none")
	// } else {
	// formGroupContentInputVisualData.Style("display:none")
	// }
	formGroupContent.Children([]*hb.Tag{
		formGroupContentLabel,
		formGroupContentInput,
		// formGroupContentInputVisualData,
		// hb.NewHTML(ve.VisualeditorContent()),
	})

	tabContentContent.Children([]*hb.Tag{
		formGroupTitle,
		formGroupContent,
	})
	tabContentSeo.Children([]*hb.Tag{
		formGroupAlias,
		formGroupMetaDescription,
		formGroupMetaKeywords,
		formGroupMetaRobots,
		formGroupCanonicalURL,
	})
	tabContentSettings.Children([]*hb.Tag{
		formGroupStatus,
		formGroupTemplate,
		formGroupName,
		formGroupEditor,
	})

	container.Children([]*hb.Tag{
		hb.NewHTML(header),
		heading,
		hb.NewHTML(breadcrumbs),
		tabNavigation,
		tabContent,
	})

	h := container.ToHTML()

	alias, _ := page.GetString("alias", "")
	content, _ := page.GetString("content", "")
	contentEditor, _ := page.GetString("content_editor", "")
	name, _ := page.GetString("name", "")
	status, _ := page.GetString("status", "")
	templateID, _ := page.GetString("template_id", "")
	title, _ := page.GetString("title", "")
	metaDescription, _ := page.GetString("meta_description", "")
	metaKeywords, _ := page.GetString("meta_keywords", "")
	metaRobots, _ := page.GetString("meta_robots", "")
	canonicalURL, _ := page.GetString("canonical_url", "")

	canonicalURLJSON, _ := json.Marshal(canonicalURL)
	contentJSON, _ := json.Marshal(content)
	contentEditorJSON, _ := json.Marshal(contentEditor)
	nameJSON, _ := json.Marshal(name)
	templateIDJSON, _ := json.Marshal(templateID)
	titleJSON, _ := json.Marshal(title)

	inlineScript := `
var pageUpdateUrl = "` + endpoint + `?path=pages/page-update-ajax";
var pageId = "` + pageID + `";
var alias = "` + alias + `";
var canonicalUrl = "` + canonicalURL + `";
var metaDescription = "` + metaDescription + `";
var metaKeywords = "` + metaKeywords + `";
var metaRobots = "` + metaRobots + `";
var name = ` + string(nameJSON) + `;
var status = "` + status + `";
var title = ` + string(titleJSON) + `;
var canonicalUrl = ` + string(canonicalURLJSON) + `;
var content = ` + string(contentJSON) + `;
var contentEditor = ` + string(contentEditorJSON) + `;
var templateId = ` + string(templateIDJSON) + `;
const PageUpdate = {
	data() {
		return {
			pageModel:{
				pageId: pageId,
				alias: alias,
				canonicalUrl:canonicalUrl,
				content: content,
				contentEditor: contentEditor,
				metaDescription:metaDescription,
				metaKeywords:metaKeywords,
				metaRobots:metaRobots,
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
			var canonicalUrl = this.pageModel.canonicalUrl;
			var content = this.pageModel.content;
			var contentEditor = this.pageModel.contentEditor;
			var metaDescription = this.pageModel.metaDescription;
			var metaKeywords = this.pageModel.metaKeywords;
			var metaRobots = this.pageModel.metaRobots;
			var name = this.pageModel.name;
			var pageId = this.pageModel.pageId;
			var status = this.pageModel.status;
			var templateId = this.pageModel.templateId;
			var title = this.pageModel.title;
			
			$.post(pageUpdateUrl, {
				page_id:pageId,
				alias: alias,
				content: content,
				content_editor: contentEditor,
				canonical_url:canonicalUrl,
				meta_description:metaDescription,
				meta_keywords:metaKeywords,
				meta_robots:metaRobots,
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
	// webpage.AddScript(ve.VisualeditorScripts())

	responses.HTMLResponse(w, r, cms.funcLayout(webpage.ToHTML()))
}

func (cms Cms) pagePagesPageManager(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	// log.Println(endpoint)

	pages, err := cms.EntityStore.EntityList(entitystore.EntityQueryOptions{
		EntityType: ENTITY_TYPE_PAGE,
		Offset:     0,
		Limit:      200,
		SortBy:     "id",
		SortOrder:  "asc",
	})

	if err != nil {
		api.Respond(w, r, api.Error("Page list failed to be retrieved "+err.Error()))
		return
	}

	header := cms.cmsHeader(endpoint)
	breadcrums := cms.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  endpoint,
			Name: "Home",
		},
		{
			URL:  (endpoint + "?path=" + PathPagesPageManager),
			Name: "Pages",
		},
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "page-manager")
	heading := hb.NewHeading1().HTML("Page Manager")
	button := hb.NewButton().HTML("New page").Attr("class", "btn btn-success float-end").Attr("v-on:click", "showPageCreateModal")
	heading.AddChild(button)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))

	// modal := hb.NewDiv().Attr("id", "ModalPageCreate").Attr("class", "modal fade")
	// modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	// modalContent := hb.NewDiv().Attr("class", "modal-content")
	// modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("New Page"))
	// modalBody := hb.NewDiv().Attr("class", "modal-body")
	// modalBody.AddChild(hb.NewDiv().Attr("class", "form-group").AddChild(hb.NewLabel().HTML("Name")).AddChild(hb.NewInput().Attr("class", "form-control").Attr("v-model", "pageCreateModel.name")))
	// modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	// modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-prsecondary").Attr("data-bs-dismiss", "modal"))
	// modalFooter.AddChild(hb.NewButton().HTML("Create & Continue").Attr("class", "btn btn-primary").Attr("v-on:click", "pageCreate"))
	// modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	// modalDialog.AddChild(modalContent)
	// modal.AddChild(modalDialog)
	container.AddChild(pagePagesPageCreateModal())
	container.AddChild(pagePagesPageTrashModal())

	table := hb.NewTable().Attr("id", "TablePages").Attr("class", "table table-responsive table-striped mt-3")
	thead := hb.NewThead()
	tbody := hb.NewTbody()
	table.AddChild(thead).AddChild(tbody)

	tr := hb.NewTR()
	th1 := hb.NewTD().HTML("Name")
	th2 := hb.NewTD().HTML("Status")
	th3 := hb.NewTD().HTML("Actions").Attr("style", "width:150px;")
	thead.AddChild(tr.AddChild(th1).AddChild(th2).AddChild(th3))

	for _, page := range pages {
		name, _ := page.GetString("name", "n/a")
		alias, _ := page.GetString("alias", "n/a")
		status, _ := page.GetString("status", "n/a")
		buttonEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Attr("class", "btn btn-primary btn-sm").Attr("v-on:click", "pageEdit('"+page.ID()+"')").Attr("style", "margin-right:5px")
		buttonTrash := hb.NewButton().HTML("Trash").Attr("type", "button").Attr("class", "btn btn-danger btn-sm").Attr("v-on:click", "showPageTrashModal('"+page.ID()+"')")

		tr := hb.NewTR()
		td1 := hb.NewTD().HTML(name).AddChild(hb.NewDiv().HTML(alias).Attr("style", "font-size:11px;"))
		td2 := hb.NewTD().HTML(status)
		td3 := hb.NewTD().AddChild(buttonEdit).AddChild(buttonTrash)

		tbody.AddChild(tr.AddChild(td1).AddChild(td2).AddChild(td3))
	}
	container.AddChild(table)

	h := container.ToHTML()

	inlineScript := `
var pageCreateUrl = "` + endpoint + `?path=pages/page-create-ajax"
var pageTrashUrl = "` + endpoint + `?path=pages/page-trash-ajax"
var pageUpdateUrl = "` + endpoint + `?path=pages/page-update"
const PageManager = {
	data() {
		return {
		  pageCreateModel:{
			  name:""
		  },
		  pageTrashModel:{
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
				$('#TablePages').DataTable({
					"order": [[ 0, "asc" ]] // 1st column
				});
			});
		},
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
		},
		showPageTrashModal(pageId){
			this.pageTrashModel.id = pageId;
			var modalPageTrash = new bootstrap.Modal(document.getElementById('ModalPageTrash'));
			modalPageTrash.show();
		},
		pageTrash(){
            let pageId = this.pageTrashModel.id;
			$.post(pageTrashUrl, {page_id: pageId}).done((result)=>{
				if (result.status==="success"){
					var ModalPageTrash = new bootstrap.Modal(document.getElementById('ModalPageTrash'));
				    ModalPageTrash.hide();
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
Vue.createApp(PageManager).mount('#page-manager')
	`

	webpage := Webpage("Page Manager", h)
	webpage.AddStyleURL("https://cdnjs.cloudflare.com/ajax/libs/datatables/1.10.21/css/jquery.dataTables.css")
	webpage.AddScriptURL("https://cdnjs.cloudflare.com/ajax/libs/datatables/1.10.21/js/jquery.dataTables.js")
	webpage.AddScript(inlineScript)

	responses.HTMLResponse(w, r, cms.funcLayout(webpage.ToHTML()))
}

// pagePagesPageTrashAjax - moves the template to the trash
func (cms Cms) pagePagesPageTrashAjax(w http.ResponseWriter, r *http.Request) {
	pageID := strings.Trim(utils.Req(r, "page_id", ""), " ")

	if pageID == "" {
		api.Respond(w, r, api.Error("Page ID is required"))
		return
	}

	page, _ := cms.EntityStore.EntityFindByID(pageID)

	if page == nil {
		api.Respond(w, r, api.Error("Page NOT FOUND with ID "+pageID))
		return
	}

	isOk, err := cms.EntityStore.EntityTrash(pageID)

	if err != nil {
		api.Respond(w, r, api.Error("Entity failed to be trashed "+err.Error()))
		return
	}

	if !isOk {
		api.Respond(w, r, api.Error("Page failed to be moved to trash"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Page moved to trash successfully", map[string]interface{}{"page_id": page.ID()}))
}

func pagePagesPageTrashModal() *hb.Tag {
	modal := hb.NewDiv().Attr("id", "ModalPageTrash").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("Trash Template"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewParagraph().HTML("Are you sure you want to move this page to trash bin?"))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-secondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Move to trash bin").Attr("class", "btn btn-danger").Attr("v-on:click", "pageTrash"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	return modal
}

func pagePagesPageCreateModal() *hb.Tag {
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
	return modal
}
