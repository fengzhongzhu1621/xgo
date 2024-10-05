package print

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dablelv/cyan/encoding"
)

// Student 学生信息
type Student struct {
	Name string
	Addr HomeInfo
	M    map[string]string
}

// HomeInfo 家庭住址
type HomeInfo struct {
	Province     string
	City         string
	County       string
	Street       string
	DetailedAddr string
}

var student = Student{
	Name: "dablelv",
	Addr: HomeInfo{
		Province:     "Guangdong",
		City:         "Shenzhen",
		County:       "Baoan",
		Street:       "Xixiang",
		DetailedAddr: "Shengtianqi",
	},
	M: map[string]string{
		"hobby": "pingpopng",
	},
}

func TestPrintStruct(t *testing.T) {
	s, _ := encoding.ToIndentJSON(&student)
	fmt.Printf("student = %v\n", s)
}

func TestPrintStruct2(t *testing.T) {
	bs, _ := json.Marshal(student)
	var out bytes.Buffer
	json.Indent(&out, bs, "", "\t")
	fmt.Printf("student=%v\n", out.String())
}
