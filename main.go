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

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Help  key.Binding
	Quit  key.Binding
	Enter key.Binding
	Plus  key.Binding
	Minus key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Plus, k.Minus, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Plus, k.Down, k.Enter}, // first column
		{k.Help, k.Quit},                        // second column
	}
}

var titleKeys = keyMap{
	Plus: key.NewBinding(
		key.WithKeys("+"),
		key.WithHelp("+", "add round"),
	),
	Minus: key.NewBinding(
		key.WithKeys("-"),
		key.WithHelp("-", "subtract round"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "start game"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "quit"),
		key.WithHelp("q", "quit"),
	),
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
	width        int
	height       int
	keys         keyMap
	help         help.Model
	list         input.Input
	game         vocab.Game
	cursor       int
	reveal       bool
	numerOfRound int
}

func newModel() Main {
	return Main{
		width:        20,
		height:       20,
		cursor:       0,
		keys:         titleKeys,
		help:         help.New(),
		numerOfRound: 10,
	}
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

		//* Handle plus and minus keys
		case "+":
			if len(m.game.Rounds) == 0 || m.game.Complete {
				m.numerOfRound++
			}

		case "-":
			if (len(m.game.Rounds) == 0 && m.numerOfRound > 1) || m.game.Complete && m.numerOfRound > 1 {
				m.numerOfRound--
			}

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

			//* First press

			//* If the game has not started, start it
			if len(m.game.Rounds) == 0 {

				m.game = vocab.NewGame(m.numerOfRound)
				round := m.game.Rounds[0]

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

				m.list = input.NewInput(items)

				return m, nil
			}

			//* If game is complete, reset the game
			if m.game.Complete {
				m.game = vocab.NewGame(m.numerOfRound)
				round := m.game.Rounds[0]

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

				m.list = input.NewInput(items)

				return m, nil
			}

			//* If the game is in progress, mark the answer
			m.game.MarkAnswer(m.list.GetSelected().Value)
			m.list.Reveal()
			m.reveal = true

			return m, nil
		}
	}

	return m, nil
}

// Colors
const primary = "#5de4c7"
const primaryText = "#000000"

var subduedColor = lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
var headingStyles = lipgloss.NewStyle().Padding(0, 1).Margin(0, 1, 0, 0).Background(lipgloss.Color(primary)).Foreground(lipgloss.Color(primaryText))
var flashcardStyles = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#ffffff"))
var contentStyles = lipgloss.NewStyle().Margin(1, 0, 0, 1)

var secondaryText = lipgloss.NewStyle().
	Foreground(subduedColor)

func (m Main) View() string {
	if m.game.Complete || len(m.game.Rounds) == 0 {
		m.keys = titleKeys
	} else {
		m.keys = keys
	}

	helpView := m.help.View(m.keys)
	content := ""

	roundControls := lipgloss.JoinVertical(
		lipgloss.Left,
		secondaryText.MarginTop(3).Render("Select number of rounds"),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			secondaryText.MarginRight(1).MarginBottom(3).Render("-"),
			fmt.Sprintf("%d", m.numerOfRound),
			secondaryText.MarginLeft(1).Render("+"),
		),
	)

	if len(m.game.Rounds) == 0 {
		//* Game is not being played yet - render the main menu
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			headingStyles.Render("你好!"),
		)

		description := lipgloss.JoinHorizontal(
			lipgloss.Left,
			"Welcome to ",
			lipgloss.NewStyle().Background(lipgloss.Color(primary)).Foreground(lipgloss.Color(primaryText)).Render("电闪卡!"),
		)

		description = lipgloss.JoinVertical(
			lipgloss.Left,
			description,
			secondaryText.MarginTop(1).Render("A flashy little terminal app to practise Chinese characters."),
		)

		description = lipgloss.NewStyle().Width(50).MarginTop(3).Render(description)

		content += lipgloss.JoinVertical(
			lipgloss.Left,
			description,
			roundControls,
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

	current := m.game.Rounds[m.game.Round]

	if m.game.Complete {
		//* Game is finished - render the score
		score := m.game.GetScore() * 100

		scoreCard := lipgloss.JoinHorizontal(
			lipgloss.Left,
			"You scored ",
			flashcardStyles.Render(fmt.Sprintf("%.f", score)),
			" points.",
		)

		scoreCard = lipgloss.
			NewStyle().
			Width(50).
			MarginTop(3).
			Render(scoreCard)

		content += lipgloss.JoinVertical(
			lipgloss.Left,
			headingStyles.Render("下课!"),
			scoreCard,
			secondaryText.MarginTop(2).Render("Want to play again?"),
			roundControls,
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

	card := lipgloss.JoinVertical(
		lipgloss.Left,
		secondaryText.Render(fmt.Sprintf("Round %d of %d", m.game.Round+1, m.game.NumberOfRounds)),
		flashcardStyles.MarginTop(1).Render(current.Card.Chinese),
		secondaryText.Margin(1, 0, 0, 0).Render("Select the correct English translation"),
		cards,
		reveal,
	)

	card = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(subduedColor).
		MarginTop(3).
		Width(50).
		Render(card)

	content = content + lipgloss.JoinVertical(
		lipgloss.Left,
		headingStyles.Render("考试!"),
		card,
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
	programme := tea.NewProgram(newModel())
	_, err := programme.Run()

	if err != nil {
		fmt.Println("Error creating program:", err)
		os.Exit(1)
	}
}
