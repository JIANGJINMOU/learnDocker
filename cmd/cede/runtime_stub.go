//go:build !linux

package main

import "fmt"

func runContainer(image, command string, args []string, hostname, netPlugin, cpuMax, memMax string, pidsMax int) error {
	return fmt.Errorf("run is only supported on linux")
}

func childInit() error {
	return fmt.Errorf("init is only supported on linux")
}
