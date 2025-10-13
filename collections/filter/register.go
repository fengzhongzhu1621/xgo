package filter

import "sync"

var (
	lock          = sync.RWMutex{}
	serverFilters = make(map[string]ServerFilter)
	clientFilters = make(map[string]ClientFilter)
)

// Register registers server/client filters by name.
func Register(name string, s ServerFilter, c ClientFilter) {
	lock.Lock()
	defer lock.Unlock()

	serverFilters[name] = s
	clientFilters[name] = c
}

// GetServer gets the ServerFilter by name.
func GetServer(name string) ServerFilter {
	lock.RLock()
	f := serverFilters[name]
	lock.RUnlock()

	return f
}

// GetClient gets the ClientFilter by name.
func GetClient(name string) ClientFilter {
	lock.RLock()
	f := clientFilters[name]
	lock.RUnlock()

	return f
}
