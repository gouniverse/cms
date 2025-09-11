package cms

import (
	"errors"

	"github.com/dracory/entitystore"
	"github.com/dracory/str"
	"github.com/gouniverse/cms/types"
	"github.com/samber/lo"
)

type WebPageQueryOptions struct {
	ID        string
	IDIn      []string
	Handle    string
	Status    string
	StatusIn  []string
	Offset    int
	Limit     int
	SortOrder string
	OrderBy   string
	CountOnly bool
}

func (cms *Cms) WebPageCreate(page types.WebPageInterface) error {
	name := page.Name()
	status := page.Status()
	alias := "/" + str.Slugify(name+"-"+str.Random(16), '-')

	if name == "" {
		return errors.New("page name is empty")
	}

	pageEntity, err := cms.EntityStore.EntityCreateWithTypeAndAttributes(ENTITY_TYPE_PAGE, map[string]string{
		"name":   name,
		"status": status,
		"title":  name,
		"alias":  alias,
	})

	if err != nil {
		return err
	}

	if pageEntity == nil {
		return errors.New("page entity is nil")
	}

	page.SetID(pageEntity.ID())

	return nil
}

func (cms *Cms) WebPageFindByID(pageID string) (types.WebPageInterface, error) {
	if pageID == "" {
		return nil, errors.New("page id is empty")
	}

	list, err := cms.WebPageList(WebPageQueryOptions{
		ID:    pageID,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (cms *Cms) WebPageCount(options WebPageQueryOptions) (int64, error) {
	options.CountOnly = true
	records, err := cms.WebPageList(options)

	if err != nil {
		return 0, err
	}

	return int64(len(records)), nil
}

func (cms *Cms) WebPageList(options WebPageQueryOptions) ([]types.WebPageInterface, error) {
	entityList, errEntityList := cms.EntityStore.EntityList(entitystore.EntityQueryOptions{
		EntityType:   ENTITY_TYPE_PAGE,
		ID:           options.ID,
		EntityHandle: options.Handle,
		Limit:        uint64(options.Limit),
		Offset:       uint64(options.Offset),
		CountOnly:    options.CountOnly,
	})

	if errEntityList != nil {
		return []types.WebPageInterface{}, errEntityList
	}

	list := []types.WebPageInterface{}
	var errMap error = nil

	lo.ForEach(entityList, func(entity entitystore.Entity, index int) {
		attrs, err := entity.GetAttributes()
		if err != nil {
			errMap = err
		}
		cmsPageMap := map[string]string{
			"id":         entity.ID(),
			"handle":     entity.Handle(),
			"created_at": entity.CreatedAt().String(),
			"updated_at": entity.UpdatedAt().String(),
		}

		lo.ForEach(attrs, func(attr entitystore.Attribute, index int) {
			cmsPageMap[attr.AttributeKey()] = attr.AttributeValue()
		})

		cmsPage := types.NewWebPageFromExistingData(cmsPageMap)
		list = append(list, cmsPage)
	})

	if errMap != nil {
		return []types.WebPageInterface{}, errEntityList
	}

	return list, nil
}

func (cms *Cms) WebPageUpdate(page types.WebPageInterface) error {
	if page == nil {
		return errors.New("page is nil")
	}

	if page.ID() == "" {
		return errors.New("page id is empty")
	}

	dataChanged := page.DataChanged()

	if len(dataChanged) == 0 {
		return nil
	}

	err := cms.EntityStore.AttributesSet(page.ID(), dataChanged)
	if err != nil {
		return err
	}

	return nil
}
