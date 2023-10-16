package cms

import (
	"database/sql"
	"net/http"

	sqldb "github.com/gouniverse/sql"
)

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
	TasksEnable                bool
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
