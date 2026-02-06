package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMerge(t *testing.T) {
	tmp := t.TempDir()
	p1 := filepath.Join(tmp, "a.out")
	p2 := filepath.Join(tmp, "b.out")
	data := "mode: atomic\nexample.com/containeredu/x/x.go:1.1,2.2 1 1\n"
	if err := os.WriteFile(p1, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p2, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	old := os.Args
	defer func() { os.Args = old }()
	os.Args = []string{"covermerge", p1, p2}
	main()
}
