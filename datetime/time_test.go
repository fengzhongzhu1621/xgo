package datetime

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeStampToLocalString(t *testing.T) {
	timestamp := int64(1633093200)

	str := TimeStampToLocalString(timestamp, time.RFC3339)
	fmt.Println(str)
}
