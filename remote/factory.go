package remote

import "io"

// IRemoteConfigFactory 远程key/value存储的客户端
type IRemoteConfigFactory interface {
	Get(rp IRemoteConfig) (io.Reader, error)
	Watch(rp IRemoteConfig) (io.Reader, error)
	WatchChannel(rp IRemoteConfig) (<-chan *RemoteResponse, chan bool)
}
