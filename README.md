# macglab

Trying to gather your work on gitlab.com takes too much time. Let's automate that instead.

This program prints all GitLab Merge Requests meeting the following criteria:
- Open
- Authored by given users
- Part of a given group

![Static Badge](https://img.shields.io/badge/version-2.1.1-66023c)

## Table of contents

- [Usage](#usage)
- [Installation](#installation)
- [Configuration](#configuration)

## Usage

Run `macglab` in a shell. Use the `-browser` flag to open MRs in the browser.

## Installation

1. Clone this repository, move into it, and run the install script:

```sh
git clone https://github.com/mjburtenshaw/macglab.git
cd macglab
go run install/install.go
```

2. Re-source your shell or open a new terminal to run the `macglab` command!

## Configuration

Create a `config.yml` that has the following values:
- `ACCESS_TOKEN`: [A GitLab personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#create-a-personal-access-token)
- `GROUP_ID`: [A GitLab group ID](https://docs.gitlab.com/ee/api/groups.html)
- `USERNAMES`: A list of GitLab usernames of users part of the group you wish to follow.
