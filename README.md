
This project is being outdated, use the new project:

New URL: https://github.com/gouniverse/cmsstore

---

# GoLang CMS <a href="https://gitpod.io/#https://github.com/gouniverse/cms" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

![tests](https://github.com/gouniverse/cms/workflows/tests/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/cms)](https://goreportcard.com/report/github.com/gouniverse/cms)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/cms)](https://pkg.go.dev/github.com/gouniverse/cms)

A "plug-and-play" content managing system (CMS) for GoLang that does its job and stays out of your way.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). You can find a copy of the license at [https://www.gnu.org/licenses/agpl-3.0.en.html](https://www.gnu.org/licenses/agpl-3.0.txt)

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

## Introduction

All of the existing GoLang CMSs require a full installations from scratch. Its impossible to just add them to an exiting Go application, and even when added feel like you don't get what you hoped for.

This package allows to add a content management system as a module dependency, which can be easily updated or removed as required to ANY Go app. It is fully self contained, and does not require any additional packages or dependencies. Removal is also a breeze just remove the module.

## Features
- Entity types
- Templates (CMS)
- Pages (CMS)
- Blocks (CMS)
- Menus (CMS)
- Settings (CMS)
- Translations (CMS)
- Custom Types
- Cache Store
- Log Store
- Session Store
- Task Store (queue for background tasks)

# Simplest Initialization

In its simplest initialization the CMS package accepts a standard DB instance.

However with this simplest initialization, the CMS basically has no capabilities (i.e no database stores can be accessed, no migrations are run, etc).

```go
db, err := mainDb(utils.Env("DB_DRIVER"), utils.Env("DB_HOST"), utils.Env("DB_PORT"), utils.Env("DB_DATABASE"), utils.Env("DB_USERNAME"), utils.Env("DB_PASSWORD"))

if err != nil {
	log.Panic("Database is NIL: " + err.Error())
	return
}

if db == nil {
	log.Panic("Database is NIL")
	return
}

myCms, errCms := cms.NewCms(cms.Config{
	DbInstance:           db,
})
```

# Initialization with entity types

```go
myCms, errCms := cms.NewCms(cms.Config{
	DbInstance:           db,
	EntitiesAutomigrate:  true,
})
```


# Initialization with CMS types

```go
myCms, errCms := cms.NewCms(cms.Config{
    DbInstance:           db,
    EntitiesAutomigrate:  true,
    BlocksEnable:         true,
    MenusEnable:          true,
    PagesEnable:          true,
    TemplatesEnable:      true,
    WidgetsEnable:        true,
    Prefix:               "cms_"
})
```

# Initialization with Settings

```go
myCms, errCms := cms.NewCms(cms.Config{
    DbInstance:           db,
    SettingsAutomigrate:  true,
    SettingsEnable:       true,
})
```

# Initialization with Custom Entity types

```go
func entityList() []cms.CustomEntityStructure {
	list := []cms.CustomEntityStructure{}
	list = append(list, cms.CustomEntityStructure{
		Type:      "product",
		TypeLabel: "Product",
		AttributeList: []cms.CustomAttributeStructure{
			{
				Name:             "title",
				Type:             "string",
				FormControlLabel: "Title",
				FormControlType:  "input",
				FormControlHelp:  "The title which will be displayed to the customer",
			},
			{
				Name:             "description",
				Type:             "string",
				FormControlLabel: "Description",
				FormControlType:  "textarea",
				FormControlHelp:  "The description which will be displayed to the customer",
			},
			{
				Name:             "price",
				Type:             "string",
				FormControlLabel: "Price",
				FormControlType:  "input",
				FormControlHelp:  "The price of the product",
			},
			{
				Name:             "image_url",
				Type:             "string",
				FormControlLabel: "Image URL",
				FormControlType:  "input",
				FormControlHelp:  "The image of the product",
			},
		},
	})
	list = append(list, cms.CustomEntityStructure{
		Type:      "order",
		TypeLabel: "Order",
		AttributeList: []cms.CustomAttributeStructure{
			{
				Name:             "user_id",
				Type:             "string",
				FormControlLabel: "User ID",
				FormControlType:  "input",
				FormControlHelp:  "The ID of the user who made the purchase",
				BelongsToType:    "user",
			},
			{
				Name:             "total",
				Type:             "string",
				FormControlLabel: "Total",
				FormControlType:  "input",
				FormControlHelp:  "Total amount of the order",
			},
		},
	})
	list = append(list, cms.CustomEntityStructure{
		Type:      "order_item",
		TypeLabel: "Order",
		AttributeList: []cms.CustomAttributeStructure{
			{
				Name:             "order_id",
				Type:             "string",
				FormControlLabel: "Order ID",
				FormControlType:  "input",
				FormControlHelp:  "The ID of the order the item belongs to",
			},
			{
				Name:             "product_id",
				Type:             "string",
				FormControlLabel: "Product ID",
				FormControlType:  "input",
				FormControlHelp:  "The ID of the product that is ordered",
				BelongsToType:    "product"
			},
			{
				Name:             "quantity",
				Type:             "string",
				FormControlLabel: "Quantity",
				FormControlType:  "input",
				FormControlHelp:  "How many products are ordered (quantity) in this order item",
			},
			{
				Name:             "subtotal",
				Type:             "string",
				FormControlLabel: "Subtotal",
				FormControlType:  "input",
				FormControlHelp:  "Subtotal amount of the order item",
			},
		},
	})
	return list
}

myCms, errCms := cms.NewCms(cms.Config{
    DbInstance:           db,
    EntitiesAutomigrate:  true,
    CustomEntityList:     entityList(),
})
```

## Cache Store

Some of the data retrieval or processing tasks performed by your application could be CPU intensive or take several seconds to complete. When this is the case, it is common to cache the retrieved data for a time so it can be retrieved quickly on subsequent requests for the same data. 

CMS comes out of the box with and SQL based cache store that can be enabled on demand. The cache store is based on the following project:
https://github.com/gouniverse/cachestore

1. Initialization with Cache Store

```go
myCms, errCms := cms.NewCms(cms.Config{
    DbInstance:        db,
    CacheAutomigrate:  true,
    CacheEnable:       true,
})
```

2. Setting a cache key

```go
isSaved, err := cms.CacheStore.Set("token", "ABCD", 60*60) // 1 hour (60 min * 60 sec)
if isSaved == false {
	log.Println("Saving failed")
	return
}
```

3. Getting a cache key

```go
token, err := cms.CacheStore.Get("token", "") // "" (default)
if token == "" {
	log.Println("Token does not exist or expired")
	return
}
```

## CMS Setup

- Example router (using the Chi router)

```golang
package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes returns the routes of the application
func Routes(cmsRouter http.HandlerFunc) *chi.Mux {
	router := chi.NewRouter()

	router.Route("/admin", func(router chi.Router) {
		router.Use(AdminOnlyMiddleware)
		router.Get("/cms", cmsRouter)
		router.Get("/cms/{catchall:.*}", cmsRouter)
		router.Post("/cms", cmsRouter)
		router.Post("/cms/{catchall:.*}", cmsRouter)
	})

	router.Get("/", CmsController{}.Frontend)
	router.Get("/{catchall:.*}", CmsController{}.Frontend)
	router.Post("/{catchall:.*}", CmsController{}.Frontend)
	return router
}
```

### CMS URL Patterns

The following URL patterns are supported:

- :any - ([^/]+)
- :num - ([0-9]+)
- :all - (.*)
- :string - ([a-zA-Z]+)
- :number - ([0-9]+)
- :numeric - ([0-9-.]+)
- :alpha - ([a-zA-Z0-9-_]+)

Example:

```
/blog/:num/:any
/shop/product/:num/:any
```

## Development Instructions

There is a development directory that allows you to quickly start working on the project or simply to preview

Instructions how to start can be found <a href="https://github.com/gouniverse/cms/tree/main/development">in the directory</a>


## Similar Projects Built in GoLang

- https://github.com/ponzu-cms/ponzu - last updated 2020
- https://github.com/dionyself/golang-cms - last update 2018
- https://github.com/tejo/boxed - last updated 2018
- https://github.com/fragmenta/fragmenta-cms - last updated 2018
- https://github.com/ngocphuongnb/tetua
- https://github.com/ketchuphq/ketchup - last update 2018
- https://github.com/monsti/monsti - last update 2018
- https://github.com/xiusin/pinecms
- https://github.com/gmemstr/pogo - last update 2018
- https://github.com/xushuhui/lin-cms-go - last update 2021
- https://github.com/uberswe/beubo - last update 2021
- https://github.com/digimakergo
- https://github.com/jlelse/GoBlog
- https://github.com/daptin/daptin
- https://github.com/fesiong/goblog

# Notable
- https://github.com/tenox7/wfm - stabdalone file manager



