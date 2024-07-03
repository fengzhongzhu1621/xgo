package nethttp

import "testing"

var tellHostPortTests = []struct {
	src    string
	ssl    bool
	server string
	port   string
}{
	{"localhost", false, "localhost", "80"},
	{"localhost:8080", false, "localhost", "8080"},
	{"localhost", true, "localhost", "443"},
	{"localhost:8080", true, "localhost", "8080"},
}

func TestTellHostPort(t *testing.T) {
	for _, testcase := range tellHostPortTests {
		s, p, e := TellHostPort(testcase.src, testcase.ssl)
		if testcase.server == "" {
			if e == nil {
				t.Errorf("test case for %#v failed, error was not returned", testcase.src)
			}
		} else if e != nil {
			t.Errorf("test case for %#v failed, error should not happen", testcase.src)
		}
		if testcase.server != s || testcase.port != p {
			t.Errorf("test case for %#v failed, server or port mismatch to expected values (%s:%s)", testcase.src, s, p)
		}
	}
}
