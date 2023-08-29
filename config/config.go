package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AccessToken string   `yaml:"ACCESS_TOKEN"`
	GroupId     string   `yaml:"GROUP_ID"`
	Usernames   []string `yaml:"USERNAMES"`
}

var (
	config     *Config
	AccessToken string
	GroupId     string
	Usernames   []string
)

func init() {
	configFile, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("💀 Failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("💀 Failed to unmarshal config file: %v", err)
	}

	AccessToken = config.AccessToken
	GroupId = config.GroupId
	Usernames = config.Usernames
}
