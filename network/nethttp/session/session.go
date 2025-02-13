package session

import (
	"net/http"

	"github.com/fengzhongzhu1621/xgo/crypto/randutils"

	"github.com/gorilla/sessions"
)

const (
	SecureSessionName = "secure-session"
)

// SessionManager 管理会话的结构体，包含一个 store 字段，用于存储会话数据。
// 使用了 sessions.Store 接口，使其具备灵活性，可以更换不同的存储后端（如 Cookie Store、Redis Store 等）。
type SessionManager struct {
	store sessions.Store
}

func NewSessionManager(secret []byte) *SessionManager {
	return &SessionManager{
		store: sessions.NewCookieStore(secret), // 初始化了一个基于 Cookie 的会话存储
	}
}

// ConfigureStore 配置会话存储
func (sm *SessionManager) ConfigureStore() {
	store := sm.store.(*sessions.CookieStore)
	store.Options = &sessions.Options{
		Path:     "/",                     // 设置会话 Cookie 的路径为根路径，使其在整个网站有效。
		MaxAge:   3600,                    // 设置 Cookie 的最大存活时间为 3600 秒（1 小时）。
		HttpOnly: true,                    // 设置为 true，防止 JavaScript 访问该 Cookie，增加安全性。
		Secure:   true,                    // 设置为 true，确保 Cookie 只在 HTTPS 连接中传输。
		SameSite: http.SameSiteStrictMode, // 设置为 Strict，防止跨站请求伪造（CSRF）攻击。
	}
}

// StartSession 启动会话
func (sm *SessionManager) StartSession(w http.ResponseWriter, r *http.Request) (*sessions.Session, error) {
	// 获取或创建一个会话
	session, err := sm.store.Get(r, SecureSessionName)
	if err != nil {
		return nil, err
	}

	// 生成新的session ID
	session.ID = randutils.GenerateSecureID()
	return session, nil
}
