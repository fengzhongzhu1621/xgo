package newset

import (
	"encoding/json"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func TestThreadUnsafeSet_MarshalJSON(t *testing.T) {
	expected := mapset.NewThreadUnsafeSet[int64](1, 2, 3)
	actual := mapset.NewThreadUnsafeSet[int64]()

	// test Marshal from Set method
	b, err := expected.MarshalJSON()
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	err = json.Unmarshal(b, actual)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	if !expected.Equal(actual) {
		t.Errorf("Expected no difference, got: %v", expected.Difference(actual))
	}

	// test Marshal from json package
	b, err = json.Marshal(expected)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	err = json.Unmarshal(b, actual)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	if !expected.Equal(actual) {
		t.Errorf("Expected no difference, got: %v", expected.Difference(actual))
	}
}

func TestThreadUnsafeSet_UnmarshalJSON(t *testing.T) {
	expected := mapset.NewThreadUnsafeSet[int64](1, 2, 3)
	actual := mapset.NewThreadUnsafeSet[int64]()

	// test Unmarshal from Set method
	err := actual.UnmarshalJSON([]byte(`[1, 2, 3]`))
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}
	if !expected.Equal(actual) {
		t.Errorf("Expected no difference, got: %v", expected.Difference(actual))
	}

	// test Unmarshal from json package
	actual = mapset.NewThreadUnsafeSet[int64]()
	err = json.Unmarshal([]byte(`[1, 2, 3]`), actual)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}
	if !expected.Equal(actual) {
		t.Errorf("Expected no difference, got: %v", expected.Difference(actual))
	}
}

func TestThreadUnsafeSet_MarshalJSON_Struct(t *testing.T) {
	expected := &testStruct{"test", mapset.NewThreadUnsafeSet("a")}

	b, err := json.Marshal(&testStruct{"test", mapset.NewThreadUnsafeSet("a")})
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	actual := &testStruct{}
	err = json.Unmarshal(b, actual)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	if !expected.Set.Equal(actual.Set) {
		t.Errorf("Expected no difference, got: %v", expected.Set.Difference(actual.Set))
	}
}

func TestThreadUnsafeSet_UnmarshalJSON_Struct(t *testing.T) {
	expected := &testStruct{"test", mapset.NewThreadUnsafeSet("a", "b", "c")}
	actual := &testStruct{}

	err := json.Unmarshal([]byte(`{"other":"test", "set":["a", "b", "c"]}`), actual)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}
	if !expected.Set.Equal(actual.Set) {
		t.Errorf("Expected no difference, got: %v", expected.Set.Difference(actual.Set))
	}

	expectedComplex := mapset.NewThreadUnsafeSet(
		struct{ Val string }{Val: "a"},
		struct{ Val string }{Val: "b"},
	)
	actualComplex := mapset.NewThreadUnsafeSet[struct{ Val string }]()

	err = actualComplex.UnmarshalJSON([]byte(`[{"Val": "a"}, {"Val": "b"}]`))
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}
	if !expectedComplex.Equal(actualComplex) {
		t.Errorf("Expected no difference, got: %v", expectedComplex.Difference(actualComplex))
	}

	actualComplex = mapset.NewThreadUnsafeSet[struct{ Val string }]()
	err = json.Unmarshal([]byte(`[{"Val": "a"}, {"Val": "b"}]`), actualComplex)
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}
	if !expectedComplex.Equal(actualComplex) {
		t.Errorf("Expected no difference, got: %v", expectedComplex.Difference(actualComplex))
	}
}

// this serves as an example of how to correctly unmarshal a struct with a Set property
type testStruct struct {
	Other string
	Set   mapset.Set[string]
}

func (t *testStruct) UnmarshalJSON(b []byte) error {
	raw := struct {
		Other string
		Set   []string
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	t.Other = raw.Other
	t.Set = mapset.NewThreadUnsafeSet(raw.Set...)

	return nil
}
