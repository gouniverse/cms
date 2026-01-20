package cms

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

func (m UiManager) MenuCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := req.GetStringTrimmed(r, "name")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	menu, err := m.entityStore.EntityCreateWithType(r.Context(), m.menuEntityType)

	if err != nil {
		api.Respond(w, r, api.Error("Menu failed to be created "+err.Error()))
		return
	}

	if menu == nil {
		api.Respond(w, r, api.Error("Menu failed to be created"))
		return
	}

	menu.SetString("name", name)

	api.Respond(w, r, api.SuccessWithData("Menu saved successfully", map[string]interface{}{"menu_id": menu.ID()}))
}
