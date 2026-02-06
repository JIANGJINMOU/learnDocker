package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "ContainerEdu (cede) - a simplified Docker-like engine\n")
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  cede run --image <name> [--cmd <path>] [args...]\n")
	fmt.Fprintf(os.Stderr, "  cede build --dockerfile <path> --tag <name>\n")
	fmt.Fprintf(os.Stderr, "  cede ps\n")
	fmt.Fprintf(os.Stderr, "  cede pull --tar <path>\n")
	fmt.Fprintf(os.Stderr, "  cede net ls | release --id <containerID>\n")
	fmt.Fprintf(os.Stderr, "  cede net config --cidr <CIDR> --gateway <IP>\n")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	cmd := os.Args[1]
	switch cmd {
	case "run":
		runCmd := flag.NewFlagSet("run", flag.ExitOnError)
		image := runCmd.String("image", "", "image name to run")
		command := runCmd.String("cmd", "/bin/sh", "command to execute in container")
		hostname := runCmd.String("hostname", "cede", "UTS hostname inside container")
		netPlugin := runCmd.String("net", "", "optional network plugin name")
		cpuMax := runCmd.String("cpu", "100000 100000", "cgroup v2 cpu.max (quota period)")
		memMax := runCmd.String("mem", "256M", "cgroup v2 memory.max")
		pidsMax := runCmd.Int("pids", 64, "cgroup v2 pids.max")
		_ = hostname
		_ = netPlugin
		_ = cpuMax
		_ = memMax
		_ = pidsMax
		runCmd.Parse(os.Args[2:])
		args := runCmd.Args()
		if *image == "" {
			fmt.Fprintf(os.Stderr, "run: --image is required\n")
			os.Exit(2)
		}
		if err := runContainer(*image, *command, args, *hostname, *netPlugin, *cpuMax, *memMax, *pidsMax); err != nil {
			fmt.Fprintf(os.Stderr, "run error: %v\n", err)
			os.Exit(1)
		}
	case "build":
		buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
		dockerfile := buildCmd.String("dockerfile", "Dockerfile.cede", "path to simplified Dockerfile")
		tag := buildCmd.String("tag", "", "image tag")
		buildCmd.Parse(os.Args[2:])
		if *tag == "" {
			fmt.Fprintf(os.Stderr, "build: --tag is required\n")
			os.Exit(2)
		}
		if err := buildImage(*dockerfile, *tag); err != nil {
			fmt.Fprintf(os.Stderr, "build error: %v\n", err)
			os.Exit(1)
		}
	case "ps":
		if err := listContainers(); err != nil {
			fmt.Fprintf(os.Stderr, "ps error: %v\n", err)
			os.Exit(1)
		}
	case "pull":
		pullCmd := flag.NewFlagSet("pull", flag.ExitOnError)
		tar := pullCmd.String("tar", "", "path to docker save tarball")
		name := pullCmd.String("name", "", "image name to register")
		pullCmd.Parse(os.Args[2:])
		if *tar == "" || *name == "" {
			fmt.Fprintf(os.Stderr, "pull: --tar and --name are required\n")
			os.Exit(2)
		}
		if err := importImageTar(*tar, *name); err != nil {
			fmt.Fprintf(os.Stderr, "pull error: %v\n", err)
			os.Exit(1)
		}
	case "net":
		if len(os.Args) < 3 {
			usage()
			os.Exit(2)
		}
		sub := os.Args[2]
		switch sub {
		case "ls":
			if err := netList(); err != nil {
				fmt.Fprintf(os.Stderr, "net ls error: %v\n", err)
				os.Exit(1)
			}
		case "release":
			netCmd := flag.NewFlagSet("release", flag.ExitOnError)
			id := netCmd.String("id", "", "container id to release ip")
			netCmd.Parse(os.Args[3:])
			if *id == "" {
				fmt.Fprintf(os.Stderr, "net release: --id is required\n")
				os.Exit(2)
			}
			if err := netRelease(*id); err != nil {
				fmt.Fprintf(os.Stderr, "net release error: %v\n", err)
				os.Exit(1)
			}
		case "config":
			netCmd := flag.NewFlagSet("config", flag.ExitOnError)
			cidr := netCmd.String("cidr", "", "CIDR (e.g., 10.0.0.0/24)")
			gw := netCmd.String("gateway", "", "gateway IP (e.g., 10.0.0.1)")
			netCmd.Parse(os.Args[3:])
			if *cidr == "" && *gw == "" {
				fmt.Fprintf(os.Stderr, "net config: --cidr or --gateway required\n")
				os.Exit(2)
			}
			if err := netConfig(*cidr, *gw); err != nil {
				fmt.Fprintf(os.Stderr, "net config error: %v\n", err)
				os.Exit(1)
			}
		default:
			usage()
			os.Exit(2)
		}
	case "init":
		if err := childInit(); err != nil {
			fmt.Fprintf(os.Stderr, "init error: %v\n", err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(1)
	}
}
