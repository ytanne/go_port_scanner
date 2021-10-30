package nmap

import (
	"regexp"

	"github.com/ytanne/go_nessus/pkg/repository/nmap"
)

type NmapScanner interface {
	ScanPorts(target string) ([]string, error)
	ScanWebPorts(target string) ([]string, error)
	ScanNetwork(target string) ([]string, error)
}

type nmapScan struct {
	repo      nmap.NmapScanner
	reNetwork *regexp.Regexp
	rePorts   *regexp.Regexp
}

func NewNmapService(repo nmap.NmapScanner) NmapScanner {
	rePorts, _ := regexp.Compile(`(\d{1,5})\/(tcp|udp)[ \t]+open[ \t]+(\S+)[ \t]*(.*)?`)
	reNetwork, _ := regexp.Compile(`(?:[0-9]{1,3}\.){3}[0-9]{1,3}`)
	return &nmapScan{
		repo:      repo,
		reNetwork: reNetwork,
		rePorts:   rePorts,
	}
}

func (ps *nmapScan) ScanPorts(target string) ([]string, error) {
	ports, err := ps.repo.ScanPorts(target)
	if err != nil {
		return nil, err
	}
	result := ps.rePorts.FindAllString(string(ports), -1)
	return result, err
}

func (ps *nmapScan) ScanWebPorts(target string) ([]string, error) {
	ports, err := ps.repo.ScanWebPorts(target)
	if err != nil {
		return nil, err
	}
	result := ps.rePorts.FindAllString(string(ports), -1)
	return result, err
}

func (ns *nmapScan) ScanNetwork(target string) ([]string, error) {
	ips, err := ns.repo.ScanNetwork(target)
	if err != nil {
		return nil, err
	}
	return ns.reNetwork.FindAllString(string(ips), -1), err
}
