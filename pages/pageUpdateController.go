package cms

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/blockeditor"
	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/entitystore"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/gouniverse/cms/types"
	"github.com/gouniverse/router"
	"github.com/samber/lo"
)

func (m UiManager) PageUpdate(w http.ResponseWriter, r *http.Request) {
	result := NewPageUpdateController(m).Handler(w, r)
	w.Write([]byte(result))
}

const VIEW_SETTINGS = "settings"
const VIEW_CONTENT = "content"
const VIEW_SEO = "seo"
const ACTION_BLOCKEDITOR_HANDLE = "blockeditor_handle"

type pageUpdateController struct {
	m UiManager
}

var _ router.HTMLControllerInterface = (*pageUpdateController)(nil)

func NewPageUpdateController(m UiManager) *pageUpdateController {
	return &pageUpdateController{m: m}
}

func (controller pageUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareDataAndValidate(r)

	if errorMessage != "" {
		//return helpers.ToFlashError(w, r, errorMessage, shared.NewLinks().Pages(map[string]string{}), 10)
		return api.Error(errorMessage).ToString()
	}

	if data.action == ACTION_BLOCKEDITOR_HANDLE {
		return blockeditor.Handle(w, r, controller.m.blockEditorDefinitions)
	}

	if r.Method == http.MethodPost {
		return controller.form(data).ToHTML()
	}

	h := controller.page(data)

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

	if controller.m.funcLayout("") != "" {
		out := hb.NewWrap().
			Children([]hb.TagInterface{
				hb.NewStyleURL(codemirrorCss),
				hb.NewStyleURL(cdn.TrumbowygCss_2_27_3()),
				hb.NewStyle(`.CodeMirror {
				border: 1px solid #eee;
				height: auto;
			}`),
				h,
				hb.NewScriptURL(cdn.Jquery_3_6_4()),
				hb.NewScriptURL(cdn.VueJs_3()),
				hb.NewScriptURL(cdn.Sweetalert2_11()),
				hb.NewScriptURL(cdn.Htmx_2_0_0()),
				hb.NewScriptURL(cdn.TrumbowygJs_2_27_3()),
				hb.NewScriptURL(codemirrorJs),
				hb.NewScriptURL(codemirrorXmlJs),
				hb.NewScriptURL(codemirrorHtmlmixedJs),
				hb.NewScriptURL(codemirrorJavascriptJs),
				hb.NewScriptURL(codemirrorCssJs),
				hb.NewScriptURL(codemirrorClikeJs),
				hb.NewScriptURL(codemirrorPhpJs),
				hb.NewScriptURL(codemirrorFormattingJs),
				hb.NewScriptURL(codemirrorMatchBracketsJs),
				hb.NewScript(controller.script()),
			}).ToHTML()
		return controller.m.funcLayout(out)
	}

	webpage := controller.m.webpageComplete("Edit Page", h.ToHTML())
	webpage.AddStyleURLs([]string{
		codemirrorCss,
		cdn.TrumbowygCss_2_27_3(),
	})
	webpage.AddScriptURLs([]string{
		cdn.Htmx_2_0_0(),
		cdn.TrumbowygJs_2_27_3(),
		codemirrorJs,
		codemirrorXmlJs,
		codemirrorHtmlmixedJs,
		codemirrorJavascriptJs,
		codemirrorCssJs,
		codemirrorClikeJs,
		codemirrorPhpJs,
		codemirrorFormattingJs,
		codemirrorMatchBracketsJs,
	})
	webpage.AddStyle(`	
.CodeMirror {
	border: 1px solid #eee;
	height: auto;
}
	`)
	webpage.AddScript(controller.script())

	return controller.m.funcLayout(webpage.ToHTML())

	// return layouts.NewAdminLayout(r, layouts.Options{
	// 	Title:   "Edit page | Blog",
	// 	Content: controller.page(data),
	// 	ScriptURLs: []string{
	// 		cdn.Jquery_3_7_1(),
	// 		cdn.TrumbowygJs_2_27_3(),
	// 		cdn.Sweetalert2_10(),
	// 		cdn.JqueryUiJs_1_13_1(), // needed for BlockArea
	// 		links.NewWebsiteLinks().Resource(`/js/blockarea_v0200.js`, map[string]string{}), // needed for BlockArea
	// 		codemirrorJs,
	// 		codemirrorXmlJs,
	// 		codemirrorHtmlmixedJs,
	// 		codemirrorJavascriptJs,
	// 		codemirrorCssJs,
	// 		codemirrorClikeJs,
	// 		codemirrorPhpJs,
	// 		codemirrorFormattingJs,
	// 		codemirrorMatchBracketsJs,
	// 	},
	// 	Scripts: []string{
	// 		controller.script(),
	// 	},
	// 	StyleURLs: []string{
	// 		cdn.TrumbowygCss_2_27_3(),
	// 		cdn.JqueryUiCss_1_13_1(), // needed for BlockArea
	// 		codemirrorCss,
	// 	},
	// 	Styles: []string{
	// 		`.CodeMirror {
	// 			border: 1px solid #eee;
	// 			height: auto;
	// 		}`,
	// 	},
	// }).ToHTML()
}

