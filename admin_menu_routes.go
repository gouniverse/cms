package cms

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/bs"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/icons"
	"github.com/gouniverse/responses"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

func (cms Cms) pageMenusMenuCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(utils.Req(r, "name", ""), " ")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	menu, err := cms.EntityStore.EntityCreateWithType(ENTITY_TYPE_MENU)

	if err != nil {
		api.Respond(w, r, api.Error("Menu failed to be created "+err.Error()))
		return
	}

	if menu == nil {
		api.Respond(w, r, api.Error("Menu failed to be created"))
		return
	}

	menu.SetString("name", name)

	api.Respond(w, r, api.SuccessWithData("Menu saved successfully", map[string]interface{}{"menu_id": menu.ID()}))
}

func (cms Cms) pageMenusMenuManager(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	//log.Println(endpoint)

	header := cms.cmsHeader(endpoint)
	breadcrumbs := cms.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  endpoint,
			Name: "Home",
		},
		{
			URL:  (endpoint + "?path=" + PathMenusMenuManager),
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

	menus, err := cms.EntityStore.EntityList(entitystore.EntityQueryOptions{
		EntityType: ENTITY_TYPE_MENU,
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
	menuCreateURL := endpoint + "?path=" + PathMenusMenuCreateAjax
	menuUpdateURL := endpoint + "?path=" + PathMenusMenuUpdate
	menuItemsUpdateURL := endpoint + "?path=" + PathMenusMenuItemsUpdate

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

	if cms.funcLayout("") != "" {
		out := hb.NewWrap().Children([]hb.TagInterface{
			hb.NewHTML(h),
			hb.NewScriptURL(cdn.Jquery_3_6_4()),
			hb.NewScriptURL(cdn.VueJs_3()),
			hb.NewScriptURL(cdn.Sweetalert2_10()),
			hb.NewScript(inlineScript),
		}).ToHTML()
		responses.HTMLResponse(w, r, cms.funcLayout(out))
		return
	}

	webpage := WebpageComplete("Menu Manager", h)
	webpage.AddScript(inlineScript)
	responses.HTMLResponse(w, r, webpage.ToHTML())
}

func (cms Cms) pageMenusMenuUpdate(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	// log.Println(endpoint)

	menuID := utils.Req(r, "menu_id", "")
	if menuID == "" {
		api.Respond(w, r, api.Error("Menu ID is required"))
		return
	}

	menu, _ := cms.EntityStore.EntityFindByID(menuID)

	if menu == nil {
		api.Respond(w, r, api.Error("Menu NOT FOUND with ID "+menuID))
		return
	}

	header := cms.cmsHeader(r.Context().Value(keyEndpoint).(string))
	breadcrumbs := cms.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  endpoint,
			Name: "Home",
		},
		{
			URL:  (endpoint + "?path=" + PathMenusMenuManager),
			Name: "Menus",
		},
		{
			URL:  (endpoint + "?path=" + PathMenusMenuUpdate + "&menu_id=" + menuID),
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
var menuUpdateUrl = "` + endpoint + `?path=menus/menu-update-ajax";
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

	if cms.funcLayout("") != "" {
		out := hb.NewWrap().Children([]hb.TagInterface{
			hb.NewHTML(h),
			hb.NewScriptURL(cdn.Jquery_3_6_4()),
			hb.NewScriptURL(cdn.VueJs_3()),
			hb.NewScriptURL(cdn.Sweetalert2_10()),
			hb.NewScript(inlineScript),
		}).ToHTML()
		responses.HTMLResponse(w, r, cms.funcLayout(out))
		return
	}

	webpage := WebpageComplete("Edit Menu", h)
	webpage.AddScript(inlineScript)
	responses.HTMLResponse(w, r, webpage.ToHTML())
}

func (cms Cms) pageMenusMenuUpdateAjax(w http.ResponseWriter, r *http.Request) {
	menuID := strings.Trim(utils.Req(r, "menu_id", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	handle := strings.Trim(utils.Req(r, "handle", ""), " ")

	if menuID == "" {
		api.Respond(w, r, api.Error("Menu ID is required"))
		return
	}

	menu, _ := cms.EntityStore.EntityFindByID(menuID)

	if menu == nil {
		api.Respond(w, r, api.Error("Menu NOT FOUND with ID "+menuID))
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

	menu.SetString("name", name)
	menu.SetString("handle", handle)
	err := menu.SetString("status", status)

	if err != nil {
		api.Respond(w, r, api.Error("Menu failed to be updated "+err.Error()))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Menu saved successfully", map[string]interface{}{"menu_id": menu.ID()}))
}

func getChildren(data []map[string]interface{}, parentID string) []map[string]interface{} {
	children := []map[string]interface{}{}
	//sequences := []string{}
	for _, node := range data {
		nodeParentID := ""
		//sequence := ""
		if keyExists(node, "parent_id") {
			nodeParentID = node["parent_id"].(string)
		}
		//if keyExists(node, "sequence") {
		//sequence = node["sequence"].(string)
		//}
		if nodeParentID == parentID {
			//sequences = append(sequences, sequence)
			children = append(children, node)
		}
	}

	sort.SliceStable(children, func(i, j int) bool {
		a, _ := strconv.ParseInt(children[i]["sequence"].(string), 10, 64)
		b, _ := strconv.ParseInt(children[j]["sequence"].(string), 10, 64)
		return a < b
	})

	return children
}

func buildTreeFromData(data []map[string]interface{}, parentID string) []map[string]interface{} {
	out := []map[string]interface{}{}

	roots := getChildren(data, parentID)

	for _, root := range roots {
		rootChildren := []map[string]interface{}{}
		rootID := root["id"].(string)
		children := getChildren(data, rootID)
		for childIndex, child := range children {
			childID := child["id"].(string)
			childrenTrees := buildTreeFromData(data, childID)
			rootChildren = append(rootChildren, childrenTrees...)
			children[childIndex]["children"] = rootChildren
		}
		root["children"] = children
		out = append(out, root)
	}

	return out
}

func (cms Cms) buildTree(menuID string) []map[string]interface{} {
	menuitems, err := cms.EntityStore.EntityListByAttribute(ENTITY_TYPE_MENUITEM, "menu_id", menuID)

	if err != nil {
		log.Panicln("Menu items failed to be retrieved " + err.Error())
		return nil
	}

	nodeList := []map[string]interface{}{}
	for _, menuitem := range menuitems {
		itemID := menuitem.ID()
		itemName, _ := menuitem.GetString("name", "n/a")
		parentID, _ := menuitem.GetString("parent_id", "")
		sequence, _ := menuitem.GetString("sequence", "")
		target, _ := menuitem.GetString("target", "")
		url, _ := menuitem.GetString("url", "")
		pageID, _ := menuitem.GetString("page_id", "")
		item := map[string]interface{}{
			"id":        itemID,
			"parent_id": parentID,
			"sequence":  sequence,
			"name":      itemName,
			"page_id":   pageID,
			"url":       url,
			"target":    target,
		}
		nodeList = append(nodeList, item)
	}

	tree := buildTreeFromData(nodeList, "")

	return tree
}

func (cms Cms) pageMenusMenuItemsFetchAjax(w http.ResponseWriter, r *http.Request) {
	menuID := strings.Trim(utils.Req(r, "menu_id", ""), " ")

	if menuID == "" {
		api.Respond(w, r, api.Error("Menu ID is required"))
		return
	}

	menu, _ := cms.EntityStore.EntityFindByID(menuID)

	if menu == nil {
		api.Respond(w, r, api.Error("Menu NOT FOUND with ID "+menuID))
		return
	}

	tree := cms.buildTree(menuID)

	api.Respond(w, r, api.SuccessWithData("Menu items found successfully", map[string]interface{}{
		"menu_id":   menu.ID(),
		"menuitems": tree,
	}))
}

func (cms Cms) pageMenusMenuItemsUpdate(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	// log.Println(endpoint+"")

	menuID := utils.Req(r, "menu_id", "")
	if menuID == "" {
		api.Respond(w, r, api.Error("Menu ID is required"))
		return
	}

	menu, _ := cms.EntityStore.EntityFindByID(menuID)

	if menu == nil {
		api.Respond(w, r, api.Error("Menu NOT FOUND with ID "+menuID))
		return
	}

	menuName, _ := menu.GetString("name", "")

	header := cms.cmsHeader(r.Context().Value(keyEndpoint).(string))
	breadcrumbs := cms.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  endpoint,
			Name: "Home",
		},
		{
			URL:  (endpoint + "?path=" + PathMenusMenuManager),
			Name: "Menus",
		},
		{
			URL:  (endpoint + "?path=" + PathMenusMenuUpdate + "&menu_id=" + menuID),
			Name: "Menu",
		},
		{
			URL:  "#",
			Name: "Edit menu items",
		},
	})

	heading := hb.NewHeading1().HTML("Edit Menu Items for Menu '" + menuName + "'")
	//button := hb.NewButton().AddChild(hb.NewHTML(icons.BootstrapCheckCircle+" ")).HTML("Save").Class( "btn btn-success float-end").Attr("v-on:click", "menuSave")
	//heading.AddChild(button)

	backURL := endpoint + "?path=" + PathMenusMenuManager
	buttonCancel := hb.NewHyperlink().
		Class("btn btn-info").
		Href(backURL).
		Child(hb.NewHTML(icons.BootstrapChevronLeft)).
		HTML(" Cancel")
	buttonAddNode := hb.NewButton().
		Attr("v-on:click", "menuItemAdd").
		Class("btn btn-success float-end").
		ID("ButtonNewMenuItem").
		Attr("disabled", "disabled").
		Child(hb.NewHTML(icons.BootstrapPlusCircle)).
		HTML(" Add item")
	buttonSaveNode := hb.NewButton().
		Attr("v-on:click", "menuItemsSave").
		Class("btn btn-success").
		ID("ButtonSaveMenuItems").
		Attr("disabled", "disabled").
		Child(hb.NewHTML(icons.BootstrapBoxArrowInDown)).
		HTML(" Save")

	actionsCard := hb.NewDiv().
		Class("card box-primary").
		Child(
			hb.NewDiv().
				Class("card-header with-border").
				Child(buttonCancel).
				Child(buttonAddNode).
				Child(buttonSaveNode)).
		Child(
			hb.NewDiv().
				Class("card-body with-border").
				Child(
					hb.NewDiv().
						ID("tree1")))

	modal := hb.NewDiv().
		ID("ModalItemUpdate").
		Class("modal fade").
		Child(hb.NewDiv().
			Class("modal-dialog").
			Child(hb.NewDiv().
				Class("modal-content").
				Child(hb.NewDiv().
					Class("modal-header").
					Child(hb.NewHeading5().HTML("Edit Menu Item"))).
				Child(hb.NewDiv().
					Class("modal-body").
					Child(hb.NewDiv().
						Class("form-group").
						Style(`margin-top:10px;`).
						Child(hb.NewLabel().HTML("Title")).
						Child(hb.NewInput().Class("form-control").Attr("v-model", "menuItemUpdateModel.title")).
						Child(hb.NewDiv().Class("text-info").HTML("Title to display"))).
					Child(hb.NewDiv().
						Class("form-group").
						Style(`margin-top:10px;`).
						Child(hb.NewLabel().HTML("Page")).
						Child(hb.NewSelect().
							Class("form-select").Attr("v-model", "menuItemUpdateModel.pageId").
							Child(hb.NewTemplate().
								Attr("v-for", "dropdown in pagesDropdownList").
								Child(hb.NewOption().Attr("v-bind:value", "dropdown.key").Attr("v-html", "dropdown.value"))),
						).
						Child(hb.NewDiv().Class("text-info").HTML("Page to link to"))).
					Child(hb.NewDiv().
						Class("form-group").
						Style(`margin-top:10px;`).
						Child(hb.NewLabel().HTML("URL")).
						Child(hb.NewInput().Class("form-control").Attr("v-model", "menuItemUpdateModel.url")).
						Child(hb.NewDiv().Class("text-info").HTML("URL to link to (if page not set)"))).
					Child(hb.NewDiv().
						Class("form-group").
						Style(`margin-top:10px;`).
						Child(hb.NewLabel().HTML("Target")).
						Child(hb.NewInput().Class("form-control").Attr("v-model", "menuItemUpdateModel.target")).
						Child(hb.NewDiv().Class("text-info").HTML("Where to open the menu item. Can be one of _blank |_self (default) | _parent | _top | frame name"))).
					Child(hb.NewDiv().
						Class("modal-footer").
						Child(hb.NewButton().HTML("Close").Class("btn btn-secondary").Attr("data-bs-dismiss", "modal")).
						Child(hb.NewButton().HTML("Update").Class("btn btn-primary").Attr("v-on:click", "menuItemUpdate").Attr("data-bs-dismiss", "modal"))))))

	container := hb.NewDiv().
		Class("container").
		ID("menu-items-update").
		Child(hb.NewHTML(header)).
		Child(heading).
		Child(hb.NewHTML(breadcrumbs)).
		Child(actionsCard).
		Child(modal)

	pagesDropdownList, errorMessage := cms.pageMenusMenuItemsPagesDropdownList()

	if errorMessage != "" {
		api.Respond(w, r, api.Error(errorMessage))
		return
	}

	menuItemsUpdateURL := endpoint + "?path=" + PathMenusMenuItemsUpdateAjax
	menuItemsFetchURL := endpoint + "?path=" + PathMenusMenuItemsFetchAjax
	pagesDropdownJSON, _ := utils.ToJSON(pagesDropdownList)
	inlineScript := `
var menuItemsSaveUrl = "` + menuItemsUpdateURL + `";
var menuItemsFetchUrl = "` + menuItemsFetchURL + `";
var menuId = "` + menuID + `";
var pagesDropdownList = ` + pagesDropdownJSON + `;
const MenuItemsUpdate = {
	data() {
		return {
			menuId: menuId,
			menuitems:[],
			pagesDropdownList: pagesDropdownList,
			menuItemUpdateModel:{
				menuItemId:null,
				pageId:null,
				url:null,
				target:null,
			}
		}
	},
	created(){
		setTimeout(()=>{
			this.menuItemsFetch();
		}, 1000);
	},
	methods: {
		menuItemAdd() {
			var uniqueId = Math.random().toString(36).substring(2) + Date.now().toString(36);
			$('#tree1').tree('appendNode', {
				name: ('new_node ' + uniqueId),
				id: uniqueId
			}, null);
			this.menuItemButtonsAddEvents();
		},
		menuItemUpdate(){
			var title = this.menuItemUpdateModel.title;
			var pageId = this.menuItemUpdateModel.pageId;
			var url = this.menuItemUpdateModel.url;
			var menuItemId = this.menuItemUpdateModel.menuItemId;
			var target = this.menuItemUpdateModel.target;
			var node = $('#tree1').tree('getNodeById', menuItemId);
			$('#tree1').tree('updateNode', node, {
				name: title,
				title: title,
				page_id: pageId,
				url: url,
				target: target
			});
			this.menuItemsLoad();			
			var modalMenuItemUpdate = new bootstrap.Modal(document.getElementById('ModalItemUpdate'));
			modalMenuItemUpdate.hide();
		},
		menuItemsSave() {
			var data = $('#tree1').tree('toJson'); // JSON string
			$('#ButtonSaveMenuItems').html('Saving menu items to server ...');
			$.post(menuItemsSaveUrl, { menu_id: this.menuId, data: data }).then((response) => {
				if (response.status === 'success') {
					$.notify("Saving items successful. Reloading from the server...", "success");
				} else {
					$.notify("Saving items failed. Reloading from the server...", "error");
				}
				$('#tree1').tree('destroy');
				this.menuItemsFetch();
			}).fail(() => {
				this.menuItemsFetch();
			});
		},
		menuItemButtonsAddEvents() {
			var self = this;
			$('.menu-item-edit-button').click(function () {
				var id = $(this).data('node-id');
				var node = $('#tree1').tree('getNodeById', id);
				console.log(node);
				self.showMenuItemUpdateModal(id, node.name, node.page_id, node.url, node.target);
			});
			$('.menu-item-delete-button').click(function () {
				var id = $(this).data('node-id');
				var node = $('#tree1').tree('getNodeById', id);
				var isConfirmed = confirm('Are you sure you want to delete node: "' + node.name + '"');
				if (isConfirmed === true) {
					console.log(node);
					$('#tree1').tree('removeNode', node);
					self.menuItemButtonsAddEvents();
					return true;
				}
				return false;
			});
		},
		menuItemsLoad() {
			$('#tree1').tree({
				data: this.menuitems,
				dragAndDrop: true,
				autoOpen: true,
				saveState: true,
				onCreateLi: function (node, $li) {
					// Append a link to the jqtree-element div.
					// The link has an url '#node-[id]' and a data property 'node-id'.
					var html = '';
					html += '<button class="menu-item-delete-button float-end btn btn-danger btn-sm" data-node-id="' + node.id + '">delete</button>';
					html += '<button class="menu-item-edit-button float-end btn btn-info btn-sm" data-node-id="' + node.id + '">edit</button>';
					$li.find('.jqtree-element').append(html);
				}
			});
			this.menuItemButtonsAddEvents();
		},
		menuItemsFetch(){
			//menuItems = [];
			// this.menuItemsLoad();
			$('#ButtonSaveMenuItems').html('<i class="fas fa-save"></i> Save nodes').prop('disabled', false);
			$('#ButtonNewMenuItem').prop('disabled', false);
			// return;

			$('#ButtonSaveMenuItems').html('Loading menu items from the server ...').prop('disabled', true);
			$('#ButtonNewMenuItem').prop('disabled', true);
			$.get(menuItemsFetchUrl, { menu_id: menuId }).then((response) => {
				if (response.status === 'success') {
					this.menuitems = [];
					this.menuitems = response.data.menuitems;
					this.menuItemsLoad();
					$('#ButtonSaveMenuItems').html('<i class="fas fa-save"></i> Save menu items').prop('disabled', false);
					$('#ButtonNewMenuItem').prop('disabled', false);
				} else {
					$.notify("Fetching menu items unsuccessful. Please refresh the page.", "error");
				}
			}).fail(() => {
				this.menuItems = [];
				this.menuItemsLoad();
				$.notify("Fetching menu items FAILED. Please refresh the page.", "error");
				$('#ButtonSaveMenuItems').html('fetching  menu items FAILED. Please, refresh page').prop('disabled', true);
				$('#ButtonNewMenuItem').prop('disabled', true);
			});
		},
		showMenuItemUpdateModal(nodeId, nodeName, nodePageId, nodeUrl, nodeTarget){
			this.menuItemUpdateModel.menuItemId = nodeId;
			this.menuItemUpdateModel.pageId = nodePageId;
			this.menuItemUpdateModel.title = nodeName;
			this.menuItemUpdateModel.url = nodeUrl;
			this.menuItemUpdateModel.target = nodeTarget;
			var modalMenuItemUpdate = new bootstrap.Modal(document.getElementById('ModalItemUpdate'));
			modalMenuItemUpdate.show();
		},
		// menuItemsSave(){
		// 	var menuId = this.menuModel.menuId;
		//	
		// 	$.post(menuItemsUpdateUrl, {
		// 		menu_id:menuId,
		// 	}).done((response)=>{
		// 		if (response.status !== "success") {
		// 			return Swal.fire({
		// 				icon: 'error',
		// 				title: 'Oops...',
		// 				text: response.message,
		// 			});
		// 		}

		// 		return Swal.fire({
		// 			icon: 'success',
		// 			title: 'Menu items saved',
		// 		});
		// 	}).fail((result)=>{
		// 		console.log(result);
		// 		return Swal.fire({
		// 			icon: 'error',
		// 			title: 'Oops...',
		// 			text: result,
		// 		});
		// 	});
		// }
	}
};
Vue.createApp(MenuItemsUpdate).mount('#menu-items-update')`

	if cms.funcLayout("") != "" {
		out := hb.NewWrap().Children([]hb.TagInterface{
			hb.NewStyleURL("https://cdnjs.cloudflare.com/ajax/libs/jqtree/1.4.12/jqtree.css"),
			hb.NewStyle(`
	ul.jqtree-tree li>.jqtree-element {
		display: block;
		padding: 10px;
	}
	ul.jqtree-tree li.jqtree_common:hover{
		background:cornsilk;
	}
		`),
			container,
			hb.NewScriptURL(cdn.Jquery_3_6_4()),
			hb.NewScriptURL(cdn.VueJs_3()),
			hb.NewScriptURL(cdn.Sweetalert2_10()),
			hb.NewScriptURL("https://cdnjs.cloudflare.com/ajax/libs/jqtree/1.4.12/tree.jquery.js"),
			hb.NewScriptURL(cdn.Notify_0_4_2()),
			hb.NewScript(inlineScript),
		}).ToHTML()
		responses.HTMLResponse(w, r, cms.funcLayout(out))
		return
	}

	webpage := WebpageComplete("Edit Menu Items", container.ToHTML())
	webpage.AddStyleURLs([]string{
		"https://cdnjs.cloudflare.com/ajax/libs/jqtree/1.4.12/jqtree.css",
	})
	webpage.AddScriptURLs([]string{
		"https://cdnjs.cloudflare.com/ajax/libs/jqtree/1.4.12/tree.jquery.js",
		cdn.Notify_0_4_2(),
	})
	webpage.AddStyle(`
ul.jqtree-tree li>.jqtree-element {
	display: block;
	padding: 10px;
}
ul.jqtree-tree li.jqtree_common:hover{
    background:cornsilk;
}
	`)
	webpage.AddScript(inlineScript)
	responses.HTMLResponse(w, r, webpage.ToHTML())
}

func (cms Cms) pageMenusMenuItemsPagesDropdownList() (pagesDropdownList []map[string]string, errorMessage string) {
	pages, err := cms.EntityStore.EntityList(entitystore.EntityQueryOptions{
		EntityType: ENTITY_TYPE_PAGE,
		Offset:     0,
		Limit:      200,
		SortBy:     "id",
		SortOrder:  "asc",
	})

	if err != nil {
		return pagesDropdownList, "Page list failed to be retrieved " + err.Error()
	}

	pagesDropdownList = make([]map[string]string, 0)

	mapPageIDTitle := map[string]string{}
	for _, page := range pages {
		title, err := page.GetString("title", "")

		if err != nil {
			return pagesDropdownList, "Page failed to be retrieved " + err.Error()
		}

		status, err := page.GetString("status", "")

		if err != nil {
			return pagesDropdownList, "Page failed to be retrieved " + err.Error()
		}

		mapPageIDTitle[page.ID()] = title + " (" + status + ")"
	}

	pageTitles := lo.Values(mapPageIDTitle)

	sort.Strings(pageTitles)

	pagesDropdownList = []map[string]string{}

	for _, title := range pageTitles {
		pageID, isFound := lo.FindKey(mapPageIDTitle, title)
		if !isFound {
			continue
		}
		pagesDropdownList = append(pagesDropdownList, map[string]string{
			"key":   pageID,
			"value": title,
		})
	}

	return pagesDropdownList, ""
}

// flattenTree flattens a JQTree data
func flattenTree(nodes []map[string]interface{}) []map[string]interface{} {
	flatTree := []map[string]interface{}{}
	for index, node := range nodes {
		children, hasChildren := node["children"]
		delete(node, "children")

		node["sequence"] = utils.ToString((index + 1))
		flatTree = append(flatTree, node)

		if !hasChildren {
			continue
		}

		childrenArray := children.([]interface{})
		childrenMapArray := []map[string]interface{}{}
		for _, child := range childrenArray {
			childMap := child.(map[string]interface{})
			childMap["parent_id"] = node["id"]
			childrenMapArray = append(childrenMapArray, childMap)
		}
		childNodesList := flattenTree(childrenMapArray)
		flatTree = append(flatTree, childNodesList...)
	}
	return flatTree
}

func (cms Cms) pageMenusMenuItemsUpdateAjax(w http.ResponseWriter, r *http.Request) {
	menuID := strings.TrimSpace(utils.Req(r, "menu_id", ""))
	data := strings.TrimSpace(utils.Req(r, "data", ""))

	if menuID == "" {
		api.Respond(w, r, api.Error("Menu ID is required"))
		return
	}

	if data == "" {
		api.Respond(w, r, api.Error("Data is required"))
		return
	}

	menu, _ := cms.EntityStore.EntityFindByID(menuID)

	if menu == nil {
		api.Respond(w, r, api.Error("Menu NOT FOUND with ID "+menuID))
		return
	}

	var nodes []map[string]interface{}
	err := json.Unmarshal([]byte(data), &nodes)
	if err != nil {
		log.Println(err)
		api.Respond(w, r, api.Error("Menu failed to be updated. IO error"))
		return
	}

	existingMenuItemIDs := []string{}
	flatNodeList := flattenTree(nodes)

	for _, node := range flatNodeList {
		id := node["id"].(string)
		name := node["name"].(string)
		pageID := lo.ValueOr(node, "page_id", "").(string)
		url := lo.ValueOr(node, "url", "").(string)
		target := lo.ValueOr(node, "target", "").(string)
		parentID := lo.ValueOr(node, "parent_id", "").(string)
		sequence := lo.ValueOr(node, "sequence", "").(string)

		menuitem, _ := cms.EntityStore.EntityFindByID(id)
		if menuitem == nil {
			menuitem, err = cms.EntityStore.EntityCreateWithType(ENTITY_TYPE_MENUITEM)
			if err != nil {
				api.Respond(w, r, api.Error("Menu item failed to be created "+err.Error()))
				return
			}
		}
		menuitem.SetString("name", name)
		menuitem.SetString("menu_id", menuID)
		menuitem.SetString("parent_id", parentID)
		menuitem.SetString("sequence", sequence)
		menuitem.SetString("page_id", pageID)
		menuitem.SetString("url", url)
		err := menuitem.SetString("target", target)

		if err != nil {
			api.Respond(w, r, api.Error("Menu items failed to be updated "+err.Error()))
			return
		}

		existingMenuItemIDs = append(existingMenuItemIDs, menuitem.ID())
	}

	errMessage := cms.cleanMenuFromNonExistingMenuItems(menuID, existingMenuItemIDs)

	if errMessage != "" {
		api.Respond(w, r, api.Error(errMessage))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Menu items saved successfully", map[string]interface{}{"menu_id": menu.ID()}))
}

func (cms Cms) cleanMenuFromNonExistingMenuItems(menuID string, existingMenuItemIDs []string) (errorMessage string) {
	allMenuItems, err := cms.EntityStore.EntityListByAttribute(ENTITY_TYPE_MENUITEM, "menu_id", menuID)

	if err != nil {
		return "Menu items failed to be fetched: " + err.Error()
	}

	// Delete old menu items
	for _, menuitem := range allMenuItems {
		exists, _ := utils.ArrayContains(existingMenuItemIDs, menuitem.ID())
		if !exists {
			cms.EntityStore.EntityDelete(menuitem.ID())
		}
	}

	return ""
}

func keyExists(decoded map[string]interface{}, key string) bool {
	val, ok := decoded[key]
	return ok && val != nil
}
