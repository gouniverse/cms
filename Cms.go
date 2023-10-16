package cms

import (
	"net/http"

	"github.com/gouniverse/cachestore"
	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/logstore"
	"github.com/gouniverse/sessionstore"
	"github.com/gouniverse/settingstore"
	sqldb "github.com/gouniverse/sql"
	"github.com/gouniverse/taskstore"
	"github.com/samber/lo"
)

type LanguageKey struct{}

// Cms defines the cms
type Cms struct {
	Database *sqldb.Database
	// DbInstance   *sql.DB
	CacheStore   *cachestore.Store
	EntityStore  *entitystore.Store
	LogStore     *logstore.Store
	SessionStore *sessionstore.Store
	SettingStore *settingstore.Store
	TaskStore    *taskstore.Store
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

	tasksEnabled        bool
	tasksAutomigrate    bool
	tasksTaskTableName  string
	tasksQueueTableName string

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
	cms.tasksAutomigrate = config.TasksAutomigrate
	cms.tasksEnabled = config.TasksEnable
	cms.tasksQueueTableName = config.TasksQueueTableName
	cms.tasksTaskTableName = config.TasksTaskTableName
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
	cms.tasksQueueTableName = cms.prefix + "tasks_queue"
	cms.tasksTaskTableName = cms.prefix + "tasks_task"
	cms.userEntityTableName = cms.prefix + "users_entity"
	cms.userAttributeTableName = cms.prefix + "users_attribute"

	return cms
}

func (c *Cms) layout(content string) string {
	return content
}

// SetFuncLayout sets a layout for the CMS to display inside
func (c *Cms) SetFuncLayout(funcLayout func(content string) string) *Cms {
	c.funcLayout = funcLayout
	return c
}
