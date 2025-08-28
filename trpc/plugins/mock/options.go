package mock

import "time"

// Config mock plugin config.
type Config []*Item

// MockItem specific mock items.
// Deprecated: use Item instead.
type MockItem = Item

// Item Specific mock items
type Item struct {
	Method        string
	Retcode       int
	Retmsg        string
	Delay         int
	delay         time.Duration
	Timeout       bool
	Body          string
	data          []byte
	Serialization int // json jce pb
	Percent       int
}

type options struct {
	mocks []*Item
}
