# fork-cleaner

[![Release](https://img.shields.io/github/release/caarlos0/fork-cleaner.svg?style=flat-square)](https://github.com/caarlos0/fork-cleaner/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Travis](https://img.shields.io/travis/caarlos0/fork-cleaner.svg?style=flat-square)](https://travis-ci.org/caarlos0/fork-cleaner)
[![Go Report Card](https://goreportcard.com/badge/github.com/caarlos0/fork-cleaner?style=flat-square)](https://goreportcard.com/report/github.com/caarlos0/fork-cleaner)
[![Godoc](https://godoc.org/github.com/caarlos0/fork-cleaner?status.svg&style=flat-square)](http://godoc.org/github.com/caarlos0/fork-cleaner)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)

Cleans up old and inactive forks on your GitHub account.

You'll need to [create a personal access token](https://github.com/settings/tokens/new?scopes=repo,delete_repo&description=fork-cleaner) with `repo` and `delete_repo`
permissions.

Then, [download the latest release](https://github.com/caarlos0/fork-cleaner/releases)
and execute the binary as in:

```sh
./fork-cleaner --token "my github token"
```

Fork-Cleaner will show you all your forks, you can then check which you want
to delete or not on a TUI:

![Screen Recording](https://user-images.githubusercontent.com/245435/104655305-4a843f80-569c-11eb-8cd5-7f55b8104375.gif)

## Install

**homebrew**:

```sh
brew install fork-cleaner
```

**snap**:

```sh
snap install fork-cleaner
```

**apt**:

```sh
echo 'deb [trusted=yes] https://repo.caarlos0.dev/apt/ /' | sudo tee /etc/apt/sources.list.d/caarlos0.list
sudo apt update
sudo apt install fork-cleaner
```

**yum**:

```sh
echo '[caarlos0]
name=caarlos0
baseurl=https://repo.caarlos0.dev/yum/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/caarlos0.repo
sudo yum install fork-cleaner
```

**deb/rpm/apk**:

Download the `.apk`, `.deb` or `.rpm` from the [releases page][releases] and install with the appropriate commands.

**manually**:

Download the pre-compiled binaries from the [releases page][releases] or clone the repo build from source.

[releases]: https://github.com/caarlos0/fork-cleaner/releases


## Stargazers over time

[![Stargazers over time](https://starchart.cc/caarlos0/fork-cleaner.svg)](https://starchart.cc/caarlos0/fork-cleaner)
