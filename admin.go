package cms

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/responses"
	"github.com/gouniverse/utils"
)

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

func (cms Cms) getRoute(route string) func(w http.ResponseWriter, r *http.Request) {
	routes := map[string]func(w http.ResponseWriter, r *http.Request){
		PathHome: cms.pageHome,

		// START: Blocks
		PathBlocksBlockCreateAjax: cms.pageBlocksBlockCreateAjax,
		PathBlocksBlockDeleteAjax: cms.pageBlocksBlockDeleteAjax,
		PathBlocksBlockManager:    cms.pageBlocksBlockManager,
		PathBlocksBlockUpdate:     cms.pageBlocksBlockUpdate,
		PathBlocksBlockTrashAjax:  cms.pageBlocksBlockTrashAjax,
		PathBlocksBlockUpdateAjax: cms.pageBlocksBlockUpdateAjax,
		// END: Blocks

		// START: Menus
		PathMenusMenuCreateAjax:      cms.pageMenusMenuCreateAjax,
		PathMenusMenuManager:         cms.pageMenusMenuManager,
		PathMenusMenuUpdate:          cms.pageMenusMenuUpdate,
		PathMenusMenuItemsFetchAjax:  cms.pageMenusMenuItemsFetchAjax,
		PathMenusMenuItemsUpdate:     cms.pageMenusMenuItemsUpdate,
		PathMenusMenuItemsUpdateAjax: cms.pageMenusMenuItemsUpdateAjax,
		PathMenusMenuUpdateAjax:      cms.pageMenusMenuUpdateAjax,
		// END: Menus

		// START: Pages
		PathPagesPageCreateAjax: cms.pagePagesPageCreateAjax,
		PathPagesPageManager:    cms.pagePagesPageManager,
		PathPagesPageTrashAjax:  cms.pagePagesPageTrashAjax,
		PathPagesPageUpdate:     cms.pagePagesPageUpdate,
		PathPagesPageUpdateAjax: cms.pagePagesPageUpdateAjax,
		// END: Pages

		// START: Templates
		PathTemplatesTemplateCreateAjax: cms.pageTemplatesTemplateCreateAjax,
		PathTemplatesTemplateManager:    cms.pageTemplatesTemplateManager,
		PathTemplatesTemplateTrashAjax:  cms.pageTemplatesTemplateTrashAjax,
		PathTemplatesTemplateUpdate:     cms.pageTemplatesTemplateUpdate,
		PathTemplatesTemplateUpdateAjax: cms.pageTemplatesTemplateUpdateAjax,
		// END: Templates

		// START: Widgets
		PathWidgetsWidgetCreateAjax: cms.pageWidgetsWidgetCreateAjax,
		PathWidgetsWidgetManager:    cms.pageWidgetsWidgetManager,
		PathWidgetsWidgetUpdate:     cms.pageWidgetsWidgetUpdate,
		PathWidgetsWidgetUpdateAjax: cms.pageWidgetsWidgetUpdateAjax,
		// END: Widgets

		// START: Settings
		PathSettingsSettingCreateAjax: cms.pageSettingsSettingCreateAjax,
		PathSettingsSettingDeleteAjax: cms.pageSettingsSettingDeleteAjax,
		PathSettingsSettingManager:    cms.pageSettingsSettingManager,
		PathSettingsSettingUpdate:     cms.pageSettingsSettingUpdate,
		PathSettingsSettingUpdateAjax: cms.pageSettingsSettingUpdateAjax,
		// END: Settings

		// START: Users
		PathUsersUserCreateAjax: cms.pageUsersUserCreateAjax,
		PathUsersUserTrashAjax:  cms.pageUsersUserTrashAjax,
		PathUsersUserManager:    cms.pageUsersUserManager,
		PathUsersUserUpdate:     cms.pageUsersUserUpdate,
		PathUsersUserUpdateAjax: cms.pageUsersUserUpdateAjax,
		// END: Users

		// START: Custom Entities
		PathEntitiesEntityCreateAjax: cms.pageEntitiesEntityCreateAjax,
		PathEntitiesEntityManager:    cms.pageEntitiesEntityManager,
		PathEntitiesEntityUpdate:     cms.pageEntitiesEntityUpdate,
		PathEntitiesEntityUpdateAjax: cms.pageEntitiesEntityUpdateAjax,
		// END: Custom Entities

	}
	// log.Println(route)
	if val, ok := routes[route]; ok {
		return val
	}

	return routes[PathHome]
}

