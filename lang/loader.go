package lang

import (
	"embed"
	"errors"
)

//go:embed *.json
var f embed.FS

// LoadLocaleText will load the translations from the locale json files according to the locale
// 从当前目录读取静态文件的内容
func LoadLocaleText(l string) ([]byte, error) {
	// 读取静态文件内容
	lText, err := f.ReadFile(l + ".json")

	if err != nil {
		return nil, errors.New("not able to read the lang file:" + err.Error())
	}

	return lText, nil
}
