package cms

import (
	"github.com/dracory/entitystore"
)

func (cms *Cms) TemplateFindByID(templateID string) (*entitystore.Entity, error) {
	return cms.EntityStore.EntityFindByID(templateID)
}
