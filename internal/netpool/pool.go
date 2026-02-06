package netpool

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"example.com/containeredu/internal/paths"
)

type fileData struct {
	CIDR        string            `json:"cidr"`
	Gateway     string            `json:"gateway"`
	Assignments map[string]string `json:"assignments"`
}

var (
	mu   sync.Mutex
	data *fileData
)

func filePath() string {
	return filepath.Join(paths.DataRoot(), "network", "netpool.json")
}

func ensureLoaded() error {
	if data != nil {
		return nil
	}
	fp := filePath()
	if err := os.MkdirAll(filepath.Dir(fp), 0o755); err != nil {
		return err
	}
	b, err := os.ReadFile(fp)
	if err != nil {
		data = &fileData{
			CIDR:        "10.0.0.0/24",
			Gateway:     "10.0.0.1",
			Assignments: map[string]string{},
		}
		return save()
	}
	var f fileData
	if err := json.Unmarshal(b, &f); err != nil {
		return err
	}
	data = &f
	return nil
}

func save() error {
	fp := filePath()
	b, _ := json.MarshalIndent(data, "", "  ")
	return os.WriteFile(fp, b, 0o644)
}

func Allocate(id string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	if err := ensureLoaded(); err != nil {
		return "", err
	}
	if ip, ok := data.Assignments[id]; ok && ip != "" {
		return ip, nil
	}
	base := "10.0.0."
	for i := 2; i <= 254; i++ {
		ip := base + strconv.Itoa(i)
		if ip == data.Gateway {
			continue
		}
		if used(ip) {
			continue
		}
		data.Assignments[id] = ip
		if err := save(); err != nil {
			return "", err
		}
		return ip, nil
	}
	return "", fmt.Errorf("no free ip")
}

func Release(id string) error {
	mu.Lock()
	defer mu.Unlock()
	if err := ensureLoaded(); err != nil {
		return err
	}
	delete(data.Assignments, id)
	return save()
}

func used(ip string) bool {
	for _, v := range data.Assignments {
		if strings.TrimSpace(v) == ip {
			return true
		}
	}
	return false
}

func CIDR() string {
	mu.Lock()
	defer mu.Unlock()
	_ = ensureLoaded()
	return data.CIDR
}

func Gateway() string {
	mu.Lock()
	defer mu.Unlock()
	_ = ensureLoaded()
	return data.Gateway
}

func SetCIDRGateway(cidr, gateway string) error {
	mu.Lock()
	defer mu.Unlock()
	if err := ensureLoaded(); err != nil {
		return err
	}
	if cidr != "" {
		data.CIDR = cidr
	}
	if gateway != "" {
		data.Gateway = gateway
	}
	return save()
}
