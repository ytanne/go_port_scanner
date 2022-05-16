package composites

import (
	"github.com/ytanne/go_port_scanner/pkg/app"
	"github.com/ytanne/go_port_scanner/pkg/config"
	serv "github.com/ytanne/go_port_scanner/pkg/service/scanning"
)

type ScannerComposite struct {
	ServPort   app.PortScanner
	ServNuclei app.NucleiScanner
}

func NewScannerComposite(cfg config.Config) ScannerComposite {
	scanServ := serv.NewScanService()
	nucleiServ := serv.NewNucleiService(cfg)

	return ScannerComposite{
		ServPort:   scanServ,
		ServNuclei: nucleiServ,
	}
}
