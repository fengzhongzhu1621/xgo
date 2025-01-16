package randutils

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/samber/lo"

	"github.com/stretchr/testify/assert"

	"github.com/duke-git/lancet/v2/random"
)

// Generate random given length string. only contains letter (a-zA-Z)
// func RandString(length int) string
func TestRandomRandString(t *testing.T) {
	{
		randStr := random.RandString(6)
		fmt.Println(randStr) // ACMbJt
	}

	{
		t.Parallel()
		is := assert.New(t)

		rand.Seed(time.Now().UnixNano())

		str1 := lo.RandomString(100, lo.LowerCaseLettersCharset)
		is.Equal(100, lo.RuneLength(str1))
		is.Subset(lo.LowerCaseLettersCharset, []rune(str1))

		str2 := lo.RandomString(100, lo.LowerCaseLettersCharset)
		is.NotEqual(str1, str2)

		noneUtf8Charset := []rune("明1好休2林森")
		str3 := lo.RandomString(100, noneUtf8Charset)
		is.Equal(100, lo.RuneLength(str3))
		is.Subset(noneUtf8Charset, []rune(str3))

		is.PanicsWithValue("lo.RandomString: Charset parameter must not be empty", func() { lo.RandomString(100, []rune{}) })
		is.PanicsWithValue("lo.RandomString: Size parameter must be greater than 0", func() { lo.RandomString(0, lo.LowerCaseLettersCharset) })
	}
}

// TestRandStringSlice Generate a slice of random string of length strLen based on charset.
// charset should be one of the following:
// random.Numeral, random.LowwerLetters, random.UpperLetters random.Letters, random.SymbolChars, random.AllChars. or a combination of them.
// func RandStringSlice(charset string, sliceLen, strLen int) []string
func TestRandStringSlice(t *testing.T) {
	strs := random.RandStringSlice(random.Letters, 4, 6)
	fmt.Println(strs)

	// output random string slice like below:
	//[CooSMq RUFjDz FAeMPf heRyGv]
}

// TestRandString Generate a random upper case string
// func RandUpper(length int) string
func TestRandUpper(t *testing.T) {
	randStr := random.RandUpper(6)
	fmt.Println(randStr)
}

// TestRandLower Generate a random lower case string
// func RandLower(length int) string
func TestRandLower(t *testing.T) {
	randStr := random.RandLower(6)
	fmt.Println(randStr)
}

// TestRandNumeral Generate a random numeral string
// func RandNumeral(length int) string
func TestRandNumeral(t *testing.T) {
	randStr := random.RandNumeral(6)
	fmt.Println(randStr)
}

// TestRandNumeralOrLetter generate a random numeral or letter string
// func RandNumeralOrLetter(length int) string
func TestRandNumeralOrLetter(t *testing.T) {
	randStr := random.RandNumeralOrLetter(6)
	fmt.Println(randStr) //0aW7cQ
}

// TestRandSymbolChar Generate a random symbol char of specified length.
// func RandSymbolChar(length int) string
func TestRandSymbolChar(t *testing.T) {
	randStr := random.RandSymbolChar(6)
	fmt.Println(randStr) // random string like: @#(_")
}
