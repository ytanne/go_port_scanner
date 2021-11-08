package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Nessus struct {
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
		URL       string `yaml:"url"`
	} `yaml:"nessus"`
	Telegram struct {
		APItoken string `yaml:"api_token"`
		ChatID   int64  `yaml:"chat_id"`
	} `yaml:"telegram"`
	DB struct {
		Type     string `yaml:"type"`
		Path     string `yaml:"path"`
		InitSQL  string `yaml:"init_sql"`
		AlterSQL string `yaml:"alter_sql"`
	} `yaml:"db"`
	Discord struct {
		Token string `yaml:"token"`
	} `yaml:"discord"`
}

func InitConfig(filepath string) *Config {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Could not obtain file from %s. Error: %s", filepath, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Could not unmarshal %s. Error: %s", filepath, err)
	}
	if cfg.Discord.Token = os.Getenv("DISCORD_BOT"); cfg.Discord.Token == "" {
		log.Fatalf("Could not obtain discord token")
	}
	return &cfg
}
