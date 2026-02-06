package storage

import "sync"

type Driver interface {
	Name() string
}

var (
	mu       sync.RWMutex
	registry = map[string]Driver{}
)

func Register(d Driver) {
	mu.Lock()
	registry[d.Name()] = d
	mu.Unlock()
}

func Get(name string) Driver {
	mu.RLock()
	defer mu.RUnlock()
	return registry[name]
}
