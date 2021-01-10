package cms

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB
var prefix string

// Init initializes the CMS
func Init(driverName string, dsn string) {
	prefix = "cms_"
	//sqlDB, err := sql.Open("mysql", "mydb_dsn")
	// gormDB, err := gorm.Open(mysql.New(mysql.Config{
	// 	Conn: sqlDB,
	// }), &gorm.Config{})
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Entity{})
	db.AutoMigrate(&EntityAttribute{})

	dbInstance = db
}

// GetDb returns an instance to the CMS database
func GetDb() *gorm.DB {
	return dbInstance
}
