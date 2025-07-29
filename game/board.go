package game

import (
	"math/rand"
	"time"
)

type Position struct {
	Row int
	Col int
}

type Board struct {
	Width  int
	Height int
	Mines  int
	Cells  [][]*Cell
}

func NewBoard(width, height, mines int) *Board {
	if mines > width*height {
		mines = width * height
	}

	cells := make([][]*Cell, height)
	for i := range cells {
		cells[i] = make([]*Cell, width)
		for j := range cells[i] {
			cells[i][j] = NewCell()
		}
	}

	return &Board{
		Width:  width,
		Height: height,
		Mines:  mines,
		Cells:  cells,
	}
}

func (b *Board) Initialize(firstClick Position) {
	rand.Seed(time.Now().UnixNano())

	mineCount := 0
	for mineCount < b.Mines {
		row := rand.Intn(b.Height)
		col := rand.Intn(b.Width)

		if row == firstClick.Row && col == firstClick.Col {
			continue
		}

		if abs(row-firstClick.Row) <= 1 && abs(col-firstClick.Col) <= 1 {
			continue
		}

		if !b.Cells[row][col].IsMine {
			b.Cells[row][col].SetMine()
			mineCount++
		}
	}

	for i := 0; i < b.Height; i++ {
		for j := 0; j < b.Width; j++ {
			if !b.Cells[i][j].IsMine {
				count := b.countAdjacentMines(Position{i, j})
				b.Cells[i][j].SetAdjacent(count)
			}
		}
	}
}

func (b *Board) GetCell(pos Position) *Cell {
	if b.IsValidPosition(pos) {
		return b.Cells[pos.Row][pos.Col]
	}
	return nil
}

func (b *Board) IsValidPosition(pos Position) bool {
	return pos.Row >= 0 && pos.Row < b.Height && pos.Col >= 0 && pos.Col < b.Width
}

func (b *Board) GetAdjacentPositions(pos Position) []Position {
	positions := []Position{}
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			newPos := Position{pos.Row + dr, pos.Col + dc}
			if b.IsValidPosition(newPos) {
				positions = append(positions, newPos)
			}
		}
	}
	return positions
}

func (b *Board) countAdjacentMines(pos Position) int {
	count := 0
	for _, adjPos := range b.GetAdjacentPositions(pos) {
		if b.Cells[adjPos.Row][adjPos.Col].IsMine {
			count++
		}
	}
	return count
}

func (b *Board) RevealCell(pos Position) bool {
	cell := b.GetCell(pos)
	if cell == nil || cell.IsRevealed || cell.IsFlagged {
		return false
	}

	cell.Reveal()

	if cell.IsMine {
		return true
	}

	if cell.Adjacent == 0 {
		for _, adjPos := range b.GetAdjacentPositions(pos) {
			b.RevealCell(adjPos)
		}
	}

	return false
}

func (b *Board) CountUnrevealedSafeCells() int {
	count := 0
	for i := 0; i < b.Height; i++ {
		for j := 0; j < b.Width; j++ {
			cell := b.Cells[i][j]
			if !cell.IsRevealed && !cell.IsMine {
				count++
			}
		}
	}
	return count
}

func (b *Board) GetAllUnrevealedPositions() []Position {
	positions := []Position{}
	for i := 0; i < b.Height; i++ {
		for j := 0; j < b.Width; j++ {
			if !b.Cells[i][j].IsRevealed {
				positions = append(positions, Position{i, j})
			}
		}
	}
	return positions
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
