//go:build !linux

package bridge

import (
	nreg "example.com/containeredu/internal/plugins/net"
)

type Plugin struct{}

func (Plugin) Name() string                                      { return "bridge0" }
func (Plugin) Setup(containerID string, pid int) (string, error) { return "", nil }

func init() {
	nreg.Register(Plugin{})
}
