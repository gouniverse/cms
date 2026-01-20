package cms

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

// TemplateTrashAjax - moves the template to the trash
func (m UiManager) TemplateTrashAjax(w http.ResponseWriter, r *http.Request) {
	templateID := req.GetStringTrimmed(r, "template_id")

	if templateID == "" {
		api.Respond(w, r, api.Error("Template ID is required"))
		return
	}

	template, _ := m.entityStore.EntityFindByID(r.Context(), templateID)

	if template == nil {
		api.Respond(w, r, api.Error("Template NOT FOUND with ID "+templateID))
		return
	}

	isOk, err := m.entityStore.EntityTrash(r.Context(), templateID)

	if err != nil {
		api.Respond(w, r, api.Error("Template failed to be moved to trash "+err.Error()))
		return
	}

	if !isOk {
		api.Respond(w, r, api.Error("Template failed to be moved to trash"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Template moved to trash successfully", map[string]interface{}{"template_id": template.ID()}))
}
