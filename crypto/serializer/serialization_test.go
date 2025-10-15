package serializer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJson(t *testing.T) {
	type Data struct {
		A int
		B string
	}
	s := &JSONSerialization{}
	body := []byte("{\"A\":1,\"B\":\"bb\"}")
	data := &Data{}

	err := s.Unmarshal(body, data)
	assert.Nil(t, err)
	assert.Equal(t, 1, data.A)
	assert.Equal(t, "bb", data.B)

	bytes, err := s.Marshal(data)
	assert.Nil(t, err)
	assert.Equal(t, body, bytes)

	// json-iterator issue https://github.com/golang/go/issues/48238#issuecomment-917321536
	m := make(map[string]string)
	m["a"] = "hello"
	bytes, err = s.Marshal(m)
	body = []byte("{\"a\":\"hello\"}")
	assert.Nil(t, err)
	assert.Equal(t, body, bytes)
}

func TestJsonPBNotImplProto(t *testing.T) {
	type Data struct {
		A int
		B string
	}
	s := &JSONPBSerialization{}
	data := &Data{A: 1, B: "test"}

	bytes, err := s.Marshal(data)
	assert.Nil(t, err)

	var data1 Data
	err = s.Unmarshal(bytes, &data1)
	assert.Nil(t, err)
	assert.Equal(t, data.A, data1.A)
	assert.Equal(t, data.B, data1.B)
}

func TestXML(t *testing.T) {
	type Data struct {
		A int
		B string
	}
	var tests = []struct {
		In Data
	}{
		{In: Data{1, "1"}},
		{In: Data{2, "2"}},
	}

	for _, tt := range tests {
		buf, err := Marshal(SerializationTypeXML, tt.In)
		assert.Nil(t, err)

		got := &Data{}
		err = Unmarshal(SerializationTypeXML, buf, got)
		assert.Nil(t, err)

		assert.Equal(t, tt.In.A, got.A)
		assert.Equal(t, tt.In.B, got.B)
	}
}
