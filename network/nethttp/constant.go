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
	RidHeader    = "X-Bkapi-Request-Id"
	BkAuthHeader = "X-Bkapi-Authorization"
	BkJWTHeader  = "X-Bkapi-JWT"
)

const (
	BKHTTPHeaderUser      = "BK_User"
	UserHeader            = "X-Xgo-User"
	LanguageHeader        = "X-Xgo-Language"
	SupplierAccountHeader = "X-Xgo-Supplier-Account"
	ReqFromWebHeader      = "X-Xgo-Request-From-Web"
	AppCodeHeader         = "X-Xgo-App-Code"
	UserTokenHeader       = "X-Xgo-User-Token"
	UserTicketHeader      = "X-Xgo-User-Ticket"
)

const (
	Success      = 0
	SuccessStr   = "success"
	NoPermission = 9900403
)
