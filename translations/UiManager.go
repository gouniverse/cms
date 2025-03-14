package cms

import (
	"github.com/gouniverse/bs"
	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/hb"
)

type Config struct {
	Endpoint                           string
	EntityStore                        entitystore.StoreInterface
	TranslationEntityType              string
	PathTranslationsTranslationManager string
	PathTranslationsTranslationUpdate  string
	TranslationLanguages               map[string]string
	TranslationLanguageDefault         string
	WebpageComplete                    func(string, string) *hb.HtmlWebpage
	FuncLayout                         func(string) string
	CmsHeader                          func(string) string
	CmsBreadcrumbs                     func([]bs.Breadcrumb) string
}

func NewUiManager(config Config) UiManager {
	return UiManager{
		endpoint:                           config.Endpoint,
		entityStore:                        config.EntityStore,
		translationEntityType:              config.TranslationEntityType,
		pathTranslationsTranslationManager: config.PathTranslationsTranslationManager,
		pathTranslationsTranslationUpdate:  config.PathTranslationsTranslationUpdate,
		translationLanguages:               config.TranslationLanguages,
		translationLanguageDefault:         config.TranslationLanguageDefault,
		webpageComplete:                    config.WebpageComplete,
		funcLayout:                         config.FuncLayout,
		cmsHeader:                          config.CmsHeader,
		cmsBreadcrumbs:                     config.CmsBreadcrumbs,
	}
}

type UiManager struct {
	endpoint                           string
	entityStore                        entitystore.StoreInterface
	translationEntityType              string
	pathTranslationsTranslationManager string
	pathTranslationsTranslationUpdate  string
	translationLanguages               map[string]string
	translationLanguageDefault         string
	webpageComplete                    func(string, string) *hb.HtmlWebpage
	funcLayout                         func(string) string
	cmsHeader                          func(string) string
	cmsBreadcrumbs                     func([]bs.Breadcrumb) string
}
