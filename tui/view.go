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
	return titleStyle.Render("AIマインスイーパー - AIにネタバレされるマインスイーパー")
}

func (m Model) renderHeader() string {
	remainingMines := m.game.GetRemainingMines()
	elapsed := int(time.Since(time.Unix(m.game.StartTime, 0)).Seconds())
	if m.game.StartTime == 0 {
		elapsed = 0
	}

	header := fmt.Sprintf("地雷: %d  時間: %02d:%02d  難易度: %s",
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
		status = gameWonStyle.Render("🎉 おめでとうございます！クリアしました！")
	case game.Lost:
		status = gameOverStyle.Render("💥 ゲームオーバー！地雷を踏みました！")
	default:
		if m.aiThinking {
			status = headerStyle.Render("🤖 AIが考え中...")
		} else {
			status = headerStyle.Render("あなたの番です！運命の選択を...")
		}
	}
	return status
}

func (m Model) renderHelp() string {
	help := []string{
		"[↑↓←→] カーソル移動",
		"[スペース] マスを開く",
		"[f] 旗を立てる",
		"[r] 新しいゲーム",
		"[1/2/3] 難易度変更",
		"[q] 終了",
	}
	return helpStyle.Render(strings.Join(help, "  "))
}
