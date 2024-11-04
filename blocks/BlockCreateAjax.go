package cms

import (
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

func (m UiManager) BlockCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(utils.Req(r, "name", ""), " ")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	block, err := m.entityStore.EntityCreateWithType(m.blockEntityType)

	if err != nil {
		api.Respond(w, r, api.Error("Block failed to be created: "+err.Error()))
		return
	}

	if block == nil {
		api.Respond(w, r, api.Error("Block failed to be created"))
		return
	}

	block.SetString("name", name)

	api.Respond(w, r, api.SuccessWithData("Block saved successfully", map[string]interface{}{"block_id": block.ID()}))
}
