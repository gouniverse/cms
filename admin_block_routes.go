package cms

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/icons"
	"github.com/gouniverse/utils"
)

func pageBlocksBlockCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(utils.Req(r, "name", ""), " ")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	block, err := GetEntityStore().EntityCreate("block")

	if err != nil {
		api.Respond(w, r, api.Error("Block failed to be created"))
		return
	}

	if block == nil {
		api.Respond(w, r, api.Error("Block failed to be created"))
		return
	}

	block.SetString("name", name)

	api.Respond(w, r, api.SuccessWithData("Block saved successfully", map[string]interface{}{"block_id": block.ID}))
	return
}

func pageBlocksBlockManager(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	log.Println(endpoint)

	header := cmsHeader(endpoint)
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathBlocksBlockManager): "Blocks",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "block-manager")
	heading := hb.NewHeading1().HTML("Block Manager")
	button := hb.NewButton().HTML("New block").Attr("class", "btn btn-success float-end").Attr("v-on:click", "showBlockCreateModal")
	heading.AddChild(button)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))

	container.AddChild(pageBlocksBlockCreateModal())
	container.AddChild(pageBlocksBlockTrashModal())

	blocks, err := GetEntityStore().EntityList("block", 0, 200, "", "id", "asc")

	if err != nil {
		api.Respond(w, r, api.Error("Blocks failed to be listed"))
		return
	}

	table := hb.NewTable().Attr("id", "TableBlocks").Attr("class", "table table-responsive table-striped mt-3")
	thead := hb.NewThead()
	tbody := hb.NewTbody()
	table.AddChild(thead).AddChild(tbody)

	tr := hb.NewTR()
	th1 := hb.NewTD().HTML("Name")
	th2 := hb.NewTD().HTML("Status")
	th3 := hb.NewTD().HTML("Actions").Attr("style", "width:120px;")
	thead.AddChild(tr.AddChild(th1).AddChild(th2).AddChild(th3))

	for _, block := range blocks {
		name := block.GetString("name", "n/a")
		status := block.GetString("status", "n/a")
		//buttonDelete := hb.NewButton().HTML("Delete").Attr("class", "btn btn-danger float-end").Attr("v-on:click", "showBlockDeleteModal('"+block.ID+"')")
		buttonEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Attr("class", "btn btn-primary btn-sm").Attr("v-on:click", "blockEdit('"+block.ID+"')").Attr("style", "margin-right:5px")
		buttonTrash := hb.NewButton().HTML("Trash").Attr("class", "btn btn-danger btn-sm").Attr("v-on:click", "showBlockTrashModal('"+block.ID+"')")

		tr := hb.NewTR()
		td1 := hb.NewTD().HTML(name)
		td2 := hb.NewTD().HTML(status)
		td3 := hb.NewTD().SetAttribute("style", "white-space:nowrap;").AddChild(buttonEdit).AddChild(buttonTrash)

		tbody.AddChild(tr.AddChild(td1).AddChild(td2).AddChild(td3))
	}
	container.AddChild(table)

	h := container.ToHTML()

	inlineScript := `
var blockCreateUrl = "` + endpoint + `?path=blocks/block-create-ajax"
var blockDeleteUrl = "` + endpoint + `?path=blocks/block-delete-ajax"
var blockTrashUrl = "` + endpoint + `?path=blocks/block-trash-ajax";
var blockUpdateUrl = "` + endpoint + `?path=blocks/block-update"
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

	webpage := Webpage("Block Manager", h)
	webpage.AddStyleURL("https://cdnjs.cloudflare.com/ajax/libs/datatables/1.10.21/css/jquery.dataTables.css")
	webpage.AddScriptURL("https://cdnjs.cloudflare.com/ajax/libs/datatables/1.10.21/js/jquery.dataTables.js")
	webpage.AddScript(inlineScript)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

func pageBlocksBlockUpdate(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	log.Println(endpoint)

	blockID := utils.Req(r, "block_id", "")
	if blockID == "" {
		api.Respond(w, r, api.Error("Block ID is required"))
		return
	}

	block := GetEntityStore().EntityFindByID(blockID)

	if block == nil {
		api.Respond(w, r, api.Error("Block NOT FOUND with ID "+blockID))
		return
	}

	header := cmsHeader(r.Context().Value(keyEndpoint).(string))
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathBlocksBlockManager):                         "Blocks",
		(endpoint + "?path=" + PathBlocksBlockUpdate + "&block_id=" + blockID): "Edit block",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "block-update")
	heading := hb.NewHeading1().HTML("Edit Block")
	button := hb.NewButton().AddChild(hb.NewHTML(icons.BootstrapCheckCircle+" ")).HTML("Save").Attr("class", "btn btn-success float-end").Attr("v-on:click", "blockSave")
	heading.AddChild(button)

	formGroupStatus := hb.NewDiv().Attr("class", "form-group")
	formGroupStatusLabel := hb.NewLabel().HTML("Status").Attr("class", "form-label")
	formGroupStatusSelect := hb.NewSelect().Attr("class", "form-select").Attr("v-model", "blockModel.status")
	formGroupOptionsActive := hb.NewOption().Attr("value", "active").HTML("Active")
	formGroupOptionsInactive := hb.NewOption().Attr("value", "inactive").HTML("Inactive")
	formGroupOptionsTrash := hb.NewOption().Attr("value", "trash").HTML("Trash")
	formGroupStatus.AddChild(formGroupStatusLabel)
	formGroupStatus.AddChild(formGroupStatusSelect.AddChild(formGroupOptionsActive).AddChild(formGroupOptionsInactive).AddChild(formGroupOptionsTrash))

	formGroupName := hb.NewDiv().Attr("class", "form-group")
	formGroupNameLabel := hb.NewLabel().HTML("Name").Attr("class", "form-label")
	formGroupNameInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "blockModel.name")
	formGroupName.AddChild(formGroupNameLabel)
	formGroupName.AddChild(formGroupNameInput)

	formGroupContent := hb.NewDiv().Attr("class", "form-group")
	formGroupContentLabel := hb.NewLabel().HTML("Content").Attr("class", "form-label")
	formGroupContentInput := hb.NewTextArea().Attr("class", "form-control CodeMirror").Attr("v-model", "blockModel.content")
	formGroupContent.AddChild(formGroupContentLabel)
	formGroupContent.AddChild(formGroupContentInput)

	paragraphUsage := hb.NewParagraph().Attr("class", "text-info mt-5").AddChild(hb.NewHTML("To use this block in your website use the following shortcode:"))
	code := hb.NewCode().AddChild(hb.NewPRE().HTML(`&lt;!-- START: Block: ` + block.GetString("name", "") + ` -->
