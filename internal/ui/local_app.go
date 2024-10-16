package ui

import (
	"log"
	"sort"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v50/github"
)

// AppModel is the UI when the CLI starts, basically loading the repos.
type LocalAppModel struct {
	err    error
	login  string
	client *github.Client
	path   string
	list   list.Model
}

// NewAppModel creates a new AppModel with required fields.
func NewLocalAppModel(client *github.Client, login, path string) LocalAppModel {
	list := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Fork Cleaner - local mode"
	list.SetSpinner(spinner.MiniDot)
	list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keySelectToggle,
			keyDeletedSelected,
			keySortBySize,
			keySortByName,
		}
	}
	list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keySelectAll,
			keySelectNone,
			keySortBySize,
			keySortByName,
		}
	}

	return LocalAppModel{
		client: client,
		login:  login,
		path:   path,
		list:   list,
	}
}

func (m LocalAppModel) Init() tea.Cmd {
	return tea.Batch(enqueueGetLocalReposCmd, m.list.StartSpinner())
}

func (m LocalAppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case getLocalRepoListMsg:
		log.Println("getLocalRepoListMsg")
		cmds = append(cmds, m.list.StartSpinner(), getLocalReposCmd(m.client, m.path))
	case gotLocalRepoListMsg:
		log.Println("gotLocalRepoListMsg")
		m.list.StopSpinner()
		cmds = append(cmds, m.list.SetItems(localReposToItems(msg.repos)))
	case localReposDeletedMsg:
		log.Println("localReposDeletedMsg")
	case requestDeleteSelectedLocalReposMsg:
		log.Println("requestDeleteSelectedLocalReposMsg")
		selected, unselected := splitLocalBySelection(m.list.Items())
		cmds = append(
			cmds,
			m.list.SetItems(localReposToItems(unselected)),
			deleteLocalReposCmd(selected),
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
			cmds = append(cmds, m.list.StartSpinner(), requestDeleteLocalReposCmd)
		}

		if key.Matches(msg, keySortBySize) {
			log.Println("tea.KeyMsg -> sortBySize")
			items := m.list.Items()
			sort.Sort(bySizeDesc(items))
			cmds = append(cmds, m.list.SetItems(items))
		}

		if key.Matches(msg, keySortByName) {
			log.Println("tea.KeyMsg -> sortByName")
			items := m.list.Items()
			sort.Sort(byName(items))
			cmds = append(cmds, m.list.SetItems(items))
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m LocalAppModel) View() string {
	if m.err != nil {
		return errorStyle.Bold(true).Render("Error gathering the repository list") +
			"\n" +
			errorStyle.Render(m.err.Error())
	}
	return m.list.View()
}

func (m LocalAppModel) toggleSelection() tea.Cmd {
	idx := m.list.Index()
	item := m.list.SelectedItem().(localItem)
	item.selected = !item.selected
	m.list.RemoveItem(idx)
	return m.list.InsertItem(idx, item)
}

func (m LocalAppModel) changeSelect(selected bool) []tea.Cmd {
	var cmds []tea.Cmd
	for idx, i := range m.list.Items() {
		item := i.(localItem)
		item.selected = selected
		m.list.RemoveItem(idx)
		cmds = append(cmds, m.list.InsertItem(idx, item))
	}
	return cmds
}
