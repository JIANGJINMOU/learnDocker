package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: covermerge <profile1> <profile2> ...")
		os.Exit(2)
	}
	mode := ""
	type entry struct {
		stmts int
		count int
	}
	data := map[string]entry{}
	for _, path := range os.Args[1:] {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "open:", err)
			os.Exit(1)
		}
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := sc.Text()
			if strings.HasPrefix(line, "mode:") {
				if mode == "" {
					mode = line
				}
				continue
			}
			parts := strings.Fields(line)
			if len(parts) != 3 {
				continue
			}
			key := parts[0]
			var stmts, cnt int
			fmt.Sscanf(parts[1], "%d", &stmts)
			fmt.Sscanf(parts[2], "%d", &cnt)
			old := data[key]
			if old.stmts == 0 {
				old.stmts = stmts
			}
			old.count += cnt
			data[key] = old
		}
		f.Close()
	}
	if mode == "" {
		mode = "mode: atomic"
	}
	fmt.Println(mode)
	for k, v := range data {
		fmt.Printf("%s %d %d\n", k, v.stmts, v.count)
	}
}
