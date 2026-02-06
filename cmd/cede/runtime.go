//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"example.com/containeredu/internal/cgroups"
	netplug "example.com/containeredu/internal/plugins/net"
	"example.com/containeredu/internal/overlay"
	"example.com/containeredu/internal/paths"
	"example.com/containeredu/internal/state"
	"example.com/containeredu/internal/id"
)

func runContainer(image, command string, args []string, hostname, netPlugin, cpuMax, memMax string, pidsMax int) error {
	if err := paths.EnsureDirs(); err != nil {
		return err
	}
	idStr := id.New()
	imgRoot := filepath.Join(paths.ImagesRoot(), image)
	layersRoot := filepath.Join(imgRoot, "layers")
	entries, err := os.ReadDir(layersRoot)
	if err != nil {
		return fmt.Errorf("image %s not found: %w", image, err)
	}
	var lowers []string
	for _, e := range entries {
		if e.IsDir() {
			lowers = append(lowers, filepath.Join(layersRoot, e.Name()))
		}
	}
	containerRoot := filepath.Join(paths.ContainersRoot(), idStr)
	upper := filepath.Join(containerRoot, "upper")
	work := filepath.Join(containerRoot, "work")
	mountDir := filepath.Join(containerRoot, "rootfs")
	if err := overlay.Prepare(overlay.MountSpec{
		LowerDirs: lowers,
		UpperDir:  upper,
		WorkDir:   work,
		MountDir:  mountDir,
	}); err != nil {
		return fmt.Errorf("overlay mount: %w", err)
	}
	initArgs := []string{"init", "--rootfs", mountDir, "--cmd", command, "--hostname", hostname}
	initArgs = append(initArgs, args...)
	cmd := exec.Command("/proc/self/exe", initArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNET | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	_ = cgroups.ApplyV2(idStr, cmd.Process.Pid, cgroups.Limits{
		CPUMax:  cpuMax,
		MemMax:  memMax,
		PidsMax: pidsMax,
	})
	var ip string
	if netPlugin != "" {
		if p := netplug.Get(netPlugin); p != nil {
			ip, _ = p.Setup(idStr, cmd.Process.Pid)
		}
	}
	st := state.ContainerState{
		ID:        idStr,
		Image:     image,
		Pid:       cmd.Process.Pid,
		Command:   command,
		Args:      args,
		CreatedAt: time.Now(),
		Hostname:  hostname,
		IP:        ip,
		Status:    "running",
		MountDir:  mountDir,
	}
	_ = state.Save(st)
	return cmd.Wait()
}

func childInit() error {
	var rootfs string
	var cmdPath string
	var hostname string
	rest := []string{}
	for i := 0; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--rootfs":
			i++
			if i < len(os.Args) {
				rootfs = os.Args[i]
			}
		case "--cmd":
			i++
			if i < len(os.Args) {
				cmdPath = os.Args[i]
			}
		case "--hostname":
			i++
			if i < len(os.Args) {
				hostname = os.Args[i]
			}
		default:
			rest = append(rest, os.Args[i])
		}
	}
	if rootfs == "" || cmdPath == "" {
		return fmt.Errorf("init: missing --rootfs or --cmd")
	}
	if err := syscall.Chroot(rootfs); err != nil {
		return fmt.Errorf("chroot: %w", err)
	}
	if err := os.Chdir("/"); err != nil {
		return err
	}
	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		return fmt.Errorf("mount proc: %w", err)
	}
	if hostname != "" {
		_ = syscall.Sethostname([]byte(hostname))
	}
	cmd := exec.Command(cmdPath, rest...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

 
