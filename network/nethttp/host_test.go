package nethttp

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/netutil"
	"github.com/stretchr/testify/assert"
)

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
			t.Errorf(
				"test case for %#v failed, server or port mismatch to expected values (%s:%s)",
				testcase.src,
				s,
				p,
			)
		}
	}
}

// TestConvertMapToQueryString Convert map to url query string.
// func ConvertMapToQueryString(param map[string]any) string
func TestConvertMapToQueryString(t *testing.T) {
	m := map[string]any{
		"c": 3,
		"a": 1,
		"b": 2,
	}
	qs := netutil.ConvertMapToQueryString(m)
	assert.Equal(t, "a=1&b=2&c=3", qs)
}

// TestEncodeUrl Encode url query string values.
// func EncodeUrl(urlStr string) (string, error)
func TestEncodeUrl(t *testing.T) {
	urlAddr := "http://www.lancet.com?a=1&b=[2]"
	encodedUrl, err := netutil.EncodeUrl(urlAddr)
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, "http://www.lancet.com?a=1&b=%5B2%5D", encodedUrl)
}
