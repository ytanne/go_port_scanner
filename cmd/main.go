package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ytanne/go_nessus/pkg/app"
	"github.com/ytanne/go_nessus/pkg/config"
	nmapRepository "github.com/ytanne/go_nessus/pkg/repository/nmap"
	dbRepository "github.com/ytanne/go_nessus/pkg/repository/sqlite"
	tgRepository "github.com/ytanne/go_nessus/pkg/repository/telegram"
	nmapService "github.com/ytanne/go_nessus/pkg/service/nmap"
	dbService "github.com/ytanne/go_nessus/pkg/service/sqlite"
	tgService "github.com/ytanne/go_nessus/pkg/service/telegram"
)

func main() {
	cfg := config.InitConfig("./assets/config.yaml")

	dbRepo, err := dbRepository.NewDatabaseRepository(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	dbServ := dbService.NewDatabaseService(dbRepo)

	nmapRepo := nmapRepository.NewScannerRepository()
	nmapServ := nmapService.NewNmapService(nmapRepo)

	tgRepo, err := tgRepository.NewCommunicatorRepository(cfg.Telegram.APItoken, cfg.Telegram.ChatID)
	if err != nil {
		log.Fatalln(err)
	}
	tgServ := tgService.NewCommunicatorService(tgRepo)

	a := app.NewApp(tgServ, dbServ, nmapServ)

	if err := a.Run(); err != nil {
		log.Fatalf("Error occured. Exiting...")
	}

	log.Println("Exiting the module")
}
