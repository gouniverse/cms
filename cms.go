package cms

import (
	"log"

	"github.com/gouniverse/cachestore"
	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/settingstore"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB
var prefix string

// Config contains the configurations for the auth package
type Config struct {
	DbInstance         *gorm.DB
	DbDriver           string
	DbDsn              string
	CustomEntityList   []CustomEntityStructure
	EnableBlocks       bool
	EnableCache        bool
	EnableMenus        bool
	EnablePages        bool
	EnableSettings     bool
	EnableTemplates    bool
	EnableTranslations bool
	EnableWidgets      bool
}

var (
	configuration       Config
	entityStore         *entitystore.Store
	settingStore        *settingstore.Store
	cacheStore          *cachestore.Store
	cacheEnabled        bool
	settingsEnabled     bool
	translationsEnabled bool
	widgetsEnabled      bool
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

	entityStore = entitystore.NewStore(entitystore.WithGormDb(config.DbInstance), entitystore.WithEntityTableName("cms_entities_entity"), entitystore.WithAttributeTableName("cms_entities_attribute"), entitystore.WithAutoMigrate(true))

	if config.EnableSettings {
		settingsEnabled = true
		settingStore = settingstore.NewStore(settingstore.WithGormDb(config.DbInstance), settingstore.WithTableName("cms_settings"), settingstore.WithAutoMigrate(true))
	}

	if config.EnableCache {
		cacheEnabled = true
		cacheStore = cachestore.NewStore(cachestore.WithGormDb(config.DbInstance), cachestore.WithTableName("cms_cache"), cachestore.WithAutoMigrate(true))
		go cacheStore.ExpireCacheGoroutine()
	}

	// Migrate the schema
	// config.DbInstance.AutoMigrate(&Entity{})
	// config.DbInstance.AutoMigrate(&EntityAttribute{})
	configuration = config
}

// GetDb returns an instance to the CMS database
func GetDb() *gorm.DB {
	return configuration.DbInstance
}

// GetEntityStore returns the entity store
func GetEntityStore() *entitystore.Store {
	return entityStore
}

// GetSettingStore returns the setting store
func GetSettingStore() *settingstore.Store {
	if settingsEnabled {
		return settingStore
	}

	return nil
}

// GetCacheStore returns the cache store
func GetCacheStore() *cachestore.Store {
	if cacheEnabled {
		return cacheStore
	}

	return nil
}

type CustomEntityStructure struct {
	// Type of the entity
	Type string
	// Label to display referencing the entity
	TypeLabel string
	// Name of the entity
	Name string
	// AttributeList list of attributes
	AttributeList []CustomAttributeStructure
	// Group to which this entity belongs (i.e. Shop, Users, etc)
	Group string
}

type CustomAttributeStructure struct {
	// Name the name of the attribute
	Name string
	// Type of the attribute - string, float, int
	Type string
	// FormControlLabel label to display for the control
	FormControlLabel string
	// FormControlType the type of form control - input, textarea. etc
	FormControlType string
	// FormControlHelp help message to display for the control
	FormControlHelp string
	// BelongsToType describes a Belong To relationsip
	BelongsToType string
	// HasOneType describes a Has One relationsip
	HasOneType string
	// HasManyType describes a Has Many relationsip
	HasManyType string
}
