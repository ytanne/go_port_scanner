package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ytanne/go_nessus/pkg/app"
	"github.com/ytanne/go_nessus/pkg/config"
	discordRepo "github.com/ytanne/go_nessus/pkg/repository/discord"
	nmapRepository "github.com/ytanne/go_nessus/pkg/repository/nmap"
	dbRepository "github.com/ytanne/go_nessus/pkg/repository/sqlite"
	discordServ "github.com/ytanne/go_nessus/pkg/service/discord"
	nmapService "github.com/ytanne/go_nessus/pkg/service/nmap"
	dbService "github.com/ytanne/go_nessus/pkg/service/sqlite"
)

func main() {
	cfg := config.InitConfig("./assets/config.yaml")

	dbRepo, err := dbRepository.NewDatabaseRepository(cfg)
	if err != nil {
		log.Fatalln("Could not create new database repository:", err)
	}

	dbServ := dbService.NewDatabaseService(dbRepo)

	nmapRepo := nmapRepository.NewScannerRepository()
	nmapServ := nmapService.NewNmapService(nmapRepo)

	// tgRepo, err := tgRepository.NewCommunicatorRepository(cfg.Telegram.APItoken, cfg.Telegram.ChatID)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// tgServ := tgService.NewCommunicatorService(tgRepo)
	log.Println("Obtained discord token:", cfg.Discord.Token)
	discordRepo, err := discordRepo.NewDiscordBot(cfg.Discord.Token)
	if err != nil {
		log.Fatalln("Could not create new discord bot:", err)
	}
	discordServ := discordServ.NewDiscordService(discordRepo)

	a := app.NewApp(discordServ, dbServ, nmapServ)
	a.SetUpChannels(cfg.Discord.ARPChannelID, cfg.Discord.PSChannelID, cfg.Discord.WPSChannelID)

	if err := a.Run(); err != nil {
		log.Fatalf("Error occured. Exiting...")
	}

	log.Println("Exiting the module")
}
