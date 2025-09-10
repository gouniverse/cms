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

func (m UiManager) BlockUpdate(w http.ResponseWriter, r *http.Request) {
	blockID := utils.Req(r, "block_id", "")
	if blockID == "" {
		api.Respond(w, r, api.Error("Block ID is required"))
		return
	}

	block, err := m.entityStore.EntityFindByID(blockID)

	if err != nil {
		api.Respond(w, r, api.Error("Block failed to be retrieved: "+err.Error()))
		return
	}

	if block == nil {
		api.Respond(w, r, api.Error("Block NOT FOUND with ID "+blockID))
		return
	}

	header := m.cmsHeader(m.endpoint)
	breadcrumbs := m.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  m.endpoint,
			Name: "Home",
		},
		{
			URL:  (m.endpoint + "?path=" + m.pathBlocksBlockManager),
			Name: "Blocks",
		},
		{
			URL:  (m.endpoint + "?path=" + m.pathBlocksBlockUpdate + "&block_id=" + blockID),
			Name: "Edit block",
		},
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "block-update")
	heading := hb.NewHeading1().HTML("Edit Block")
	button := hb.NewButton().AddChild(hb.NewHTML(icons.BootstrapCheckCircle+" ")).HTML("Save").Attr("class", "btn btn-success float-end").Attr("v-on:click", "blockSave")
	heading.AddChild(button)

	formGroupStatus := hb.NewDiv().Attr("class", "form-group")
	formGroupStatusLabel := hb.NewLabel().HTML("Status").Attr("class", "form-label")
	formGroupStatusSelect := hb.NewSelect().Attr("class", "form-select").Attr("v-model", "blockModel.status")
	formGroupOptionsActive := hb.NewOption().Attr("value", "active").HTML("Active")
	formGroupOptionsInactive := hb.NewOption().Attr("value", "inactive").HTML("Inactive")
	formGroupOptionsTrash := hb.NewOption().Attr("value", "trash").HTML("Trash")
	formGroupStatus.AddChild(formGroupStatusLabel)
	formGroupStatus.AddChild(formGroupStatusSelect.AddChild(formGroupOptionsActive).AddChild(formGroupOptionsInactive).AddChild(formGroupOptionsTrash))

	formGroupName := hb.NewDiv().Attr("class", "form-group")
	formGroupNameLabel := hb.NewLabel().HTML("Name").Attr("class", "form-label")
	formGroupNameInput := hb.NewInput().Attr("class", "form-control").Attr("v-model", "blockModel.name")
	formGroupName.AddChild(formGroupNameLabel)
	formGroupName.AddChild(formGroupNameInput)

	formGroupContent := hb.NewDiv().Attr("class", "form-group")
	formGroupContentLabel := hb.NewLabel().HTML("Content").Attr("class", "form-label")
	formGroupContentInput := hb.NewTextArea().Attr("class", "form-control CodeMirror").Attr("v-model", "blockModel.content")
	formGroupContent.AddChild(formGroupContentLabel)
	formGroupContent.AddChild(formGroupContentInput)

	paragraphUsage := hb.NewParagraph().Attr("class", "text-info mt-5").AddChild(hb.NewHTML("To use this block in your website use the following shortcode:"))

	blockName, _ := block.GetString("name", "")
	code := hb.NewCode().AddChild(hb.NewPRE().HTML(`&lt;!-- START: Block: ` + blockName + ` -->
[[BLOCK_` + block.ID() + `]]
&lt;!-- END: Block: ` + blockName + ` -->`))
	paragraphUsage.AddChild(code)

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrumbs))
	container.AddChild(formGroupStatus).AddChild(formGroupName).AddChild(formGroupContent)
	container.AddChild(paragraphUsage)

	h := container.ToHTML()

	name, err := block.GetString("name", "")
	if err != nil {
		api.Respond(w, r, api.Error("Name failed to be retrieved: "+err.Error()))
		return
	}
	statusAttribute, err := m.entityStore.AttributeFind(block.ID(), "status")

	if err != nil {
		api.Respond(w, r, api.Error("IO Error. Attribute failed to be pulled"))
		return
	}

	status := ""
	if statusAttribute != nil {
		status = statusAttribute.GetString()
	}
	contentAttribute, err := m.entityStore.AttributeFind(block.ID(), "content")

	if err != nil {
		api.Respond(w, r, api.Error("IO Error. Attribute failed to be fetched"))
		return
	}

	content := ""
	if contentAttribute != nil {
		content = contentAttribute.GetString()
	}

	contentJSON, _ := json.Marshal(content)
	nameJSON, _ := json.Marshal(name)
	statusJSON, _ := json.Marshal(status)

	inlineScript := `
var blockUpdateUrl = "` + m.endpoint + `?path=blocks/block-update-ajax";
var blockId = "` + blockID + `";
var name = ` + string(nameJSON) + `;
var status = ` + string(statusJSON) + `;
var content = ` + string(contentJSON) + `;
const BlockUpdate = {
	data() {
		return {
			blockModel:{
				blockId: blockId,
				content: content,
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
					self.blockModel.content = editor.getValue();
				});
				$(document).on('change', '.CodeMirror', function() {
					self.blockModel.content = editor.getValue();
				});
				setInterval(()=>{
					self.blockModel.content = editor.getValue();
				}, 1000)
			}
		}, 500);
	},
	methods: {
		blockSave(){
			var content = this.blockModel.content;
			var name = this.blockModel.name;
			var blockId = this.blockModel.blockId;
			var status = this.blockModel.status;
			
			$.post(blockUpdateUrl, {
				block_id:blockId,
				content: content,
				name: name,
				status: status,
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
Vue.createApp(BlockUpdate).mount('#block-update')
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

	webpage := m.webpageComplete("Edit Block", h).
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
