package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"ownkng.dev/cli/input"
	"ownkng.dev/cli/vocab"
)

const listHeight = 10

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Help  key.Binding
	Quit  key.Binding
	Enter key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},   // first column
		{k.Help, k.Quit}, // second column
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

type Main struct {
	width  int
	height int
	keys   keyMap
	help   help.Model
	list   input.Input
	choice string
	game   vocab.Game
	cursor int
	reveal bool
}

func (m Main) Init() tea.Cmd {
	return nil
}

func (m Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.game.Complete {
		return m, tea.Quit
	}

	switch msg := msg.(type) {

	//_ Handle resize events
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width

	//_ Handle key events
	case tea.KeyMsg:
		switch msg.String() {

		//* Help
		case "?":
			m.help.ShowAll = !m.help.ShowAll

		//* Quit the application
		case "q", "ctrl+c":

			return m, tea.Quit

		//* Handle key up and down - locked if in a reveal state
		case "up", "k":
			if m.cursor > 0 && !m.reveal {
				m.cursor--
				m.list.Down()
			}

		case "down", "j":
			if m.cursor < len(m.game.Rounds[m.game.Round].Cards)-1 && !m.reveal {
				m.cursor++
				m.list.Up()
			}

		//* Make a selection
		case "enter":

			//* If this is the second press
			if m.reveal {
				m.game.NextRound()

				//* Update the inputs
				items := []input.Item{}

				for _, card := range m.game.Rounds[m.game.Round].Cards {
					correct := false

					if card.Chinese == m.game.Rounds[m.game.Round].Card.Chinese {
						correct = true
					}

					items = append(items,
						input.Item{Title: card.English, Subtitle: card.Pinyin, Value: card.Chinese, Correct: correct},
					)
				}

				m.list = input.NewInput(items)
				m.cursor = 0

				//* Reset reveal
				m.reveal = false

				return m, nil
			}

			//* First press - reveal the correct answer
			m.game.MarkAnswer(m.list.GetSelected().Value)
			m.list.Reveal()
			m.reveal = true

			return m, nil
		}
	}

	return m, nil
}

// Big heading
var subduedColor = lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}

var headingStyles = lipgloss.NewStyle().Padding(0, 1).Margin(0, 1, 0, 0).Background(lipgloss.Color("#5de4c7")).Foreground(lipgloss.Color("#000000"))
var flashcardStyles = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#ffffff"))
var contentStyles = lipgloss.NewStyle().Margin(1, 0, 0, 1)

var secondaryText = lipgloss.NewStyle().
	Foreground(subduedColor)

func (m Main) View() string {
	current := m.game.Rounds[m.game.Round]

	helpView := m.help.View(m.keys)

	//* Headings
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		headingStyles.Render("你好!"),
		lipgloss.NewStyle().Margin(1, 0, 0, 0).Width(50).Render("Welcome to 金闪卡, a flashy little terminal app to practise Chinese."),
	)

	if m.game.Complete {
		score := m.game.GetScore() * 100
		content += lipgloss.JoinVertical(
			lipgloss.Left,
			secondaryText.Margin(3, 0, 3, 0).Render(fmt.Sprintf("Game complete. You scored %.2f%% \n", score)),
			helpView,
		)

		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Left,
			lipgloss.Top,
			contentStyles.Render(content),
		)
	}

	cards := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Padding(1, 0).Render(m.list.View()),
	)

	//* Cards
	reveal := ""

	if m.reveal {
		reveal = lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				flashcardStyles.MarginRight(2).Render(current.Card.Chinese),
				flashcardStyles.MarginRight(2).Render(current.Card.Pinyin),
				flashcardStyles.MarginRight(2).Render(current.Card.English),
				fmt.Sprintln(""),
			),
		)
	}

	content = content + lipgloss.JoinVertical(
		lipgloss.Left,
		flashcardStyles.MarginTop(3).Render(current.Card.Chinese),
		secondaryText.Margin(1, 0, 0, 0).Render("Select the correct English translation"),
		cards,
		reveal,
		helpView,
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Left,
		lipgloss.Top,
		contentStyles.Render(content),
	)
}

func main() {
	game := vocab.NewGame(5)
	round := game.Rounds[0]

	items := []input.Item{}

	for _, card := range round.Cards {
		correct := false

		if card.Chinese == round.Card.Chinese {
			correct = true
		}

		items = append(items,
			input.Item{Title: card.English, Subtitle: card.Pinyin, Value: card.Chinese, Correct: correct},
		)
	}

	l := input.NewInput(items)

	m := Main{
		width:  20,
		height: listHeight,
		list:   l,
		choice: "",
		game:   game,
		cursor: 0,
		help:   help.New(),
		keys:   keys,
	}

	programme := tea.NewProgram(m)
	_, err := programme.Run()

	if err != nil {
		fmt.Println("Error creating program:", err)
		os.Exit(1)
	}
}
