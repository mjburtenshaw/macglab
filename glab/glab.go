package glab

import (
	"log"

	"github.com/xanzy/go-gitlab"

	"github.com/mjburtenshaw/macglab/config"
)

var Client *gitlab.Client

func init() {
	client, err := gitlab.NewClient(config.AccessToken)
	if err != nil {
		log.Fatalf("💀 Failed to create client: %v", err)
	}

	Client = client
}
