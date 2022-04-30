package nmap

import (
	"context"
	"fmt"
	"testing"

	"github.com/ytanne/go_nessus/pkg/repository/nmap"
)

func TestPortScan(t *testing.T) {
	repo := nmap.NewScannerRepository()
	nmap := NewNmapService(repo)
	ports, err := nmap.ScanPorts(context.Background(), "cert.kz")
	if err != nil {
		t.Fatalf("Could not scan ports of localhost: %s", err)
	}
	fmt.Println(ports)
}
