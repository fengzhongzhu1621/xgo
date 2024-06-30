package file

import "testing"

func TestSanitizedName(t *testing.T) {
	tests := []struct {
		orig   string
		expect string
	}{
		// {"C:\\hello", "C:/hello"}, // Only works in windows
		{"", "."},
		{"//../foo", "foo"},
		{"/../../", ""},
		{"/hello/world/..", "hello"},
		{"/..", ""},
		{"/foo/..", ""},
		{"/-/foo", "-/foo"},
	}
	for _, v := range tests {
		res := SlashAndCleanPath(v.orig)
		if res != v.expect {
			t.Fatalf("Clean path(%v) expect(%v) but got(%v)", v.orig, v.expect, res)
		}
	}
}