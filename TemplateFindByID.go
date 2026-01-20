package cms

import (
	"context"

	"github.com/dracory/entitystore"
)

func (cms *Cms) TemplateFindByID(templateID string) (*entitystore.Entity, error) {
	return cms.EntityStore.EntityFindByID(context.Background(), templateID)
}
