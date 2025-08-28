package remote

import "io"

// IRemoteConfigFactory 远程key/value存储的客户端.
type IRemoteConfigFactory interface {
	Get(rp IRemoteProvider) (io.Reader, error)
	Watch(rp IRemoteProvider) (io.Reader, error)
	WatchChannel(rp IRemoteProvider) (<-chan *RemoteResponse, chan bool)
}
