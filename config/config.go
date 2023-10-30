package config

import (
	"fmt"
	"io"
	"os"

	"github.com/mjburtenshaw/macglab/files"
	"gopkg.in/yaml.v2"
)

type Config struct {
	AccessToken string              `yaml:"access_token"`
	GroupId     string              `yaml:"group_id"`
	Me          int                 `yaml:"me"`
	Projects    map[string][]string `yaml:"projects"`
	Usernames   []string            `yaml:"usernames"`
}

func Read(configUrl string) (*Config, error) {
	if err := files.CheckFileExists(configUrl); err != nil {
		return nil, fmt.Errorf("couldn't find %s: %w", configUrl, err)
	}

	configFile, err := os.ReadFile(configUrl)
	if err != nil {
		return nil, fmt.Errorf("couldn't read %s: %w", configUrl, err)
	}

	var config *Config
	if err = yaml.Unmarshal(configFile, &config); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal %s: %w", configUrl, err)
	}

	return config, nil
}

func Create(sampleConfigUrl string, configUrl string) (err error) {
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

func Update(configUrl string, key string, nextValue string) (err error) {
	configFile, err := os.OpenFile(configUrl, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("couldn't open config file: %w", err)
	}
	defer func() {
		if cerr := configFile.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("couldn't close config file: %s", cerr)
		}
	}()

	data, err := io.ReadAll(configFile)
	if err != nil {
		return fmt.Errorf("couldn't read config file: %w", err)
	}

	var content map[interface{}]interface{}
	if err = yaml.Unmarshal(data, &content); err != nil {
		return fmt.Errorf("couldn't unmarshal config file: %w", err)
	}

	if _, ok := content[key]; ok {
		content[key] = nextValue
	} else {
		return fmt.Errorf("invalid key: %s", key)
	}

	output, err := yaml.Marshal(content)
	if err != nil {
		return fmt.Errorf("couldn't marshal config file: %w", err)
	}

	if err = os.WriteFile(configUrl, output, 0644); err != nil {
		return fmt.Errorf("couldn't update config file: %w", err)
	}

	return nil
}
