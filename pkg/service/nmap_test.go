package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/ytanne/go_nessus/pkg/repository"
)

func TestPortScan(t *testing.T) {
	repo := repository.NewRepository(nil, nil)
	nmap := NewNmapScanner(repo)
	ports, err := nmap.ScanPorts(context.Background(), "cert.kz")
	if err != nil {
		t.Fatalf("Could not scan ports of localhost: %s", err)
	}
	fmt.Println(ports)
}
