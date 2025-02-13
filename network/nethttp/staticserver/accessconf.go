package staticserver

import (
	"net/http"
	"regexp"

	"github.com/codeskyblue/openid-go"
	"github.com/gorilla/sessions"
)

var (
	// NonceStore 创建一个新的SimpleNonceStore实例。SimpleNonceStore是一个用于存储和检索随机数（nonce）的简单容器
	// 在密码学中，nonce 是一个只使用一次的数值，通常与密码散列函数一起使用，以确保先前的散列不会在未来被重复使用。这对于防止重放攻击特别重要。
	NonceStore = openid.NewSimpleNonceStore()
	// DiscoveryCache 用于创建一个新的SimpleDiscoveryCache实例。SimpleDiscoveryCache是一个用于缓存OpenID提供者元数据的简单容器。
	// 这些元数据包括提供者的端点URL、公钥等，用于在OpenID认证过程中进行验证和授权。
	DiscoveryCache = openid.NewSimpleDiscoveryCache()
	// Store 用于创建一个新的基于 cookie 的 session 存储引擎。
	// 接受一个密钥（key）切片作为参数，这些密钥用于对 cookie 数据进行签名和加密，以确保数据的安全性和完整性。
	Store              = sessions.NewCookieStore([]byte("something-very-secret"))
	DefaultSessionName = "ghs-session"
)

type UserInfo struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	NickName string `json:"nickName"`
}

// 用于判断登录用户的权限
type UserControl struct {
	Email string
	// Access bool
	Upload bool
	Delete bool
	Token  string
}

// 根据正则表达式判断文件的权限
type AccessTable struct {
	Regex string `yaml:"regex"`
	Allow bool   `yaml:"allow"`
}

// 配置访问类
type AccessConf struct {
	// 非登录用户的上传文件的权限
	Upload bool `yaml:"upload" json:"upload"`
	// 非登录用户的的删除权限
	Delete bool          `yaml:"delete" json:"delete"`
	Users  []UserControl `yaml:"users" json:"users"`
	// 根据正则表达式判断是否有某个文件的访问权限
	AccessTables []AccessTable `yaml:"accessTables"`
}

// 正则表达式缓冲
var reCache = make(map[string]*regexp.Regexp)

// 判断文件权限，与用户无关；使用正则表达式匹配文件名，确认文件是否可以被访问
func (c *AccessConf) CanAccess(fileName string) bool {
	for _, table := range c.AccessTables {
		// 缓存正则表达式
		pattern, ok := reCache[table.Regex]
		if !ok {
			pattern, _ = regexp.Compile(table.Regex)
			reCache[table.Regex] = pattern
		}
		// skip wrong format regex
		if pattern == nil {
			continue
		}
		// MatchString 函数来进行正则表达式匹配。
		// MatchString 函数接收两个参数：一个正则表达式字符串和一个待匹配的字符串，
		// 如果待匹配的字符串符合正则表达式的规则，该函数就会返回 true，否则返回 false。
		if pattern.MatchString(fileName) {
			return table.Allow
		}
	}
	return true
}

// CanDelete 判断登录用户是否有删除权限
func (c *AccessConf) CanDelete(r *http.Request) bool {
	// 用户没有登录则返回默认权限
	session, err := Store.Get(r, DefaultSessionName)
	if err != nil {
		return c.Delete
	}
	val := session.Values["user"]
	if val == nil {
		return c.Delete
	}

	// 用户已登录，则返回用户自定义的权限
	userInfo := val.(*UserInfo)
	for _, rule := range c.Users {
		if rule.Email == userInfo.Email {
			return rule.Delete
		}
	}
	return c.Delete
}

// 判断是否有根据 token 上传的权限
func (c *AccessConf) CanUploadByToken(token string) bool {
	for _, rule := range c.Users {
		if rule.Token == token {
			return rule.Upload
		}
	}
	return c.Upload
}

// CanUpload 判断登录用户是否有上传权限
func (c *AccessConf) CanUpload(r *http.Request) bool {
	token := r.FormValue("token")
	if token != "" {
		// 判断是否有根据 token 上传的权限
		return c.CanUploadByToken(token)
	}

	// 用户没有登录则返回默认权限
	session, err := Store.Get(r, DefaultSessionName)
	if err != nil {
		return c.Upload
	}
	val := session.Values["user"]
	if val == nil {
		return c.Upload
	}

	// 用户登录了，则返回用户的自定义权限
	userInfo := val.(*UserInfo)
	for _, rule := range c.Users {
		if rule.Email == userInfo.Email {
			return rule.Upload
		}
	}
	return c.Upload
}
