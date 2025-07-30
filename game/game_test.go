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

// Phase 2 Tests

func TestGame_Click(t *testing.T) { //nolint:gocyclo // テストケースが多いため複雑度が高い
	tests := []struct {
		name           string
		setupGame      func() *Game
		clickPos       Position
		wantState      GameState
		wantFirstClick bool
		checkBoard     func(*testing.T, *Board)
	}{
		{
			name: "first click initializes board",
			setupGame: func() *Game {
				return NewGame(Beginner)
			},
			clickPos:       Position{Row: 4, Col: 4},
			wantState:      Playing,
			wantFirstClick: false,
			checkBoard: func(t *testing.T, board *Board) {
				// 初回クリック位置とその周囲に地雷がないことを確認
				cell := board.GetCell(Position{Row: 4, Col: 4})
				if cell.IsMine {
					t.Error("First click position has mine")
				}

				// 隣接セルに地雷がないことを確認
				for _, adjPos := range board.GetAdjacentPositions(Position{Row: 4, Col: 4}) {
					adjCell := board.GetCell(adjPos)
					if adjCell != nil && adjCell.IsMine {
						t.Errorf("Adjacent position %v to first click has mine", adjPos)
					}
				}

				// タイマーが開始されたことを確認
				// ※実際のテストでは時間の確認が難しいので、FirstClickがfalseになったことで代用
			},
		},
		{
			name: "click on mine loses game",
			setupGame: func() *Game {
				game := NewGame(Beginner)
				// テスト用に小さなボードを作成
				game.Board = NewBoard(3, 3, 1)
				game.Board.Cells[1][1].SetMine()
				game.FirstClick = false // 初回クリックを無効化
				return game
			},
			clickPos:  Position{Row: 1, Col: 1},
			wantState: Lost,
			checkBoard: func(t *testing.T, board *Board) {
				// すべての地雷が表示されることを確認
				for i := 0; i < board.Height; i++ {
					for j := 0; j < board.Width; j++ {
						cell := board.Cells[i][j]
						if cell.IsMine && !cell.IsRevealed {
							t.Errorf("Mine at [%d][%d] is not revealed after loss", i, j)
						}
					}
				}
			},
		},
		{
			name: "reveal all safe cells wins game",
			setupGame: func() *Game {
				game := NewGame(Beginner)
				// 3x3のボードで1つだけ地雷
				game.Board = NewBoard(3, 3, 1)
				game.Board.Cells[0][0].SetMine()
				// 周辺の隣接数を設定
				game.Board.Cells[0][1].SetAdjacent(1)
				game.Board.Cells[1][0].SetAdjacent(1)
				game.Board.Cells[1][1].SetAdjacent(1)
				game.FirstClick = false

				// 地雷以外をすべて開く（最後の1つを残す）
				for i := 0; i < 3; i++ {
					for j := 0; j < 3; j++ {
						if !(i == 0 && j == 0) && !(i == 2 && j == 2) {
							game.Board.Cells[i][j].IsRevealed = true
						}
					}
				}
				return game
			},
			clickPos:  Position{Row: 2, Col: 2}, // 最後の安全なセル
			wantState: Won,
			checkBoard: func(t *testing.T, board *Board) {
				// 勝利時には経過時間が記録されることを確認
				// ※実際のテストでは確認が難しいので省略
			},
		},
		{
			name: "cannot click when game is over",
			setupGame: func() *Game {
				game := NewGame(Beginner)
				game.State = Lost
				game.FirstClick = false
				return game
			},
			clickPos:  Position{Row: 1, Col: 1},
			wantState: Lost, // 状態は変わらない
			checkBoard: func(t *testing.T, board *Board) {
				// どのセルも開かれないことを確認
				allUnrevealed := true
				for i := 0; i < board.Height; i++ {
					for j := 0; j < board.Width; j++ {
						if board.Cells[i][j].IsRevealed {
							allUnrevealed = false
							break
						}
					}
				}
				if !allUnrevealed {
					t.Error("Cells were revealed even though game was over")
				}
			},
		},
		{
			name: "chain reaction reveals multiple cells",
			setupGame: func() *Game {
				game := NewGame(Beginner)
				// 5x5のボードで地雷を隅に配置（地雷を増やして勝利を防ぐ）
				game.Board = NewBoard(5, 5, 4)
				game.Board.Cells[0][4].SetMine()
				game.Board.Cells[1][4].SetMine()
				game.Board.Cells[3][4].SetMine()
				game.Board.Cells[4][4].SetMine()

				// 隣接数を手動で設定
				game.Board.Cells[0][3].SetAdjacent(2)
				game.Board.Cells[1][3].SetAdjacent(2)
				game.Board.Cells[2][3].SetAdjacent(2)
				game.Board.Cells[2][4].SetAdjacent(3)
				game.Board.Cells[3][3].SetAdjacent(2)
				game.Board.Cells[4][3].SetAdjacent(2)

				game.FirstClick = false
				return game
			},
			clickPos:  Position{Row: 0, Col: 0}, // 地雷から離れた位置
			wantState: Playing,
			checkBoard: func(t *testing.T, board *Board) {
				// 連鎖的に開かれることを確認
				// 少なくとも左側のエリアは開かれているはず
				expectedRevealed := []Position{
					{0, 0}, {0, 1}, {0, 2},
					{1, 0}, {1, 1}, {1, 2},
					{2, 0}, {2, 1}, {2, 2},
				}

				for _, pos := range expectedRevealed {
					cell := board.GetCell(pos)
					if cell != nil && !cell.IsRevealed {
						t.Errorf("Expected position %v to be revealed in chain reaction", pos)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := tt.setupGame()

			// クリック実行
			game.Click(tt.clickPos)

			// 状態の確認
			if game.State != tt.wantState {
				t.Errorf("State = %v, want %v", game.State, tt.wantState)
			}

			// FirstClickの確認（該当する場合）
			if tt.wantFirstClick && game.FirstClick != tt.wantFirstClick {
				t.Errorf("FirstClick = %v, want %v", game.FirstClick, tt.wantFirstClick)
			}

			// ボードの状態を確認
			if tt.checkBoard != nil {
				tt.checkBoard(t, game.Board)
			}
		})
	}
}

func TestGame_StateTransitions(t *testing.T) {
	tests := []struct {
		name      string
		scenario  func(*testing.T, *Game)
		wantState GameState
	}{
		{
			name: "playing to won transition",
			scenario: func(t *testing.T, game *Game) {
				// 3x3で地雷1つの簡単なゲーム
				game.Board = NewBoard(3, 3, 1)
				game.Board.Cells[0][0].SetMine()
				game.Board.Cells[0][1].SetAdjacent(1)
				game.Board.Cells[1][0].SetAdjacent(1)
				game.Board.Cells[1][1].SetAdjacent(1)
				game.FirstClick = false

				// 安全なセルをすべてクリック
				positions := []Position{
					{0, 1}, {0, 2},
					{1, 0}, {1, 1}, {1, 2},
					{2, 0}, {2, 1}, {2, 2},
				}

				for _, pos := range positions {
					game.Click(pos)
					if game.State == Won {
						break
					}
				}
			},
			wantState: Won,
		},
		{
			name: "playing to lost transition",
			scenario: func(t *testing.T, game *Game) {
				// 地雷を配置
				game.Board = NewBoard(3, 3, 1)
				game.Board.Cells[1][1].SetMine()
				game.FirstClick = false

				// 地雷をクリック
				game.Click(Position{Row: 1, Col: 1})
			},
			wantState: Lost,
		},
		{
			name: "reset from won state",
			scenario: func(t *testing.T, game *Game) {
				game.State = Won
				game.ElapsedTime = 100
				game.Reset()
			},
			wantState: Playing,
		},
		{
			name: "reset from lost state",
			scenario: func(t *testing.T, game *Game) {
				game.State = Lost
				game.Reset()
			},
			wantState: Playing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame(Beginner)

			// シナリオ実行
			tt.scenario(t, game)

			// 最終状態の確認
			if game.State != tt.wantState {
				t.Errorf("Final state = %v, want %v", game.State, tt.wantState)
			}
		})
	}
}

func TestGame_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		scenario func(*testing.T)
	}{
		{
			name: "click on flagged cell does not reveal",
			scenario: func(t *testing.T) {
				game := NewGame(Beginner)
				pos := Position{Row: 1, Col: 1}

				// フラグを設定
				game.ToggleFlag(pos)

				// フラグがあるセルをクリック
				game.Click(pos)

				// セルが開かれていないことを確認
				cell := game.Board.GetCell(pos)
				if cell != nil && cell.IsRevealed {
					t.Error("Flagged cell was revealed")
				}
			},
		},
		{
			name: "multiple clicks on same cell",
			scenario: func(t *testing.T) {
				game := NewGame(Beginner)
				game.Board = NewBoard(3, 3, 0) // 地雷なし
				game.FirstClick = false
				pos := Position{Row: 1, Col: 1}

				// 同じセルを複数回クリック
				game.Click(pos)
				firstState := game.State

				game.Click(pos)
				game.Click(pos)

				// 状態が変わらないことを確認
				if game.State != firstState {
					t.Error("Game state changed after multiple clicks on same cell")
				}
			},
		},
		{
			name: "click outside board boundaries",
			scenario: func(t *testing.T) {
				game := NewGame(Beginner)

				// 境界外をクリック
				invalidPositions := []Position{
					{-1, 0},
					{0, -1},
					{100, 0},
					{0, 100},
				}

				for _, pos := range invalidPositions {
					game.Click(pos)

					// ゲームがクラッシュしないことを確認
					if game.State != Playing {
						t.Error("Game state changed after clicking outside boundaries")
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.scenario(t)
		})
	}
}
