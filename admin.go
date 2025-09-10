package cms

import (
	"context"
	"maps"
	"net/http"
	"strconv"

	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/entitystore"
	"github.com/dracory/hb"
	"github.com/gouniverse/responses"
	"github.com/gouniverse/utils"

	cmsBlocks "github.com/gouniverse/cms/blocks"
	cmsMenus "github.com/gouniverse/cms/menus"
	cmsPages "github.com/gouniverse/cms/pages"
	cmsSettings "github.com/gouniverse/cms/settings"
	cmsTemplates "github.com/gouniverse/cms/templates"
	cmsTranslations "github.com/gouniverse/cms/translations"
	cmsWidgets "github.com/gouniverse/cms/widgets"
)

func (cms Cms) blockUiManager(r *http.Request) cmsBlocks.UiManager {
	endpoint := r.Context().Value(keyEndpoint).(string)

	ui := cmsBlocks.NewUiManager(cmsBlocks.Config{
		Endpoint:               endpoint,
		EntityStore:            cms.EntityStore,
		BlockEntityType:        string(ENTITY_TYPE_BLOCK),
		PathBlocksBlockManager: string(PathBlocksBlockManager),
		PathBlocksBlockUpdate:  string(PathBlocksBlockUpdate),
		WebpageComplete:        WebpageComplete,
		FuncLayout:             cms.funcLayout,
		CmsHeader:              cms.cmsHeader,
		CmsBreadcrumbs:         cms.cmsBreadcrumbs,
	})

	return ui
}

func (cms Cms) menuUiManager(r *http.Request) cmsMenus.UiManager {
	endpoint := r.Context().Value(keyEndpoint).(string)

	ui := cmsMenus.NewUiManager(cmsMenus.Config{
		Endpoint:                     endpoint,
		EntityStore:                  cms.EntityStore,
		MenuEntityType:               string(ENTITY_TYPE_MENU),
		PathMenusMenuManager:         string(PathMenusMenuManager),
		PathMenusMenuUpdate:          string(PathMenusMenuUpdate),
		PathMenusMenuCreateAjax:      string(PathMenusMenuCreateAjax),
		PathMenusMenuItemsUpdate:     string(PathMenusMenuItemsUpdate),
		PathMenusMenuItemsUpdateAjax: string(PathMenusMenuItemsUpdateAjax),
		PathMenusMenuItemsFetchAjax:  string(PathMenusMenuItemsFetchAjax),
		WebpageComplete:              WebpageComplete,
		FuncLayout:                   cms.funcLayout,
		CmsHeader:                    cms.cmsHeader,
		CmsBreadcrumbs:               cms.cmsBreadcrumbs,
	})

	return ui
}

func (cms Cms) pagesUiManager(r *http.Request) cmsPages.UiManager {
	endpoint := r.Context().Value(keyEndpoint).(string)

	ui := cmsPages.NewUiManager(cmsPages.Config{
		BlockEditorDefinitions: cms.blockEditorDefinitions,
		Endpoint:               endpoint,
		EntityStore:            cms.EntityStore,
		PageEntityType:         string(ENTITY_TYPE_PAGE),
		PathPagesPageManager:   string(PathPagesPageManager),
		PathPagesPageUpdate:    string(PathPagesPageUpdate),
		WebpageComplete:        WebpageComplete,
		FuncLayout:             cms.funcLayout,
		CmsHeader:              cms.cmsHeader,
		CmsBreadcrumbs:         cms.cmsBreadcrumbs,
		WebPageFindByID:        cms.WebPageFindByID,
		WebPageUpdate:          cms.WebPageUpdate,
	})

	return ui
}

