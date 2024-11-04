package cms

import (
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

func (m UiManager) BlockTrashAjax(w http.ResponseWriter, r *http.Request) {
	blockID := strings.Trim(utils.Req(r, "block_id", ""), " ")

	if blockID == "" {
		api.Respond(w, r, api.Error("Block ID is required"))
		return
	}

	block, err := m.entityStore.EntityFindByID(blockID)

	if err != nil {
		api.Respond(w, r, api.Error("Error: "+err.Error()))
		return
	}

	if block == nil {
		api.Respond(w, r, api.Success("Block already deleted"))
		return
	}

	isOk, err := m.entityStore.EntityTrash(blockID)

	if err != nil {
		api.Respond(w, r, api.Error("Block failed to be trashed"))
		return
	}

	if !isOk {
		api.Respond(w, r, api.Error("Block failed to be trashed"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Block trashed successfully", map[string]interface{}{"block_id": block.ID()}))
}
