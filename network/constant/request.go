package constant

const (
	ErrorIDKey         = "error"
	RequestIDKey       = "request_id"
	RequestIDHeaderKey = "X-Request-Id"
	ClientIDKey        = "app_code"
	ClientUsernameKey  = "username"
	UserIDKey          = "bk_uid"
	UserTokenKey       = "ticket"
)

const (
	DefaultBackendOperator = "admin"
	BackendUserKey         = "backend_user"
)

const (
	// ContextRequestIDField TODO
	ContextRequestIDField = "request_id"
	// ContextRequestUserField TODO
	ContextRequestUserField = "request_user"
	// ContextRequestOwnerField TODO
	ContextRequestOwnerField = "request_owner"

	// DefaultOwnerID the default owner value
	DefaultOwnerID         = "0"
	SystemOperatorUserName = "system"
	// ReadReferenceKey read preference key
	ReadReferenceKey = "Read_Preference"
)
