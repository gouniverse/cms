package cms

import "github.com/gouniverse/entitystore"

func (cms *Cms) BlockFindByID(blockID string) (*entitystore.Entity, error) {
	return cms.EntityStore.EntityFindByID(blockID)
}
