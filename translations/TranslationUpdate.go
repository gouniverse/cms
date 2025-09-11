package cms

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/gouniverse/icons"
	"github.com/gouniverse/responses"
	"github.com/samber/lo"
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

func (m UiManager) TranslationUpdate(w http.ResponseWriter, r *http.Request) {
	translationID := req.GetStringTrimmed(r, "translation_id")
	if translationID == "" {
		api.Respond(w, r, api.Error("Translation ID is required"))
		return
	}

	translation, err := m.entityStore.EntityFindByID(translationID)

	if err != nil {
		api.Respond(w, r, api.Error("Translation failed to be retrieved: "+err.Error()))
		return
	}

	if translation == nil {
		api.Respond(w, r, api.Error("Translation NOT FOUND with ID "+translationID))
		return
	}

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
		{
			URL:  (m.endpoint + "?path=" + m.pathTranslationsTranslationUpdate + "&translation_id=" + translationID),
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
		Children(lo.Map(lo.Keys(m.translationLanguages), func(key string, index int) hb.TagInterface {
			language := m.translationLanguages[key]
			isDefault := m.translationLanguageDefault == key
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
	lo.ForEach(lo.Keys(m.translationLanguages), func(key string, index int) {
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
var translationUpdateUrl = "` + m.endpoint + `?path=translations/translation-update-ajax";
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

	webpage := m.webpageComplete("Edit Translation", h).
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
