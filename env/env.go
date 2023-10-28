package env

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mjburtenshaw/macglab/files"
)

func Update(shConfigUrl string) (err error) { 
    if didUpdateEnv, err := check(shConfigUrl); err != nil {
        return fmt.Errorf("couldn't check %s for environment variables: %w", shConfigUrl, err)
    } else if didUpdateEnv {
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

    if strings.Contains(string(contents), "macglab") {
        return true, nil
    }

    return false, nil
}
