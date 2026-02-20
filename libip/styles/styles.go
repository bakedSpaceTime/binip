package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	purple       = lipgloss.Color("99")
	gray         = lipgloss.Color("245")
	lightGray    = lipgloss.Color("241")
	red          = lipgloss.Color("196")
	green        = lipgloss.Color("42")
	lightGreen   = lipgloss.Color("112")
	adaptiveGray = lipgloss.AdaptiveColor{Light: "241", Dark: "245"}

	tableHeaderStyle = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	HeaderStyle      = lipgloss.NewStyle().Background(purple).Bold(true).Align(lipgloss.Left)
	FooterStyle      = lipgloss.NewStyle().
				Align(lipgloss.Center, lipgloss.Bottom)
	ErrorStyle   = lipgloss.NewStyle().Foreground(red).Bold(true)
	StatusStyle  = lipgloss.NewStyle().Foreground(green)
	AccentStyle  = lipgloss.NewStyle().Foreground(lightGreen)
	InfoStyle    = lipgloss.NewStyle().Foreground(adaptiveGray)
	cellStyle    = lipgloss.NewStyle().Padding(0, 1)
	oddRowStyle  = cellStyle.Foreground(gray)
	evenRowStyle = cellStyle.Foreground(lightGray)
)

func StyledTable() *table.Table {
	return table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return tableHeaderStyle
			case row%2 == 0:
				return evenRowStyle
			default:
				return oddRowStyle
			}
		}).
		Headers("key", "value")
}

func BorderlessTable() *table.Table {
	return table.New().
		Border(lipgloss.NormalBorder()).
		BorderRow(false).
		BorderColumn(true).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false)
}