func (controller pageUpdateController) script() string {
	js := ``
	return js
}

func (controller pageUpdateController) page(data pageUpdateControllerData) hb.TagInterface {
	buttonSave := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-save").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Save").
		HxInclude("#FormpageUpdate").
		HxPost(controller.m.url(controller.m.pathPagesPageUpdate, map[string]string{"page_id": data.pageID})).
		HxTarget("#FormpageUpdate")

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(controller.m.url(controller.m.pathPagesPageManager, map[string]string{}))

	heading := hb.Heading1().
		Text("Edit page:").
		Text(" ").
		Text(data.page.Title()).
		Child(buttonSave).
		Child(buttonCancel)

	card := hb.Div().
		Class("card").
		Child(
			hb.Div().
				Class("card-header").
				Style(`display:flex;justify-content:space-between;align-items:center;`).
				Child(hb.Heading4().
					HTMLIf(data.view == VIEW_CONTENT, "Web Page Contents").
					HTMLIf(data.view == VIEW_SEO, "Web Page SEO").
					HTMLIf(data.view == VIEW_SETTINGS, "Web Page Settings").
					Style("margin-bottom:0;display:inline-block;")).
				Child(buttonSave),
		).
		Child(
			hb.Div().
				Class("card-body").
				Child(controller.form(data)))

	tabs := bs.NavTabs().
		Class("mb-3").
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_CONTENT, "active").
				Href(controller.m.url(controller.m.pathPagesPageUpdate, map[string]string{
					"page_id": data.pageID,
					"view":    VIEW_CONTENT,
				})).
				HTML("Content"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_SEO, "active").
				Href(controller.m.url(controller.m.pathPagesPageUpdate, map[string]string{
					"page_id": data.pageID,
					"view":    VIEW_SEO,
				})).
				HTML("SEO"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_SETTINGS, "active").
				Href(controller.m.url(controller.m.pathPagesPageUpdate, map[string]string{
					"page_id": data.pageID,
					"view":    VIEW_SETTINGS,
				})).
				HTML("Settings")))

	header := controller.m.cmsHeader(controller.m.endpoint)
	breadcrumbs := controller.m.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  controller.m.url("", map[string]string{}),
			Name: "Home",
		},
		{
			URL:  controller.m.url(controller.m.pathPagesPageManager, map[string]string{}),
			Name: "Webpage Manager",
		},
		{
			URL:  controller.m.url(controller.m.pathPagesPageUpdate, map[string]string{"page_id": data.pageID}),
			Name: "Edit page",
		},
	})

	return hb.Div().
		Class("container").
		HTML(header).
		Child(heading).
		HTML(breadcrumbs).
		// Child(pageTitle).
		Child(tabs).
		Child(card)
}

