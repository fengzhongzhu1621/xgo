package csv

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
)

// Write content to target csv file.
// func WriteCsvFile(filepath string, records [][]string, append bool, delimiter ...rune) error
func TestWriteCsvFile(t *testing.T) {
	fpath := "./test.csv"
	fileutil.CreateFile(fpath)

	f, _ := os.OpenFile(fpath, os.O_WRONLY|os.O_TRUNC, 0o777)
	defer f.Close()

	data := [][]string{
		{"Lili", "22", "female"},
		{"Jim", "21", "male"},
	}
	err := fileutil.WriteCsvFile(fpath, data, false)
	if err != nil {
		return
	}

	content, err := fileutil.ReadCsvFile(fpath)
	if err != nil {
		return
	}
	fmt.Println(content)

	// Output:
	// [[Lili 22 female] [Jim 21 male]]
}

// Write slice of map to csv file.
// filepath: path of the CSV file.
// records: slice of maps to be written. the value of map should be basic type. The maps will be sorted by key in alphabeta order, then be written into csv file.
// appendToExistingFile: If true, data will be appended to the file if it exists.
// delimiter: Delimiter to use in the CSV file.
// headers: order of the csv column headers, needs to be consistent with the key of the map.
// func WriteMapsToCsv(filepath string, records []map[string]any, appendToExistingFile bool, delimiter rune, headers ...[]string) error
func TestWriteMapsToCsv(t *testing.T) {
	fpath := "./test.csv"
	fileutil.CreateFile(fpath)

	f, _ := os.OpenFile(fpath, os.O_WRONLY|os.O_TRUNC, 0o777)
	defer f.Close()

	records := []map[string]any{
		{"Name": "Lili", "Age": "22", "Gender": "female"},
		{"Name": "Jim", "Age": "21", "Gender": "male"},
	}

	headers := []string{"Name", "Age", "Gender"}
	err := fileutil.WriteMapsToCsv(fpath, records, false, ';', headers)
	if err != nil {
		log.Fatal(err)
	}

	content, err := fileutil.ReadCsvFile(fpath, ';')

	fmt.Println(content)

	// Output:
	// [[Name Age Gender] [Lili 22 female] [Jim 21 male]]
}