func (cms Cms) settingsUiManager(r *http.Request) cmsSettings.UiManager {
	endpoint := r.Context().Value(keyEndpoint).(string)

	ui := cmsSettings.NewUiManager(cmsSettings.Config{
		Endpoint:                   endpoint,
		SettingStore:               cms.SettingStore,
		PathSettingsSettingManager: string(PathSettingsSettingManager),
		PathSettingsSettingUpdate:  string(PathSettingsSettingUpdate),
		WebpageComplete:            WebpageComplete,
		FuncLayout:                 cms.funcLayout,
		CmsHeader:                  cms.cmsHeader,
		CmsBreadcrumbs:             cms.cmsBreadcrumbs,
	})

	return ui
}
func (cms Cms) templatesUiManager(r *http.Request) cmsTemplates.UiManager {
	endpoint := r.Context().Value(keyEndpoint).(string)

	ui := cmsTemplates.NewUiManager(cmsTemplates.Config{
		Endpoint:                     endpoint,
		EntityStore:                  cms.EntityStore,
		TemplateEntityType:           string(ENTITY_TYPE_TEMPLATE),
		PathTemplatesTemplateManager: string(PathTemplatesTemplateManager),
		PathTemplatesTemplateUpdate:  string(PathTemplatesTemplateUpdate),
		WebpageComplete:              WebpageComplete,
		FuncLayout:                   cms.funcLayout,
		CmsHeader:                    cms.cmsHeader,
		CmsBreadcrumbs:               cms.cmsBreadcrumbs,
	})

	return ui
}

func (cms Cms) translationsUiManager(r *http.Request) cmsTranslations.UiManager {
	endpoint := r.Context().Value(keyEndpoint).(string)

	ui := cmsTranslations.NewUiManager(cmsTranslations.Config{
		Endpoint:                           endpoint,
		EntityStore:                        cms.EntityStore,
		TranslationEntityType:              string(ENTITY_TYPE_TRANSLATION),
		PathTranslationsTranslationManager: string(PathTranslationsTranslationManager),
		PathTranslationsTranslationUpdate:  string(PathTranslationsTranslationUpdate),
		TranslationLanguageDefault:         cms.translationLanguageDefault,
		TranslationLanguages:               cms.translationLanguages,
		WebpageComplete:                    WebpageComplete,
		FuncLayout:                         cms.funcLayout,
		CmsHeader:                          cms.cmsHeader,
		CmsBreadcrumbs:                     cms.cmsBreadcrumbs,
	})

	return ui
}

func (cms Cms) widgetsUiManager(r *http.Request) cmsWidgets.UiManager {
	endpoint := r.Context().Value(keyEndpoint).(string)

	ui := cmsWidgets.NewUiManager(cmsWidgets.Config{
		Endpoint:                 endpoint,
		EntityStore:              cms.EntityStore,
		WidgetEntityType:         string(ENTITY_TYPE_WIDGET),
		PathWidgetsWidgetManager: string(PathWidgetsWidgetManager),
		PathWidgetsWidgetUpdate:  string(PathWidgetsWidgetUpdate),
		WebpageComplete:          WebpageComplete,
		FuncLayout:               cms.funcLayout,
		CmsHeader:                cms.cmsHeader,
		CmsBreadcrumbs:           cms.cmsBreadcrumbs,
	})

	return ui
}

// Router shows the admin page
func (cms Cms) Router(w http.ResponseWriter, r *http.Request) {
	path := utils.Req(r, "path", "home")

	if path == "" {
		path = PathHome
	}

	ctx := context.WithValue(r.Context(), keyEndpoint, r.URL.Path)

	routeFunc := cms.getRoute(path)
	routeFunc(w, r.WithContext(ctx))
}

