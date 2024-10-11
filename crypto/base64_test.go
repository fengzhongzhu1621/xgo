package crypto

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func StdBase64(s string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(s))
}

func TestBase64(t *testing.T) {
	msg := []byte("Hello world！ URL safe 编码，相当于替换掉字符串中的特殊字符，+ 和 /")

	// 1、标准编码
	encoded := base64.StdEncoding.EncodeToString(msg)
	assert.Equal(t, "SGVsbG8gd29ybGTvvIEgVVJMIHNhZmUg57yW56CB77yM55u45b2T5LqO5pu/5o2i5o6J5a2X56ym5Liy5Lit55qE54m55q6K5a2X56ym77yMKyDlkowgLw==", encoded)
	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	assert.Equal(t, msg, decoded)

	// 2、常规编码，末尾不补 =
	encoded = base64.RawStdEncoding.EncodeToString(msg)
	assert.Equal(t, "SGVsbG8gd29ybGTvvIEgVVJMIHNhZmUg57yW56CB77yM55u45b2T5LqO5pu/5o2i5o6J5a2X56ym5Liy5Lit55qE54m55q6K5a2X56ym77yMKyDlkowgLw", encoded)
	decoded, _ = base64.RawStdEncoding.DecodeString(encoded)
	assert.Equal(t, msg, decoded)

	// 3、URL safe 编码	, 替换掉字符串中的特殊字符，+ 和 /
	encoded = base64.URLEncoding.EncodeToString(msg)
	assert.Equal(t, "SGVsbG8gd29ybGTvvIEgVVJMIHNhZmUg57yW56CB77yM55u45b2T5LqO5pu_5o2i5o6J5a2X56ym5Liy5Lit55qE54m55q6K5a2X56ym77yMKyDlkowgLw==", encoded)
	decoded, _ = base64.URLEncoding.DecodeString(encoded)
	assert.Equal(t, msg, decoded)

	// 4、URL safe 编码, 替换掉字符串中的特殊字符，+ 和 /，末尾不补 =
	encoded = base64.RawURLEncoding.EncodeToString(msg)
	assert.Equal(t, "SGVsbG8gd29ybGTvvIEgVVJMIHNhZmUg57yW56CB77yM55u45b2T5LqO5pu_5o2i5o6J5a2X56ym5Liy5Lit55qE54m55q6K5a2X56ym77yMKyDlkowgLw", encoded)
	decoded, _ = base64.RawURLEncoding.DecodeString(encoded)
	assert.Equal(t, msg, decoded)
}
