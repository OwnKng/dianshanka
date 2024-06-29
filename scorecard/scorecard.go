package scorecard

import (
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	rows    []Row
	columns []Column
}

type Row []string

type Column struct {
	Title string
	Width int
}

var subduedColor = lipgloss.
	AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}

var headingStyle = lipgloss.
	NewStyle().
	Foreground(lipgloss.Color("#FFFFFF"))

var text = lipgloss.
	NewStyle().
	Foreground(subduedColor)

func NewModel(rows []Row, columns []Column) *Model {
	m := &Model{
		rows:    rows,
		columns: columns,
	}

	return m
}

func (m *Model) renderRow(r int) string {
	s := make([]string, 0, len(m.columns))

	for i, value := range m.rows[r] {
		style := text.Width(m.columns[i].Width).MaxWidth(m.columns[i].Width).Inline(true)
		s = append(s, style.Render(string(value)))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Left, s...)

	return row
}

func (m *Model) renderHeading() string {
	s := make([]string, 0, len(m.columns))

	for _, column := range m.columns {
		style := headingStyle.Width(column.Width).MaxWidth(column.Width).Bold(true)
		s = append(s, style.Render(column.Title))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, s...)
}

func (m *Model) Render() string {
	s := make([]string, 0, len(m.rows))

	for i, _ := range m.rows {
		s = append(s, m.renderRow(i))
	}

	headings := m.renderHeading()

	table := lipgloss.JoinVertical(lipgloss.Left, headings, lipgloss.JoinVertical(lipgloss.Left, s...))

	return lipgloss.
		NewStyle().
		Padding(0, 2, 0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(subduedColor).
		Render(table)
}
