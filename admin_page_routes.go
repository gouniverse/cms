package cms

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/utils"
)

func pagePagesPageCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(utils.Req(r, "name", ""), " ")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	page := GetEntityStore().EntityCreate("page")

	log.Println(page)

	if page == nil {
		api.Respond(w, r, api.Error("Page failed to be created"))
		return
	}

	page.SetString("name", name)
	page.SetString("status", "inactive")
	page.SetString("title", name)
	page.SetString("alias", "/"+utils.Slugify(name+"-"+utils.RandStr(16), '-'))

	api.Respond(w, r, api.SuccessWithData("Page saved successfully", map[string]interface{}{"page_id": page.ID}))
	return
}

func pagePagesPageUpdateAjax(w http.ResponseWriter, r *http.Request) {
	pageID := strings.Trim(utils.Req(r, "page_id", ""), " ")
	alias := strings.Trim(utils.Req(r, "alias", ""), " ")
	canonicalURL := strings.Trim(utils.Req(r, "canonical_url", ""), " ")
	content := strings.Trim(utils.Req(r, "content", ""), " ")
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

	page := GetEntityStore().EntityFindByID(pageID)

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
	page.SetString("meta_description", metaDescription)
	page.SetString("meta_keywords", metaKeywords)
	page.SetString("meta_robots", metaRobots)
	page.SetString("name", name)
	page.SetString("status", status)
	page.SetString("template_id", templateID)
	page.SetString("handle", handle)
	isOk := page.SetString("title", title)

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

	page := GetEntityStore().EntityFindByID(pageID)

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

	// "PageMetaDescription": pageMetaDescription,
	// "PageMetaKeywords":    pageMetaKeywords,
	// "PageRobots":          pageMetaRobots,

	// Canonical Url
	formGroupCanonicalURL := hb.NewDiv().Attr("class", "form-group")
	formGroupCanonicalURLLabel := hb.NewLabel().HTML("Canonical Url").Attr("class", "form-label")
	formGroupCanonicalURLInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "pageModel.canonicalUrl")
	formGroupCanonicalURL.AddChild(formGroupCanonicalURLLabel).AddChild(formGroupCanonicalURLInput)

	// Meta Description
	formGroupMetaDescription := hb.NewDiv().Attr("class", "form-group")
	formGroupMetaDescriptionLabel := hb.NewLabel().HTML("Meta Description").Attr("class", "form-label")
	formGroupMetaDescriptionInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "pageModel.metaDescription")
	formGroupMetaDescription.AddChild(formGroupMetaDescriptionLabel).AddChild(formGroupMetaDescriptionInput)

	// Meta Keywords
	formGroupMetaKeywords := hb.NewDiv().Attr("class", "form-group")
	formGroupMetaKeywordsLabel := hb.NewLabel().HTML("Meta Keywords").Attr("class", "form-label")
	formGroupMetaKeywordsInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "pageModel.metaKeywords")
	formGroupMetaKeywords.AddChild(formGroupMetaKeywordsLabel).AddChild(formGroupMetaKeywordsInput)

	// Robots
	formGroupMetaRobots := hb.NewDiv().Attr("class", "form-group")
	formGroupMetaRobotsLabel := hb.NewLabel().HTML("Meta Robots").Attr("class", "form-label")
	formGroupMetaRobotsInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "pageModel.metaRobots")
	formGroupMetaRobots.AddChild(formGroupMetaRobotsLabel).AddChild(formGroupMetaRobotsInput)

	// Template
	templateList := GetEntityStore().EntityList("template", 0, 100, "", "id", "asc")
	formGroupTemplate := hb.NewDiv().Attr("class", "form-group")
	formGroupTemplateLabel := hb.NewLabel().HTML("Template").Attr("class", "form-label")
	formGroupTemplateSelect := hb.NewSelect().Attr("class", "form-control").Attr("v-model", "pageModel.templateId")
	formGroupTemplateOptionsEmpty := hb.NewOption().Attr("value", "").HTML("- none -")
	formGroupTemplateSelect.AddChild(formGroupTemplateOptionsEmpty)
	for _, template := range templateList {
		formGroupTemplateOptionsTemplate := hb.NewOption().Attr("value", template.ID).HTML(template.GetString("name", "n/a"))
		formGroupTemplateSelect.AddChild(formGroupTemplateOptionsTemplate)
	}
	formGroupTemplate.AddChild(formGroupTemplateLabel).AddChild(formGroupTemplateSelect)

	// Title
	formGroupTitle := hb.NewDiv().Attr("class", "form-group")
	formGroupTitleLabel := hb.NewLabel().HTML("Title").Attr("class", "form-label")
	formGroupTitleInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "pageModel.title")
	formGroupTitle.AddChild(formGroupTitleLabel).AddChild(formGroupTitleInput)

	formGroupContent := hb.NewDiv().Attr("class", "form-group")
	formGroupContentLabel := hb.NewLabel().HTML("Content").Attr("class", "form-label")
	formGroupContentInput := hb.NewTextArea().Attr("class", "form-control CodeMirror").Attr("v-model", "pageModel.content")
	formGroupContent.AddChild(formGroupContentLabel)
	formGroupContent.AddChild(formGroupContentInput)

	tabContentContent.AddChild(formGroupTitle).AddChild(formGroupContent)
	tabContentSeo.AddChild(formGroupAlias).AddChild(formGroupMetaDescription).AddChild(formGroupMetaKeywords).AddChild(formGroupMetaRobots).AddChild(formGroupCanonicalURL)
	tabContentSettings.AddChild(formGroupStatus).AddChild(formGroupTemplate).AddChild(formGroupName)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))
	container.AddChild(tabNavigation)
	container.AddChild(tabContent)

	h := container.ToHTML()

	alias := page.GetString("alias", "")
	content := page.GetString("content", "")
	name := page.GetString("name", "")
	status := page.GetString("status", "")
	templateID := page.GetString("template_id", "")
	title := page.GetString("title", "")
	metaDescription := page.GetString("meta_description", "")
	metaKeywords := page.GetString("meta_keywords", "")
	metaRobots := page.GetString("meta_robots", "")
	canonicalURL := page.GetString("canonical_url", "")

	canonicalURLJSON, _ := json.Marshal(canonicalURL)
	contentJSON, _ := json.Marshal(content)
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
var templateId = ` + string(templateIDJSON) + `;
const PageUpdate = {
	data() {
		return {
			pageModel:{
				pageId: pageId,
				alias: alias,
				canonicalUrl:canonicalUrl,
				content: content,
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

	pages := GetEntityStore().EntityList("page", 0, 200, "", "id", "asc")

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
		name := page.GetString("name", "n/a")
		alias := page.GetString("alias", "n/a")
		status := page.GetString("status", "n/a")
		buttonEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Attr("class", "btn btn-primary btn-sm").Attr("v-on:click", "pageEdit('"+page.ID+"')").Attr("style", "margin-right:5px")
		buttonTrash := hb.NewButton().HTML("Trash").Attr("type", "button").Attr("class", "btn btn-danger btn-sm").Attr("v-on:click", "showPageTrashModal('"+page.ID+"')")

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
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

// pagePagesPageTrashAjax - moves the template to the trash
func pagePagesPageTrashAjax(w http.ResponseWriter, r *http.Request) {
	pageID := strings.Trim(utils.Req(r, "page_id", ""), " ")

	if pageID == "" {
		api.Respond(w, r, api.Error("Page ID is required"))
		return
	}

	page := GetEntityStore().EntityFindByID(pageID)

	if page == nil {
		api.Respond(w, r, api.Error("Page NOT FOUND with ID "+pageID))
		return
	}

	isOk := GetEntityStore().EntityTrash(pageID)

	if isOk == false {
		api.Respond(w, r, api.Error("Page failed to be moved to trash"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Page moved to trash successfully", map[string]interface{}{"page_id": page.ID}))
	return
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
