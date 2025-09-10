package types

import (
	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/uid"
)

type WebBlock struct {
	dataobject.DataObject
}

var _ WebBlockInterface = (*WebBlock)(nil)

func NewWebBlock() *WebBlock {
	o := &WebBlock{}
	o.SetID(uid.HumanUid())
	// o.SetStatus(WEBPAGE_STATUS_DRAFT)
	return o
}

func NewWebBlockFromExistingData(data map[string]string) *WebBlock {
	o := &WebBlock{}
	o.Hydrate(data)
	return o
}

// == SETTERS AND GETTERS ===================================================

func (o *WebBlock) Content() string {
	return o.Get(COLUMN_CONTENT)
}

func (o *WebBlock) SetContent(content string) {
	o.Set(COLUMN_CONTENT, content)
}

func (o *WebBlock) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *WebBlock) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *WebBlock) SetCreatedAt(createdAt string) {
	o.Set(COLUMN_CREATED_AT, createdAt)
}

func (o *WebBlock) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

func (o *WebBlock) SetHandle(handle string) {
	o.Set(COLUMN_HANDLE, handle)
}

func (o *WebBlock) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *WebBlock) SetName(name string) {
	o.Set(COLUMN_NAME, name)
}

func (o *WebBlock) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *WebBlock) SetStatus(status string) {
	o.Set(COLUMN_STATUS, status)
}
