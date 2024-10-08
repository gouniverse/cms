package cms

import (
	"encoding/json"
	"net/http"
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

func (cms Cms) pageTranslationsTranslationCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(utils.Req(r, "name", ""), " ")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	translation, err := cms.EntityStore.EntityCreateWithType(ENTITY_TYPE_TRANSLATION)

	if err != nil {
		api.Respond(w, r, api.Error("Translation failed to be created: "+err.Error()))
		return
	}

	if translation == nil {
		api.Respond(w, r, api.Error("Translation failed to be created"))
		return
	}

	translation.SetString("name", name)

	api.Respond(w, r, api.SuccessWithData("Translation saved successfully", map[string]interface{}{"translation_id": translation.ID()}))
}

func (cms Cms) pageTranslationsTranslationManager(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	// log.Println(endpoint)

	header := cms.cmsHeader(endpoint)
	breadcrumbs := cms.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  endpoint,
			Name: "Home",
		},
		{
			URL:  (endpoint + "?path=" + PathTranslationsTranslationManager),
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
			cms.pageTranslationsTranslationCreateModal(),
			cms.pageTranslationsTranslationTrashModal(),
		})

	translations, err := cms.EntityStore.EntityList(entitystore.EntityQueryOptions{
		EntityType: ENTITY_TYPE_TRANSLATION,
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
var translationCreateUrl = "` + endpoint + `?path=translations/translation-create-ajax"
var translationDeleteUrl = "` + endpoint + `?path=translations/translation-delete-ajax"
var translationTrashUrl = "` + endpoint + `?path=translations/translation-trash-ajax";
var translationUpdateUrl = "` + endpoint + `?path=translations/translation-update"
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
	if cms.funcLayout("") != "" {
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
		responses.HTMLResponse(w, r, cms.funcLayout(out))
		return
	}

	webpage := WebpageComplete("Translation Manager", h)
	webpage.AddStyleURL(cdn.JqueryDataTablesCss_1_13_4())
	webpage.AddScriptURL(cdn.JqueryDataTablesJs_1_13_4())
	webpage.AddScript(inlineScript)
	responses.HTMLResponse(w, r, webpage.ToHTML())
}

func (cms Cms) pageTranslationsTranslationUpdate(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	// log.Println(endpoint)

	translationID := utils.Req(r, "translation_id", "")
	if translationID == "" {
		api.Respond(w, r, api.Error("Translation ID is required"))
		return
	}

	translation, err := cms.EntityStore.EntityFindByID(translationID)

	if err != nil {
		api.Respond(w, r, api.Error("Translation failed to be retrieved: "+err.Error()))
		return
	}

	if translation == nil {
		api.Respond(w, r, api.Error("Translation NOT FOUND with ID "+translationID))
		return
	}

	header := cms.cmsHeader(r.Context().Value(keyEndpoint).(string))
	breadcrumbs := cms.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  endpoint,
			Name: "Home",
		},
		{
			URL:  (endpoint + "?path=" + PathTranslationsTranslationManager),
			Name: "Translations",
		},
		{
			URL:  (endpoint + "?path=" + PathTranslationsTranslationUpdate + "&translation_id=" + translationID),
			Name: "Edit translation",
		},
	})

	heading := hb.NewHeading1().HTML("Edit Translation")
	button := hb.NewButton().AddChild(hb.NewHTML(icons.BootstrapCheckCircle+" ")).HTML("Save").Attr("class", "btn btn-success float-end").Attr("v-on:click", "translationSave")
	heading.AddChild(button)

	formGroupStatus := hb.NewDiv().Class("form-group").Children([]hb.TagInterface{
		hb.NewLabel().HTML("Status").Class("form-label"),
		hb.NewSelect().Class("form-select").Attr("v-model", "translationModel.status").Children([]hb.TagInterface{
			hb.NewOption().Value("active").HTML("Active"),
			hb.NewOption().Value("inactive").HTML("Inactive"),
			hb.NewOption().Value("trash").HTML("Trash"),
		}),
	})

	formGroupName := hb.NewDiv().Class("form-group mt-3").Children([]hb.TagInterface{
		hb.NewLabel().HTML("Name").Class("form-label"),
		hb.NewInput().Attr("class", "form-control").Attr("v-model", "translationModel.name"),
		hb.NewParagraph().Class("text-info fs-6").HTML("This is the display name, so that you can easily find it"),
	})

	formGroupHandle := hb.NewDiv().Class("form-group mt-3").Children([]hb.TagInterface{
		hb.NewLabel().HTML("Key").Class("form-label"),
		hb.NewInput().Class("form-control").Attr("v-model", "translationModel.handle"),
		hb.NewParagraph().Class("text-info fs-6").HTML("Must be lowercase, no space, hyphens and dots allowed"),
	})

	formGroupComment := hb.NewDiv().Class("form-group mt-3").Children([]hb.TagInterface{
		hb.NewLabel().HTML("Comment").Class("form-label"),
		hb.NewTextArea().Class("form-control CodeMirror").Attr("v-model", "translationModel.comment"),
	})

	paragraphUsage := hb.NewParagraph().
		Class("text-info mt-5").
		HTML("To use this translation in your website use the following shortcode:")

	translationName, _ := translation.GetString("name", "")
	code := hb.NewCode().Children([]hb.TagInterface{
		hb.NewPRE().
			HTML(`&lt;!-- START: Translation: ` + translationName + ` -->
[[TRANSLATION_` + translation.ID() + `]]
&lt;!-- END: Translation: ` + translationName + ` -->`),
		hb.NewPRE().
			HTML(`&lt;!-- START: Translation: ` + translation.Handle() + ` -->
[[TRANSLATION_` + translation.Handle() + `]]
&lt;!-- END: Translation: ` + translationName + ` -->`),
	})
	paragraphUsage.Child(code)

	container := hb.NewDiv().ID("translation-update").Class("container").
		Children([]hb.TagInterface{
			hb.NewHTML(header),
			heading,
			hb.NewParagraph().Class("alert alert-info").
				HTML("Tools: ").
				Child(hb.NewHyperlink().HTML("Google Translate").Href("https://translate.google.com").Target("_blank")).
				HTML(", ").
				Child(hb.NewHyperlink().HTML("Bing Translate").Href("https://www.bing.com/translator").Target("_blank")).
				HTML(", ").
				Child(hb.NewHyperlink().HTML("Translateking").Href("https://translateking.com/").Target("_blank")).
				HTML(", ").
				Child(hb.NewHyperlink().HTML("Baidu Translate").Href("https://fanyi.baidu.com/").Target("_blank")).
				HTML(", ").
				Child(hb.NewHyperlink().HTML("Yandex Translate").Href("https://translate.yandex.com").Target("_blank")).
				HTML(", ").
				Child(hb.NewHyperlink().HTML("Yandex Translate").Href("https://www.reverso.net/text-translation").Target("_blank")),
			hb.NewHTML(breadcrumbs),
			formGroupStatus,
			formGroupName,
			formGroupHandle,
			formGroupComment,
		}).
		Children(lo.Map(lo.Keys(cms.translationLanguages), func(key string, index int) hb.TagInterface {
			language := cms.translationLanguages[key]
			isDefault := cms.translationLanguageDefault == key
			defaultText := lo.Ternary(isDefault, " (Default)", "")
			formGroupTranslation := hb.NewDiv().Class("form-group mt-3").Children([]hb.TagInterface{
				hb.NewLabel().
					HTML(language + " Translation" + defaultText).
					Class("form-label"),
				hb.NewTextArea().
					Class("form-control CodeMirror").
					Attr("v-model", "translationModel.translations['"+key+"']"),
			})
			return formGroupTranslation
		})).
		Child(paragraphUsage)

	h := container.ToHTML()

	handle := translation.Handle()

	name, err := translation.GetString("name", "")
	if err != nil {
		api.Respond(w, r, api.Error("Name failed to be retrieved: "+err.Error()))
		return
	}

	status, err := translation.GetString("status", "")
	if err != nil {
		api.Respond(w, r, api.Error("Status failed to be retrieved: "+err.Error()))
		return
	}

	translationContents := map[string]string{}
	lo.ForEach(lo.Keys(cms.translationLanguages), func(key string, index int) {
		translationContent, _ := translation.GetString(key, "none")
		translationContents[key] = translationContent
	})

	comment, err := translation.GetString("comment", "")
	if err != nil {
		api.Respond(w, r, api.Error("Comment failed to be retrieved: "+err.Error()))
		return
	}

	commentJSON, _ := json.Marshal(comment)
	handleJSON, _ := json.Marshal(handle)
	nameJSON, _ := json.Marshal(name)
	statusJSON, _ := json.Marshal(status)
	translationsJSON, _ := json.Marshal(translationContents)

	inlineScript := `
var translationUpdateUrl = "` + endpoint + `?path=translations/translation-update-ajax";
var translationId = "` + translationID + `";
var comment = ` + string(commentJSON) + `;
var handle = ` + string(handleJSON) + `;
var name = ` + string(nameJSON) + `;
var status = ` + string(statusJSON) + `;
var translations = ` + string(translationsJSON) + `;
const TranslationUpdate = {
	data() {
		return {
			translationModel:{
				translationId: translationId,
				comment: comment,
				handle: handle,
				translations: translations,
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
					self.translationModel.comment = editor.getValue();
				});
				$(document).on('change', '.CodeMirror', function() {
					self.translationModel.comment = editor.getValue();
				});
				setInterval(()=>{
					self.translationModel.comment = editor.getValue();
				}, 1000)
			}
		}, 500);
	},
	methods: {
		translationSave(){
			var comment = this.translationModel.comment;
			var handle = this.translationModel.handle;
			var name = this.translationModel.name;
			var translationId = this.translationModel.translationId;
			var status = this.translationModel.status;
			var translations = this.translationModel.translations;
			
			$.post(translationUpdateUrl, {
				translation_id:translationId,
				comment: comment,
				handle: handle,
				name: name,
				status: status,
				translations: translations,
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
					title: 'Translation saved',
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
Vue.createApp(TranslationUpdate).mount('#translation-update')
	`

	if cms.funcLayout("") != "" {
		out := hb.NewWrap().Children([]hb.TagInterface{
			hb.NewStyleURL(codemirrorCss),
			hb.NewStyle(`.CodeMirror {
				border: 1px solid #eee;
				height: auto;
			}`),
			hb.NewHTML(h),
			hb.NewScriptURL(cdn.Jquery_3_6_4()),
			hb.NewScriptURL(cdn.VueJs_3()),
			hb.NewScriptURL(cdn.Sweetalert2_10()),
			hb.NewScriptURL(codemirrorJs),
			hb.NewScriptURL(codemirrorXmlJs),
			hb.NewScriptURL(codemirrorHtmlmixedJs),
			hb.NewScriptURL(codemirrorJavascriptJs),
			hb.NewScriptURL(codemirrorCssJs),
			hb.NewScriptURL(codemirrorClikeJs),
			hb.NewScriptURL(codemirrorPhpJs),
			hb.NewScriptURL(codemirrorFormattingJs),
			hb.NewScriptURL(codemirrorMatchBracketsJs),
			hb.NewScript(inlineScript),
		}).ToHTML()
		responses.HTMLResponse(w, r, cms.funcLayout(out))
		return
	}

	webpage := WebpageComplete("Edit Translation", h).
		AddStyleURLs([]string{
			codemirrorCss,
		}).
		AddScriptURLs([]string{
			codemirrorJs,
			codemirrorXmlJs,
			codemirrorHtmlmixedJs,
			codemirrorJavascriptJs,
			codemirrorCssJs,
			codemirrorClikeJs,
			codemirrorPhpJs,
			codemirrorFormattingJs,
			codemirrorMatchBracketsJs,
		}).
		AddStyle(`	
.CodeMirror {
	border: 1px solid #eee;
	height: auto;
}`).
		AddScript(inlineScript)

	responses.HTMLResponse(w, r, webpage.ToHTML())
}

