package netpool

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAllocateRelease(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	p := filepath.Join(tmp, ".local", "share", "cede", "network")
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
	ip1, err := Allocate("a")
	if err != nil {
		t.Fatal(err)
	}
	ip2, err := Allocate("b")
	if err != nil {
		t.Fatal(err)
	}
	if ip1 == ip2 {
		t.Fatalf("duplicate allocation: %s", ip1)
	}
	if err := Release("a"); err != nil {
		t.Fatal(err)
	}
	ip1b, err := Allocate("a")
	if err != nil {
		t.Fatal(err)
	}
	if ip1b == "" {
		t.Fatalf("empty ip")
	}
}
