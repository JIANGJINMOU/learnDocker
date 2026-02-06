//go:build linux

package bridge

import (
	"fmt"
	"strings"
	"testing"
)

func TestRunExecutesCommand(t *testing.T) {
	// 保存原始的runner函数
	oldRunner := runner
	defer func() {
		runner = oldRunner
	}()
	
	// 模拟runner函数，验证命令是否被正确调用
	var called bool
	var calledCmd string
	var calledArgs []string
	
	runner = func(cmd string, args ...string) error {
		called = true
		calledCmd = cmd
		calledArgs = args
		return nil
	}
	
	// 调用run函数
	err := run("echo", "hello", "world")
	
	// 验证run函数是否返回nil错误
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	
	// 验证runner是否被调用
	if !called {
		t.Fatalf("runner not called")
	}
	
	// 验证命令和参数是否正确
	if calledCmd != "echo" {
		t.Fatalf("expected command 'echo', got '%s'", calledCmd)
	}
	if len(calledArgs) != 2 {
		t.Fatalf("expected 2 arguments, got %d", len(calledArgs))
	}
	if calledArgs[0] != "hello" {
		t.Fatalf("expected first argument 'hello', got '%s'", calledArgs[0])
	}
	if calledArgs[1] != "world" {
		t.Fatalf("expected second argument 'world', got '%s'", calledArgs[1])
	}
}

func TestRunReturnsError(t *testing.T) {
	// 保存原始的runner函数
	oldRunner := runner
	defer func() {
		runner = oldRunner
	}()
	
	// 模拟runner函数，返回一个错误
	var called bool
	var calledCmd string
	var calledArgs []string
	
	runner = func(cmd string, args ...string) error {
		called = true
		calledCmd = cmd
		calledArgs = args
		// 模拟runImpl函数的错误格式
		return fmt.Errorf("%s %s: command failed", cmd, strings.Join(args, " "))
	}
	
	// 调用run函数
	err := run("echo", "hello")
	
	// 验证run函数是否返回错误
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	
	// 验证runner是否被调用
	if !called {
		t.Fatalf("runner not called")
	}
	
	// 验证错误信息是否包含命令和参数
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "echo") {
		t.Fatalf("error message does not contain command: %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "hello") {
		t.Fatalf("error message does not contain argument: %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "command failed") {
		t.Fatalf("error message does not contain original error: %s", errorMsg)
	}
}
