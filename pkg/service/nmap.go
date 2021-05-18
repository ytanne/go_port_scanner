package service

import (
	"regexp"

	"github.com/ytanne/go_nessus/pkg/repository"
)

type NmapScanner struct {
	repo repository.NmapScan
}

func NewNmapScanner(repo repository.NmapScan) *NmapScanner {
	return &NmapScanner{
		repo: repo,
	}
}

func (ps *NmapScanner) ScanPorts(target string) (string, error) {
	ports, err := ps.repo.ScanPorts(target)
	if err != nil {
		return "", err
	}
	return string(ports), err
}

func (ns *NmapScanner) ScanNetwork(target string) ([]string, error) {
	ips, err := ns.repo.ScanNetwork(target)
	if err != nil {
		return nil, err
	}
	re, _ := regexp.Compile(`(?:[0-9]{1,3}\.){3}[0-9]{1,3}`)
	return re.FindAllString(string(ips), -1), err
}
