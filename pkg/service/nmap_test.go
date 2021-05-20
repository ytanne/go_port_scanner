package service

import (
	"fmt"
	"testing"

	"github.com/ytanne/go_nessus/pkg/repository"
)

func TestPortScan(t *testing.T) {
	repo := repository.NewRepository(nil, nil, "", "", "")
	nmap := NewNmapScanner(repo)
	ports, err := nmap.ScanPorts("cert.kz")
	if err != nil {
		t.Fatalf("Could not scan ports of localhost: %s", err)
	}
	fmt.Println(ports)
}
