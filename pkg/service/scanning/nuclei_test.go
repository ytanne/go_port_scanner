package scanning

import (
	"context"
	"github.com/ytanne/go_port_scanner/pkg/config"
	"strings"
	"testing"
)

func TestServNucleiScan_ScanURL(t *testing.T) {
	newNuclei := NewNucleiService(config.Config{
		Nuclei: struct {
			BinaryPath string `yaml:"bin_path"`
		}(struct{ BinaryPath string }{BinaryPath: "/home/fl4ssh/go/src/github.com/ytanne/go_port_scanner/nuclei"}),
	})

	result, err := newNuclei.ScanURL(context.TODO(), "http://192.168.1.1")
	if err != nil {
		t.Fatal(err)
	}

	if result == "" || !strings.HasSuffix(result, ".txt") {
		t.Fatal("wrong result")
	}
}
