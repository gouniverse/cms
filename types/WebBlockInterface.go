package types

import (
	"github.com/golang-module/carbon/v2"
	"github.com/gouniverse/dataobject"
)

type WebBlockInterface interface {
	dataobject.DataObjectInterface

	CreatedAt() string
	SetCreatedAt(createdAt string)
	CreatedAtCarbon() carbon.Carbon
	Handle() string
	SetHandle(handle string)
	ID() string
	SetID(id string)
	Name() string
	SetName(name string)
	Status() string
	SetStatus(status string)
}