func (controller pageUpdateController) form(data pageUpdateControllerData) hb.TagInterface {
	fieldsSettings := []form.FieldInterface{
		&form.Field{
			Label: "Status",
			Name:  "page_status",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formStatus,
			Help:  "The status of this webpage. Published pages will be displayed on the website.",
			Options: []form.FieldOption{
				{
					Value: "- not selected -",
					Key:   "",
				},
				{
					Value: "Draft",
					Key:   types.WEBPAGE_STATUS_DRAFT,
				},
				{
					Value: "Published",
					Key:   types.WEBPAGE_STATUS_ACTIVE,
				},
				{
					Value: "Unpublished",
					Key:   types.WEBPAGE_STATUS_INACTIVE,
				},
				{
					Value: "In Trash Bin",
					Key:   types.WEBPAGE_STATUS_DELETED,
				},
			},
		},
		&form.Field{
			Label: "Template ID",
			Name:  "page_template_id",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formTemplateID,
			Help:  "The template that this page content will be displayed in. This feature is useful if you want to implement consistent layouts. Leaving the template empty will display the page content as it is, standalone",
			OptionsF: func() []form.FieldOption {
				options := []form.FieldOption{
					{
						Value: "- not template selected, page content will be displayed as it is -",
						Key:   "",
					},
				}
				for _, template := range data.templateList {
					name, _ := template.GetString("name", "")
					options = append(options, form.FieldOption{
						Value: name,
						Key:   template.ID(),
					})
				}
				return options

			},
		},
		// {
		// 	Label: "Image URL",
		// 	Name:  "page_image_url",
		// 	Type:  form.FORM_FIELD_TYPE_IMAGE,
		// 	Value: data.formImageUrl,
		// 	Help:  "The image that will be displayed on the blog page. If left empty, the default image will be used.",
		// },
		// {
		// 	Label: "Featured",
		// 	Name:  "page_featured",
		// 	Type:  form.FORM_FIELD_TYPE_SELECT,
		// 	Value: data.formFeatured,
		// 	Help:  "Is this blog page featured? Featured pages will be displayed on the home page.",
		// 	Options: []form.FieldOption{
		// 		{
		// 			Value: "- not selected -",
		// 			Key:   "",
		// 		},
		// 		{
		// 			Value: "No",
		// 			Key:   "no",
		// 		},
		// 		{
		// 			Value: "Yes",
		// 			Key:   "yes",
		// 		},
		// 	},
		// },
		// {
		// 	Label: "Published At",
		// 	Name:  "page_published_at",
		// 	Type:  form.FORM_FIELD_TYPE_DATETIME,
		// 	Value: data.formPublishedAt,
		// 	Help:  "The date this blog page was published.",
		// },
		&form.Field{
			Label: "Editor",
			Name:  "page_editor",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formEditor,
			Help:  "The content editor that will be used while editing this webpage content. Once set, this should not be changed, or the content may be lost. If left empty, the default editor (textarea) will be used. Note you will need to save and refresh to activate",
			OptionsF: func() []form.FieldOption {
				options := []form.FieldOption{
					{
						Value: "- not selected -",
						Key:   "",
					},
				}

				options = append(options, form.FieldOption{
					Value: "CodeMirror (HTML Source Editor)",
					Key:   types.WEBPAGE_EDITOR_CODEMIRROR,
				})

				if len(controller.m.blockEditorDefinitions) > 0 {
					options = append(options, form.FieldOption{
						Value: "BlockEditor (Visual Editor using Blocks)",
						Key:   types.WEBPAGE_EDITOR_BLOCKEDITOR,
					})
				}

				options = append(options, form.FieldOption{
					Value: "Markdown (Simple Textarea)",
					Key:   types.WEBPAGE_EDITOR_MARKDOWN,
				})

				options = append(options, form.FieldOption{
					Value: "HTML Area (WYSIWYG)",
					Key:   types.WEBPAGE_EDITOR_HTMLAREA,
				})

				options = append(options, form.FieldOption{
					Value: "Text Area",
					Key:   types.WEBPAGE_EDITOR_TEXTAREA,
				})

				return options
			},
		},
		&form.Field{
			Label: "Webpage Name",
			Name:  "page_name",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formName,
			Help:  "The name of the page as displayed in the admin panel. This is not vsible to the page vistors",
		},
		// {
		// 	Label: "Admin Notes",
		// 	Name:  "page_memo",
		// 	Type:  form.FORM_FIELD_TYPE_TEXTAREA,
		// 	Value: data.formMemo,
		// 	Help:  "Admin notes for this blogpage. These notes will not be visible to the public.",
		// },
		&form.Field{
			Label:    "Webpage ID",
			Name:     "page_id",
			Type:     form.FORM_FIELD_TYPE_STRING,
			Value:    data.pageID,
			Readonly: true,
			Help:     "The reference number (ID) of the webpage. This is used to identify the webpage in the system and should not be changed.",
		},
		&form.Field{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    data.view,
			Readonly: true,
		},
	}

	editor := lo.IfF(data.page != nil, func() string { return data.page.Editor() }).Else("")

	fieldContent := form.Field{
		Label:   "Content",
		Name:    "page_content",
		Type:    form.FORM_FIELD_TYPE_TEXTAREA,
		Value:   data.formContent,
		Help:    "The content of this webpage. This will be displayed in the browser. If template is selected, the content will be displayed inside the template.",
		Options: []form.FieldOption{},
	}

	if editor == types.WEBPAGE_EDITOR_CODEMIRROR {
		//fieldContent.Type = form.FORM_FIELD_TYPE_CODEMIRROR
		fieldContent.Options = []form.FieldOption{}
	}

	// For HTML Area editor, configure the Trumbowyg editor
	if editor == types.WEBPAGE_EDITOR_HTMLAREA {
		htmlAreaFieldOptions := []form.FieldOption{
			{
				Key: "config",
				Value: `{
	btns: [
		['viewHTML'],
		['undo', 'redo'],
		['formatting'],
		['strong', 'em', 'del'],
		['superscript', 'subscript'],
		['link','justifyLeft','justifyRight','justifyCenter','justifyFull'],
		['unorderedList', 'orderedList'],
		['insertImage'],
		['removeformat'],
		['horizontalRule'],
		['fullscreen'],
	],
	autogrow: true,
	removeformatPasted: true,
	tagsToRemove: ['script', 'link', 'embed', 'iframe', 'input'],
	tagsToKeep: ['hr', 'img', 'i'],
	autogrowOnEnter: true,
	linkTargets: ['_blank'],
	}`,
			}}
		fieldContent.Type = form.FORM_FIELD_TYPE_HTMLAREA
		fieldContent.Options = htmlAreaFieldOptions
	}

	if editor == types.WEBPAGE_EDITOR_BLOCKEDITOR {
		value := fieldContent.Value

		if value == "" {
			value = `[]`
		}

		editor, err := blockeditor.NewEditor(blockeditor.NewEditorOptions{
			// ID:    "blockeditor" + uid.HumanUid(),
			Name:  fieldContent.Name,
			Value: value,
			HandleEndpoint: controller.m.url(controller.m.pathPagesPageUpdate, map[string]string{
				"page_id": data.pageID,
				"action":  ACTION_BLOCKEDITOR_HANDLE,
			}),
			BlockDefinitions: controller.m.blockEditorDefinitions,
		})

		if err != nil {
			return hb.Div().Class("alert alert-danger").Text("Error creating blockeditor: ").Text(err.Error())
		}

		fieldContent.Type = form.FORM_FIELD_TYPE_BLOCKEDITOR
		fieldContent.CustomInput = editor
	}

	fieldsContent := []form.FieldInterface{
		&form.Field{
			Label: "Title",
			Name:  "page_title",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formTitle,
			Help:  "The title of this blog as will be seen everywhere",
		},
		&fieldContent,
		&form.Field{
			Label:    "page ID",
			Name:     "page_id",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    data.pageID,
			Readonly: true,
		},
		&form.Field{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    VIEW_CONTENT,
			Readonly: true,
		},
	}

	if editor == types.WEBPAGE_EDITOR_MARKDOWN {
		contentScript := hb.Script(`
setTimeout(() => {
	const textArea = document.querySelector('textarea[name="page_content"]');
	textArea.style.height = '300px';
}, 2000)
			`).
			ToHTML()

		fieldsContent = append(fieldsContent, &form.Field{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: contentScript,
		})
	}

	if editor == types.WEBPAGE_EDITOR_CODEMIRROR {
		contentScript := hb.Script(`
function codeMirrorSelector() {
	return 'textarea[name="page_content"]';
}
function getCodeMirrorEditor() {
	return document.querySelector(codeMirrorSelector());
}
setTimeout(function () {
    console.log(getCodeMirrorEditor());
	if (getCodeMirrorEditor()) {
		var editor = CodeMirror.fromTextArea(getCodeMirrorEditor(), {
			lineNumbers: true,
			matchBrackets: true,
			mode: "application/x-httpd-php",
			indentUnit: 4,
			indentWithTabs: true,
			enterMode: "keep", tabMode: "shift"
		});
		$(document).on('mouseup', codeMirrorSelector(), function() {
			getCodeMirrorEditor().value = editor.getValue();
		});
		$(document).on('change', codeMirrorSelector(), function() {
			getCodeMirrorEditor().value = editor.getValue();
		});
		setInterval(()=>{
			getCodeMirrorEditor().value = editor.getValue();
		}, 1000)
	}
}, 500);
		`).ToHTML()

		fieldsContent = append(fieldsContent, &form.Field{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: contentScript,
		})
	}

	fieldsSEO := controller.fieldsSEO(data)

	formpageUpdate := form.NewForm(form.FormOptions{
		ID: "FormpageUpdate",
	})

	if data.view == VIEW_SETTINGS {
		formpageUpdate.SetFields(fieldsSettings)
	}

	if data.view == VIEW_CONTENT {
		formpageUpdate.SetFields(fieldsContent)
	}

	if data.view == VIEW_SEO {
		formpageUpdate.SetFields(fieldsSEO)
	}

	if data.formErrorMessage != "" {
		formpageUpdate.AddField(&form.Field{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{Icon: "error", Text: data.formErrorMessage}).ToHTML(),
		})
	}

	if data.formSuccessMessage != "" {
		formpageUpdate.AddField(&form.Field{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{Icon: "success", Text: data.formSuccessMessage}).ToHTML(),
		})
	}

	return formpageUpdate.Build()

	// required := hb.Sup().HTML("required").Style("color:red;margin-left:10px;")

	// // Status
	// fomrGroupStatus := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Status").Child(required)).
	// 	Child(bs.FormSelect().
	// 		Name("page_status").
	// 		Child(bs.FormSelectOption("", "").
	// 			AttrIf(data.formStatus == "", "selected", "")).
	// 		Child(bs.FormSelectOption(blogstore.page_STATUS_DRAFT, "Draft").
	// 			AttrIf(data.formStatus == blogstore.page_STATUS_DRAFT, "selected", "selected")).
	// 		Child(bs.FormSelectOption(blogstore.page_STATUS_PUBLISHED, "Published").
	// 			AttrIf(data.formStatus == blogstore.page_STATUS_PUBLISHED, "selected", "selected")).
	// 		Child(bs.FormSelectOption(blogstore.page_STATUS_UNPUBLISHED, "Unpublished").
	// 			AttrIf(data.formStatus == blogstore.page_STATUS_UNPUBLISHED, "selected", "selected")).
	// 		Child(bs.FormSelectOption(blogstore.page_STATUS_TRASH, "Trashed").
	// 			AttrIf(data.formStatus == blogstore.page_STATUS_TRASH, "selected", "selected")),
	// 	)

	// // Admin Notes
	// formGroupMemo := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Admin Notes")).
	// 	Child(bs.FormTextArea().
	// 		Name("page_memo").
	// 		HTML(data.formMemo).
	// 		Style("height:100px;"),
	// 	)

	// // page ID
	// formGrouppageId := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("page ID")).
	// 	Child(bs.FormInput().
	// 		Type(hb.TYPE_TEXT).
	// 		Name("page_id").
	// 		Value(data.pageID).
	// 		Attr("readonly", ""),
	// 	)

	// // Title
	// formGroupTitle := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Title").Child(required)).
	// 	Child(bs.FormInput().
	// 		Type("text").
	// 		Name("page_title").
	// 		Value(data.formTitle).
	// 		Style("width:100%;"),
	// 	)

	// // Summary
	// formGroupSummary := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Summary")).
	// 	Child(bs.FormTextArea().
	// 		Type("text").
	// 		Name("page_summary").
	// 		HTML(data.formSummary).
	// 		Style("width:100%;"),
	// 	)

	// // Published Date
	// formGroupPublishedAt := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Published Date")).
	// 	Child(bs.FormInput().
	// 		Type(hb.TYPE_TEXT).
	// 		Name("page_published_at").
	// 		Value(data.formPublishedAt).
	// 		Style("width:100%;"),
	// 	)

	// // Featured
	// formGroupFeatured := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Featured")).
	// 	Child(bs.FormSelect().
	// 		Name("page_featured").
	// 		Child(bs.FormSelectOption("", "").
	// 			AttrIf(data.formFeatured == "", "selected", "")).
	// 		Child(bs.FormSelectOption("yes", "Yes").
	// 			AttrIf(data.formFeatured == "yes", "selected", "selected")).
	// 		Child(bs.FormSelectOption("no", "No").
	// 			AttrIf(data.formFeatured == "no", "selected", "selected")),
	// 	)

	// form := hb.Form().
	// 	ID("FormpageUpdate").
	// 	Child(formGroupTitle).
	// 	Child(fomrGroupStatus).
	// 	Child(formGroupSummary).
	// 	Child(formGroupPublishedAt).
	// 	Child(formGroupFeatured).
	// 	Child(formGroupMemo).
	// 	Child(formGrouppageId)

	// if data.formErrorMessage != "" {
	// 	form.Child(hb.Swal(hb.SwalOptions{Icon: "error", Text: data.formErrorMessage}))
	// }

	// if data.formSuccessMessage != "" {
	// 	form.Child(hb.Swal(hb.SwalOptions{Icon: "success", Text: data.formSuccessMessage}))
	// }

	// return form
}

