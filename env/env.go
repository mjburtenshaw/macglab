package env

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mjburtenshaw/macglab/files"
)

const (
	macglabZshContent = `# [macglab](https://github.com/mjburtenshaw/macglab)

export MACGLAB="${HOME}/.macglab"
export PATH="${GOPATH}/bin/macglab:${PATH}"

`
	shContent  = "source ${HOME}/.macglab/macglab.zsh"
	createMode = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	appendMode = os.O_WRONLY | os.O_APPEND
)

// Checks if we've already installed environment variables.
// If not, it will create a file we manage to define environment variables.
// Then, we'll update the shell config file to source the first one.
// This way, if we need to push updates in the future, we can do so without
// reaching into the main shell config file, and not introduce a breaking change.
func Update(shConfigUrl string, macglabZshConfigUrl string) (err error) {
	if didUpdateEnv, err := check(shConfigUrl); err != nil {
		return fmt.Errorf("couldn't check %s for environment variables: %w", shConfigUrl, err)
	} else if didUpdateEnv {
		return nil // We already did the stuff below. Exit early.
	}

	if err := writeFile(macglabZshConfigUrl, createMode, macglabZshContent); err != nil {
		return err
	}

	return writeFile(shConfigUrl, appendMode, shContent)
}

func writeFile(fileUrl string, flag int, content string) error {
	file, err := os.OpenFile(fileUrl, flag, 0644)
	if err != nil {
		return fmt.Errorf("couldn't open or create %s: %w", fileUrl, err)
	}
	defer file.Close()

	if _, err = file.WriteString(content); err != nil {
		return fmt.Errorf("couldn't write to %s: %w", fileUrl, err)
	}

	return nil
}

func check(shConfigUrl string) (didUpdateEnv bool, err error) {
	if err := files.CheckFileExists(shConfigUrl); err != nil {
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

	return strings.Contains(string(contents), "macglab"), nil
}
