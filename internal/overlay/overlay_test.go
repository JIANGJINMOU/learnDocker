//go:build linux

package overlay

import (
	"os"
	"path/filepath"
	"testing"
)

func TestJoinLower(t *testing.T) {
	dirs := []string{"/a", "/b", "/c"}
	s := joinLower(dirs)
	if s != "/a:/b:/c" {
		t.Fatalf("unexpected join: %s", s)
	}
}

func TestPrepareCreatesDirs(t *testing.T) {
	tmp := t.TempDir()
	lower := filepath.Join(tmp, "lower1")
	upper := filepath.Join(tmp, "upper")
	work := filepath.Join(tmp, "work")
	mount := filepath.Join(tmp, "mount")
	
	if err := os.MkdirAll(lower, 0o755); err != nil {
		t.Fatal(err)
	}
	
	spec := MountSpec{
		LowerDirs: []string{lower},
		UpperDir:  upper,
		WorkDir:   work,
		MountDir:  mount,
	}
	
	// 测试Prepare函数是否创建了必要的目录
	if err := Prepare(spec); err != nil {
		t.Fatal(err)
	}
	
	// 验证目录是否创建
	if _, err := os.Stat(upper); os.IsNotExist(err) {
		t.Fatalf("upper dir not created: %v", err)
	}
	if _, err := os.Stat(work); os.IsNotExist(err) {
		t.Fatalf("work dir not created: %v", err)
	}
	if _, err := os.Stat(mount); os.IsNotExist(err) {
		t.Fatalf("mount dir not created: %v", err)
	}
	
	// 清理：卸载mount目录
	if err := Unmount(mount); err != nil {
		t.Logf("unmount error: %v", err)
	}
}
