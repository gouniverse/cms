package cms

import (
	"net/http"
	"strings"

	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

// FrontendHandler is the main handler for the CMS frontend.
//
// It handles the routing of the request to the appropriate page.
//
// If the URI ends with ".ico", it will return a blank response, as the browsers
// (at least Chrome and Firefox) will always request the favicon even if
// it's not present in the HTML.
func (cms *Cms) FrontendHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(cms.FrontendHandlerRenderAsString(w, r)))
}

// FrontendHandlerRenderAsString is the same as FrontendHandler but returns a string
// instead of writing to the http.ResponseWriter.
//
// It handles the routing of the request to the appropriate page.
//
// If the URI ends with ".ico", it will return a blank response, as the browsers
// (at least Chrome and Firefox) will always request the favicon even if
// it's not present in the HTML.
//
// If the translations are enabled, it will use the language from the request context.
// If the language is not valid, it will use the default language for the translations.
func (cms *Cms) FrontendHandlerRenderAsString(w http.ResponseWriter, r *http.Request) string {
	uri := r.RequestURI

	if strings.HasSuffix(uri, ".ico") {
		return ""
	}

	languageAny := r.Context().Value(LanguageKey{})
	language := utils.ToString(languageAny)

	if cms.translationsEnabled {
		isValidLanguage := lo.Contains(lo.Keys(cms.translationLanguages), language)

		if !isValidLanguage {
			language = cms.translationLanguageDefault
		}
	}

	return cms.PageRenderHtmlByAlias(r, r.URL.Path, language)
}
