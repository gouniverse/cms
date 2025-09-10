package cms

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dracory/bs"
	"github.com/dracory/entitystore"
	"github.com/dracory/hb"
	"github.com/gouniverse/api"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/responses"
	"github.com/gouniverse/utils"
)

func (cms Cms) pageEntitiesEntityCreateAjax(w http.ResponseWriter, r *http.Request) {
	entityType := utils.Req(r, "type", "")

	if entityType == "" {
		api.Respond(w, r, api.Error("Entity Type is required"))
		return
	}

	name := strings.Trim(utils.Req(r, "name", ""), " ")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	entity, err := cms.EntityStore.EntityCreateWithType(entityType)

	if err != nil {
		api.Respond(w, r, api.Error("Entity failed to be created"))
		return
	}

	if entity == nil {
		api.Respond(w, r, api.Error("Entity failed to be created"))
		return
	}

	entity.SetString("name", name)

	api.Respond(w, r, api.SuccessWithData("Entity saved successfully", map[string]interface{}{"entity_id": entity.ID()}))
}

func (cms Cms) pageEntitiesEntityManager(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	// log.Println(endpoint)

	entityType := utils.Req(r, "type", "")
	if entityType == "" {
		api.Respond(w, r, api.Error("Entity Type is required"))
		return
	}

	header := cms.cmsHeader(endpoint)
	breadcrums := cms.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  endpoint,
			Name: "Home",
		},
		{
			URL:  (endpoint + "?path=" + PathWidgetsWidgetManager),
			Name: "Custom Entities",
		},
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "entity-manager")
	heading := hb.NewHeading1().HTML("Entity Manager - " + entityType)
	button := hb.NewButton().HTML("New entity").Attr("class", "btn btn-success float-end").Attr("v-on:click", "showEntityCreateModal")
	heading.AddChild(button)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))

	// modal := hb.NewDiv().Attr("id", "ModalEntityCreate").Attr("class", "modal fade")
	// modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	// modalContent := hb.NewDiv().Attr("class", "modal-content")
	// modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("New Enity - " + entityType))
	// modalBody := hb.NewDiv().Attr("class", "modal-body")
	// modalBody.AddChild(hb.NewDiv().Attr("class", "form-group").AddChild(hb.NewLabel().HTML("Name")).AddChild(hb.NewInput().Attr("class", "form-control").Attr("v-model", "entityCreateModel.name")))
	// modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	// modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-prsecondary").Attr("data-bs-dismiss", "modal"))
	// modalFooter.AddChild(hb.NewButton().HTML("Create & Continue").Attr("class", "btn btn-primary").Attr("v-on:click", "entityCreate"))
	// modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	// modalDialog.AddChild(modalContent)
	// modal.AddChild(modalDialog)
	container.AddChild(cms.pageEntitiesEntityCreateModal())
	container.AddChild(cms.pageEntitiesEntityTrashModal())

	entities, err := cms.EntityStore.EntityList(entitystore.EntityQueryOptions{
		EntityType: entityType,
		Offset:     0,
		Limit:      200,
		SortBy:     "id",
		SortOrder:  "asc",
	})

	if err != nil {
		api.Respond(w, r, api.Error("Entities failed to be retrieved"))
		return
	}

	table := hb.NewTable().Attr("id", "TableEntities").Attr("class", "table table-responsive table-striped mt-3")
	thead := hb.NewThead()
	tbody := hb.NewTbody()
	table.AddChild(thead).AddChild(tbody)

	tr := hb.NewTR()
	th1 := hb.NewTD().HTML("Name")
	th2 := hb.NewTD().HTML("Status")
	th3 := hb.NewTD().HTML("Actions").Attr("style", "width:120px;")
	thead.AddChild(tr.AddChild(th1).AddChild(th2).AddChild(th3))

	for _, entity := range entities {
		name, _ := entity.GetString("name", "n/a")
		status, _ := entity.GetString("status", "n/a")
		buttonEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Attr("class", "btn btn-primary btn-sm").Attr("v-on:click", "entityEdit('"+entity.ID()+"')").Attr("style", "margin-right:5px")
		buttonTrash := hb.NewButton().HTML("Trash").Attr("type", "button").Attr("class", "btn btn-danger btn-sm").Attr("v-on:click", "showEntityTrashModal('"+entity.ID()+"')")

		tr := hb.NewTR()
		td1 := hb.NewTD().HTML(name)
		td2 := hb.NewTD().HTML(status)
		td3 := hb.NewTD().AddChild(buttonEdit).AddChild(buttonTrash)

		tbody.AddChild(tr.AddChild(td1).AddChild(td2).AddChild(td3))
	}
	container.AddChild(table)

	h := container.ToHTML()

	inlineScript := `
var entityCreateUrl = "` + endpoint + `?path=entities/entity-create-ajax";
var entityTrashUrl = "` + endpoint + `?path=entities/entity-trash-ajax";
var entityUpdateUrl = "` + endpoint + `?path=entities/entity-update";
const EntityManager = {
	data() {
		return {
		  entityCreateModel:{
			  name:"",
			  type:"` + entityType + `",
		  },
		  entityTrashModel:{
			entityId:null,
		  }
		}
	},
	created(){
		this.initDataTable();
	},
	methods: {
		initDataTable(){
			$(() => {
				$('#TableEntities').DataTable({
					"order": [[ 0, "asc" ]] // 1st column
				});
			});
		},
        showEntityCreateModal(){
			var modalEntityCreate = new bootstrap.Modal(document.getElementById('ModalEntityCreate'));
			modalEntityCreate.show();
		},
		showEntityTrashModal(entityId){
			this.entityTrashModel.entityId = entityId;
			var modalEntityDelete = new bootstrap.Modal(document.getElementById('ModalEntityDelete'));
			modalEntitykDelete.show();
		},
		entityCreate(){
			var name = this.entityCreateModel.name;
			var type = this.entityCreateModel.type;
		    $.post(entityCreateUrl, {type:type, name: name}).done((result)=>{
				if (result.status==="success"){
					var modalEntityCreate = new bootstrap.Modal(document.getElementById('ModalEntityCreate'));
			        modalEntityCreate.hide();

					return location.href = entityUpdateUrl+ "&entity_id=" + result.data.entity_id;
				}
				alert("Failed. " + result.message)
			}).fail((result)=>{
				alert("Failed" + result)
			});
		},
		entityEdit(entityId){
			return location.href = entityUpdateUrl+ "&entity_id=" + entityId;
		}
	}
};
setTimeout(()=>{
	Vue.createApp(EntityManager).mount('#entity-manager')
}, 1000);
	`

	if cms.funcLayout("") != "" {
		out := hb.NewWrap().Children([]hb.TagInterface{
			hb.NewStyleURL(cdn.JqueryDataTablesCss_1_13_4()),
			hb.NewHTML(h),
			hb.NewScriptURL(cdn.Jquery_3_6_4()),
			hb.NewScriptURL(cdn.VueJs_3()),
			hb.NewScriptURL(cdn.Sweetalert2_11()),
			hb.NewScriptURL(cdn.JqueryDataTablesJs_1_13_4()),
			hb.NewScript(inlineScript),
		}).ToHTML()
		responses.HTMLResponse(w, r, cms.funcLayout(out))
		return
	}

	webpage := WebpageComplete("Custom Entity Manager", h)
	webpage.AddStyleURL(cdn.JqueryDataTablesCss_1_13_4())
	webpage.AddScriptURL(cdn.JqueryDataTablesJs_1_13_4())
	webpage.AddScript(inlineScript)
	responses.HTMLResponse(w, r, webpage.ToHTML())
}

