package config

import (
	"fmt"
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
	configHome := os.Getenv("MACGLAB")
	if configHome == "" {
		log.Fatal("💀 Couldn't find MACGLAB environment variable")
	}

	configFile, err := os.ReadFile(fmt.Sprintf("%s/config.yml", configHome))
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
