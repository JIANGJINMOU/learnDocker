//go:build linux

package overlay

import "testing"

func TestJoinLowerAlt(t *testing.T) {
	s := joinLower([]string{"a", "b", "c"})
	if s != "a:b:c" {
		t.Fatalf("unexpected: %s", s)
	}
}

func TestUnmountNonMount(t *testing.T) {
	err := Unmount("/tmp/not-a-mount")
	if err == nil {
		t.Fatalf("expected error")
	}
}
