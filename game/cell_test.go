package game

import (
	"testing"
)

func TestNewCell(t *testing.T) {
	cell := NewCell()

	if cell == nil {
		t.Fatal("NewCell() returned nil")
	}

	if cell.IsMine {
		t.Error("NewCell() should create a cell without mine")
	}

	if cell.IsRevealed {
		t.Error("NewCell() should create an unrevealed cell")
	}

	if cell.IsFlagged {
		t.Error("NewCell() should create an unflagged cell")
	}

	if cell.Adjacent != 0 {
		t.Error("NewCell() should create a cell with 0 adjacent mines")
	}
}

func TestCell_Reveal(t *testing.T) {
	tests := []struct {
		name       string
		isFlagged  bool
		wantReveal bool
	}{
		{
			name:       "reveal unflagged cell",
			isFlagged:  false,
			wantReveal: true,
		},
		{
			name:       "cannot reveal flagged cell",
			isFlagged:  true,
			wantReveal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := NewCell()
			if tt.isFlagged {
				cell.IsFlagged = true
			}

			cell.Reveal()

			if cell.IsRevealed != tt.wantReveal {
				t.Errorf("Reveal() result = %v, want %v", cell.IsRevealed, tt.wantReveal)
			}
		})
	}
}

func TestCell_ToggleFlag(t *testing.T) {
	tests := []struct {
		name         string
		isRevealed   bool
		initialFlag  bool
		wantFlag     bool
	}{
		{
			name:         "toggle flag on unrevealed cell",
			isRevealed:   false,
			initialFlag:  false,
			wantFlag:     true,
		},
		{
			name:         "untoggle flag on unrevealed cell",
			isRevealed:   false,
			initialFlag:  true,
			wantFlag:     false,
		},
		{
			name:         "cannot toggle flag on revealed cell",
			isRevealed:   true,
			initialFlag:  false,
			wantFlag:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := NewCell()
			cell.IsRevealed = tt.isRevealed
			cell.IsFlagged = tt.initialFlag

			cell.ToggleFlag()

			if cell.IsFlagged != tt.wantFlag {
				t.Errorf("ToggleFlag() result = %v, want %v", cell.IsFlagged, tt.wantFlag)
			}
		})
	}
}

func TestCell_SetMine(t *testing.T) {
	cell := NewCell()

	if cell.IsMine {
		t.Error("Cell should not have mine initially")
	}

	cell.SetMine()

	if !cell.IsMine {
		t.Error("SetMine() should set IsMine to true")
	}
}

func TestCell_SetAdjacent(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{"zero adjacent", 0},
		{"one adjacent", 1},
		{"eight adjacent", 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := NewCell()

			cell.SetAdjacent(tt.count)

			if cell.Adjacent != tt.count {
				t.Errorf("SetAdjacent(%d) result = %d, want %d", tt.count, cell.Adjacent, tt.count)
			}
		})
	}
}