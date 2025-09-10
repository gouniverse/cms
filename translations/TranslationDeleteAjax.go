package cms

import (
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/gouniverse/utils"
)

func (m UiManager) TranslationDeleteAjax(w http.ResponseWriter, r *http.Request) {
	translationID := strings.Trim(utils.Req(r, "translation_id", ""), " ")

	if translationID == "" {
		api.Respond(w, r, api.Error("Translation ID is required"))
		return
	}

	translation, err := m.entityStore.EntityFindByID(translationID)

	if err != nil {
		api.Respond(w, r, api.Error("Database error: "+err.Error()))
		return
	}

	if translation == nil {
		api.Respond(w, r, api.Success("Translation already deleted"))
		return
	}

	isOk, err := m.entityStore.EntityDelete(translationID)

	if err != nil {
		api.Respond(w, r, api.Error("Translation failed to be deleted: "+err.Error()))
		return
	}

	if !isOk {
		api.Respond(w, r, api.Error("Translation failed to be deleted"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Translation deleted successfully", map[string]interface{}{"translation_id": translation.ID()}))
}
