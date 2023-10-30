package glab

import (
	"github.com/xanzy/go-gitlab"
)

var Client *gitlab.Client

func Initialize(accessToken string) error {
	client, err := gitlab.NewClient(accessToken)
	if err != nil {
		return err
	}

	Client = client

	return nil
}
