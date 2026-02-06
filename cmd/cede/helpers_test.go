package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIOCopyAndCopyFile(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src.txt")
	dst := filepath.Join(tmp, "dst.txt")
	if err := os.WriteFile(src, []byte("abc123"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := copyFile(src, dst, 0o644); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "abc123" {
		t.Fatalf("content mismatch: %q", string(b))
	}
}

func TestNetConfigAndList(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	os.Setenv("USERPROFILE", tmp)
	p := filepath.Join(tmp, ".local", "share", "cede", "network")
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := netConfig("10.3.0.0/24", "10.3.0.1"); err != nil {
		t.Fatal(err)
	}
	ass := netpoolList()
	if len(ass) != 0 {
		t.Fatalf("unexpected assignments")
	}
}
