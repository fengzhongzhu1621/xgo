package remote

// 远程服务器的返回.
type RemoteResponse struct {
	Value []byte
	Error error
}
