package cms

import (
	"github.com/dracory/blockeditor"
	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/gouniverse/cms/types"
	"github.com/gouniverse/entitystore"
)

type Config struct {
	BlockEditorDefinitions []blockeditor.BlockDefinition
	Endpoint               string
	EntityStore            entitystore.StoreInterface
	PageEntityType         string
	PathPagesPageManager   string
	PathPagesPageUpdate    string
	WebpageComplete        func(string, string) *hb.HtmlWebpage
	FuncLayout             func(string) string
	CmsHeader              func(string) string
	CmsBreadcrumbs         func([]bs.Breadcrumb) string
	WebPageFindByID        func(string) (types.WebPageInterface, error)
	WebPageUpdate          func(types.WebPageInterface) error
}

func NewUiManager(config Config) UiManager {
	return UiManager{
		blockEditorDefinitions: config.BlockEditorDefinitions,
		endpoint:               config.Endpoint,
		entityStore:            config.EntityStore,
		pageEntityType:         config.PageEntityType,
		pathPagesPageManager:   config.PathPagesPageManager,
		pathPagesPageUpdate:    config.PathPagesPageUpdate,
		webpageComplete:        config.WebpageComplete,
		funcLayout:             config.FuncLayout,
		cmsHeader:              config.CmsHeader,
		cmsBreadcrumbs:         config.CmsBreadcrumbs,
		webPageFindByID:        config.WebPageFindByID,
		webPageUpdate:          config.WebPageUpdate,
	}
}

type UiManager struct {
	blockEditorDefinitions []blockeditor.BlockDefinition
	endpoint               string
	entityStore            entitystore.StoreInterface
	pageEntityType         string
	pathPagesPageManager   string
	pathPagesPageUpdate    string
	webpageComplete        func(string, string) *hb.HtmlWebpage
	funcLayout             func(string) string
	cmsHeader              func(string) string
	cmsBreadcrumbs         func([]bs.Breadcrumb) string
	webPageFindByID        func(string) (types.WebPageInterface, error)
	webPageUpdate          func(types.WebPageInterface) error
}
