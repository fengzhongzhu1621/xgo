package watch

type WatchEventOptions struct {
	// event types you want to care, empty means all.
	EventTypes []EventType `json:"bk_event_types"`
	// the fields you only care, if nil, means all.
	Fields []string `json:"bk_fields"`
	// unix seconds time to where you want to watch from.
	// it's like Cursor, but StartFrom and Cursor can not use at the same time.
	StartFrom int64 `json:"bk_start_from"`
	// the cursor you hold previous, means you want to watch event form here.
	Cursor string `json:"bk_cursor"`
	// the resource kind you want to watch
	Resource CursorType       `json:"bk_resource"`
	Filter   WatchEventFilter `json:"bk_filter"`
}

type WatchEventFilter struct {
	// SubResource the sub resource you want to watch, eg. object ID of the instance resource, watch all if not set
	SubResource string `json:"bk_sub_resource,omitempty"`
	// SubResources is the sub resources you want to watch, NOTE: this is a special parameter for internal use only
	SubResources []string `json:"-"`
}
