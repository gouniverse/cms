package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gouniverse/cms"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/utils"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	log.Println("1. Initializing environment variables...")
	utils.EnvInitialize()

	log.Println("2. Initializing database...")
	db := mainDb(utils.Env("DB_DRIVER"), utils.Env("DB_HOST"), utils.Env("DB_PORT"), utils.Env("DB_DATABASE"), utils.Env("DB_USERNAME"), utils.Env("DB_PASSWORD"))

	if db == nil {
		log.Println(utils.FileExists(".env"))
		log.Panic("Database is NIL")
		return
	}

	log.Println("3. Initializing CMS...")
	cms.Init(cms.Config{
		DbInstance:       db,
		EnableTemplates:  true,
		EnablePages:      true,
		EnableBlocks:     true,
		EnableSettings:   true,
		CustomEntityList: entityList(),
	})

	log.Println("4. Starting server on http://" + utils.Env("SERVER_HOST") + ":" + utils.Env("SERVER_PORT") + " ...")
	log.Println("URL: http://" + utils.Env("APP_URL") + " ...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", cms.Router)
	mux.HandleFunc("/cms", cms.Router)
	mux.HandleFunc("/embeddedcms", pageDashboardWithEmbeddedCms)

	srv := &http.Server{
		Handler: mux,
		Addr:    utils.Env("SERVER_HOST") + ":" + utils.Env("SERVER_PORT"),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func pageDashboardWithEmbeddedCms(w http.ResponseWriter, r *http.Request) {
	leftMenu := hb.NewHTML("<a href='/embeddedcms'>Embedded CMS</a>")
	iframe := hb.NewHTML("<iframe src=\"/\" style='width:100%;height:2000px;border:none;' scrolling='no'></iframe>")
	layout := hb.NewHTML("<table style='width:100%;height:100%;'><tr><td style='width:300px;vertical-align:top;'>" + leftMenu.ToHTML() + "</td><td style='vertical-align:top;'>" + iframe.ToHTML() + "</td></tr></table>")
	webpage := hb.NewWebpage().AddChild(layout)
	w.Write([]byte(webpage.ToHTML()))
}

func mainDb(driverName string, dbHost string, dbPort string, dbName string, dbUser string, dbPass string) *gorm.DB {
	var db *gorm.DB
	var err error

	if driverName == "sqlite" {
		dsn := dbName
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	}
	if driverName == "mysql" {
		// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
		dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
		log.Println(dsn)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}
	//defer db.Close()

	if err != nil {
		panic("Failed to connect to the database")
	}

	return db
}

func entityList() []cms.CustomEntityStructure {
	list := []cms.CustomEntityStructure{}
	list = append(list, cms.CustomEntityStructure{
		Group:     "Users",
		Type:      "user",
		TypeLabel: "User",
		AttributeList: []cms.CustomAttributeStructure{
			{
				Name:             "first_name",
				Type:             "string",
				FormControlLabel: "First Name",
				FormControlType:  "input",
				FormControlHelp:  "The first name of the user",
			},
			{
				Name:             "last_name",
				Type:             "string",
				FormControlLabel: "Last Name",
				FormControlType:  "input",
				FormControlHelp:  "The last name of the user",
			},
			{
				Name:             "email",
				Type:             "string",
				FormControlLabel: "E-mail",
				FormControlType:  "input",
				FormControlHelp:  "The e-mail address of the user",
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
		Group:     "Shop",
		Type:      "shop_product",
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
		Type:      "makeawish",
		TypeLabel: "Make-a-Wish",
		AttributeList: []cms.CustomAttributeStructure{
			{
				Name:             "wish",
				Type:             "string",
				FormControlLabel: "Wish",
				FormControlType:  "textarea",
				FormControlHelp:  "The wish that was made",
			},
			{
				Name:             "referral",
				Type:             "string",
				FormControlLabel: "Referral",
				FormControlType:  "input",
				FormControlHelp:  "Whare the wish was made from",
			},
		},
	})
	list = append(list, cms.CustomEntityStructure{
		Group:     "Shop",
		Type:      "shop_order",
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
		Group:     "Shop",
		Type:      "shop_order_line_item",
		TypeLabel: "Order Line Item",
		AttributeList: []cms.CustomAttributeStructure{
			{
				Name:             "order_id",
				Type:             "string",
				FormControlLabel: "Order ID",
				FormControlType:  "input",
				FormControlHelp:  "The order the item belongs to",
				BelongsToType:    "shop_order",
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
	// {
	// 	"type":"product",
	// 	"price":"12.00",
	// }
	// list := []map[string]interface{}{
	// 	{
	// 		"type": "product",
	// 		"attributes": []map[string]interface{}{
	// 			{
	// 				"name": "title",
	// 				"type": "string",
	// 				"rule": "required",
	// 				// type - one of text, textarea, select, hidden, html
	// 				// name - name of the input field as seen in the request
	// 				// label - publicly visible name
	// 				// width - width of the field - min 1, max 12
	// 				// rule - rules for the field, used when validating
	// 				// value - value of the field
	// 				// options - array of options (used by the select type)
	// 				// html - raw HTML to be displayed as-is (used by the html type)

	// 			},
	// 			{
	// 				"name": "price",
	// 				"type": "float",
	// 				"rule": "required",
	// 			},
	// 			{
	// 				"name": "image_url",
	// 				"type": "string",
	// 			},
	// 		},
	// 	},
	// }
	// return list
}
