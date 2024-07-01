package file

import "testing"

func TestSanitizedName(t *testing.T) {
	tests := []struct {
		orig   string
		expect string
	}{
		{"/a.log", "a.log"},
		{"//a.log", "a.log"},
		{"/../a.log", "../a.log"},
		{"/a/..b/c.log", "a/..b/c.log"},
		{"/a/b/../c.log", "a/c.log"},
		{"a/b/../c.log", "a/c.log"},
		{"/a/b/c/d.log", "a/b/c/d.log"},
		{"a/b/c/d.log", "a/b/c/d.log"},
	}
	for _, v := range tests {
		res := SanitizedName(v.orig)
		if res != v.expect {
			t.Fatalf("Clean path(%v) expect(%v) but got(%v)", v.orig, v.expect, res)
		}
	}
}
