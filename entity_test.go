package cms

import (
	"testing"
	//"database/sql"
   _ "github.com/mattn/go-sqlite3"
  "gorm.io/driver/sqlite"
  "gorm.io/gorm"
)

func InitDB(filepath string) *gorm.DB /**sql.DB*/ {

  db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{})
  if err != nil { panic(err) }
  // Auto Migrate
  db.AutoMigrate(&Entity{})
	// db, err := sql.Open("sqlite3", filepath)
	// if err != nil { panic(err) }
	// if db == nil { panic("db nil") }
	// return db
	return db
}
func TestCreation(t *testing.T) {
	InitDB("entity_test.db")
	//if sql != "SELECT * FROM 'user';" {
		//t.Fatalf(sql)
	//}
	entity := EntityCreate("post")
	if entity == nil{
		t.Fatalf("Entity could not be created")
	}
}
