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
