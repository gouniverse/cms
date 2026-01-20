package cms

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

func (m UiManager) MenuItemsFetchAjax(w http.ResponseWriter, r *http.Request) {
	menuID := req.GetStringTrimmed(r, "menu_id")

	if menuID == "" {
		api.Respond(w, r, api.Error("Menu ID is required"))
		return
	}

	menu, _ := m.entityStore.EntityFindByID(r.Context(), menuID)

	if menu == nil {
		api.Respond(w, r, api.Error("Menu NOT FOUND with ID "+menuID))
		return
	}

	tree := m.buildTree(menuID)

	api.Respond(w, r, api.SuccessWithData("Menu items found successfully", map[string]interface{}{
		"menu_id":   menu.ID(),
		"menuitems": tree,
	}))
}
