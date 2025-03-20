package discovery

// NewMockDiscoveryInterface TODO
func NewMockDiscoveryInterface() IDiscoveryInterface {
	return &MockDiscovery{}
}

// MockDiscovery TODO
type MockDiscovery struct{}

// CoreService TODO
func (d *MockDiscovery) CoreService() Interface {
	return &mockServer{}
}

// IsMaster TODO
func (d *MockDiscovery) IsMaster() bool {
	return true
}

// Server TODO
func (d *MockDiscovery) Server(name string) Interface {
	return emptyServerInst
}

type mockServer struct{}

// GetServers TODO
func (*mockServer) GetServers() ([]string, error) {
	return []string{"http://127.0.0.1:8080"}, nil
}

// IsMaster TODO
func (*mockServer) IsMaster(string) bool {
	return true
}

// GetServersChan TODO
func (s *mockServer) GetServersChan() chan []string {
	return nil
}
