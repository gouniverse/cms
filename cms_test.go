package cms

import (
	"log"
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

type CmsTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *CmsTestSuite) SetupTest() {
	suite.VariableThatShouldStartAtFive = 5

	// Setup()
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *CmsTestSuite) TestExample() {
	assert.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)
}

//TestAuth tests the auth page
func (suite *CmsTestSuite) TestCmsInit() {
	_, err := mainDb("sqlite", "", "", "./test_init.db", "", "")
	assert.Nil(suite.T(), err, "DB error %s", err.Error())
	// Init(Config{
	// 	DbInstance:      db,
	// 	EnableTemplates: true,
	// 	EnablePages:     true,
	// 	EnableBlocks:    true,
	// 	EnableSettings:  true,
	// 	// CustomEntityList: entityList(),
	// })
	// pages, err := EntityStore.EntityList("page", 0, 10, "", "name", "ASC")
	// assert.NotNil(suite.T(), err, "DB error %s", err.Error())
	// assert.Equal(suite.T(), 0, len(pages), "Pages must be 0 - %s found", len(pages))
	//assert.True(suite.T(), cms.Config.EnableBlocks, "DB error %s", err.Error())
	//assert.HTTPBodyContainsf(suite.T(), routes.Routes().ServeHTTP, "GET", "/auth", url.Values{}, "api key is required", "%")
}

func mainDb(driverName string, dbHost string, dbPort string, dbName string, dbUser string, dbPass string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	if driverName == "sqlite" {
		dsn := dbName
		log.Println(dsn)
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

func TestCmsTestSuite(t *testing.T) {
	suite.Run(t, new(CmsTestSuite))
}
