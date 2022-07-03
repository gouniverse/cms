package cms

import (
	"database/sql"
	"errors"

	//"log"
	"time"

	"github.com/gouniverse/cachestore"
	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/logstore"
	"github.com/gouniverse/sessionstore"
	"github.com/gouniverse/settingstore"
)

// Cms defines the cms
type Cms struct {
	DbInstance *sql.DB
	// DbDriver           string
	// DbDsn              string
	customEntityList    []CustomEntityStructure
	blocksEnabled       bool
	menusEnabled        bool
	pagesEnabled        bool
	templatesEnabled    bool
	translationsEnabled bool
	widgetsEnabled      bool
	//AutomigrateEnabled bool
	CacheStore      *cachestore.Store
	EntityStore     *entitystore.Store
	LogStore        *logstore.Store
	SessionStore    *sessionstore.Store
	SettingStore    *settingstore.Store
	cacheEnabled    bool
	debug           bool
	prefix          string
	logsEnabled     bool
	sessionEnabled  bool
	settingsEnabled bool
}

// Cmsption defines an option for the CMS store
type CmsOption func(*Cms)

// WithAutoMigrate sets the table name for the cache store
// func WithAutoMigrate(automigrateEnabled bool) CmsOption {
// 	return func(cms *Cms) {
// 		cms.AutomigrateEnabled = automigrateEnabled
// 	}
// }

// WithDb sets the database for the CMS
func WithDb(db *sql.DB) CmsOption {
	return func(cms *Cms) {
		cms.DbInstance = db
		//cms.DbDriverName = cms.DriverName(cms.db)
	}
}

// WithDebug prints the SQL queries
func WithDebug(debug bool) CmsOption {
	return func(cms *Cms) {
		cms.debug = debug
	}
}

// WithBlocks enables blocks
func WithBlocks() CmsOption {
	return func(cms *Cms) {
		cms.blocksEnabled = true
	}
}

// WithBlocks enables blocks
func WithCache() CmsOption {
	return func(cms *Cms) {
		cms.cacheEnabled = true
	}
}

// WithBlocks enables blocks
func WithLogs() CmsOption {
	return func(cms *Cms) {
		cms.logsEnabled = true
	}
}

// WithMenus enables menus
func WithMenus() CmsOption {
	return func(cms *Cms) {
		cms.menusEnabled = true
	}
}

// WithPages enables pages
func WithPages() CmsOption {
	return func(cms *Cms) {
		cms.pagesEnabled = true
	}
}

// WithPrefix enables pages
func WithPrefix(prefix string) CmsOption {
	return func(cms *Cms) {
		cms.prefix = prefix
	}
}

// WithSession enables session
func WithSession() CmsOption {
	return func(cms *Cms) {
		cms.sessionEnabled = true
	}
}

// WithSettings enables settings
func WithSettings() CmsOption {
	return func(cms *Cms) {
		cms.settingsEnabled = true
	}
}

// WithTemplates enables templates
func WithTemplates() CmsOption {
	return func(cms *Cms) {
		cms.templatesEnabled = true
	}
}

// WithWidgets enables widgets
func WithWidgets() CmsOption {
	return func(cms *Cms) {
		cms.widgetsEnabled = true
	}
}

// WithCustomEntityList adds custom entities
func WithCustomEntityList(customEntityList []CustomEntityStructure) CmsOption {
	return func(cms *Cms) {
		cms.customEntityList = customEntityList
	}
}

// WithTableName sets the table name for the cache store
//func WithTableName(settingsTableName string) StoreOption {
//	return func(cms *Cms) {
//		cms.SettingsTableName = settingsTableName
//	}
//}

// NewCms creates a new CMS
func NewCms(opts ...CmsOption) (*Cms, error) {
	cms := &Cms{}
	for _, opt := range opts {
		opt(cms)
	}

	//if cms.settingsTableName == "" {
	//	return nil, error.New("Setting store: settingTableName is required")
	//}

	//if cms.automigrateEnabled {
	//	cms.AutoMigrate()
	//}

	if cms.DbInstance == nil {
		return nil, errors.New("DbInstance is required field")
	}

	if cms.prefix == "" {
		cms.prefix = "cms_"
	}

	var err error
	cms.EntityStore, err = entitystore.NewStore(entitystore.WithDb(cms.DbInstance), entitystore.WithEntityTableName("cms_entities_entity"), entitystore.WithAttributeTableName("cms_entities_attribute"), entitystore.WithAutoMigrate(true))

	if err != nil {
		// log.Panicln("Entity store failed to be intiated")
		return nil, err
	}

	if cms.cacheEnabled {
		cms.CacheStore, err = cachestore.NewStore(cachestore.WithDb(cms.DbInstance), cachestore.WithTableName("cms_cache"), cachestore.WithAutoMigrate(true))

		if err != nil {
			// log.Panicln("Cache store failed to be intiated")
			return nil, err
		}

		time.AfterFunc(3*time.Second, func() {
			go cms.CacheStore.ExpireCacheGoroutine()
		})
	}

	if cms.logsEnabled {
		cms.LogStore, err = logstore.NewStore(logstore.WithDb(cms.DbInstance), logstore.WithTableName("cms_log"), logstore.WithAutoMigrate(true))

		if err != nil {
			// log.Panicln("Log store failed to be intiated")
			return nil, err
		}
	}

	if cms.sessionEnabled {
		cms.SessionStore, err = sessionstore.NewStore(sessionstore.WithDb(cms.DbInstance), sessionstore.WithTableName("cms_session"), sessionstore.WithAutoMigrate(true))

		if err != nil {
			// log.Panicln("Session store failed to be intiated")
			return nil, err
		}

		time.AfterFunc(3*time.Second, func() {
			go cms.SessionStore.ExpireSessionGoroutine()
		})
	}

	if cms.settingsEnabled {
		cms.SettingStore, err = settingstore.NewStore(settingstore.WithDb(cms.DbInstance), settingstore.WithTableName("cms_settings"), settingstore.WithAutoMigrate(true))

		if err != nil {
			// log.Panicln("Setting store failed to be intiated")
			return nil, err
		}
	}

	// 	// Migrate the schema
	// 	// cms.DbInstance.AutoMigrate(&Entity{})
	// 	// cms.DbInstance.AutoMigrate(&EntityAttribute{})
	// 	configuration = config

	return cms, nil
}

