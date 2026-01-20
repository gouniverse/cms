package cms

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

// pageTemplatesTemplateUpdateAjax - saves the template via Ajax
func (m UiManager) TemplateUpdateAjax(w http.ResponseWriter, r *http.Request) {
	templateID := req.GetStringTrimmed(r, "template_id")
	content := req.GetStringTrimmed(r, "content")
	name := req.GetStringTrimmed(r, "name")
	status := req.GetStringTrimmed(r, "status")
	handle := req.GetStringTrimmed(r, "handle")

	if templateID == "" {
		api.Respond(w, r, api.Error("Template ID is required"))
		return
	}

	template, _ := m.entityStore.EntityFindByID(r.Context(), templateID)

	if template == nil {
		api.Respond(w, r, api.Error("Template NOT FOUND with ID "+templateID))
		return
	}

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	if status == "" {
		api.Respond(w, r, api.Error("status is required field"))
		return
	}

	err := m.entityStore.AttributeSetString(r.Context(), template.ID(), "content", content)
	if err != nil {
		api.Respond(w, r, api.Error("Content failed to be updated: "+err.Error()))
		return
	}

	err = m.entityStore.AttributeSetString(r.Context(), template.ID(), "name", name)
	if err != nil {
		api.Respond(w, r, api.Error("Name failed to be updated: "+err.Error()))
		return
	}

	err = m.entityStore.AttributeSetString(r.Context(), template.ID(), "status", status)
	if err != nil {
		api.Respond(w, r, api.Error("Status failed to be updated: "+err.Error()))
		return
	}

	template.SetHandle(handle)
	errUpdate := m.entityStore.EntityUpdate(r.Context(), *template)

	if errUpdate != nil {
		api.Respond(w, r, api.Error("Template failed to be updated: "+errUpdate.Error()))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Template saved successfully", map[string]interface{}{"template_id": template.ID()}))
}
