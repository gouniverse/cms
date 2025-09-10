package cms

import (
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/gouniverse/strutils"
	"github.com/gouniverse/utils"
)

func (m UiManager) PageCreateAjax(w http.ResponseWriter, r *http.Request) {
	name := strings.Trim(utils.Req(r, "name", ""), " ")

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
	page.SetString("alias", "/"+strutils.Slugify(name+"-"+strutils.Random(16), '-'))

	api.Respond(w, r, api.SuccessWithData("Page saved successfully", map[string]interface{}{"page_id": page.ID()}))
}
