package strcase

import (
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/stretchr/testify/assert"
)

func toCamel(tb testing.TB) {
	cases := [][]string{
		{"test_case", "TestCase"},
		{"test.case", "TestCase"},
		{"test", "Test"},
		{"TestCase", "TestCase"},
		{" test  case ", "TestCase"},
		{"", ""},
		{"many_many_words", "ManyManyWords"},
		{"AnyKind of_string", "AnyKindOfString"},
		{"odd-fix", "OddFix"},
		{"numbers2And55with000", "Numbers2And55With000"},
		{"ID", "Id"},
	}
	for _, i := range cases {
		in := i[0]
		out := i[1]
		result := ToCamel(in)
		if result != out {
			tb.Errorf("%q (%q != %q)", in, result, out)
		}
	}
}

func TestToCamel(t *testing.T) {
	toCamel(t)
}

func BenchmarkToCamel(b *testing.B) {
	benchmarkCamelTest(b, toCamel)
}

func toLowerCamel(tb testing.TB) {
	cases := [][]string{
		{"foo-bar", "fooBar"},
		{"TestCase", "testCase"},
		{"", ""},
		{"AnyKind of_string", "anyKindOfString"},
		{"AnyKind.of-string", "anyKindOfString"},
		{"ID", "id"},
		{"some string", "someString"},
		{" some string", "someString"},
	}
	for _, i := range cases {
		in := i[0]
		out := i[1]
		result := ToLowerCamel(in)
		if result != out {
			tb.Errorf("%q (%q != %q)", in, result, out)
		}
	}
}

func TestToLowerCamel(t *testing.T) {
	toLowerCamel(t)
}

func TestCustomAcronymsToCamel(t *testing.T) {
	tests := []struct {
		name         string
		acronymKey   string
		acronymValue string
		expected     string
	}{
		{
			name:         "API Custom Acronym",
			acronymKey:   "API",
			acronymValue: "api",
			expected:     "Api",
		},
		{
			name:         "ABCDACME Custom Acroynm",
			acronymKey:   "ABCDACME",
			acronymValue: "AbcdAcme",
			expected:     "AbcdAcme",
		},
		{
			name:         "PostgreSQL Custom Acronym",
			acronymKey:   "PostgreSQL",
			acronymValue: "PostgreSQL",
			expected:     "PostgreSQL",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ConfigureAcronym(test.acronymKey, test.acronymValue)
			if result := ToCamel(test.acronymKey); result != test.expected {
				t.Errorf("expected custom acronym result %s, got %s", test.expected, result)
			}
		})
	}
}

func TestCustomAcronymsToLowerCamel(t *testing.T) {
	tests := []struct {
		name         string
		acronymKey   string
		acronymValue string
		expected     string
	}{
		{
			name:         "API Custom Acronym",
			acronymKey:   "API",
			acronymValue: "api",
			expected:     "api",
		},
		{
			name:         "ABCDACME Custom Acroynm",
			acronymKey:   "ABCDACME",
			acronymValue: "AbcdAcme",
			expected:     "abcdAcme",
		},
		{
			name:         "PostgreSQL Custom Acronym",
			acronymKey:   "PostgreSQL",
			acronymValue: "PostgreSQL",
			expected:     "postgreSQL",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ConfigureAcronym(test.acronymKey, test.acronymValue)
			if result := ToLowerCamel(test.acronymKey); result != test.expected {
				t.Errorf("expected custom acronym result %s, got %s", test.expected, result)
			}
		})
	}
}

func BenchmarkToLowerCamel(b *testing.B) {
	benchmarkCamelTest(b, toLowerCamel)
}

func benchmarkCamelTest(b *testing.B, fn func(testing.TB)) {
	for n := 0; n < b.N; n++ {
		fn(b)
	}
}

