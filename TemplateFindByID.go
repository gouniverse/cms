package cms

import (
	"github.com/gouniverse/entitystore"
)

func (cms *Cms) TemplateFindByID(templateID string) (*entitystore.Entity, error) {
	return cms.EntityStore.EntityFindByID(templateID)
}
