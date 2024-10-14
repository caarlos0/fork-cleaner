package ui

import (
	"fmt"
	"strings"

	forkcleaner "github.com/caarlos0/fork-cleaner/v2"
	"github.com/charmbracelet/bubbles/list"
)

type localItem struct {
	repo     *forkcleaner.LocalRepoState
	selected bool
}

func (i localItem) Title() string {
	clean := " (DIRTY)"
	if i.repo.Clean() {
		clean = " (clean)"
	}
	if i.selected {
		return iconSelected + " " + i.repo.Path + clean
	}
	return iconNotSelected + " " + i.repo.Path + clean
}

func (i localItem) Description() string {
	var details []string
	if i.repo.StatusClean {
		details = append(details, "git status: clean")
	} else {
		details = append(details, "git status: dirty")
	}
	if i.repo.StashClean {
		details = append(details, "git stash:  clean")
	} else {
		details = append(details, "git stash:  dirty")
	}
	details = append(details, fmt.Sprintf("%d unmerged branches", len(i.repo.Unmerged)))

	return detailsStyle.Render(strings.Join(details, separator))
}

func (i localItem) FilterValue() string { return "  " + i.repo.Path }

func splitLocalBySelection(localItems []list.Item) ([]*forkcleaner.LocalRepoState, []*forkcleaner.LocalRepoState) {
	var selected, unselected []*forkcleaner.LocalRepoState
	for _, it := range localItems {
		localItem := it.(localItem)
		if localItem.selected {
			selected = append(selected, localItem.repo)
		} else {
			unselected = append(unselected, localItem.repo)
		}
	}
	return selected, unselected
}

func localReposToItems(repos []*forkcleaner.LocalRepoState) []list.Item {
	var localItems = make([]list.Item, 0, len(repos))
	for _, repo := range repos {
		localItems = append(localItems, localItem{
			repo: repo,
		})
	}
	return localItems
}
