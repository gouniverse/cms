package cms

import (
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/gouniverse/utils"
)

func (m UiManager) TranslationTrashAjax(w http.ResponseWriter, r *http.Request) {
	translationID := strings.Trim(utils.Req(r, "translation_id", ""), " ")

	if translationID == "" {
		api.Respond(w, r, api.Error("Translation ID is required"))
		return
	}

	translation, err := m.entityStore.EntityFindByID(translationID)

	if err != nil {
		api.Respond(w, r, api.Error("Error: "+err.Error()))
		return
	}

	if translation == nil {
		api.Respond(w, r, api.Success("Translation already deleted"))
		return
	}

	isOk, err := m.entityStore.EntityTrash(translationID)

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
