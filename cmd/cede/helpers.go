package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"example.com/containeredu/internal/images"
	"example.com/containeredu/internal/netpool"
	"example.com/containeredu/internal/paths"
	"example.com/containeredu/internal/state"
)

type dfLine struct {
	Keyword string
	Arg     string
}

func parseDockerfile(s string) []dfLine {
	var out []dfLine
	for _, line := range splitLines(s) {
		if line == "" || line[0] == '#' {
			continue
		}
		kw, rest := splitKeyword(line)
		out = append(out, dfLine{Keyword: kw, Arg: rest})
	}
	return out
}

func buildImage(dockerfile, tag string) error {
	content, err := os.ReadFile(dockerfile)
	if err != nil {
		return err
	}
	lines := parseDockerfile(string(content))
	if len(lines) == 0 || lines[0].Keyword != "FROM" {
		return fmt.Errorf("first line must be FROM")
	}
	if lines[0].Arg != "scratch" {
		return fmt.Errorf("only FROM scratch supported")
	}
	root := filepath.Join(paths.ImagesRoot(), tag, "layers", "00")
	if err := os.MkdirAll(root, 0o755); err != nil {
		return err
	}
	for _, l := range lines {
		if l.Keyword == "ADD" {
			parts := splitTwo(l.Arg)
			if len(parts) != 2 {
				return fmt.Errorf("ADD expects: ADD src dest")
			}
			src, dest := parts[0], parts[1]
			target := filepath.Join(root, dest)
			if err := copyPath(src, target); err != nil {
				return err
			}
		}
	}
	meta := struct {
		Name   string   `json:"name"`
		Layers []string `json:"layers"`
	}{
		Name:   tag,
		Layers: []string{"00"},
	}
	imgMeta := filepath.Join(paths.ImagesRoot(), tag, "metadata.json")
	if err := os.MkdirAll(filepath.Dir(imgMeta), 0o755); err != nil {
		return err
	}
	b, _ := json.MarshalIndent(meta, "", "  ")
	return os.WriteFile(imgMeta, b, 0o644)
}

func importImageTar(tarPath, name string) error {
	return images.ImportDockerSaveTar(tarPath, name)
}

func listContainers() error {
	items, err := state.List()
	if err != nil {
		return err
	}
	fmt.Printf("ID\tIMAGE\tPID\tSTATUS\tIP\tCMD\n")
	for _, it := range items {
		fmt.Printf("%s\t%s\t%d\t%s\t%s\t%s %v\n", it.ID, it.Image, it.Pid, it.Status, it.IP, it.Command, it.Args)
	}
	return nil
}

func netList() error {
	// print current assignments
	fmt.Printf("ID\tIP\n")
	ass := netpoolList()
	for id, ip := range ass {
		fmt.Printf("%s\t%s\n", id, ip)
	}
	return nil
}

func netRelease(id string) error {
	return netpool.Release(id)
}

func netConfig(cidr, gateway string) error {
	return netpool.SetCIDRGateway(cidr, gateway)
}

func netpoolList() map[string]string {
	// read file directly
	fp := filepath.Join(paths.DataRoot(), "network", "netpool.json")
	b, err := os.ReadFile(fp)
	if err != nil {
		return map[string]string{}
	}
	var v struct {
		Assignments map[string]string `json:"assignments"`
	}
	_ = json.Unmarshal(b, &v)
	if v.Assignments == nil {
		return map[string]string{}
	}
	return v.Assignments
}

func splitLines(s string) []string {
	lines := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, trimSpace(s[start:i]))
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, trimSpace(s[start:]))
	}
	return lines
}

func trimSpace(s string) string {
	b := 0
	e := len(s)
	for b < e && (s[b] == ' ' || s[b] == '\t' || s[b] == '\r') {
		b++
	}
	for e > b && (s[e-1] == ' ' || s[e-1] == '\t' || s[e-1] == '\r') {
		e--
	}
	return s[b:e]
}

func splitKeyword(line string) (string, string) {
	for i := 0; i < len(line); i++ {
		if line[i] == ' ' || line[i] == '\t' {
			return line[:i], trimSpace(line[i+1:])
		}
	}
	return line, ""
}

func splitTwo(s string) []string {
	out := []string{}
	cur := ""
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
			continue
		}
		cur += string(s[i])
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}

func copyPath(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst, info.Mode())
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	ents, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range ents {
		s := filepath.Join(src, e.Name())
		if skipInternalDataRoot(s) {
			continue
		}
		d := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := copyDir(s, d); err != nil {
				return err
			}
		} else {
			info, _ := os.Stat(s)
			if err := copyFile(s, d, info.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string, mode os.FileMode) error {
	if skipInternalDataRoot(src) {
		return nil
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := ioCopy(out, in); err != nil {
		out.Close()
		return err
	}
	out.Close()
	return os.Chmod(dst, mode)
}

func ioCopy(dst *os.File, src *os.File) (int64, error) {
	buf := make([]byte, 1<<20)
	var total int64
	for {
		n, err := src.Read(buf)
		if n > 0 {
			w, werr := dst.Write(buf[:n])
			total += int64(w)
			if werr != nil {
				return total, werr
			}
		}
		if err != nil {
			if err == io.EOF {
				return total, nil
			}
			return total, err
		}
	}
}

func skipInternalDataRoot(p string) bool {
	dr := paths.DataRoot()
	drAbs, _ := filepath.Abs(dr)
	pAbs, _ := filepath.Abs(p)
	// skip copying anything under our own data root to avoid recursive self-copy during build
	return strings.HasPrefix(pAbs, drAbs)
}
