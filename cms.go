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

type Config struct {
	DbInstance          *sql.DB
	DbDriver            string
	DbDsn               string
	CustomEntityList    []CustomEntityStructure
	Prefix              string
	BlocksEnable        bool
	CacheAutomigrate    bool
	CacheEnable         bool
	EntitiesAutomigrate bool
	LogsEnable          bool
	LogsAutomigrate     bool
	MenusEnable         bool
	PagesEnable         bool
	SessionAutomigrate  bool
	SessionEnable       bool
	SettingsAutomigate  bool
	SettingsEnable      bool
	TemplatesEnable     bool
	TranslationsEnable  bool
	WidgetsEnable       bool
}

// Cms defines the cms
type Cms struct {
	DbInstance   *sql.DB
	CacheStore   *cachestore.Store
	EntityStore  *entitystore.Store
	LogStore     *logstore.Store
	SessionStore *sessionstore.Store
	SettingStore *settingstore.Store

	entitiesAutoMigrate bool
	entityTableName     string
	attributeTableName  string

	blocksEnabled bool

	cacheAutoMigrate bool
	cacheEnabled     bool
	cacheTableName   string

	customEntityList []CustomEntityStructure
	// debug               bool

	menusEnabled        bool
	pagesEnabled        bool
	templatesEnabled    bool
	translationsEnabled bool
	widgetsEnabled      bool

	prefix string

	logsEnabled     bool
	logsAutomigrate bool
	logTableName    string

	sessionEnabled     bool
	sessionAutomigrate bool
	sessionTableName   string

	settingsEnabled     bool
	settingsAutomigrate bool
	settingsTableName   string
}

func configToCms(config Config) *Cms {

	cms := &Cms{}

	if config.Prefix == "" {
		cms.prefix = "cms_"
	}

	cms.blocksEnabled = config.BlocksEnable
	cms.cacheAutoMigrate = config.CacheAutomigrate
	cms.cacheEnabled = config.CacheEnable
	cms.customEntityList = config.CustomEntityList
	cms.entitiesAutoMigrate = config.EntitiesAutomigrate
	cms.logsAutomigrate = config.LogsAutomigrate
	cms.logsEnabled = config.LogsEnable
	cms.menusEnabled = config.MenusEnable
	cms.pagesEnabled = config.PagesEnable
	cms.sessionAutomigrate = config.SessionAutomigrate
	cms.sessionEnabled = config.SessionEnable
	cms.settingsAutomigrate = config.SettingsAutomigate
	cms.settingsEnabled = config.SettingsEnable
	cms.templatesEnabled = config.TemplatesEnable
	cms.translationsEnabled = config.TranslationsEnable
	cms.widgetsEnabled = config.WidgetsEnable
	cms.DbInstance = config.DbInstance
	cms.prefix = config.Prefix

	// Table Names
	cms.attributeTableName = cms.prefix + "entities_attribute"
	cms.cacheTableName = cms.prefix + "cache"
	cms.entityTableName = cms.prefix + "entities_entity"
	cms.logTableName = cms.prefix + "log"
	cms.sessionTableName = cms.prefix + "session"
	cms.settingsTableName = cms.prefix + "settings"

	return cms
}

// NewCms creates a new CMS
func NewCms(config Config) (*Cms, error) {

	if config.DbInstance == nil {
		return nil, errors.New("DbInstance is required field")
	}

	cms := configToCms(config)

	var err error
	cms.EntityStore, err = entitystore.NewStore(entitystore.WithDb(cms.DbInstance), entitystore.WithEntityTableName(cms.entityTableName), entitystore.WithAttributeTableName(cms.attributeTableName))

	if err != nil {
		return nil, err
	}

	if cms.entitiesAutoMigrate {
		err = cms.EntityStore.AutoMigrate()
		if err != nil {
			return nil, err
		}
	}

	if cms.cacheEnabled {
		cms.CacheStore, err = cachestore.NewStore(cachestore.WithDb(cms.DbInstance), cachestore.WithTableName(cms.cacheTableName))

		if err != nil {
			// log.Panicln("Cache store failed to be intiated")
			return nil, err
		}

		if cms.cacheAutoMigrate {
			err = cms.CacheStore.AutoMigrate()
			if err != nil {
				return nil, err
			}
		}

		time.AfterFunc(3*time.Second, func() {
			go cms.CacheStore.ExpireCacheGoroutine()
		})
	}

	if cms.logsEnabled {
		cms.LogStore, err = logstore.NewStore(logstore.WithDb(cms.DbInstance), logstore.WithTableName(cms.logTableName))

		if err != nil {
			return nil, err
		}

		if cms.cacheAutoMigrate {
			err = cms.LogStore.AutoMigrate()
			if err != nil {
				return nil, err
			}
		}

	}

	if cms.sessionEnabled {
		cms.SessionStore, err = sessionstore.NewStore(sessionstore.WithDb(cms.DbInstance), sessionstore.WithTableName(cms.sessionTableName))

		if err != nil {
			// log.Panicln("Session store failed to be intiated")
			return nil, err
		}

		if cms.sessionAutomigrate {
			err = cms.SessionStore.AutoMigrate()
			if err != nil {
				return nil, err
			}
		}

		time.AfterFunc(3*time.Second, func() {
			go cms.SessionStore.ExpireSessionGoroutine()
		})
	}

	if cms.settingsEnabled {
		cms.SettingStore, err = settingstore.NewStore(settingstore.WithDb(cms.DbInstance), settingstore.WithTableName(cms.settingsTableName), settingstore.WithAutoMigrate(true))

		if err != nil {
			// log.Panicln("Setting store failed to be intiated")
			return nil, err
		}

		if cms.settingsAutomigrate {
			err = cms.SettingStore.AutoMigrate()
			if err != nil {
				return nil, err
			}
		}
	}

	return cms, nil
}

// func NewCms(opts ...CmsOption) (*Cms, error) {
// 	cms := &Cms{}
// 	for _, opt := range opts {
// 		opt(cms)
// 	}
// }

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
