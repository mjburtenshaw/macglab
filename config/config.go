package config

import (
	"fmt"
	"io"
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
    if err := CheckFileExists(shConfigUrl); err != nil {
        return false, fmt.Errorf("couldn't find %s: %w", shConfigUrl, err)
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
