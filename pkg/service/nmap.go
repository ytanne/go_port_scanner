package service

import (
	"regexp"

	"github.com/ytanne/go_nessus/pkg/repository"
)

type NmapScanner struct {
	repo      repository.NmapScan
	reNetwork *regexp.Regexp
	rePorts   *regexp.Regexp
}

func NewNmapScanner(repo repository.NmapScan) *NmapScanner {
	rePorts, _ := regexp.Compile(`(\d{1,5})\/(tcp|udp)[ \t]+open[ \t]+(\S+)[ \t]*(.*)?`)
	reNetwork, _ := regexp.Compile(`(?:[0-9]{1,3}\.){3}[0-9]{1,3}`)
	return &NmapScanner{
		repo:      repo,
		reNetwork: reNetwork,
		rePorts:   rePorts,
	}
}

func (ps *NmapScanner) ScanPorts(target string) ([]string, error) {
	ports, err := ps.repo.ScanPorts(target)
	if err != nil {
		return nil, err
	}
	result := ps.rePorts.FindAllString(string(ports), -1)
	return result, err
}

func (ps *NmapScanner) ScanWebPorts(target string) ([]string, error) {
	ports, err := ps.repo.ScanWebPorts(target)
	if err != nil {
		return nil, err
	}
	result := ps.rePorts.FindAllString(string(ports), -1)
	return result, err
}

func (ns *NmapScanner) ScanNetwork(target string) ([]string, error) {
	ips, err := ns.repo.ScanNetwork(target)
	if err != nil {
		return nil, err
	}
	return ns.reNetwork.FindAllString(string(ips), -1), err
}