// func (cms Cms) pageUserHome(w http.ResponseWriter, r *http.Request) {
// 	endpoint := r.Context().Value(keyEndpoint).(string)
// 	// log.Println(endpoint)
// 	header := cms.cmsHeader(endpoint)
// 	breadcrumbs := cms.cmsBreadcrumbs(map[string]string{
// 		endpoint: "Home",
// 	})
// 	container := hb.NewDiv().Attr("class", "container").Attr("id", "page-manager")
// 	heading := hb.NewHeading1().HTML("User Dashboard")
// 	container.AddChild(hb.NewHTML(header))
// 	container.AddChild(heading)
// 	container.AddChild(hb.NewHTML(breadcrumbs))
// 	h := container.ToHTML()
// 	webpage := Webpage("Home", h)
// 	w.WriteHeader(200)
// 	w.Header().Set("Content-Type", "text/html")
// 	w.Write([]byte(webpage.ToHTML()))
// }

func (cms Cms) pageHome(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	// log.Println(endpoint)

	header := cms.cmsHeader(endpoint)
	breadcrumbs := cms.cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "page-manager")
	heading := hb.NewHeading1().HTML("CMS Dashboard")

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrumbs))

	h := container.ToHTML()

	webpage := Webpage("Home", h).ToHTML()

	responses.HTMLResponse(w, r, cms.funcLayout(webpage))
}

func (cms Cms) cmsBreadcrumbs(breadcrumbs map[string]string) string {
	nav := hb.NewNav().Attr("aria-label", "breadcrumb")
	ol := hb.NewOL().Attr("class", "breadcrumb")
	for k, v := range breadcrumbs {
		li := hb.NewLI().Attr("class", "breadcrumb-item")
		link := hb.NewHyperlink().HTML(v).Attr("href", k)

		li.AddChild(link)

		ol.AddChild(li)
	}
	nav.AddChild(ol)
	return nav.ToHTML()
}

