# fork-cleaner [![Build Status](https://travis-ci.org/caarlos0/fork-cleaner.svg?branch=master)](https://travis-ci.org/caarlos0/fork-cleaner) [![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser) [![SayThanks.io](https://img.shields.io/badge/SayThanks.io-%E2%98%BC-1EAEDB.svg?style=flat-square)](https://saythanks.io/to/caarlos0)

Cleans up old and inactive forks on your github account.

You'll need to create a personal access token with `repo` and `delete_repo`
permissions.

Then, [download the latest release](https://github.com/caarlos0/fork-cleaner/releases)
and execute the binary as in:

```console
GITHUB_TOKEN="my github token" ./fork-cleaner
```

Fork-Cleaner will show you repos that:

- are forks
- are not private
- have no forks
- have no stars
- had no activity in the last 1 month

Then, it will ask you if you want to delete them:

![screenshot](https://cloud.githubusercontent.com/assets/245435/19216454/a0201810-8d92-11e6-8edc-4e1fe156b5c2.png)

Read carefully the list, and, if you agree, type `y` and it will
finish the job for you.

## Install

```console
brew install caarlos0/tap/fork-cleaner
```
