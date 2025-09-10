package cms

import (
	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/gouniverse/entitystore"
)

type Config struct {
	Endpoint                 string
	EntityStore              entitystore.StoreInterface
	WidgetEntityType         string
	PathWidgetsWidgetManager string
	PathWidgetsWidgetUpdate  string
	WebpageComplete          func(string, string) *hb.HtmlWebpage
	FuncLayout               func(string) string
	CmsHeader                func(string) string
	CmsBreadcrumbs           func([]bs.Breadcrumb) string
}

func NewUiManager(config Config) UiManager {
	return UiManager{
		// keyEndpoint:          config.KeyEndpoint,
		endpoint:                 config.Endpoint,
		entityStore:              config.EntityStore,
		widgetEntityType:         config.WidgetEntityType,
		pathWidgetsWidgetManager: config.PathWidgetsWidgetManager,
		pathWidgetsWidgetUpdate:  config.PathWidgetsWidgetUpdate,
		webpageComplete:          config.WebpageComplete,
		funcLayout:               config.FuncLayout,
		cmsHeader:                config.CmsHeader,
		cmsBreadcrumbs:           config.CmsBreadcrumbs,
	}
}

type UiManager struct {
	// keyEndpoint          string
	endpoint         string
	entityStore      entitystore.StoreInterface
	widgetEntityType string
	// pathPagesPageManager string
	// pathPagesPageUpdate  string
	pathWidgetsWidgetManager string
	pathWidgetsWidgetUpdate  string
	webpageComplete          func(string, string) *hb.HtmlWebpage
	funcLayout               func(string) string
	cmsHeader                func(string) string
	cmsBreadcrumbs           func([]bs.Breadcrumb) string
}