[[BLOCK_` + block.ID + `]]
&lt;!-- END: Block: ` + block.GetString("name", "") + ` -->`))
	paragraphUsage.AddChild(code)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))
	container.AddChild(formGroupStatus).AddChild(formGroupName).AddChild(formGroupContent)
	container.AddChild(paragraphUsage)

	h := container.ToHTML()

	name := block.GetString("name", "")
	statusAttribute, err := GetEntityStore().AttributeFind(block.ID, "status")

	if err != nil {
		api.Respond(w, r, api.Error("IO Error. Attribute failed to be pulled"))
		return
	}

	status := ""
	if statusAttribute != nil {
		status = statusAttribute.GetString()
	}
	contentAttribute, err := GetEntityStore().AttributeFind(block.ID, "content")

	if err != nil {
		api.Respond(w, r, api.Error("IO Error. Attribute failed to be fetched"))
		return
	}

	content := ""
	if contentAttribute != nil {
		content = contentAttribute.GetString()
	}

	contentJSON, _ := json.Marshal(content)
	nameJSON, _ := json.Marshal(name)
	statusJSON, _ := json.Marshal(status)

	inlineScript := `
var blockUpdateUrl = "` + endpoint + `?path=blocks/block-update-ajax";
var blockId = "` + blockID + `";
var name = ` + string(nameJSON) + `;
var status = ` + string(statusJSON) + `;
var content = ` + string(contentJSON) + `;
const BlockUpdate = {
	data() {
		return {
			blockModel:{
				blockId: blockId,
				content: content,
				name: name,
				status: status,
		    }
		}
	},
	created(){
		var self = this;
		setTimeout(function () {
			if ($('.CodeMirror').length > 0) {
				var editor = CodeMirror.fromTextArea($('.CodeMirror').get(0), {
					lineNumbers: true,
					matchBrackets: true,
					mode: "application/x-httpd-php",
					indentUnit: 4,
					indentWithTabs: true,
					enterMode: "keep", tabMode: "shift"
				});
				$(document).on('mouseup', '.CodeMirror', function() {
					self.blockModel.content = editor.getValue();
				});
				$(document).on('change', '.CodeMirror', function() {
					self.blockModel.content = editor.getValue();
				});
				setInterval(()=>{
					self.blockModel.content = editor.getValue();
				}, 1000)
			}
		}, 500);
	},
	methods: {
		blockSave(){
			var content = this.blockModel.content;
			var name = this.blockModel.name;
			var blockId = this.blockModel.blockId;
			var status = this.blockModel.status;
			
			$.post(blockUpdateUrl, {
				block_id:blockId,
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
					title: 'Block saved',
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
Vue.createApp(BlockUpdate).mount('#block-update')
	`

	webtemplate := Webpage("Edit Block", h)
	webtemplate.AddStyleURLs([]string{
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/codemirror.min.css",
	})
	webtemplate.AddScriptURLs([]string{
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/codemirror.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/xml/xml.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/htmlmixed/htmlmixed.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/javascript/javascript.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/css/css.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/clike/clike.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/php/php.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/2.36.0/formatting.min.js",
		"//cdnjs.cloudflare.com/ajax/libs/codemirror/3.22.0/addon/edit/matchbrackets.min.js",
	})
	webtemplate.AddStyle(`	
.CodeMirror {
	border: 1px solid #eee;
	height: auto;
}
	`)
	webtemplate.AddScript(inlineScript)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webtemplate.ToHTML()))
}

func pageBlocksBlockUpdateAjax(w http.ResponseWriter, r *http.Request) {
	blockID := strings.Trim(utils.Req(r, "block_id", ""), " ")
	content := strings.Trim(utils.Req(r, "content", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	handle := strings.Trim(utils.Req(r, "handle", ""), " ")

	if blockID == "" {
		api.Respond(w, r, api.Error("Block ID is required"))
		return
	}

	block := GetEntityStore().EntityFindByID(blockID)

	if block == nil {
		api.Respond(w, r, api.Error("Block NOT FOUND with ID "+blockID))
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

	block.SetString("content", content)
	block.SetString("name", name)
	block.SetString("handle", handle)
	isOk, _ := block.SetString("status", status)

	if isOk == false {
		api.Respond(w, r, api.Error("Block failed to be updated"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Block saved successfully", map[string]interface{}{"block_id": block.ID}))
	return
}

func pageBlocksBlockDeleteAjax(w http.ResponseWriter, r *http.Request) {
	blockID := strings.Trim(utils.Req(r, "block_id", ""), " ")

	if blockID == "" {
		api.Respond(w, r, api.Error("Block ID is required"))
		return
	}

	block := GetEntityStore().EntityFindByID(blockID)

	if block == nil {
		api.Respond(w, r, api.Success("Block already deleted"))
		return
	}

	isOk := GetEntityStore().EntityDelete(blockID)

	if isOk == false {
		api.Respond(w, r, api.Error("Block failed to be deleted"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Block deleted successfully", map[string]interface{}{"block_id": block.ID}))
	return
}

func pageBlocksBlockTrashAjax(w http.ResponseWriter, r *http.Request) {
	blockID := strings.Trim(utils.Req(r, "block_id", ""), " ")

	if blockID == "" {
		api.Respond(w, r, api.Error("Block ID is required"))
		return
	}

	block := GetEntityStore().EntityFindByID(blockID)

	if block == nil {
		api.Respond(w, r, api.Success("Block already deleted"))
		return
	}

	isOk, err := GetEntityStore().EntityTrash(blockID)

	if err != nil {
		api.Respond(w, r, api.Error("Block failed to be trashed"))
		return
	}

	if isOk == false {
		api.Respond(w, r, api.Error("Block failed to be trashed"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Block trashed successfully", map[string]interface{}{"block_id": block.ID}))
	return
}

func pageBlocksBlockTrashModal() *hb.Tag {
	modal := hb.NewDiv().Attr("id", "ModalBlockTrash").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("Trash Block"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewParagraph().HTML("Are you sure you want to move this block to trash bin?"))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-secondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Move to trash bin").Attr("class", "btn btn-danger").Attr("v-on:click", "blockTrash"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	return modal
}

func pageBlocksBlockCreateModal() *hb.Tag {
	modal := hb.NewDiv().Attr("id", "ModalBlockCreate").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("New Block"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewDiv().Attr("class", "form-group").AddChild(hb.NewLabel().HTML("Name")).AddChild(hb.NewInput().Attr("class", "form-control").Attr("v-model", "blockCreateModel.name")))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-prsecondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Create & Continue").Attr("class", "btn btn-primary").Attr("v-on:click", "blockCreate"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	return modal
}
