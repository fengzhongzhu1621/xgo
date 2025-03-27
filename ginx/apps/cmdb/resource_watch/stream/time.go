package stream

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type TimeStamp struct {
	// the most significant 32 bits are a time_t value (seconds since the Unix epoch)
	Sec int64 `json:"sec" bson:"sec"`
	// the least significant 32 bits are an incrementing ordinal for operations within a given second.
	Nano int64 `json:"nano" bson:"nano"`
}

// String 用于打印
func (t TimeStamp) String() string {
	return time.Unix(int64(t.Sec), int64(t.Nano)).Format("2006-01-02/15:04:05")
}

func (t TimeStamp) MarshalBSONValue() (bsontype.Type, []byte, error) {
	// 1. 构造 time.Time 对象
	goTime := time.Unix(t.Sec, t.Nano)

	// 2. 使用 bson.MarshalValue 直接序列化 time.Time
	return bson.MarshalValue(goTime)
}
