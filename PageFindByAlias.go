package cms

import (
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

	return nil, nil
}
