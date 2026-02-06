package main

import (
	"os"
	"testing"
)

func TestUsage(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stderr
	os.Stderr = w
	usage()
	w.Close()
	os.Stderr = old
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	if n == 0 {
		t.Fatalf("no usage output")
	}
}

// 由于os.Exit不能被直接赋值，我们使用一个包装函数来测试main函数的行为
// 注意：这个测试会实际调用os.Exit，所以测试进程会退出
// 为了避免影响其他测试，我们只测试usage函数的行为

func TestMainUsageOutput(t *testing.T) {
	// 测试usage函数的输出
	r, w, _ := os.Pipe()
	old := os.Stderr
	os.Stderr = w
	usage()
	w.Close()
	os.Stderr = old
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	if n == 0 {
		t.Fatalf("no usage output")
	}
	// 验证输出包含预期的内容
	output := string(buf[:n])
	if len(output) == 0 {
		t.Fatalf("empty usage output")
	}
	if output[0:1] != "C" {
		t.Fatalf("usage output doesn't start with expected character")
	}
}
