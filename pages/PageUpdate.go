package cms

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/entitystore"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/gouniverse/responses"
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

func (m UiManager) PageUpdateV1(w http.ResponseWriter, r *http.Request) {
	pageID := req.GetStringTrimmed(r, "page_id")
	if pageID == "" {
		api.Respond(w, r, api.Error("Page ID is required"))
		return
	}

	page, _ := m.entityStore.EntityFindByID(r.Context(), pageID)

	if page == nil {
		api.Respond(w, r, api.Error("Page NOT FOUND with ID "+pageID))
		return
	}

	header := m.cmsHeader(m.endpoint)
	breadcrumbs := m.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  m.endpoint,
			Name: "Home",
		},
		{
			URL:  (m.endpoint + "?path=" + m.pathPagesPageManager),
			Name: "Pages",
		},
		{
			URL:  (m.endpoint + "?path=" + m.pathPagesPageUpdate + "&page_id=" + pageID),
			Name: "Edit page",
		},
	})

	container := hb.NewDiv().ID("page-update").Class("container")
	heading := hb.NewHeading1().HTML("Edit Page")
	saveButton := hb.NewButton().HTML("Save").Class("btn btn-success float-end").Attr("v-on:click", "pageSave")
	heading.Child(saveButton)

	tabNavigation := bs.NavTabs().Style("margin-bottom: 3px;").Attr("role", "tablist")
	tabNavigationContent := bs.NavItem().Child(bs.NavLink().ID("TabContent-tab").Attr("href", "#TabContent").Attr("v-on:click", "tab('TabContent')").HTML("Content"))
	tabNavigationSeo := bs.NavItem().Child(bs.NavLink().ID("TabSeo-tab").Attr("href", "#TabSeo").Attr("v-on:click", "tab('TabSeo')").HTML("SEO"))
	tabNavigationSettings := bs.NavItem().Child(bs.NavLink().ID("TabSettings-tab").Attr("href", "#TabSettings").Attr("v-on:click", "tab('TabSettings')").HTML("Settings"))
	tabNavigation.AddChild(tabNavigationContent).AddChild(tabNavigationSeo).AddChild(tabNavigationSettings)

	tabContent := hb.NewDiv().Class("tab-content")
	tabContentContent := hb.NewDiv().ID("TabContent").Attr("class", "tab-pane fade show active").Attr("data-bs-toggle", "tab")
	tabContentSeo := hb.NewDiv().ID("TabSeo").Attr("class", "tab-pane fade").Attr("data-bs-toggle", "tab")
	tabContentSettings := hb.NewDiv().ID("TabSettings").Attr("class", "tab-pane fade").Attr("data-bs-toggle", "tab")
	tabContent.AddChild(tabContentContent).AddChild(tabContentSeo).AddChild(tabContentSettings)

	// <div class="form-group well" style="display:table;width:100%;margin-top:10px;padding:5px 10px;">
	//                     Page address: <a href="<?php echo $page->url(); ?>" target="_blank"><?php echo $page->url(); ?></a> &nbsp;&nbsp;&nbsp; (to change click on Settings tab)
	// 				</div>

	// Status
	formGroupStatus := bs.FormGroup().Children([]hb.TagInterface{
		bs.FormLabel("Status"),
		bs.FormSelect().Attr("v-model", "pageModel.status").Children([]hb.TagInterface{
			bs.FormSelectOption("active", "Active"),
			bs.FormSelectOption("inactive", "Inactive"),
			bs.FormSelectOption("trash", "Trash"),
		}),
	})

	// Content Editor
	formGroupEditor := bs.FormGroup().Children([]hb.TagInterface{
		bs.FormLabel("Content Editor"),
		bs.FormSelect().Attr("v-model", "pageModel.contentEditor").Children([]hb.TagInterface{
			bs.FormSelectOption("", "- none -"),
			bs.FormSelectOption("codemirror", "CodeMirror"),
			// bs.FormSelectOption("visual", "Visual Editor (Experimental)"),
		}),
		bs.FormText("The content editor allows you to select the mode for editing the content. Note you will need to save and refresh to activate"),
	})

	// Name
	formGroupName := bs.FormGroup().Children([]hb.TagInterface{
		bs.FormLabel("Name"),
		bs.FormInput().Attr("v-model", "pageModel.name"),
		bs.FormText("The name of the page as visible in the admin panel. This is not vsible to the page vistors"),
	})

	// Alias
	formGroupAlias := bs.FormGroup().Children([]hb.TagInterface{
		bs.FormLabel("Alias / Path"),
		bs.FormInput().Attr("v-model", "pageModel.alias"),
		bs.FormText("The relative path on the website where this page will be visible to the vistors"),
	})

	// Canonical Url
	formGroupCanonicalURL := bs.FormGroup().Children([]hb.TagInterface{
		bs.FormLabel("Canonical Url"),
		bs.FormInput().Attr("v-model", "pageModel.canonicalUrl"),
	})

	// Meta Description
	formGroupMetaDescription := bs.FormGroup().Children([]hb.TagInterface{
		bs.FormLabel("Meta Description"),
		bs.FormInput().Attr("v-model", "pageModel.metaDescription"),
	})

	// Meta Keywords
	formGroupMetaKeywords := hb.NewDiv().Class("form-group").Children([]hb.TagInterface{
		bs.FormLabel("Meta Keywords"),
		bs.FormInput().Attr("v-model", "pageModel.metaKeywords"),
	})

	// Robots
	formGroupMetaRobots := hb.NewDiv().Class("form-group").Children([]hb.TagInterface{
		bs.FormLabel("Meta Robots"),
		bs.FormInput().Attr("v-model", "pageModel.metaRobots"),
	})

	// Template
	templateList, err := m.entityStore.EntityList(r.Context(), entitystore.EntityQueryOptions{
		EntityType: "template",
		Offset:     0,
		Limit:      100,
		SortBy:     "id",
		SortOrder:  "asc",
	})
	if err != nil {
		api.Respond(w, r, api.Error("Entity list failed to be retrieved "+err.Error()))
		return
	}
	formGroupTemplateSelect := bs.FormSelect().Attr("v-model", "pageModel.templateId")
	formGroupTemplateOptionsEmpty := bs.FormSelectOption("", "- none -")
	formGroupTemplateSelect.Child(formGroupTemplateOptionsEmpty)
	for _, template := range templateList {
		templateName, _ := template.GetString("name", "n/a")
		formGroupTemplateOptionsTemplate := bs.FormSelectOption(template.ID(), templateName)
		formGroupTemplateSelect.Child(formGroupTemplateOptionsTemplate)
	}
	formGroupTemplate := bs.FormGroup().Children([]hb.TagInterface{
		bs.FormLabel("Template").Class("form-label"),
		formGroupTemplateSelect,
		bs.FormText("Select the template that this page content will be displayed in. This feature is useful if you want to implement consistent layouts. Leaving the template field empty will display page as it is standalone"),
	})

	// Title
	formGroupTitle := hb.NewDiv().Class("form-group")
	formGroupTitleLabel := hb.NewLabel().HTML("Title").Class("form-label")
	formGroupTitleInput := hb.NewInput().Class("form-control").Attr("v-model", "pageModel.title")
	formGroupTitle.Child(formGroupTitleLabel).Child(formGroupTitleInput)

	// Content
	// editor, _ := page.GetString("content_editor", "")
	formGroupContent := bs.FormGroup()
	formGroupContentLabel := bs.FormLabel("Content")
	formGroupContentInput := hb.NewTextArea().Class("form-control CodeMirror").Attr("v-model", "pageModel.content")
	// formGroupContentInputVisualData := hb.NewTextArea().Class("form-control").Attr("v-model", "pageModel.contentVisual")
	// if editor == "visualeditor" {
	// formGroupContentInput.Style("display:none")
	// } else {
	// formGroupContentInputVisualData.Style("display:none")
	// }
	formGroupContent.Children([]hb.TagInterface{
		formGroupContentLabel,
		formGroupContentInput,
		// formGroupContentInputVisualData,
		// hb.NewHTML(ve.VisualeditorContent()),
	})

	tabContentContent.Children([]hb.TagInterface{
		formGroupTitle,
		formGroupContent,
	})
	tabContentSeo.Children([]hb.TagInterface{
		formGroupAlias,
		formGroupMetaDescription,
		formGroupMetaKeywords,
		formGroupMetaRobots,
		formGroupCanonicalURL,
	})
	tabContentSettings.Children([]hb.TagInterface{
		formGroupStatus,
		formGroupTemplate,
		formGroupName,
		formGroupEditor,
	})

	container.Children([]hb.TagInterface{
		hb.NewHTML(header),
		heading,
		hb.NewHTML(breadcrumbs),
		tabNavigation,
		tabContent,
	})

	h := container.ToHTML()

	alias, _ := page.GetString("alias", "")
	content, _ := page.GetString("content", "")
	contentEditor, _ := page.GetString("content_editor", "")
	name, _ := page.GetString("name", "")
	status, _ := page.GetString("status", "")
	templateID, _ := page.GetString("template_id", "")
	title, _ := page.GetString("title", "")
	metaDescription, _ := page.GetString("meta_description", "")
	metaKeywords, _ := page.GetString("meta_keywords", "")
	metaRobots, _ := page.GetString("meta_robots", "")
	canonicalURL, _ := page.GetString("canonical_url", "")

	canonicalURLJSON, _ := json.Marshal(canonicalURL)
	contentJSON, _ := json.Marshal(content)
	contentEditorJSON, _ := json.Marshal(contentEditor)
	nameJSON, _ := json.Marshal(name)
	templateIDJSON, _ := json.Marshal(templateID)
	titleJSON, _ := json.Marshal(title)

	inlineScript := `
var pageUpdateUrl = "` + m.endpoint + `?path=pages/page-update-ajax";
var pageId = "` + pageID + `";
var alias = "` + alias + `";
var canonicalUrl = "` + canonicalURL + `";
var metaDescription = "` + metaDescription + `";
var metaKeywords = "` + metaKeywords + `";
var metaRobots = "` + metaRobots + `";
var name = ` + string(nameJSON) + `;
var status = "` + status + `";
var title = ` + string(titleJSON) + `;
var canonicalUrl = ` + string(canonicalURLJSON) + `;
var content = ` + string(contentJSON) + `;
var contentEditor = ` + string(contentEditorJSON) + `;
var templateId = ` + string(templateIDJSON) + `;
const PageUpdate = {
	data() {
		return {
			pageModel:{
				pageId: pageId,
				alias: alias,
				canonicalUrl:canonicalUrl,
				content: content,
				contentEditor: contentEditor,
				metaDescription:metaDescription,
				metaKeywords:metaKeywords,
				metaRobots:metaRobots,
				name: name,
				status: status,
				title: title,
				templateId: templateId,
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
					self.pageModel.content = editor.getValue();
				});
				$(document).on('change', '.CodeMirror', function() {
					self.pageModel.content = editor.getValue();
				});
				setInterval(()=>{
					self.pageModel.content = editor.getValue();
				}, 1000)
			}
		}, 500);
	},
	methods: {
		tab(id){
			$(".nav-link").removeClass("show active")
			$(".tab-pane").removeClass("show active")
			$("#"+id).addClass("show active")
			$("#"+id+"-tab").addClass("active")
		},
		pageSave(){
			var alias = this.pageModel.alias;
			var canonicalUrl = this.pageModel.canonicalUrl;
			var content = this.pageModel.content;
			var contentEditor = this.pageModel.contentEditor;
			var metaDescription = this.pageModel.metaDescription;
			var metaKeywords = this.pageModel.metaKeywords;
			var metaRobots = this.pageModel.metaRobots;
			var name = this.pageModel.name;
			var pageId = this.pageModel.pageId;
			var status = this.pageModel.status;
			var templateId = this.pageModel.templateId;
			var title = this.pageModel.title;
			
			$.post(pageUpdateUrl, {
				page_id:pageId,
				alias: alias,
				content: content,
				content_editor: contentEditor,
				canonical_url:canonicalUrl,
				meta_description:metaDescription,
				meta_keywords:metaKeywords,
				meta_robots:metaRobots,
				name: name,
				status: status,
				title: title,
				template_id: templateId,
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
					title: 'Page saved',
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
Vue.createApp(PageUpdate).mount('#page-update')
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

	webpage := m.webpageComplete("Edit Page", h)
	webpage.AddStyleURLs([]string{
		codemirrorCss,
	})
	webpage.AddScriptURLs([]string{
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
	webpage.AddScript(inlineScript)
	// webpage.AddScript(ve.VisualeditorScripts())

	responses.HTMLResponse(w, r, m.funcLayout(webpage.ToHTML()))
}
