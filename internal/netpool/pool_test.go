package netpool

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAllocateRelease(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	p := filepath.Join(tmp, ".local", "share", "cede", "network")
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
	ip1, err := Allocate("a")
	if err != nil {
		t.Fatal(err)
	}
	ip2, err := Allocate("b")
	if err != nil {
		t.Fatal(err)
	}
	if ip1 == ip2 {
		t.Fatalf("duplicate allocation: %s", ip1)
	}
	if err := Release("a"); err != nil {
		t.Fatal(err)
	}
	ip1b, err := Allocate("a")
	if err != nil {
		t.Fatal(err)
	}
	if ip1b == "" {
		t.Fatalf("empty ip")
	}
}

func TestEnsureLoadedWithNoFile(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	// 确保network目录不存在，测试ensureLoaded的边界情况
	// 注意：这里我们需要重置全局变量data，以测试ensureLoaded的完整逻辑
	mu.Lock()
	data = nil
	mu.Unlock()
	
	// 调用CIDR函数，它会间接调用ensureLoaded
	cidr := CIDR()
	if cidr == "" {
		t.Fatalf("empty CIDR")
	}
	// 验证gateway是否正确设置
	gateway := Gateway()
	if gateway == "" {
		t.Fatalf("empty gateway")
	}
	// 只验证CIDR和gateway不为空，不验证具体值，因为可能会被其他测试影响
}
