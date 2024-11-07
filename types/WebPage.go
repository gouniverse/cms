package types

import (
	"github.com/golang-module/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/uid"
)

// var _ dataobject.DataObjectFluentInterface = (*WebPage)(nil) // verify it extends the data object interface

const WEBPAGE_STATUS_DELETED = "deleted"
const WEBPAGE_STATUS_DRAFT = "draft"
const WEBPAGE_STATUS_ACTIVE = "active"
const WEBPAGE_STATUS_INACTIVE = "inactive"

const WEBPAGE_EDITOR_BLOCKAREA = "blockarea"
const WEBPAGE_EDITOR_BLOCKEDITOR = "blockeditor"
const WEBPAGE_EDITOR_CODEMIRROR = "codemirror"
const WEBPAGE_EDITOR_HTMLAREA = "htmlarea"
const WEBPAGE_EDITOR_MARKDOWN = "markdown"
const WEBPAGE_EDITOR_TEXTAREA = "textarea"

type WebPage struct {
	dataobject.DataObject
}

var _ dataobject.DataObjectInterface = (*WebPage)(nil)

func NewWebPage() *WebPage {
	o := &WebPage{}
	o.SetID(uid.HumanUid())
	o.SetStatus(WEBPAGE_STATUS_DRAFT)
	return o
}

func NewWebPageFromExistingData(data map[string]string) *WebPage {
	o := &WebPage{}
	o.Hydrate(data)
	return o
}

var _ dataobject.DataObjectInterface = (*WebPage)(nil)

// == METHODS ===============================================================

func (o *WebPage) IsActive() bool {
	return o.Status() == WEBPAGE_STATUS_ACTIVE
}

func (o *WebPage) IsInactive() bool {
	return o.Status() == WEBPAGE_STATUS_INACTIVE
}

func (o *WebPage) IsDeleted() bool {
	return o.Status() == WEBPAGE_STATUS_DELETED
}

// == SETTERS AND GETTERS =====================================================

func (o *WebPage) Alias() string {
	return o.Get(COLUMN_ALIAS)
}

func (o *WebPage) SetAlias(alias string) {
	o.Set(COLUMN_ALIAS, alias)
}

func (o *WebPage) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *WebPage) SetCreatedAt(createdAt string) {
	o.Set(COLUMN_CREATED_AT, createdAt)
}

func (o *WebPage) CreatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *WebPage) CanonicalUrl() string {
	return o.Get(COLUMN_CANONICAL_URL)
}

func (o *WebPage) SetCanonicalUrl(canonicalUrl string) {
	o.Set(COLUMN_CANONICAL_URL, canonicalUrl)
}

func (o *WebPage) Content() string {
	return o.Get(COLUMN_CONTENT)
}

func (o *WebPage) SetContent(content string) {
	o.Set(COLUMN_CONTENT, content)
}

func (o *WebPage) Editor() string {
	return o.Get(COLUMN_EDITOR)
}

func (o *WebPage) SetEditor(editor string) {
	o.Set(COLUMN_EDITOR, editor)
}

func (o *WebPage) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *WebPage) SetID(id string) {
	o.Set(COLUMN_ID, id)
}

func (o *WebPage) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

func (o *WebPage) SetHandle(handle string) {
	o.Set(COLUMN_HANDLE, handle)
}

func (o *WebPage) MetaDescription() string {
	return o.Get(COLUMN_META_DESCRIPTION)
}

func (o *WebPage) SetMetaDescription(metaDescription string) {
	o.Set(COLUMN_META_DESCRIPTION, metaDescription)
}

func (o *WebPage) MetaKeywords() string {
	return o.Get(COLUMN_META_KEYWORDS)
}

func (o *WebPage) SetMetaKeywords(metaKeywords string) {
	o.Set(COLUMN_META_KEYWORDS, metaKeywords)
}

func (o *WebPage) MetaRobots() string {
	return o.Get(COLUMN_META_ROBOTS)
}

func (o *WebPage) SetMetaRobots(metaRobots string) {
	o.Set(COLUMN_META_ROBOTS, metaRobots)
}

func (o *WebPage) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *WebPage) SetName(name string) {
	o.Set(COLUMN_NAME, name)
}

func (o *WebPage) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *WebPage) SetStatus(status string) {
	o.Set(COLUMN_STATUS, status)
}

func (o *WebPage) Title() string {
	return o.Get(COLUMN_TITLE)
}

func (o *WebPage) SetTitle(title string) {
	o.Set(COLUMN_TITLE, title)
}

func (o *WebPage) TemplateID() string {
	return o.Get(COLUMN_TEMPLATE_ID)
}

func (o *WebPage) SetTemplateID(templateID string) {
	o.Set(COLUMN_TEMPLATE_ID, templateID)
}

func (o *WebPage) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *WebPage) SetUpdatedAt(updatedAt string) {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
}

func (o *WebPage) UpdatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.UpdatedAt())
}
