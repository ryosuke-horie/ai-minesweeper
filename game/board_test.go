package game

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		mines  int
		want   struct {
			width  int
			height int
			mines  int
		}
	}{
		{
			name:   "standard board",
			width:  9,
			height: 9,
			mines:  10,
			want:   struct{ width, height, mines int }{9, 9, 10},
		},
		{
			name:   "large board",
			width:  30,
			height: 16,
			mines:  99,
			want:   struct{ width, height, mines int }{30, 16, 99},
		},
		{
			name:   "mines exceed cells",
			width:  3,
			height: 3,
			mines:  100,
			want:   struct{ width, height, mines int }{3, 3, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := NewBoard(tt.width, tt.height, tt.mines)

			if board == nil {
				t.Fatal("NewBoard() returned nil")
			}

			if board.Width != tt.want.width {
				t.Errorf("Width = %d, want %d", board.Width, tt.want.width)
			}

			if board.Height != tt.want.height {
				t.Errorf("Height = %d, want %d", board.Height, tt.want.height)
			}

			if board.Mines != tt.want.mines {
				t.Errorf("Mines = %d, want %d", board.Mines, tt.want.mines)
			}

			if len(board.Cells) != tt.want.height {
				t.Errorf("Cells rows = %d, want %d", len(board.Cells), tt.want.height)
			}

			for i, row := range board.Cells {
				if len(row) != tt.want.width {
					t.Errorf("Cells[%d] cols = %d, want %d", i, len(row), tt.want.width)
				}

				for j, cell := range row {
					if cell == nil {
						t.Errorf("Cell[%d][%d] is nil", i, j)
					}
				}
			}
		})
	}
}

func TestBoard_GetCell(t *testing.T) {
	board := NewBoard(5, 5, 5)

	tests := []struct {
		name     string
		pos      Position
		wantCell bool
	}{
		{
			name:     "valid position",
			pos:      Position{Row: 2, Col: 3},
			wantCell: true,
		},
		{
			name:     "top-left corner",
			pos:      Position{Row: 0, Col: 0},
			wantCell: true,
		},
		{
			name:     "bottom-right corner",
			pos:      Position{Row: 4, Col: 4},
			wantCell: true,
		},
		{
			name:     "negative row",
			pos:      Position{Row: -1, Col: 2},
			wantCell: false,
		},
		{
			name:     "negative col",
			pos:      Position{Row: 2, Col: -1},
			wantCell: false,
		},
		{
			name:     "row overflow",
			pos:      Position{Row: 5, Col: 2},
			wantCell: false,
		},
		{
			name:     "col overflow",
			pos:      Position{Row: 2, Col: 5},
			wantCell: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := board.GetCell(tt.pos)

			if tt.wantCell && cell == nil {
				t.Errorf("GetCell(%v) = nil, want cell", tt.pos)
			}

			if !tt.wantCell && cell != nil {
				t.Errorf("GetCell(%v) = cell, want nil", tt.pos)
			}
		})
	}
}

func TestBoard_IsValidPosition(t *testing.T) {
	board := NewBoard(5, 5, 5)

	tests := []struct {
		name  string
		pos   Position
		want  bool
	}{
		{
			name: "valid center",
			pos:  Position{Row: 2, Col: 2},
			want: true,
		},
		{
			name: "valid top-left",
			pos:  Position{Row: 0, Col: 0},
			want: true,
		},
		{
			name: "valid bottom-right",
			pos:  Position{Row: 4, Col: 4},
			want: true,
		},
		{
			name: "invalid negative row",
			pos:  Position{Row: -1, Col: 2},
			want: false,
		},
		{
			name: "invalid negative col",
			pos:  Position{Row: 2, Col: -1},
			want: false,
		},
		{
			name: "invalid row overflow",
			pos:  Position{Row: 5, Col: 2},
			want: false,
		},
		{
			name: "invalid col overflow",
			pos:  Position{Row: 2, Col: 5},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := board.IsValidPosition(tt.pos); got != tt.want {
				t.Errorf("IsValidPosition(%v) = %v, want %v", tt.pos, got, tt.want)
			}
		})
	}
}

