package ui

import forkcleaner "github.com/caarlos0/fork-cleaner/v2"

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

type getRepoListMsg struct{}

type gotRepoListMsg struct {
	repos []*forkcleaner.RepositoryWithDetails
}

type reposDeletedMsg struct{}

type requestDeleteSelectedReposMsg struct{}

type requestArchiveSelectedReposMsg struct{}
