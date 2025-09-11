package cms

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

func (m UiManager) PageUpdateAjax(w http.ResponseWriter, r *http.Request) {
	pageID := req.GetStringTrimmed(r, "page_id")
	alias := req.GetStringTrimmed(r, "alias")
	canonicalURL := req.GetStringTrimmed(r, "canonical_url")
	content := req.GetStringTrimmed(r, "content")
	contentEditor := req.GetStringTrimmed(r, "content_editor")
	metaDescription := req.GetStringTrimmed(r, "meta_description")
	metaKeywords := req.GetStringTrimmed(r, "meta_keywords")
	metaRobots := req.GetStringTrimmed(r, "meta_robots")
	name := req.GetStringTrimmed(r, "name")
	status := req.GetStringTrimmed(r, "status")
	title := req.GetStringTrimmed(r, "title")
	templateID := req.GetStringTrimmed(r, "template_id")
	handle := req.GetStringTrimmed(r, "handle")

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
