package cms

import (
	"net/http"
	"strings"

	"github.com/dracory/settingstore"
	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

func (m UiManager) SettingCreateAjax(w http.ResponseWriter, r *http.Request) {
	key := strings.Trim(utils.Req(r, "key", ""), " ")

	if key == "" {
		api.Respond(w, r, api.Error("Key is required field"))
		return
	}

	setting := settingstore.NewSetting().SetKey(key).SetValue("")
	err := m.settingStore.SettingCreate(r.Context(), setting)

	if err != nil {
		api.Respond(w, r, api.Error("Setting failed to be created: "+err.Error()))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Setting saved successfully", map[string]interface{}{"setting_key": key}))
}
