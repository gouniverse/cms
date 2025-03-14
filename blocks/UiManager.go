package cms

import (
	"github.com/gouniverse/bs"
	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/hb"
)

type Config struct {
	Endpoint               string
	EntityStore            entitystore.StoreInterface
	BlockEntityType        string
	PathBlocksBlockManager string
	PathBlocksBlockUpdate  string
	WebpageComplete        func(string, string) *hb.HtmlWebpage
	FuncLayout             func(string) string
	CmsHeader              func(string) string
	CmsBreadcrumbs         func([]bs.Breadcrumb) string
}

func NewUiManager(config Config) UiManager {
	return UiManager{
		// keyEndpoint:          config.KeyEndpoint,
		endpoint:               config.Endpoint,
		entityStore:            config.EntityStore,
		blockEntityType:        config.BlockEntityType,
		pathBlocksBlockManager: config.PathBlocksBlockManager,
		pathBlocksBlockUpdate:  config.PathBlocksBlockUpdate,
		webpageComplete:        config.WebpageComplete,
		funcLayout:             config.FuncLayout,
		cmsHeader:              config.CmsHeader,
		cmsBreadcrumbs:         config.CmsBreadcrumbs,
	}
}

type UiManager struct {
	// keyEndpoint          string
	endpoint        string
	entityStore     entitystore.StoreInterface
	blockEntityType string
	// pathPagesPageManager string
	// pathPagesPageUpdate  string
	pathBlocksBlockManager string
	pathBlocksBlockUpdate  string
	webpageComplete        func(string, string) *hb.HtmlWebpage
	funcLayout             func(string) string
	cmsHeader              func(string) string
	cmsBreadcrumbs         func([]bs.Breadcrumb) string
}
