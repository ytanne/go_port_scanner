package scanning

import (
	"context"
	"errors"
	"fmt"
	"github.com/ytanne/go_port_scanner/pkg/config"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type servNucleiScan struct {
	binPath string
}

func NewNucleiService(cfg config.Config) servNucleiScan {
	binPath := cfg.Nuclei.BinaryPath
	if binPath == "" {
		log.Println("nuclei binary path is empty")

		return servNucleiScan{}
	}

	_, err := os.Stat(binPath)
	if os.IsNotExist(err) {
		log.Println("could not find nuclei bin at start:", err)

		return servNucleiScan{}
	}

	// Installing templates to local machine
	if err := exec.Command(binPath).Run(); err != nil {
		log.Println("could not run nuclei at start:", err)

		return servNucleiScan{}
	}

	return servNucleiScan{
		binPath: cfg.Nuclei.BinaryPath,
	}
}

func (s servNucleiScan) ScanURL(ctx context.Context, targetURL string) (string, error) {
	if s.binPath == "" {
		return "", errors.New("nuclei binary path is empty")
	}

	_, err := os.Stat(s.binPath)
	if os.IsNotExist(err) {
		return "", err
	}

	const (
		HTTP  = "http://"
		HTTPs = "https://"
	)

	if !strings.HasPrefix(targetURL, HTTP) && !strings.HasPrefix(targetURL, HTTPs) {
		targetURL = fmt.Sprintf("http://%s", targetURL)
	}

	fileName := fmt.Sprintf("%d.txt", rand.Int())
	if err := exec.CommandContext(ctx, s.binPath, "-o", fileName, "-u", targetURL).Run(); err != nil {
		return "", fmt.Errorf("running nuclei failed: %w", err)
	}

	return fileName, nil
}
