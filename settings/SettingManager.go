package cms

import (
	"net/http"

	"github.com/gouniverse/api"
	"github.com/gouniverse/bs"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/responses"
	"github.com/gouniverse/settingstore"
)

func (m UiManager) SettingManager(w http.ResponseWriter, r *http.Request) {
	header := m.cmsHeader(m.endpoint)
	breadcrums := m.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  m.endpoint,
			Name: "Home",
		},
		{
			URL:  (m.endpoint + "?path=" + m.pathSettingsSettingManager),
			Name: "Settings",
		},
	})

	container := bs.Container().ID("setting-manager")
	heading := hb.NewHeading1().HTML("Setting Manager")
	button := bs.Button().HTML("New setting").Class("btn-success float-end").Attr("v-on:click", "showSettingCreateModal")
	heading.AddChild(button)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))

	container.AddChild(m.settingDeleteModal())
	container.AddChild(m.settingCreateModal())

	settings, err := m.settingStore.SettingList(r.Context(), settingstore.SettingQuery())

	if err != nil {
		api.Respond(w, r, api.Error("Setting keys failed to be retrieved "+err.Error()))
		return
	}

	settingKeys := []string{}

	for _, setting := range settings {
		settingKeys = append(settingKeys, setting.GetKey())
	}

	table := hb.NewTable().ID("TableSettings").Class("table table-responsive table-striped mt-3")
	thead := hb.NewThead()
	tbody := hb.NewTbody()
	table.Child(thead).Child(tbody)

	tr := hb.NewTR()
	th1 := hb.NewTD().HTML("Name")
	th3 := hb.NewTD().HTML("Actions").Attr("style", "width:1px;")
	thead.AddChild(tr.AddChild(th1).AddChild(th3))

	for _, settingKey := range settingKeys {
		name := settingKey
		buttonEdit := bs.Button().HTML("Edit").Attr("type", "button").Class("btn-primary btn-sm").Attr("v-on:click", "settingEdit('"+settingKey+"')")
		buttonDelete := bs.Button().HTML("Delete").Attr("type", "button").Class("btn-danger btn-sm").Attr("v-on:click", "showSettingDeleteModal('"+settingKey+"')")

		tr := hb.NewTR()
		td1 := hb.NewTD().HTML(name)
		td3 := hb.NewTD().Attr("style", "width:150px;").AddChild(buttonEdit).AddChild(buttonDelete)

		tbody.AddChild(tr.AddChild(td1).AddChild(td3))
	}
	container.AddChild(table)

	h := container.ToHTML()

	inlineScript := `
var settingCreateUrl = "` + m.endpoint + `?path=settings/setting-create-ajax"
var settingDeleteUrl = "` + m.endpoint + `?path=settings/setting-delete-ajax";
var settingUpdateUrl = "` + m.endpoint + `?path=settings/setting-update"
const SettingManager = {
	data() {
		return {
		  settingCreateModel:{
			  key:""
		  },
		  settingDeleteModel:{
			  key:"",
		  },
		}
	},
	created(){
		this.initDataTable();
	},
	methods: {
		initDataTable(){
			$(() => {
				$('#TableSettings').DataTable({
					"order": [[ 0, "asc" ]] // 1st column
				});
			});
		},
        showSettingCreateModal(){
			var modalSettingCreate = new bootstrap.Modal(document.getElementById('ModalSettingCreate'));
			modalSettingCreate.show();
		},
		settingCreate(){
			var key = this.settingCreateModel.key;
		    $.post(settingCreateUrl, {key: key}).done((result)=>{
				if (result.status==="success"){
					var modalSettingCreate = new bootstrap.Modal(document.getElementById('ModalSettingCreate'));
			        modalSettingCreate.hide();

					return location.href = settingUpdateUrl+ "&setting_key=" + result.data.setting_key;
				}
				alert("Failed. " + result.message)
			}).fail((result)=>{
				alert("Failed" + result)
			});
		},
		settingEdit(settingKey){
			return location.href = settingUpdateUrl+ "&setting_key=" + settingKey;
		},
		showSettingDeleteModal(settingKey){
			this.settingDeleteModel.key = settingKey;
			var modalSettingDelete = new bootstrap.Modal(document.getElementById('ModalSettingDelete'));
			modalSettingDelete.show();
		},
		settingDelete(){
            let settingKey = this.settingDeleteModel.key;
			$.post(settingDeleteUrl, {setting_key: settingKey}).done((result)=>{
				if (result.status==="success"){
					var ModalSettingDelete = new bootstrap.Modal(document.getElementById('ModalSettingDelete'));
				    ModalSettingDelete.hide();
					location.href = location.href;
					return;
				}
				alert("Failed. " + result.message)
			}).fail((result)=>{
				alert("Failed" + result)
			});
		}
	}
};
Vue.createApp(SettingManager).mount('#setting-manager')
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

	webpage := m.webpageComplete("Setting Manager", h)
	webpage.AddStyleURL(cdn.JqueryDataTablesCss_1_13_4())
	webpage.AddScriptURL(cdn.JqueryDataTablesJs_1_13_4())
	webpage.AddScript(inlineScript)

	responses.HTMLResponse(w, r, m.funcLayout(webpage.ToHTML()))
}

func (m UiManager) settingCreateModal() *hb.Tag {
	modal := hb.NewDiv().Attr("id", "ModalSettingCreate").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("New Setting"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewDiv().Attr("class", "form-group").AddChild(hb.NewLabel().HTML("Key")).AddChild(hb.NewInput().Attr("class", "form-control").Attr("v-model", "settingCreateModel.key")))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-prsecondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Create & Continue").Attr("class", "btn btn-primary").Attr("v-on:click", "settingCreate"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	return modal
}

func (m UiManager) settingDeleteModal() *hb.Tag {
	modal := hb.NewDiv().Attr("id", "ModalSettingDelete").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("Delete Setting"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewParagraph().HTML("Are you sure you want to delete this setting?"))
	modalBody.AddChild(hb.NewParagraph().HTML("Note!. This action cannot be reveresed"))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Cancel").Attr("class", "btn btn-secondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Delete").Attr("class", "btn btn-danger").Attr("v-on:click", "settingDelete"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	return modal
}
