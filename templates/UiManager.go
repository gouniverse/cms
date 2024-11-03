package cms

import (
	"github.com/gouniverse/bs"
	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/hb"
)

type Config struct {
	Endpoint             string
	EntityStore          *entitystore.Store
	TemplateEntityType   string
	PathPagesPageManager string
	PathPagesPageUpdate  string
	WebpageComplete      func(string, string) *hb.HtmlWebpage
	FuncLayout           func(string) string
	CmsHeader            func(string) string
	CmsBreadcrumbs       func([]bs.Breadcrumb) string
}

func NewUiManager(config Config) UiManager {
	return UiManager{
		// keyEndpoint:          config.KeyEndpoint,
		endpoint:             config.Endpoint,
		entityStore:          config.EntityStore,
		templateEntityType:   config.TemplateEntityType,
		pathPagesPageManager: config.PathPagesPageManager,
		pathPagesPageUpdate:  config.PathPagesPageUpdate,
		webpageComplete:      config.WebpageComplete,
		funcLayout:           config.FuncLayout,
		cmsHeader:            config.CmsHeader,
		cmsBreadcrumbs:       config.CmsBreadcrumbs,
	}
}

type UiManager struct {
	// keyEndpoint          string
	endpoint             string
	entityStore          *entitystore.Store
	templateEntityType   string
	pathPagesPageManager string
	pathPagesPageUpdate  string
	webpageComplete      func(string, string) *hb.HtmlWebpage
	funcLayout           func(string) string
	cmsHeader            func(string) string
	cmsBreadcrumbs       func([]bs.Breadcrumb) string
}
