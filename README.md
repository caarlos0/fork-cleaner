# fork-cleaner

[![Release](https://img.shields.io/github/release/caarlos0/fork-cleaner.svg?style=for-the-badge)](https://github.com/caarlos0/fork-cleaner/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](LICENSE.md)
[![Build Status](https://img.shields.io/github/actions/workflow/status/caarlos0/fork-cleaner/build.yml?style=for-the-badge)](https://github.com/caarlos0/fork-cleaner/actions?workflow=build)
[![Go Report Card](https://goreportcard.com/badge/github.com/caarlos0/fork-cleaner?style=for-the-badge)](https://goreportcard.com/report/github.com/caarlos0/fork-cleaner)
[![Godoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://pkg.go.dev/github.com/caarlos0/fork-cleaner)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge)](https://github.com/goreleaser)

Quickly clean up old and inactive forks on your GitHub account.

![](https://user-images.githubusercontent.com/245435/104655305-4a843f80-569c-11eb-8cd5-7f55b8104375.gif)

## Installation

### Homebrew

```sh
brew install caarlos0/tap/fork-cleaner
```

### snap

```sh
snap install fork-cleaner
```

### apt

```sh
echo 'deb [trusted=yes] https://repo.caarlos0.dev/apt/ /' | sudo tee /etc/apt/sources.list.d/caarlos0.list
sudo apt update
sudo apt install fork-cleaner
```

### yum

```sh
echo '[caarlos0]
name=caarlos0
baseurl=https://repo.caarlos0.dev/yum/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/caarlos0.repo
sudo yum install fork-cleaner
```

### deb/rpm/apk

Download the `.apk`, `.deb` or `.rpm` from the [latest release](https://github.com/caarlos0/fork-cleaner/releases/latest) and install with the appropriate commands.

### Manually

Download the binaries from the [latest release](https://github.com/caarlos0/fork-cleaner/releases/latest) or clone the repository and build from source.

## Usage

You'll need to [create a personal access token](https://github.com/settings/tokens/new?scopes=repo,delete_repo&description=fork-cleaner) with `repo` and `delete_repo`
permissions. You'll need to pass this token to `fork-cleaner` with the `--token` flag.

### Local mode

This is a newly added mode, which scans one or more git repositories that you have checked out (cloned) locally.
It marks each repository, as either "clean" (safe to delete), or "dirty" (not safe to delete).

For a repository to be marked clean, it needs to meet all of the following conditions:

* the are no uncommitted changes.
* all branches are found in a remote named "upstream" or "origin".
* there is nothing in the stash.

Note:

* if your local branches are found in remote that goes by another name, the local repository is still marked "dirty". This could perhaps
  be considered as "clean" in the future. (with an optional flag)

### Remote mode

This is the original fork-cleaner mode.

```sh
fork-cleaner --token "<token>"
```

`fork-cleaner` will load your forked repositories, displaying the oldest first. This can take a little while as `fork-cleaner` will iterate over the page of forks and check the upstream repository's status (e.g. checking for active PRs).

## Troubleshooting

### Taking forever to load?

The app hits various endpoints in order to collect information on the upstream repository, this can take a while if you have a lot of forks. Setting `-skip-upstream=true` will skip checking commits, issues, PRs, etc on each upstream repository, potentially alleviating this issue.

### I've hit the rate limit.

You can check your current limits by calling GitHub's API:

```sh
curl -L \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer <YOUR-TOKEN>" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/rate_limit
```

## Stargazers

[![Stargazers over time](https://starchart.cc/caarlos0/fork-cleaner.svg)](https://starchart.cc/caarlos0/fork-cleaner)
