//go:build linux

package bridge

import (
	"fmt"
	"os/exec"
	"strings"

	nreg "example.com/containeredu/internal/plugins/net"
	"example.com/containeredu/internal/netpool"
)

type Plugin struct {
	name    string
	bridge  string
	gateway string
	cidr    string
}

func (p Plugin) Name() string { return p.name }

func (p Plugin) Setup(containerID string, pid int) (string, error) {
	br := p.bridge
	// ensure bridge exists and up
	if err := runner("ip", "link", "show", br); err != nil {
		if err := runner("ip", "link", "add", "name", br, "type", "bridge"); err != nil {
			return "", fmt.Errorf("bridge add: %w", err)
		}
		if err := runner("ip", "addr", "add", p.gateway+"/24", "dev", br); err != nil {
			// ignore if address exists
		}
		if err := runner("ip", "link", "set", br, "up"); err != nil {
			return "", fmt.Errorf("bridge up: %w", err)
		}
	}
	short := containerID
	if len(short) > 8 {
		short = short[:8]
	}
	hostV := "veth-" + short + "-h"
	ctV := "veth-" + short + "-c"
	// create veth pair
	if err := runner("ip", "link", "add", hostV, "type", "veth", "peer", "name", ctV); err != nil {
		return "", fmt.Errorf("veth add: %w", err)
	}
	// move container end into NET ns by pid
	if err := runner("ip", "link", "set", ctV, "netns", fmt.Sprintf("%d", pid)); err != nil {
		return "", fmt.Errorf("move to netns: %w", err)
	}
	// connect host end to bridge and up
	_ = runner("ip", "link", "set", hostV, "master", br)
	if err := runner("ip", "link", "set", hostV, "up"); err != nil {
		return "", fmt.Errorf("host veth up: %w", err)
	}
	ip, _ := netpool.Allocate(containerID)
	if err := runner("nsenter", "-t", fmt.Sprintf("%d", pid), "-n", "ip", "link", "set", "dev", ctV, "name", "eth0"); err != nil {
		return "", fmt.Errorf("rename eth0: %w", err)
	}
	if err := runner("nsenter", "-t", fmt.Sprintf("%d", pid), "-n", "ip", "addr", "add", ip+"/24", "dev", "eth0"); err != nil {
		return "", fmt.Errorf("addr add: %w", err)
	}
	if err := runner("nsenter", "-t", fmt.Sprintf("%d", pid), "-n", "ip", "link", "set", "eth0", "up"); err != nil {
		return "", fmt.Errorf("eth0 up: %w", err)
	}
	if err := runner("nsenter", "-t", fmt.Sprintf("%d", pid), "-n", "ip", "route", "add", "default", "via", p.gateway); err != nil {
		// ignore if route exists
	}
	if runner("iptables", "-t", "nat", "-C", "POSTROUTING", "-s", p.cidr, "!", "-o", br, "-j", "MASQUERADE") != nil {
		_ = runner("iptables", "-t", "nat", "-A", "POSTROUTING", "-s", p.cidr, "!", "-o", br, "-j", "MASQUERADE")
	}
	return ip, nil
}

func run(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	out, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s %s: %v (%s)", cmd, strings.Join(args, " "), err, string(out))
	}
	return nil
}

var runner = run

func init() {
	nreg.Register(Plugin{
		name:    "bridge0",
		bridge:  "cede0",
		gateway: "10.0.0.1",
		cidr:    "10.0.0.0/24",
	})
}