func (cms Cms) pageEntitiesEntityUpdate(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	// log.Println(endpoint)

	entityID := utils.Req(r, "entity_id", "")
	if entityID == "" {
		api.Respond(w, r, api.Error("Entity ID is required"))
		return
	}

	entity, _ := cms.EntityStore.EntityFindByID(entityID)

	if entity == nil {
		api.Respond(w, r, api.Error("Entity NOT FOUND with ID "+entityID))
		return
	}

	entityAttributeList := cms.customEntityAttributeList(entity.Type())

	header := cms.cmsHeader(r.Context().Value(keyEndpoint).(string))
	breadcrums := cms.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  endpoint,
			Name: "Home",
		},
		{
			URL:  (endpoint + "?path=" + PathEntitiesEntityManager + "&type=" + entity.Type()),
			Name: "Custom Entities",
		},
		{
			URL:  (endpoint + "?path=" + PathEntitiesEntityUpdate + "&entity_id=" + entityID),
			Name: "Edit Entity",
		},
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "entity-update")
	heading := hb.NewHeading1().HTML("Edit Custom Entity (type: " + entity.Type() + ")")
	button := hb.NewButton().HTML("Save").Attr("class", "btn btn-success float-end").Attr("v-on:click", "entitySave")
	heading.AddChild(button)

	formGroupStatus := hb.NewDiv().Attr("class", "form-group")
	formGroupStatusLabel := hb.NewLabel().HTML("Status").Attr("class", "form-label")
	formGroupStatusSelect := hb.NewSelect().Attr("class", "form-select").Attr("v-model", "entityModel.status")
	formGroupOptionsActive := hb.NewOption().Attr("value", "active").HTML("Active")
	formGroupOptionsInactive := hb.NewOption().Attr("value", "inactive").HTML("Inactive")
	formGroupOptionsTrash := hb.NewOption().Attr("value", "trash").HTML("Trash")
	formGroupStatusHelp := hb.NewParagraph().Attr("class", "text-info").HTML("What is the current status of the entity")
	formGroupStatus.AddChild(formGroupStatusLabel)
	formGroupStatus.AddChild(formGroupStatusSelect.AddChild(formGroupOptionsActive).AddChild(formGroupOptionsInactive).AddChild(formGroupOptionsTrash))
	formGroupStatus.AddChild(formGroupStatusHelp)

	formGroupName := hb.NewDiv().Attr("class", "form-group")
	formGroupNameLabel := hb.NewLabel().HTML("Name").Attr("class", "form-label")
	formGroupNameInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "entityModel.name")
	formGroupNameHelp := hb.NewParagraph().Attr("class", "text-info").HTML("What is the name of the entity")
	formGroupName.AddChild(formGroupNameLabel)
	formGroupName.AddChild(formGroupNameInput)
	formGroupName.AddChild(formGroupNameHelp)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))
	container.AddChild(formGroupStatus).AddChild(formGroupName)

	customAttrValues := make(map[string]string)
	for _, attr := range entityAttributeList {
		attrName := attr.Name
		attrFormControlLabel := attr.FormControlLabel
		if attrFormControlLabel == "" {
			attrFormControlLabel = attrName
		}
		formGroupAttr := hb.NewDiv().Attr("class", "form-group mt-3")
		formGroupAttrLabel := hb.NewLabel().HTML(attrFormControlLabel).Attr("class", "form-label")
		formGroupAttrInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "entityModel."+attrName)
		if attr.FormControlType == "textarea" {
			formGroupAttrInput = hb.NewTextArea().Attr("class", "form-control").Attr("v-model", "entityModel."+attrName)
		}
		if attr.BelongsToType != "" {
			entities, _ := cms.EntityStore.EntityList(entitystore.EntityQueryOptions{
				EntityType: attr.BelongsToType,
				Offset:     0,
				Limit:      300,
				SortBy:     "name",
				SortOrder:  "ASC",
			})
			formGroupAttrInput = hb.NewSelect().Attr("class", "form-select").Attr("v-model", "entityModel."+attrName)
			for _, ent := range entities {
				entName, _ := ent.GetString("name", "")
				formGroupAttrOption := hb.NewOption().Attr("value", ent.ID()).HTML(entName + " (" + ent.ID() + ")")
				formGroupAttrInput.AddChild(formGroupAttrOption)
			}
		}
		formGroupAttr.AddChild(formGroupAttrLabel)
		formGroupAttr.AddChild(formGroupAttrInput)

		// Add help
		if attr.FormControlHelp != "" {
			formGroupAttrHelp := hb.NewParagraph().Attr("class", "text-info").HTML(attr.FormControlHelp)
			formGroupAttr.AddChild(formGroupAttrHelp)
		}

		container.AddChild(formGroupAttr)

		customAttrValues[attrName], _ = entity.GetString(attrName, "")
	}

	h := container.ToHTML()

	name, _ := entity.GetString("name", "")
	status, _ := entity.GetString("status", "")
	jsonCustomValues, _ := json.Marshal(customAttrValues)

	inlineScript := `
var entityUpdateUrl = "` + endpoint + `?path=entities/entity-update-ajax";
var entityId = "` + entityID + `";
var name = "` + name + `";
var status = "` + status + `";
var customValues = ` + string(jsonCustomValues) + `;
const EntityUpdate = {
	data() {
		return {
			entityModel:{
				entityId: entityId,
				name: name,
				status: status,
				...customValues
		    }
		}
	},
	methods: {
		entitySave(){
			var name = this.entityModel.name;
			var entityId = this.entityModel.entityId;
			var status = this.entityModel.status;
			var data = JSON.parse(JSON.stringify(this.entityModel));
			data["entity_id"] = data["entityId"];
			delete data["entityId"];
			
			$.post(entityUpdateUrl, data).done((response)=>{
				if (response.status !== "success") {
					return Swal.fire({
						icon: 'error',
						title: 'Oops...',
						text: response.message,
					});
				}

				return Swal.fire({
					icon: 'success',
					title: 'Entity saved',
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

setTimeout(() => {
	Vue.createApp(EntityUpdate).mount('#entity-update')	
}, 1000);
	`

	if cms.funcLayout("") != "" {
		out := hb.NewWrap().Children([]hb.TagInterface{
			hb.NewHTML(h),
			hb.NewScriptURL(cdn.Jquery_3_6_4()),
			hb.NewScriptURL(cdn.VueJs_3()),
			hb.NewScriptURL(cdn.Sweetalert2_11()),
			hb.NewScript(inlineScript),
		}).ToHTML()
		responses.HTMLResponse(w, r, cms.funcLayout(out))
		return
	}

	webpage := WebpageComplete("Edit Custom Entity", h)
	webpage.AddScript(inlineScript)
	responses.HTMLResponse(w, r, webpage.ToHTML())
}

