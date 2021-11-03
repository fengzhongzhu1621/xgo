package cache

type Cache interface {
	Parse() error
	Scan() (map[string][]string, error)
	Get(key string) string
	GetMany(key string) []string
	Set(key string, value string)
	SetMany(key string, value []string)
}
