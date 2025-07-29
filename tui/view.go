package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/r-horie/ai-minesweeper/game"
)

func (m Model) View() string {
	var sections []string

	sections = append(sections, m.renderTitle())
	sections = append(sections, m.renderHeader())
	sections = append(sections, m.renderBoard())
	sections = append(sections, m.renderStatus())
	sections = append(sections, m.renderHelp())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderTitle() string {
	return titleStyle.Render("AI Minesweeper - Spoiled by AI")
}

func (m Model) renderHeader() string {
	remainingMines := m.game.GetRemainingMines()
	elapsed := int(time.Since(time.Unix(m.game.StartTime, 0)).Seconds())
	if m.game.StartTime == 0 {
		elapsed = 0
	}

	header := fmt.Sprintf("Mines: %d  Time: %02d:%02d  Difficulty: %s",
		remainingMines,
		elapsed/60,
		elapsed%60,
		m.game.Difficulty.Name,
	)
	return headerStyle.Render(header)
}

func (m Model) renderBoard() string {
	var rows []string

	for row := 0; row < m.game.Board.Height; row++ {
		var cells []string
		for col := 0; col < m.game.Board.Width; col++ {
			pos := game.Position{Row: row, Col: col}
			cells = append(cells, m.renderCell(pos))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}

	board := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return lipgloss.NewStyle().PaddingLeft(1).Render(board)
}

func (m Model) renderCell(pos game.Position) string {
	cell := m.game.Board.GetCell(pos)
	if cell == nil {
		return cellStyle.Render(" ")
	}

	isCursor := pos.Row == m.cursor.Row && pos.Col == m.cursor.Col

	var style lipgloss.Style
	var content string

	if m.game.State == game.Lost && cell.IsMine && cell.IsRevealed {
		style = mineStyle
		content = "*"
	} else if cell.IsFlagged {
		style = flagStyle
		if isCursor {
			style = style.Copy().Background(lipgloss.Color("33"))
		}
		content = "F"
	} else if !cell.IsRevealed {
		style = unrevealedStyle
		if isCursor {
			style = cursorStyle
		}
		content = " "
	} else if cell.IsMine {
		style = mineStyle
		content = "*"
	} else if cell.Adjacent == 0 {
		style = revealedStyle
		if isCursor {
			style = style.Copy().Background(lipgloss.Color("239"))
		}
		content = " "
	} else {
		style = getNumberStyle(cell.Adjacent)
		if isCursor {
			style = style.Copy().Background(lipgloss.Color("239"))
		}
		content = fmt.Sprintf("%d", cell.Adjacent)
	}

	return style.Render(content)
}

func (m Model) renderStatus() string {
	var status string
	switch m.game.State {
	case game.Won:
		status = gameWonStyle.Render("🎉 Congratulations! You won!")
	case game.Lost:
		status = gameOverStyle.Render("💥 Game Over! You hit a mine!")
	default:
		if m.aiThinking {
			status = headerStyle.Render("🤖 AI is thinking...")
		} else {
			status = headerStyle.Render("Your turn! Choose wisely...")
		}
	}
	return status
}

func (m Model) renderHelp() string {
	help := []string{
		"[↑↓←→] Move cursor",
		"[Space] Reveal cell",
		"[f] Toggle flag",
		"[r] New game",
		"[1/2/3] Difficulty",
		"[q] Quit",
	}
	return helpStyle.Render(strings.Join(help, "  "))
}