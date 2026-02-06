package net

import "testing"

type dummy struct{}

func (d dummy) Name() string                                       { return "d" }
func (d dummy) Setup(containerID string, pid int) (string, error)  { return "10.0.0.2", nil }

func TestRegistry(t *testing.T) {
	Register(dummy{})
	p := Get("d")
	if p == nil || p.Name() != "d" {
		t.Fatalf("bad registry")
	}
}
