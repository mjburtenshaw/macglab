macglab
=======

Automate gathering your work on gitlab.com to save time.

![Static Badge](https://img.shields.io/badge/version-4.2.1-66023c)

Table of Contents
------------------

- [Installation](#installation)
    - [Requirements](#requirements)
    - [Updating](#updating)
- [Usage](#usage)
    - [Commands](#commands)
- [Configuration](#configuration)
    - [`access_token`](#access_token)
    - [`group_id`](#group_id)
    - [`me`](#me)
    - [`projects`](#projects)
    - [`usernames`](#usernames)

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

### Updating

To update to the latest version, pull the latest from the repository and reinstall the binary:

```sh
git checkout main
git pull
go install
```

Usage
-----

### Commands

- [`list`](#list)

#### `list`

Prints GitLab Merge Request (MRs) authors and URLs to the terminal.

```
macglab list [OPTIONS...]
```

`list` fetches MRs meeting ALL the following criteria:
- State is open.
- Belongs to [the configured group ID](#group_id).
- Is NOT a draft.
- Meets ANY of the following criteria:
    - The author is listed in [the configured usernames](#usernames).
    - The author is listed in ANY of [the configured projects](#projects); but it only returns MRs for projects the author is listed under.
    - [You](#me) are listed as a [reviewer](https://docs.gitlab.com/ee/user/project/merge_requests/reviews/#request-a-review).

`list` then excludes MRs meeting the following criteria:
- Approved by [you](#me).
- Mergeable MRs where [you](#me) are NOT the author.

##### Options

- `-a, --approved`: Include MRs [you](#me) approved.
- `-b, --browser`: Open MRs in the browser.
- `-d, --drafts`: Include draft MRs.
- `-g, --group`: ONLY include MRs where the author is listed in the provided users (*see `-u, --users`*) or [the configured usernames](#usernames).
- `-i <string>, --group-id=<string>`: Override [the configured group ID](#group_id) with the given string.
- `-m <number>, --me <number>`: Override [the configured `me`](#me) user ID with the given number.
- `-p, --projects`: ONLY include MRs where the author is listed in ANY of [the configured projects](#projects); but it only returns MRs for projects the author is listed under.
- `-r, --ready`: Include mergeable MRs.
- `-t <string>, --access-token <string>`: Override [the configured access token](#access_token).
- `-u <string>, --users=<string>`: Override [configured usernames](#usernames) and ONLY filter on usernames you provided. Accepts a CSV of usernames.

> 👯‍♀️ **Note:** `group` and `projects` are not mutually exclusive. If neither are provided, the program will run as if both are provided.

Configuration
--------------

See [the sample config](/config.sample.yml) for a full example.

### `access_token`

A [GitLab personal access tokens](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#create-a-personal-access-token).

### `group_id`

A [GitLab group ID](https://docs.gitlab.com/ee/api/groups.html).

### `me`

Your GitLab user ID (though it doesn't *have* to be yours). It's used for the following:
- Filter MRs based on approval.
- Include MRs where the given user ID is a reviewer.

### `projects`

A map of [GitLab project IDs](https://stackoverflow.com/questions/39559689/where-do-i-find-the-project-id-for-the-gitlab-api) having a list associated usernames you wish to follow. For example:

```yaml
projects:
    all: # usernames listed under the "all" entry will apply to every project.
        - username1
    123: # projectA
        - username2
        - username3
    456: # projectB
        - username3
        - username4
    789: # projectC
        # if left blank, this will inherit from `all`.
    101112:
        # projectD
        - username4
```

### `usernames`

A list of GitLab usernames in the group you wish to follow.
