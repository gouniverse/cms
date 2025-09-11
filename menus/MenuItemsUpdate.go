package cms

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/gouniverse/responses"
)

func (m UiManager) MenuItemsUpdate(w http.ResponseWriter, r *http.Request) {
	menuID := req.GetStringTrimmed(r, "menu_id")
	if menuID == "" {
		api.Respond(w, r, api.Error("Menu ID is required"))
		return
	}

	menu, _ := m.entityStore.EntityFindByID(menuID)

	if menu == nil {
		api.Respond(w, r, api.Error("Menu NOT FOUND with ID "+menuID))
		return
	}

	menuName, _ := menu.GetString("name", "")

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

	backURL := m.endpoint + "?path=" + m.pathMenusMenuManager
	buttonCancel := hb.NewHyperlink().
		Class("btn btn-info").
		Href(backURL).
		Child(hb.I().Class("bi bi-chevron-left")).
		HTML(" Cancel")
	buttonAddNode := hb.NewButton().
		Attr("v-on:click", "menuItemAdd").
		Class("btn btn-success float-end").
		ID("ButtonNewMenuItem").
		Attr("disabled", "disabled").
		Child(hb.I().Class("bi bi-plus-circle")).
		HTML(" Add item")
	buttonSaveNode := hb.NewButton().
		Attr("v-on:click", "menuItemsSave").
		Class("btn btn-success").
		ID("ButtonSaveMenuItems").
		Attr("disabled", "disabled").
		Child(hb.I().Class("bi bi-box-arrow-in-down")).
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

	pagesDropdownList, errorMessage := m.pageMenusMenuItemsPagesDropdownList()

	if errorMessage != "" {
		api.Respond(w, r, api.Error(errorMessage))
		return
	}

	menuItemsUpdateURL := m.endpoint + "?path=" + m.pathMenusMenuItemsUpdateAjax
	menuItemsFetchURL := m.endpoint + "?path=" + m.pathMenusMenuItemsFetchAjax
	pagesDropdownJSON, _ := json.Marshal(pagesDropdownList)
	inlineScript := `
var menuItemsSaveUrl = "` + menuItemsUpdateURL + `";
var menuItemsFetchUrl = "` + menuItemsFetchURL + `";
var menuId = "` + menuID + `";
var pagesDropdownList = ` + string(pagesDropdownJSON) + `;
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

	if m.funcLayout("") != "" {
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
		responses.HTMLResponse(w, r, m.funcLayout(out))
		return
	}

	webpage := m.webpageComplete("Edit Menu Items", container.ToHTML())
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
