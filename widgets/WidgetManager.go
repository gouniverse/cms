package cms

import (
	"net/http"

	"github.com/gouniverse/api"
	"github.com/gouniverse/bs"
	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/responses"
)

func (m UiManager) WidgetManager(w http.ResponseWriter, r *http.Request) {
	header := m.cmsHeader(m.endpoint)
	breadcrumbs := m.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  m.url("", map[string]string{}),
			Name: "Home",
		},
		{
			URL:  m.url(m.pathWidgetsWidgetManager, map[string]string{}),
			Name: "Widgets",
		},
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "widget-manager")
	heading := hb.NewHeading1().HTML("Widget Manager")
	button := hb.NewButton().HTML("New widget").Attr("class", "btn btn-success float-end").Attr("v-on:click", "showWidgetCreateModal")
	heading.AddChild(button)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrumbs))

	modal := hb.NewDiv().Attr("id", "ModalWidgetCreate").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("New Widget"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewDiv().Attr("class", "form-group").AddChild(hb.NewLabel().HTML("Name")).AddChild(hb.NewInput().Attr("class", "form-control").Attr("v-model", "widgetCreateModel.name")))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-prsecondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Create & Continue").Attr("class", "btn btn-primary").Attr("v-on:click", "widgetCreate"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	container.AddChild(modal)

	widgets, err := m.entityStore.EntityList(entitystore.EntityQueryOptions{
		EntityType: "widget",
		Offset:     0,
		Limit:      200,
		SortBy:     "id",
		SortOrder:  "asc",
	})

	if err != nil {
		api.Respond(w, r, api.Error("Widgets failed to be fetched: "+err.Error()))
		return
	}

	table := hb.NewTable().Attr("class", "table table-responsive table-striped mt-3")
	thead := hb.NewThead()
	tbody := hb.NewTbody()
	table.AddChild(thead).AddChild(tbody)

	tr := hb.NewTR()
	th1 := hb.NewTD().HTML("Name")
	th2 := hb.NewTD().HTML("Status")
	th3 := hb.NewTD().HTML("Actions").Attr("style", "width:1px;")
	thead.AddChild(tr.AddChild(th1).AddChild(th2).AddChild(th3))

	for _, widget := range widgets {
		name, _ := widget.GetString("name", "n/a")
		status, _ := widget.GetString("status", "n/a")
		buttonEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Attr("class", "btn btn-primary").Attr("v-on:click", "widgetEdit('"+widget.ID()+"')")

		tr := hb.NewTR()
		td1 := hb.NewTD().HTML(name)
		td2 := hb.NewTD().HTML(status)
		td3 := hb.NewTD().AddChild(buttonEdit)

		tbody.AddChild(tr.AddChild(td1).AddChild(td2).AddChild(td3))
	}
	container.AddChild(table)

	h := container.ToHTML()

	inlineScript := `
var widgetCreateUrl = "` + m.url("widgets/widget-create-ajax", map[string]string{}) + `"
var widgetUpdateUrl = "` + m.url("widgets/widget-update", map[string]string{}) + `"
const WidgetManager = {
	data() {
		return {
		  widgetCreateModel:{
			  name:""
		  }
		}
	},
	methods: {
        showWidgetCreateModal(){
			var modalWidgetCreate = new bootstrap.Modal(document.getElementById('ModalWidgetCreate'));
			modalWidgetCreate.show();
		},
		widgetCreate(){
			var name = this.widgetCreateModel.name;
		    $.post(widgetCreateUrl, {name: name}).done((result)=>{
				if (result.status==="success"){
					var modalWidgetCreate = new bootstrap.Modal(document.getElementById('ModalWidgetCreate'));
			        modalWidgetCreate.hide();

					return location.href = widgetUpdateUrl+ "&widget_id=" + result.data.widget_id;
				}
				alert("Failed. " + result.message)
			}).fail((result)=>{
				alert("Failed" + result)
			});
		},
		widgetEdit(widgetId){
			return location.href = widgetUpdateUrl+ "&widget_id=" + widgetId;
		}
	}
};
Vue.createApp(WidgetManager).mount('#widget-manager')
	`

	if m.funcLayout("") != "" {
		out := hb.NewWrap().Children([]hb.TagInterface{
			hb.NewHTML(h),
			hb.NewScript(inlineScript),
		}).ToHTML()
		responses.HTMLResponse(w, r, m.funcLayout(out))
		return
	}

	webpage := m.webpageComplete("Widget Manager", h)
	webpage.AddScript(inlineScript)

	responses.HTMLResponse(w, r, webpage.ToHTML())
}
