package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureDirs(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	if err := EnsureDirs(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(DataRoot()); err != nil {
		t.Fatalf("data root missing: %v", err)
	}
	if _, err := os.Stat(ImagesRoot()); err != nil {
		t.Fatalf("images root missing: %v", err)
	}
	if _, err := os.Stat(ContainersRoot()); err != nil {
		t.Fatalf("containers root missing: %v", err)
	}
	if !filepath.IsAbs(DataRoot()) {
		t.Fatalf("data root not abs: %s", DataRoot())
	}
}