func (pageUpdateController) fieldsSEO(data pageUpdateControllerData) []form.FieldInterface {
	fieldsSEO := []form.FieldInterface{
		&form.Field{
			Label: "Alias / Path / User Friendly URL",
			Name:  "page_alias",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formAlias,
			Help:  "The relative path on the website where this page will be visible to the vistors. Once set do not change it as search engines will look for this path.",
		},
		&form.Field{
			Label: "Meta Description",
			Name:  "page_meta_description",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formMetaDescription,
			Help:  "The description of this webpage as will be seen in search engines.",
		},
		&form.Field{
			Label: "Meta Keywords",
			Name:  "page_meta_keywords",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formMetaKeywords,
			Help:  "Specifies the keywords that will be used by the search engines to find this webpage. Separate keywords with commas.",
		},
		&form.Field{
			Label: "Meta Robots",
			Name:  "page_meta_robots",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formMetaRobots,
			Help:  "Specifies if this webpage should be indexed by the search engines. Index, Follow, means all. NoIndex, NoFollow means none.",
			Options: []form.FieldOption{
				{
					Value: "- not selected -",
					Key:   "",
				},
				{
					Value: "INDEX, FOLLOW",
					Key:   "INDEX, FOLLOW",
				},
				{
					Value: "NOINDEX, FOLLOW",
					Key:   "NOINDEX, FOLLOW",
				},
				{
					Value: "INDEX, NOFOLLOW",
					Key:   "INDEX, NOFOLLOW",
				},
				{
					Value: "NOINDEX, NOFOLLOW",
					Key:   "NOINDEX, NOFOLLOW",
				},
			},
		},
		&form.Field{
			Label: "Canonical URL",
			Name:  "page_canonical_url",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formCanonicalURL,
			Help:  "The canonical URL for this webpage. This is used by the search engines to display the preferred version of the web page in search results.",
		},
		&form.Field{
			Label:    "Webpage ID",
			Name:     "page_id",
			Type:     form.FORM_FIELD_TYPE_STRING,
			Value:    data.pageID,
			Readonly: true,
		},
		&form.Field{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    VIEW_SEO,
			Readonly: true,
		},
	}
	return fieldsSEO
}

