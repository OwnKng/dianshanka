package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"ownkng.dev/cli/types"
	"ownkng.dev/cli/vocab"
)

const listHeight = 10

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type Styles struct {
	BorderColor lipgloss.Color
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("#FF00FF")

	return s
}

type Main struct {
	index    int
	width    int
	height   int
	complete bool
	styles   *Styles
	list     list.Model
	choice   string
	game     types.Game
}

func (m Main) Init() tea.Cmd {
	return nil
}

func (m Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	//_ Handle resize events
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	//_ Handle key events
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":

			return m, tea.Quit

		case "enter":

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Main) View() string {
	current := m.game.Rounds[m.index]

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			current.Character,
			m.list.View(),
		),
	)
}

func main() {
	game := vocab.NewGame(10)
	round := game.Rounds[0]

	items := []list.Item{
		item{title: round.Cards[0].English, desc: round.Cards[0].Pinyin},
		item{title: round.Cards[1].English, desc: round.Cards[1].Pinyin},
		item{title: round.Cards[2].English, desc: round.Cards[2].Pinyin},
		item{title: round.Cards[3].English, desc: round.Cards[3].Pinyin},
	}

	l := list.New(items, list.NewDefaultDelegate(), 20, 20)
	l.Title = "Select the matching meaning"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	m := Main{
		index:    0,
		width:    20,
		height:   listHeight,
		complete: false,
		styles:   DefaultStyles(),
		list:     l,
		choice:   "",
		game:     game,
	}

	programme := tea.NewProgram(m)
	_, err := programme.Run()

	if err != nil {
		fmt.Println("Error creating program:", err)
		os.Exit(1)
	}
}
