package main

import (
	"fmt"
	"net/http"
	"os"

	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type statusMsg int
type errorMsg struct{ err error }

type model struct {
	status int
	err    error
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {

	return checkServer
}

const url = "https://ownkng.dev"

func checkServer() tea.Msg {

	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(url)

	if err != nil {
		return errorMsg{err}
	}

	return statusMsg(res.StatusCode)
}

func (e errorMsg) Error() string { return e.err.Error() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case statusMsg:
		m.status = int(msg)
		return m, tea.Quit

	case errorMsg:
		m.err = msg
		return m, tea.Quit

	case tea.KeyMsg:

		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\n We had some trouble: %v\n\n", m.err)
	}

	s := fmt.Sprintf("Checking %s ...", url)

	if m.status > 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}

	return "\n" + s + "\n\n"
}

func main() {

	p := tea.NewProgram(initialModel())

	_, err := p.Run()

	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
