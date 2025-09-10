package cms

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

func (m UiManager) MenuItemsUpdateAjax(w http.ResponseWriter, r *http.Request) {
	menuID := strings.TrimSpace(utils.Req(r, "menu_id", ""))
	data := strings.TrimSpace(utils.Req(r, "data", ""))

	if menuID == "" {
		api.Respond(w, r, api.Error("Menu ID is required"))
		return
	}

	if data == "" {
		api.Respond(w, r, api.Error("Data is required"))
		return
	}

	menu, _ := m.entityStore.EntityFindByID(menuID)

	if menu == nil {
		api.Respond(w, r, api.Error("Menu NOT FOUND with ID "+menuID))
		return
	}

	var nodes []map[string]interface{}
	err := json.Unmarshal([]byte(data), &nodes)
	if err != nil {
		log.Println(err)
		api.Respond(w, r, api.Error("Menu failed to be updated. IO error"))
		return
	}

	existingMenuItemIDs := []string{}
	flatNodeList := flattenTree(nodes)

	for _, node := range flatNodeList {
		id := node["id"].(string)
		name := node["name"].(string)
		pageID := lo.ValueOr(node, "page_id", "").(string)
		url := lo.ValueOr(node, "url", "").(string)
		target := lo.ValueOr(node, "target", "").(string)
		parentID := lo.ValueOr(node, "parent_id", "").(string)
		sequence := lo.ValueOr(node, "sequence", "").(string)

		menuitem, _ := m.entityStore.EntityFindByID(id)
		if menuitem == nil {
			menuitem, err = m.entityStore.EntityCreateWithType(m.menuEntityType)
			if err != nil {
				api.Respond(w, r, api.Error("Menu item failed to be created "+err.Error()))
				return
			}
		}
		menuitem.SetString("name", name)
		menuitem.SetString("menu_id", menuID)
		menuitem.SetString("parent_id", parentID)
		menuitem.SetString("sequence", sequence)
		menuitem.SetString("page_id", pageID)
		menuitem.SetString("url", url)
		err := menuitem.SetString("target", target)

		if err != nil {
			api.Respond(w, r, api.Error("Menu items failed to be updated "+err.Error()))
			return
		}

		existingMenuItemIDs = append(existingMenuItemIDs, menuitem.ID())
	}

	errMessage := m.cleanMenuFromNonExistingMenuItems(menuID, existingMenuItemIDs)

	if errMessage != "" {
		api.Respond(w, r, api.Error(errMessage))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Menu items saved successfully", map[string]interface{}{"menu_id": menu.ID()}))
}
