package scanning

import (
	"context"
	"fmt"
	"testing"

	"github.com/ytanne/go_nessus/pkg/repository/scanning"
)

func TestPortScan(t *testing.T) {
	repo := scanning.NewScannerRepository()
	nmap := NewScanService(repo)
	ports, err := nmap.ScanPorts(context.Background(), "cert.kz")
	if err != nil {
		t.Fatalf("Could not scan ports of localhost: %s", err)
	}
	fmt.Println(ports)
}