func TestBoard_GetAdjacentPositions(t *testing.T) {
	board := NewBoard(5, 5, 5)

	tests := []struct {
		name      string
		pos       Position
		wantCount int
	}{
		{
			name:      "center cell has 8 adjacent",
			pos:       Position{Row: 2, Col: 2},
			wantCount: 8,
		},
		{
			name:      "corner cell has 3 adjacent",
			pos:       Position{Row: 0, Col: 0},
			wantCount: 3,
		},
		{
			name:      "edge cell has 5 adjacent",
			pos:       Position{Row: 0, Col: 2},
			wantCount: 5,
		},
		{
			name:      "bottom-right corner",
			pos:       Position{Row: 4, Col: 4},
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			positions := board.GetAdjacentPositions(tt.pos)

			if len(positions) != tt.wantCount {
				t.Errorf("GetAdjacentPositions(%v) returned %d positions, want %d", tt.pos, len(positions), tt.wantCount)
			}

			// 隣接位置が正しいかチェック
			for _, adjPos := range positions {
				dr := abs(adjPos.Row - tt.pos.Row)
				dc := abs(adjPos.Col - tt.pos.Col)

				if dr > 1 || dc > 1 {
					t.Errorf("Position %v is not adjacent to %v", adjPos, tt.pos)
				}

				if dr == 0 && dc == 0 {
					t.Errorf("Position %v is the same as original %v", adjPos, tt.pos)
				}

				if !board.IsValidPosition(adjPos) {
					t.Errorf("Adjacent position %v is not valid", adjPos)
				}
			}
		})
	}
}

// Phase 2 Tests

func TestBoard_Initialize(t *testing.T) {
	tests := []struct {
		name        string
		boardSize   struct{ width, height, mines int }
		firstClick  Position
		wantSuccess bool
	}{
		{
			name: "standard initialization",
			boardSize: struct{ width, height, mines int }{
				width: 9, height: 9, mines: 10,
			},
			firstClick:  Position{Row: 4, Col: 4},
			wantSuccess: true,
		},
		{
			name: "corner first click",
			boardSize: struct{ width, height, mines int }{
				width: 5, height: 5, mines: 5,
			},
			firstClick:  Position{Row: 0, Col: 0},
			wantSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := NewBoard(tt.boardSize.width, tt.boardSize.height, tt.boardSize.mines)
			board.Initialize(tt.firstClick)

			// 地雷の総数をカウント
			mineCount := 0
			for i := 0; i < board.Height; i++ {
				for j := 0; j < board.Width; j++ {
					if board.Cells[i][j].IsMine {
						mineCount++
					}
				}
			}

			// 地雷数が正しいか
			if mineCount != board.Mines {
				t.Errorf("Mine count = %d, want %d", mineCount, board.Mines)
			}

			// 初回クリック位置とその周囲に地雷がないか
			firstClickCell := board.GetCell(tt.firstClick)
			if firstClickCell != nil && firstClickCell.IsMine {
				t.Error("First click position has mine")
			}

			// 初回クリック位置の周囲8マスに地雷がないか
			for _, adjPos := range board.GetAdjacentPositions(tt.firstClick) {
				adjCell := board.GetCell(adjPos)
				if adjCell != nil && adjCell.IsMine {
					t.Errorf("Adjacent position %v to first click has mine", adjPos)
				}
			}

			// 隣接地雷数が正しく計算されているか
			for i := 0; i < board.Height; i++ {
				for j := 0; j < board.Width; j++ {
					cell := board.Cells[i][j]
					if !cell.IsMine {
						expectedAdjacent := countAdjacentMinesForTest(board, Position{i, j})
						if cell.Adjacent != expectedAdjacent {
							t.Errorf("Cell[%d][%d] adjacent = %d, want %d", i, j, cell.Adjacent, expectedAdjacent)
						}
					}
				}
			}
		})
	}
}

// テスト用のヘルパー関数
func countAdjacentMinesForTest(board *Board, pos Position) int {
	count := 0
	for _, adjPos := range board.GetAdjacentPositions(pos) {
		if board.Cells[adjPos.Row][adjPos.Col].IsMine {
			count++
		}
	}
	return count
}

