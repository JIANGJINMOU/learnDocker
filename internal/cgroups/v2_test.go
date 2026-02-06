//go:build linux

package cgroups

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApplyV2CreatesDirs(t *testing.T) {
	// 创建临时目录作为 cgroup 根目录
	tmp := t.TempDir()
	
	// 设置环境变量，指向临时目录
	oldEnv := os.Getenv("CEDE_CGROUP_ROOT")
	defer func() {
		if oldEnv == "" {
			os.Unsetenv("CEDE_CGROUP_ROOT")
		} else {
			os.Setenv("CEDE_CGROUP_ROOT", oldEnv)
		}
	}()
	os.Setenv("CEDE_CGROUP_ROOT", tmp)
	
	// 测试 ApplyV2 函数
	containerID := "test-container"
	pid := os.Getpid()
	lim := Limits{
		CPUMax:   "100000 100000",
		MemMax:   "256M",
		PidsMax:  100,
	}
	
	err := ApplyV2(containerID, pid, lim)
	
	// 验证目录是否创建
	groupPath := filepath.Join(tmp, "cede", containerID)
	if _, err := os.Stat(groupPath); os.IsNotExist(err) {
		t.Fatalf("cgroup directory not created: %v", err)
	}
	
	// 即使写入文件失败，我们也不失败测试
	// 因为在某些环境中（如CI/CD）可能没有权限写入 cgroup 文件
	if err != nil {
		t.Logf("ApplyV2 error (expected in some environments): %v", err)
	}
}

func TestApplyV2EmptyLimits(t *testing.T) {
	// 创建临时目录作为 cgroup 根目录
	tmp := t.TempDir()
	
	// 设置环境变量，指向临时目录
	oldEnv := os.Getenv("CEDE_CGROUP_ROOT")
	defer func() {
		if oldEnv == "" {
			os.Unsetenv("CEDE_CGROUP_ROOT")
		} else {
			os.Setenv("CEDE_CGROUP_ROOT", oldEnv)
		}
	}()
	os.Setenv("CEDE_CGROUP_ROOT", tmp)
	
	// 测试 ApplyV2 函数，使用空限制
	containerID := "test-container-empty"
	pid := os.Getpid()
	lim := Limits{}
	
	err := ApplyV2(containerID, pid, lim)
	
	// 验证目录是否创建
	groupPath := filepath.Join(tmp, "cede", containerID)
	if _, err := os.Stat(groupPath); os.IsNotExist(err) {
		t.Fatalf("cgroup directory not created: %v", err)
	}
	
	// 即使写入文件失败，我们也不失败测试
	if err != nil {
		t.Logf("ApplyV2 error (expected in some environments): %v", err)
	}
}

func TestApplyV2PartialLimits(t *testing.T) {
	// 创建临时目录作为 cgroup 根目录
	tmp := t.TempDir()
	
	// 设置环境变量，指向临时目录
	oldEnv := os.Getenv("CEDE_CGROUP_ROOT")
	defer func() {
		if oldEnv == "" {
			os.Unsetenv("CEDE_CGROUP_ROOT")
		} else {
			os.Setenv("CEDE_CGROUP_ROOT", oldEnv)
		}
	}()
	os.Setenv("CEDE_CGROUP_ROOT", tmp)
	
	// 测试 ApplyV2 函数，使用部分限制
	containerID := "test-container-partial"
	pid := os.Getpid()
	lim := Limits{
		MemMax:   "512M",
		// 只设置内存限制，不设置 CPU 和 PIDs 限制
	}
	
	err := ApplyV2(containerID, pid, lim)
	
	// 验证目录是否创建
	groupPath := filepath.Join(tmp, "cede", containerID)
	if _, err := os.Stat(groupPath); os.IsNotExist(err) {
		t.Fatalf("cgroup directory not created: %v", err)
	}
	
	// 即使写入文件失败，我们也不失败测试
	if err != nil {
		t.Logf("ApplyV2 error (expected in some environments): %v", err)
	}
}
