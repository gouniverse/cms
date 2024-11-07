package cms

import (
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

func (m UiManager) TranslationCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(utils.Req(r, "name", ""), " ")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	translation, err := m.entityStore.EntityCreateWithType(m.translationEntityType)

	if err != nil {
		api.Respond(w, r, api.Error("Translation failed to be created: "+err.Error()))
		return
	}

	if translation == nil {
		api.Respond(w, r, api.Error("Translation failed to be created"))
		return
	}

	translation.SetString("name", name)

	api.Respond(w, r, api.SuccessWithData("Translation saved successfully", map[string]interface{}{"translation_id": translation.ID()}))
}
