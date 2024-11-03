package cms

import (
	"net/http"

	cmsPages "github.com/gouniverse/cms/pages"
	// "github.com/gouniverse/cms/ve"
)

func (cms Cms) pagesUiManager(r *http.Request) cmsPages.UiManager {
	endpoint := r.Context().Value(keyEndpoint).(string)

	ui := cmsPages.NewUiManager(cmsPages.Config{
		Endpoint:             endpoint,
		EntityStore:          cms.EntityStore,
		PageEntityType:       string(ENTITY_TYPE_PAGE),
		PathPagesPageManager: string(PathPagesPageManager),
		PathPagesPageUpdate:  string(PathPagesPageUpdate),
		WebpageComplete:      WebpageComplete,
		FuncLayout:           cms.funcLayout,
		CmsHeader:            cms.cmsHeader,
		CmsBreadcrumbs:       cms.cmsBreadcrumbs,
	})

	return ui
}

func (cms Cms) pagePagesPageCreateAjax(w http.ResponseWriter, r *http.Request) {
	cms.pagesUiManager(r).PageCreateAjax(w, r)
}

func (cms Cms) pagePagesPageUpdateAjax(w http.ResponseWriter, r *http.Request) {
	cms.pagesUiManager(r).PageUpdateAjax(w, r)
}

func (cms Cms) pagePagesPageUpdate(w http.ResponseWriter, r *http.Request) {
	cms.pagesUiManager(r).PageUpdate(w, r)
}

func (cms Cms) pagePagesPageManager(w http.ResponseWriter, r *http.Request) {
	cms.pagesUiManager(r).PageManager(w, r)
}

// pagePagesPageTrashAjax - moves the template to the trash
func (cms Cms) pagePagesPageTrashAjax(w http.ResponseWriter, r *http.Request) {
	cms.pagesUiManager(r).PageTrashAjax(w, r)
}
