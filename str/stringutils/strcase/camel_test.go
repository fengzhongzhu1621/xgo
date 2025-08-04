package strcase

import (
	"testing"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/samber/lo"
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

func TestAllCase(t *testing.T) {
	type output struct {
		PascalCase string
		CamelCase  string
		KebabCase  string
		SnakeCase  string
	}
	name := ""
	tests := []struct {
		name   string
		input  string
		output output
	}{
		{name: name, output: output{}},
		{name: name, input: ".", output: output{}},
		{name: name, input: "Hello world!", output: output{
			PascalCase: "HelloWorld",
			CamelCase:  "helloWorld",
			KebabCase:  "hello-world",
			SnakeCase:  "hello_world",
		}},
		{name: name, input: "A", output: output{
			PascalCase: "A",
			CamelCase:  "a",
			KebabCase:  "a",
			SnakeCase:  "a",
		}},
		{name: name, input: "a", output: output{
			PascalCase: "A",
			CamelCase:  "a",
			KebabCase:  "a",
			SnakeCase:  "a",
		}},
		{name: name, input: "foo", output: output{
			PascalCase: "Foo",
			CamelCase:  "foo",
			KebabCase:  "foo",
			SnakeCase:  "foo",
		}},
		{name: name, input: "snake_case", output: output{
			PascalCase: "SnakeCase",
			CamelCase:  "snakeCase",
			KebabCase:  "snake-case",
			SnakeCase:  "snake_case",
		}},
		{name: name, input: "SNAKE_CASE", output: output{
			PascalCase: "SnakeCase",
			CamelCase:  "snakeCase",
			KebabCase:  "snake-case",
			SnakeCase:  "snake_case",
		}},
		{name: name, input: "kebab-case", output: output{
			PascalCase: "KebabCase",
			CamelCase:  "kebabCase",
			KebabCase:  "kebab-case",
			SnakeCase:  "kebab_case",
		}},
		{name: name, input: "PascalCase", output: output{
			PascalCase: "PascalCase",
			CamelCase:  "pascalCase",
			KebabCase:  "pascal-case",
			SnakeCase:  "pascal_case",
		}},
		{name: name, input: "camelCase", output: output{
			PascalCase: "CamelCase",
			CamelCase:  "camelCase",
			KebabCase:  `camel-case`,
			SnakeCase:  "camel_case",
		}},
		{name: name, input: "Title Case", output: output{
			PascalCase: "TitleCase",
			CamelCase:  "titleCase",
			KebabCase:  "title-case",
			SnakeCase:  "title_case",
		}},
		{name: name, input: "point.case", output: output{
			PascalCase: "PointCase",
			CamelCase:  "pointCase",
			KebabCase:  "point-case",
			SnakeCase:  "point_case",
		}},
		{name: name, input: "snake_case_with_more_words", output: output{
			PascalCase: "SnakeCaseWithMoreWords",
			CamelCase:  "snakeCaseWithMoreWords",
			KebabCase:  "snake-case-with-more-words",
			SnakeCase:  "snake_case_with_more_words",
		}},
		{name: name, input: "SNAKE_CASE_WITH_MORE_WORDS", output: output{
			PascalCase: "SnakeCaseWithMoreWords",
			CamelCase:  "snakeCaseWithMoreWords",
			KebabCase:  "snake-case-with-more-words",
			SnakeCase:  "snake_case_with_more_words",
		}},
		{name: name, input: "kebab-case-with-more-words", output: output{
			PascalCase: "KebabCaseWithMoreWords",
			CamelCase:  "kebabCaseWithMoreWords",
			KebabCase:  "kebab-case-with-more-words",
			SnakeCase:  "kebab_case_with_more_words",
		}},
		{name: name, input: "PascalCaseWithMoreWords", output: output{
			PascalCase: "PascalCaseWithMoreWords",
			CamelCase:  "pascalCaseWithMoreWords",
			KebabCase:  "pascal-case-with-more-words",
			SnakeCase:  "pascal_case_with_more_words",
		}},
		{name: name, input: "camelCaseWithMoreWords", output: output{
			PascalCase: "CamelCaseWithMoreWords",
			CamelCase:  "camelCaseWithMoreWords",
			KebabCase:  "camel-case-with-more-words",
			SnakeCase:  "camel_case_with_more_words",
		}},
		{name: name, input: "Title Case With More Words", output: output{
			PascalCase: "TitleCaseWithMoreWords",
			CamelCase:  "titleCaseWithMoreWords",
			KebabCase:  "title-case-with-more-words",
			SnakeCase:  "title_case_with_more_words",
		}},
		{name: name, input: "point.case.with.more.words", output: output{
			PascalCase: "PointCaseWithMoreWords",
			CamelCase:  "pointCaseWithMoreWords",
			KebabCase:  "point-case-with-more-words",
			SnakeCase:  "point_case_with_more_words",
		}},
		{name: name, input: "snake_case__with___multiple____delimiters", output: output{
			PascalCase: "SnakeCaseWithMultipleDelimiters",
			CamelCase:  "snakeCaseWithMultipleDelimiters",
			KebabCase:  "snake-case-with-multiple-delimiters",
			SnakeCase:  "snake_case_with_multiple_delimiters",
		}},
		{name: name, input: "SNAKE_CASE__WITH___multiple____DELIMITERS", output: output{
			PascalCase: "SnakeCaseWithMultipleDelimiters",
			CamelCase:  "snakeCaseWithMultipleDelimiters",
			KebabCase:  "snake-case-with-multiple-delimiters",
			SnakeCase:  "snake_case_with_multiple_delimiters",
		}},
		{name: name, input: "kebab-case--with---multiple----delimiters", output: output{
			PascalCase: "KebabCaseWithMultipleDelimiters",
			CamelCase:  "kebabCaseWithMultipleDelimiters",
			KebabCase:  "kebab-case-with-multiple-delimiters",
			SnakeCase:  "kebab_case_with_multiple_delimiters",
		}},
		{name: name, input: "Title Case  With   Multiple    Delimiters", output: output{
			PascalCase: "TitleCaseWithMultipleDelimiters",
			CamelCase:  "titleCaseWithMultipleDelimiters",
			KebabCase:  "title-case-with-multiple-delimiters",
			SnakeCase:  "title_case_with_multiple_delimiters",
		}},
		{name: name, input: "point.case..with...multiple....delimiters", output: output{
			PascalCase: "PointCaseWithMultipleDelimiters",
			CamelCase:  "pointCaseWithMultipleDelimiters",
			KebabCase:  "point-case-with-multiple-delimiters",
			SnakeCase:  "point_case_with_multiple_delimiters",
		}},
		{name: name, input: " leading space", output: output{
			PascalCase: "LeadingSpace",
			CamelCase:  "leadingSpace",
			KebabCase:  "leading-space",
			SnakeCase:  "leading_space",
		}},
		{name: name, input: "   leading spaces", output: output{
			PascalCase: "LeadingSpaces",
			CamelCase:  "leadingSpaces",
			KebabCase:  "leading-spaces",
			SnakeCase:  "leading_spaces",
		}},
		{name: name, input: "\t\t\r\n leading whitespaces", output: output{
			PascalCase: "LeadingWhitespaces",
			CamelCase:  "leadingWhitespaces",
			KebabCase:  "leading-whitespaces",
			SnakeCase:  "leading_whitespaces",
		}},
		{name: name, input: "trailing space ", output: output{
			PascalCase: "TrailingSpace",
			CamelCase:  "trailingSpace",
			KebabCase:  "trailing-space",
			SnakeCase:  "trailing_space",
		}},
		{name: name, input: "trailing spaces   ", output: output{
			PascalCase: "TrailingSpaces",
			CamelCase:  "trailingSpaces",
			KebabCase:  "trailing-spaces",
			SnakeCase:  "trailing_spaces",
		}},
		{name: name, input: "trailing whitespaces\t\t\r\n", output: output{
			PascalCase: "TrailingWhitespaces",
			CamelCase:  "trailingWhitespaces",
			KebabCase:  "trailing-whitespaces",
			SnakeCase:  "trailing_whitespaces",
		}},
		{name: name, input: " on both sides ", output: output{
			PascalCase: "OnBothSides",
			CamelCase:  "onBothSides",
			KebabCase:  "on-both-sides",
			SnakeCase:  "on_both_sides",
		}},
		{name: name, input: "    many on both sides  ", output: output{
			PascalCase: "ManyOnBothSides",
			CamelCase:  "manyOnBothSides",
			KebabCase:  "many-on-both-sides",
			SnakeCase:  "many_on_both_sides",
		}},
		{name: name, input: "\r whitespaces on both sides\t\t\r\n", output: output{
			PascalCase: "WhitespacesOnBothSides",
			CamelCase:  "whitespacesOnBothSides",
			KebabCase:  "whitespaces-on-both-sides",
			SnakeCase:  "whitespaces_on_both_sides",
		}},
		{name: name, input: "  extraSpaces in_This TestCase Of MIXED_CASES\t", output: output{
			PascalCase: "ExtraSpacesInThisTestCaseOfMixedCases",
			CamelCase:  "extraSpacesInThisTestCaseOfMixedCases",
			KebabCase:  "extra-spaces-in-this-test-case-of-mixed-cases",
			SnakeCase:  "extra_spaces_in_this_test_case_of_mixed_cases",
		}},
		{name: name, input: "CASEBreak", output: output{
			PascalCase: "CaseBreak",
			CamelCase:  "caseBreak",
			KebabCase:  "case-break",
			SnakeCase:  "case_break",
		}},
		{name: name, input: "ID", output: output{
			PascalCase: "Id",
			CamelCase:  "id",
			KebabCase:  "id",
			SnakeCase:  "id",
		}},
		{name: name, input: "userID", output: output{
			PascalCase: "UserId",
			CamelCase:  "userId",
			KebabCase:  "user-id",
			SnakeCase:  "user_id",
		}},
		{name: name, input: "JSON_blob", output: output{
			PascalCase: "JsonBlob",
			CamelCase:  "jsonBlob",
			KebabCase:  "json-blob",
			SnakeCase:  "json_blob",
		}},
		{name: name, input: "HTTPStatusCode", output: output{
			PascalCase: "HttpStatusCode",
			CamelCase:  "httpStatusCode",
			KebabCase:  "http-status-code",
			SnakeCase:  "http_status_code",
		}},
		{name: name, input: "FreeBSD and SSLError are not golang initialisms", output: output{
			PascalCase: "FreeBsdAndSslErrorAreNotGolangInitialisms",
			CamelCase:  "freeBsdAndSslErrorAreNotGolangInitialisms",
			KebabCase:  "free-bsd-and-ssl-error-are-not-golang-initialisms",
			SnakeCase:  "free_bsd_and_ssl_error_are_not_golang_initialisms",
		}},
		{name: name, input: "David's Computer", output: output{
			PascalCase: "DavidSComputer",
			CamelCase:  "davidSComputer",
			KebabCase:  "david-s-computer",
			SnakeCase:  "david_s_computer",
		}},
		{name: name, input: "http200", output: output{
			PascalCase: "Http200",
			CamelCase:  "http200",
			KebabCase:  "http-200",
			SnakeCase:  "http_200",
		}},
		{name: name, input: "NumberSplittingVersion1.0r3", output: output{
			PascalCase: "NumberSplittingVersion10R3",
			CamelCase:  "numberSplittingVersion10R3",
			KebabCase:  "number-splitting-version-1-0-r3",
			SnakeCase:  "number_splitting_version_1_0_r3",
		}},
		{name: name, input: "When you have a comma, odd results", output: output{
			PascalCase: "WhenYouHaveACommaOddResults",
			CamelCase:  "whenYouHaveACommaOddResults",
			KebabCase:  "when-you-have-a-comma-odd-results",
			SnakeCase:  "when_you_have_a_comma_odd_results",
		}},
		{name: name, input: "Ordinal numbers work: 1st 2nd and 3rd place", output: output{
			PascalCase: "OrdinalNumbersWork1St2NdAnd3RdPlace",
			CamelCase:  "ordinalNumbersWork1St2NdAnd3RdPlace",
			KebabCase:  "ordinal-numbers-work-1-st-2-nd-and-3-rd-place",
			SnakeCase:  "ordinal_numbers_work_1_st_2_nd_and_3_rd_place",
		}},
		{name: name, input: "BadUTF8\xe2\xe2\xa1", output: output{
			PascalCase: "BadUtf8",
			CamelCase:  "badUtf8",
			KebabCase:  "bad-utf-8",
			SnakeCase:  "bad_utf_8",
		}},
		{name: name, input: "IDENT3", output: output{
			PascalCase: "Ident3",
			CamelCase:  "ident3",
			KebabCase:  "ident-3",
			SnakeCase:  "ident_3",
		}},
		{name: name, input: "LogRouterS3BucketName", output: output{
			PascalCase: "LogRouterS3BucketName",
			CamelCase:  "logRouterS3BucketName",
			KebabCase:  "log-router-s3-bucket-name",
			SnakeCase:  "log_router_s3_bucket_name",
		}},
		{name: name, input: "PINEAPPLE", output: output{
			PascalCase: "Pineapple",
			CamelCase:  "pineapple",
			KebabCase:  "pineapple",
			SnakeCase:  "pineapple",
		}},
		{name: name, input: "Int8Value", output: output{
			PascalCase: "Int8Value",
			CamelCase:  "int8Value",
			KebabCase:  "int-8-value",
			SnakeCase:  "int_8_value",
		}},
		{name: name, input: "first.last", output: output{
			PascalCase: "FirstLast",
			CamelCase:  "firstLast",
			KebabCase:  "first-last",
			SnakeCase:  "first_last",
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pascal := lo.PascalCase(test.input)
			if pascal != test.output.PascalCase {
				t.Errorf(
					"PascalCase(%q) = %q; expected %q",
					test.input,
					pascal,
					test.output.PascalCase,
				)
			}
			camel := lo.CamelCase(test.input)
			if camel != test.output.CamelCase {
				t.Errorf(
					"CamelCase(%q) = %q; expected %q",
					test.input,
					camel,
					test.output.CamelCase,
				)
			}
			kebab := lo.KebabCase(test.input)
			if kebab != test.output.KebabCase {
				t.Errorf(
					"KebabCase(%q) = %q; expected %q",
					test.input,
					kebab,
					test.output.KebabCase,
				)
			}
			snake := lo.SnakeCase(test.input)
			if snake != test.output.SnakeCase {
				t.Errorf(
					"SnakeCase(%q) = %q; expected %q",
					test.input,
					snake,
					test.output.SnakeCase,
				)
			}
		})
	}
}
