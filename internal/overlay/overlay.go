//go:build linux

package overlay

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

type MountSpec struct {
	LowerDirs []string
	UpperDir  string
	WorkDir   string
	MountDir  string
}

func Prepare(spec MountSpec) error {
	for _, d := range []string{spec.UpperDir, spec.WorkDir, spec.MountDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}
	opts := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s",
		joinLower(spec.LowerDirs), spec.UpperDir, spec.WorkDir)
	if err := syscall.Mount("overlay", spec.MountDir, "overlay", 0, opts); err != nil {
		return err
	}
	return nil
}

func joinLower(dirs []string) string {
	return strings.Join(dirs, ":")
}

func Unmount(mountDir string) error {
	return syscall.Unmount(mountDir, 0)
}
