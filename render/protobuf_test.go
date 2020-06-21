package render

import (
	testdata "github.com/gin-gonic/gin/testdata/protoexample"
	"github.com/stretchr/testify/assert"
	"github.com/golang/protobuf/proto"
	"net/http/httptest"
	"testing"
)

func TestRenderProtoBuf(t *testing.T) {
	w := httptest.NewRecorder()
	reps := []int64{int64(1), int64(2)}
	label := "test"
	data := &testdata.Test{
		Label: &label,
		Reps:  reps,
	}

	(ProtoBuf{data}).WriteContentType(w)
	protoData, err := proto.Marshal(data)
	assert.NoError(t, err)
	assert.Equal(t, "application/x-protobuf", w.Header().Get("Content-Type"))

	err = (ProtoBuf{data}).Render(w)

	assert.NoError(t, err)
	assert.Equal(t, string(protoData), w.Body.String())
	assert.Equal(t, "application/x-protobuf", w.Header().Get("Content-Type"))
}

func TestRenderProtoBufFail(t *testing.T) {
	w := httptest.NewRecorder()
	data := &testdata.Test{}
	err := (ProtoBuf{data}).Render(w)
	assert.Nil(t, err)
}
