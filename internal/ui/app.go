package ui

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v50/github"
)

// AppModel is the UI when the CLI starts, basically loading the repos.
type AppModel struct {
	err    error
	login  string
	client *github.Client
	list   list.Model
}

// NewAppModel creates a new AppModel with required fields.
func NewAppModel(client *github.Client, login string) AppModel {
	list := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Fork Cleaner"
	list.SetSpinner(spinner.MiniDot)
	list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keySelectToggle,
			keyDeletedSelected,
		}
	}
	list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keySelectAll,
			keySelectNone,
		}
	}

	return AppModel{
		client: client,
		login:  login,
		list:   list,
	}
}

func (m AppModel) Init() tea.Cmd {
	return tea.Batch(enqueueGetReposCmd, m.list.StartSpinner())
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m AppModel) View() string {
	if m.err != nil {
		return errorStyle.Bold(true).Render("Error gathering the repository list") +
			"\n" +
			errorStyle.Render(m.err.Error())
	}
	return m.list.View()
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
