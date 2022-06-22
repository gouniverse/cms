# GoLang CMS

![tests](https://github.com/gouniverse/cms/workflows/tests/badge.svg)

PREVIEW ONLY. NOT STABLE API CAN AND WILL CHANGE FREQUENTLY

A "plug-and-play" content managing system (CMS) for GoLang that does its job and stays out of your way.

## Introduction

All of the existing GoLang CMSs require a full installations from scratch. Its impossible to just add them to an exiting Go application, and even when added feel like you don't get what you hoped for.

This package allows to add a content management system as a module dependency, which can be easily updated or removed as required to ANY Go app. It is fully self contained, and does not require any additional packages or dependencies. Removal is also a breeze just remove the module.

## Features

- Templates (CMS)
- Pages (CMS)
- Blocks (CMS)
- Menus (CMS)
- Settings (CMS)
- Custom Types

# Simple Initialization

In its simplest initialization the CMS package accepts a standard DB instance

```
db, err := mainDb(utils.Env("DB_DRIVER"), utils.Env("DB_HOST"), utils.Env("DB_PORT"), utils.Env("DB_DATABASE"), utils.Env("DB_USERNAME"), utils.Env("DB_PASSWORD"))

if err != nil {
	log.Panic("Database is NIL: " + err.Error())
	return
}

if db == nil {
	log.Panic("Database is NIL")
	return
}

cms.Init(cms.Config{
	DbInstance: db,
})
```

# Initialization with CMS types

```
cms.Init(cms.Config{
    DbInstance:      db,
    EnableTemplates: true,
    EnablePages:     true,
    EnableBlocks:    true,
})
```

# Initialization with Settings

```
cms.Init(cms.Config{
    DbInstance:      db,
    EnableSettings:  true,
})
```

# Initialization with Custom Entity types

```
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

cms.Init(cms.Config{
    DbInstance:      db,
    CustomEntityList: entityList(),
})
```

## Development Instructions

There is a development directory that allows you to quickly start working on the project or simply to preview

Instructions how to start are in the README file in the directory

[README] https://github.com/gounverse/development/README.md


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

# Notable
- https://github.com/tenox7/wfm - stabdalone file manager



