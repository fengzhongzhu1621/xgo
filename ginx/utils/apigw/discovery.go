package apigw

import (
	"errors"
	"sync"
)

type ApiGWBaseResponse struct {
	Result  bool   `json:"result"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ApiGWDiscovery api gateway discovery struct
type ApiGWDiscovery struct {
	Servers []string
	index   int
	sync.Mutex
}

// GetServers get api gateway server
func (s *ApiGWDiscovery) GetServers() ([]string, error) {
	s.Lock()
	defer s.Unlock()

	num := len(s.Servers)
	if num == 0 {
		return []string{}, errors.New("oops, there is no server can be used")
	}

	if s.index < num-1 {
		s.index = s.index + 1
		return append(s.Servers[s.index-1:], s.Servers[:s.index-1]...), nil
	}

	s.index = 0
	return append(s.Servers[num-1:], s.Servers[:num-1]...), nil
}

// GetServersChan get api gateway chan
func (s *ApiGWDiscovery) GetServersChan() chan []string {
	return nil
}
