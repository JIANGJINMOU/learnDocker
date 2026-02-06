//go:build linux

package bridge

import (
	"sync"
	"testing"
)

func TestSetupRunsCommandsAndReturnsIP(t *testing.T) {
	var mu sync.Mutex
	var calls [][]string
	runner = func(cmd string, args ...string) error {
		mu.Lock()
		defer mu.Unlock()
		one := append([]string{cmd}, args...)
		calls = append(calls, one)
		return nil
	}
	p := Plugin{name: "bridge0", bridge: "cede0", gateway: "10.0.0.1", cidr: "10.0.0.0/24"}
	ip, err := p.Setup("abcdef012345", 9999)
	if err != nil {
		t.Fatal(err)
	}
	if ip == "" {
		t.Fatalf("empty ip")
	}
	if len(calls) == 0 {
		t.Fatalf("no commands executed")
	}
}
