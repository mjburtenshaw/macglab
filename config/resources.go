package config

import (
	"fmt"
	"log"
	"os"
)

var (
	HomeUri          string
	MacglabConfigUrl string
	MacglabUri       string
	ShConfigUrl      string
)

func init() {
	HomeUri = os.Getenv("HOME")
	if HomeUri == "" {
		log.Fatal("macglab: 🏚️ Couldn't find HOME environment variable.")
	}

	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		log.Fatal("macglab: 🏚️ Couldn't find GOPATH environment variable.")
	}

	ShConfigUrl = fmt.Sprintf("%s/.zshrc", HomeUri)
	MacglabUri = fmt.Sprintf("%s/.macglab", HomeUri)
	MacglabConfigUrl = fmt.Sprintf("%s/config.yml", MacglabUri)
}

func CheckFileExists(fileUrl string) error {
	info, err := os.Stat(fileUrl)
	if err != nil {
		return fmt.Errorf("%s doesn't exist: %w", fileUrl, err)
	} else if info.IsDir() {
		return fmt.Errorf("%s exists but is a directory", fileUrl)
	}
	return nil
}
