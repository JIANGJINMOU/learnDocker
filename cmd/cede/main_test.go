package main

import (
	"os"
	"testing"
)

func TestUsage(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stderr
	os.Stderr = w
	usage()
	w.Close()
	os.Stderr = old
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	if n == 0 {
		t.Fatalf("no usage output")
	}
}
