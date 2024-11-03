package cms

import (
	"net/http"
	"strings"

	"github.com/gouniverse/api"
	"github.com/gouniverse/utils"
)

func (m UiManager) PageUpdateAjax(w http.ResponseWriter, r *http.Request) {
	pageID := strings.Trim(utils.Req(r, "page_id", ""), " ")
	alias := strings.Trim(utils.Req(r, "alias", ""), " ")
	canonicalURL := strings.Trim(utils.Req(r, "canonical_url", ""), " ")
	content := strings.Trim(utils.Req(r, "content", ""), " ")
	contentEditor := strings.Trim(utils.Req(r, "content_editor", ""), " ")
	metaDescription := strings.Trim(utils.Req(r, "meta_description", ""), " ")
	metaKeywords := strings.Trim(utils.Req(r, "meta_keywords", ""), " ")
	metaRobots := strings.Trim(utils.Req(r, "meta_robots", ""), " ")
	name := strings.Trim(utils.Req(r, "name", ""), " ")
	status := strings.Trim(utils.Req(r, "status", ""), " ")
	title := strings.Trim(utils.Req(r, "title", ""), " ")
	templateID := strings.Trim(utils.Req(r, "template_id", ""), " ")
	handle := strings.Trim(utils.Req(r, "handle", ""), " ")

	if pageID == "" {
		api.Respond(w, r, api.Error("Page ID is required"))
		return
	}

	page, _ := m.entityStore.EntityFindByID(pageID)

	if page == nil {
		api.Respond(w, r, api.Error("Page NOT FOUND with ID "+pageID))
		return
	}

	if alias == "" {
		api.Respond(w, r, api.Error("alias is required field"))
		return
	}

	if name == "" {
		api.Respond(w, r, api.Error("name is required field"))
		return
	}

	if status == "" {
		api.Respond(w, r, api.Error("status is required field"))
		return
	}

	if title == "" {
		api.Respond(w, r, api.Error("title is required field"))
		return
	}

	page.SetString("alias", alias)
	page.SetString("canonical_url", canonicalURL)
	page.SetString("content", content)
	page.SetString("content_editor", contentEditor)
	page.SetString("meta_description", metaDescription)
	page.SetString("meta_keywords", metaKeywords)
	page.SetString("meta_robots", metaRobots)
	page.SetString("name", name)
	page.SetString("status", status)
	page.SetString("template_id", templateID)
	page.SetString("handle", handle)
	err := page.SetString("title", title)

	if err != nil {
		api.Respond(w, r, api.Error("Page failed to be updated: "+err.Error()))
		return
	}

	api.Respond(w, r, api.SuccessWithData("Page saved successfully", map[string]interface{}{"page_id": page.ID()}))
}
