package cms

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

func (m UiManager) BlockDeleteAjax(w http.ResponseWriter, r *http.Request) {
	blockID := req.GetStringTrimmed(r, "block_id")

	if blockID == "" {
		api.Respond(w, r, api.Error("Block ID is required"))
		return
	}

	block, err := m.entityStore.EntityFindByID(blockID)

	if err != nil {
		api.Respond(w, r, api.Error("Database error: "+err.Error()))
		return
	}

	if block == nil {
		api.Respond(w, r, api.Success("Block already deleted"))
		return
	}

	isOk, err := m.entityStore.EntityDelete(blockID)

	if err != nil {
		api.Respond(w, r, api.Error("Block failed to be deleted: "+err.Error()))
		return
	}

	if !isOk {
		api.Respond(w, r, api.Error("Block failed to be deleted"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Block deleted successfully", map[string]interface{}{"block_id": block.ID()}))
}
