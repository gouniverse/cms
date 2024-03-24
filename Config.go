package cms

import (
	"database/sql"

	"github.com/gouniverse/sb"
)

type Config struct {
	Database                   *sb.Database
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
	Shortcodes                 []ShortcodeInterface
	TasksEnable                bool
	TasksAutomigrate           bool
	TasksQueueTableName        string
	TasksTaskTableName         string
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
