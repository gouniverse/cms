package cms

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/icons"
	"github.com/gouniverse/utils"
)

func pageMenusMenuCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(utils.Req(r, "name", ""), " ")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	menu, err := GetEntityStore().EntityCreate("menu")

	if err != nil {
		api.Respond(w, r, api.Error("Menu failed to be created "+err.Error()))
		return
	}

	if menu == nil {
		api.Respond(w, r, api.Error("Menu failed to be created"))
		return
	}

	menu.SetString("name", name)

	api.Respond(w, r, api.SuccessWithData("Menu saved successfully", map[string]interface{}{"menu_id": menu.ID}))
	return
}

func pageMenusMenuManager(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	//log.Println(endpoint)

	header := cmsHeader(endpoint)
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathMenusMenuManager): "Menus",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "menu-manager")
	heading := hb.NewHeading1().HTML("Menu Manager")
	button := hb.NewButton().HTML("New menu").Attr("class", "btn btn-success float-end").Attr("v-on:click", "showMenuCreateModal")
	heading.AddChild(button)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))

	modal := hb.NewDiv().Attr("id", "ModalMenuCreate").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("New Menu"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewDiv().Attr("class", "form-group").AddChild(hb.NewLabel().HTML("Name")).AddChild(hb.NewInput().Attr("class", "form-control").Attr("v-model", "menuCreateModel.name")))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-prsecondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Create & Continue").Attr("class", "btn btn-primary").Attr("v-on:click", "menuCreate"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	container.AddChild(modal)

	menus, err := entityStore.EntityList("menu", 0, 200, "", "id", "asc")

	if err != nil {
		api.Respond(w, r, api.Error("Entity list failed to be retrieved "+err.Error()))
		return
	}

	table := hb.NewTable().Attr("class", "table table-responsive table-striped mt-3")
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
		buttonEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Attr("class", "btn btn-primary").Attr("v-on:click", "menuEdit('"+menu.ID+"')")
		buttonMenuItemsEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Attr("class", "btn btn-primary").Attr("v-on:click", "menuItemsEdit('"+menu.ID+"')")

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

	webmenu := Webpage("Menu Manager", h)
	webmenu.AddScript(inlineScript)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webmenu.ToHTML()))
}

func pageMenusMenuUpdate(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	//log.Println(endpoint)

	menuID := utils.Req(r, "menu_id", "")
	if menuID == "" {
		api.Respond(w, r, api.Error("Menu ID is required"))
		return
	}

	menu, _ := GetEntityStore().EntityFindByID(menuID)

	if menu == nil {
		api.Respond(w, r, api.Error("Menu NOT FOUND with ID "+menuID))
		return
	}

	header := cmsHeader(r.Context().Value(keyEndpoint).(string))
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathMenusMenuManager):                       "Menus",
		(endpoint + "?path=" + PathMenusMenuUpdate + "&menu_id=" + menuID): "Edit menu",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "menu-update")
	heading := hb.NewHeading1().HTML("Edit Menu")
	button := hb.NewButton().AddChild(hb.NewHTML(icons.BootstrapCheckCircle+" ")).HTML("Save").Attr("class", "btn btn-success float-end").Attr("v-on:click", "menuSave")
	heading.AddChild(button)

	formGroupStatus := hb.NewDiv().Attr("class", "form-group")
	formGroupStatusLabel := hb.NewLabel().HTML("Status").Attr("class", "form-label")
	formGroupStatusSelect := hb.NewSelect().Attr("class", "form-select").Attr("v-model", "menuModel.status")
	formGroupOptionsActive := hb.NewOption().Attr("value", "active").HTML("Active")
	formGroupOptionsInactive := hb.NewOption().Attr("value", "inactive").HTML("Inactive")
	formGroupOptionsTrash := hb.NewOption().Attr("value", "trash").HTML("Trash")
	formGroupStatus.AddChild(formGroupStatusLabel)
	formGroupStatus.AddChild(formGroupStatusSelect.AddChild(formGroupOptionsActive).AddChild(formGroupOptionsInactive).AddChild(formGroupOptionsTrash))

	formGroupName := hb.NewDiv().Attr("class", "form-group")
	formGroupNameLabel := hb.NewLabel().HTML("Name").Attr("class", "form-label")
	formGroupNameInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "menuModel.name")
	formGroupName.AddChild(formGroupNameLabel)
	formGroupName.AddChild(formGroupNameInput)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))
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

	webmenu := Webpage("Edit Menu", h)
	webmenu.AddScript(inlineScript)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webmenu.ToHTML()))
}

