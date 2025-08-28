package nethttp

import (
	"fmt"
	"net/http"
	"net/textproto"
	"strings"
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

// CloneHeader clone http header
func CloneHeader(src http.Header) http.Header {
	tar := http.Header{}
	for key := range src {
		tar.Set(key, src.Get(key))
	}

	return tar
}

// CopyHeader copy http header into target
func CopyHeader(src http.Header, target http.Header) {
	for key := range src {
		target.Set(key, src.Get(key))
	}
}

func GetRid(header http.Header) string {
	return header.Get(RidHeader)
}

func SetRid(header http.Header, value string) {
	header.Set(RidHeader, value)
}

// AddRid add request id to http header
func AddRid(header http.Header, value string) {
	if GetRid(header) != value {
		header.Add(RidHeader, value)
	}
}

func GetUser(header http.Header) string {
	return header.Get(UserHeader)
}

// GetLanguage get language from http header
func GetLanguage(header http.Header) string {
	return header.Get(LanguageHeader)
}

func SetLanguage(header http.Header, value string) {
	header.Set(LanguageHeader, value)
}

// GetSupplierAccount get supplier account from http header
func GetSupplierAccount(header http.Header) string {
	return header.Get(SupplierAccountHeader)
}

func SetSupplierAccount(header http.Header, value string) {
	header.Set(SupplierAccountHeader, value)
}

func AddUser(header http.Header, value string) {
	if GetUser(header) != value {
		header.Add(UserHeader, value)
	}
}

func SetUser(header http.Header, value string) {
	header.Set(UserHeader, value)
}

func AddLanguage(header http.Header, value string) {
	if GetLanguage(header) != value {
		header.Add(LanguageHeader, value)
	}
}

func AddSupplierAccount(header http.Header, value string) {
	if GetSupplierAccount(header) != value {
		header.Add(SupplierAccountHeader, value)
	}
}

// IsReqFromWeb check if request is from web server
func IsReqFromWeb(header http.Header) bool {
	return header.Get(ReqFromWebHeader) == "true"
}

func SetReqFromWeb(header http.Header) {
	header.Set(ReqFromWebHeader, "true")
}

func GetBkJWT(header http.Header) string {
	return header.Get(BkJWTHeader)
}

func SetBkJWT(header http.Header, value string) {
	header.Set(BkJWTHeader, value)
}

func GetAppCode(header http.Header) string {
	return header.Get(AppCodeHeader)
}

func SetAppCode(header http.Header, value string) {
	header.Set(AppCodeHeader, value)
}

func GetUserToken(header http.Header) string {
	return header.Get(UserTokenHeader)
}

func GetUserTicket(header http.Header) string {
	return header.Get(UserTicketHeader)
}

func SetBkAuth(header http.Header, value string) http.Header {
	header.Set(BkAuthHeader, value)
	return header
}

func SetUserToken(header http.Header, value string) {
	header.Set(UserTokenHeader, value)
}

func SetUserTicket(header http.Header, value string) {
	header.Set(UserTicketHeader, value)
}

// GetReqRealIP get request real ip from http header
func GetReqRealIP(header http.Header) string {
	return header.Get(ReqRealIPHeader)
}

func SetReqRealIP(header http.Header, value string) {
	header.Set(ReqRealIPHeader, value)
}

// IsInnerReq check if request is inner request
func IsInnerReq(header http.Header) bool {
	return header.Get(IsInnerReqHeader) == "true"
}

// SetIsInnerReqHeader set the request is inner flag to http header
func SetIsInnerReqHeader(header http.Header) {
	header.Set(IsInnerReqHeader, "true")
}
