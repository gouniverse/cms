package cms

// import "database/sql"

// // Cmsption defines an option for the CMS store
// type CmsOption func(*Cms)

// // WithDb sets the database for the CMS
// func WithDb(db *sql.DB) CmsOption {
// 	return func(cms *Cms) {
// 		cms.DbInstance = db
// 		//cms.DbDriverName = cms.DriverName(cms.db)
// 	}
// }

// // WithDebug prints the SQL queries
// // func WithDebug(debug bool) CmsOption {
// // 	return func(cms *Cms) {
// // 		cms.debug = debug
// // 	}
// // }

// // WithBlocks enables blocks
// func WithBlocks() CmsOption {
// 	return func(cms *Cms) {
// 		cms.blocksEnabled = true
// 	}
// }

// // WithBlocks enables blocks
// func WithCache() CmsOption {
// 	return func(cms *Cms) {
// 		cms.cacheEnabled = true
// 	}
// }

// // WithBlocks enables blocks
// func WithLogs() CmsOption {
// 	return func(cms *Cms) {
// 		cms.logsEnabled = true
// 	}
// }

// // WithMenus enables menus
// func WithMenus() CmsOption {
// 	return func(cms *Cms) {
// 		cms.menusEnabled = true
// 	}
// }

// // WithPages enables pages
// func WithPages() CmsOption {
// 	return func(cms *Cms) {
// 		cms.pagesEnabled = true
// 	}
// }

// // WithPrefix enables pages
// func WithPrefix(prefix string) CmsOption {
// 	return func(cms *Cms) {
// 		cms.prefix = prefix
// 	}
// }

// // WithSession enables session
// func WithSession() CmsOption {
// 	return func(cms *Cms) {
// 		cms.sessionEnabled = true
// 	}
// }

// // WithSettings enables settings
// func WithSettings() CmsOption {
// 	return func(cms *Cms) {
// 		cms.settingsEnabled = true
// 	}
// }

// // WithTemplates enables templates
// func WithTemplates() CmsOption {
// 	return func(cms *Cms) {
// 		cms.templatesEnabled = true
// 	}
// }

// // WithWidgets enables widgets
// func WithWidgets() CmsOption {
// 	return func(cms *Cms) {
// 		cms.widgetsEnabled = true
// 	}
// }

// // WithCustomEntityList adds custom entities
// func WithCustomEntityList(customEntityList []CustomEntityStructure) CmsOption {
// 	return func(cms *Cms) {
// 		cms.customEntityList = customEntityList
// 	}
// }
