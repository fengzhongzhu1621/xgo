package csv

import (
	"fmt"
	"os"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

// Reads file content into slice.
// func ReadCsvFile(filepath string, delimiter ...rune) ([][]string, error)
func TestReadCsvFile(t *testing.T) {
	fname := "./test.csv"
	fileutil.CreateFile(fname)

	f, _ := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC, 0o777)
	defer f.Close()

	f.WriteString("Bob, 12, male\nDuke, 14, male\nLucy, 16, female")

	content, err := fileutil.ReadCsvFile(fname)

	fmt.Println(content)
	fmt.Println(err)

	// Output:
	// [[Bob  12  male] [Duke  14  male] [Lucy  16  female]]
	// <nil>

	os.Remove(fname)
}
