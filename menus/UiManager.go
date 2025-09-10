package cms

import (
	"github.com/dracory/bs"
	"github.com/dracory/entitystore"
	"github.com/dracory/hb"
)

type Config struct {
	Endpoint                     string
	EntityStore                  entitystore.StoreInterface
	MenuEntityType               string
	PageEntityType               string
	PathMenusMenuManager         string
	PathMenusMenuUpdate          string
	PathMenusMenuCreateAjax      string
	PathMenusMenuItemsUpdate     string
	PathMenusMenuItemsUpdateAjax string
	PathMenusMenuItemsFetchAjax  string
	WebpageComplete              func(string, string) *hb.HtmlWebpage
	FuncLayout                   func(string) string
	CmsHeader                    func(string) string
	CmsBreadcrumbs               func([]bs.Breadcrumb) string
}

func NewUiManager(config Config) UiManager {
	return UiManager{
		endpoint:                     config.Endpoint,
		entityStore:                  config.EntityStore,
		menuEntityType:               config.MenuEntityType,
		pageEntityType:               config.PageEntityType,
		pathMenusMenuManager:         config.PathMenusMenuManager,
		pathMenusMenuUpdate:          config.PathMenusMenuUpdate,
		pathMenusMenuCreateAjax:      config.PathMenusMenuCreateAjax,
		pathMenusMenuItemsUpdate:     config.PathMenusMenuItemsUpdate,
		pathMenusMenuItemsUpdateAjax: config.PathMenusMenuItemsUpdateAjax,
		pathMenusMenuItemsFetchAjax:  config.PathMenusMenuItemsFetchAjax,
		webpageComplete:              config.WebpageComplete,
		funcLayout:                   config.FuncLayout,
		cmsHeader:                    config.CmsHeader,
		cmsBreadcrumbs:               config.CmsBreadcrumbs,
	}
}

type UiManager struct {
	// keyEndpoint          string
	endpoint                     string
	entityStore                  entitystore.StoreInterface
	menuEntityType               string
	pageEntityType               string
	pathMenusMenuManager         string
	pathMenusMenuUpdate          string
	pathMenusMenuCreateAjax      string
	pathMenusMenuItemsUpdate     string
	pathMenusMenuItemsUpdateAjax string
	pathMenusMenuItemsFetchAjax  string
	webpageComplete              func(string, string) *hb.HtmlWebpage
	funcLayout                   func(string) string
	cmsHeader                    func(string) string
	cmsBreadcrumbs               func([]bs.Breadcrumb) string
}
