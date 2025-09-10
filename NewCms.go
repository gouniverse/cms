package cms

import (
	"context"
	"errors"
	"time"

	"github.com/dracory/cachestore"
	"github.com/dracory/entitystore"
	"github.com/dracory/logstore"
	"github.com/dracory/sessionstore"
	"github.com/dracory/settingstore"
	"github.com/dracory/taskstore"
)

// NewCms creates a new CMS instance based on the given configuration
//
// Parameters:
// - config Config - the CMS configuration
//
// Returns:
// - *Cms - the new CMS instance
// - error - any error if occurred, nil otherwise
func NewCms(config Config) (cms *Cms, err error) {

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

	cms = configToCms(config)

	err = cmsEntitiesSetup(cms)

	if err != nil {
		return nil, err
	}

	err = cmsCacheSetup(cms)

	if err != nil {
		return nil, err
	}

	err = cmsLogsSetup(cms)

	if err != nil {
		return nil, err
	}

	err = cmsSessionSetup(cms)

	if err != nil {
		return nil, err
	}

	err = cmsSettingsSetup(cms)

	if err != nil {
		return nil, err
	}

	err = cmsTasksSetup(cms)

	if err != nil {
		return nil, err
	}

	err = cmsUsersSetup(cms)

	if err != nil {
		return nil, err
	}

	return cms, nil
}

// cmsCacheSetup sets up the cache
func cmsCacheSetup(cms *Cms) (err error) {
	if !cms.cacheEnabled {
		return nil
	}

	cms.CacheStore, err = cachestore.NewStore(cachestore.NewStoreOptions{
		DB:             cms.Database.DB(),
		CacheTableName: cms.cacheTableName,
	})

	if err != nil {
		return err
	}

	if cms.cacheAutoMigrate {
		err = cms.CacheStore.AutoMigrate()
		if err != nil {
			return err
		}
	}

	time.AfterFunc(3*time.Second, func() {
		go cms.CacheStore.ExpireCacheGoroutine()
	})

	return nil
}

// cmsEntitiesSetup sets up the entities
func cmsEntitiesSetup(cms *Cms) (err error) {
	cms.EntityStore, err = entitystore.NewStore(entitystore.NewStoreOptions{
		Database:           cms.Database,
		EntityTableName:    cms.entityTableName,
		AttributeTableName: cms.attributeTableName,
	})

	if err != nil {
		return err
	}

	if cms.entitiesAutoMigrate {
		err = cms.EntityStore.AutoMigrate()
		if err != nil {
			return err
		}
	}

	return nil
}

// cmsLogsSetup sets up the logs
func cmsLogsSetup(cms *Cms) (err error) {
	if !cms.logsEnabled {
		return nil
	}

	cms.LogStore, err = logstore.NewStore(logstore.NewStoreOptions{
		DB:           cms.Database.DB(),
		LogTableName: cms.logTableName,
	})

	if err != nil {
		return err
	}

	if cms.logsAutomigrate {
		err = cms.LogStore.AutoMigrate()
		if err != nil {
			return err
		}
	}

	return nil
}

// cmsSessionSetup sets up the session
func cmsSessionSetup(cms *Cms) (err error) {
	if !cms.sessionEnabled {
		return nil
	}

	cms.SessionStore, err = sessionstore.NewStore(sessionstore.NewStoreOptions{
		DB:               cms.Database.DB(),
		SessionTableName: cms.sessionTableName,
	})

	if err != nil {
		return err
	}

	if cms.sessionAutomigrate {
		err = cms.SessionStore.AutoMigrate(context.Background())
		if err != nil {
			return err
		}
	}

	time.AfterFunc(3*time.Second, func() {
		go cms.SessionStore.SessionExpiryGoroutine()
	})

	return nil
}

// cmsSettingsSetup sets up the settings
func cmsSettingsSetup(cms *Cms) (err error) {
	if !cms.settingsEnabled {
		return nil
	}

	cms.SettingStore, err = settingstore.NewStore(settingstore.NewStoreOptions{
		SettingTableName:   cms.settingsTableName,
		DB:                 cms.Database.DB(),
		AutomigrateEnabled: true,
	})

	if err != nil {
		return err
	}

	if cms.settingsAutomigrate {
		err = cms.SettingStore.AutoMigrate(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}

// cmsTasksSetup sets up the tasks
func cmsTasksSetup(cms *Cms) (err error) {
	if !cms.tasksEnabled {
		return nil
	}

	cms.TaskStore, err = taskstore.NewStore(taskstore.NewStoreOptions{
		DB:             cms.Database.DB(),
		TaskTableName:  cms.tasksTaskTableName,
		QueueTableName: cms.tasksQueueTableName,
	})

	if err != nil {
		return err
	}

	if cms.tasksAutomigrate {
		err = cms.TaskStore.AutoMigrate()
		if err != nil {
			return err
		}
	}

	return nil
}

// cmsUsersSetup sets up the users
func cmsUsersSetup(cms *Cms) (err error) {
	if !cms.usersEnabled {
		return nil
	}

	cms.UserStore, err = entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                 cms.Database.DB(),
		EntityTableName:    cms.userAttributeTableName,
		AttributeTableName: cms.userAttributeTableName,
	})

	if err != nil {
		return err
	}

	if cms.usersAutoMigrate {
		err = cms.UserStore.AutoMigrate()
		if err != nil {
			return err
		}
	}

	return nil
}
