package cms

import (
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

func (m UiManager) PageTrashAjax(w http.ResponseWriter, r *http.Request) {
	pageID := strings.Trim(req.GetStringTrimmed(r, "page_id"), " ")

	if pageID == "" {
		api.Respond(w, r, api.Error("Page ID is required"))
		return
	}

	page, _ := m.entityStore.EntityFindByID(pageID)

	if page == nil {
		api.Respond(w, r, api.Error("Page NOT FOUND with ID "+pageID))
		return
	}

	isOk, err := m.entityStore.EntityTrash(pageID)

	if err != nil {
		api.Respond(w, r, api.Error("Entity failed to be trashed "+err.Error()))
		return
	}

	if !isOk {
		api.Respond(w, r, api.Error("Page failed to be moved to trash"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Page moved to trash successfully", map[string]interface{}{"page_id": page.ID()}))
}
