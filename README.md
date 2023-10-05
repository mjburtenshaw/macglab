macglab
=======

Automate gathering your work on gitlab.com to save time.

This program lists all GitLab Merge Requests (MRs) based on:
- Open state
- Specified usernames and/or projects
- Specified group

![Static Badge](https://img.shields.io/badge/version-2.2.0-66023c)

Table of Contents
------------------

- [Usage](#usage)
- [Installation](#installation)
- [Configuration](#configuration)

Usage
-----

Run `macglab` in a shell:
- Use the `-browser` flag to open MRs in the browser.
- Use the `-group` flag to filter output to the usernames configuration.
- Use the `-projects` flag to filter output to the projects configuration.

> 👯‍♀️ *`group` and `projects` are not mutually exclusive. If neither are provided, the program will run as if both are provided.*

Installation
-------------

1. Clone this repository, move into it, and run the install script:

```sh
git clone https://github.com/mjburtenshaw/macglab.git
cd macglab
go run install/install.go
```

2. Re-source your shell or open a new terminal to run the `macglab` command!

Configuration
--------------

Update `config.yml` with:
- `ACCESS_TOKEN`: [GitLab personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#create-a-personal-access-token)
- `GROUP_ID`: [GitLab group ID](https://docs.gitlab.com/ee/api/groups.html)
- `USERNAMES`: List of GitLab usernames in the group you wish to follow.
- `PROJECTS`: Map of projects and associated usernames you wish to follow. For example:

```yaml
# usernames listed under the "all" entry will apply to every project listed below

PROJECTS:
    - all:
        - username1
    - project1:
        - username2
        - username3
    - project2:
        - username3
        - username4
```
