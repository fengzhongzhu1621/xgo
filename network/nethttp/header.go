package nethttp

import (
	"fmt"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/fengzhongzhu1621/xgo"
	"github.com/gin-gonic/gin"
)

var HeaderNewlineToSpace = strings.NewReplacer("\n", " ", "\r", " ") // 换行字符替换器

var HeaderDashToUnderscore = strings.NewReplacer("-", "_") // 短横线字符替换器

// HeaderSet http header set
type HeaderSet struct {
	Key   string
	Value string
}

// AppendEnv 构造header, 返回一个新数组.
func AppendEnv(env []string, k string, v ...string) []string {
	if len(v) == 0 {
		return env
	}

	// 创建一个字符串空数组
	vCleaned := make([]string, 0, len(v))
	// 将数组元素 v 去掉换行符和首尾的空白字符
	for _, val := range v {
		vCleaned = append(vCleaned, strings.TrimSpace(HeaderNewlineToSpace.Replace(val)))
	}
	return append(env, fmt.Sprintf("%s=%s",
		strings.ToUpper(k),
		strings.Join(vCleaned, ", ")))
}

func SplitMimeHeader(s string) (string, string) {
	p := strings.IndexByte(s, ':')
	if p < 0 {
		return s, ""
	}
	key := textproto.CanonicalMIMEHeaderKey(s[:p])

	for p++; p < len(s); p++ {
		if s[p] != ' ' {
			break
		}
	}
	return key, s[p:]
}

func PushHeaders(h http.Header, hdrs []string) {
	for _, hstr := range hdrs {
		h.Add(SplitMimeHeader(hstr))
	}
}

// GetBluekingLanguageFromHeader 从请求头获取语言
func GetBluekingLanguageFromHeader(c *gin.Context) (string, error) {
	header := c.Request.Header.Get("X-BkApi-Blueking-Language")
	if len(header) == 0 {
		return "", xgo.JwtTokenNoneErr
	}
	strs := strings.Split(header, " ")

	return strs[0], nil
}

// GetEnvFromHeader 从请求投获取 env 的值，如果找不到则从 Get 请求参数中获取
func GetEnvFromHeader(c *gin.Context) string {
	env := c.Request.Header.Get("env")
	if env == "" {
		env = c.Query("env")
	}

	return env
}

// GetUsernameFromHeader ...
func GetUsernameFromHeader(c *gin.Context) string {
	username := c.Request.Header.Get("username")
	if username == "" {
		username = c.Query("username")
	}

	return username
}

// GetTokenFromHeader ...
func GetTokenFromHeader(c *gin.Context) string {
	token := c.Request.Header.Get("token")
	if token == "" {
		token = c.Query("token")
	}

	return token
}
