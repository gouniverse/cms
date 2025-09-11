package cms

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/req"
	"github.com/dracory/str"
)

func (m UiManager) PageCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := req.GetStringTrimmed(r, "name")

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	page, err := m.entityStore.EntityCreateWithType(m.pageEntityType)

	if err != nil {
		api.Respond(w, r, api.Error("Page failed to be created: "+err.Error()))
		return
	}

	if page == nil {
		api.Respond(w, r, api.Error("Page failed to be created"))
		return
	}

	page.SetString("name", name)
	page.SetString("status", "inactive")
	page.SetString("title", name)
	page.SetString("alias", "/"+str.Slugify(name+"-"+str.Random(16), '-'))

	api.Respond(w, r, api.SuccessWithData("Page saved successfully", map[string]interface{}{"page_id": page.ID()}))
}
