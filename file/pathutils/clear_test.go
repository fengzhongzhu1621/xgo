package pathutils

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type cleanPathTest struct {
	path, result string
}

var cleanTests = []cleanPathTest{
	// Already clean
	{"/", "/"},
	{"/abc", "/abc"},
	{"/a/b/c", "/a/b/c"},
	{"/abc/", "/abc/"},
	{"/a/b/c/", "/a/b/c/"},

	// missing root
	{"", "/"},
	{"a/", "/a/"},
	{"abc", "/abc"},
	{"abc/def", "/abc/def"},
	{"a/b/c", "/a/b/c"},

	// Remove doubled slash
	{"//", "/"},
	{"/abc//", "/abc/"},
	{"/abc/def//", "/abc/def/"},
	{"/a/b/c//", "/a/b/c/"},
	{"/abc//def//ghi", "/abc/def/ghi"},
	{"//abc", "/abc"},
	{"///abc", "/abc"},
	{"//abc//", "/abc/"},

	// Remove . elements
	{".", "/"},
	{"./", "/"},
	{"/abc/./def", "/abc/def"},
	{"/./abc/def", "/abc/def"},
	{"/abc/.", "/abc/"},

	// Remove .. elements
	{"..", "/"},
	{"../", "/"},
	{"../../", "/"},
	{"../..", "/"},
	{"../../abc", "/abc"},
	{"/abc/def/ghi/../jkl", "/abc/def/jkl"},
	{"/abc/def/../ghi/../jkl", "/abc/jkl"},
	{"/abc/def/..", "/abc"},
	{"/abc/def/../..", "/"},
	{"/abc/def/../../..", "/"},
	{"/abc/def/../../..", "/"},
	{"/abc/def/../../../ghi/jkl/../../../mno", "/mno"},

	// Combinations
	{"abc/./../def", "/def"},
	{"abc//./../def", "/def"},
	{"abc/../../././../def", "/def"},
}

func TestPathClean(t *testing.T) {
	for _, test := range cleanTests {
		assert.Equal(t, test.result, CleanPath(test.path))
		assert.Equal(t, test.result, CleanPath(test.result))
	}
}

func TestPathCleanMallocs(t *testing.T) {
	// 跳过耗时函数
	if testing.Short() {
		t.Skip("skipping malloc count in short mode")
	}

	// CPU逻辑核心数 > 1 即多核处理
	if runtime.GOMAXPROCS(0) > 1 {
		t.Skip("skipping malloc count; GOMAXPROCS>1")
	}

	for _, test := range cleanTests {
		// AllocsPerRun 用于测量每次运行测试时的内存分配次数。这个方法应该返回一个整数，表示每次运行时的内存分配次数。
		allocs := testing.AllocsPerRun(100, func() { CleanPath(test.result) })
		assert.EqualValues(t, allocs, 0)
	}
}

func BenchmarkPathClean(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, test := range cleanTests {
			CleanPath(test.path)
		}
	}
}

func genLongPaths() (testPaths []cleanPathTest) {
	for i := 1; i <= 1234; i++ {
		ss := strings.Repeat("a", i)

		correctPath := "/" + ss
		testPaths = append(testPaths, cleanPathTest{
			path:   correctPath,
			result: correctPath,
		}, cleanPathTest{
			path:   ss,
			result: correctPath,
		}, cleanPathTest{
			path:   "//" + ss,
			result: correctPath,
		}, cleanPathTest{
			path:   "/" + ss + "/b/..",
			result: correctPath,
		})
	}
	return
}

func TestPathCleanLong(t *testing.T) {
	cleanTests := genLongPaths()

	for _, test := range cleanTests {
		assert.Equal(t, test.result, CleanPath(test.path))
		assert.Equal(t, test.result, CleanPath(test.result))
	}
}

func BenchmarkPathCleanLong(b *testing.B) {
	cleanTests := genLongPaths()
	// ResetTimer函数用于重置基准测试的计时器。
	// 它会将已经消耗的时间和内存分配次数清零，重新开始计时和统计。
	// 这个函数通常用于在基准测试函数中排除一些初始化或准备工作的时间和内存开销，使得基准测试的结果更加准确和公平。
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, test := range cleanTests {
			CleanPath(test.path)
		}
	}
}

func TestSlashAndCleanPath(t *testing.T) {
	tests := []struct {
		orig   string
		expect string
	}{
		// {"C:\\hello", "C:/hello"}, // Only works in windows
		{"", "."},
		{"//../foo", "/foo"},
		{"/../../", "/"},
		{"/hello/world/..", "/hello"},
		{"/..", "/"},
		{"/foo/..", "/"},
		{"/-/foo", "/-/foo"},
	}
	for _, v := range tests {
		res := SlashAndCleanPath(v.orig)
		if res != v.expect {
			t.Fatalf("Clean path(%v) expect(%v) but got(%v)", v.orig, v.expect, res)
		}
	}
}
