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
	// 清理之前的状态文件
	root := paths.ContainersRoot()
	ents, _ := os.ReadDir(root)
	for _, e := range ents {
		_ = os.Remove(filepath.Join(root, e.Name()))
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

func TestListEmptyDirectory(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	if err := paths.EnsureDirs(); err != nil {
		t.Fatal(err)
	}
	// 确保containers目录为空
	root := paths.ContainersRoot()
	ents, _ := os.ReadDir(root)
	for _, e := range ents {
		_ = os.Remove(filepath.Join(root, e.Name()))
	}
	// 测试空目录的情况
	items, err := List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty list, got %d items", len(items))
	}
}

func TestSaveCreatesDirs(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	// 确保containers目录不存在
	root := paths.ContainersRoot()
	_ = os.RemoveAll(root)
	// 测试Save函数是否创建了必要的目录
	s := ContainerState{
		ID:        "test",
		Image:     "busybox",
		Pid:       123,
		Command:   "/bin/sh",
		Args:      []string{},
		CreatedAt: time.Now(),
		Status:    "running",
	}
	if err := Save(s); err != nil {
		t.Fatalf("save: %v", err)
	}
	// 验证状态文件是否创建
	p := filepath.Join(root, "test.json")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		t.Fatalf("state file not created: %v", err)
	}
}

func TestListNonExistentDirectory(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	// 确保containers目录不存在
	root := paths.ContainersRoot()
	_ = os.RemoveAll(root)
	// 测试List函数在目录不存在时的行为
	items, err := List()
	if err == nil {
		t.Fatalf("expected error when directory does not exist, got nil")
	}
	if items != nil {
		t.Fatalf("expected nil items when directory does not exist, got %+v", items)
	}
}