func pageMenusMenuUpdateAjax(w http.ResponseWriter, r *http.Request) {
	menuID := strings.Trim(utils.Req(r, "menu_id", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	handle := strings.Trim(utils.Req(r, "handle", ""), " ")

	if menuID == "" {
		api.Respond(w, r, api.Error("Menu ID is required"))
		return
	}

	menu, _ := GetEntityStore().EntityFindByID(menuID)

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
	isOk, err := menu.SetString("status", status)

	if err != nil {
		api.Respond(w, r, api.Error("Menu failed to be updated "+err.Error()))
		return
	}

	if isOk == false {
		api.Respond(w, r, api.Error("Menu failed to be updated"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Menu saved successfully", map[string]interface{}{"menu_id": menu.ID}))
	return
}

func getChildren(data []map[string]interface{}, parentID string) []map[string]interface{} {
	children := []map[string]interface{}{}
	sequences := []string{}
	for _, node := range data {
		nodeParentID := ""
		sequence := ""
		if keyExists(node, "parent_id") {
			nodeParentID = node["parent_id"].(string)
		}
		if keyExists(node, "sequence") {
			sequence = node["sequence"].(string)
		}
		if nodeParentID == parentID {
			sequences = append(sequences, sequence)
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
			for _, childTree := range childrenTrees {
				rootChildren = append(rootChildren, childTree)
			}
			children[childIndex]["children"] = rootChildren
		}
		root["children"] = children
		out = append(out, root)
	}

	return out
}

func buildTree(menuID string) []map[string]interface{} {
	menuitems, err := GetEntityStore().EntityListByAttribute("menuitem", "menu_id", menuID)

	if err != nil {
		log.Panicln("Menu items failed to be retrieved " + err.Error())
		return nil
	}

	nodeList := []map[string]interface{}{}
	for _, menuitem := range menuitems {
		itemID := menuitem.ID
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

func pageMenusMenuItemsFetchAjax(w http.ResponseWriter, r *http.Request) {
	menuID := strings.Trim(utils.Req(r, "menu_id", ""), " ")

	if menuID == "" {
		api.Respond(w, r, api.Error("Menu ID is required"))
		return
	}

	menu, _ := GetEntityStore().EntityFindByID(menuID)

	if menu == nil {
		api.Respond(w, r, api.Error("Menu NOT FOUND with ID "+menuID))
		return
	}

	tree := buildTree(menuID)

	api.Respond(w, r, api.SuccessWithData("Menu items found successfully", map[string]interface{}{
		"menu_id":   menu.ID,
		"menuitems": tree,
	}))
	return
}

func pageMenusMenuItemsUpdate(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	//log.Println(endpoint+"")

	menuID := utils.Req(r, "menu_id", "")
	if menuID == "" {
		api.Respond(w, r, api.Error("Menu ID is required"))
		return
	}

	menu, _ := GetEntityStore().EntityFindByID(menuID)

	if menu == nil {
		api.Respond(w, r, api.Error("Menu NOT FOUND with ID "+menuID))
		return
	}

	menuName, _ := menu.GetString("name", "")

	header := cmsHeader(r.Context().Value(keyEndpoint).(string))
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathMenusMenuManager):                       "Menus",
		(endpoint + "?path=" + PathMenusMenuUpdate + "&menu_id=" + menuID): "Menu",
		"#": "Edit menu items",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "menu-update")
	heading := hb.NewHeading1().HTML("Edit Menu Items for Menu '" + menuName + "'")
	//button := hb.NewButton().AddChild(hb.NewHTML(icons.BootstrapCheckCircle+" ")).HTML("Save").Attr("class", "btn btn-success float-end").Attr("v-on:click", "menuSave")
	//heading.AddChild(button)

	backURL := endpoint + "?path=" + PathMenusMenuManager
	actionsCard := hb.NewDiv().Attr("class", "card box-primary")
	actionsCardHeader := hb.NewDiv().Attr("class", "card-header with-border")
	actionsCardBody := hb.NewDiv().Attr("class", "card-body with-border")

	buttonCancel := hb.NewHyperlink().Attr("class", "btn btn-info").Attr("href", backURL).AddChild(hb.NewHTML(icons.BootstrapChevronLeft)).HTML(" Cancel")
	buttonAddNode := hb.NewButton().Attr("v-on:click", "menuItemAdd").Attr("class", "btn btn-success float-end").Attr("id", "ButtonNewMenuItem").Attr("disabled", "disabled").AddChild(hb.NewHTML(icons.BootstrapPlusCircle)).HTML(" Add item")
	buttonSaveNode := hb.NewButton().Attr("v-on:click", "menuItemsSave").Attr("class", "btn btn-success").Attr("id", "ButtonSaveMenuItems").Attr("disabled", "disabled").AddChild(hb.NewHTML(icons.BootstrapBoxArrowInDown)).HTML(" Save")

	divTree := hb.NewDiv().Attr("id", "tree1")
	actionsCardBody.AddChild(divTree)
	actionsCardHeader.AddChild(buttonCancel).AddChild(buttonAddNode).AddChild(buttonSaveNode)
	actionsCard.AddChild(actionsCardHeader).AddChild(actionsCardBody)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))
	container.AddChild(actionsCard)

	modal := hb.NewDiv().Attr("id", "ModalItemUpdate").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("Edit Menu Item"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewDiv().Attr("class", "form-group").AddChild(hb.NewLabel().HTML("Title")).AddChild(hb.NewInput().Attr("class", "form-control").Attr("v-model", "menuItemUpdateModel.title")))
	modalBody.AddChild(hb.NewDiv().Attr("class", "form-group").AddChild(hb.NewLabel().HTML("Page")).AddChild(hb.NewInput().Attr("class", "form-control").Attr("v-model", "menuItemUpdateModel.pageId")))
	modalBody.AddChild(hb.NewDiv().Attr("class", "form-group").AddChild(hb.NewLabel().HTML("URL")).AddChild(hb.NewInput().Attr("class", "form-control").Attr("v-model", "menuItemUpdateModel.url")))
	modalBody.AddChild(hb.NewDiv().Attr("class", "form-group").AddChild(hb.NewLabel().HTML("Target")).AddChild(hb.NewInput().Attr("class", "form-control").Attr("v-model", "menuItemUpdateModel.target")))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-prsecondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Update").Attr("class", "btn btn-primary").Attr("v-on:click", "menuItemUpdate").Attr("data-bs-dismiss", "modal"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	container.AddChild(modal)

	h := container.ToHTML()

	menuItemsUpdateURL := endpoint + "?path=" + PathMenusMenuItemsUpdateAjax
	menuItemsFetchURL := endpoint + "?path=" + PathMenusMenuItemsFetchAjax
	inlineScript := `
var menuItemsSaveUrl = "` + menuItemsUpdateURL + `";
var menuItemsFetchUrl = "` + menuItemsFetchURL + `";
var menuId = "` + menuID + `";
const MenuUpdate = {
	data() {
		return {
			menuId: menuId,
			menuitems:[],
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
Vue.createApp(MenuUpdate).mount('#menu-update')`
	// <script src="https://cdn.jsdelivr.net/npm/bootstrap@4.5.3/dist/js/bootstrap.bundle.min.js"></script>
	//         <script src="https://cdnjs.cloudflare.com/ajax/libs/jqtree/1.4.12/tree.jquery.js"></script>
	//         <script src="https://cdnjs.cloudflare.com/ajax/libs/notify/0.4.2/notify.min.js"></script>
	//         <script>var menuId = 1;</script>
	//         <script>var menuItemsSaveUrl = "/save";</script>
	//         <script>var menuItemsFetchUrl = "/menus.json";</script>

	webpage := Webpage("Edit Menu Items", h)
	webpage.AddStyleURLs([]string{
		"https://cdnjs.cloudflare.com/ajax/libs/jqtree/1.4.12/jqtree.css",
	})
	webpage.AddScriptURLs([]string{
		"https://cdnjs.cloudflare.com/ajax/libs/jqtree/1.4.12/tree.jquery.js",
		"https://cdnjs.cloudflare.com/ajax/libs/notify/0.4.2/notify.min.js",
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
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

// flattenTree flattens a JQTree data
func flattenTree(nodes []map[string]interface{}) []map[string]interface{} {
	flatTree := []map[string]interface{}{}
	for index, node := range nodes {
		children, hasChildren := node["children"]
		delete(node, "children")

		node["sequence"] = utils.ToString((index + 1))
		flatTree = append(flatTree, node)

		if hasChildren == false {
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
		for _, childNode := range childNodesList {
			//childNode["parent_id"] = node["id"]
			flatTree = append(flatTree, childNode)
		}
	}
	return flatTree
}

func pageMenusMenuItemsUpdateAjax(w http.ResponseWriter, r *http.Request) {
	menuID := strings.Trim(utils.Req(r, "menu_id", ""), " ")
	data := strings.Trim(utils.Req(r, "data", ""), " ")

	if menuID == "" {
		api.Respond(w, r, api.Error("Menu ID is required"))
		return
	}

	if data == "" {
		api.Respond(w, r, api.Error("Data is required"))
		return
	}

	menu, _ := GetEntityStore().EntityFindByID(menuID)

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

	newIDs := []string{}
	flatNodeList := flattenTree(nodes)

	for _, node := range flatNodeList {
		//log.Println(node)
		id := node["id"].(string)
		name := node["name"].(string)
		pageID := ""
		url := ""
		target := ""
		parentID := ""
		sequence := ""
		if keyExists(node, "page_id") {
			pageID = node["page_id"].(string)
		}
		if keyExists(node, "url") {
			url = node["url"].(string)
		}
		if keyExists(node, "target") {
			target = node["target"].(string)
		}
		if keyExists(node, "parent_id") {
			parentID = node["parent_id"].(string)
		}
		if keyExists(node, "sequence") {
			sequence = node["sequence"].(string)
		}

		menuitem, _ := GetEntityStore().EntityFindByID(id)
		if menuitem == nil {
			menuitem, err = GetEntityStore().EntityCreate("menuitem")
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
		isOk, err := menuitem.SetString("target", target)

		if err != nil {
			api.Respond(w, r, api.Error("Menu items failed to be updated "+err.Error()))
			return
		}

		if isOk == false {
			api.Respond(w, r, api.Error("Menu items failed to be updated"))
			return
		}

		newIDs = append(newIDs, menuitem.ID)
	}

	allMenuItems, err := GetEntityStore().EntityListByAttribute("menuitem", "menu_id", menuID)

	if err != nil {
		api.Respond(w, r, api.Error("Menu items failed to be fetched: "+err.Error()))
		return
	}

	for _, menuitem := range allMenuItems {
		//allIDs = append(allIDs, menuitem.ID)
		exists, _ := utils.ArrayContains(newIDs, menuitem.ID)
		if exists == false {
			GetEntityStore().EntityDelete(menuitem.ID)
		}
	}

	api.Respond(w, r, api.SuccessWithData("Menu saved successfully", map[string]interface{}{"menu_id": menu.ID}))
	return
}

func keyExists(decoded map[string]interface{}, key string) bool {
	val, ok := decoded[key]
	return ok && val != nil
}
