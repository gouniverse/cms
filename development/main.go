package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gouniverse/cms"
	"github.com/gouniverse/utils"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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
		DbInstance: db,
		//CustomEntityList: entityList(),
	})

	log.Println("4. Starting server on http://" + utils.Env("SERVER_HOST") + ":" + utils.Env("SERVER_PORT") + " ...")
	log.Println("URL: http://" + utils.Env("APP_URL") + " ...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", cms.Router)
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
