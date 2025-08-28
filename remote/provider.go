package remote

// IRemoteProvider stores the configuration necessary
// to connect to a remote key/value store.
// Optional secretKeyring to unencrypt encrypted values
// can be provided.
// 远程key/value存储的客户端配置接口.
type IRemoteProvider interface {
	// Provider is a string value: "etcd", "consul" or "firestore" are currently supported.
	Provider() string
	// Endpoint is the url.  etcd requires http://ip:port  consul requires ip:port
	Endpoint() string
	// Path is the path in the k/v store to retrieve configuration
	Path() string
	// SecretKeyring is the filepath to your openpgp secret keyring.  e.g. /etc/secrets/myring.gpg
	SecretKeyring() string
}
