package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/r-horie/ai-minesweeper/game"
	"github.com/r-horie/ai-minesweeper/solver"
)

type tickMsg time.Time

type solverMsg struct {
	result solver.SolverResult
}

type revealCellMsg struct {
	positions []game.Position
	index     int
}

type Model struct {
	game           *game.Game
	solver         *solver.Solver
	cursor         game.Position
	aiThinking     bool
	lastUpdate     time.Time
	pendingReveals []game.Position
}

func NewModel() Model {
	g := game.NewGame(game.Beginner)
	return Model{
		game:           g,
		solver:         nil,
		cursor:         game.Position{Row: 0, Col: 0},
		aiThinking:     false,
		lastUpdate:     time.Now(),
		pendingReveals: []game.Position{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.ClearScreen,
		tickCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/10, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *Model) runSolver() tea.Cmd {
	return func() tea.Msg {
		s := solver.NewSolver(m.game.Board)
		result := s.Solve()
		return solverMsg{result: result}
	}
}

func revealNextCell(positions []game.Position, index int) tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return revealCellMsg{positions: positions, index: index}
	})
}