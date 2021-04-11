package cms

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	hb "github.com/gouniverse/html"
	"github.com/gouniverse/utils"
)

func pageEntitiesEntityCreateAjax(w http.ResponseWriter, r *http.Request) {
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

	entity := GetEntityStore().EntityCreate(entityType)

	if entity == nil {
		api.Respond(w, r, api.Error("Entity failed to be created"))
		return
	}

	entity.SetString("name", name)

	api.Respond(w, r, api.SuccessWithData("Entity saved successfully", map[string]interface{}{"entity_id": entity.ID}))
	return
}

func pageEntitiesEntityManager(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	log.Println(endpoint)

	entityType := utils.Req(r, "type", "")
	if entityType == "" {
		api.Respond(w, r, api.Error("Entity Type is required"))
		return
	}

	header := cmsHeader(endpoint)
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathWidgetsWidgetManager): "Custom Entities",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "entity-manager")
	heading := hb.NewHeading1().HTML("Entity Manager - " + entityType)
	button := hb.NewButton().HTML("New entity").Attr("class", "btn btn-success float-end").Attr("v-on:click", "showEntityCreateModal")
	heading.AddChild(button)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))

	modal := hb.NewDiv().Attr("id", "ModalEntityCreate").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("New Enity - " + entityType))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewDiv().Attr("class", "form-group").AddChild(hb.NewLabel().HTML("Name")).AddChild(hb.NewInput().Attr("class", "form-control").Attr("v-model", "entityCreateModel.name")))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-prsecondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Create & Continue").Attr("class", "btn btn-primary").Attr("v-on:click", "entityCreate"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	container.AddChild(modal)

	entities := GetEntityStore().EntityList(entityType, 0, 200, "", "id", "asc")

	table := hb.NewTable().Attr("class", "table table-responsive table-striped mt-3")
	thead := hb.NewThead()
	tbody := hb.NewTbody()
	table.AddChild(thead).AddChild(tbody)

	tr := hb.NewTR()
	th1 := hb.NewTD().HTML("Name")
	th2 := hb.NewTD().HTML("Status")
	th3 := hb.NewTD().HTML("Actions").Attr("style", "width:1px;")
	thead.AddChild(tr.AddChild(th1).AddChild(th2).AddChild(th3))

	for _, entity := range entities {
		name := entity.GetString("name", "n/a")
		status := entity.GetString("status", "n/a")
		buttonEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Attr("class", "btn btn-primary").Attr("v-on:click", "entityEdit('"+entity.ID+"')")

		tr := hb.NewTR()
		td1 := hb.NewTD().HTML(name)
		td2 := hb.NewTD().HTML(status)
		td3 := hb.NewTD().AddChild(buttonEdit)

		tbody.AddChild(tr.AddChild(td1).AddChild(td2).AddChild(td3))
	}
	container.AddChild(table)

	h := container.ToHTML()

	inlineScript := `
var entityCreateUrl = "` + endpoint + `?path=entities/entity-create-ajax"
var entityUpdateUrl = "` + endpoint + `?path=entities/entity-update"
const EntityManager = {
	data() {
		return {
		  entityCreateModel:{
			  name:"",
			  type:"` + entityType + `",
		  }
		}
	},
	methods: {
        showEntityCreateModal(){
			var modalEntityCreate = new bootstrap.Modal(document.getElementById('ModalEntityCreate'));
			modalEntityCreate.show();
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
Vue.createApp(EntityManager).mount('#entity-manager')
	`

	webwidget := Webpage("Custom Entity Manager", h)
	webwidget.AddScript(inlineScript)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webwidget.ToHTML()))
}

