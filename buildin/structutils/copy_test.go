package structutils

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/convertor"
)

// Copies each field from the source struct into the destination struct.
// Use json.Marshal/Unmarshal, so json tag should be set for fields of dst and src struct.
// func CopyProperties[T, U any](dst T, src U) (err error)
func TestCopyProperties(t *testing.T) {
	type Disk struct {
		Name    string  `json:"name"`
		Total   string  `json:"total"`
		Used    string  `json:"used"`
		Percent float64 `json:"percent"`
	}

	type DiskVO struct {
		Name    string  `json:"name"`
		Total   string  `json:"total"`
		Used    string  `json:"used"`
		Percent float64 `json:"percent"`
	}

	type Indicator struct {
		Id      string    `json:"id"`
		Ip      string    `json:"ip"`
		UpTime  string    `json:"upTime"`
		LoadAvg string    `json:"loadAvg"`
		Cpu     int       `json:"cpu"`
		Disk    []Disk    `json:"disk"`
		Stop    chan bool `json:"-"`
	}

	type IndicatorVO struct {
		Id      string   `json:"id"`
		Ip      string   `json:"ip"`
		UpTime  string   `json:"upTime"`
		LoadAvg string   `json:"loadAvg"`
		Cpu     int64    `json:"cpu"`
		Disk    []DiskVO `json:"disk"`
	}

	indicator := &Indicator{Id: "001", Ip: "127.0.0.1", Cpu: 1, Disk: []Disk{
		{Name: "disk-001", Total: "100", Used: "1", Percent: 10},
		{Name: "disk-002", Total: "200", Used: "1", Percent: 20},
		{Name: "disk-003", Total: "300", Used: "1", Percent: 30},
	}}

	indicatorVO := IndicatorVO{}

	err := convertor.CopyProperties(&indicatorVO, indicator)

	if err != nil {
		return
	}

	fmt.Println(indicatorVO.Id)
	fmt.Println(indicatorVO.Ip)
	fmt.Println(len(indicatorVO.Disk))

	// Output:
	// 001
	// 127.0.0.1
	// 3
}

// Creates a deep copy of passed item, can't clone unexported field of struct.
// func DeepClone[T any](src T) T
func TestDeepClone(t *testing.T) {
	type Struct struct {
		Str        string
		Int        int
		Float      float64
		Bool       bool
		Nil        interface{}
		unexported string
	}

	cases := []interface{}{
		true,
		1,
		0.1,
		map[string]int{
			"a": 1,
			"b": 2,
		},
		&Struct{
			Str:   "test",
			Int:   1,
			Float: 0.1,
			Bool:  true,
			Nil:   nil,
		},
	}

	for _, item := range cases {
		cloned := convertor.DeepClone(item)

		isPointerEqual := &cloned == &item
		fmt.Println(cloned, isPointerEqual)
	}

	// Output:
	// true false
	// 1 false
	// 0.1 false
	// map[a:1 b:2] false
	// &{test 1 0.1 true <nil> } false
}
