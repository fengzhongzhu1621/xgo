package newset

import (
	"fmt"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func Test_Each(t *testing.T) {
	a := mapset.NewSet[string]()

	a.Add("Z")
	a.Add("Y")
	a.Add("X")
	a.Add("W")

	b := mapset.NewSet[string]()
	a.Each(func(elem string) bool {
		b.Add(elem)
		return false
	})
	fmt.Println(a) // Set{Y, X, W, Z}
	fmt.Println(b) // Set{Z, Y, X, W}

	// true
	if !a.Equal(b) {
		t.Error("The sets are not equal after iterating (Each) through the first set")
	}

	var count int
	a.Each(func(elem string) bool {
		if count == 2 {
			return true
		}
		count++
		return false
	})
	if count != 2 {
		t.Error("Iteration should stop on the way")
	}
}

func Test_Iter(t *testing.T) {
	a := mapset.NewSet[string]()

	a.Add("Z")
	a.Add("Y")
	a.Add("X")
	a.Add("W")

	b := mapset.NewSet[string]()
	for val := range a.Iter() {
		b.Add(val)
	}

	if !a.Equal(b) {
		t.Error("The sets are not equal after iterating (Iter) through the first set")
	}
}

func Test_Iterator(t *testing.T) {
	a := mapset.NewSet[string]()

	a.Add("Z")
	a.Add("Y")
	a.Add("X")
	a.Add("W")

	b := mapset.NewSet[string]()
	for val := range a.Iterator().C {
		b.Add(val)
	}

	if !a.Equal(b) {
		t.Error("The sets are not equal after iterating (Iterator) through the first set")
	}
}

func Test_IteratorStop(t *testing.T) {
	a := mapset.NewSet[string]()

	a.Add("Z")
	a.Add("Y")
	a.Add("X")
	a.Add("W")

	it := a.Iterator()
	it.Stop()
	for range it.C {
		t.Error("The iterating (Iterator) did not stop after Stop() has been called")
	}
}

type yourType struct {
	name string
}

func Test_ExampleIterator(t *testing.T) {
	s := mapset.NewSet(
		[]*yourType{
			{name: "Alise"},
			{name: "Bob"},
			{name: "John"},
			{name: "Nick"},
		}...,
	)

	var found *yourType
	it := s.Iterator()

	for elem := range it.C {
		if elem.name == "John" {
			found = elem
			it.Stop()
		}
	}

	if found == nil || found.name != "John" {
		t.Fatalf("expected iterator to have found `John` record but got nil or something else")
	}
}
