package paths

import (
	"os"
	"path/filepath"
)

func DataRoot() string {
	home, _ := os.UserHomeDir()
	if home == "" {
		home = "/var/lib"
	}
	return filepath.Join(home, ".local", "share", "cede")
}

func ImagesRoot() string {
	return filepath.Join(DataRoot(), "images")
}

func ContainersRoot() string {
	return filepath.Join(DataRoot(), "containers")
}

func EnsureDirs() error {
	dirs := []string{DataRoot(), ImagesRoot(), ContainersRoot()}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}
	return nil
}
