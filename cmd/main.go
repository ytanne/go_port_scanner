package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ytanne/go_port_scanner/pkg/app"
	"github.com/ytanne/go_port_scanner/pkg/composites"
	"github.com/ytanne/go_port_scanner/pkg/config"
)

func main() {
	cfg := config.InitConfig("./assets/config.yaml")

	// dbComp, err := composites.NewDBComposite(*cfg)
	// if err != nil {
	// 	log.Println("could not initialize new database composite:", err)

	// 	return
	// }

	dbComp, err := composites.NewMongoComposite(*cfg)
	if err != nil {
		log.Println("could not initialize new database composite:", err)

		return
	}

	comComp, err := composites.NewCommunicationComposite(*cfg)
	if err != nil {
		log.Println("could not initialize new communication composite:", err)

		return
	}

	scanComp := composites.NewScannerComposite(*cfg)

	a := app.NewApp(comComp.Serv, dbComp.Serv, scanComp.ServPort, scanComp.ServNuclei)
	a.SetUpChannels(cfg.Discord.ARPChannelID, cfg.Discord.PSChannelID, cfg.Discord.WPSChannelID)

	if err := a.Run(); err != nil {
		log.Fatalf("Error occured. Exiting...")
	}

	log.Println("Exiting the module")
}
