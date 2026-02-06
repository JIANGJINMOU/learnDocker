package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"example.com/containeredu/internal/paths"
)

type ContainerState struct {
	ID        string    `json:"id"`
	Image     string    `json:"image"`
	Pid       int       `json:"pid"`
	Command   string    `json:"command"`
	Args      []string  `json:"args"`
	CreatedAt time.Time `json:"created_at"`
	Hostname  string    `json:"hostname"`
	IP        string    `json:"ip"`
	Status    string    `json:"status"`
	MountDir  string    `json:"mount_dir"`
}

func Save(s ContainerState) error {
	if err := paths.EnsureDirs(); err != nil {
		return err
	}
	root := paths.ContainersRoot()
	p := filepath.Join(root, s.ID+".json")
	b, _ := json.MarshalIndent(s, "", "  ")
	return os.WriteFile(p, b, 0o644)
}

func List() ([]ContainerState, error) {
	root := paths.ContainersRoot()
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var out []ContainerState
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(root, e.Name()))
		if err != nil {
			continue
		}
		var s ContainerState
		if err := json.Unmarshal(b, &s); err != nil {
			continue
		}
		out = append(out, s)
	}
	return out, nil
}
