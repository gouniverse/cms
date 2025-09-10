package cms

import (
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/gouniverse/utils"
)

func (m UiManager) BlockUpdateAjax(w http.ResponseWriter, r *http.Request) {
	blockID := strings.Trim(utils.Req(r, "block_id", ""), " ")
	content := strings.Trim(utils.Req(r, "content", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	handle := strings.Trim(utils.Req(r, "handle", ""), " ")

	if blockID == "" {
		api.Respond(w, r, api.Error("Block ID is required"))
		return
	}

	block, err := m.entityStore.EntityFindByID(blockID)

	if err != nil {
		api.Respond(w, r, api.Error("Block not found: "+err.Error()))
		return
	}

	if block == nil {
		api.Respond(w, r, api.Error("Block NOT FOUND with ID "+blockID))
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

	block.SetString("content", content)
	block.SetString("name", name)
	block.SetString("handle", handle)
	errSetString := block.SetString("status", status)

	if errSetString != nil {
		api.Respond(w, r, api.Error("Block failed to be updated"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Block saved successfully", map[string]interface{}{"block_id": block.ID()}))
}
