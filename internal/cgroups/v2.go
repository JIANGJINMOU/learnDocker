//go:build linux

package cgroups

import (
	"fmt"
	"os"
	"path/filepath"
)

type Limits struct {
	CPUMax   string // e.g., "100000 100000"
	MemMax   string // e.g., "256M"
	PidsMax  int
}

func ApplyV2(containerID string, pid int, lim Limits) error {
	root := "/sys/fs/cgroup"
	group := filepath.Join(root, "cede", containerID)
	if err := os.MkdirAll(group, 0o755); err != nil {
		return err
	}
	if lim.CPUMax != "" {
		if err := os.WriteFile(filepath.Join(group, "cpu.max"), []byte(lim.CPUMax), 0o644); err != nil {
			return fmt.Errorf("cpu.max: %w", err)
		}
	}
	if lim.MemMax != "" {
		if err := os.WriteFile(filepath.Join(group, "memory.max"), []byte(lim.MemMax), 0o644); err != nil {
			return fmt.Errorf("memory.max: %w", err)
		}
	}
	if lim.PidsMax > 0 {
		if err := os.WriteFile(filepath.Join(group, "pids.max"), []byte(fmt.Sprintf("%d", lim.PidsMax)), 0o644); err != nil {
			return fmt.Errorf("pids.max: %w", err)
		}
	}
	if err := os.WriteFile(filepath.Join(group, "cgroup.procs"), []byte(fmt.Sprintf("%d", pid)), 0o644); err != nil {
		return fmt.Errorf("cgroup.procs: %w", err)
	}
	return nil
}