// TestCamelCase 将字符串转换为驼峰式字符串，非字母和数字将被忽略。
func TestCamelCase(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "test1",
			args: "",
			want: "",
		},
		{
			name: "test2",
			args: "foobar",
			want: "foobar",
		},
		{
			name: "test3",
			args: "&FOO:BAR$BAZ",
			want: "fooBarBaz",
		},
		{
			name: "test4",
			args: "$foo%",
			want: "foo",
		},
		{
			name: "test5",
			args: "Foo-#1😄$_%^&*(1bar",
			want: "foo11Bar",
		},
		{
			name: "test6",
			args: "convert_test.go",
			want: "convertTestGo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expect := strutil.CamelCase(tt.args)
			assert.Equal(t, tt.want, expect)
		})
	}
}

// TestKebabCase 将字符串转换为短横线分隔（kebab - case）字符串，非字母和数字将被忽略。
func TestKebabCase(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "test1",
			args: "",
			want: "",
		},
		{
			name: "test2",
			args: "foo-bar",
			want: "foo-bar",
		},
		{
			name: "test3",
			args: "Foo Bar-",
			want: "foo-bar",
		},
		{
			name: "test4",
			args: "FOOBAR",
			want: "foobar",
		},
		{
			name: "test5",
			args: "Foo-#1😄$_%^&*(1bar",
			want: "foo-1-1-bar",
		},
		{
			name: "test6",
			args: "convertTestGo",
			want: "convert-test-go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expect := strutil.KebabCase(tt.args)
			assert.Equal(t, tt.want, expect)
		})
	}
}

// TestUpperKebabCase 将字符串转换为大写短横线分隔（UPPER KEBAB - CASE），非字母和数字将被忽略。
func TestUpperKebabCase(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "test1",
			args: "",
			want: "",
		},
		{
			name: "test2",
			args: "foo-bar",
			want: "FOO-BAR",
		},
		{
			name: "test3",
			args: "Foo Bar-",
			want: "FOO-BAR",
		},
		{
			name: "test4",
			args: "FooBAR",
			want: "FOO-BAR",
		},
		{
			name: "test5",
			args: "Foo-#1😄$_%^&*(1bar",
			want: "FOO-1-1-BAR",
		},
		{
			name: "test6",
			args: "convertTestGo",
			want: "CONVERT-TEST-GO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expect := strutil.UpperKebabCase(tt.args)
			assert.Equal(t, tt.want, expect)
		})
	}
}

// TestSnakeCase 将字符串转换为蛇形命名法（snake_case），非字母和数字的字符将被忽略。
func TestSnakeCase(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "test1",
			args: "",
			want: "",
		},
		{
			name: "test2",
			args: "foo-bar",
			want: "foo_bar",
		},
		{
			name: "test3",
			args: "Foo Bar-",
			want: "foo_bar",
		},
		{
			name: "test4",
			args: "FOOBAR",
			want: "foobar",
		},
		{
			name: "test5",
			args: "Foo-#1😄$_%^&*(1bar",
			want: "foo_1_1_bar",
		},
		{
			name: "test6",
			args: "convertTestGo",
			want: "convert_test_go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expect := strutil.SnakeCase(tt.args)
			assert.Equal(t, tt.want, expect)
		})
	}
}

// TestUpperSnakeCase 将字符串转换为全大写蛇形命名法，非字母和数字将被忽略。
func TestUpperSnakeCase(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "test1",
			args: "",
			want: "",
		},
		{
			name: "test2",
			args: "foo-bar",
			want: "FOO_BAR",
		},
		{
			name: "test3",
			args: "Foo Bar-",
			want: "FOO_BAR",
		},
		{
			name: "test4",
			args: "FooBAR",
			want: "FOO_BAR",
		},
		{
			name: "test5",
			args: "Foo-#1😄$_%^&*(1bar",
			want: "FOO_1_1_BAR",
		},
		{
			name: "test6",
			args: "convertTestGo",
			want: "CONVERT_TEST_GO",
		},
		{
			name: "test7",
			args: "convertTest.Go",
			want: "CONVERT_TEST_GO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expect := strutil.UpperSnakeCase(tt.args)
			assert.Equal(t, tt.want, expect, tt.name)
		})
	}
}
