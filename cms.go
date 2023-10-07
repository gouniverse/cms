package cms

import (
	"database/sql"
	"errors"
	"net/http"

	"time"

	"github.com/gouniverse/cachestore"
	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/logstore"
	"github.com/gouniverse/sessionstore"
	"github.com/gouniverse/settingstore"
	sqldb "github.com/gouniverse/sql"
	"github.com/samber/lo"
)

type LanguageKey struct{}

type Config struct {
	Database                   *sqldb.Database
	DbInstance                 *sql.DB
	DbDriver                   string
	DbDsn                      string
	CustomEntityList           []CustomEntityStructure
	Prefix                     string
	BlocksEnable               bool
	CacheAutomigrate           bool
	CacheEnable                bool
	EntitiesAutomigrate        bool
	LogsEnable                 bool
	LogsAutomigrate            bool
	MenusEnable                bool
	PagesEnable                bool
	SessionAutomigrate         bool
	SessionEnable              bool
	SettingsAutomigrate        bool
	SettingsEnable             bool
	Shortcodes                 map[string]func(*http.Request, string, map[string]string) string
	TemplatesEnable            bool
	TranslationsEnable         bool
	TranslationLanguageDefault string
	TranslationLanguages       map[string]string
	UsersEnable                bool
	UsersAutomigrate           bool
	DashboardEnable            bool
	WidgetsEnable              bool
	FuncLayout                 func(content string) string
}

// Cms defines the cms
type Cms struct {
	Database *sqldb.Database
	// DbInstance   *sql.DB
	CacheStore   *cachestore.Store
	EntityStore  *entitystore.Store
	LogStore     *logstore.Store
	SessionStore *sessionstore.Store
	SettingStore *settingstore.Store
	UserStore    *entitystore.Store

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

	shortcodes map[string]func(*http.Request, string, map[string]string) string

	translationLanguageDefault string
	translationLanguages       map[string]string

	usersEnabled           bool
	usersAutoMigrate       bool
	userEntityTableName    string
	userAttributeTableName string

	funcLayout func(content string) string
}

func configToCms(config Config) *Cms {

	cms := &Cms{}

	if config.Prefix == "" {
		cms.prefix = "cms_"
	}

	if config.FuncLayout == nil {
		config.FuncLayout = cms.layout
	}

	if config.Database == nil && config.DbInstance != nil {
		config.Database = sqldb.NewDatabase(config.DbInstance, sqldb.DatabaseDriverName(config.DbInstance))
	}

	if config.Database == nil && (config.DbDriver != "" && config.DbDsn != "") {
		var errDatabase error
		config.Database, errDatabase = sqldb.NewDatabaseFromDriver(config.DbDriver, config.DbDsn)
		if errDatabase != nil {
			panic("At CMS: " + errDatabase.Error())
		}
	}

	if config.Shortcodes == nil {
		config.Shortcodes = map[string]func(*http.Request, string, map[string]string) string{}
	}

	if config.TranslationLanguageDefault == "" && len(config.TranslationLanguages) > 0 {
		config.TranslationLanguageDefault = config.TranslationLanguages[lo.Keys(config.TranslationLanguages)[0]]
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
	cms.settingsAutomigrate = config.SettingsAutomigrate
	cms.settingsEnabled = config.SettingsEnable
	cms.shortcodes = config.Shortcodes
	cms.templatesEnabled = config.TemplatesEnable
	cms.translationsEnabled = config.TranslationsEnable
	cms.translationLanguageDefault = config.TranslationLanguageDefault
	cms.translationLanguages = config.TranslationLanguages
	cms.widgetsEnabled = config.WidgetsEnable
	cms.usersEnabled = config.UsersEnable
	cms.usersAutoMigrate = config.UsersAutomigrate
	cms.Database = config.Database
	// cms.DbInstance = config.DbInstance
	cms.prefix = config.Prefix
	cms.funcLayout = config.FuncLayout

	// Table Names
	cms.attributeTableName = cms.prefix + "entities_attribute"
	cms.cacheTableName = cms.prefix + "cache"
	cms.entityTableName = cms.prefix + "entities_entity"
	cms.logTableName = cms.prefix + "log"
	cms.sessionTableName = cms.prefix + "session"
	cms.settingsTableName = cms.prefix + "setting"
	cms.userEntityTableName = cms.prefix + "users_entity"
	cms.userAttributeTableName = cms.prefix + "users_attribute"

	return cms
}

// NewCms creates a new CMS
func NewCms(config Config) (*Cms, error) {

	if config.DbInstance == nil && config.Database == nil && (config.DbDriver == "" || config.DbDsn == "") {
		return nil, errors.New("database (preferred) OR db instance OR (driver & dsn) are required field")
	}

	if config.DbInstance != nil && config.Database != nil && (config.DbDriver != "" && config.DbDsn != "") {
		return nil, errors.New("only one of database (preferred) OR db instance OR (driver & dsn) are required field")
	}

	if config.TranslationsEnable && len(config.TranslationLanguages) < 1 {
		return nil, errors.New("translations enabled but no translation languages specified")
	}

	if config.TranslationsEnable && len(config.TranslationLanguageDefault) < 1 {
		return nil, errors.New("translations enabled but no default translation language specified")
	}

	cms := configToCms(config)

	var err error
	cms.EntityStore, err = entitystore.NewStore(entitystore.NewStoreOptions{
		Database:           cms.Database,
		EntityTableName:    cms.entityTableName,
		AttributeTableName: cms.attributeTableName,
	})

	if err != nil {
		return nil, err
	}

	if cms.entitiesAutoMigrate {
		err = cms.EntityStore.AutoMigrate()
		if err != nil {
			return nil, err
		}
	}

	if cms.usersEnabled {
		cms.UserStore, err = entitystore.NewStore(entitystore.NewStoreOptions{
			DB:                 cms.Database.DB(),
			EntityTableName:    cms.userAttributeTableName,
			AttributeTableName: cms.userAttributeTableName,
		})

		if err != nil {
			return nil, err
		}

		if cms.usersAutoMigrate {
			err = cms.UserStore.AutoMigrate()
			if err != nil {
				return nil, err
			}
		}
	}

	if cms.cacheEnabled {
		cms.CacheStore, err = cachestore.NewStore(cachestore.NewStoreOptions{
			DB:             cms.Database.DB(),
			CacheTableName: cms.cacheTableName,
		})

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
		cms.LogStore, err = logstore.NewStore(logstore.NewStoreOptions{
			DB:           cms.Database.DB(),
			LogTableName: cms.logTableName,
		})

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
		cms.SessionStore, err = sessionstore.NewStore(sessionstore.NewStoreOptions{
			DB:               cms.Database.DB(),
			SessionTableName: cms.sessionTableName,
		})

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
		cms.SettingStore, err = settingstore.NewStore(settingstore.WithDb(cms.Database.DB()), settingstore.WithTableName(cms.settingsTableName), settingstore.WithAutoMigrate(true))

		if err != nil {
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

func (c *Cms) layout(content string) string {
	return content
}

// SetFuncLayout sets a layout for the CMS to display inside
func (c *Cms) SetFuncLayout(funcLayout func(content string) string) *Cms {
	c.funcLayout = funcLayout
	return c
}
