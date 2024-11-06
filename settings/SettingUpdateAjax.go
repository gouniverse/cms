package cms

import (
	"net/http"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

func (m UiManager) SettingUpdateAjax(w http.ResponseWriter, r *http.Request) {
	settingKey := utils.Req(r, "setting_key", "")
	settingValue := utils.Req(r, "setting_value", "%%NOTSENT%%")
	if settingKey == "" {
		api.Respond(w, r, api.Error("Setting key is required"))
		return
	}
	if settingValue == "%%NOTSENT%%" {
		api.Respond(w, r, api.Error("Setting value is required"))
		return
	}

	isOk, err := m.settingStore.Set(settingKey, settingValue)

	if err != nil {
		api.Respond(w, r, api.Error("Settings failed to be updated "+err.Error()))
		return
	}

	if !isOk {
		api.Respond(w, r, api.Error("Setting failed to be updated"))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Setting saved successfully", map[string]interface{}{"setting_key": settingKey}))
}
