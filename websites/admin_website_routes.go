package cms

// import (
// 	"net/http"

// 	"github.com/gouniverse/bs"
// 	"github.com/gouniverse/cdn"
// 	. "github.com/gouniverse/hb"
// 	"github.com/gouniverse/responses"
// )

// func (cms Cms) pageWebsitesWebsiteManager(w http.ResponseWriter, r *http.Request) {
// 	endpoint := r.Context().Value(keyEndpoint).(string)

// 	// websites, err := cms.EntityStore.EntityList(entitystore.EntityQueryOptions{
// 	// 	EntityType: ENTITY_TYPE_WEBSITE,
// 	// 	Offset:     0,
// 	// 	Limit:      200,
// 	// 	SortBy:     COLUMN_ID,
// 	// 	SortOrder:  sb.ASC,
// 	// })

// 	// if err != nil {
// 	// 	api.Respond(w, r, api.Error("Page list failed to be retrieved "+err.Error()))
// 	// 	return
// 	// }

// 	header := cms.cmsHeader(endpoint)
// 	breadcrumbs := cms.cmsBreadcrumbs([]bs.Breadcrumb{
// 		{
// 			URL:  endpoint,
// 			Name: "Home",
// 		},
// 		{
// 			URL:  (endpoint + "?path=" + PathPagesPageManager),
// 			Name: "Websites",
// 		},
// 	})

// 	button := NewButton().
// 		Text("New page").
// 		Class("btn btn-success float-end").
// 		Attr("v-on:click", "showPageCreateModal")

// 	heading := NewHeading1().
// 		Text("Page Manager").
// 		Child(button)

// 	table := NewTable().
// 		ID("TableWebsites").
// 		Class("table table-responsive table-striped mt-3").
// 		Child(NewThead().
// 			Child(NewTR().
// 				Child(NewTH().
// 					Text("Name")).
// 				Child(NewTH().
// 					Text("Status")).
// 				Child(NewTH().
// 					Text("Actions").
// 					Style("width:150px;"))).
// 			Child(NewTbody()))

// 	container := NewDiv().
// 		Class("container").
// 		ID("page-manager").
// 		HTML(header).
// 		Child(heading).
// 		HTML(breadcrumbs).
// 		Child(table)
// 		// Child(pagePagesPageCreateModal()).
// 		// Child(pagePagesPageTrashModal())

// 	if cms.funcLayout("") != "" {
// 		out := NewWrap().
// 			Child(NewScriptURL(cdn.JqueryDataTablesCss_1_13_4())).
// 			Child(NewHTML(container.ToHTML())).
// 			Child(NewScriptURL(cdn.Jquery_3_6_4())).
// 			Child(NewScriptURL(cdn.VueJs_3())).
// 			Child(NewScriptURL(cdn.Sweetalert2_10())).
// 			Child(NewScriptURL(cdn.JqueryDataTablesJs_1_13_4())).
// 			// Child(NewScript(inlineScript)).
// 			ToHTML()
// 		responses.HTMLResponse(w, r, cms.funcLayout(out))
// 		return
// 	}

// 	webpage := WebpageComplete("Website Manager", container.ToHTML())
// 	webpage.AddStyleURL(cdn.JqueryDataTablesCss_1_13_4())
// 	webpage.AddScriptURL(cdn.JqueryDataTablesJs_1_13_4())
// 	// webpage.AddScript(inlineScript)

// 	responses.HTMLResponse(w, r, cms.funcLayout(webpage.ToHTML()))
// }
