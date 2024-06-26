package input

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var subduedColor = lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}

var StatusBar = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#ffffff"})

var StatusEmpty = lipgloss.NewStyle().
	Foreground(subduedColor).
	BorderForeground(lipgloss.Color("#ffffff"))

var correct = lipgloss.NewStyle().Foreground(lipgloss.Color("#5de4c7"))
var incorrect = lipgloss.NewStyle().Foreground(lipgloss.Color("#F15152"))

type Input struct {
	Cursor   int
	Selected Item
	Items    []Item
	reveal   bool
}

type Item struct {
	Title    string
	Subtitle string
	Value    string
	Correct  bool
}

func NewInput(items []Item) Input {
	return Input{Cursor: 0, Items: items, reveal: false}
}

func (i *Input) GetSelected() Item {
	return i.Items[i.Cursor]
}

func (i *Input) SetCursor(cursor int) {
	i.Cursor = cursor
}

func (i *Input) Reveal() {
	i.reveal = true
}

func (i *Input) Up() {
	i.Cursor++
}

func (i *Input) Down() {
	i.Cursor--
}

func (input *Input) View() string {
	// Print out the items
	s := ""

	for i, item := range input.Items {
		cursor := " "

		if input.reveal {
			if item.Correct {
				s += fmt.Sprintf("%s\n", correct.Render("✔", " ", item.Title))
			} else if input.Cursor == i {
				s += fmt.Sprintf("%s\n", incorrect.Render("✘", " ", item.Title))
			} else {
				s += fmt.Sprintf("%s\n", StatusEmpty.Render(" ", " ", item.Title))
			}
		} else {

			if input.Cursor == i {
				cursor = ">"
				s += fmt.Sprintf("%s\n", StatusBar.Render(cursor, " ", StatusBar.Render(item.Title)))
			} else {
				s += fmt.Sprintf("%s\n", StatusEmpty.Render(cursor, " ", item.Title))
			}
		}
	}

	return s
}
