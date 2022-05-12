package composites

import (
	"github.com/ytanne/go_nessus/pkg/app"
	repo "github.com/ytanne/go_nessus/pkg/repository/scanning"
	serv "github.com/ytanne/go_nessus/pkg/service/scanning"
)

type ScannerComposite struct {
	Serv app.PortScanner
}

func NewScannerComposite() ScannerComposite {
	scanRepo := repo.NewScannerRepository()
	scanServ := serv.NewScanService(scanRepo)

	return ScannerComposite{
		Serv: scanServ,
	}
}
