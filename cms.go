package cms

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB
var prefix string

// Config contains the configurations for the auth package
type Config struct {
	DbInstance *gorm.DB
	DbDriver   string
	DbDsn      string
}

var (
	configuration Config
)

// Init initializes the CMS
func Init(config Config) {
	if config.DbInstance == nil && (config.DbDriver == "" || config.DbDsn == "") {
		log.Panicln("Either DbInstance or DnDriver and DbDsn are required field")
	}

	prefix = "cms_"
	
	if config.DbDriver != "" {
		db, err := gorm.Open(sqlite.Open(config.DbDsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		// dbInstance = db
		config.DbInstance = db
	}

	// Migrate the schema
	config.DbInstance.AutoMigrate(&Entity{})
	config.DbInstance.AutoMigrate(&EntityAttribute{})
	configuration = config
}

// GetDb returns an instance to the CMS database
func GetDb() *gorm.DB {
	return configuration.DbInstance
}
