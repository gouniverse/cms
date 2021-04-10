package cms

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	hb "github.com/gouniverse/html"
	"github.com/gouniverse/icons"
	"github.com/gouniverse/utils"
)

func pageSettingsSettingCreateAjax(w http.ResponseWriter, r *http.Request) {
	key := strings.Trim(utils.Req(r, "key", ""), " ")

	if key == "" {
		api.Respond(w, r, api.Error("Key is required field"))
		return
	}

	isOk := GetSettingStore().Set(key, "")

	if isOk == false {
		api.Respond(w, r, api.Error("Setting failed to be created"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Setting saved successfully", map[string]interface{}{"setting_key": key}))
	return
}

func pageSettingsSettingManager(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	log.Println(endpoint)

	header := cmsHeader(endpoint)
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathSettingsSettingManager): "Settings",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "setting-manager")
	heading := hb.NewHeading1().HTML("Setting Manager")
	button := hb.NewButton().HTML("New setting").Attr("class", "btn btn-success float-end").Attr("v-on:click", "showSettingCreateModal")
	heading.AddChild(button)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))

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
	container.AddChild(modal)

	settingKeys := GetSettingStore().Keys()

	table := hb.NewTable().Attr("class", "table table-responsive table-striped mt-3")
	thead := hb.NewThead()
	tbody := hb.NewTbody()
	table.AddChild(thead).AddChild(tbody)

	tr := hb.NewTR()
	th1 := hb.NewTD().HTML("Name")
	th3 := hb.NewTD().HTML("Actions").Attr("style", "width:1px;")
	thead.AddChild(tr.AddChild(th1).AddChild(th3))

	for _, settingKey := range settingKeys {
		name := settingKey
		buttonEdit := hb.NewButton().HTML("Edit").Attr("type", "button").Attr("class", "btn btn-primary").Attr("v-on:click", "settingEdit('"+settingKey+"')")

		tr := hb.NewTR()
		td1 := hb.NewTD().HTML(name)
		td3 := hb.NewTD().AddChild(buttonEdit)

		tbody.AddChild(tr.AddChild(td1).AddChild(td3))
	}
	container.AddChild(table)

	h := container.ToHTML()

	inlineScript := `
var settingCreateUrl = "` + endpoint + `?path=settings/setting-create-ajax"
var settingUpdateUrl = "` + endpoint + `?path=settings/setting-update"
const SettingManager = {
	data() {
		return {
		  settingCreateModel:{
			  key:""
		  }
		}
	},
	methods: {
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
		}
	}
};
Vue.createApp(SettingManager).mount('#setting-manager')
	`

	webpage := Webpage("Setting Manager", h)
	webpage.AddScript(inlineScript)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

func pageSettingsSettingUpdate(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	log.Println(endpoint)

	settingKey := utils.Req(r, "setting_key", "")
	if settingKey == "" {
		api.Respond(w, r, api.Error("Setting key is required"))
		return
	}

	settingValue := GetSettingStore().Get(settingKey, "%%NOTFOUND%%")

	if settingValue == "%%NOTFOUND%%" {
		api.Respond(w, r, api.Error("Setting NOT FOUND with key "+settingKey))
		return
	}

	header := cmsHeader(r.Context().Value(keyEndpoint).(string))
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
		(endpoint + "?path=" + PathSettingsSettingManager):                               "Settings",
		(endpoint + "?path=" + PathSettingsSettingUpdate + "&setting_key=" + settingKey): "Edit setting",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "setting-update")
	heading := hb.NewHeading1().HTML("Edit Setting")
	button := hb.NewButton().AddChild(hb.NewHTML(icons.BootstrapCheckCircle+" ")).HTML("Save").Attr("class", "btn btn-success float-end").Attr("v-on:click", "settingSave")
	heading.AddChild(button)

	formGroupName := hb.NewDiv().Attr("class", "form-group")
	formGroupNameLabel := hb.NewLabel().HTML("Key").Attr("class", "form-label")
	formGroupNameInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "settingModel.settingKey").Attr("disabled", "disabled")
	formGroupName.AddChild(formGroupNameLabel)
	formGroupName.AddChild(formGroupNameInput)

	formGroupContent := hb.NewDiv().Attr("class", "form-group")
	formGroupContentLabel := hb.NewLabel().HTML("Value").Attr("class", "form-label")
	formGroupContentInput := hb.NewTextArea().Attr("class", "form-control CodeMirror").Attr("v-model", "settingModel.settingValue")
	formGroupContent.AddChild(formGroupContentLabel)
	formGroupContent.AddChild(formGroupContentInput)

	paragraphUsage := hb.NewParagraph().Attr("class", "text-info mt-5").AddChild(hb.NewHTML("To use this block in your website use the following shortcode:"))
	code := hb.NewCode().AddChild(hb.NewPRE().HTML(`&lt;!-- START: Setting: ` + settingKey + ` -->
[[SETTING_` + settingKey + `]]
&lt;!-- END: Setting: ` + settingKey + ` -->`))
	paragraphUsage.AddChild(code)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))
	container.AddChild(formGroupName).AddChild(formGroupContent)

	h := container.ToHTML()

	settingValueJSON, _ := json.Marshal(settingValue)

	inlineScript := `
var settingUpdateUrl = "` + endpoint + `?path=settings/setting-update-ajax";
var settingKey = "` + settingKey + `";
var settingValue = ` + string(settingValueJSON) + `;
const SettingUpdate = {
	data() {
		return {
			settingModel:{
				settingKey: settingKey,
				settingValue: settingValue,
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
					self.settingModel.settingValue = editor.getValue();
				});
				$(document).on('change', '.CodeMirror', function() {
					self.settingModel.settingValue = editor.getValue();
				});
				setInterval(()=>{
					self.settingModel.settingValue = editor.getValue();
				}, 1000)
			}
		}, 500);
	},
	methods: {
		settingSave(){
			var settingKey = this.settingModel.settingKey;
			var settingValue = this.settingModel.settingValue;
			
			$.post(settingUpdateUrl, {
				setting_key:settingKey,
				setting_value: settingValue,
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
Vue.createApp(SettingUpdate).mount('#setting-update')
	`

	webtemplate := Webpage("Edit Setting", h)
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

func pageSettingsSettingUpdateAjax(w http.ResponseWriter, r *http.Request) {
	settingKey := utils.Req(r, "setting_key", "")
	settingValue := utils.Req(r, "setting_value", "%%NOTSENT%%")
	if settingKey == "" {
		api.Respond(w, r, api.Error("Setting key is required"))
		return
	}
	if settingValue == "%%NOTSENT%%" {
		api.Respond(w, r, api.Error("Setting value is required"))
		return
	}

	isOk := GetSettingStore().Set(settingKey, settingValue)

	if isOk == false {
		api.Respond(w, r, api.Error("Setting failed to be updated"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Setting saved successfully", map[string]interface{}{"setting_key": settingKey}))
	return
}
