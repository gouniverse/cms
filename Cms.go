package cms

import (
	"github.com/dracory/blockeditor"
	"github.com/dracory/cachestore"
	"github.com/dracory/entitystore"
	"github.com/dracory/logstore"
	"github.com/dracory/sb"
	"github.com/dracory/sessionstore"
	"github.com/dracory/settingstore"
	"github.com/dracory/taskstore"
	"github.com/dracory/ui"
	"github.com/samber/lo"
)

type LanguageKey struct{}

// Cms defines the cms
type Cms struct {
	Database     sb.DatabaseInterface
	CacheStore   cachestore.StoreInterface
	EntityStore  entitystore.StoreInterface
	LogStore     logstore.StoreInterface
	SessionStore sessionstore.StoreInterface
	SettingStore settingstore.StoreInterface
	TaskStore    taskstore.StoreInterface
	UserStore    entitystore.StoreInterface

	entitiesAutoMigrate bool
	entityTableName     string
	attributeTableName  string

	blocksEnabled bool

	blockEditorDefinitions []blockeditor.BlockDefinition
	blockEditorRenderer    func(blocks []ui.BlockInterface) string

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

	shortcodes []ShortcodeInterface

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
		config.Database = sb.NewDatabase(config.DbInstance, sb.DatabaseDriverName(config.DbInstance))
	}

	if config.Database == nil && (config.DbDriver != "" && config.DbDsn != "") {
		var errDatabase error
		config.Database, errDatabase = sb.NewDatabaseFromDriver(config.DbDriver, config.DbDsn)
		if errDatabase != nil {
			panic("At CMS: " + errDatabase.Error())
		}
	}

	if config.Shortcodes == nil {
		config.Shortcodes = []ShortcodeInterface{}
	}

	if config.TranslationLanguageDefault == "" && len(config.TranslationLanguages) > 0 {
		config.TranslationLanguageDefault = config.TranslationLanguages[lo.Keys(config.TranslationLanguages)[0]]
	}

	cms.blocksEnabled = config.BlocksEnable
	cms.blockEditorDefinitions = config.BlockEditorDefinitions
	cms.blockEditorRenderer = config.BlockEditorRenderer
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

func (c *Cms) logErrorWithContext(message string, context any) {
	if c.LogStore == nil {
		return
	}

	c.LogStore.ErrorWithContext(message, context)
}
