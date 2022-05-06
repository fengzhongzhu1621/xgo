package ini

import (
	"bytes"
	"sort"
	"strings"

	"xgo/utils"

	"github.com/spf13/cast"
	"gopkg.in/ini.v1"
)

// LoadOptions contains all customized options used for load data source(s).
// This type is added here for convenience: this way consumers can import a single package called "ini".
type LoadOptions = ini.LoadOptions

// Codec implements the encoding.Encoder and encoding.Decoder interfaces for INI encoding.
type Codec struct {
	KeyDelimiter string      // 字典key之间的分隔符
	LoadOptions  LoadOptions // 读取init文件的初始化配置
}

func (c Codec) Encode(v map[string]interface{}) ([]byte, error) {
	cfg := ini.Empty()
	ini.PrettyFormat = false

	// 摊平字典
	flattened := utils.FlattenAndMergeMap(nil, v, "", c.keyDelimiter())

	// 获得字典的key
	keys := make([]string, 0, len(flattened))
	for key := range flattened {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		sectionName, keyName := "", key
		// 查询字符 . 最后一次出现的位置
		lastSep := strings.LastIndex(key, ".")
		if lastSep != -1 {
			sectionName = key[:(lastSep)]
			keyName = key[(lastSep + 1):]
		}

		// TODO: is this a good idea?
		if sectionName == "default" {
			sectionName = ""
		}
		// 设置key的值
		cfg.Section(sectionName).Key(keyName).SetValue(cast.ToString(flattened[key]))
	}

	// 将字典转换为ini配置文件的内容
	var buf bytes.Buffer
	_, err := cfg.WriteTo(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c Codec) Decode(b []byte, v map[string]interface{}) error {
	// 加载 ini 配置内容
	cfg := ini.Empty(c.LoadOptions)
	err := cfg.Append(b)
	if err != nil {
		return err
	}

	// 遍历所有的配置段
	sections := cfg.Sections()
	for i := 0; i < len(sections); i++ {
		// 获得每个段中的所有key值
		section := sections[i]
		keys := section.Keys()

		for j := 0; j < len(keys); j++ {
			// 获得key对应的值（字符串格式）
			key := keys[j]
			value := cfg.Section(section.Name()).Key(key.Name()).String()
			// 生成深度遍历字典
			deepestMap := utils.DeepSearch(v, strings.Split(section.Name(), c.keyDelimiter()))
			// set innermost value
			deepestMap[key.Name()] = value
		}
	}

	return nil
}

func (c Codec) keyDelimiter() string {
	if c.KeyDelimiter == "" {
		return "."
	}

	return c.KeyDelimiter
}
