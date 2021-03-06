package cms

import (
	"log"
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/utils"
)

func pageWidgetsWidgetCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(utils.Req(r, "name", ""), " ")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	widget := GetEntityStore().EntityCreate("widget")

	if widget == nil {
		api.Respond(w, r, api.Error("Widget failed to be created"))
		return
	}

	widget.SetString("name", name)

	api.Respond(w, r, api.SuccessWithData("Widget saved successfully", map[string]interface{}{"widget_id": widget.ID}))
	return
}

func pageWidgetsWidgetManager(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	log.Println(endpoint)

	header := cmsHeader(endpoint)
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathWidgetsWidgetManager): "Widgets",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "widget-manager")
	heading := hb.NewHeading1().HTML("Widget Manager")
	button := hb.NewButton().HTML("New widget").Attr("class", "btn btn-success float-end").Attr("v-on:click", "showWidgetCreateModal")
	heading.AddChild(button)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))

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

	widgets := GetEntityStore().EntityList("widget", 0, 200, "", "id", "asc")

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
		name := widget.GetString("name", "n/a")
		status := widget.GetString("status", "n/a")
		buttonEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Attr("class", "btn btn-primary").Attr("v-on:click", "widgetEdit('"+widget.ID+"')")

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

	webwidget := Webpage("Widget Manager", h)
	webwidget.AddScript(inlineScript)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webwidget.ToHTML()))
}

func pageWidgetsWidgetUpdate(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	log.Println(endpoint)

	widgetID := utils.Req(r, "widget_id", "")
	if widgetID == "" {
		api.Respond(w, r, api.Error("Widget ID is required"))
		return
	}

	widget := GetEntityStore().EntityFindByID(widgetID)

	if widget == nil {
		api.Respond(w, r, api.Error("Widget NOT FOUND with ID "+widgetID))
		return
	}

	header := cmsHeader(r.Context().Value(keyEndpoint).(string))
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathWidgetsWidgetManager):                           "Widgets",
		(endpoint + "?path=" + PathWidgetsWidgetUpdate + "&widget_id=" + widgetID): "Edit widget",
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
	code := hb.NewCode().AddChild(hb.NewPRE().HTML(`&lt;!-- START: Widget: ` + widget.GetString("name", "") + ` -->
[[BLOCK_` + widget.ID + `]]
&lt;!-- END: Widget: ` + widget.GetString("name", "") + ` -->`))
	paragraphUsage.AddChild(code)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))
	container.AddChild(formGroupStatus).AddChild(formGroupName).AddChild(formGroupContent)
	container.AddChild(paragraphUsage)

	h := container.ToHTML()

	name := widget.GetString("name", "")
	statusAttribute := GetEntityStore().AttributeFind(widget.ID, "status")
	status := ""
	if statusAttribute != nil {
		status = statusAttribute.GetString()
	}
	contentAttribute := GetEntityStore().AttributeFind(widget.ID, "content")
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

	webwidget := Webpage("Edit Widget", h)
	webwidget.AddScript(inlineScript)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webwidget.ToHTML()))
}

func pageWidgetsWidgetUpdateAjax(w http.ResponseWriter, r *http.Request) {
	widgetID := strings.Trim(utils.Req(r, "widget_id", ""), " ")
	content := strings.Trim(utils.Req(r, "content", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	handle := strings.Trim(utils.Req(r, "handle", ""), " ")

	if widgetID == "" {
		api.Respond(w, r, api.Error("Widget ID is required"))
		return
	}

	widget := GetEntityStore().EntityFindByID(widgetID)

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
	isOk := widget.SetString("status", status)

	if isOk == false {
		api.Respond(w, r, api.Error("Widget failed to be updated"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Widget saved successfully", map[string]interface{}{"widget_id": widget.ID}))
	return
}
