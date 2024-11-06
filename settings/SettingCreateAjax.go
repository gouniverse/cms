package cms

import (
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

func (m UiManager) SettingCreateAjax(w http.ResponseWriter, r *http.Request) {
	key := strings.Trim(utils.Req(r, "key", ""), " ")

	if key == "" {
		api.Respond(w, r, api.Error("Key is required field"))
		return
	}

	isOk, err := m.settingStore.Set(key, "")

	if err != nil {
		api.Respond(w, r, api.Error("Setting failed to be created: "+err.Error()))
		return
	}

	if !isOk {
		api.Respond(w, r, api.Error("Setting failed to be created"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Setting saved successfully", map[string]interface{}{"setting_key": key}))
}
