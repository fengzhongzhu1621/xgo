package watch

import (
	"encoding/json"

	"github.com/fengzhongzhu1621/xgo/network/nethttp"
)

type WatchEventResp struct {
	nethttp.BaseResp `json:",inline"`
	Data             *WatchResp `json:"data"`
}

type WatchResp struct {
	// watched events or not
	Watched bool                `json:"bk_watched"`
	Events  []*WatchEventDetail `json:"bk_events"`
}

type WatchEventDetail struct {
	Cursor    string     `json:"bk_cursor"`
	Resource  CursorType `json:"bk_resource"`
	EventType EventType  `json:"bk_event_type"`
	// Default instance is JsonString type
	Detail DetailInterface `json:"bk_detail"`

	// ChainNode is the chain node of this watch event
	// NOTE: this is a special return value for internal use only
	ChainNode *ChainNode `json:"-"`
}

type jsonWatchEventDetail struct {
	Cursor    string          `json:"bk_cursor"`
	Resource  CursorType      `json:"bk_resource"`
	EventType EventType       `json:"bk_event_type"`
	Detail    json.RawMessage `json:"bk_detail"`
}

type DetailInterface interface {
	Name() string
}

// JsonString TODO
type JsonString string

// Name TODO
func (j JsonString) Name() string {
	return "JsonString"
}

// MarshalJSON TODO
func (j JsonString) MarshalJSON() ([]byte, error) {
	if j == "" {
		j = "{}"
	}
	return []byte(j), nil
}
