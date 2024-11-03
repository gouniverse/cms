package cms

import (
	"net/http"

	"github.com/gouniverse/api"
	"github.com/gouniverse/bs"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/responses"
)

func (m UiManager) PageManager(w http.ResponseWriter, r *http.Request) {
	UiManager, err := m.entityStore.EntityList(entitystore.EntityQueryOptions{
		EntityType: m.pageEntityType,
		Offset:     0,
		Limit:      200,
		SortBy:     "id",
		SortOrder:  "asc",
	})

	if err != nil {
		api.Respond(w, r, api.Error("Page list failed to be retrieved "+err.Error()))
		return
	}

	header := m.cmsHeader(m.endpoint)
	breadcrumbs := m.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  m.endpoint,
			Name: "Home",
		},
		{
			URL:  (m.endpoint + "?path=" + m.pathPagesPageManager),
			Name: "UiManager",
		},
	})

	container := hb.NewDiv().Class("container").ID("page-manager")
	heading := hb.NewHeading1().HTML("Page Manager")
	button := hb.NewButton().HTML("New page").Attr("class", "btn btn-success float-end").Attr("v-on:click", "showPageCreateModal")
	heading.AddChild(button)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrumbs))

	container.AddChild(pageUiManagerPageCreateModal())
	container.AddChild(pageUiManagerPageTrashModal())

	table := hb.NewTable().Attr("id", "TableUiManager").Attr("class", "table table-responsive table-striped mt-3")
	thead := hb.NewThead()
	tbody := hb.NewTbody()
	table.AddChild(thead).AddChild(tbody)

	tr := hb.NewTR()
	th1 := hb.NewTD().HTML("Name")
	th2 := hb.NewTD().HTML("Status")
	th3 := hb.NewTD().HTML("Actions").Attr("style", "width:150px;")
	thead.AddChild(tr.AddChild(th1).AddChild(th2).AddChild(th3))

	for _, page := range UiManager {
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
var pageCreateUrl = "` + m.endpoint + `?path=pages/page-create-ajax"
var pageTrashUrl = "` + m.endpoint + `?path=pages/page-trash-ajax"
var pageUpdateUrl = "` + m.endpoint + `?path=pages/page-update"
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
				$('#TableUiManager').DataTable({
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

	if m.funcLayout("") != "" {
		out := hb.NewWrap().Children([]hb.TagInterface{
			hb.NewHTML(h),
			hb.NewScriptURL(cdn.Jquery_3_6_4()),
			hb.NewScriptURL(cdn.VueJs_3()),
			hb.NewScriptURL(cdn.Sweetalert2_10()),
			hb.NewScriptURL(cdn.JqueryDataTablesCss_1_13_4()),
			hb.NewScript(inlineScript),
		}).ToHTML()
		responses.HTMLResponse(w, r, m.funcLayout(out))
		return
	}

	webpage := m.webpageComplete("Page Manager", h)
	webpage.AddStyleURL(cdn.JqueryDataTablesCss_1_13_4())
	webpage.AddScriptURL(cdn.JqueryDataTablesJs_1_13_4())
	webpage.AddScript(inlineScript)

	responses.HTMLResponse(w, r, m.funcLayout(webpage.ToHTML()))
}

func pageUiManagerPageTrashModal() *hb.Tag {
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

func pageUiManagerPageCreateModal() *hb.Tag {
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