func (cms Cms) menuRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	menuRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		PathMenusMenuCreateAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.menuUiManager(r).MenuCreateAjax(w, r)
		},
		PathMenusMenuItemsFetchAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.menuUiManager(r).MenuItemsFetchAjax(w, r)
		},
		PathMenusMenuItemsUpdate: func(w http.ResponseWriter, r *http.Request) {
			cms.menuUiManager(r).MenuItemsUpdate(w, r)
		},
		PathMenusMenuItemsUpdateAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.menuUiManager(r).MenuItemsUpdateAjax(w, r)
		},
		PathMenusMenuManager: func(w http.ResponseWriter, r *http.Request) {
			cms.menuUiManager(r).MenuManager(w, r)
		},
		PathMenusMenuUpdate: func(w http.ResponseWriter, r *http.Request) {
			cms.menuUiManager(r).MenuUpdate(w, r)
		},
		PathMenusMenuUpdateAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.menuUiManager(r).MenuUpdateAjax(w, r)
		},
	}

	return menuRoutes
}

func (cms Cms) pageRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	pageRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		PathPagesPageCreateAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.pagesUiManager(r).PageCreateAjax(w, r)
		},
		PathPagesPageManager: func(w http.ResponseWriter, r *http.Request) {
			cms.pagesUiManager(r).PageManager(w, r)
		},
		PathPagesPageTrashAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.pagesUiManager(r).PageTrashAjax(w, r)
		},
		PathPagesPageUpdate: func(w http.ResponseWriter, r *http.Request) {
			cms.pagesUiManager(r).PageUpdate(w, r)
		},
		PathPagesPageUpdateAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.pagesUiManager(r).PageUpdateAjax(w, r)
		},
	}

	return pageRoutes
}

func (cms Cms) settingRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	settingRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		PathSettingsSettingCreateAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.settingsUiManager(r).SettingCreateAjax(w, r)
		},
		PathSettingsSettingDeleteAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.settingsUiManager(r).SettingDeleteAjax(w, r)
		},
		PathSettingsSettingManager: func(w http.ResponseWriter, r *http.Request) {
			cms.settingsUiManager(r).SettingManager(w, r)
		},
		PathSettingsSettingUpdate: func(w http.ResponseWriter, r *http.Request) {
			cms.settingsUiManager(r).SettingUpdate(w, r)
		},
		PathSettingsSettingUpdateAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.settingsUiManager(r).SettingUpdateAjax(w, r)
		},
	}

	return settingRoutes
}

func (cms Cms) templateRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	templateRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		PathTemplatesTemplateCreateAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.templatesUiManager(r).TemplateCreateAjax(w, r)
		},
		PathTemplatesTemplateManager: func(w http.ResponseWriter, r *http.Request) {
			cms.templatesUiManager(r).TemplateManager(w, r)
		},
		PathTemplatesTemplateTrashAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.templatesUiManager(r).TemplateTrashAjax(w, r)
		},
		PathTemplatesTemplateUpdate: func(w http.ResponseWriter, r *http.Request) {
			cms.templatesUiManager(r).TemplateUpdate(w, r)
		},
		PathTemplatesTemplateUpdateAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.templatesUiManager(r).TemplateUpdateAjax(w, r)
		},
	}

	return templateRoutes
}

func (cms Cms) translationRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	translationRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		PathTranslationsTranslationCreateAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.translationsUiManager(r).TranslationCreateAjax(w, r)
		},
		PathTranslationsTranslationDeleteAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.translationsUiManager(r).TranslationDeleteAjax(w, r)
		},
		PathTranslationsTranslationManager: func(w http.ResponseWriter, r *http.Request) {
			cms.translationsUiManager(r).TranslationManager(w, r)
		},
		PathTranslationsTranslationTrashAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.translationsUiManager(r).TranslationTrashAjax(w, r)
		},
		PathTranslationsTranslationUpdate: func(w http.ResponseWriter, r *http.Request) {
			cms.translationsUiManager(r).TranslationUpdate(w, r)
		},
		PathTranslationsTranslationUpdateAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.translationsUiManager(r).TranslationUpdateAjax(w, r)
		},
	}

	return translationRoutes
}

