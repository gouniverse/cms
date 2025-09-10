package cms

import (
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/gouniverse/utils"
)

func (m UiManager) TemplateCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(utils.Req(r, "name", ""), " ")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	template, err := m.entityStore.EntityCreateWithType(m.templateEntityType)

	if err != nil {
		api.Respond(w, r, api.Error("Template failed to be created: "+err.Error()))
		return
	}

	if template == nil {
		api.Respond(w, r, api.Error("Template failed to be created"))
		return
	}

	template.SetString("name", name)
	template.SetString("status", "inactive")
	m.entityStore.EntityUpdate(*template)

	api.Respond(w, r, api.SuccessWithData("Template saved successfully", map[string]interface{}{"template_id": template.ID()}))
}
