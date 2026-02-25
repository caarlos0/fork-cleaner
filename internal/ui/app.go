// Package ui contains the UI components for the app.
package ui

import (
	"log"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/google/go-github/v83/github"
)

// AppModel is the UI when the CLI starts, basically loading the repos.
type AppModel struct {
	err           error
	login         string
	client        *github.Client
	skipUpstream  bool
	list          list.Model
	lightdarkFunc lipgloss.LightDarkFunc
}

// NewAppModel creates a new AppModel with required fields.
func NewAppModel(
	client *github.Client,
	login string,
	skipUpstream bool,
) AppModel {
	list := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Fork Cleaner"
	list.SetSpinner(spinner.MiniDot)
	list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keySelectToggle,
			keyDeletedSelected,
			keyArchiveSelected,
		}
	}
	list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keySelectAll,
			keySelectNone,
		}
	}

	return AppModel{
		client:       client,
		login:        login,
		skipUpstream: skipUpstream,
		list:         list,
	}
}

func (m AppModel) Init() tea.Cmd {
	return tea.Batch(enqueueGetReposCmd, m.list.StartSpinner())
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		log.Println("tea.BackgroundColorMsg")
		m.lightdarkFunc = lipgloss.LightDark(msg.IsDark())
	case tea.WindowSizeMsg:
		log.Println("tea.WindowSizeMsg")
		top, right, bottom, left := listStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	case errMsg:
		log.Println("errMsg")
		m.err = msg.error
	case getRepoListMsg:
		log.Println("getRepoListMsg")
		cmds = append(cmds, m.list.StartSpinner(), getReposCmd(m.client, m.login, m.skipUpstream))
	case gotRepoListMsg:
		log.Println("gotRepoListMsg")
		m.list.StopSpinner()
		cmds = append(cmds, m.list.SetItems(reposToItems(msg.repos)))
	case reposDeletedMsg:
		log.Println("reposDeletedMsg")
		cmds = append(cmds, m.list.StartSpinner(), enqueueGetReposCmd)
	case requestDeleteSelectedReposMsg:
		log.Println("requestDeleteSelectedReposMsg")
		selected, unselected := splitBySelection(m.list.Items())
		cmds = append(
			cmds,
			m.list.SetItems(reposToItems(unselected)),
			deleteReposCmd(m.client, selected),
		)
	case requestArchiveSelectedReposMsg:
		log.Println("requestArchiveSelectedReposMsg")
		selected, unselected := splitBySelection(m.list.Items())
		cmds = append(
			cmds,
			m.list.SetItems(reposToItems(unselected)),
			archiveReposCmd(m.client, selected),
		)
	case tea.KeyMsg:
		if m.list.SettingFilter() {
			break
		}

		if key.Matches(msg, keySelectAll) {
			log.Println("tea.KeyMsg -> selectAll")
			cmds = append(cmds, m.changeSelect(true)...)
		}

		if key.Matches(msg, keySelectNone) {
			log.Println("tea.KeyMsg -> selectNone")
			cmds = append(cmds, m.changeSelect(false)...)
		}

		if key.Matches(msg, keySelectToggle) {
			log.Println("tea.KeyMsg -> selectToggle")
			cmds = append(cmds, m.toggleSelection())
		}

		if key.Matches(msg, keyDeletedSelected) {
			log.Println("tea.KeyMsg -> deleteSelected")
			cmds = append(cmds, m.list.StartSpinner(), requestDeleteReposCmd)
		}

		if key.Matches(msg, keyArchiveSelected) {
			log.Println("tea.KeyMsg -> archiveSelected")
			cmds = append(cmds, m.list.StartSpinner(), requestArchiveReposCmd)
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m AppModel) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true
	if m.err != nil {
		errStyle := errorStyle(m.lightdarkFunc)
		v.SetContent(
			errStyle.Bold(true).Render("Error gathering the repository list") +
				"\n" +
				errStyle.Render(m.err.Error()),
		)
	} else {
		v.SetContent(m.list.View())
	}
	return v
}

func (m AppModel) toggleSelection() tea.Cmd {
	idx := m.list.Index()
	item := m.list.SelectedItem().(item)
	item.selected = !item.selected
	m.list.RemoveItem(idx)
	return m.list.InsertItem(idx, item)
}

func (m AppModel) changeSelect(selected bool) []tea.Cmd {
	var cmds []tea.Cmd
	for idx, i := range m.list.Items() {
		item := i.(item)
		item.selected = selected
		m.list.RemoveItem(idx)
		cmds = append(cmds, m.list.InsertItem(idx, item))
	}
	return cmds
}