func TestBoard_RevealCell(t *testing.T) {
	tests := []struct {
		name           string
		setupBoard     func() *Board
		revealPos      Position
		wantHitMine    bool
		wantRevealed   []Position
	}{
		{
			name: "reveal safe cell with adjacent mines",
			setupBoard: func() *Board {
				// 3x3の盤面で中央に数字のマスを作る
				board := NewBoard(3, 3, 1)
				// 手動で地雷を配置
				board.Cells[0][0].SetMine()
				board.Cells[1][1].SetAdjacent(1)
				return board
			},
			revealPos:    Position{Row: 1, Col: 1},
			wantHitMine:  false,
			wantRevealed: []Position{{Row: 1, Col: 1}},
		},
		{
			name: "reveal empty cell triggers chain reaction",
			setupBoard: func() *Board {
				// 5x5の盤面で連鎖開放をテスト
				board := NewBoard(5, 5, 2)
				// 右上と右下に地雷を配置
				board.Cells[0][4].SetMine()
				board.Cells[4][4].SetMine()
				// 隣接数を更新
				for i := 0; i < 5; i++ {
					for j := 0; j < 5; j++ {
						if !board.Cells[i][j].IsMine {
							count := 0
							for _, adj := range board.GetAdjacentPositions(Position{i, j}) {
								if board.Cells[adj.Row][adj.Col].IsMine {
									count++
								}
							}
							board.Cells[i][j].SetAdjacent(count)
						}
					}
				}
				return board
			},
			revealPos:   Position{Row: 0, Col: 0},
			wantHitMine: false,
			// 左側のエリアが連鎖的に開く
			wantRevealed: []Position{
				{0, 0}, {0, 1}, {0, 2},
				{1, 0}, {1, 1}, {1, 2},
				{2, 0}, {2, 1}, {2, 2},
			},
		},
		{
			name: "reveal mine",
			setupBoard: func() *Board {
				board := NewBoard(3, 3, 1)
				board.Cells[1][1].SetMine()
				return board
			},
			revealPos:    Position{Row: 1, Col: 1},
			wantHitMine:  true,
			wantRevealed: []Position{{Row: 1, Col: 1}},
		},
		{
			name: "cannot reveal flagged cell",
			setupBoard: func() *Board {
				board := NewBoard(3, 3, 0)
				board.Cells[1][1].IsFlagged = true
				return board
			},
			revealPos:    Position{Row: 1, Col: 1},
			wantHitMine:  false,
			wantRevealed: []Position{}, // フラグがあるので開かない
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := tt.setupBoard()
			hitMine := board.RevealCell(tt.revealPos)

			if hitMine != tt.wantHitMine {
				t.Errorf("RevealCell() hit mine = %v, want %v", hitMine, tt.wantHitMine)
			}

			// 期待される位置が開いているか確認
			for _, pos := range tt.wantRevealed {
				cell := board.GetCell(pos)
				if cell != nil && !cell.IsRevealed {
					t.Errorf("Expected position %v to be revealed", pos)
				}
			}
		})
	}
}

func TestBoard_CountUnrevealedSafeCells(t *testing.T) {
	tests := []struct {
		name       string
		setupBoard func() *Board
		want       int
	}{
		{
			name: "new board",
			setupBoard: func() *Board {
				board := NewBoard(3, 3, 2)
				// NewBoardは地雷を配置しないので、すべてのセルが安全
				return board
			},
			want: 9, // すべてのセルが安全（地雷なし）
		},
		{
			name: "partially revealed board",
			setupBoard: func() *Board {
				board := NewBoard(3, 3, 2)
				board.Cells[0][0].SetMine()
				board.Cells[0][1].SetMine()
				// 2つのセルを開く
				board.Cells[1][0].IsRevealed = true
				board.Cells[1][1].IsRevealed = true
				return board
			},
			want: 5, // 7 safe cells - 2 revealed
		},
		{
			name: "all safe cells revealed",
			setupBoard: func() *Board {
				board := NewBoard(3, 3, 1)
				board.Cells[0][0].SetMine()
				// 地雷以外を全て開く
				for i := 0; i < 3; i++ {
					for j := 0; j < 3; j++ {
						if !board.Cells[i][j].IsMine {
							board.Cells[i][j].IsRevealed = true
						}
					}
				}
				return board
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := tt.setupBoard()
			got := board.CountUnrevealedSafeCells()

			if got != tt.want {
				t.Errorf("CountUnrevealedSafeCells() = %d, want %d", got, tt.want)
			}
		})
	}
}