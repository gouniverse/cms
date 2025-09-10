package cms

import (
	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/dracory/settingstore"
)

type Config struct {
	Endpoint                   string
	SettingStore               settingstore.StoreInterface
	PathSettingsSettingManager string
	PathSettingsSettingUpdate  string
	WebpageComplete            func(string, string) *hb.HtmlWebpage
	FuncLayout                 func(string) string
	CmsHeader                  func(string) string
	CmsBreadcrumbs             func([]bs.Breadcrumb) string
}

func NewUiManager(config Config) UiManager {
	return UiManager{
		endpoint: config.Endpoint,
		// entityStore:                  config.EntityStore,
		settingStore:               config.SettingStore,
		pathSettingsSettingManager: config.PathSettingsSettingManager,
		pathSettingsSettingUpdate:  config.PathSettingsSettingUpdate,
		// pathMenusMenuManager:         config.PathMenusMenuManager,
		// pathMenusMenuUpdate:          config.PathMenusMenuUpdate,
		// pathMenusMenuCreateAjax:      config.PathMenusMenuCreateAjax,
		// pathMenusMenuItemsUpdate:     config.PathMenusMenuItemsUpdate,
		// pathMenusMenuItemsUpdateAjax: config.PathMenusMenuItemsUpdateAjax,
		// pathMenusMenuItemsFetchAjax:  config.PathMenusMenuItemsFetchAjax,
		webpageComplete: config.WebpageComplete,
		funcLayout:      config.FuncLayout,
		cmsHeader:       config.CmsHeader,
		cmsBreadcrumbs:  config.CmsBreadcrumbs,
	}
}

type UiManager struct {
	endpoint                   string
	settingStore               settingstore.StoreInterface
	pathSettingsSettingManager string
	pathSettingsSettingUpdate  string
	webpageComplete            func(string, string) *hb.HtmlWebpage
	funcLayout                 func(string) string
	cmsHeader                  func(string) string
	cmsBreadcrumbs             func([]bs.Breadcrumb) string
}
