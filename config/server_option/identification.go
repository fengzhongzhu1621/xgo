package server_option

const (
	MODULE_WEBSERVER = "webserver"
	MODULE_APISERVER = "apiserver"
	MODULE_MIGRATE   = "migrate"
)

var identification = "unknown"

// SetIdentification TODO
func SetIdentification(id string) {
	if identification == "unknown" {
		identification = id
	}
}

// GetIdentification TODO
func GetIdentification() string {
	return identification
}
