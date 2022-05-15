package scanning

import (
	"context"
	"regexp"

	"github.com/ytanne/go_port_scanner/pkg/service"
)

type serviceScan struct {
	repo      service.PortScanner
	reNetwork *regexp.Regexp
	rePorts   *regexp.Regexp
}

func NewScanService(scanner service.PortScanner) *serviceScan {
	rePorts, _ := regexp.Compile(`(\d{1,5})\/(tcp|udp)[ \t]+open[ \t]+(\S+)[ \t]*(.*)?`)
	reNetwork, _ := regexp.Compile(`(?:[0-9]{1,3}\.){3}[0-9]{1,3}`)

	return &serviceScan{
		repo:      scanner,
		reNetwork: reNetwork,
		rePorts:   rePorts,
	}
}

func (ps *serviceScan) ScanPorts(ctx context.Context, target string) ([]string, error) {
	ports, err := ps.repo.ScanPorts(ctx, target)
	if err != nil {
		return nil, err
	}

	return ps.rePorts.FindAllString(string(ports), -1), nil
}

func (ps *serviceScan) ScanWebPorts(ctx context.Context, target string) ([]string, error) {
	ports, err := ps.repo.ScanWebPorts(ctx, target)
	if err != nil {
		return nil, err
	}

	return ps.rePorts.FindAllString(string(ports), -1), nil
}

func (ns *serviceScan) ScanNetwork(ctx context.Context, target string) ([]string, error) {
	ips, err := ns.repo.ScanNetwork(ctx, target)
	if err != nil {
		return nil, err
	}

	return ns.reNetwork.FindAllString(string(ips), -1), nil
}