func (cms Cms) pageEntitiesEntityUpdateAjax(w http.ResponseWriter, r *http.Request) {
	entityID := strings.Trim(utils.Req(r, "entity_id", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	handle := strings.Trim(utils.Req(r, "handle", ""), " ")

	if entityID == "" {
		api.Respond(w, r, api.Error("Entity ID is required"))
		return
	}

	entity, _ := cms.EntityStore.EntityFindByID(entityID)

	if entity == nil {
		api.Respond(w, r, api.Error("Entity NOT FOUND with ID "+entityID))
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

	entityAttributeList := cms.customEntityAttributeList(entity.Type())
	for _, attr := range entityAttributeList {
		attrValue := strings.Trim(utils.Req(r, attr.Name, ""), " ")
		// attrLabel := attr.Label
		entity.SetString(attr.Name, attrValue)
	}

	entity.SetString("name", name)
	entity.SetString("handle", handle)
	err := entity.SetString("status", status)

	if err != nil {
		api.Respond(w, r, api.Error("Entity failed to be updated"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Entity saved successfully", map[string]interface{}{"entity_id": entity.ID()}))
}

func (cms Cms) customEntityAttributeList(entityType string) []CustomAttributeStructure {
	for _, entity := range cms.customEntityList {
		if entity.Type == entityType {
			if entity.AttributeList == nil {
				return []CustomAttributeStructure{}
			}
			return entity.AttributeList
		}
	}
	return []CustomAttributeStructure{}
}

func (cms Cms) pageEntitiesEntityTrashModal() *hb.Tag {
	modal := hb.NewDiv().Attr("id", "ModalEntityTrash").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("Trash Entity"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewParagraph().HTML("Are you sure you want to move this entity to trash bin?"))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-secondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Move to trash bin").Attr("class", "btn btn-danger").Attr("v-on:click", "entityTrash"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	return modal
}

func (cms Cms) pageEntitiesEntityCreateModal() *hb.Tag {
	modal := hb.NewDiv().Attr("id", "ModalEntityCreate").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("New Entity"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewDiv().Attr("class", "form-group").AddChild(hb.NewLabel().HTML("Name")).AddChild(hb.NewInput().Attr("class", "form-control").Attr("v-model", "entityCreateModel.name")))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-prsecondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Create & Continue").Attr("class", "btn btn-primary").Attr("v-on:click", "entityCreate"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	return modal
}
