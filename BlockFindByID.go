package cms

import (
	"errors"

	"github.com/dracory/entitystore"
	"github.com/gouniverse/cms/types"
	"github.com/samber/lo"
)

// func (cms *Cms) BlockFindByID(blockID string) (*entitystore.Entity, error) {
// 	return cms.EntityStore.EntityFindByID(blockID)
// }

func (cms *Cms) BlockFindByID(blockID string) (types.WebBlockInterface, error) {
	entity, err := cms.EntityStore.EntityFindByID(blockID)

	if err != nil {
		return nil, err
	}

	if entity == nil {
		return nil, nil
	}

	if entity.Type() != ENTITY_TYPE_BLOCK {
		return nil, errors.New("entity is not a block")
	}

	block := &types.WebBlock{}
	block.SetID(entity.ID())
	block.SetHandle(entity.Handle())

	attrs, err := entity.GetAttributes()
	if err != nil {
		return nil, err
	}

	lo.ForEach(attrs, func(attr entitystore.Attribute, index int) {
		switch attr.AttributeKey() {
		case "name":
			block.SetName(attr.AttributeValue())
		case "status":
			block.SetStatus(attr.AttributeValue())
		case "content":
			block.SetContent(attr.AttributeValue())
		}
	})

	return block, nil
}
