package nmap

import (
	"os/exec"
)

type NmapScanner interface {
	ScanPorts(target string) ([]byte, error)
	ScanWebPorts(target string) ([]byte, error)
	ScanNetwork(target string) ([]byte, error)
}

type nmapScan struct {
}

func NewScannerRepository() NmapScanner {
	return &nmapScan{}
}

func (ps *nmapScan) ScanPorts(target string) ([]byte, error) {
	return getPorts(target)
}

func (ps *nmapScan) ScanWebPorts(target string) ([]byte, error) {
	return getWebPorts(target)
}

func (ns *nmapScan) ScanNetwork(target string) ([]byte, error) {
	return getIPs(target)
}

func getIPs(target string) ([]byte, error) {
	return exec.Command("nmap", "-sn", "-T3", target).Output()
}

func getPorts(address string) ([]byte, error) {
	return exec.Command("nmap", "--top-ports", "10000", "-T3", address).Output()
}

func getWebPorts(address string) ([]byte, error) {
	return exec.Command("nmap", "-p", "80,443,8080,8282,8181", "-T3", address).Output()
}
