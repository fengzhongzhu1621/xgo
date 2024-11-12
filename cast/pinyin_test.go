package cast

import (
	"testing"

	"fmt"
	"strings"

	"github.com/mozillazg/go-pinyin"
)

type Response struct {
	FullPinyin string `json:"full_pinyin"`
	Initials   string `json:"initials"`
}

// ConvertToPinyin 接收中文姓名并返回全拼和拼音首字母
func ConvertToPinyin(name string) (fullPinyin, initials string) {
	// 使用默认的汉字转换选项
	args := pinyin.NewArgs()

	// 获取拼音的二维数组
	py := pinyin.Pinyin(name, args)

	var fullPinyinList []string
	var initialsList []string
	for _, syllable := range py {
		// 全拼音
		fullPinyinList = append(fullPinyinList, syllable[0])
		// 首字母
		initialsList = append(initialsList, string(syllable[0][0]))
	}

	fullPinyin = strings.Join(fullPinyinList, "")
	initials = strings.Join(initialsList, "")

	return fullPinyin, initials
}

func TestConvertToPinyin(t *testing.T) {
	name := "中国人"
	fullPinyin, initials := ConvertToPinyin(name)

	// 姓名: 中国人
	// 全拼音: zhongguoren
	// 拼音首字母: zgr
	fmt.Println("姓名:", name)
	fmt.Println("全拼音:", fullPinyin)
	fmt.Println("拼音首字母:", initials)
}
