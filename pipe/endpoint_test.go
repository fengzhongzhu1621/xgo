package pipe

import (
	"strconv"
	"testing"
	"time"
)

type TestEndpoint struct {
	limit  int
	prefix string
	c      chan []byte
	result []string
}

func (e *TestEndpoint) StartReading() {
	go func() {
		for i := 0; i < e.limit; i++ {
			e.c <- []byte(e.prefix + strconv.Itoa(i))
		}
		time.Sleep(time.Millisecond) // should be enough for smaller channel to catch up with long one
		close(e.c)
	}()
}

func (e *TestEndpoint) Terminate() {
}

func (e *TestEndpoint) Output() chan []byte {
	return e.c
}

func (e *TestEndpoint) Send(msg []byte) bool {
	e.result = append(e.result, string(msg))
	return true
}

func TestEndpointPipe(t *testing.T) {
	one := &TestEndpoint{2, "one:", make(chan []byte), make([]string, 0)}
	two := &TestEndpoint{4, "two:", make(chan []byte), make([]string, 0)}
	PipeEndpoints(one, two)
	if len(one.result) != 4 || len(two.result) != 2 {
		t.Errorf("Invalid lengths, should be 4 and 2: %v %v", one.result, two.result)
	} else if one.result[0] != "two:0" || two.result[0] != "one:0" {
		t.Errorf("Invalid first results, should be two:0 and one:0: %#v %#v", one.result[0], two.result[0])
	}
}
