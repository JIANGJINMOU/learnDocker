//go:build linux

package overlay

import "testing"

func TestJoinLower(t *testing.T) {
	dirs := []string{"/a", "/b", "/c"}
	s := joinLower(dirs)
	if s != "/a:/b:/c" {
		t.Fatalf("unexpected join: %s", s)
	}
}
