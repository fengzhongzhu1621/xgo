package bcrypt

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	hashValue, _ := Encrypt("123456")
	fmt.Println(hashValue) // +$2a$10$.ovap27/i/VgBwbAz1DXF.QlrtZhCbV9LbLW5sS5/RqB36r9/b6OG
}
