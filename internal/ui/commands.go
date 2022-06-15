package ui

import (
	"context"
	"log"
	"strings"

	forkcleaner "github.com/caarlos0/fork-cleaner/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v45/github"
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

func getReposCmd(client *github.Client, login string) tea.Cmd {
	return func() tea.Msg {
		log.Println("getReposCmd")
		repos, err := forkcleaner.FindAllForks(context.Background(), client, login)
		if err != nil {
			return errMsg{err}
		}
		return gotRepoListMsg{repos}
	}
}