func (cms Cms) pageTranslationsTranslationUpdateAjax(w http.ResponseWriter, r *http.Request) {
	translationID := strings.Trim(utils.Req(r, "translation_id", ""), " ")
	comment := strings.Trim(utils.Req(r, "comment", ""), " ")
	handle := strings.Trim(utils.Req(r, "handle", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")
	translationContents := map[string]string{}
	lo.ForEach(lo.Keys(cms.translationLanguages), func(key string, index int) {
		translationContent := strings.Trim(utils.Req(r, "translations["+key+"]", ""), " ")
		translationContents[key] = translationContent
	})

	if translationID == "" {
		api.Respond(w, r, api.Error("Translation ID is required"))
		return
	}

	translation, err := cms.EntityStore.EntityFindByID(translationID)

	if err != nil {
		api.Respond(w, r, api.Error("Translation not found: "+err.Error()))
		return
	}

	if translation == nil {
		api.Respond(w, r, api.Error("Translation NOT FOUND with ID "+translationID))
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

	translation.SetHandle(handle)
	err = cms.EntityStore.EntityUpdate(*translation)

	if err != nil {
		api.Respond(w, r, api.Error("Translation failed to be updated: "+err.Error()))
		return
	}

	translation.SetString("comment", comment)
	translation.SetString("name", name)
	translation.SetAll(translationContents)
	errSetString := translation.SetString("status", status)

	if errSetString != nil {
		api.Respond(w, r, api.Error("Translation failed to be updated"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Translation saved successfully", map[string]interface{}{"translation_id": translation.ID()}))
}

func (cms Cms) pageTranslationsTranslationDeleteAjax(w http.ResponseWriter, r *http.Request) {
	translationID := strings.Trim(utils.Req(r, "translation_id", ""), " ")

	if translationID == "" {
		api.Respond(w, r, api.Error("Translation ID is required"))
		return
	}

	translation, err := cms.EntityStore.EntityFindByID(translationID)

	if err != nil {
		api.Respond(w, r, api.Error("Database error: "+err.Error()))
		return
	}

	if translation == nil {
		api.Respond(w, r, api.Success("Translation already deleted"))
		return
	}

	isOk, err := cms.EntityStore.EntityDelete(translationID)

	if err != nil {
		api.Respond(w, r, api.Error("Translation failed to be deleted: "+err.Error()))
		return
	}

	if !isOk {
		api.Respond(w, r, api.Error("Translation failed to be deleted"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Translation deleted successfully", map[string]interface{}{"translation_id": translation.ID()}))
}

func (cms Cms) pageTranslationsTranslationTrashAjax(w http.ResponseWriter, r *http.Request) {
	translationID := strings.Trim(utils.Req(r, "translation_id", ""), " ")

	if translationID == "" {
		api.Respond(w, r, api.Error("Translation ID is required"))
		return
	}

	translation, err := cms.EntityStore.EntityFindByID(translationID)

	if err != nil {
		api.Respond(w, r, api.Error("Error: "+err.Error()))
		return
	}

	if translation == nil {
		api.Respond(w, r, api.Success("Translation already deleted"))
		return
	}

	isOk, err := cms.EntityStore.EntityTrash(translationID)

	if err != nil {
		api.Respond(w, r, api.Error("Translation failed to be trashed"))
		return
	}

	if !isOk {
		api.Respond(w, r, api.Error("Translation failed to be trashed"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Translation trashed successfully", map[string]any{
		"translation_id": translation.ID(),
	}))
}

func (cms Cms) pageTranslationsTranslationTrashModal() *hb.Tag {
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

func (cms Cms) pageTranslationsTranslationCreateModal() *hb.Tag {
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
