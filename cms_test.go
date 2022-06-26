package cms

import (
	"testing"

	"database/sql"
	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

func init() {

	// mailServer := smtpmock.New(smtpmock.ConfigurationAttr{
	// 	LogToStdout:       false, // enable if you have errors sending emails
	// 	LogServerActivity: true,
	// 	PortNumber:        32435,
	// 	HostAddress:       "127.0.0.1",
	// })

	// if err := mailServer.Start(); err != nil {
	// 	fmt.Println(err)
	// }
}

func TestCmsTestSuite(t *testing.T) {
	suite.Run(t, new(CmsTestSuite))
}

type CmsTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *CmsTestSuite) SetupTest() {
	// suite.VariableThatShouldStartAtFive = 5

	// Setup()
}

// All methods that begin with "Test" are run as tests within a
// suite.
// func (suite *CmsTestSuite) TestExample() {
// 	assert.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)
// }

//TestAuth tests the auth page
func (suite *CmsTestSuite) TestCmsInit() {
	db, err := mainDb("sqlite", "", "", "test_init.db", "", "")
	defer db.Close()
	assert.Nil(suite.T(), err, "DB error")

	assert.False(suite.T(), configuration.EnableBlocks, "Enable blocks MUST BE false before init")
	assert.False(suite.T(), configuration.EnableCache, "Enable cache MUST BE false before init")
	assert.False(suite.T(), configuration.EnableLogs, "Enable logs MUST BE false before init")
	assert.False(suite.T(), configuration.EnablePages, "Enable pages MUST BE false before init")
	assert.False(suite.T(), configuration.EnableSettings, "Enable pages MUST BE false before init")
	assert.False(suite.T(), configuration.EnableSession, "Enable pages MUST BE false before init")
	assert.False(suite.T(), configuration.EnableTemplates, "Enable templates MUST BE false before init")

	Init(Config{
		DbInstance:      db,
		EnableCache:     true,
		EnableLogs:      true,
		EnablePages:     true,
		EnableBlocks:    true,
		EnableSettings:  true,
		EnableSession:   true,
		EnableTemplates: true,
		// CustomEntityList: entityList(),
	})

	assert.True(suite.T(), configuration.EnableBlocks, "Enable blocks MUST BE true after init")
	assert.True(suite.T(), configuration.EnableCache, "Enable cache MUST BE true after init")
	assert.True(suite.T(), configuration.EnableLogs, "Enable logs MUST BE true after init")
	assert.True(suite.T(), configuration.EnablePages, "Enable pages MUST BE true after init")
	assert.True(suite.T(), configuration.EnableSettings, "Enable pages MUST BE true after init")
	assert.True(suite.T(), configuration.EnableSession, "Enable pages MUST BE true after init")
	assert.True(suite.T(), configuration.EnableTemplates, "Enable templates MUST BE true after init")

	// pages, err := EntityStore.EntityList("page", 0, 10, "", "name", "ASC")
	// assert.Nil(suite.T(), err, "Entity list MUST NOT throw errors")
	// assert.Equal(suite.T(), 0, len(pages), "Pages must be 0 - %s found", len(pages))
	//assert.HTTPBodyContainsf(suite.T(), routes.Routes().ServeHTTP, "GET", "/auth", url.Values{}, "api key is required", "%")
}

// func (suite *CmsTestSuite) TestCmsPages() {
// 	db, err := mainDb("sqlite", "", "", "test_pages.db", "", "")
// 	assert.Nil(suite.T(), err, "DB error")
// 	defer db.Close()

// 	Init(Config{
// 		DbInstance:      db,
// 		EnableCache:     true,
// 		EnablePages:     true,
// 		EnableBlocks:    true,
// 		EnableSettings:  true,
// 		EnableSession:   true,
// 		EnableTemplates: true,
// 		// CustomEntityList: entityList(),
// 	})

// 	pages, err := EntityStore.EntityList("page", 0, 10, "", "name", "ASC")
// 	assert.Nil(suite.T(), err, "Entity list MUST NOT throw errors")
// 	assert.Equal(suite.T(), 0, len(pages), "Pages must be 0 - %s found", len(pages))
// }

func mainDb(driverName string, dbHost string, dbPort string, dbName string, dbUser string, dbPass string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	if driverName == "sqlite" {
		dsn := dbName
		db, err = sql.Open("sqlite", dsn)
	}
	if driverName == "mysql" {
		dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = sql.Open("mysql", dsn)
	}
	if driverName == "postgres" {
		dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=Europe/London"
		db, err = sql.Open("postgres", dsn)
	}
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, errors.New("database for driver " + driverName + " could not be intialized")
	}
	return db, nil
}
