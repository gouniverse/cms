package cms

import (
	"regexp"
	"strings"

	"github.com/gouniverse/entitystore"
)

// PageFindByAlias helper method to find a page by alias
//
// =====================================================================
//  1. It will attempt to find the page by the provided alias exactly
//     as provided
//  2. It will attempt to find the page with the alias prefixed with "/"
//     in case of error
//
// =====================================================================
func (cms *Cms) PageFindByAlias(alias string) (*entitystore.Entity, error) {
	// Try to find by "alias"
	page, err := cms.EntityStore.EntityFindByAttribute(ENTITY_TYPE_PAGE, "alias", ""+alias+"")

	if err != nil {
		return nil, err
	}

	if page != nil {
		return page, nil
	}

	// Try to find by "/alias"
	page, err = cms.EntityStore.EntityFindByAttribute(ENTITY_TYPE_PAGE, "alias", "/"+alias+"")

	if err != nil {
		return nil, err
	}

	if page != nil {
		return page, nil
	}

	page, err = cms.PageFindByAliasWithPatterns(alias)

	if err != nil {
		return nil, err
	}

	if page != nil {
		return page, nil
	}

	return nil, nil
}

// PageFindByAliasWithPatterns helper method to find a page by matching patterns
//
// =====================================================================
//
//	The following patterns are supported:
//	:any
//	:num
//	:all
//	:string
//	:number
//	:numeric
//	:alpha
//
// =====================================================================
func (cms *Cms) PageFindByAliasWithPatterns(alias string) (*entitystore.Entity, error) {
	patterns := map[string]string{
		":any":     "([^/]+)",
		":num":     "([0-9]+)",
		":all":     "(.*)",
		":string":  "([a-zA-Z]+)",
		":number":  "([0-9]+)",
		":numeric": "([0-9-.]+)",
		":alpha":   "([a-zA-Z0-9-_]+)",
	}

	attributes, err := cms.EntityStore.AttributeList(entitystore.AttributeQueryOptions{
		EntityType:   ENTITY_TYPE_PAGE,
		AttributeKey: "alias",
	})
	if err != nil {
		return nil, err
	}

	pageAliasMap := make(map[string]string, len(attributes))
	for _, attribute := range attributes {
		pageAliasMap[attribute.EntityID()] = attribute.AttributeValue()
	}

	for pageID, pageAlias := range pageAliasMap {
		if !strings.Contains(pageAlias, ":") {
			continue
		}

		for pattern, replacement := range patterns {
			pageAlias = strings.ReplaceAll(pageAlias, pattern, replacement)
		}

		matcher := regexp.MustCompile("^" + pageAlias + "$")
		if matcher.MatchString(alias) {
			return cms.EntityStore.EntityFindByID(pageID)
		}
	}

	return nil, nil
}
