package cms

import (
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

func (m UiManager) WidgetCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(utils.Req(r, "name", ""), " ")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	widget, err := m.entityStore.EntityCreateWithType(m.widgetEntityType)

	if err != nil {
		api.Respond(w, r, api.Error("Widget failed to be created: "+err.Error()))
		return
	}

	if widget == nil {
		api.Respond(w, r, api.Error("Widget failed to be created"))
		return
	}

	widget.SetString("name", name)

	api.Respond(w, r, api.SuccessWithData("Widget saved successfully", map[string]interface{}{"widget_id": widget.ID()}))
}
