package emoj

import (
	"fmt"
	"testing"

	"github.com/kyokomi/emoji/v2"
)

func TestPrintln(t *testing.T) {
	emoji.Println(":beer: Beer")

}

func TestSprint(t *testing.T) {
	message := emoji.Sprint("I like a :pizza: and :sushi: !")
	fmt.Println(message)
}
