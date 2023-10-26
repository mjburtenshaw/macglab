package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AccessToken string              `yaml:"ACCESS_TOKEN"`
	GroupId     string              `yaml:"GROUP_ID"`
	Me          int                 `yaml:"ME"`
	Projects    map[string][]string `yaml:"PROJECTS"`
	Usernames   []string            `yaml:"USERNAMES"`
}

var (
	config      *Config
	AccessToken string
	GroupId     string
	Me          int
	Projects    map[string][]string
	Usernames   []string
)

func Read() error {
	configFile, err := os.ReadFile(MacglabConfigUrl)
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

func DemandConfigDir() error {
    info, err := os.Stat(MacglabUri)
    if err != nil {
        if os.IsNotExist(err) {
            log.Println("macglab: 🏠 Making home directory for macglab...")
            err = os.MkdirAll(MacglabUri, 0755)
            return err
        }
        return err
    } else if !info.IsDir() {
        return fmt.Errorf("%s exists but is not a directory", MacglabUri)
    }
    return nil
}

