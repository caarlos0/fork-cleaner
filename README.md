# fork-cleaner

[![Release](https://img.shields.io/github/release/caarlos0/fork-cleaner.svg?style=flat-square)](https://github.com/caarlos0/fork-cleaner/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Travis](https://img.shields.io/travis/caarlos0/fork-cleaner.svg?style=flat-square)](https://travis-ci.org/caarlos0/fork-cleaner)
[![Go Report Card](https://goreportcard.com/badge/github.com/caarlos0/fork-cleaner?style=flat-square)](https://goreportcard.com/report/github.com/caarlos0/fork-cleaner)
[![Godoc](https://godoc.org/github.com/caarlos0/fork-cleaner?status.svg&style=flat-square)](http://godoc.org/github.com/caarlos0/fork-cleaner)
[![SayThanks.io](https://img.shields.io/badge/SayThanks.io-%E2%98%BC-1EAEDB.svg?style=flat-square)](https://saythanks.io/to/caarlos0)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)


Cleans up old and inactive forks on your GitHub account.

You'll need to create a personal access token with `repo` and `delete_repo`
permissions.

Then, [download the latest release](https://github.com/caarlos0/fork-cleaner/releases)
and execute the binary as in:

```console
./fork-cleaner --token "my github token"
```

Fork-Cleaner will show you repos that:

- are forks
- are not private (customizable via `--include-private,`)
- have no forks
- have no stars
- had no activity in the last 1 month (customizable via `--since`)
- are not blacklisted (customizable via `--blacklist`)

fork-cleaner will list them and ask if you want to remove them! Simple as that.

## Install

```console
brew install fork-cleaner
```
