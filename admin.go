package cms

import (
	"context"
	"log"
	"net/http"
	"sinevia/app/helpers"
	"strings"

	"github.com/gouniverse/api"
	hb "github.com/gouniverse/html"
	"github.com/gouniverse/utils"
)

const (
	keyEndpoint             string = "endpoint"
	PathHome                string = "home"
	PathPagesPageManager    string = "pages/page-manager"
	PathPagesPageCreateAjax string = "pages/page-create-ajax"
	PathPagesPageUpdate     string = "pages/page-update"
	PathPagesPageUpdateAjax string = "pages/page-update-ajax"
)

func CmsAdminPage(w http.ResponseWriter, r *http.Request) {
	path := utils.Req(r, "path", "home")

	if path == "" {
		path = "home"
	}

	ctx := context.WithValue(r.Context(), keyEndpoint, r.URL.Path)

	routeFunc := getRoute(path)
	routeFunc(w, r.WithContext(ctx))
	return
}

func getRoute(route string) func(w http.ResponseWriter, r *http.Request) {
	routes := map[string]func(w http.ResponseWriter, r *http.Request){
		PathHome:                pageHome,
		PathPagesPageCreateAjax: pagePagesPageCreateAjax,
		PathPagesPageManager:    pagePagesPageManager,
		PathPagesPageUpdate:     pagePagesPageUpdate,
		PathPagesPageUpdateAjax: pagePagesPageUpdateAjax,
	}
	log.Println(route)
	if val, ok := routes[route]; ok {
		return val
	}
	return routes["home"]
}

func pageHome(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	log.Println(endpoint)

	header := cmsHeader(endpoint)
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "page-manager")
	heading := hb.NewHeading1().HTML("CMS Dashboard")

	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))
	container.AddChild(hb.NewHTML(header))

	h := container.ToHTML()

	webpage := Webpage("Home", h)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

