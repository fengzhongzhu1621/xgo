package stream

import (
	"time"
)

type TimeStamp struct {
	// the most significant 32 bits are a time_t value (seconds since the Unix epoch)
	Sec uint32 `json:"sec" bson:"sec"`
	// the least significant 32 bits are an incrementing ordinal for operations within a given second.
	Nano uint32 `json:"nano" bson:"nano"`
}

// String 用于打印
func (t TimeStamp) String() string {
	return time.Unix(int64(t.Sec), int64(t.Nano)).Format(time.DateTime)
}
