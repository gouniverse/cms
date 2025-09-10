package cms

import (
	"net/http"

	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/gouniverse/api"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/responses"
)

func (m UiManager) BlockManager(w http.ResponseWriter, r *http.Request) {
	header := m.cmsHeader(m.endpoint)
	breadcrumbs := m.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  m.endpoint,
			Name: "Home",
		},
		{
			URL:  (m.endpoint + "?path=" + m.pathBlocksBlockManager),
			Name: "Blocks",
		},
	})

	container := hb.NewDiv().Class("container").Attr("id", "block-manager")
	heading := hb.NewHeading1().HTML("Block Manager")
	button := hb.NewButton().HTML("New block").Attr("class", "btn btn-success float-end").Attr("v-on:click", "showBlockCreateModal")
	heading.AddChild(button)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrumbs))

	container.AddChild(m.blockCreateModal())
	container.AddChild(m.blockTrashModal())

	blocks, err := m.entityStore.EntityList(entitystore.EntityQueryOptions{
		EntityType: m.blockEntityType,
		Offset:     0,
		Limit:      200,
		SortBy:     "id",
		SortOrder:  "asc",
	})

	if err != nil {
		api.Respond(w, r, api.Error("Blocks failed to be listed"))
		return
	}

	table := hb.NewTable().ID("TableBlocks").Class("table table-responsive table-striped mt-3")
	thead := hb.NewThead()
	tbody := hb.NewTbody()
	table.AddChild(thead).AddChild(tbody)

	tr := hb.NewTR()
	th1 := hb.NewTD().HTML("Name")
	th2 := hb.NewTD().HTML("Status")
	th3 := hb.NewTD().HTML("Actions").Attr("style", "width:120px;")
	thead.AddChild(tr.AddChild(th1).AddChild(th2).AddChild(th3))

	for _, block := range blocks {
		name, err := block.GetString("name", "n/a")
		if err != nil {
			api.Respond(w, r, api.Error("Name failed to be retrieved: "+err.Error()))
			return
		}
		status, err := block.GetString("status", "n/a")
		if err != nil {
			api.Respond(w, r, api.Error("Status failed to be retrieved: "+err.Error()))
			return
		}
		//buttonDelete := hb.NewButton().HTML("Delete").Attr("class", "btn btn-danger float-end").Attr("v-on:click", "showBlockDeleteModal('"+block.ID+"')")
		buttonEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Attr("class", "btn btn-primary btn-sm").Attr("v-on:click", "blockEdit('"+block.ID()+"')").Attr("style", "margin-right:5px")
		buttonTrash := hb.NewButton().HTML("Trash").Attr("class", "btn btn-danger btn-sm").Attr("v-on:click", "showBlockTrashModal('"+block.ID()+"')")

		tr := hb.NewTR()
		td1 := hb.NewTD().HTML(name)
		td2 := hb.NewTD().HTML(status)
		td3 := hb.NewTD().SetAttribute("style", "white-space:nowrap;").AddChild(buttonEdit).AddChild(buttonTrash)

		tbody.AddChild(tr.AddChild(td1).AddChild(td2).AddChild(td3))
	}
	container.AddChild(table)

	h := container.ToHTML()

	inlineScript := `
var blockCreateUrl = "` + m.endpoint + `?path=blocks/block-create-ajax"
var blockDeleteUrl = "` + m.endpoint + `?path=blocks/block-delete-ajax"
var blockTrashUrl = "` + m.endpoint + `?path=blocks/block-trash-ajax";
var blockUpdateUrl = "` + m.endpoint + `?path=blocks/block-update"
const BlockManager = {
	data() {
		return {
		  blockCreateModel:{
			  name:"",
		  },
		  blockDeleteModel:{
			blockId:null,
		  },
		  blockTrashModel:{
			blockId:null,
		  }
		}
	},
	created(){
		this.initDataTable();
	},
	methods: {
		initDataTable(){
			$(() => {
				$('#TableBlocks').DataTable({
					"order": [[ 0, "asc" ]] // 1st column
				});
			});
		},
        showBlockCreateModal(){
			var modalBlockCreate = new bootstrap.Modal(document.getElementById('ModalBlockCreate'));
			modalBlockCreate.show();
		},
		// showBlockDeleteModal(blockId){
		// 	this.blockDeleteModel.blockId = blockId;
		// 	var modalBlockDelete = new bootstrap.Modal(document.getElementById('ModalBlockDelete'));
		// 	modalBlockDelete.show();
		// },
		showBlockTrashModal(blockId){
			this.blockTrashModel.blockId = blockId;
			var modalBlockTrash = new bootstrap.Modal(document.getElementById('ModalBlockTrash'));
			modalBlockTrash.show();
		},
		blockCreate(){
			var name = this.blockCreateModel.name;
		    $.post(blockCreateUrl, {name: name}).done((result)=>{
				if (result.status==="success"){
					var modalBlockCreate = new bootstrap.Modal(document.getElementById('ModalBlockCreate'));
			        modalBlockCreate.hide();

					return location.href = blockUpdateUrl+ "&block_id=" + result.data.block_id;
				}
				alert("Failed. " + result.message)
			}).fail((result)=>{
				alert("Failed" + result)
			});
		},
		// blockDelete(){
		// 	var blockId = this.blockDeleteModel.blockId;
		//     $.post(blockDeleteUrl, {block_id: blockId}).done((result)=>{
		// 		if (result.status==="success"){
		// 			var modalBlockDelete = new bootstrap.Modal(document.getElementById('ModalBlockDelete'));
		// 	        modalBlockDelete.hide();

		// 			return location.href = location.href;
		// 		}
		// 		alert("Failed. " + result.message);
		// 	}).fail((result)=>{
		// 		alert("Failed" + result);
		// 	});
		// },
		blockTrash(){
			var blockId = this.blockTrashModel.blockId;
			$.post(blockTrashUrl, {block_id: blockId}).done((result)=>{
				if (result.status==="success"){
					var modalBlockDelete = new bootstrap.Modal(document.getElementById('ModalBlockTrash'));
					modalBlockDelete.hide();

					return location.href = location.href;
				}
				alert("Failed. " + result.message);
			}).fail((result)=>{
				alert("Failed" + result);
			});
		},
		blockEdit(blockId){
			return location.href = blockUpdateUrl+ "&block_id=" + blockId;
		}
	}
};
Vue.createApp(BlockManager).mount('#block-manager')
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

	webpage := m.webpageComplete("Block Manager", h)
	webpage.AddStyleURL(cdn.JqueryDataTablesCss_1_13_4())
	webpage.AddScriptURL(cdn.JqueryDataTablesJs_1_13_4())
	webpage.AddScript(inlineScript)
	responses.HTMLResponse(w, r, webpage.ToHTML())
}
