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

<img width="1486" alt="image" src="https://user-images.githubusercontent.com/245435/104636064-8958cc00-5681-11eb-999b-dcde7d3eb300.png">

## Install

On macOS:

```sh
brew install fork-cleaner
```

On Linux:

```sh
snap install fork-cleaner
```

Or download one of the archives from the releases page.
