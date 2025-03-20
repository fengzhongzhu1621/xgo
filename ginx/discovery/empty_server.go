package discovery

var (
	emptyServerInst = &emptyServer{}
)

// emptyServer 适配服务不存在的情况， 当服务不存在的时候，返回空的服务
type emptyServer struct {
}

// GetServers TODO
func (es *emptyServer) GetServers() ([]string, error) {
	return []string{}, nil
}

// GetServersChan TODO
func (es *emptyServer) GetServersChan() chan []string {
	return make(chan []string, 20)
}
