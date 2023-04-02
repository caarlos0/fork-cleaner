package ui

import (
	"fmt"
	"strings"
	"time"

	forkcleaner "github.com/caarlos0/fork-cleaner/v2"
	timeago "github.com/caarlos0/timea.go"
	"github.com/charmbracelet/bubbles/list"
)

type item struct {
	repo     *forkcleaner.RepositoryWithDetails
	selected bool
}

func (i item) Title() string {
	var forked string
	if i.repo.ParentName != "" {
		forked = fmt.Sprintf(" (forked from %s)", i.repo.ParentName)
	}
	if i.selected {
		return iconSelected + " " + i.repo.Name + forked
	}
	return iconNotSelected + " " + i.repo.Name + forked
}

func (i item) Description() string {
	repo := i.repo
	var details []string
	if repo.ParentDeleted {
		details = append(details, "parent was deleted")
	}
	if repo.ParentDMCATakeDown {
		details = append(details, "parent was taken down by DMCA")
	}
	if repo.Private {
		details = append(details, "is private")
	}
	if repo.CommitsAhead > 0 {
		details = append(details, fmt.Sprintf("%d commit%s ahead", repo.CommitsAhead, maybePlural(repo.CommitsAhead)))
	}
	if repo.Forks > 0 {
		details = append(details, fmt.Sprintf("has %d fork%s", repo.Forks, maybePlural(repo.Forks)))
	}
	if repo.Stars > 0 {
		details = append(details, fmt.Sprintf("has %d star%s", repo.Stars, maybePlural(repo.Stars)))
	}
	if repo.OpenPRs > 0 {
		details = append(details, fmt.Sprintf("has %d open PR%s to upstream", repo.OpenPRs, maybePlural(repo.OpenPRs)))
	}
	if time.Now().Add(-30 * 24 * time.Hour).Before(repo.LastUpdate) {
		details = append(details, fmt.Sprintf("recently updated (%s)", timeago.Of(repo.LastUpdate)))
	}

	return detailsStyle.Render(strings.Join(details, separator))
}

func maybePlural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func (i item) FilterValue() string { return "  " + i.repo.Name }

func splitBySelection(items []list.Item) ([]*forkcleaner.RepositoryWithDetails, []*forkcleaner.RepositoryWithDetails) {
	var selected, unselected []*forkcleaner.RepositoryWithDetails
	for _, it := range items {
		item := it.(item)
		if item.selected {
			selected = append(selected, item.repo)
		} else {
			unselected = append(unselected, item.repo)
		}
	}
	return selected, unselected
}

func reposToItems(repos []*forkcleaner.RepositoryWithDetails) []list.Item {
	var items = make([]list.Item, 0, len(repos))
	for _, repo := range repos {
		items = append(items, item{
			repo: repo,
		})
	}
	return items
}
