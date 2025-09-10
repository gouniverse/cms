package cms

import (
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/gouniverse/utils"
)

func (m UiManager) WidgetUpdateAjax(w http.ResponseWriter, r *http.Request) {
	widgetID := strings.Trim(utils.Req(r, "widget_id", ""), " ")
	content := strings.Trim(utils.Req(r, "content", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	handle := strings.Trim(utils.Req(r, "handle", ""), " ")

	if widgetID == "" {
		api.Respond(w, r, api.Error("Widget ID is required"))
		return
	}

	widget, _ := m.entityStore.EntityFindByID(widgetID)

	if widget == nil {
		api.Respond(w, r, api.Error("Widget NOT FOUND with ID "+widgetID))
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

	widget.SetString("content", content)
	widget.SetString("name", name)
	widget.SetString("handle", handle)
	err := widget.SetString("status", status)

	if err != nil {
		api.Respond(w, r, api.Error("Widget failed to be updated: "+err.Error()))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Widget saved successfully", map[string]interface{}{"widget_id": widget.ID()}))
}
