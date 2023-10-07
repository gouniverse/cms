package cms

import (
	"log"

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
func (cms *Cms) PageFindByAlias(alias string) *entitystore.Entity {
	// Try to find by "alias"
	page, err := cms.EntityStore.EntityFindByAttribute(ENTITY_TYPE_PAGE, "alias", ""+alias+"")

	if err != nil {
		log.Println(err.Error())
		return nil
	}

	if page != nil {
		return page
	}

	// Try to find by "/alias"
	page, err = cms.EntityStore.EntityFindByAttribute(ENTITY_TYPE_PAGE, "alias", "/"+alias+"")

	if err != nil {
		log.Println(err.Error())
		return nil
	}

	if page != nil {
		return page
	}

	return nil
}
