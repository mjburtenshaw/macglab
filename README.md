macglab
=======

Automate gathering your work on gitlab.com to save time.

This program lists all GitLab Merge Requests (MRs) based on:
- Open state
- Specified usernames and/or projects
- Specified group

![Static Badge](https://img.shields.io/badge/version-3.0.0-66023c)

Table of Contents
------------------

- [Usage](#usage)
    - [List](#list)
- [Installation](#installation)
    - [Requirements](#requirements)
- [Configuration](#configuration)
- [See Also](#see-also)

Usage
-----

### List

Run `macglab list` in a shell:
- Use the `-b, --browser` flag to open MRs in the browser.
- Use the `-d, --drafts` flag to include draft MRs.
- Use the `-g, --group` flag to filter output to the usernames configuration.
- Use the `-p, --projects` flag to filter output to the projects configuration.
- Use the `-u, --users` flag to override configured usernames and only filter on usernames you provided. Accepts a CSV string of usernames. For example:

```sh
macglab list --users=harry,hermoine,ron
```

> 👯‍♀️ *`group` and `projects` are not mutually exclusive. If neither are provided, the program will run as if both are provided.*

Installation
-------------

### Requirements

1. Verify you have [installed Go](https://go.dev/doc/install): `go version`
2. Verify you have `GOPATH` in your shell environment: `echo "${GOPATH}"`
3. Verify you have added Go binaries to your `PATH`: `export PATH="${GOPATH}/bin:${PATH}"`

> 🐚 ***Tip:** it might be worth it to add the last command to your shell config file.*

--------------------------------------------------------------------------------------------

1. Clone this repository, move into it, install the binary, and run the install script:

```sh
git clone https://github.com/mjburtenshaw/macglab.git
cd macglab
go install
macglab init
```

2. Define values in the config file at `$HOME/.macglab/`. See [configuration](#configuration) for details.

3. Re-source your shell or open a new terminal to run the `macglab list` command!

Configuration
--------------

Update `config.yml` with:
- `ACCESS_TOKEN`: A GitLab personal access token[^1].
- `GROUP_ID`: A GitLab group ID[^2].
- `USERNAMES`: A list of GitLab usernames in the group you wish to follow.
- `PROJECTS`: A map of GitLab project IDs[^3] having a list associated usernames you wish to follow. For example:

```yaml
# usernames listed under the "all" entry will apply to every project listed below.

PROJECTS:
    all:
        - username1
    123:
        # projectA
        - username2
        - username3
    456:
        # projectB
        - username3
        - username4
    789:
        # projectC
        # if left blank, this will inherit from `all`.
    101112:
        # projectD
        - username4
```

See [the sample config](/config.sample.yml) for a full example.

See Also
---------

[^1]: [GitLab personal access tokens](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#create-a-personal-access-token)
[^2]: [GitLab groups](https://docs.gitlab.com/ee/api/groups.html)
[^3]: [GitLab project IDs](https://stackoverflow.com/questions/39559689/where-do-i-find-the-project-id-for-the-gitlab-api)
