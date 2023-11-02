package glab

import (
	"github.com/xanzy/go-gitlab"
)

type TGitlabClient = gitlab.Client

func Initialize(accessToken string) (*TGitlabClient, error) {
	return gitlab.NewClient(accessToken)
}
