package scanning

import (
	"context"
	"fmt"
	"testing"
)

func TestPortScan(t *testing.T) {
	nmap := NewScanService()
	ports, err := nmap.ScanPorts(context.Background(), "cert.kz")
	if err != nil {
		t.Fatalf("Could not scan ports of localhost: %s", err)
	}
	fmt.Println(ports)
}