func (controller pageUpdateController) savepage(r *http.Request, data pageUpdateControllerData) (d pageUpdateControllerData, errorMessage string) {
	data.formAlias = req.GetStringTrimmed(r, "page_alias")
	data.formCanonicalURL = req.GetStringTrimmed(r, "page_canonical_url")
	data.formContent = req.GetStringTrimmed(r, "page_content")
	data.formEditor = req.GetStringTrimmed(r, "page_editor")
	data.formFeatured = req.GetStringTrimmed(r, "page_featured")
	data.formImageUrl = req.GetStringTrimmed(r, "page_image_url")
	data.formMemo = req.GetStringTrimmed(r, "page_memo")
	data.formMetaDescription = req.GetStringTrimmed(r, "page_meta_description")
	data.formMetaKeywords = req.GetStringTrimmed(r, "page_meta_keywords")
	data.formMetaRobots = req.GetStringTrimmed(r, "page_meta_robots")
	data.formName = req.GetStringTrimmed(r, "page_name")
	data.formPublishedAt = req.GetStringTrimmed(r, "page_published_at")
	data.formSummary = req.GetStringTrimmed(r, "page_summary")
	data.formStatus = req.GetStringTrimmed(r, "page_status")
	data.formTitle = req.GetStringTrimmed(r, "page_title")
	data.formTemplateID = req.GetStringTrimmed(r, "page_template_id")

	if data.view == VIEW_SETTINGS {
		if data.formStatus == "" {
			data.formErrorMessage = "Status is required"
			return data, ""
		}
	}

	if data.view == VIEW_CONTENT {
		if data.formTitle == "" {
			data.formErrorMessage = "Title is required"
			return data, ""
		}
	}

	if data.view == VIEW_SETTINGS {
		// make sure the date is in the correct format
		// data.formPublishedAt = lo.Substring(strings.ReplaceAll(data.formPublishedAt, " ", "T")+":00", 0, 19)
		// publishedAt := lo.Ternary(data.formPublishedAt == "", sb.NULL_DATE, carbon.Parse(data.formPublishedAt).ToDateTimeString(carbon.UTC))
		data.page.SetEditor(data.formEditor)
		// data.page.SetFeatured(data.formFeatured)
		// data.page.SetImageUrl(data.formImageUrl)
		// data.page.SetMemo(data.formMemo)
		data.page.SetName(data.formName)
		// data.page.SetPublishedAt(publishedAt)
		data.page.SetStatus(data.formStatus)
		data.page.SetTemplateID(data.formTemplateID)
	}

	if data.view == VIEW_CONTENT {
		data.page.SetContent(data.formContent)
		data.page.SetTitle(data.formTitle)
	}

	if data.view == VIEW_SEO {
		data.page.SetAlias(data.formAlias)
		data.page.SetCanonicalUrl(data.formCanonicalURL)
		data.page.SetMetaDescription(data.formMetaDescription)
		data.page.SetMetaKeywords(data.formMetaKeywords)
		data.page.SetMetaRobots(data.formMetaRobots)
	}

	err := controller.m.webPageUpdate(data.page)

	if err != nil {
		//config.LogStore.ErrorWithContext("At pageUpdateController > prepareDataAndValidate", err.Error())
		data.formErrorMessage = "System error. Saving page failed. " + err.Error()
		return data, ""
	}

	data.formSuccessMessage = "page saved successfully"

	return data, ""
}

