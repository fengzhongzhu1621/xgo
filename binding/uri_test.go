package binding

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestUriBinding(t *testing.T) {
	b := Uri
	assert.Equal(t, "uri", b.Name())

	type Tag struct {
		Name string `uri:"name"`
	}
	var tag Tag
	m := make(map[string][]string)
	m["name"] = []string{"thinkerou"}
	assert.NoError(t, b.BindUri(m, &tag))
	assert.Equal(t, "thinkerou", tag.Name)

	type NotSupportStruct struct {
		Name map[string]interface{} `uri:"name"`
	}
	var not NotSupportStruct
	assert.Error(t, b.BindUri(m, &not))
	assert.Equal(t, map[string]interface{}(nil), not.Name)
}

func TestUriInnerBinding(t *testing.T) {
	type Tag struct {
		Name string `uri:"name"`
		S    struct {
			Age int `uri:"age"`
		}
	}

	expectedName := "mike"
	expectedAge := 25

	m := map[string][]string{
		"name": {expectedName},
		"age":  {strconv.Itoa(expectedAge)},
	}

	var tag Tag
	assert.NoError(t, Uri.BindUri(m, &tag))
	assert.Equal(t, tag.Name, expectedName)
	assert.Equal(t, tag.S.Age, expectedAge)
}
