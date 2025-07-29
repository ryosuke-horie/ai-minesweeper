package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("226")).
			PaddingLeft(1)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250")).
			PaddingLeft(1)

	cellStyle = lipgloss.NewStyle().
			Width(3).
			Height(1).
			Align(lipgloss.Center)

	unrevealedStyle = cellStyle.
			Background(lipgloss.Color("238")).
			Foreground(lipgloss.Color("250"))

	revealedStyle = cellStyle.
			Background(lipgloss.Color("235"))

	cursorStyle = cellStyle.
			Background(lipgloss.Color("33")).
			Foreground(lipgloss.Color("231"))

	mineStyle = cellStyle.
			Background(lipgloss.Color("196")).
			Foreground(lipgloss.Color("231"))

	flagStyle = cellStyle.
			Background(lipgloss.Color("238")).
			Foreground(lipgloss.Color("226"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			PaddingLeft(1)

	gameOverStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196")).
			PaddingLeft(1)

	gameWonStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("46")).
			PaddingLeft(1)

	numberColors = map[int]lipgloss.Color{
		1: lipgloss.Color("21"),
		2: lipgloss.Color("22"),
		3: lipgloss.Color("196"),
		4: lipgloss.Color("20"),
		5: lipgloss.Color("88"),
		6: lipgloss.Color("51"),
		7: lipgloss.Color("16"),
		8: lipgloss.Color("236"),
	}
)

func getNumberStyle(num int) lipgloss.Style {
	color, ok := numberColors[num]
	if !ok {
		color = lipgloss.Color("250")
	}
	return revealedStyle.Foreground(color)
}
