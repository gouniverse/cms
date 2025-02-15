package cms

import (
	"net/http"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

func (m UiManager) SettingDeleteAjax(w http.ResponseWriter, r *http.Request) {
	settingKey := utils.Req(r, "setting_key", "")

	if settingKey == "" {
		api.Respond(w, r, api.Error("Setting key is required"))
		return
	}

	setting, err := m.settingStore.SettingFindByKey(r.Context(), settingKey)

	if err != nil {
		api.Respond(w, r, api.Error(err.Error()))
		return
	}

	m.settingStore.SettingDelete(r.Context(), setting)

	// if isOk == false {
	// 	api.Respond(w, r, api.Error("Setting failed to be deleted"))
	// 	return
	// }

	api.Respond(w, r, api.SuccessWithData("Setting deleted successfully", map[string]interface{}{"setting_key": settingKey}))
}