func (controller pageUpdateController) prepareDataAndValidate(r *http.Request) (data pageUpdateControllerData, errorMessage string) {
	data.action = req.GetStringTrimmed(r, "action")
	data.pageID = req.GetStringTrimmed(r, "page_id")
	data.view = req.GetStringTrimmedOr(r, "view", VIEW_CONTENT)

	if data.view == "" {
		data.view = VIEW_CONTENT
	}

	if data.pageID == "" {
		return data, "page ID is required"
	}

	var err error
	data.page, err = controller.m.webPageFindByID(data.pageID)

	if err != nil {
		// config.LogStore.ErrorWithContext("At pageUpdateController > prepareDataAndValidate", err.Error())
		//return data, "page not found"
		return data, err.Error()

	}

	if data.page == nil {
		return data, "page not found"
	}

	data.formAlias = data.page.Alias()
	data.formCanonicalURL = data.page.CanonicalUrl()
	data.formContent = data.page.Content()
	data.formEditor = data.page.Editor()
	// data.formImageUrl = data.page.ImageUrl()
	// data.formFeatured = data.page.Featured()
	data.formMetaDescription = data.page.MetaDescription()
	data.formMetaKeywords = data.page.MetaKeywords()
	data.formMetaRobots = data.page.MetaRobots()
	data.formName = data.page.Name()
	// data.formMemo = data.page.Memo()
	// data.formPublishedAt = data.page.PublishedAtCarbon().ToDateTimeString()
	// data.formSummary = data.page.Summary()
	data.formStatus = data.page.Status()
	data.formTemplateID = data.page.TemplateID()
	data.formTitle = data.page.Title()

	templateList, err := controller.m.entityStore.EntityList(r.Context(), entitystore.EntityQueryOptions{
		EntityType: "template",
		Offset:     0,
		Limit:      100,
		SortBy:     "id",
		SortOrder:  "asc",
	})
	if err != nil {
		return data, "Template list failed to be retrieved" + err.Error()
	}

	data.templateList = templateList

	if r.Method != http.MethodPost {
		return data, ""
	}

	return controller.savepage(r, data)
}

type pageUpdateControllerData struct {
	action string
	pageID string
	page   types.WebPageInterface
	view   string

	templateList []entitystore.Entity

	formErrorMessage    string
	formSuccessMessage  string
	formAlias           string
	formCanonicalURL    string
	formContent         string
	formName            string
	formEditor          string
	formFeatured        string
	formImageUrl        string
	formMemo            string
	formMetaDescription string
	formMetaKeywords    string
	formMetaRobots      string
	formPublishedAt     string
	formStatus          string
	formTemplateID      string
	formSummary         string
	formTitle           string
}
