package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DB struct {
		Type     string `yaml:"type"`
		Path     string `yaml:"path"`
		InitSQL  string `yaml:"init_sql"`
		AlterSQL string `yaml:"alter_sql"`
	} `yaml:"db"`
	Discord struct {
		Token        string `yaml:"token"`
		ARPChannelID string `yaml:"arp_channel_id"`
		PSChannelID  string `yaml:"ps_channel_id"`
		WPSChannelID string `yaml:"wps_channel_id"`
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
	return &cfg
}
