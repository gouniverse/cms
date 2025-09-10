package cms

import (
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/gouniverse/utils"
)

func (m UiManager) MenuUpdateAjax(w http.ResponseWriter, r *http.Request) {
	menuID := strings.Trim(utils.Req(r, "menu_id", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	handle := strings.Trim(utils.Req(r, "handle", ""), " ")

	if menuID == "" {
		api.Respond(w, r, api.Error("Menu ID is required"))
		return
	}

	menu, _ := m.entityStore.EntityFindByID(menuID)

	if menu == nil {
		api.Respond(w, r, api.Error("Menu NOT FOUND with ID "+menuID))
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

	menu.SetString("name", name)
	menu.SetString("handle", handle)
	err := menu.SetString("status", status)

	if err != nil {
		api.Respond(w, r, api.Error("Menu failed to be updated "+err.Error()))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Menu saved successfully", map[string]interface{}{"menu_id": menu.ID()}))
}
