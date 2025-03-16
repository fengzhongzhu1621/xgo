package backbone

import (
	"encoding/json"
	"errors"

	"github.com/fengzhongzhu1621/xgo/db/zookeeper/registerdiscover"
	"github.com/fengzhongzhu1621/xgo/ginx/server_info"
)

// ServiceRegisterInterface TODO
type ServiceRegisterInterface interface {
	// Ping to ping server
	Ping() error
	// Register local server info, it can only be called for once.
	Register(path string, c server_info.ServerInfo) error
	// Cancel to stop server register and discover
	Cancel()
	// ClearRegisterPath to delete server register path from zk
	ClearRegisterPath() error
}

// NewServiceRegister TODO
func NewServiceRegister(client *registerdiscover.ZkClient) (ServiceRegisterInterface, error) {
	s := new(serviceRegister)
	s.client = registerdiscover.NewRegDiscoverEx(client)
	return s, nil
}

type serviceRegister struct {
	client *registerdiscover.RegDiscover
}

// Register TODO
func (s *serviceRegister) Register(path string, c server_info.ServerInfo) error {
	if c.RegisterIP == "0.0.0.0" {
		return errors.New("register ip can not be 0.0.0.0")
	}

	js, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return s.client.RegisterAndWatchService(path, js)
}

// Ping to ping server
func (s *serviceRegister) Ping() error {
	return s.client.Ping()
}

// Cancel to stop server register and discover
func (s *serviceRegister) Cancel() {
	s.client.Cancel()
}

// ClearRegisterPath to delete server register path from zk
func (s *serviceRegister) ClearRegisterPath() error {
	return s.client.ClearRegisterPath()
}
