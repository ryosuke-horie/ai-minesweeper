package game

import "time"

type GameState int

const (
	Playing GameState = iota
	Won
	Lost
)

type Difficulty struct {
	Name   string
	Width  int
	Height int
	Mines  int
}

var (
	Beginner     = Difficulty{"Beginner", 9, 9, 10}
	Intermediate = Difficulty{"Intermediate", 16, 16, 40}
	Expert       = Difficulty{"Expert", 30, 16, 99}
)

type Game struct {
	Board       *Board
	State       GameState
	FirstClick  bool
	Difficulty  Difficulty
	StartTime   int64
	ElapsedTime int64
}

func NewGame(difficulty Difficulty) *Game {
	return &Game{
		Board:      NewBoard(difficulty.Width, difficulty.Height, difficulty.Mines),
		State:      Playing,
		FirstClick: true,
		Difficulty: difficulty,
	}
}

func (g *Game) Click(pos Position) {
	if g.State != Playing {
		return
	}

	if g.FirstClick {
		g.Board.Initialize(pos)
		g.FirstClick = false
		g.StartTime = getCurrentTime()
	}

	hitMine := g.Board.RevealCell(pos)

	if hitMine {
		g.State = Lost
		g.revealAllMines()
	} else if g.Board.CountUnrevealedSafeCells() == 0 {
		g.State = Won
		g.ElapsedTime = getCurrentTime() - g.StartTime
	}
}

func (g *Game) ToggleFlag(pos Position) {
	if g.State != Playing {
		return
	}

	cell := g.Board.GetCell(pos)
	if cell != nil {
		cell.ToggleFlag()
	}
}

func (g *Game) Reset() {
	g.Board = NewBoard(g.Difficulty.Width, g.Difficulty.Height, g.Difficulty.Mines)
	g.State = Playing
	g.FirstClick = true
	g.StartTime = 0
	g.ElapsedTime = 0
}

func (g *Game) GetRemainingMines() int {
	flaggedCount := 0
	for i := 0; i < g.Board.Height; i++ {
		for j := 0; j < g.Board.Width; j++ {
			if g.Board.Cells[i][j].IsFlagged {
				flaggedCount++
			}
		}
	}
	return g.Board.Mines - flaggedCount
}

func (g *Game) revealAllMines() {
	for i := 0; i < g.Board.Height; i++ {
		for j := 0; j < g.Board.Width; j++ {
			cell := g.Board.Cells[i][j]
			if cell.IsMine {
				cell.IsRevealed = true
			}
		}
	}
}

func getCurrentTime() int64 {
	return time.Now().Unix()
}