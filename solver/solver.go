package solver

import (
	"github.com/r-horie/ai-minesweeper/game"
)

type SolverResult struct {
	SafeCells   []game.Position
	MineCells   []game.Position
	CanProgress bool
}

type Solver struct {
	board         *game.Board
	virtualRevealed map[game.Position]bool
}

func NewSolver(board *game.Board) *Solver {
	return &Solver{
		board:           board,
		virtualRevealed: make(map[game.Position]bool),
	}
}

func (s *Solver) Solve() SolverResult {
	result := SolverResult{
		SafeCells:   []game.Position{},
		MineCells:   []game.Position{},
		CanProgress: false,
	}

	iteration := 0
	maxIterations := 100

	for iteration < maxIterations {
		iteration++
		
		mines := s.findDefiniteMines()
		for _, mine := range mines {
			if !containsPosition(result.MineCells, mine) {
				result.MineCells = append(result.MineCells, mine)
			}
		}

		safes := s.findDefiniteSafeCells()
		newSafes := []game.Position{}
		for _, safe := range safes {
			if !s.virtualRevealed[safe] && !containsPosition(result.SafeCells, safe) {
				newSafes = append(newSafes, safe)
				s.virtualRevealed[safe] = true
			}
		}

		if len(newSafes) == 0 {
			break
		}

		result.SafeCells = append(result.SafeCells, newSafes...)
	}

	result.CanProgress = len(result.SafeCells) > 0
	return result
}

func (s *Solver) findDefiniteMines() []game.Position {
	mines := []game.Position{}

	for row := 0; row < s.board.Height; row++ {
		for col := 0; col < s.board.Width; col++ {
			pos := game.Position{Row: row, Col: col}
			cell := s.board.GetCell(pos)

			if cell == nil || !cell.IsRevealed || cell.IsMine {
				continue
			}

			if cell.Adjacent == 0 {
				continue
			}

			unrevealed, flagged := s.getUnrevealedAndFlaggedCounts(pos)

			if unrevealed == cell.Adjacent-flagged && unrevealed > 0 {
				for _, adjPos := range s.board.GetAdjacentPositions(pos) {
					adjCell := s.board.GetCell(adjPos)
					if adjCell != nil && !adjCell.IsRevealed && !adjCell.IsFlagged && !s.virtualRevealed[adjPos] {
						if !containsPosition(mines, adjPos) {
							mines = append(mines, adjPos)
						}
					}
				}
			}
		}
	}

	return mines
}

func (s *Solver) findDefiniteSafeCells() []game.Position {
	safes := []game.Position{}

	for row := 0; row < s.board.Height; row++ {
		for col := 0; col < s.board.Width; col++ {
			pos := game.Position{Row: row, Col: col}
			cell := s.board.GetCell(pos)

			if cell == nil || !cell.IsRevealed || cell.IsMine {
				continue
			}

			if cell.Adjacent == 0 {
				continue
			}

			unrevealed, _ := s.getUnrevealedAndFlaggedCounts(pos)
			mineCount := s.getKnownMineCount(pos)

			if mineCount == cell.Adjacent && unrevealed > 0 {
				for _, adjPos := range s.board.GetAdjacentPositions(pos) {
					adjCell := s.board.GetCell(adjPos)
					if adjCell != nil && !adjCell.IsRevealed && !adjCell.IsFlagged && !s.isKnownMine(adjPos) && !s.virtualRevealed[adjPos] {
						if !containsPosition(safes, adjPos) {
							safes = append(safes, adjPos)
						}
					}
				}
			}
		}
	}

	return safes
}

func (s *Solver) getUnrevealedAndFlaggedCounts(pos game.Position) (unrevealed, flagged int) {
	for _, adjPos := range s.board.GetAdjacentPositions(pos) {
		cell := s.board.GetCell(adjPos)
		if cell != nil && !cell.IsRevealed && !s.virtualRevealed[adjPos] {
			unrevealed++
			if cell.IsFlagged {
				flagged++
			}
		}
	}
	return
}

func (s *Solver) getKnownMineCount(pos game.Position) int {
	count := 0
	for _, adjPos := range s.board.GetAdjacentPositions(pos) {
		if s.isKnownMine(adjPos) {
			count++
		}
	}
	return count
}

func (s *Solver) isKnownMine(pos game.Position) bool {
	cell := s.board.GetCell(pos)
	return cell != nil && (cell.IsFlagged || (cell.IsRevealed && cell.IsMine))
}

func containsPosition(positions []game.Position, pos game.Position) bool {
	for _, p := range positions {
		if p.Row == pos.Row && p.Col == pos.Col {
			return true
		}
	}
	return false
}