package scanning

import (
	"context"
	"os/exec"
	"regexp"
)

type nmapScan struct {
	reNetwork *regexp.Regexp
	rePorts   *regexp.Regexp
}

func NewScanService() *nmapScan {
	rePorts, _ := regexp.Compile(`(\d{1,5})\/(tcp|udp)[ \t]+open[ \t]+(\S+)[ \t]*(.*)?`)
	reNetwork, _ := regexp.Compile(`(?:[0-9]{1,3}\.){3}[0-9]{1,3}`)

	return &nmapScan{
		reNetwork: reNetwork,
		rePorts:   rePorts,
	}
}

func (ps *nmapScan) ScanPorts(ctx context.Context, target string) ([]string, error) {
	ports, err := getPorts(ctx, target)
	if err != nil {
		return nil, err
	}

	return ps.rePorts.FindAllString(string(ports), -1), nil
}

func (ps *nmapScan) ScanWebPorts(ctx context.Context, target string) ([]string, error) {
	ports, err := getWebPorts(ctx, target)
	if err != nil {
		return nil, err
	}

	return ps.rePorts.FindAllString(string(ports), -1), nil
}

func (ns *nmapScan) ScanNetwork(ctx context.Context, target string) ([]string, error) {
	ips, err := getIPs(ctx, target)
	if err != nil {
		return nil, err
	}

	return ns.reNetwork.FindAllString(string(ips), -1), nil
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
