package nmap

import (
	"context"
	"os/exec"
)

type NmapScanner interface {
	ScanPorts(ctx context.Context, target string) ([]byte, error)
	ScanWebPorts(ctx context.Context, target string) ([]byte, error)
	ScanNetwork(ctx context.Context, target string) ([]byte, error)
}

type nmapScan struct {
}

func NewScannerRepository() NmapScanner {
	return &nmapScan{}
}

func (ps *nmapScan) ScanPorts(ctx context.Context, target string) ([]byte, error) {
	return getPorts(ctx, target)
}

func (ps *nmapScan) ScanWebPorts(ctx context.Context, target string) ([]byte, error) {
	return getWebPorts(ctx, target)
}

func (ns *nmapScan) ScanNetwork(ctx context.Context, target string) ([]byte, error) {
	return getIPs(ctx, target)
}

func getIPs(ctx context.Context, target string) ([]byte, error) {
	return exec.CommandContext(ctx, "nmap", "-sn", "-T3", target).Output()
}

func getPorts(ctx context.Context, address string) ([]byte, error) {
	return exec.CommandContext(ctx, "nmap", "--top-ports", "1000", "-T3", address).Output()
}

func getWebPorts(ctx context.Context, address string) ([]byte, error) {
	return exec.CommandContext(ctx, "nmap", "-p", "80,443,8080,8282,8181", "-T3", address).Output()
}
