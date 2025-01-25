package newset

import (
	"encoding/json"
	"fmt"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func Test_UnmarshalJSON(t *testing.T) {
	s := []byte(`["test", "1", "2", "3"]`) //,["4,5,6"]]`)
	expected := mapset.NewSet(
		[]string{
			string(json.Number("1")),
			string(json.Number("2")),
			string(json.Number("3")),
			"test",
		}...,
	)

	actual := mapset.NewSet[string]()
	err := json.Unmarshal(s, actual)
	fmt.Println(actual) // Set{test, 1, 2, 3}
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	if !expected.Equal(actual) {
		t.Errorf("Expected no difference, got: %v", expected.Difference(actual))
	}
}

func Test_MarshalJSON(t *testing.T) {
	expected := mapset.NewSet(
		[]string{
			string(json.Number("1")),
			"test",
		}...,
	)

	b, err := json.Marshal(
		mapset.NewSet(
			[]string{
				"1",
				"test",
			}...,
		),
	)
	fmt.Println(b) // [91 34 49 34 44 34 116 101 115 116 34 93]
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	actual := mapset.NewSet[string]()
	err = json.Unmarshal(b, actual)
	fmt.Println(actual) // Set{1, test}
	if err != nil {
		t.Errorf("Error should be nil: %v", err)
	}

	if !expected.Equal(actual) {
		t.Errorf("Expected no difference, got: %v", expected.Difference(actual))
	}
}
