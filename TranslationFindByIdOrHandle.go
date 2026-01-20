package cms

import (
	"context"

	"github.com/dracory/entitystore"
	"github.com/samber/lo"
)

func (cms *Cms) TranslationFindByIdOrHandle(idOrHandle string, language string) (string, error) {
	if idOrHandle == "" {
		return "", nil
	}

	handle := ""
	id := ""
	if !isNumeric(idOrHandle) {
		handle = idOrHandle
		entity, err := cms.EntityStore.EntityFindByHandle(context.Background(), ENTITY_TYPE_TRANSLATION, handle)
		if err != nil {
			return "", err
		}
		if entity == nil {
			return "", nil
		}
		id = entity.ID()
	} else {
		id = idOrHandle
	}

	attributes, err := cms.EntityStore.EntityAttributeList(context.Background(), id)
	if err != nil {
		return "", err
	}

	defaultTranslation := ""
	translation := ""

	lo.ForEach(attributes, func(attr entitystore.Attribute, index int) {
		if attr.AttributeKey() == cms.translationLanguageDefault {
			defaultTranslation = attr.AttributeValue()
			return
		}
		if attr.AttributeKey() == language {
			translation = attr.AttributeValue()
			return
		}
	})

	return lo.Ternary(translation == "", defaultTranslation, translation), nil
}
