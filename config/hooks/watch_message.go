package hooks

// WatchMessage change message
type WatchMessage struct {
	Provider  string // provider name
	Path      string // config path
	ExpandEnv bool   // expend env status
	Codec     string // codec
	Watch     bool   // status for start watch
	Value     []byte // config content diff ? 配置原始内容
	Error     error  // load error message, success is empty string
}

type WatchMessageHookFunc func(message WatchMessage)

// defaultNotifyChange default hook for notify config changed
var defaultWatchHook = func(message WatchMessage) {}

// SetDefaultWatchHook set default hook notify when config changed
func SetDefaultWatchHook(f WatchMessageHookFunc) {
	defaultWatchHook = f
}

func GetWatchMessageHook() WatchMessageHookFunc {
	return defaultWatchHook
}
