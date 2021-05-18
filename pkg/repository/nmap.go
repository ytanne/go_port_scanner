package repository

import (
	"os/exec"
)

type NmapScanner struct {
}

func NewScanner() *NmapScanner {
	return &NmapScanner{}
}

func (ps *NmapScanner) ScanPorts(target string) ([]byte, error) {
	return getPorts(target)
}

func getPorts(address string) ([]byte, error) {
	return exec.Command("nmap", "--top-ports", "10000", "-sV", "-T3", address).Output()
}

func (ns *NmapScanner) ScanNetwork(target string) ([]byte, error) {
	return getIPs(target)
}

func getIPs(target string) ([]byte, error) {
	return exec.Command("nmap", "-sn", target).Output()
}
