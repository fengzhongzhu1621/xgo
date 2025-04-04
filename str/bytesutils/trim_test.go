package bytesutils

import "testing"

var eolTests = []string{
	"", "\n", "\r\n", "ok\n", "ok\n",
	"quite long string for our test\n",
	"quite long string for our test\r\n",
}

var eolAnswers = []string{
	"", "", "", "ok", "ok",
	"quite long string for our test", "quite long string for our test",
}

func TestTrimEOL(t *testing.T) {
	for n := 0; n < len(eolTests); n++ {
		answer := TrimEOL([]byte(eolTests[n]))
		if string(answer) != eolAnswers[n] {
			t.Errorf("Answer '%s' did not match predicted '%s'", answer, eolAnswers[n])
		}
	}
}

func BenchmarkTrimEOL(b *testing.B) {
	for n := 0; n < b.N; n++ {
		TrimEOL([]byte(eolTests[n%len(eolTests)]))
	}
}
