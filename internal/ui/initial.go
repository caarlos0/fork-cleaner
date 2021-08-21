package ui

import (
	"context"
	"fmt"
	"log"

	forkcleaner "github.com/caarlos0/fork-cleaner/v2"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v33/github"
)

var keySelectAll = key.NewBinding(
	key.WithKeys("a"),
	key.WithHelp("a", "select all"),
)

var keySelectNone = key.NewBinding(
	key.WithKeys("n"),
	key.WithHelp("n", "select none"),
)

var keySelectToggle = key.NewBinding(
	key.WithKeys(" "),
	key.WithHelp("space", "toggle selected item"),
)

var keyDeletedSelected = key.NewBinding(
	key.WithKeys("d"),
	key.WithHelp("d", "delete selected forks"),
)

var keyConfirmDelete = key.NewBinding(
	key.WithKeys("y"),
	key.WithHelp("y", "confirm deleting"),
)

var keyCancelDelete = key.NewBinding(
	key.WithKeys("n", "esc"),
	key.WithHelp("n/esc", "go back"),
)

var regularKeys = []key.Binding{
	keySelectToggle,
	keySelectAll,
	keySelectNone,
	keyDeletedSelected,
}

var confirmingDeleteKeys = []key.Binding{
	keyConfirmDelete,
	keyCancelDelete,
}

// NewInitialModel creates a new InitialModel with required fields.
func NewInitialModel(client *github.Client, login string) InitialModel {
	return InitialModel{
		client: client,
		login:  login,
		list:   newList(false),
	}
}

func newList(confirmingDelete bool) list.Model {
	list := list.NewModel([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Fork Cleaner"
	list.SetSpinner(spinner.MiniDot)
	list.AdditionalShortHelpKeys = func() []key.Binding {
		if confirmingDelete {
			return confirmingDeleteKeys
		}
		return regularKeys
	}

	return list
}

// InitialModel is the UI when the CLI starts, basically loading the repos.
type InitialModel struct {
	err              error
	login            string
	client           *github.Client
	list             list.Model
	confirmingDelete bool
}

func (m InitialModel) Init() tea.Cmd {
	return enqueueGetReposCmd()
}

func (m InitialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		log.Println("tea.WindowSizeMsg")
		top, right, bottom, left := listStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	case errMsg:
		log.Println("errMsg")
		m.err = msg.error
	case getRepoListMsg:
		log.Println("getRepoListMsg")
		cmds = append(cmds, m.list.StartSpinner(), getReposCmd(m.client, m.login))
	case gotRepoListMsg:
		log.Println("gotRepoListMsg")
		var items = make([]list.Item, 0, len(msg.repos))
		for _, repo := range msg.repos {
			items = append(items, item{
				repo: repo,
			})
		}

		m.list.StopSpinner()
		cmds = append(cmds, m.list.SetItems(items))
	case askDeleteConfirmationMsg:
		log.Println("askDeleteConfirmationMsg")
		var items = make([]list.Item, 0, len(msg.repos))
		for _, repo := range msg.repos {
			items = append(items, item{
				repo:     repo,
				selected: true,
			})
		}
		m.list = newList(true)
		m.confirmingDelete = true
		cmds = append(cmds, m.list.SetItems(items))
	case tea.KeyMsg:
		log.Println("tea.KeyMsg")
		if !m.list.SettingFilter() {
			log.Println("tea.KeyMsg -> !settingFilter")

			if key.Matches(msg, keySelectAll) {
				log.Println("tea.KeyMsg -> !settingFilter -> selectAll")
				for idx, i := range m.list.Items() {
					item := i.(item)
					item.selected = true
					m.list.RemoveItem(idx)
					cmds = append(cmds, m.list.InsertItem(item, idx))
				}
			}

			if key.Matches(msg, keySelectNone) {
				log.Println("tea.KeyMsg -> !settingFilter -> selectNone")
				for idx, i := range m.list.Items() {
					item := i.(item)
					item.selected = false
					m.list.RemoveItem(idx)
					cmds = append(cmds, m.list.InsertItem(item, idx))
				}
			}

			if key.Matches(msg, keySelectToggle) {
				log.Println("tea.KeyMsg -> !settingFilter -> selectToggle")
				item := m.list.SelectedItem().(item)
				item.selected = !item.selected
				idx := m.list.Index()
				m.list.RemoveItem(idx)
				cmds = append(cmds, m.list.InsertItem(item, idx))
			}

			if key.Matches(msg, keyDeletedSelected) {
				log.Println("tea.KeyMsg -> !settingFilter -> deleteSelected")
				var selected []*forkcleaner.RepositoryWithDetails
				for _, i := range m.list.Items() {
					item := i.(item)
					if item.selected {
						selected = append(selected, item.repo)
					}
				}
				cmds = append(cmds, askDeleteConfirmationCmd(selected))
			}
		}
	}

	log.Println("default")
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m InitialModel) View() string {
	if m.err != nil {
		return errorStyle.Bold(true).Render("Error gathering the repository list") +
			"\n" +
			errorStyle.Render(m.err.Error())
	}
	return m.list.View()
}

// msgs

type getRepoListMsg struct{}

type gotRepoListMsg struct {
	repos []*forkcleaner.RepositoryWithDetails
}

type askDeleteConfirmationMsg struct {
	repos []*forkcleaner.RepositoryWithDetails
}

// cmds

func enqueueGetReposCmd() tea.Cmd {
	return func() tea.Msg {
		log.Println("enqueueGetReposCmd")
		return getRepoListMsg{}
	}
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

func askDeleteConfirmationCmd(repos []*forkcleaner.RepositoryWithDetails) tea.Cmd {
	return func() tea.Msg {
		return askDeleteConfirmationMsg{repos}
	}
}

// models

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
	return detailsStyle.Render(viewRepositoryDetails(i.repo))
}

func (i item) FilterValue() string { return i.repo.Name }
