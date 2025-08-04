package dotenv

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/fengzhongzhu1621/xgo/collections/maps"
	"github.com/subosito/gotenv"
)

const keyDelimiter = "_"

// Codec implements the encoding.Encoder and encoding.Decoder interfaces for encoding data containing environment variables
// (commonly called as dotenv format).
type Codec struct{}

func (Codec) Encode(v map[string]interface{}) ([]byte, error) {
	// 摊平字典，key转换为小些
	flattened := maps.FlattenAndMergeMap(nil, v, "", keyDelimiter)

	// 获得字典的key
	keys := make([]string, 0, len(flattened))
	for key := range flattened {
		keys = append(keys, key)
	}

	// 对字典的key排序
	sort.Strings(keys)

	var buf bytes.Buffer

	for _, key := range keys {
		_, err := buf.WriteString(fmt.Sprintf("%v=%v\n", strings.ToUpper(key), flattened[key]))
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (Codec) Decode(b []byte, v map[string]interface{}) error {
	var buf bytes.Buffer

	_, err := buf.Write(b)
	if err != nil {
		return err
	}

	// 解析环境变量
	env, err := gotenv.StrictParse(&buf)
	if err != nil {
		return err
	}

	for key, value := range env {
		v[key] = value
	}

	return nil
}
