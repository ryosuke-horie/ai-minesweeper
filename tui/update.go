package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/r-horie/ai-minesweeper/game"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "ctrl+q":
			return m, tea.Quit
		}

		if m.aiThinking {
			return m, nil
		}

		switch msg.String() {

		case "up", "k":
			if m.cursor.Row > 0 {
				m.cursor.Row--
			}

		case "down", "j":
			if m.cursor.Row < m.game.Board.Height-1 {
				m.cursor.Row++
			}

		case "left", "h":
			if m.cursor.Col > 0 {
				m.cursor.Col--
			}

		case "right", "l":
			if m.cursor.Col < m.game.Board.Width-1 {
				m.cursor.Col++
			}

		case " ", "enter":
			if m.game.State == game.Playing {
				m.game.Click(m.cursor)
				if m.game.State == game.Playing {
					m.aiThinking = true
					return m, m.runSolver()
				}
			}

		case "f":
			if m.game.State == game.Playing {
				m.game.ToggleFlag(m.cursor)
			}

		case "r":
			m.game.Reset()
			m.solver = nil
			m.cursor = game.Position{Row: 0, Col: 0}
			m.aiThinking = false

		case "1":
			m.game.Difficulty = game.Beginner
			m.game.Reset()
			m.solver = nil
			m.cursor = game.Position{Row: 0, Col: 0}
			m.aiThinking = false

		case "2":
			m.game.Difficulty = game.Intermediate
			m.game.Reset()
			m.solver = nil
			m.cursor = game.Position{Row: 0, Col: 0}
			m.aiThinking = false

		case "3":
			m.game.Difficulty = game.Expert
			m.game.Reset()
			m.solver = nil
			m.cursor = game.Position{Row: 0, Col: 0}
			m.aiThinking = false
		}

	case solverMsg:
		m.aiThinking = false
		result := msg.result

		for _, minePos := range result.MineCells {
			cell := m.game.Board.GetCell(minePos)
			if cell != nil && !cell.IsFlagged {
				cell.IsFlagged = true
			}
		}

		for _, safePos := range result.SafeCells {
			m.game.Board.RevealCell(safePos)
		}

		if m.game.Board.CountUnrevealedSafeCells() == 0 {
			m.game.State = game.Won
		} else if len(result.SafeCells) > 0 {
			m.aiThinking = true
			return m, m.runSolver()
		}

	case tickMsg:
		return m, tickCmd()
	}

	return m, nil
}