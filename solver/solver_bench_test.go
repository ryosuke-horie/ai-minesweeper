package solver

import (
	"testing"

	"github.com/r-horie/ai-minesweeper/game"
)

func BenchmarkSolver_Solve(b *testing.B) {
	// 典型的なゲーム状況を設定
	board := game.NewBoard(16, 16, 40) // 中級

	// いくつかのセルを開く
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if (i+j)%3 == 0 {
				board.Cells[i][j].IsRevealed = true
				board.Cells[i][j].SetAdjacent((i + j) % 4)
			}
		}
	}

	solver := NewSolver(board)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solver.Solve()
	}
}

func BenchmarkSolver_findDefiniteMines(b *testing.B) {
	board := createComplexBoard()
	solver := NewSolver(board)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solver.findDefiniteMines()
	}
}

func BenchmarkSolver_findDefiniteSafeCells(b *testing.B) {
	board := createComplexBoard()
	solver := NewSolver(board)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		solver.findDefiniteSafeCells()
	}
}

// 複雑な盤面を作成するヘルパー関数..
func createComplexBoard() *game.Board {
	board := game.NewBoard(30, 16, 99) // 上級

	// パターンを作成
	patterns := []struct {
		row, col int
		adjacent int
	}{
		{5, 5, 1}, {5, 6, 2}, {5, 7, 3},
		{10, 10, 2}, {10, 11, 3}, {10, 12, 2},
		{15, 15, 1}, {15, 16, 1}, {15, 17, 1},
	}

	for _, p := range patterns {
		if p.row < board.Height && p.col < board.Width {
			board.Cells[p.row][p.col].IsRevealed = true
			board.Cells[p.row][p.col].SetAdjacent(p.adjacent)

			// 周囲にいくつかフラグを配置
			if p.adjacent > 1 && p.row > 0 && p.col > 0 {
				board.Cells[p.row-1][p.col-1].IsFlagged = true
			}
		}
	}

	return board
}

func BenchmarkContainsPosition(b *testing.B) {
	// 最悪ケース: 大きなリストの最後を探す
	positions := make([]game.Position, 100)
	for i := 0; i < 100; i++ {
		positions[i] = game.Position{Row: i / 10, Col: i % 10}
	}

	target := game.Position{Row: 9, Col: 9} // 最後の要素

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		containsPosition(positions, target)
	}
}
