package cms

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/gouniverse/api"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/icons"
	"github.com/gouniverse/responses"
	"github.com/gouniverse/utils"
)

const codemirrorCss = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/codemirror.min.css"
const codemirrorJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/codemirror.min.js"
const codemirrorXmlJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/xml/xml.min.js"
const codemirrorHtmlmixedJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/htmlmixed/htmlmixed.min.js"
const codemirrorJavascriptJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/javascript/javascript.js"
const codemirrorCssJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/css/css.js"
const codemirrorClikeJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/clike/clike.min.js"
const codemirrorPhpJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.20.0/mode/php/php.min.js"
const codemirrorFormattingJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/2.36.0/formatting.min.js"
const codemirrorMatchBracketsJs = "//cdnjs.cloudflare.com/ajax/libs/codemirror/3.22.0/addon/edit/matchbrackets.min.js"

func (m UiManager) SettingUpdate(w http.ResponseWriter, r *http.Request) {
	settingKey := utils.Req(r, "setting_key", "")
	if settingKey == "" {
		api.Respond(w, r, api.Error("Setting key is required"))
		return
	}

	settingValue, err := m.settingStore.Get(r.Context(), settingKey, "%%NOTFOUND%%")

	if err != nil {
		api.Respond(w, r, api.Error("There was an error: "+err.Error()))
		return
	}

	if settingValue == "%%NOTFOUND%%" {
		api.Respond(w, r, api.Error("Setting NOT FOUND with key "+settingKey))
		return
	}

	header := m.cmsHeader(m.endpoint)
	breadcrumbs := m.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  m.endpoint,
			Name: "Home",
		},
		{
			URL:  (m.endpoint + "?path=" + m.pathSettingsSettingManager),
			Name: "Settings",
		},
		{
			URL:  (m.endpoint + "?path=" + m.pathSettingsSettingUpdate + "&setting_key=" + settingKey),
			Name: "Edit setting",
		},
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
	container.AddChild(hb.NewHTML(breadcrumbs))
	container.AddChild(formGroupName).AddChild(formGroupContent).AddChild(paragraphUsage)

	h := container.ToHTML()

	settingValueJSON, _ := json.Marshal(settingValue)

	inlineScript := `
var settingUpdateUrl = "` + m.endpoint + `?path=settings/setting-update-ajax";
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
					return Swal.fire({icon: 'error',title: 'Oops...',text: response.message,});
				}

				return Swal.fire({icon: 'success',title: 'Setting saved',});
			}).fail((result)=>{
				console.log(result);
				return Swal.fire({icon: 'error',title: 'Oops...',text: result,});
			});
		}
	}
};
Vue.createApp(SettingUpdate).mount('#setting-update')
	`

	if m.funcLayout("") != "" {
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
		responses.HTMLResponse(w, r, m.funcLayout(out))
		return
	}

	webpage := m.webpageComplete("Edit Setting", h).
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
