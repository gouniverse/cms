package cms

import (
	"net/http"

	"github.com/gouniverse/api"
	"github.com/gouniverse/bs"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/icons"
	"github.com/gouniverse/responses"
	"github.com/gouniverse/utils"
)

func (m UiManager) MenuUpdate(w http.ResponseWriter, r *http.Request) {
	menuID := utils.Req(r, "menu_id", "")
	if menuID == "" {
		api.Respond(w, r, api.Error("Menu ID is required"))
		return
	}

	menu, _ := m.entityStore.EntityFindByID(menuID)

	if menu == nil {
		api.Respond(w, r, api.Error("Menu NOT FOUND with ID "+menuID))
		return
	}

	header := m.cmsHeader(m.endpoint)
	breadcrumbs := m.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  m.endpoint,
			Name: "Home",
		},
		{
			URL:  (m.endpoint + "?path=" + m.pathMenusMenuManager),
			Name: "Menus",
		},
		{
			URL:  (m.endpoint + "?path=" + m.pathMenusMenuUpdate + "&menu_id=" + menuID),
			Name: "Edit menu",
		},
	})

	container := hb.NewDiv().Class("container").Attr("id", "menu-update")
	heading := hb.NewHeading1().HTML("Edit Menu")
	button := hb.NewButton().AddChild(hb.NewHTML(icons.BootstrapCheckCircle+" ")).HTML("Save").Class("btn btn-success float-end").Attr("v-on:click", "menuSave")
	heading.AddChild(button)

	formGroupStatus := hb.NewDiv().Class("form-group")
	formGroupStatusLabel := hb.NewLabel().HTML("Status").Class("form-label")
	formGroupStatusSelect := hb.NewSelect().Class("form-select").Attr("v-model", "menuModel.status")
	formGroupOptionsActive := hb.NewOption().Attr("value", "active").HTML("Active")
	formGroupOptionsInactive := hb.NewOption().Attr("value", "inactive").HTML("Inactive")
	formGroupOptionsTrash := hb.NewOption().Attr("value", "trash").HTML("Trash")
	formGroupStatus.AddChild(formGroupStatusLabel)
	formGroupStatus.AddChild(formGroupStatusSelect.AddChild(formGroupOptionsActive).AddChild(formGroupOptionsInactive).AddChild(formGroupOptionsTrash))

	formGroupName := hb.NewDiv().Class("form-group")
	formGroupNameLabel := hb.NewLabel().HTML("Name").Class("form-label")
	formGroupNameInput := hb.NewInput().Class("form-control").Attr("v-model", "menuModel.name")
	formGroupName.AddChild(formGroupNameLabel)
	formGroupName.AddChild(formGroupNameInput)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrumbs))
	container.AddChild(formGroupStatus).AddChild(formGroupName)

	h := container.ToHTML()

	name, _ := menu.GetString("name", "")
	status, _ := menu.GetString("status", "")
	content, _ := menu.GetString("content", "")

	inlineScript := `
var menuUpdateUrl = "` + m.endpoint + `?path=menus/menu-update-ajax";
var menuId = "` + menuID + `";
var name = "` + name + `";
var status = "` + status + `";
var content = "` + content + `";
const MenuUpdate = {
	data() {
		return {
			menuModel:{
				menuId: menuId,
				name: name,
				status: status,
		    }
		}
	},
	methods: {
		menuSave(){
			var content = this.menuModel.content;
			var name = this.menuModel.name;
			var menuId = this.menuModel.menuId;
			var status = this.menuModel.status;
			
			$.post(menuUpdateUrl, {
				menu_id:menuId,
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
					title: 'Menu saved',
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
Vue.createApp(MenuUpdate).mount('#menu-update')
	`

	if m.funcLayout("") != "" {
		out := hb.NewWrap().Children([]hb.TagInterface{
			hb.NewHTML(h),
			hb.NewScriptURL(cdn.Jquery_3_6_4()),
			hb.NewScriptURL(cdn.VueJs_3()),
			hb.NewScriptURL(cdn.Sweetalert2_10()),
			hb.NewScript(inlineScript),
		}).ToHTML()
		responses.HTMLResponse(w, r, m.funcLayout(out))
		return
	}

	webpage := m.webpageComplete("Edit Menu", h)
	webpage.AddScript(inlineScript)
	responses.HTMLResponse(w, r, webpage.ToHTML())
}
