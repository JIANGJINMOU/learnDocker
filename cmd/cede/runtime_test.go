package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSplitLines(t *testing.T) {
	in := "FROM scratch\nADD a b\n# comment\nADD c d"
	lines := parseDockerfile(in)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0].Keyword != "FROM" || lines[0].Arg != "scratch" {
		t.Fatalf("bad FROM parse: %+v", lines[0])
	}
	if lines[1].Keyword != "ADD" || lines[1].Arg != "a b" {
		t.Fatalf("bad ADD parse: %+v", lines[1])
	}
}

func TestBuildImageScratch(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	df := filepath.Join(tmp, "Dockerfile.cede")
	if err := os.WriteFile(df, []byte("FROM scratch\nADD "+tmp+" x"), 0o644); err != nil {
		t.Fatal(err)
	}
	tag := "testtag"
	if err := buildImage(df, tag); err != nil {
		t.Fatalf("build error: %v", err)
	}
	home, _ := os.UserHomeDir()
	meta := filepath.Join(home, ".local", "share", "cede", "images", tag, "metadata.json")
	if _, err := os.Stat(meta); err != nil {
		t.Fatalf("metadata missing: %v", err)
	}
}

func TestBuildImageAddFile(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	f := filepath.Join(tmp, "f.txt")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	df := filepath.Join(tmp, "Dockerfile.cede")
	if err := os.WriteFile(df, []byte("FROM scratch\nADD "+f+" etc/file.txt"), 0o644); err != nil {
		t.Fatal(err)
	}
	tag := "filetag"
	if err := buildImage(df, tag); err != nil {
		t.Fatalf("build error: %v", err)
	}
	home, _ := os.UserHomeDir()
	target := filepath.Join(home, ".local", "share", "cede", "images", tag, "layers", "00", "etc", "file.txt")
	b, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("file not copied: %v", err)
	}
	if string(b) != "x" {
		t.Fatalf("content mismatch: %q", string(b))
	}
}

func TestBuildImageErrors(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	df := filepath.Join(tmp, "Dockerfile.cede")
	if err := os.WriteFile(df, []byte("ADD a b"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := buildImage(df, "t1"); err == nil {
		t.Fatalf("expected FROM error")
	}
	if err := os.WriteFile(df, []byte("FROM busybox"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := buildImage(df, "t2"); err == nil {
		t.Fatalf("expected scratch error")
	}
	if err := os.WriteFile(df, []byte("FROM scratch\nADD a"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := buildImage(df, "t3"); err == nil {
		t.Fatalf("expected ADD args error")
	}
}

func TestSplitLinesTrim(t *testing.T) {
	s := "  FROM scratch  \n\nADD a b\n"
	lines := splitLines(s)
	if len(lines) != 3 {
		t.Fatalf("splitLines count: %d", len(lines))
	}
	if lines[0] != "FROM scratch" {
		t.Fatalf("trim failed: %q", lines[0])
	}
}

func TestSplitTwo(t *testing.T) {
	parts := splitTwo("a b")
	if len(parts) != 2 || parts[0] != "a" || parts[1] != "b" {
		t.Fatalf("bad splitTwo: %+v", parts)
	}
}

func TestSplitTwoMultiSpaces(t *testing.T) {
	parts := splitTwo("a   b   c")
	if len(parts) != 3 || parts[0] != "a" || parts[1] != "b" || parts[2] != "c" {
		t.Fatalf("bad splitTwo multi: %+v", parts)
	}
}

func TestTrimSpace(t *testing.T) {
	out := trimSpace(" \t abc \r ")
	if out != "abc" {
		t.Fatalf("trimSpace: %q", out)
	}
}

func TestSplitKeywordNoSpace(t *testing.T) {
	kw, rest := splitKeyword("FROM")
	if kw != "FROM" || rest != "" {
		t.Fatalf("splitKeyword no space: %s %s", kw, rest)
	}
}

func TestListContainersOutput(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	os.Setenv("USERPROFILE", tmp)
	// create state file
	if err := os.MkdirAll(filepath.Join(tmp, ".local", "share", "cede", "containers"), 0o755); err != nil {
		t.Fatal(err)
	}
	home := tmp
	stateFile := filepath.Join(home, ".local", "share", "cede", "containers", "x.json")
	if err := os.WriteFile(stateFile, []byte(`{"id":"x","image":"busybox","pid":1,"command":"/bin/sh","args":["-c","echo"],"status":"running"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	// capture stdout
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	if err := listContainers(); err != nil {
		t.Fatalf("list: %v", err)
	}
	w.Close()
	os.Stdout = old
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	out := string(buf[:n])
	if n == 0 || out == "" {
		t.Fatalf("no output")
	}
	if out == "" || out[0:2] != "ID" {
		t.Fatalf("header missing: %q", out)
	}
}

func TestTrimAndKeyword(t *testing.T) {
	line := "  ADD   foo bar  "
	kw, rest := splitKeyword(trimSpace(line))
	if kw != "ADD" {
		t.Fatalf("kw: %s", kw)
	}
	if rest != "foo bar" {
		t.Fatalf("rest: %s", rest)
	}
}

func TestCopyFileAndDir(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("HOME", tmp)
	srcDir := filepath.Join(tmp, "src")
	dstDir := filepath.Join(tmp, "dst")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	srcFile := filepath.Join(srcDir, "f.txt")
	if err := os.WriteFile(srcFile, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := copyPath(srcDir, dstDir); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dstDir, "f.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("copy mismatch: %q", string(data))
	}
}
