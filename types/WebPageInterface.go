package types

import "github.com/dromara/carbon/v2"

type WebPageInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	ID() string
	SetID(id string)
	Alias() string
	SetAlias(alias string)
	CreatedAt() string
	SetCreatedAt(createdAt string)
	CreatedAtCarbon() *carbon.Carbon
	CanonicalUrl() string
	SetCanonicalUrl(canonicalUrl string)
	Content() string
	SetContent(content string)
	Editor() string
	SetEditor(editor string)
	Handle() string
	SetHandle(handle string)
	MetaDescription() string
	SetMetaDescription(metaDescription string)
	MetaKeywords() string
	SetMetaKeywords(metaKeywords string)
	MetaRobots() string
	SetMetaRobots(metaRobots string)
	Name() string
	SetName(name string)
	Status() string
	SetStatus(status string)
	Title() string
	SetTitle(title string)
	TemplateID() string
	SetTemplateID(templateID string)
	UpdatedAt() string
	SetUpdatedAt(updatedAt string)
	UpdatedAtCarbon() *carbon.Carbon

	IsActive() bool
	IsInactive() bool
	IsDeleted() bool
}
