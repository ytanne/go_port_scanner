package scanning

import (
	"context"
	"os/exec"
)

type nmapScan struct {
}

func NewScannerRepository() *nmapScan {
	return &nmapScan{}
}

func (ps *nmapScan) ScanPorts(ctx context.Context, target string) (string, error) {
	return getPorts(ctx, target)
}

func (ps *nmapScan) ScanWebPorts(ctx context.Context, target string) (string, error) {
	return getWebPorts(ctx, target)
}

func (ns *nmapScan) ScanNetwork(ctx context.Context, target string) (string, error) {
	return getIPs(ctx, target)
}

func getIPs(ctx context.Context, target string) (string, error) {
	output, err := exec.CommandContext(ctx, "nmap", "-sn", "-T3", target).Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func getPorts(ctx context.Context, address string) (string, error) {
	output, err := exec.CommandContext(ctx, "nmap", "--top-ports", "1000", "-T3", address).Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func getWebPorts(ctx context.Context, address string) (string, error) {
	output, err := exec.CommandContext(ctx, "nmap", "-p", "80,443,8080,8282,8181", "-T3", address).Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
