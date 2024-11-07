package types

import (
	"github.com/golang-module/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/uid"
)

type WebWidget struct {
	dataobject.DataObject
}

var _ WebWidgetInterface = (*WebWidget)(nil)

func NewWebWidget() *WebWidget {
	o := &WebWidget{}
	o.SetID(uid.HumanUid())
	// o.SetStatus(WEBPAGE_STATUS_DRAFT)
	return o
}

func NewWebWidgetFromExistingData(data map[string]string) *WebWidget {
	o := &WebWidget{}
	o.Hydrate(data)
	return o
}

// == SETTERS AND GETTERS ===================================================
func (o *WebWidget) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *WebWidget) CreatedAtCarbon() carbon.Carbon {
	return carbon.Parse(o.CreatedAt())
}

func (o *WebWidget) SetCreatedAt(createdAt string) {
	o.Set(COLUMN_CREATED_AT, createdAt)
}

func (o *WebWidget) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

func (o *WebWidget) SetHandle(handle string) {
	o.Set(COLUMN_HANDLE, handle)
}

func (o *WebWidget) Name() string {
	return o.Get(COLUMN_NAME)
}

func (o *WebWidget) SetName(name string) {
	o.Set(COLUMN_NAME, name)
}

func (o *WebWidget) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *WebWidget) SetStatus(status string) {
	o.Set(COLUMN_STATUS, status)
}
