package watch

import (
	"errors"

	"github.com/fengzhongzhu1621/xgo/ginx/apps/cmdb/resource_watch/stream"
)

// EventType TODO
type EventType string

const (
	// Create TODO
	Create EventType = "create"
	// Update TODO
	Update EventType = "update"
	// Delete TODO
	Delete EventType = "delete"
	// Unknown TODO
	Unknown EventType = "unknown"
)

// Validate TODO
func (e EventType) Validate() error {
	switch e {
	case Create, Update, Delete:
		return nil
	default:
		return errors.New("unsupported event type")
	}
}

func ConvertOperateType(typ stream.OperType) EventType {
	switch typ {
	case stream.Insert:
		return Create
	case stream.Replace, stream.Update:
		return Update
	case stream.Delete:
		return Delete
	default:
		return Unknown
	}
}
