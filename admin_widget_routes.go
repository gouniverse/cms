package cms

import (
	"log"
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/bs"
	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/responses"
	"github.com/gouniverse/utils"
)

func (cms Cms) pageWidgetsWidgetCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(utils.Req(r, "name", ""), " ")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	widget, err := cms.EntityStore.EntityCreate("widget")

	if err != nil {
		api.Respond(w, r, api.Error("Widget failed to be created: "+err.Error()))
		return
	}

	if widget == nil {
		api.Respond(w, r, api.Error("Widget failed to be created"))
		return
	}

	widget.SetString("name", name)

	api.Respond(w, r, api.SuccessWithData("Widget saved successfully", map[string]interface{}{"widget_id": widget.ID()}))
}

func (cms Cms) pageWidgetsWidgetManager(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	// log.Println(endpoint)

	header := cms.cmsHeader(endpoint)
	breadcrumbs := cms.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  endpoint,
			Name: "Home",
		},
		{
			URL:  (endpoint + "?path=" + PathWidgetsWidgetManager),
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

	widgets, err := cms.EntityStore.EntityList(entitystore.EntityQueryOptions{
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
var widgetCreateUrl = "` + endpoint + `?path=widgets/widget-create-ajax"
var widgetUpdateUrl = "` + endpoint + `?path=widgets/widget-update"
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

	webpage := Webpage("Widget Manager", h)
	webpage.AddScript(inlineScript)

	responses.HTMLResponse(w, r, cms.funcLayout(webpage.ToHTML()))
}

func (cms Cms) pageWidgetsWidgetUpdate(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	log.Println(endpoint)

	widgetID := utils.Req(r, "widget_id", "")
	if widgetID == "" {
		api.Respond(w, r, api.Error("Widget ID is required"))
		return
	}

	widget, _ := cms.EntityStore.EntityFindByID(widgetID)

	if widget == nil {
		api.Respond(w, r, api.Error("Widget NOT FOUND with ID "+widgetID))
		return
	}

	header := cms.cmsHeader(r.Context().Value(keyEndpoint).(string))
	breadcrumbs := cms.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  endpoint,
			Name: "Home",
		},
		{
			URL:  (endpoint + "?path=" + PathWidgetsWidgetManager),
			Name: "Widgets",
		},
		{
			URL:  (endpoint + "?path=" + PathWidgetsWidgetUpdate + "&widget_id=" + widgetID),
			Name: "Edit widget",
		},
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "widget-update")
	heading := hb.NewHeading1().HTML("Edit Widget")
	button := hb.NewButton().HTML("Save").Attr("class", "btn btn-success float-end").Attr("v-on:click", "widgetSave")
	heading.AddChild(button)

	formGroupStatus := hb.NewDiv().Attr("class", "form-group")
	formGroupStatusLabel := hb.NewLabel().HTML("Status").Attr("class", "form-label")
	formGroupStatusSelect := hb.NewSelect().Attr("class", "form-select").Attr("v-model", "widgetModel.status")
	formGroupOptionsActive := hb.NewOption().Attr("value", "active").HTML("Active")
	formGroupOptionsInactive := hb.NewOption().Attr("value", "inactive").HTML("Inactive")
	formGroupOptionsTrash := hb.NewOption().Attr("value", "trash").HTML("Trash")
	formGroupStatus.AddChild(formGroupStatusLabel)
	formGroupStatus.AddChild(formGroupStatusSelect.AddChild(formGroupOptionsActive).AddChild(formGroupOptionsInactive).AddChild(formGroupOptionsTrash))

	formGroupName := hb.NewDiv().Attr("class", "form-group")
	formGroupNameLabel := hb.NewLabel().HTML("Name").Attr("class", "form-label")
	formGroupNameInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "widgetModel.name")
	formGroupName.AddChild(formGroupNameLabel)
	formGroupName.AddChild(formGroupNameInput)

	formGroupContent := hb.NewDiv().Attr("class", "form-group")
	formGroupContentLabel := hb.NewLabel().HTML("Content").Attr("class", "form-label")
	formGroupContentInput := hb.NewTextArea().Attr("class", "form-control").Attr("v-model", "widgetModel.content")
	formGroupContent.AddChild(formGroupContentLabel)
	formGroupContent.AddChild(formGroupContentInput)

	paragraphUsage := hb.NewParagraph().Attr("class", "text-info mt-5").AddChild(hb.NewHTML("To use this widget in your website use the following shortcode:"))
	widgetName, _ := widget.GetString("name", "")
	code := hb.NewCode().AddChild(hb.NewPRE().HTML(`&lt;!-- START: Widget: ` + widgetName + ` -->
[[BLOCK_` + widget.ID() + `]]
&lt;!-- END: Widget: ` + widgetName + ` -->`))
	paragraphUsage.AddChild(code)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrumbs))
	container.AddChild(formGroupStatus).AddChild(formGroupName).AddChild(formGroupContent)
	container.AddChild(paragraphUsage)

	h := container.ToHTML()

	name, _ := widget.GetString("name", "")
	statusAttribute, err := cms.EntityStore.AttributeFind(widget.ID(), "status")
	if err != nil {
		api.Respond(w, r, api.Error("Status failed to be found: "+err.Error()))
		return
	}

	status := ""
	if statusAttribute != nil {
		status = statusAttribute.GetString()
	}
	contentAttribute, err := cms.EntityStore.AttributeFind(widget.ID(), "content")

	if err != nil {
		api.Respond(w, r, api.Error("Content failed to be found: "+err.Error()))
		return
	}

	content := ""
	if contentAttribute != nil {
		content = contentAttribute.GetString()
	}

	inlineScript := `
var widgetUpdateUrl = "` + endpoint + `?path=widgets/widget-update-ajax";
var widgetId = "` + widgetID + `";
var name = "` + name + `";
var status = "` + status + `";
var content = "` + content + `";
const WidgetUpdate = {
	data() {
		return {
			widgetModel:{
				widgetId: widgetId,
				content: content,
				name: name,
				status: status,
		    }
		}
	},
	methods: {
		widgetSave(){
			var content = this.widgetModel.content;
			var name = this.widgetModel.name;
			var widgetId = this.widgetModel.widgetId;
			var status = this.widgetModel.status;
			
			$.post(widgetUpdateUrl, {
				widget_id:widgetId,
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
					title: 'Widget saved',
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
Vue.createApp(WidgetUpdate).mount('#widget-update')
	`

	webpage := Webpage("Edit Widget", h)
	webpage.AddScript(inlineScript)
	responses.HTMLResponse(w, r, cms.funcLayout(webpage.ToHTML()))
}

func (cms Cms) pageWidgetsWidgetUpdateAjax(w http.ResponseWriter, r *http.Request) {
	widgetID := strings.Trim(utils.Req(r, "widget_id", ""), " ")
	content := strings.Trim(utils.Req(r, "content", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	handle := strings.Trim(utils.Req(r, "handle", ""), " ")

	if widgetID == "" {
		api.Respond(w, r, api.Error("Widget ID is required"))
		return
	}

	widget, _ := cms.EntityStore.EntityFindByID(widgetID)

	if widget == nil {
		api.Respond(w, r, api.Error("Widget NOT FOUND with ID "+widgetID))
		return
	}

	if status == "" {
		api.Respond(w, r, api.Error("status is required field"))
		return
	}

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	widget.SetString("content", content)
	widget.SetString("name", name)
	widget.SetString("handle", handle)
	err := widget.SetString("status", status)

	if err != nil {
		api.Respond(w, r, api.Error("Widget failed to be updated: "+err.Error()))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Widget saved successfully", map[string]interface{}{"widget_id": widget.ID()}))
}
