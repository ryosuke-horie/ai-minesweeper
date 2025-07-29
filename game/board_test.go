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