func (cms Cms) widgetRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	widgetRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		PathWidgetsWidgetCreateAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.widgetsUiManager(r).WidgetCreateAjax(w, r)
		},
		PathWidgetsWidgetManager: func(w http.ResponseWriter, r *http.Request) {
			cms.widgetsUiManager(r).WidgetManager(w, r)
		},
		// PathWidgetsWidgetTrashAjax: func(w http.ResponseWriter, r *http.Request) {
		// 	cms.widgetsUiManager(r).WidgetTrashAjax(w, r)
		// },
		PathWidgetsWidgetUpdate: func(w http.ResponseWriter, r *http.Request) {
			cms.widgetsUiManager(r).WidgetUpdate(w, r)
		},
		PathWidgetsWidgetUpdateAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.widgetsUiManager(r).WidgetUpdateAjax(w, r)
		},
	}

	return widgetRoutes
}

func (cms Cms) blocksRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {

	blockRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		PathBlocksBlockCreateAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.blockUiManager(r).BlockCreateAjax(w, r)
		},
		PathBlocksBlockDeleteAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.blockUiManager(r).BlockDeleteAjax(w, r)
		},
		PathBlocksBlockManager: func(w http.ResponseWriter, r *http.Request) {
			cms.blockUiManager(r).BlockManager(w, r)
		},
		PathBlocksBlockUpdate: func(w http.ResponseWriter, r *http.Request) {
			cms.blockUiManager(r).BlockUpdate(w, r)
		},
		PathBlocksBlockTrashAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.blockUiManager(r).BlockTrashAjax(w, r)
		},
		PathBlocksBlockUpdateAjax: func(w http.ResponseWriter, r *http.Request) {
			cms.blockUiManager(r).BlockUpdateAjax(w, r)
		},
	}

	return blockRoutes

}

func (cms Cms) getRoute(route string) func(w http.ResponseWriter, r *http.Request) {

	routes := map[string]func(w http.ResponseWriter, r *http.Request){
		PathHome: cms.pageHome,

		// START: Users
		PathUsersUserCreateAjax: cms.pageUsersUserCreateAjax,
		PathUsersUserTrashAjax:  cms.pageUsersUserTrashAjax,
		PathUsersUserManager:    cms.pageUsersUserManager,
		PathUsersUserUpdate:     cms.pageUsersUserUpdate,
		PathUsersUserUpdateAjax: cms.pageUsersUserUpdateAjax,
		// END: Users

		// START: Websites
		// PathWebsitesWebsiteManager: cms.pageWebsitesWebsiteManager,
		// END: Websites

		// START: Custom Entities
		PathEntitiesEntityCreateAjax: cms.pageEntitiesEntityCreateAjax,
		PathEntitiesEntityManager:    cms.pageEntitiesEntityManager,
		PathEntitiesEntityUpdate:     cms.pageEntitiesEntityUpdate,
		PathEntitiesEntityUpdateAjax: cms.pageEntitiesEntityUpdateAjax,
		// END: Custom Entities

	}

	maps.Copy(routes, cms.blocksRoutes())
	maps.Copy(routes, cms.menuRoutes())
	maps.Copy(routes, cms.pageRoutes())
	maps.Copy(routes, cms.settingRoutes())
	maps.Copy(routes, cms.templateRoutes())
	maps.Copy(routes, cms.translationRoutes())
	maps.Copy(routes, cms.widgetRoutes())

	// log.Println(route)
	if val, ok := routes[route]; ok {
		return val
	}

	return routes[PathHome]
}

func (cms Cms) pageHome(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	// log.Println(endpoint)

	header := cms.cmsHeader(endpoint)
	breadcrumbs := cms.cmsBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  endpoint,
			Name: "Home",
		},
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "page-manager")
	heading := hb.NewHeading1().HTML("CMS Dashboard")

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrumbs))

	h := container.ToHTML()

	if cms.funcLayout("") != "" {
		responses.HTMLResponse(w, r, cms.funcLayout(h))
		return
	}

	webpage := WebpageComplete("Home", h).ToHTML()
	responses.HTMLResponse(w, r, webpage)
}

