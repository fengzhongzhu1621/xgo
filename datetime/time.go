package datetime

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/fengzhongzhu1621/xgo/cast"
	"github.com/fengzhongzhu1621/xgo/validator"
)

// Time TODO
type Time struct {
	time.Time `bson:",inline" json:",inline"`
}

// MarshalJSON TODO
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Time) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}

	dataStr := strings.Trim(string(data), "\"")

	timeType, isTime := validator.IsTime(dataStr)
	if isTime {
		t.Time = cast.Str2Time(dataStr, timeType)
		return nil
	}

	return json.Unmarshal(data, &t.Time)
}