func pagePagesPageCreateAjax(w http.ResponseWriter, r *http.Request) {
	//status := strings.Trim(helpers.Req(r, "status", ""), " ")
	title := strings.Trim(helpers.Req(r, "title", ""), " ")

	if title == "" {
		api.Respond(w, r, api.Error("title is required field"))
		return
	}

	entity := EntityCreateWithAttributes("page", map[string]interface{}{
		"title": title,
		"alias": "/" + utils.Slugify(title+"-"+utils.RandStr(8), '-'),
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
	pageID := strings.Trim(helpers.Req(r, "page_id", ""), " ")
	alias := strings.Trim(helpers.Req(r, "alias", ""), " ")
	content := strings.Trim(helpers.Req(r, "content", ""), " ")
	name := strings.Trim(helpers.Req(r, "name", ""), " ")
	status := strings.Trim(helpers.Req(r, "status", ""), " ")
	title := strings.Trim(helpers.Req(r, "title", ""), " ")

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
		"alias":   alias,
		"content": content,
		"name":    name,
		"status":  status,
		"title":   title,
	})

	if isOk == false {
		api.Respond(w, r, api.Error("Page failed to be updated"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Page saved successfully", map[string]interface{}{"page_id": page.ID}))
	return
}

func cmsBreadcrumbs(breadcrumbs map[string]string) string {
	nav := hb.NewNav().Attr("aria-label", "breadcrumb")
	ol := hb.NewOL().Attr("class", "breadcrumb")
	for k, v := range breadcrumbs {
		li := hb.NewLI().Attr("class", "breadcrumb-item")
		link := hb.NewHyperlink().HTML(v).Attr("href", k)

		li.AddChild(link)

		ol.AddChild(li)
	}
	nav.AddChild(ol)
	return nav.ToHTML()
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

	formGroupTitle := hb.NewDiv().Attr("class", "form-group")
	formGroupTitleLabel := hb.NewLabel().HTML("Title").Attr("class", "form-label")
	formGroupTitleInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "pageModel.title")
	formGroupTitle.AddChild(formGroupTitleLabel)
	formGroupTitle.AddChild(formGroupTitleInput)

	formGroupContent := hb.NewDiv().Attr("class", "form-group")
	formGroupContentLabel := hb.NewLabel().HTML("Content").Attr("class", "form-label")
	formGroupContentInput := hb.NewTextArea().Attr("class", "form-control").Attr("v-model", "pageModel.content")
	formGroupContent.AddChild(formGroupContentLabel)
	formGroupContent.AddChild(formGroupContentInput)

	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))
	container.AddChild(hb.NewHTML(header))
	container.AddChild(formGroupStatus).AddChild(formGroupName).AddChild(formGroupAlias).AddChild(formGroupTitle).AddChild(formGroupContent)

	h := container.ToHTML()

	titleAttribute := EntityAttributeFind(page.ID, "title")
	title := ""
	if titleAttribute != nil {
		title = titleAttribute.GetValue().(string)
	}
	nameAttribute := EntityAttributeFind(page.ID, "name")
	name := ""
	if nameAttribute != nil {
		name = nameAttribute.GetValue().(string)
	}
	statusAttribute := EntityAttributeFind(page.ID, "status")
	status := ""
	if statusAttribute != nil {
		status = statusAttribute.GetValue().(string)
	}
	contentAttribute := EntityAttributeFind(page.ID, "content")
	content := ""
	if contentAttribute != nil {
		content = contentAttribute.GetValue().(string)
	}
	aliasAttribute := EntityAttributeFind(page.ID, "alias")
	alias := ""
	if aliasAttribute != nil {
		alias = aliasAttribute.GetValue().(string)
	}

	inlineScript := `
var pageUpdateUrl = "` + endpoint + `?path=pages/page-update-ajax";
var pageId = "` + pageID + `";
var alias = "` + alias + `";
var name = "` + name + `";
var status = "` + status + `";
var title = "` + title + `";
var content = "` + content + `";
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
		    }
		}
	},
	methods: {
		pageSave(){
			var alias = this.pageModel.alias;
			var content = this.pageModel.content;
			var name = this.pageModel.name;
			var pageId = this.pageModel.pageId;
			var status = this.pageModel.status;
			var title = this.pageModel.title;
			
			$.post(pageUpdateUrl, {
				page_id:pageId,
				alias: alias,
				content: content,
				name: name,
				status: status,
				title: title,
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

	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))
	container.AddChild(hb.NewHTML(header))

	modal := hb.NewDiv().Attr("id", "ModalPageCreate").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("New Page"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewDiv().Attr("class", "form-group").AddChild(hb.NewLabel().HTML("Title")).AddChild(hb.NewInput().Attr("class", "form-control").Attr("v-model", "pageCreateModel.title")))
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
	th1 := hb.NewTD().HTML("Title")
	th2 := hb.NewTD().HTML("Status")
	th3 := hb.NewTD().HTML("Actions").Attr("style", "width:1px;")
	thead.AddChild(tr.AddChild(th1).AddChild(th2).AddChild(th3))

	for _, page := range pages {
		attributeTitle := page.GetAttribute("title")
		title := "n/a"
		if attributeTitle != nil {
			title = attributeTitle.GetValue().(string)
		}
		attributeAlias := page.GetAttribute("alias")
		alias := "n/a"
		if attributeAlias != nil {
			alias = attributeAlias.GetValue().(string)
		}
		attributeStatus := page.GetAttribute("status")
		status := "n/a"
		if attributeStatus != nil {
			status = attributeStatus.GetValue().(string)
		}
		buttonEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Attr("class", "btn btn-primary").Attr("v-on:click", "pageEdit('"+page.ID+"')")

		tr := hb.NewTR()
		td1 := hb.NewTD().HTML(title).AddChild(hb.NewDiv().HTML(alias).Attr("style", "font-size:11px;"))
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
			  title:"Test"
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
			var title = this.pageCreateModel.title;
		    $.post(pageCreateUrl, {title: title}).done((result)=>{
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

func cmsHeader(endpoint string) string {
	linkHome := hb.NewHyperlink().HTML("Dashboard").Attr("href", endpoint+"").Attr("class", "nav-link").ToHTML()
	linkPages := hb.NewHyperlink().HTML("Pages").Attr("href", endpoint+"?path="+PathPagesPageManager).Attr("class", "nav-link").ToHTML()

	h := `<div class="card card-default">
    <div class="card-body" style="padding: 2px;">
        <ul class="nav  nav-pills justify-content-center">
            <li class="nav-item">
                ` + linkHome + `
            </li>
            <li class="nav-item">
                <a class="nav-link" href="#">
                    Templates
                    <span class="badge">1</span>
                </a>
            </li>
            <li class="nav-item">
                ` + linkPages + `
            </li>
            <li class="nav-item">
                <a class="nav-link" href="#">
                    Blocks
                    <span class="badge">9</span>
                </a>
            </li>
            <li class="nav-item">
                <a class="nav-link" href="#">
                    Widgets
                    <span class="badge">3</span>
                </a>
            </li>
            <li class="nav-item">
                <a class="nav-link" href="#">
                    Translations
                    <span class="badge">2</span>
                </a>
            </li>
            <li class="nav-item" class="nav-item">
                <a class="nav-link" href="#" target="_blank">
                    Media
                </a>
            </li>
        </ul>
    </div>
</div>`
	return h
}

// Webpage returns the webpage template for the website
func Webpage(title, content string) *hb.Webpage {
	faviconImgCms := `data:image/x-icon;base64,AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAmzKzAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABEQEAAQERAAEAAQABAAEAAQABAQEBEQABAAEREQEAAAERARARAREAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD//wAA//8AAP//AAD//wAA//8AAP//AAD//wAAi6MAALu7AAC6owAAuC8AAIkjAAD//wAA//8AAP//AAD//wAA`
	app := ""
	app += `APP_ID="lobreid-com";WEBSITE_URL="` + helpers.AppURL() + `"`
	webpage := hb.NewWebpage()
	webpage.SetTitle(title + " | " + helpers.AppName())
	webpage.SetFavicon(faviconImgCms)

	webpage.AddStyleURLs([]string{
		"https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/css/bootstrap.min.css",
	})
	webpage.AddScriptURLs([]string{
		"https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/js/bootstrap.bundle.min.js",
		"http://code.jquery.com/jquery-3.5.1.min.js",
		"https://unpkg.com/vue@next",
		"https://cdn.jsdelivr.net/npm/sweetalert2@9",
	})
	webpage.AddScripts([]string{
		app,
	})
	webpage.AddStyle(`html,body{height:100%;font-family: Ubuntu, sans-serif;}`)
	webpage.AddStyle(`body {
		font-family: "Nunito", sans-serif;
		font-size: 0.9rem;
		font-weight: 400;
		line-height: 1.6;
		color: #212529;
		text-align: left;
		background-color: #f8fafc;
	}`)
	webpage.AddChild(hb.NewHTML(content))
	return webpage
}
