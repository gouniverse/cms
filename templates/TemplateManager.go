package cms

import (
	"net/http"

	"github.com/dracory/bs"
	"github.com/dracory/entitystore"
	"github.com/dracory/hb"
	"github.com/gouniverse/api"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/responses"
)

func (m UiManager) TemplateManager(w http.ResponseWriter, r *http.Request) {
	templates, err := m.entityStore.EntityList(entitystore.EntityQueryOptions{
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

	header := m.cmsHeader(m.endpoint)
	breadcrumbs := m.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  m.endpoint,
			Name: "Home",
		},
		{
			URL:  (m.endpoint + "?path=" + m.pathTemplatesTemplateManager),
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

	container.AddChild(templateTrashModal())
	container.AddChild(templateCreateModal())

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
var templateCreateUrl = "` + m.endpoint + `?path=templates/template-create-ajax";
var templateTrashUrl = "` + m.endpoint + `?path=templates/template-trash-ajax";
var templateUpdateUrl = "` + m.endpoint + `?path=templates/template-update";
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

	if m.funcLayout("") != "" {
		out := hb.NewWrap().Children([]hb.TagInterface{
			hb.NewStyleURL(cdn.JqueryDataTablesCss_1_13_4()),
			hb.NewHTML(h),
			hb.NewScriptURL(cdn.Jquery_3_6_4()),
			hb.NewScriptURL(cdn.VueJs_3()),
			hb.NewScriptURL(cdn.Sweetalert2_10()),
			hb.NewScriptURL(cdn.JqueryDataTablesJs_1_13_4()),
			hb.NewScript(inlineScript),
		}).ToHTML()
		responses.HTMLResponse(w, r, m.funcLayout(out))
		return
	}

	webpage := m.webpageComplete("Template Manager", h)
	webpage.AddStyleURL(cdn.JqueryDataTablesCss_1_13_4())
	webpage.AddScriptURL(cdn.JqueryDataTablesJs_1_13_4())
	webpage.AddScript(inlineScript)
	responses.HTMLResponse(w, r, webpage.ToHTML())
}
