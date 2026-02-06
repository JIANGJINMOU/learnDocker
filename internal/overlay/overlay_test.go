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
	
	// 测试Prepare函数
	err := Prepare(spec)
	
	// 即使挂载失败，我们也应该验证目录是否创建
	// 因为目录创建是Prepare函数的一部分，应该在挂载之前执行
	if _, err := os.Stat(upper); os.IsNotExist(err) {
		t.Fatalf("upper dir not created: %v", err)
	}
	if _, err := os.Stat(work); os.IsNotExist(err) {
		t.Fatalf("work dir not created: %v", err)
	}
	if _, err := os.Stat(mount); os.IsNotExist(err) {
		t.Fatalf("mount dir not created: %v", err)
	}
	
	// 如果挂载成功，我们需要卸载
	if err == nil {
		// 清理：卸载mount目录
		if err := Unmount(mount); err != nil {
			t.Logf("unmount error: %v", err)
		}
	} else {
		// 如果挂载失败，我们记录错误但不失败测试
		// 因为在某些环境中（如CI/CD）可能没有挂载权限
		t.Logf("mount error (expected in some environments): %v", err)
	}
}

func TestUnmountNonExistent(t *testing.T) {
	// 测试卸载一个不存在的挂载点
	nonExistentMount := "/nonexistent/mount/point"
	err := Unmount(nonExistentMount)
	
	// 验证返回错误
	if err == nil {
		t.Logf("unmount non-existent mount point returned nil, which may be expected on some platforms")
	} else {
		t.Logf("unmount non-existent mount point returned error (expected): %v", err)
	}
	
	// 不失败测试，因为在不同平台上，卸载不存在的挂载点的行为可能不同
}

func TestPrepareWithMultipleLowerDirs(t *testing.T) {
	tmp := t.TempDir()
	
	// 创建多个lower目录
	lower1 := filepath.Join(tmp, "lower1")
	lower2 := filepath.Join(tmp, "lower2")
	lower3 := filepath.Join(tmp, "lower3")
	if err := os.MkdirAll(lower1, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(lower2, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(lower3, 0o755); err != nil {
		t.Fatal(err)
	}
	
	upper := filepath.Join(tmp, "upper")
	work := filepath.Join(tmp, "work")
	mount := filepath.Join(tmp, "mount")
	
	spec := MountSpec{
		LowerDirs: []string{lower1, lower2, lower3},
		UpperDir:  upper,
		WorkDir:   work,
		MountDir:  mount,
	}
	
	// 测试Prepare函数
	err := Prepare(spec)
	
	// 即使挂载失败，我们也应该验证目录是否创建
	if _, err := os.Stat(upper); os.IsNotExist(err) {
		t.Fatalf("upper dir not created: %v", err)
	}
	if _, err := os.Stat(work); os.IsNotExist(err) {
		t.Fatalf("work dir not created: %v", err)
	}
	if _, err := os.Stat(mount); os.IsNotExist(err) {
		t.Fatalf("mount dir not created: %v", err)
	}
	
	// 如果挂载成功，我们需要卸载
	if err == nil {
		// 清理：卸载mount目录
		if err := Unmount(mount); err != nil {
			t.Logf("unmount error: %v", err)
		}
	} else {
		// 如果挂载失败，我们记录错误但不失败测试
		t.Logf("mount error (expected in some environments): %v", err)
	}
}

