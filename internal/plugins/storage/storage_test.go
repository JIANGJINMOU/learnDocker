package storage

import "testing"

type sdriver struct{}

func (s sdriver) Name() string { return "s" }

func TestRegistry(t *testing.T) {
	Register(sdriver{})
	d := Get("s")
	if d == nil || d.Name() != "s" {
		t.Fatalf("bad registry")
	}
}
