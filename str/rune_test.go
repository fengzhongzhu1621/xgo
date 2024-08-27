package stringutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRune(t *testing.T) {

	temp := []rune{20320, 22909, 32, 19990, 30028}
	fmt.Println(string(temp)) // 你好 世界

	var str string = "hello world"
	fmt.Println("byte=", []byte(str))    // [104 101 108 108 111 32 119 111 114 108 100]
	fmt.Println("byte=", []rune(str))    // [104 101 108 108 111 32 119 111 114 108 100]
	fmt.Println(str[:2])                 // he
	fmt.Println(string([]rune(str)[:2])) // he

	var str2 string = "你好 世界"
	fmt.Println("byte=", []byte(str2)) // [228 189 160 229 165 189 32 228 184 150 231 149 140]
	fmt.Println("byte=", []rune(str2)) // [20320 22909 32 19990 30028]
	fmt.Println(str2[:2])
	fmt.Println(string([]rune(str2)[:2])) // 你好
}

func TestLen(t *testing.T) {
	s := "Hello王"
	sHello := "Hello"
	sWang := "王"
	//len()获得的是 byte 字节的数量
	assert.Equal(t, len(s), 8)
	assert.Equal(t, len(sHello), 5)
	assert.Equal(t, len(sWang), 3)
}

func TestRange(t *testing.T) {
	s := "Hello王"
	for _, c := range s {
		fmt.Printf("%c\n", c) //%c 字符
	}
	// H
	// e
	// l
	// l
	// o
	// 王
}

// TestChangeString 字符串修改是不能直接修改的，需要转成rune切片后再修改
func TestChangeString(t *testing.T) {
	s2 := "大白兔"
	s3 := []rune(s2)        // 把字符串强制转成rune切片
	s3[0] = '小'             // 注意 这里需要使用单引号的字符，而不是双引号的字符串
	fmt.Println(string(s3)) // 把rune类型的s3强转成字符串
}
