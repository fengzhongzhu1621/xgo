package nethttp

const (
	LocalSesionLockPrefix = "localsess"
)

const (
	WEB_LANGUAGE_CN = "zh-cn"
	WEB_LANGUAGE_EN = "en"
)

const DEFAULT_BACKEND_OPERATOR = "system"

const (
	RidHeader             = "X-Bkapi-Request-Id"
	UserHeader            = "X-Xgo-User"
	LanguageHeader        = "X-Xgo-Language"
	SupplierAccountHeader = "X-Xgo-Supplier-Account"
	// ReqFromWebHeader is the http header key that represents if request is from web server
	ReqFromWebHeader = "X-Xgo-Request-From-Web"
	// AppCodeHeader is the blueking app code http header key, its value is from jwt info
	AppCodeHeader = "X-Xgo-App-Code"

	// BKHTTPHeaderUser current request http request header fields name for login user
	BKHTTPHeaderUser = "BK_User"
	// BkJWTHeader is the blueking api gateway jwt http header key
	BkJWTHeader = "X-Bkapi-JWT"
)

const (
	Success      = 0
	SuccessStr   = "success"
	NoPermission = 9900403
)
