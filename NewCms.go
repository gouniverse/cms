package cms

import (
	"errors"
	"time"

	"github.com/gouniverse/cachestore"
	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/logstore"
	"github.com/gouniverse/sessionstore"
	"github.com/gouniverse/settingstore"
	"github.com/gouniverse/taskstore"
)

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

	if cms.tasksEnabled {
		cms.TaskStore, err = taskstore.NewStore(taskstore.NewStoreOptions{
			DB:             cms.Database.DB(),
			TaskTableName:  cms.tasksTaskTableName,
			QueueTableName: cms.tasksQueueTableName,
		})

		if err != nil {
			return nil, err
		}

		if cms.tasksAutomigrate {
			err = cms.TaskStore.AutoMigrate()
			if err != nil {
				return nil, err
			}
		}

	}

	return cms, nil
}