func pageEntitiesEntityUpdate(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	log.Println(endpoint)

	entityID := utils.Req(r, "entity_id", "")
	if entityID == "" {
		api.Respond(w, r, api.Error("Entity ID is required"))
		return
	}

	entity := GetEntityStore().EntityFindByID(entityID)

	if entity == nil {
		api.Respond(w, r, api.Error("Entity NOT FOUND with ID "+entityID))
		return
	}

	entityAttributeList := customEntityAttributeList(entity.Type)

	header := cmsHeader(r.Context().Value(keyEndpoint).(string))
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathEntitiesEntityManager + "&type=" + entity.Type):  "Custom Entities",
		(endpoint + "?path=" + PathEntitiesEntityUpdate + "&entity_id=" + entityID): "Edit Entity",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "entity-update")
	heading := hb.NewHeading1().HTML("Edit Custom Entity (type: " + entity.Type + ")")
	button := hb.NewButton().HTML("Save").Attr("class", "btn btn-success float-end").Attr("v-on:click", "entitySave")
	heading.AddChild(button)

	formGroupStatus := hb.NewDiv().Attr("class", "form-group")
	formGroupStatusLabel := hb.NewLabel().HTML("Status").Attr("class", "form-label")
	formGroupStatusSelect := hb.NewSelect().Attr("class", "form-control").Attr("v-model", "entityModel.status")
	formGroupOptionsActive := hb.NewOption().Attr("value", "active").HTML("Active")
	formGroupOptionsInactive := hb.NewOption().Attr("value", "inactive").HTML("Inactive")
	formGroupOptionsTrash := hb.NewOption().Attr("value", "trash").HTML("Trash")
	formGroupStatus.AddChild(formGroupStatusLabel)
	formGroupStatus.AddChild(formGroupStatusSelect.AddChild(formGroupOptionsActive).AddChild(formGroupOptionsInactive).AddChild(formGroupOptionsTrash))

	formGroupName := hb.NewDiv().Attr("class", "form-group")
	formGroupNameLabel := hb.NewLabel().HTML("Name").Attr("class", "form-label")
	formGroupNameInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "entityModel.name")
	formGroupName.AddChild(formGroupNameLabel)
	formGroupName.AddChild(formGroupNameInput)

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
		formGroupAttr.AddChild(formGroupAttrLabel)
		formGroupAttr.AddChild(formGroupAttrInput)

		// Add help
		if attr.FormControlHelp != "" {
			formGroupAttrHelp := hb.NewParagraph().Attr("class", "text-info").HTML(attr.FormControlHelp)
			formGroupAttr.AddChild(formGroupAttrHelp)
		}

		container.AddChild(formGroupAttr)

		customAttrValues[attrName] = entity.GetString(attrName, "")
	}

	h := container.ToHTML()

	name := entity.GetString("name", "")
	status := entity.GetString("status", "")
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
Vue.createApp(EntityUpdate).mount('#entity-update')
	`

	webpage := Webpage("Edit Custom Entity", h)
	webpage.AddScript(inlineScript)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

func pageEntitiesEntityUpdateAjax(w http.ResponseWriter, r *http.Request) {
	entityID := strings.Trim(utils.Req(r, "entity_id", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	handle := strings.Trim(utils.Req(r, "handle", ""), " ")

	if entityID == "" {
		api.Respond(w, r, api.Error("Entity ID is required"))
		return
	}

	entity := GetEntityStore().EntityFindByID(entityID)

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

	entityAttributeList := customEntityAttributeList(entity.Type)
	for _, attr := range entityAttributeList {
		attrValue := strings.Trim(utils.Req(r, attr.Name, ""), " ")
		// attrLabel := attr.Label
		entity.SetString(attr.Name, attrValue)
	}

	entity.SetString("name", name)
	entity.SetString("handle", handle)
	isOk := entity.SetString("status", status)

	if isOk == false {
		api.Respond(w, r, api.Error("Entity failed to be updated"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Entity saved successfully", map[string]interface{}{"entity_id": entity.ID}))
	return
}

func customEntityAttributeList(entityType string) []CustomAttributeStructure {
	for _, entity := range configuration.CustomEntityList {
		if entity.Type == entityType {
			if entity.AttributeList == nil {
				return []CustomAttributeStructure{}
			}
			return entity.AttributeList
		}
	}
	return []CustomAttributeStructure{}
}
