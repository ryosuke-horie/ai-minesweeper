package solver

import (
	"testing"

	"github.com/r-horie/ai-minesweeper/game"
)

func TestNewSolver(t *testing.T) {
	board := game.NewBoard(5, 5, 5)
	solver := NewSolver(board)

	if solver == nil {
		t.Fatal("NewSolver() returned nil")
	}

	if solver.board != board {
		t.Error("NewSolver() did not set board correctly")
	}
}

func TestSolver_Solve_Empty(t *testing.T) {
	// すべて未開放の盤面では推論できない
	board := game.NewBoard(5, 5, 5)
	solver := NewSolver(board)

	result := solver.Solve()

	if len(result.SafeCells) != 0 {
		t.Errorf("Expected 0 safe cells, got %d", len(result.SafeCells))
	}

	if len(result.MineCells) != 0 {
		t.Errorf("Expected 0 mine cells, got %d", len(result.MineCells))
	}

	if result.CanProgress {
		t.Error("Expected CanProgress to be false")
	}
}

func TestSolver_findDefiniteMines(t *testing.T) { //nolint:gocyclo // テストケースが多いため複雑度が高い
	tests := []struct {
		name       string
		setupBoard func() *game.Board
		wantMines  []game.Position
	}{
		{
			name: "simple case - 1 surrounded by 1 unrevealed",
			setupBoard: func() *game.Board {
				board := game.NewBoard(3, 3, 1)
				// 中央に1を配置、右に未開放のマス
				board.Cells[1][1].IsRevealed = true
				board.Cells[1][1].SetAdjacent(1)
				// 右以外を開く
				for i := 0; i < 3; i++ {
					for j := 0; j < 3; j++ {
						if !(i == 1 && j == 1) && !(i == 1 && j == 2) {
							board.Cells[i][j].IsRevealed = true
						}
					}
				}
				return board
			},
			wantMines: []game.Position{{Row: 1, Col: 2}},
		},
		{
			name: "2 surrounded by 2 unrevealed",
			setupBoard: func() *game.Board {
				board := game.NewBoard(3, 3, 2)
				// 中央に2を配置
				board.Cells[1][1].IsRevealed = true
				board.Cells[1][1].SetAdjacent(2)
				// 右と下だけ未開放
				for i := 0; i < 3; i++ {
					for j := 0; j < 3; j++ {
						if !((i == 1 && j == 1) || (i == 1 && j == 2) || (i == 2 && j == 1)) {
							board.Cells[i][j].IsRevealed = true
						}
					}
				}
				return board
			},
			wantMines: []game.Position{{Row: 1, Col: 2}, {Row: 2, Col: 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := tt.setupBoard()
			solver := NewSolver(board)

			mines := solver.findDefiniteMines()

			if len(mines) != len(tt.wantMines) {
				t.Errorf("Expected %d mines, got %d", len(tt.wantMines), len(mines))
				// デバッグ情報を追加
				t.Logf("Found mines: %v", mines)
				t.Logf("Expected mines: %v", tt.wantMines)
				// 中央セルの周囲の状況を確認
				centerPos := game.Position{Row: 1, Col: 1}
				unrevealed, flagged := solver.getUnrevealedAndFlaggedCounts(centerPos)
				t.Logf("Center cell: Adjacent=%d, Unrevealed=%d, Flagged=%d",
					board.Cells[1][1].Adjacent, unrevealed, flagged)
				t.Logf("Condition check: unrevealed(%d) == Adjacent(%d) - flagged(%d) = %d",
					unrevealed, board.Cells[1][1].Adjacent, flagged,
					board.Cells[1][1].Adjacent-flagged)
				return
			}

			// 順序に依存しないチェック
			for _, wantMine := range tt.wantMines {
				found := false
				for _, mine := range mines {
					if mine.Row == wantMine.Row && mine.Col == wantMine.Col {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected mine at %v not found", wantMine)
				}
			}
		})
	}
}

func TestSolver_findDefiniteSafeCells(t *testing.T) {
	tests := []struct {
		name       string
		setupBoard func() *game.Board
		wantSafes  []game.Position
	}{
		{
			name: "simple case - 1 with 1 flagged adjacent",
			setupBoard: func() *game.Board {
				board := game.NewBoard(3, 3, 1)
				// 中央に1を配置
				board.Cells[1][1].IsRevealed = true
				board.Cells[1][1].SetAdjacent(1)
				// 右にフラグ
				board.Cells[1][2].IsFlagged = true
				// 他は未開放
				return board
			},
			wantSafes: []game.Position{
				{Row: 0, Col: 0}, {Row: 0, Col: 1}, {Row: 0, Col: 2},
				{Row: 1, Col: 0},
				{Row: 2, Col: 0}, {Row: 2, Col: 1}, {Row: 2, Col: 2},
			},
		},
		{
			name: "3x3 with center 8 and all flagged",
			setupBoard: func() *game.Board {
				board := game.NewBoard(3, 3, 8)
				// 中央に8を配置（周囲全て地雷）
				board.Cells[1][1].IsRevealed = true
				board.Cells[1][1].SetAdjacent(8)
				// 周囲全てにフラグ
				for i := 0; i < 3; i++ {
					for j := 0; j < 3; j++ {
						if !(i == 1 && j == 1) {
							board.Cells[i][j].IsFlagged = true
						}
					}
				}
				return board
			},
			wantSafes: []game.Position{}, // すべてフラグ済みなので安全なセルなし
		},
		{
			name: "revealed mine counts as known mine",
			setupBoard: func() *game.Board {
				board := game.NewBoard(3, 3, 1)
				// 中央に1を配置
				board.Cells[1][1].IsRevealed = true
				board.Cells[1][1].SetAdjacent(1)
				// 右に開かれた地雷
				board.Cells[1][2].IsRevealed = true
				board.Cells[1][2].SetMine()
				return board
			},
			wantSafes: []game.Position{
				{Row: 0, Col: 0}, {Row: 0, Col: 1}, {Row: 0, Col: 2},
				{Row: 1, Col: 0},
				{Row: 2, Col: 0}, {Row: 2, Col: 1}, {Row: 2, Col: 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := tt.setupBoard()
			solver := NewSolver(board)

			safes := solver.findDefiniteSafeCells()

			if len(safes) != len(tt.wantSafes) {
				t.Errorf("Expected %d safe cells, got %d", len(tt.wantSafes), len(safes))
				return
			}

			// 順序に依存しないチェック
			for _, wantSafe := range tt.wantSafes {
				found := false
				for _, safe := range safes {
					if safe.Row == wantSafe.Row && safe.Col == wantSafe.Col {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected safe cell at %v not found", wantSafe)
				}
			}
		})
	}
}

func TestSolver_Solve_ComplexScenarios(t *testing.T) {
	tests := []struct {
		name         string
		setupBoard   func() *game.Board
		wantSafes    int
		wantMines    int
		wantProgress bool
	}{
		{
			name: "multiple numbers with overlapping constraints",
			setupBoard: func() *game.Board {
				board := game.NewBoard(5, 5, 5)
				// 1-2-1 パターンの簡略版
				// より単純なケース: 1つの数字セルで確実な推論
				board.Cells[1][1].IsRevealed = true
				board.Cells[1][1].SetAdjacent(1)

				// 周囲8マスのうち7マスを開く（数字を設定）
				positions := []struct {
					pos game.Position
					adj int
				}{
					{game.Position{Row: 0, Col: 0}, 1}, {game.Position{Row: 0, Col: 1}, 1}, {game.Position{Row: 0, Col: 2}, 1},
					{game.Position{Row: 1, Col: 0}, 1}, {game.Position{Row: 1, Col: 2}, 1},
					{game.Position{Row: 2, Col: 0}, 1}, {game.Position{Row: 2, Col: 1}, 1},
				}
				for _, p := range positions {
					board.Cells[p.pos.Row][p.pos.Col].IsRevealed = true
					board.Cells[p.pos.Row][p.pos.Col].SetAdjacent(p.adj)
				}
				// [2,2]だけ未開放 → これが地雷

				return board
			},
			wantMines:    1,    // [2,2]が地雷
			wantSafes:    0,    // 安全なセルは特定できない
			wantProgress: true, // 地雷が見つかるので進捗あり
		},
		{
			name: "corner pattern",
			setupBoard: func() *game.Board {
				board := game.NewBoard(5, 5, 3)
				// コーナーに1
				board.Cells[0][0].IsRevealed = true
				board.Cells[0][0].SetAdjacent(1)

				// 隣接する2マスを開く（数字を設定）
				board.Cells[0][1].IsRevealed = true
				board.Cells[0][1].SetAdjacent(1)
				board.Cells[1][0].IsRevealed = true
				board.Cells[1][0].SetAdjacent(1)

				return board
			},
			wantMines:    1, // [1,1]が地雷
			wantSafes:    0,
			wantProgress: true, // 地雷が見つかるので進捗あり
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := tt.setupBoard()
			solver := NewSolver(board)

			result := solver.Solve()

			if len(result.SafeCells) != tt.wantSafes {
				t.Errorf("Expected %d safe cells, got %d", tt.wantSafes, len(result.SafeCells))
			}

			if len(result.MineCells) != tt.wantMines {
				t.Errorf("Expected %d mine cells, got %d", tt.wantMines, len(result.MineCells))
			}

			if result.CanProgress != tt.wantProgress {
				t.Errorf("Expected CanProgress=%v, got %v", tt.wantProgress, result.CanProgress)
			}
		})
	}
}

func TestSolver_getUnrevealedAndFlaggedCounts(t *testing.T) {
	board := game.NewBoard(3, 3, 2)
	solver := NewSolver(board)

	// 中央のセルの周囲を設定
	// 上: 開いている
	board.Cells[0][1].IsRevealed = true
	// 右: フラグ
	board.Cells[1][2].IsFlagged = true
	// 下: 未開放
	// 左: 未開放

	unrevealed, flagged := solver.getUnrevealedAndFlaggedCounts(game.Position{Row: 1, Col: 1})

	if unrevealed != 7 { // 8 - 1(開いている)
		t.Errorf("Expected 7 unrevealed, got %d", unrevealed)
	}

	if flagged != 1 {
		t.Errorf("Expected 1 flagged, got %d", flagged)
	}
}

func TestSolver_getKnownMineCount(t *testing.T) {
	board := game.NewBoard(3, 3, 3)
	solver := NewSolver(board)

	// 中央のセルの周囲を設定
	// 上: フラグ（既知の地雷）
	board.Cells[0][1].IsFlagged = true
	// 右: 開かれた地雷
	board.Cells[1][2].IsRevealed = true
	board.Cells[1][2].SetMine()
	// 下: 未開放
	// 左: 開いている（地雷でない）
	board.Cells[1][0].IsRevealed = true

	count := solver.getKnownMineCount(game.Position{Row: 1, Col: 1})

	if count != 2 {
		t.Errorf("Expected 2 known mines, got %d", count)
	}
}

func TestSolver_isKnownMine(t *testing.T) {
	board := game.NewBoard(3, 3, 2)
	solver := NewSolver(board)

	tests := []struct {
		name     string
		setup    func(pos game.Position)
		pos      game.Position
		expected bool
	}{
		{
			name: "flagged cell is known mine",
			setup: func(pos game.Position) {
				board.Cells[pos.Row][pos.Col].IsFlagged = true
			},
			pos:      game.Position{Row: 0, Col: 0},
			expected: true,
		},
		{
			name: "revealed mine is known mine",
			setup: func(pos game.Position) {
				board.Cells[pos.Row][pos.Col].IsRevealed = true
				board.Cells[pos.Row][pos.Col].SetMine()
			},
			pos:      game.Position{Row: 0, Col: 1},
			expected: true,
		},
		{
			name: "unrevealed cell is not known mine",
			setup: func(pos game.Position) {
				// 何もしない（デフォルトで未開放）
			},
			pos:      game.Position{Row: 0, Col: 2},
			expected: false,
		},
		{
			name: "revealed non-mine is not known mine",
			setup: func(pos game.Position) {
				board.Cells[pos.Row][pos.Col].IsRevealed = true
			},
			pos:      game.Position{Row: 1, Col: 0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(tt.pos)

			result := solver.isKnownMine(tt.pos)

			if result != tt.expected {
				t.Errorf("Expected isKnownMine=%v, got %v", tt.expected, result)
			}
		})
	}
}

func TestContainsPosition(t *testing.T) {
	positions := []game.Position{
		{Row: 0, Col: 0},
		{Row: 1, Col: 2},
		{Row: 3, Col: 4},
	}

	tests := []struct {
		name     string
		pos      game.Position
		expected bool
	}{
		{
			name:     "contains first position",
			pos:      game.Position{Row: 0, Col: 0},
			expected: true,
		},
		{
			name:     "contains middle position",
			pos:      game.Position{Row: 1, Col: 2},
			expected: true,
		},
		{
			name:     "does not contain position",
			pos:      game.Position{Row: 2, Col: 2},
			expected: false,
		},
		{
			name:     "does not contain with same row",
			pos:      game.Position{Row: 1, Col: 0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsPosition(positions, tt.pos)

			if result != tt.expected {
				t.Errorf("Expected containsPosition=%v, got %v", tt.expected, result)
			}
		})
	}
}
