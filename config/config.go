package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AccessToken string              `yaml:"ACCESS_TOKEN"`
	GroupId     string              `yaml:"GROUP_ID"`
	Me					string							`yaml:"ME"`
	Projects    map[string][]string `yaml:"PROJECTS"`
	Usernames   []string            `yaml:"USERNAMES"`
}

var (
	config      *Config
	AccessToken string
	GroupId     string
	Me					string
	Projects    map[string][]string
	Usernames   []string
)

func Read() error {
	configHome := os.Getenv("MACGLAB")
	if configHome == "" {
		return fmt.Errorf("💀 Couldn't find MACGLAB environment variable")
	}

	configFile, err := os.ReadFile(fmt.Sprintf("%s/config.yml", configHome))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return err
	}

	AccessToken = config.AccessToken
	GroupId = config.GroupId
	Me = config.Me
	Projects = config.Projects
	Usernames = config.Usernames
	return nil
}
