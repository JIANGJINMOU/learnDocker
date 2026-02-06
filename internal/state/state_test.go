package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"example.com/containeredu/internal/paths"
)

func TestSaveAndList(t *testing.T) {
	// redirect data root to temp by setting HOME
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	if err := paths.EnsureDirs(); err != nil {
		t.Fatal(err)
	}
	s := ContainerState{
		ID:        "abc",
		Image:     "busybox",
		Pid:       123,
		Command:   "/bin/sh",
		Args:      []string{"-c", "echo hi"},
		CreatedAt: time.Now(),
		Hostname:  "cede",
		Status:    "running",
		MountDir:  "/rootfs",
	}
	if err := Save(s); err != nil {
		t.Fatalf("save: %v", err)
	}
	p := filepath.Join(paths.ContainersRoot(), "abc.json")
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("state file missing: %v", err)
	}
	items, err := List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 1 || items[0].ID != "abc" {
		t.Fatalf("unexpected list: %+v", items)
	}
}

func TestListIgnoresNonJSON(t *testing.T) {
	if err := paths.EnsureDirs(); err != nil {
		t.Fatal(err)
	}
	root := paths.ContainersRoot()
	ents, _ := os.ReadDir(root)
	for _, e := range ents {
		if filepath.Ext(e.Name()) == ".json" {
			_ = os.Remove(filepath.Join(root, e.Name()))
		}
	}
	p := filepath.Join(root, "note.txt")
	if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	items, err := List()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty, got %d", len(items))
	}
}
