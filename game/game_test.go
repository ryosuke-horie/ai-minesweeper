package game

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	tests := []struct {
		name       string
		difficulty Difficulty
	}{
		{
			name:       "beginner difficulty",
			difficulty: Beginner,
		},
		{
			name:       "intermediate difficulty",
			difficulty: Intermediate,
		},
		{
			name:       "expert difficulty",
			difficulty: Expert,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame(tt.difficulty)

			if game == nil {
				t.Fatal("NewGame() returned nil")
			}

			if game.Board == nil {
				t.Fatal("NewGame() created game with nil board")
			}

			if game.State != Playing {
				t.Errorf("State = %v, want Playing", game.State)
			}

			if !game.FirstClick {
				t.Error("FirstClick should be true for new game")
			}

			if game.Difficulty.Name != tt.difficulty.Name {
				t.Errorf("Difficulty = %v, want %v", game.Difficulty.Name, tt.difficulty.Name)
			}

			if game.Board.Width != tt.difficulty.Width {
				t.Errorf("Board.Width = %d, want %d", game.Board.Width, tt.difficulty.Width)
			}

			if game.Board.Height != tt.difficulty.Height {
				t.Errorf("Board.Height = %d, want %d", game.Board.Height, tt.difficulty.Height)
			}

			if game.Board.Mines != tt.difficulty.Mines {
				t.Errorf("Board.Mines = %d, want %d", game.Board.Mines, tt.difficulty.Mines)
			}
		})
	}
}

func TestGame_Reset(t *testing.T) {
	game := NewGame(Beginner)

	// ゲーム状態を変更
	game.State = Lost
	game.FirstClick = false
	game.StartTime = 12345
	game.ElapsedTime = 100

	// リセット
	game.Reset()

	// リセット後の検証
	if game.Board == nil {
		t.Fatal("Reset() resulted in nil board")
	}

	if game.State != Playing {
		t.Errorf("State after Reset() = %v, want Playing", game.State)
	}

	if !game.FirstClick {
		t.Error("FirstClick should be true after Reset()")
	}

	if game.StartTime != 0 {
		t.Errorf("StartTime after Reset() = %d, want 0", game.StartTime)
	}

	if game.ElapsedTime != 0 {
		t.Errorf("ElapsedTime after Reset() = %d, want 0", game.ElapsedTime)
	}

	// 難易度が維持されているか
	if game.Difficulty.Name != Beginner.Name {
		t.Errorf("Difficulty changed after Reset(): %v", game.Difficulty.Name)
	}
}

func TestGame_GetRemainingMines(t *testing.T) {
	game := NewGame(Beginner)

	tests := []struct {
		name         string
		flaggedCells []Position
		want         int
	}{
		{
			name:         "no flags",
			flaggedCells: []Position{},
			want:         10, // Beginner has 10 mines
		},
		{
			name: "one flag",
			flaggedCells: []Position{
				{Row: 1, Col: 1},
			},
			want: 9,
		},
		{
			name: "multiple flags",
			flaggedCells: []Position{
				{Row: 1, Col: 1},
				{Row: 2, Col: 2},
				{Row: 3, Col: 3},
			},
			want: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ゲームをリセットしてフラグをクリア
			game.Reset()

			// フラグを設定
			for _, pos := range tt.flaggedCells {
				cell := game.Board.GetCell(pos)
				if cell != nil {
					cell.IsFlagged = true
				}
			}

			remaining := game.GetRemainingMines()

			if remaining != tt.want {
				t.Errorf("GetRemainingMines() = %d, want %d", remaining, tt.want)
			}
		})
	}
}

func TestGame_ToggleFlag(t *testing.T) {
	game := NewGame(Beginner)
	pos := Position{Row: 1, Col: 1}

	tests := []struct {
		name      string
		gameState GameState
		wantFlag  bool
	}{
		{
			name:      "toggle flag during playing",
			gameState: Playing,
			wantFlag:  true,
		},
		{
			name:      "cannot toggle flag when won",
			gameState: Won,
			wantFlag:  false,
		},
		{
			name:      "cannot toggle flag when lost",
			gameState: Lost,
			wantFlag:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ゲームをリセット
			game.Reset()
			game.State = tt.gameState

			// フラグをトグル
			game.ToggleFlag(pos)

			cell := game.Board.GetCell(pos)
			if cell != nil && cell.IsFlagged != tt.wantFlag {
				t.Errorf("Cell flagged = %v, want %v", cell.IsFlagged, tt.wantFlag)
			}
		})
	}
}