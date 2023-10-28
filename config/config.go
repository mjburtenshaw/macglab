package config

import (
	"fmt"
	"io"
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

func Get(configUrl string) (*Config, error) {
    if err := read(configUrl); err != nil {
        return nil, fmt.Errorf("couldn't get config at %s: %w", configUrl, err)
    }
    return config, nil
}

func read(configUrl string) error {
    if err := CheckFileExists(configUrl); err != nil {
        return fmt.Errorf("couldn't find %s: %w", configUrl, err)
    }

	configFile, err := os.ReadFile(configUrl)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configFile, &config)

	return err
}

func DemandConfigDir() error {
    info, err := os.Stat(MacglabUri)
    if err != nil {
        if os.IsNotExist(err) {
            err = os.MkdirAll(MacglabUri, 0755)
            return err
        }
        return err
    } else if !info.IsDir() {
        return fmt.Errorf("%s exists but is not a directory", MacglabUri)
    }
    return nil
}

func AddConfig(sampleConfigUrl string, configUrl string) (err error) {
    sampleConfig, err := os.Open(sampleConfigUrl)
    if err != nil {
        return fmt.Errorf("couldn't open sample config: %s", err)
    }
    defer func() {
        if cerr := sampleConfig.Close(); cerr != nil && err == nil {
            err = fmt.Errorf("couldn't close sample config: %s", cerr)
        }
    }()

    configFile, err := os.Create(configUrl)
    if err != nil {
        return fmt.Errorf("couldn't create config: %s", err)
    }
    defer func() {
        if cerr := configFile.Close(); cerr != nil && err == nil {
            err = fmt.Errorf("couldn't close config: %s", cerr)
        }
    }()

    if _, err = io.Copy(configFile, sampleConfig); err != nil {
        return fmt.Errorf("couldn't copy config: %s", err)
    }

    return nil
}
