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

func (ps *NmapScanner) ScanWebPorts(target string) ([]byte, error) {
	return getWebPorts(target)
}

func getPorts(address string) ([]byte, error) {
	return exec.Command("nmap", "--top-ports", "10000", "-T3", address).Output()
}

func getWebPorts(address string) ([]byte, error) {
	return exec.Command("nmap", "-p", "80,443,8080,8282,8181", "-T3", address).Output()
}

func (ns *NmapScanner) ScanNetwork(target string) ([]byte, error) {
	return getIPs(target)
}

func getIPs(target string) ([]byte, error) {
	return exec.Command("nmap", "-sn", "-T3", target).Output()
}
