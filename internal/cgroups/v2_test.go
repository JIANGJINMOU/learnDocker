//go:build linux

package cgroups

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApplyV2WritesFiles(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("CEDE_CGROUP_ROOT", tmp)
	id := "cid1"
	err := ApplyV2(id, 1234, Limits{
		CPUMax:  "50000 100000",
		MemMax:  "64M",
		PidsMax: 8,
	})
	if err != nil {
		t.Fatal(err)
	}
	group := filepath.Join(tmp, "cede", id)
	check := []string{"cpu.max", "memory.max", "pids.max", "cgroup.procs"}
	for _, f := range check {
		if _, err := os.Stat(filepath.Join(group, f)); err != nil {
			t.Fatalf("missing %s", f)
		}
	}
}

func TestApplyV2WithEmptyLimits(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("CEDE_CGROUP_ROOT", tmp)
	id := "cid2"
	// 测试空的Limits值
	err := ApplyV2(id, 1234, Limits{})
	if err != nil {
		t.Fatal(err)
	}
	group := filepath.Join(tmp, "cede", id)
	// 验证cgroup.procs文件是否创建（即使没有设置其他限制）
	if _, err := os.Stat(filepath.Join(group, "cgroup.procs")); err != nil {
		t.Fatalf("missing cgroup.procs")
	}
}

func TestApplyV2WithPartialLimits(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("CEDE_CGROUP_ROOT", tmp)
	id := "cid3"
	// 测试部分Limits值
	err := ApplyV2(id, 1234, Limits{
		MemMax: "128M",
		// 只设置内存限制，不设置CPU和PID限制
	})
	if err != nil {
		t.Fatal(err)
	}
	group := filepath.Join(tmp, "cede", id)
	// 验证memory.max和cgroup.procs文件是否创建
	if _, err := os.Stat(filepath.Join(group, "memory.max")); err != nil {
		t.Fatalf("missing memory.max")
	}
	if _, err := os.Stat(filepath.Join(group, "cgroup.procs")); err != nil {
		t.Fatalf("missing cgroup.procs")
	}
}
