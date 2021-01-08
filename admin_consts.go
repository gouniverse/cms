package cms

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	keyEndpoint                            = contextKey("endpoint")
	PathHome                        string = "home"
	PathBlocksBlockManager          string = "blocks/block-manager"
	PathBlocksBlockCreateAjax       string = "blocks/block-create-ajax"
	PathBlocksBlockUpdate           string = "blocks/block-update"
	PathBlocksBlockUpdateAjax       string = "blocks/block-update-ajax"
	PathMenusMenuManager            string = "menus/menu-manager"
	PathMenusMenuCreateAjax         string = "manus/menu-create-ajax"
	PathMenusMenuUpdate             string = "manus/menu-update"
	PathMenusMenuItemsFetchAjax     string = "menus/menu-items-fetch-ajax"
	PathMenusMenuItemsUpdateAjax    string = "menus/menu-items-update-ajax"
	PathMenusMenuItemsUpdate        string = "manus/menu-items-update"
	PathMenusMenuUpdateAjax         string = "menus/menu-update-ajax"
	PathPagesPageManager            string = "pages/page-manager"
	PathPagesPageCreateAjax         string = "pages/page-create-ajax"
	PathPagesPageUpdate             string = "pages/page-update"
	PathPagesPageUpdateAjax         string = "pages/page-update-ajax"
	PathTemplatesTemplateManager    string = "templates/template-manager"
	PathTemplatesTemplateCreateAjax string = "templates/template-create-ajax"
	PathTemplatesTemplateUpdate     string = "templates/template-update"
	PathTemplatesTemplateUpdateAjax string = "templates/template-update-ajax"
	PathWidgetsWidgetManager        string = "widgets/widget-manager"
	PathWidgetsWidgetCreateAjax     string = "widgets/widget-create-ajax"
	PathWidgetsWidgetUpdate         string = "widgets/widget-update"
	PathWidgetsWidgetUpdateAjax     string = "widgets/widget-update-ajax"
)