func (cms Cms) cmsBreadcrumbs(breadcrumbs []bs.Breadcrumb) string {
	return bs.Breadcrumbs(breadcrumbs).
		Style("margin-bottom:10px;").
		ToHTML()
	// nav := hb.NewNav().Attr("aria-label", "breadcrumb")
	// ol := hb.NewOL().Attr("class", "breadcrumb")
	// for k, v := range breadcrumbs {
	// 	li := hb.NewLI().Attr("class", "breadcrumb-item")
	// 	link := hb.NewHyperlink().HTML(v).Attr("href", k)

	// 	li.AddChild(link)

	// 	ol.AddChild(li)
	// }
	// nav.AddChild(ol)
	// return nav.ToHTML()
}

func (cms Cms) cmsHeader(endpoint string) string {
	linkHome := hb.NewHyperlink().
		HTML("Dashboard").
		Href(endpoint + "").
		Class("nav-link")
	linkBlocks := hb.NewHyperlink().
		HTML("Blocks ").
		Href(endpoint + "?path=" + PathBlocksBlockManager).
		Class("nav-link")
	linkMenus := hb.NewHyperlink().
		HTML("Menus ").
		Href(endpoint + "?path=" + PathMenusMenuManager).
		Class("nav-link")
	linkPages := hb.NewHyperlink().
		HTML("Pages ").
		Href(endpoint + "?path=" + PathPagesPageManager).
		Class("nav-link")
	linkTemplates := hb.NewHyperlink().
		HTML("Templates ").
		Href(endpoint + "?path=" + PathTemplatesTemplateManager).
		Class("nav-link")
	linkWidgets := hb.NewHyperlink().
		HTML("Widgets ").
		Href(endpoint + "?path=" + PathWidgetsWidgetManager).
		Class("nav-link")
	linkSettings := hb.NewHyperlink().
		HTML("Settings").
		Href(endpoint + "?path=" + PathSettingsSettingManager).
		Class("nav-link")
	linkTranslations := hb.NewHyperlink().
		HTML("Translations").
		Href(endpoint + "?path=" + PathTranslationsTranslationManager).
		Class("nav-link")

	blocksCount, _ := cms.EntityStore.EntityCount(entitystore.EntityQueryOptions{
		EntityType: ENTITY_TYPE_BLOCK,
	})
	menusCount, _ := cms.EntityStore.EntityCount(entitystore.EntityQueryOptions{
		EntityType: ENTITY_TYPE_MENU,
	})
	pagesCount, _ := cms.EntityStore.EntityCount(entitystore.EntityQueryOptions{
		EntityType: ENTITY_TYPE_PAGE,
	})
	templatesCount, _ := cms.EntityStore.EntityCount(entitystore.EntityQueryOptions{
		EntityType: ENTITY_TYPE_TEMPLATE,
	})
	translationsCount, _ := cms.EntityStore.EntityCount(entitystore.EntityQueryOptions{
		EntityType: ENTITY_TYPE_TRANSLATION,
	})
	widgetsCount, _ := cms.EntityStore.EntityCount(entitystore.EntityQueryOptions{
		EntityType: ENTITY_TYPE_WIDGET,
	})

	ulNav := hb.NewUL().Class("nav  nav-pills justify-content-center")
	ulNav.AddChild(hb.NewLI().Class("nav-item").Child(linkHome))

	if cms.templatesEnabled {
		ulNav.AddChild(hb.NewLI().Class("nav-item").AddChild(linkTemplates.AddChild(hb.NewSpan().Class("badge bg-secondary").HTML(strconv.FormatInt(templatesCount, 10)))))
	}

	if cms.pagesEnabled {
		ulNav.AddChild(hb.NewLI().Class("nav-item").AddChild(linkPages.AddChild(hb.NewSpan().Class("badge bg-secondary").HTML(strconv.FormatInt(pagesCount, 10)))))
	}

	if cms.menusEnabled {
		ulNav.AddChild(hb.NewLI().Class("nav-item").AddChild(linkMenus.AddChild(hb.NewSpan().Class("badge bg-secondary").HTML(strconv.FormatInt(menusCount, 10)))))
	}

	if cms.blocksEnabled {
		ulNav.AddChild(hb.NewLI().Class("nav-item").AddChild(linkBlocks.AddChild(hb.NewSpan().Class("badge bg-secondary").HTML(strconv.FormatInt(blocksCount, 10)))))
	}

	if cms.widgetsEnabled {
		ulNav.AddChild(hb.NewLI().Class("nav-item").AddChild(linkWidgets.AddChild(hb.NewSpan().Class("badge bg-secondary").HTML(strconv.FormatInt(widgetsCount, 10)))))
	}

	if cms.translationsEnabled {
		ulNav.AddChild(hb.NewLI().Class("nav-item").Child(linkTranslations.Child(hb.NewSpan().Class("badge bg-secondary").HTML(utils.ToString(translationsCount)))))
	}

	if cms.settingsEnabled {
		ulNav.AddChild(hb.NewLI().Class("nav-item").AddChild(linkSettings))
	}
	// add Translations

	for _, entity := range cms.customEntityList {
		linkEntity := hb.NewHyperlink().HTML(entity.TypeLabel).Href(endpoint + "?path=entities/entity-manager&type=" + entity.Type).Class("nav-link")
		ulNav.AddChild(hb.NewLI().Class("nav-item").Child(linkEntity))
	}

	divCard := hb.NewDiv().Class("card card-default mt-3 mb-3")
	divCardBody := hb.NewDiv().Class("card-body").Style("padding: 2px;")
	return divCard.AddChild(divCardBody.AddChild(ulNav)).ToHTML()
}

