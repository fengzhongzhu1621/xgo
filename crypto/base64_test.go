package crypto

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestBase64(t *testing.T) {
    msg := []byte("Hello world！ URL safe 编码，相当于替换掉字符串中的特殊字符，+ 和 /")

    // 1、标准编码
    encoded := base64.StdEncoding.EncodeToString(msg)
    fmt.Println(encoded)
    decoded, _ := base64.StdEncoding.DecodeString(encoded)
    fmt.Println(string(decoded))


    // 2、常规编码，末尾不补 =
    encoded = base64.RawStdEncoding.EncodeToString(msg)
    fmt.Println(encoded)
    decoded, _ = base64.RawStdEncoding.DecodeString(encoded)
    fmt.Println(string(decoded))

    // 3、URL safe 编码	, 替换掉字符串中的特殊字符，+ 和 /
    encoded = base64.URLEncoding.EncodeToString(msg)
    fmt.Println(encoded)
    decoded, _ = base64.URLEncoding.DecodeString(encoded)
    fmt.Println(string(decoded))

	// 4、URL safe 编码, 替换掉字符串中的特殊字符，+ 和 /，末尾不补 =
    encoded = base64.RawURLEncoding.EncodeToString(msg)
    fmt.Println(encoded)
    decoded, _ = base64.RawURLEncoding.DecodeString(encoded)
    fmt.Println(string(decoded))
}