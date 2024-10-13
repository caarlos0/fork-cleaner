package ui

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	forkcleaner "github.com/caarlos0/fork-cleaner/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v50/github"
)

func requestDeleteReposCmd() tea.Msg {
	return requestDeleteSelectedReposMsg{}
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
