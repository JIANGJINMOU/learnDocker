package id

import "testing"

func TestNew(t *testing.T) {
	s := New()
	if len(s) == 0 {
		t.Fatalf("empty uuid")
	}
}
