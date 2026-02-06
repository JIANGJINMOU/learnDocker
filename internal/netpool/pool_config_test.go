package netpool

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigCIDRGateway(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	p := filepath.Join(tmp, ".local", "share", "cede", "network")
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := SetCIDRGateway("10.2.0.0/24", "10.2.0.1"); err != nil {
		t.Fatal(err)
	}
	if CIDR() != "10.2.0.0/24" {
		t.Fatalf("cidr mismatch: %s", CIDR())
	}
	if Gateway() != "10.2.0.1" {
		t.Fatalf("gw mismatch: %s", Gateway())
	}
}
