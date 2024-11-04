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

func (m UiManager) MenuManager(w http.ResponseWriter, r *http.Request) {
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
	})

	container := hb.NewDiv().Class("container").Attr("id", "menu-manager")
	heading := hb.NewHeading1().HTML("Menu Manager")
	button := hb.NewButton().HTML("New menu").Class("btn btn-success float-end").Attr("v-on:click", "showMenuCreateModal")
	heading.AddChild(button)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrumbs))

	modal := hb.NewDiv().Attr("id", "ModalMenuCreate").Class("modal fade")
	modalDialog := hb.NewDiv().Class("modal-dialog")
	modalContent := hb.NewDiv().Class("modal-content")
	modalHeader := hb.NewDiv().Class("modal-header").AddChild(hb.NewHeading5().HTML("New Menu"))
	modalBody := hb.NewDiv().Class("modal-body")
	modalBody.AddChild(hb.NewDiv().Class("form-group").AddChild(hb.NewLabel().HTML("Name")).AddChild(hb.NewInput().Class("form-control").Attr("v-model", "menuCreateModel.name")))
	modalFooter := hb.NewDiv().Class("modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Class("btn btn-secondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Create & Continue").Class("btn btn-primary").Attr("v-on:click", "menuCreate"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	container.AddChild(modal)

	menus, err := m.entityStore.EntityList(entitystore.EntityQueryOptions{
		EntityType: m.menuEntityType,
		Offset:     0,
		Limit:      200,
		SortBy:     "id",
		SortOrder:  "asc",
	})

	if err != nil {
		api.Respond(w, r, api.Error("Entity list failed to be retrieved "+err.Error()))
		return
	}

	table := hb.NewTable().Class("table table-responsive table-striped mt-3")
	thead := hb.NewThead()
	tbody := hb.NewTbody()
	table.AddChild(thead).AddChild(tbody)

	tr := hb.NewTR()
	th1 := hb.NewTD().HTML("Name")
	th2 := hb.NewTD().HTML("Status")
	th3 := hb.NewTD().HTML("Menu Items")
	th4 := hb.NewTD().HTML("Actions").Attr("style", "width:1px;")
	thead.AddChild(tr.AddChild(th1).AddChild(th2).AddChild(th3).AddChild(th4))

	for _, menu := range menus {
		name, _ := menu.GetString("name", "n/a")
		status, _ := menu.GetString("status", "n/a")
		buttonEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Class("btn btn-primary").Attr("v-on:click", "menuEdit('"+menu.ID()+"')")
		buttonMenuItemsEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Class("btn btn-primary").Attr("v-on:click", "menuItemsEdit('"+menu.ID()+"')")

		tr := hb.NewTR()
		td1 := hb.NewTD().HTML(name)
		td2 := hb.NewTD().HTML(status)
		td3 := hb.NewTD().AddChild(buttonMenuItemsEdit)
		td4 := hb.NewTD().AddChild(buttonEdit)

		tbody.AddChild(tr.AddChild(td1).AddChild(td2).AddChild(td3).AddChild(td4))
	}
	container.AddChild(table)

	h := container.ToHTML()
	menuCreateURL := m.endpoint + "?path=" + m.pathMenusMenuCreateAjax
	menuUpdateURL := m.endpoint + "?path=" + m.pathMenusMenuUpdate
	menuItemsUpdateURL := m.endpoint + "?path=" + m.pathMenusMenuItemsUpdate

	inlineScript := `
var menuCreateUrl = "` + menuCreateURL + `"
var menuUpdateUrl = "` + menuUpdateURL + `"
var menuItemsUpdateUrl = "` + menuItemsUpdateURL + `"
const MenuManager = {
	data() {
		return {
		  menuCreateModel:{
			  name:""
		  }
		}
	},
	methods: {
        showMenuCreateModal(){
			var modalMenuCreate = new bootstrap.Modal(document.getElementById('ModalMenuCreate'));
			modalMenuCreate.show();
		},
		menuCreate(){
			var name = this.menuCreateModel.name;
		    $.post(menuCreateUrl, {name: name}).done((result)=>{
				if (result.status==="success"){
					var modalMenuCreate = new bootstrap.Modal(document.getElementById('ModalMenuCreate'));
			        modalMenuCreate.hide();

					return location.href = menuUpdateUrl+ "&menu_id=" + result.data.menu_id;
				}
				alert("Failed. " + result.message)
			}).fail((result)=>{
				alert("Failed" + result)
			});
		},
		menuEdit(menuId){
			return location.href = menuUpdateUrl+ "&menu_id=" + menuId;
		},
		menuItemsEdit(menuId){
			return location.href = menuItemsUpdateUrl+ "&menu_id=" + menuId;
		}
	}
};
Vue.createApp(MenuManager).mount('#menu-manager')
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

	webpage := m.webpageComplete("Menu Manager", h)
	webpage.AddScript(inlineScript)
	responses.HTMLResponse(w, r, webpage.ToHTML())
}
