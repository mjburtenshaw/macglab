package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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
            log.Println("macglab: making home directory for macglab...")
            err = os.MkdirAll(MacglabUri, 0755)
            return err
        }
        return err
    } else if !info.IsDir() {
        return fmt.Errorf("%s exists but is not a directory", MacglabUri)
    }
    return nil
}

func AddEnv(shConfigUrl string) (err error) { 
    if didAddEnv, err := checkAddEnv(shConfigUrl); err != nil {
        return fmt.Errorf("couldn't check %s for environment variables: %w", shConfigUrl, err)
    } else if didAddEnv {
        return nil  // We already did the stuff below. Exit early.
    }

    shConfig, err := os.OpenFile(shConfigUrl, os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return fmt.Errorf("couldn't open %s: %w", shConfigUrl, err)
    }
    defer func() {
        if cerr := shConfig.Close(); cerr != nil && err == nil {
            err = cerr
        }
    }()

    envVariables := `
    # [macglab](https://github.com/mjburtenshaw/macglab)

    export MACGLAB="${HOME}/.macglab"
    export PATH="${GOPATH}/bin/macglab:${PATH}"
    `
    if _, err := shConfig.WriteString(envVariables); err != nil {
        return fmt.Errorf("couldn't write to %s: %w", shConfigUrl, err)
    }

    return nil
}

func checkAddEnv(shConfigUrl string) (didAddEnv bool, err error) {
    info, err := os.Stat(shConfigUrl)
    if err != nil {
        return false, fmt.Errorf("%s doesn't exist: %w", shConfigUrl, err)
    } else if info.IsDir() {
        return false, fmt.Errorf("%s exists but is a directory", shConfigUrl)
    }

    shConfig, err := os.Open(shConfigUrl)
    if err != nil {
        return false, fmt.Errorf("couldn't open %s: %w", shConfigUrl, err)
    }
    defer func() {
        if cerr := shConfig.Close(); cerr != nil && err == nil {
            err = cerr
        }
    }()

    contents, err := io.ReadAll(shConfig)
    if err != nil {
        return false, fmt.Errorf("couldn't read %s: %w", shConfigUrl, err)
    }

    if strings.Contains(string(contents), "macglab") {
        return true, nil
    }

    return false, nil
}
