package ui

import forkcleaner "github.com/caarlos0/fork-cleaner/v2"

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

type getRepoListMsg struct{}
type getLocalRepoListMsg struct{}

type gotRepoListMsg struct {
	repos []*forkcleaner.RepositoryWithDetails
}

type gotLocalRepoListMsg struct {
	repos []*forkcleaner.LocalRepoState
}

type reposDeletedMsg struct{}
type localReposDeletedMsg struct{}

type requestDeleteSelectedReposMsg struct{}
type requestDeleteSelectedLocalReposMsg struct{}
