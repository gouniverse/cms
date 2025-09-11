package cms

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/gouniverse/responses"
)

func (m UiManager) WidgetUpdate(w http.ResponseWriter, r *http.Request) {
	widgetID := req.GetStringTrimmed(r, "widget_id")
	if widgetID == "" {
		api.Respond(w, r, api.Error("Widget ID is required"))
		return
	}

	widget, _ := m.entityStore.EntityFindByID(widgetID)

	if widget == nil {
		api.Respond(w, r, api.Error("Widget NOT FOUND with ID "+widgetID))
		return
	}

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
		{
			URL:  m.url(m.pathWidgetsWidgetUpdate, map[string]string{"widget_id": widgetID}),
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
	statusAttribute, err := m.entityStore.AttributeFind(widget.ID(), "status")
	if err != nil {
		api.Respond(w, r, api.Error("Status failed to be found: "+err.Error()))
		return
	}

	status := ""
	if statusAttribute != nil {
		status = statusAttribute.GetString()
	}
	contentAttribute, err := m.entityStore.AttributeFind(widget.ID(), "content")

	if err != nil {
		api.Respond(w, r, api.Error("Content failed to be found: "+err.Error()))
		return
	}

	content := ""
	if contentAttribute != nil {
		content = contentAttribute.GetString()
	}

	inlineScript := `
var widgetUpdateUrl = "` + m.endpoint + `?path=widgets/widget-update-ajax";
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

	if m.funcLayout("") != "" {
		out := hb.NewWrap().Children([]hb.TagInterface{
			hb.NewHTML(h),
			hb.NewScript(inlineScript),
		}).ToHTML()
		responses.HTMLResponse(w, r, m.funcLayout(out))
		return
	}

	webpage := m.webpageComplete("Edit Widget", h)
	webpage.AddScript(inlineScript)
	responses.HTMLResponse(w, r, webpage.ToHTML())
}
