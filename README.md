# macglab

Trying to gather your work on gitlab.com takes too much time. Let's automate that instead.

This program opens all GitLab Merge Requests in a browser meeting the following criteria:
- Open
- Authored by given users
- Part of a given group

![Static Badge](https://img.shields.io/badge/version-1.0.0-66023c)

## Table of contents

- [Config](#config)

## Config

Create a `config.yml` that has the following values:
- `ACCESS_TOKEN`: [A GitLab personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#create-a-personal-access-token)
- `GROUP_ID`: [A GitLab group ID](https://docs.gitlab.com/ee/api/groups.html)
- `USERNAMES`: A list of GitLab usernames of users part of the group you wish to follow.

## Usage

Run `go main.go` in a shell.