// import (
// 	"database/sql"
// 	"log"
// 	"time"

// 	"github.com/gouniverse/cachestore"
// 	"github.com/gouniverse/entitystore"
// 	"github.com/gouniverse/logstore"
// 	"github.com/gouniverse/sessionstore"
// 	"github.com/gouniverse/settingstore"
// )

// var dbInstance *sql.DB
// var prefix string

// // Config contains the configurations for the auth package
// type Config struct {
// 	DbInstance         *sql.DB
// 	DbDriver           string
// 	DbDsn              string
// 	CustomEntityList   []CustomEntityStructure
// 	EnableBlocks       bool
// 	EnableCache        bool
// 	EnableLogs         bool
// 	EnableMenus        bool
// 	EnablePages        bool
// 	EnableSession      bool
// 	EnableSettings     bool
// 	EnableTemplates    bool
// 	EnableTranslations bool
// 	EnableWidgets      bool
// }

// var (
// 	configuration       Config
// 	EntityStore         *entitystore.Store
// 	CacheStore          *cachestore.Store
// 	LogStore            *logstore.Store
// 	SessionStore        *sessionstore.Store
// 	SettingStore        *settingstore.Store
// 	cacheEnabled        bool
// 	logsEnabled         bool
// 	sessionEnabled      bool
// 	settingsEnabled     bool
// 	translationsEnabled bool
// 	widgetsEnabled      bool
// )

// // Init initializes the CMS
// func Init(config Config) {
// 	if cms.DbInstance == nil && (cms.DbDriver == "" || cms.DbDsn == "") {
// 		log.Panicln("Either DbInstance or DnDriver and DbDsn are required field")
// 	}

// 	prefix = "cms_"

// 	var err error
// 	EntityStore, err = entitystore.NewStore(entitystore.WithDb(cms.DbInstance), entitystore.WithEntityTableName("cms_entities_entity"), entitystore.WithAttributeTableName("cms_entities_attribute"), entitystore.WithAutoMigrate(true))

// 	if err != nil {
// 		log.Panicln("Entity store failed to be intiated")
// 		return
// 	}

// 	if cms.EnableCache {
// 		cacheEnabled = true
// 		CacheStore, err = cachestore.NewStore(cachestore.WithDb(cms.DbInstance), cachestore.WithTableName("cms_cache"), cachestore.WithAutoMigrate(true))

// 		if err != nil {
// 			log.Panicln("Cache store failed to be intiated")
// 			return
// 		}

// 		time.AfterFunc(3*time.Second, func() {
// 			go CacheStore.ExpireCacheGoroutine()
// 		})
// 	}

// 	if cms.EnableLogs {
// 		logsEnabled = true
// 		LogStore, err = logstore.NewStore(logstore.WithDb(cms.DbInstance), logstore.WithTableName("cms_log"), logstore.WithAutoMigrate(true))

// 		if err != nil {
// 			log.Panicln("Log store failed to be intiated")
// 			return
// 		}
// 	}

// 	if cms.EnableSession {
// 		sessionEnabled = true
// 		SessionStore, err = sessionstore.NewStore(sessionstore.WithDb(cms.DbInstance), sessionstore.WithTableName("cms_session"), sessionstore.WithAutoMigrate(true))

// 		if err != nil {
// 			log.Panicln("Session store failed to be intiated")
// 			return
// 		}

// 		time.AfterFunc(3*time.Second, func() {
// 			go SessionStore.ExpireSessionGoroutine()
// 		})
// 	}

// 	if cms.EnableSettings {
// 		settingsEnabled = true
// 		SettingStore, err = settingstore.NewStore(settingstore.WithDb(cms.DbInstance), settingstore.WithTableName("cms_settings"), settingstore.WithAutoMigrate(true))

// 		if err != nil {
// 			log.Panicln("Setting store failed to be intiated")
// 			return
// 		}
// 	}

// 	// Migrate the schema
// 	// cms.DbInstance.AutoMigrate(&Entity{})
// 	// cms.DbInstance.AutoMigrate(&EntityAttribute{})
// 	configuration = config
// }

// // GetDb returns an instance to the CMS database
// func GetDb() *sql.DB {
// 	return configuration.DbInstance
// }

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