func (cms Cms) cmsHeader(endpoint string) string {
	linkHome := hb.NewHyperlink().HTML("Dashboard").Attr("href", endpoint+"").Attr("class", "nav-link")
	linkBlocks := hb.NewHyperlink().HTML("Blocks ").Attr("href", endpoint+"?path="+PathBlocksBlockManager).Attr("class", "nav-link")
	linkMenus := hb.NewHyperlink().HTML("Menus ").Attr("href", endpoint+"?path="+PathMenusMenuManager).Attr("class", "nav-link")
	linkPages := hb.NewHyperlink().HTML("Pages ").Attr("href", endpoint+"?path="+PathPagesPageManager).Attr("class", "nav-link")
	linkTemplates := hb.NewHyperlink().HTML("Templates ").Attr("href", endpoint+"?path="+PathTemplatesTemplateManager).Attr("class", "nav-link")
	linkWidgets := hb.NewHyperlink().HTML("Widgets ").Attr("href", endpoint+"?path="+PathWidgetsWidgetManager).Attr("class", "nav-link")
	linkSettings := hb.NewHyperlink().HTML("Settings").Attr("href", endpoint+"?path="+PathSettingsSettingManager).Attr("class", "nav-link")
	linkTranslations := hb.NewHyperlink().HTML("Translations").Attr("href", "#").Attr("class", "nav-link")
	blocksCount, _ := cms.EntityStore.EntityCount(entitystore.EntityQueryOptions{
		EntityType: "block",
	})
	menusCount, _ := cms.EntityStore.EntityCount(entitystore.EntityQueryOptions{
		EntityType: "menu",
	})
	pagesCount, _ := cms.EntityStore.EntityCount(entitystore.EntityQueryOptions{
		EntityType: "page",
	})
	templatesCount, _ := cms.EntityStore.EntityCount(entitystore.EntityQueryOptions{
		EntityType: "template",
	})
	widgetsCount, _ := cms.EntityStore.EntityCount(entitystore.EntityQueryOptions{
		EntityType: "widget",
	})

	ulNav := hb.NewUL().Attr("class", "nav  nav-pills justify-content-center")
	ulNav.AddChild(hb.NewLI().Attr("class", "nav-item").AddChild(linkHome))

	if cms.templatesEnabled {
		ulNav.AddChild(hb.NewLI().Attr("class", "nav-item").AddChild(linkTemplates.AddChild(hb.NewSpan().Attr("class", "badge bg-secondary").HTML(strconv.FormatInt(templatesCount, 10)))))
	}

	if cms.pagesEnabled {
		ulNav.AddChild(hb.NewLI().Attr("class", "nav-item").AddChild(linkPages.AddChild(hb.NewSpan().Attr("class", "badge bg-secondary").HTML(strconv.FormatInt(pagesCount, 10)))))
	}

	if cms.menusEnabled {
		ulNav.AddChild(hb.NewLI().Attr("class", "nav-item").AddChild(linkMenus.AddChild(hb.NewSpan().Attr("class", "badge bg-secondary").HTML(strconv.FormatInt(menusCount, 10)))))
	}

	if cms.blocksEnabled {
		ulNav.AddChild(hb.NewLI().Attr("class", "nav-item").AddChild(linkBlocks.AddChild(hb.NewSpan().Attr("class", "badge bg-secondary").HTML(strconv.FormatInt(blocksCount, 10)))))
	}

	if cms.widgetsEnabled {
		ulNav.AddChild(hb.NewLI().Attr("class", "nav-item").AddChild(linkWidgets.AddChild(hb.NewSpan().Attr("class", "badge bg-secondary").HTML(strconv.FormatInt(widgetsCount, 10)))))
	}

	if cms.translationsEnabled {
		ulNav.AddChild(hb.NewLI().Attr("class", "nav-item").AddChild(linkTranslations))
	}

	if cms.settingsEnabled {
		ulNav.AddChild(hb.NewLI().Attr("class", "nav-item").AddChild(linkSettings))
	}
	// add Translations

	for _, entity := range cms.customEntityList {
		linkEntity := hb.NewHyperlink().HTML(entity.TypeLabel).Attr("href", endpoint+"?path=entities/entity-manager&type="+entity.Type).Attr("class", "nav-link")
		ulNav.AddChild(hb.NewLI().Attr("class", "nav-item").AddChild(linkEntity))
	}

	divCard := hb.NewDiv().Attr("class", "card card-default mt-3 mb-3")
	divCardBody := hb.NewDiv().Attr("class", "card-body").Attr("style", "padding: 2px;")
	return divCard.AddChild(divCardBody.AddChild(ulNav)).ToHTML()
}

// Webpage returns the webpage template for the website
func Webpage(title, content string) *hb.Webpage {
	faviconImgCms := `data:image/x-icon;base64,AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAmzKzAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABEQEAAQERAAEAAQABAAEAAQABAQEBEQABAAEREQEAAAERARARAREAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD//wAA//8AAP//AAD//wAA//8AAP//AAD//wAAi6MAALu7AAC6owAAuC8AAIkjAAD//wAA//8AAP//AAD//wAA`
	app := ""
	webpage := hb.NewWebpage()
	webpage.SetTitle(title)
	webpage.SetFavicon(faviconImgCms)

	webpage.AddStyleURLs([]string{
		"https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta3/dist/css/bootstrap.min.css",
	})
	webpage.AddScriptURLs([]string{
		"https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta3/dist/js/bootstrap.bundle.min.js",
		"https://code.jquery.com/jquery-3.6.0.min.js",
		"https://unpkg.com/vue@next",
		"https://cdn.jsdelivr.net/npm/sweetalert2@9",
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
