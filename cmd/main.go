package main

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ytanne/go_nessus/pkg/app"
	"github.com/ytanne/go_nessus/pkg/config"
	"github.com/ytanne/go_nessus/pkg/repository"
	"github.com/ytanne/go_nessus/pkg/service"
	"github.com/ytanne/go_nessus/pkg/tg"
)

func main() {
	cfg := config.InitConfig("./assets/config.yaml")

	telegram, err := tg.NewTelegramConn(cfg.Telegram.APItoken, cfg.Telegram.ChatID)
	if err != nil {
		log.Fatalf("Could not initialize telegram bot. Error: %s", err)
	}

	initSQL, err := ioutil.ReadFile(cfg.DB.InitSQL)
	if err != nil {
		log.Fatalf("Could not read file from %s to initialize DB. Error: %s", cfg.DB.InitSQL, err)
	}

	db, err := sql.Open(cfg.DB.Type, cfg.DB.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(string(initSQL)); err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(db, telegram)
	serv := service.NewService(repo)
	a := app.NewApp(serv)

	if err := a.Run(); err != nil {
		log.Fatalf("Error occured. Exiting...")
	}

	log.Println("Exiting the module")
}
