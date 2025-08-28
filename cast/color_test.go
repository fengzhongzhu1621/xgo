package cast

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/convertor"
)

// func ColorHexToRGB(colorHex string) (red, green, blue int)
func TestColorHexToRGB(t *testing.T) {
	colorHex := "#003366"
	r, g, b := convertor.ColorHexToRGB(colorHex)

	fmt.Println(r, g, b)

	// Output:
	// 0 51 102
}

// func ColorRGBToHex(red, green, blue int) string
func TestColorRGBToHex(t *testing.T) {
	r := 0
	g := 51
	b := 102
	colorHex := convertor.ColorRGBToHex(r, g, b)

	fmt.Println(colorHex)

	// Output:
	// #003366
}
