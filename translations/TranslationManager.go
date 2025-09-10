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

func (m UiManager) TranslationManager(w http.ResponseWriter, r *http.Request) {
	header := m.cmsHeader(m.endpoint)
	breadcrumbs := m.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  m.endpoint,
			Name: "Home",
		},
		{
			URL:  (m.endpoint + "?path=" + m.pathTranslationsTranslationManager),
			Name: "Translations",
		},
	})

	button := hb.NewButton().
		HTML("New translation").
		Class("btn btn-success float-end").
		Attr("v-on:click", "showTranslationCreateModal")

	heading := hb.NewHeading1().
		HTML("Translation Manager").
		Child(button)

	container := hb.NewDiv().
		ID("translation-manager").
		Class("container").
		Children([]hb.TagInterface{
			hb.NewHTML(header),
			heading,
			hb.NewHTML(breadcrumbs),
			m.translationCreateModal(),
			m.translationTrashModal(),
		})

	translations, err := m.entityStore.EntityList(entitystore.EntityQueryOptions{
		EntityType: m.translationEntityType,
		Offset:     0,
		Limit:      200,
		SortBy:     "id",
		SortOrder:  "asc",
	})

	if err != nil {
		api.Respond(w, r, api.Error("Translations failed to be listed"))
		return
	}

	table := hb.NewTable().
		ID("TableTranslations").
		Class("table table-responsive table-striped mt-3")
	thead := hb.NewThead()
	tbody := hb.NewTbody()
	table.AddChild(thead).AddChild(tbody)

	tr := hb.NewTR()
	th1 := hb.NewTD().HTML("Name")
	th2 := hb.NewTD().HTML("Status").Style("width:1px;")
	th3 := hb.NewTD().HTML("Handle").Style("width:1px;")
	th4 := hb.NewTD().HTML("Actions").Style("width:120px;")
	thead.Child(tr.Child(th1).Child(th2).Child(th3).Child(th4))

	for _, translation := range translations {
		name, err := translation.GetString("name", "n/a")
		if err != nil {
			api.Respond(w, r, api.Error("Name failed to be retrieved: "+err.Error()))
			return
		}
		status, err := translation.GetString("status", "n/a")
		if err != nil {
			api.Respond(w, r, api.Error("Status failed to be retrieved: "+err.Error()))
			return
		}
		//buttonDelete := hb.NewButton().HTML("Delete").Attr("class", "btn btn-danger float-end").Attr("v-on:click", "showTranslationDeleteModal('"+translation.ID+"')")
		buttonEdit := hb.NewButton().HTML("Edit").Type(hb.TYPE_BUTTON).Class("btn btn-primary btn-sm").Attr("v-on:click", "translationEdit('"+translation.ID()+"')").Style("margin-right:5px")
		buttonTrash := hb.NewButton().HTML("Trash").Class("btn btn-danger btn-sm").Attr("v-on:click", "showTranslationTrashModal('"+translation.ID()+"')")

		tr := hb.NewTR()
		td1 := hb.NewTD().HTML(name)
		td2 := hb.NewTD().HTML(status)
		td3 := hb.NewTD().HTML(translation.Handle())
		td4 := hb.NewTD().Style("white-space:nowrap;").Child(buttonEdit).Child(buttonTrash)

		tbody.Child(tr.Child(td1).Child(td2).Child(td3).Child(td4))
	}
	container.AddChild(table)

	h := container.ToHTML()

	inlineScript := `
var translationCreateUrl = "` + m.endpoint + `?path=translations/translation-create-ajax"
var translationDeleteUrl = "` + m.endpoint + `?path=translations/translation-delete-ajax"
var translationTrashUrl = "` + m.endpoint + `?path=translations/translation-trash-ajax";
var translationUpdateUrl = "` + m.endpoint + `?path=translations/translation-update"
const TranslationManager = {
	data() {
		return {
		  translationCreateModel:{
			  name:"",
		  },
		  translationDeleteModel:{
			translationId:null,
		  },
		  translationTrashModel:{
			translationId:null,
		  }
		}
	},
	created(){
		this.initDataTable();
	},
	methods: {
		initDataTable(){
			$(() => {
				$('#TableTranslations').DataTable({
					"order": [[ 0, "asc" ]] // 1st column
				});
			});
		},
        showTranslationCreateModal(){
			var modalTranslationCreate = new bootstrap.Modal(document.getElementById('ModalTranslationCreate'));
			modalTranslationCreate.show();
		},
		// showTranslationDeleteModal(translationId){
		// 	this.translationDeleteModel.translationId = translationId;
		// 	var modalTranslationDelete = new bootstrap.Modal(document.getElementById('ModalTranslationDelete'));
		// 	modalTranslationDelete.show();
		// },
		showTranslationTrashModal(translationId){
			this.translationTrashModel.translationId = translationId;
			var modalTranslationTrash = new bootstrap.Modal(document.getElementById('ModalTranslationTrash'));
			modalTranslationTrash.show();
		},
		translationCreate(){
			var name = this.translationCreateModel.name;
		    $.post(translationCreateUrl, {name: name}).done((result)=>{
				if (result.status==="success"){
					var modalTranslationCreate = new bootstrap.Modal(document.getElementById('ModalTranslationCreate'));
			        modalTranslationCreate.hide();

					return location.href = translationUpdateUrl+ "&translation_id=" + result.data.translation_id;
				}
				alert("Failed. " + result.message)
			}).fail((result)=>{
				alert("Failed" + result)
			});
		},
		// translationDelete(){
		// 	var translationId = this.translationDeleteModel.translationId;
		//     $.post(translationDeleteUrl, {translation_id: translationId}).done((result)=>{
		// 		if (result.status==="success"){
		// 			var modalTranslationDelete = new bootstrap.Modal(document.getElementById('ModalTranslationDelete'));
		// 	        modalTranslationDelete.hide();

		// 			return location.href = location.href;
		// 		}
		// 		alert("Failed. " + result.message);
		// 	}).fail((result)=>{
		// 		alert("Failed" + result);
		// 	});
		// },
		translationTrash(){
			var translationId = this.translationTrashModel.translationId;
			$.post(translationTrashUrl, {translation_id: translationId}).done((result)=>{
				if (result.status==="success"){
					var modalTranslationDelete = new bootstrap.Modal(document.getElementById('ModalTranslationTrash'));
					modalTranslationDelete.hide();

					return location.href = location.href;
				}
				alert("Failed. " + result.message);
			}).fail((result)=>{
				alert("Failed" + result);
			});
		},
		translationEdit(translationId){
			return location.href = translationUpdateUrl+ "&translation_id=" + translationId;
		}
	}
};
Vue.createApp(TranslationManager).mount('#translation-manager')
	`
	if m.funcLayout("") != "" {
		out := hb.NewWrap().Children([]hb.TagInterface{
			hb.NewStyleURL(cdn.JqueryDataTablesCss_1_13_4()),
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

	webpage := m.webpageComplete("Translation Manager", h)
	webpage.AddStyleURL(cdn.JqueryDataTablesCss_1_13_4())
	webpage.AddScriptURL(cdn.JqueryDataTablesJs_1_13_4())
	webpage.AddScript(inlineScript)
	responses.HTMLResponse(w, r, webpage.ToHTML())
}

func (m UiManager) translationTrashModal() *hb.Tag {
	modal := hb.NewDiv().Attr("id", "ModalTranslationTrash").Attr("class", "modal fade")
	modalDialog := hb.NewDiv().Attr("class", "modal-dialog")
	modalContent := hb.NewDiv().Attr("class", "modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("Trash Translation"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewParagraph().HTML("Are you sure you want to move this translation to trash bin?"))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-secondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Move to trash bin").Attr("class", "btn btn-danger").Attr("v-on:click", "translationTrash"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modalDialog.AddChild(modalContent)
	modal.AddChild(modalDialog)
	return modal
}

func (m UiManager) translationCreateModal() *hb.Tag {
	modal := hb.NewDiv().
		ID("ModalTranslationCreate").
		Class("modal fade")
	modalContent := hb.NewDiv().Class("modal-content")
	modalHeader := hb.NewDiv().Attr("class", "modal-header").AddChild(hb.NewHeading5().HTML("New Translation"))
	modalBody := hb.NewDiv().Attr("class", "modal-body")
	modalBody.AddChild(hb.NewDiv().Attr("class", "form-group").AddChild(hb.NewLabel().HTML("Name")).AddChild(hb.NewInput().Attr("class", "form-control").Attr("v-model", "translationCreateModel.name")))
	modalFooter := hb.NewDiv().Attr("class", "modal-footer")
	modalFooter.AddChild(hb.NewButton().HTML("Close").Attr("class", "btn btn-secondary").Attr("data-bs-dismiss", "modal"))
	modalFooter.AddChild(hb.NewButton().HTML("Create & Continue").Attr("class", "btn btn-primary").Attr("v-on:click", "translationCreate"))
	modalContent.AddChild(modalHeader).AddChild(modalBody).AddChild(modalFooter)
	modal.AddChild(hb.NewDiv().Class("modal-dialog").Children([]hb.TagInterface{
		modalContent,
	}))
	return modal
}
