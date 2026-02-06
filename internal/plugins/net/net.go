package net

import "sync"

type NetPlugin interface {
	Name() string
	Setup(containerID string, pid int) (string, error)
}

var (
	mu       sync.RWMutex
	registry = map[string]NetPlugin{}
)

func Register(p NetPlugin) {
	mu.Lock()
	registry[p.Name()] = p
	mu.Unlock()
}

func Get(name string) NetPlugin {
	mu.RLock()
	defer mu.RUnlock()
	return registry[name]
}
