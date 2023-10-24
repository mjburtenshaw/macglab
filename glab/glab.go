package glab

import (
	"github.com/xanzy/go-gitlab"

	"github.com/mjburtenshaw/macglab/config"
)

var Client *gitlab.Client

func Initialize() error {
	client, err := gitlab.NewClient(config.AccessToken)
	if err != nil {
		return err
	}

	Client = client

	return nil
}
