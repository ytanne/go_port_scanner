package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
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
	log.Println(cfg.Telegram.APItoken)
	return &cfg
}
