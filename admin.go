package cms

import (
	"context"
	"net/http"

	"strconv"

	hb "github.com/gouniverse/html"
	"github.com/gouniverse/utils"
)

// Router shows the admin page
func Router(w http.ResponseWriter, r *http.Request) {
	path := utils.Req(r, "path", "home")

	if path == "" {
		path = PathHome
	}

	ctx := context.WithValue(r.Context(), keyEndpoint, r.URL.Path)

	routeFunc := getRoute(path)
	routeFunc(w, r.WithContext(ctx))
	return
}

func getRoute(route string) func(w http.ResponseWriter, r *http.Request) {
	routes := map[string]func(w http.ResponseWriter, r *http.Request){
		PathHome:                        pageHome,
		PathBlocksBlockCreateAjax:       pageBlocksBlockCreateAjax,
		PathBlocksBlockManager:          pageBlocksBlockManager,
		PathBlocksBlockUpdate:           pageBlocksBlockUpdate,
		PathBlocksBlockUpdateAjax:       pageBlocksBlockUpdateAjax,
		PathMenusMenuCreateAjax:         pageMenusMenuCreateAjax,
		PathMenusMenuManager:            pageMenusMenuManager,
		PathMenusMenuUpdate:             pageMenusMenuUpdate,
		PathMenusMenuItemsFetchAjax:     pageMenusMenuItemsFetchAjax,
		PathMenusMenuItemsUpdate:        pageMenusMenuItemsUpdate,
		PathMenusMenuItemsUpdateAjax:    pageMenusMenuItemsUpdateAjax,
		PathMenusMenuUpdateAjax:         pageMenusMenuUpdateAjax,
		PathPagesPageCreateAjax:         pagePagesPageCreateAjax,
		PathPagesPageManager:            pagePagesPageManager,
		PathPagesPageUpdate:             pagePagesPageUpdate,
		PathPagesPageUpdateAjax:         pagePagesPageUpdateAjax,
		PathTemplatesTemplateCreateAjax: pageTemplatesTemplateCreateAjax,
		PathTemplatesTemplateManager:    pageTemplateTemplateManager,
		PathTemplatesTemplateUpdate:     pageTemplatesTemplateUpdate,
		PathTemplatesTemplateUpdateAjax: pageTemplatesTemplateUpdateAjax,
		PathWidgetsWidgetCreateAjax:     pageWidgetsWidgetCreateAjax,
		PathWidgetsWidgetManager:        pageWidgetsWidgetManager,
		PathWidgetsWidgetUpdate:         pageWidgetsWidgetUpdate,
		PathWidgetsWidgetUpdateAjax:     pageWidgetsWidgetUpdateAjax,
	}
	// log.Println(route)
	if val, ok := routes[route]; ok {
		return val
	}
	return routes[PathHome]
}

func pageHome(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Context().Value(keyEndpoint).(string)
	// log.Println(endpoint)

	header := cmsHeader(endpoint)
	breadcrums := cmsBreadcrumbs(map[string]string{
		endpoint: "Home",
	})

	container := hb.NewDiv().Attr("class", "container").Attr("id", "page-manager")
	heading := hb.NewHeading1().HTML("CMS Dashboard")

	container.AddChild(hb.NewHTML(header))
	container.AddChild(heading)
	container.AddChild(hb.NewHTML(breadcrums))

	h := container.ToHTML()

	webpage := Webpage("Home", h)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(webpage.ToHTML()))
}

func cmsBreadcrumbs(breadcrumbs map[string]string) string {
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

func cmsHeader(endpoint string) string {
	linkHome := hb.NewHyperlink().HTML("Dashboard").Attr("href", endpoint+"").Attr("class", "nav-link")
	linkBlocks := hb.NewHyperlink().HTML("Blocks ").Attr("href", endpoint+"?path="+PathBlocksBlockManager).Attr("class", "nav-link")
	linkMenus := hb.NewHyperlink().HTML("Menus ").Attr("href", endpoint+"?path="+PathMenusMenuManager).Attr("class", "nav-link")
	linkPages := hb.NewHyperlink().HTML("Pages ").Attr("href", endpoint+"?path="+PathPagesPageManager).Attr("class", "nav-link")
	linkTemplates := hb.NewHyperlink().HTML("Templates ").Attr("href", endpoint+"?path="+PathTemplatesTemplateManager).Attr("class", "nav-link")
	linkWidgets := hb.NewHyperlink().HTML("Widgets ").Attr("href", endpoint+"?path="+PathWidgetsWidgetManager).Attr("class", "nav-link")
	blocksCount := EntityCount("block")
	menusCount := EntityCount("menu")
	pagesCount := EntityCount("page")
	templatesCount := EntityCount("template")
	widgetsCount := EntityCount("widget")
	//log.Println(blocksCount)
	//log.Println(pagesCount)
	//log.Println(templatesCount)

	ulNav := hb.NewUL().Attr("class", "nav  nav-pills justify-content-center")
	ulNav.AddChild(hb.NewLI().Attr("class", "nav-item").AddChild(linkHome))
	ulNav.AddChild(hb.NewLI().Attr("class", "nav-item").AddChild(linkTemplates.AddChild(hb.NewSpan().Attr("class", "badge bg-secondary").HTML(strconv.FormatUint(templatesCount, 10)))))
	ulNav.AddChild(hb.NewLI().Attr("class", "nav-item").AddChild(linkPages.AddChild(hb.NewSpan().Attr("class", "badge bg-secondary").HTML(strconv.FormatUint(pagesCount, 10)))))
	ulNav.AddChild(hb.NewLI().Attr("class", "nav-item").AddChild(linkMenus.AddChild(hb.NewSpan().Attr("class", "badge bg-secondary").HTML(strconv.FormatUint(menusCount, 10)))))
	ulNav.AddChild(hb.NewLI().Attr("class", "nav-item").AddChild(linkBlocks.AddChild(hb.NewSpan().Attr("class", "badge bg-secondary").HTML(strconv.FormatUint(blocksCount, 10)))))
	ulNav.AddChild(hb.NewLI().Attr("class", "nav-item").AddChild(linkWidgets.AddChild(hb.NewSpan().Attr("class", "badge bg-secondary").HTML(strconv.FormatUint(widgetsCount, 10)))))
	// add Translations

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
		"https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/css/bootstrap.min.css",
	})
	webpage.AddScriptURLs([]string{
		"https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/js/bootstrap.bundle.min.js",
		"http://code.jquery.com/jquery-3.5.1.min.js",
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
	}`)
	webpage.AddChild(hb.NewHTML(content))
	return webpage
}
