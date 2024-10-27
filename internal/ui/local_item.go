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
		return iconSelected + " " + ByteCountIEC(i.repo.Size) + " " + i.repo.Path + clean
	}
	return iconNotSelected + " " + ByteCountIEC(i.repo.Size) + " " + i.repo.Path + clean
}

func (i localItem) Description() string {
	var details []string
	if i.repo.StatusClean {
		details = append(details, "status clean")
	} else {
		details = append(details, "status dirty")
	}
	if i.repo.StashClean {
		details = append(details, "stash clean")
	} else {
		details = append(details, "stash dirty")
	}
	if len(i.repo.Unmerged) > 2 || len(i.repo.Unmerged) == 0 {
		details = append(details, fmt.Sprintf("%d unmerged branches", len(i.repo.Unmerged)))
	} else {
		var keys []string
		for k := range i.repo.Unmerged {
			keys = append(keys, k)
		}
		details = append(details, fmt.Sprintf("unmerged: %s", strings.Join(keys, ", ")))
	}
	details = append(details, i.repo.RemotesChecked...)

	return detailsStyle.Render(strings.Join(details, separator))
}

func (i localItem) FilterValue() string {
	clean := "dirty"
	if i.repo.Clean() {
		clean = "clean"
	}

	return "  " + i.repo.Path + " " + clean + " " + strings.Join(i.repo.RemotesChecked, " ")
}

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

// support for sorting by size
type bySizeDesc []list.Item

func (s bySizeDesc) Less(i, j int) bool {
	return s[i].(localItem).repo.Size > s[j].(localItem).repo.Size
}
func (s bySizeDesc) Len() int      { return len(s) }
func (s bySizeDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type byName []list.Item

func (s byName) Less(i, j int) bool { return s[i].(localItem).repo.Path < s[j].(localItem).repo.Path }
func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// from https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
