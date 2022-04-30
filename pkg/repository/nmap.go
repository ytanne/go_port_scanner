package repository

import (
	"context"
	"os/exec"
)

type NmapScanner struct {
}

func NewScanner() *NmapScanner {
	return &NmapScanner{}
}

func (ps *NmapScanner) ScanPorts(ctx context.Context, target string) ([]byte, error) {
	return getPorts(ctx, target)
}

func (ps *NmapScanner) ScanWebPorts(ctx context.Context, target string) ([]byte, error) {
	return getWebPorts(ctx, target)
}

func (ns *NmapScanner) ScanNetwork(ctx context.Context, target string) ([]byte, error) {
	return getIPs(ctx, target)
}

func getPorts(ctx context.Context, address string) ([]byte, error) {
	return exec.CommandContext(ctx, "nmap", "--top-ports", "10000", "-T3", address).Output()
}

func getWebPorts(ctx context.Context, address string) ([]byte, error) {
	return exec.CommandContext(ctx, "nmap", "-p", "80,443,8080,8282,8181", "-T3", address).Output()
}

func getIPs(ctx context.Context, target string) ([]byte, error) {
	return exec.CommandContext(ctx, "nmap", "-sn", "-T3", target).Output()
}
