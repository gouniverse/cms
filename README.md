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

In its simplest initialization the CMS package accepts a GORM DB instance

```
cms.Init(cms.Config{
		DbInstance: gormDB,
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
