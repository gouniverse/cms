package cms

import (
	//"log"
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

func TestEntityCreate(t *testing.T) {
	Init("sqlite", "entity_test.db")
	entity := EntityCreate("post")
	if entity == nil{
		t.Fatalf("Entity could not be created")
	}
}

func TestEntityCreateWithAttributes(t *testing.T) {
	Init("sqlite", "entity_test.db")
	entity := EntityCreateWithAttributes("post", map[string]interface{}{
		"name":"Hello world",
	})
	if entity == nil{
		t.Fatalf("Entity could not be created")
	}
	if entity.GetAttributeValue("name","").(string) != "Hello world"{
		t.Fatalf("Entity attribute mismatch")
	}
}
