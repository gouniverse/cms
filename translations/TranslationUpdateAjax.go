package cms

import (
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

func (m UiManager) TranslationUpdateAjax(w http.ResponseWriter, r *http.Request) {
	translationID := strings.Trim(utils.Req(r, "translation_id", ""), " ")
	comment := strings.Trim(utils.Req(r, "comment", ""), " ")
	handle := strings.Trim(utils.Req(r, "handle", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")
	translationContents := map[string]string{}
	lo.ForEach(lo.Keys(m.translationLanguages), func(key string, index int) {
		translationContent := strings.Trim(utils.Req(r, "translations["+key+"]", ""), " ")
		translationContents[key] = translationContent
	})

	if translationID == "" {
		api.Respond(w, r, api.Error("Translation ID is required"))
		return
	}

	translation, err := m.entityStore.EntityFindByID(translationID)

	if err != nil {
		api.Respond(w, r, api.Error("Translation not found: "+err.Error()))
		return
	}

	if translation == nil {
		api.Respond(w, r, api.Error("Translation NOT FOUND with ID "+translationID))
		return
	}

	if status == "" {
		api.Respond(w, r, api.Error("status is required field"))
		return
	}

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	translation.SetHandle(handle)
	err = m.entityStore.EntityUpdate(*translation)

	if err != nil {
		api.Respond(w, r, api.Error("Translation failed to be updated: "+err.Error()))
		return
	}

	translation.SetString("comment", comment)
	translation.SetString("name", name)
	translation.SetAll(translationContents)
	errSetString := translation.SetString("status", status)

	if errSetString != nil {
		api.Respond(w, r, api.Error("Translation failed to be updated"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Translation saved successfully", map[string]interface{}{"translation_id": translation.ID()}))
}
