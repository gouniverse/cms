package cms

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

func (m UiManager) BlockUpdateAjax(w http.ResponseWriter, r *http.Request) {
	blockID := req.GetStringTrimmed(r, "block_id")
	content := req.GetStringTrimmed(r, "content")
	status := req.GetStringTrimmed(r, "status")
	name := req.GetStringTrimmed(r, "name")
	handle := req.GetStringTrimmed(r, "handle")

	if blockID == "" {
		api.Respond(w, r, api.Error("Block ID is required"))
		return
	}

	block, err := m.entityStore.EntityFindByID(r.Context(), blockID)

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
