package ui

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	forkcleaner "github.com/caarlos0/fork-cleaner/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v50/github"
)

func requestDeleteReposCmd() tea.Msg {
	return requestDeleteSelectedReposMsg{}
}

func requestDeleteLocalReposCmd() tea.Msg {
	return requestDeleteSelectedLocalReposMsg{}
}

func deleteReposCmd(client *github.Client, repos []*forkcleaner.RepositoryWithDetails) tea.Cmd {
	return func() tea.Msg {
		var names []string
		for _, r := range repos {
			names = append(names, r.Name)
		}
		log.Println("deleteReposCmd", strings.Join(names, ", "))
		if err := forkcleaner.Delete(context.Background(), client, repos); err != nil {
			return errMsg{err}
		}
		return reposDeletedMsg{}
	}
}

func deleteLocalReposCmd(repos []*forkcleaner.LocalRepoState) tea.Cmd {
	return func() tea.Msg {
		for _, r := range repos {
			log.Println("deleteLocalReposCmd: DELETING", r.Path)
			if err := os.RemoveAll(r.Path); err != nil {
				return errMsg{err}
			}
		}
		return localReposDeletedMsg{}
	}
}

func enqueueGetReposCmd() tea.Msg {
	return getRepoListMsg{}
}

func getReposCmd(client *github.Client, login string, skipUpstream bool) tea.Cmd {
	limits, _, err := client.RateLimits(context.Background())
	if err != nil {
		return func() tea.Msg {
			return errMsg{err}
		}
	}
	log.Println("RateLimits: ", limits)
	if limits.Core.Remaining < 1 {
		return func() tea.Msg {
			return errMsg{
				fmt.Errorf("Rate limit exceeded. Remaining: %d, Time till reset: %v",
					limits.Core.Remaining, time.Since(limits.Core.Reset.Time)),
			}
		}
	}

	return func() tea.Msg {
		log.Println("getReposCmd")
		repos, err := forkcleaner.FindAllForks(context.Background(), client, login, skipUpstream)
		if err != nil {
			return errMsg{err}
		}
		return gotRepoListMsg{repos}
	}
}

func enqueueGetLocalReposCmd() tea.Msg {
	return getLocalRepoListMsg{}
}

func getLocalReposCmd(client *github.Client, path string) tea.Cmd {
	return func() tea.Msg {

		// path should already have been validated to be a directory
		// if path has a .git directory in it, scan it
		// otherwise, find all directories inside of it that have a .git directory in them and scan them.

		_, err := os.Stat(filepath.Join(path, ".git"))

		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return errMsg{err}
		}
		if err == nil {
			lr, err := forkcleaner.NewLocalRepoState(path, client, context.Background())
			if err != nil {
				return errMsg{err}
			}
			return gotLocalRepoListMsg{[]*forkcleaner.LocalRepoState{lr}}

		}
		// we had an error but it was ErrNotExist for .git, so we assume it's a directory that contains code repos (checkouts)

		var repos []*forkcleaner.LocalRepoState

		entries, err := os.ReadDir(path)
		if err != nil {
			return errMsg{err}
		}

		for _, entry := range entries {
			if entry.IsDir() {
				gitpath := filepath.Join(path, entry.Name(), ".git")
				if _, err := os.Stat(gitpath); err == nil {
					lr, err := forkcleaner.NewLocalRepoState(filepath.Join(path, entry.Name()), client, context.Background())
					if err != nil {
						return errMsg{err}
					}
					repos = append(repos, lr)
				}
			}
		}

		return gotLocalRepoListMsg{repos}
	}
}