// WebpageComplete returns the webpage template for the website
func WebpageComplete(title, content string) *hb.HtmlWebpage {
	faviconImgCms := `data:image/x-icon;base64,AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAmzKzAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABEQEAAQERAAEAAQABAAEAAQABAQEBEQABAAEREQEAAAERARARAREAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD//wAA//8AAP//AAD//wAA//8AAP//AAD//wAAi6MAALu7AAC6owAAuC8AAIkjAAD//wAA//8AAP//AAD//wAA`
	app := ""
	webpage := hb.NewWebpage()
	webpage.SetTitle(title)
	webpage.SetFavicon(faviconImgCms)

	webpage.AddStyleURLs([]string{
		cdn.BootstrapCss_5_2_3(),
	})
	webpage.AddScriptURLs([]string{
		cdn.BootstrapJs_5_2_3(),
		cdn.Jquery_3_6_4(),
		cdn.VueJs_3(),
		cdn.Sweetalert2_10(),
	})
	webpage.AddScripts([]string{
		app,
	})
	webpage.AddStyle(`html,body{height:100%;font-family: Ubuntu, sans-serif;}`)
	webpage.AddStyle(`body {
		font-family: "Nunito", sans-serif;
		font-size: 0.9rem;
		font-weight: 400;
		line-height: 1.6;
		color: #212529;
		text-align: left;
		background-color: #f8fafc;
	}
	.form-select {
		display: block;
		width: 100%;
		padding: .375rem 2.25rem .375rem .75rem;
		font-size: 1rem;
		font-weight: 400;
		line-height: 1.5;
		color: #212529;
		background-color: #fff;
		background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16'%3e%3cpath fill='none' stroke='%23343a40' stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M2 5l6 6 6-6'/%3e%3c/svg%3e");
		background-repeat: no-repeat;
		background-position: right .75rem center;
		background-size: 16px 12px;
		border: 1px solid #ced4da;
		border-radius: .25rem;
		-webkit-appearance: none;
		-moz-appearance: none;
		appearance: none;
	}`)
	webpage.AddChild(hb.NewHTML(content))
	return webpage
}
