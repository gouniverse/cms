package types

import (
	"github.com/dracory/dataobject"
	"github.com/dromara/carbon/v2"
)

type WebWidgetInterface interface {
	dataobject.DataObjectInterface

	CreatedAt() string
	SetCreatedAt(createdAt string)
	CreatedAtCarbon() *carbon.Carbon
	Handle() string
	SetHandle(handle string)
	ID() string
	SetID(id string)
	Name() string
	SetName(name string)
	Status() string
	SetStatus(status string)
}
