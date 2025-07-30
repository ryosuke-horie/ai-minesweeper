package testutil

import (
	"strings"
	"testing"

	"github.com/r-horie/ai-minesweeper/game"
)

func TestBoardBuilder(t *testing.T) {
	t.Run("basic building", func(t *testing.T) {
		board := NewBoardBuilder(3, 3, 1).
			WithMineAt(1, 1).
			WithRevealedAt(0, 0).
			WithFlagAt(2, 2).
			WithAdjacentAt(0, 1, 1).
			Build()

		if !board.Cells[1][1].IsMine {
			t.Error("Mine not set at [1,1]")
		}
		if !board.Cells[0][0].IsRevealed {
			t.Error("Cell not revealed at [0,0]")
		}
		if !board.Cells[2][2].IsFlagged {
			t.Error("Flag not set at [2,2]")
		}
		if board.Cells[0][1].Adjacent != 1 {
			t.Error("Adjacent not set correctly at [0,1]")
		}
	})

	t.Run("pattern building", func(t *testing.T) {
		pattern := []string{
			"12*",
			"2*2",
			"*21",
		}
		board := NewBoardBuilder(3, 3, 3).
			WithPattern(pattern).
			Build()

		// 地雷の確認
		expectedMines := []game.Position{
			{Row: 0, Col: 2}, {Row: 1, Col: 1}, {Row: 2, Col: 0},
		}
		for _, pos := range expectedMines {
			if !board.Cells[pos.Row][pos.Col].IsMine {
				t.Errorf("Expected mine at %v", pos)
			}
		}

		// 数字の確認
		if board.Cells[0][0].Adjacent != 1 || !board.Cells[0][0].IsRevealed {
			t.Error("Cell [0,0] should be revealed with adjacent=1")
		}
		if board.Cells[0][1].Adjacent != 2 || !board.Cells[0][1].IsRevealed {
			t.Error("Cell [0,1] should be revealed with adjacent=2")
		}
	})

	t.Run("auto adjacent calculation", func(t *testing.T) {
		board := NewBoardBuilder(3, 3, 4).
			WithMineAt(0, 0).
			WithMineAt(0, 2).
			WithMineAt(2, 0).
			WithMineAt(2, 2).
			Build()

		// 中央のセルは4つの地雷に囲まれている
		centerCell := board.Cells[1][1]
		if centerCell.Adjacent != 4 {
			t.Errorf("Center cell adjacent = %d, want 4", centerCell.Adjacent)
		}
	})
}

func TestGameBuilder(t *testing.T) {
	t.Run("default game", func(t *testing.T) {
		g := NewGameBuilder().Build()

		if g.State != game.Playing {
			t.Error("Default game should be in Playing state")
		}
		if !g.FirstClick {
			t.Error("Default game should have FirstClick=true")
		}
		if g.Difficulty.Name != "初級" {
			t.Error("Default game should be Beginner difficulty")
		}
	})

	t.Run("custom settings", func(t *testing.T) {
		customBoard := NewBoardBuilder(5, 5, 5).Build()
		g := NewGameBuilder().
			WithDifficulty(game.Expert).
			WithCustomBoard(customBoard).
			WithState(game.Won).
			WithoutFirstClick().
			WithElapsedTime(123).
			Build()

		if g.Board != customBoard {
			t.Error("Custom board not set")
		}
		if g.State != game.Won {
			t.Error("State not set to Won")
		}
		if g.FirstClick {
			t.Error("FirstClick should be false")
		}
		if g.ElapsedTime != 123 {
			t.Error("ElapsedTime not set correctly")
		}
	})
}

func TestAssertions(t *testing.T) {
	t.Run("AssertCellState", func(t *testing.T) {
		cell := game.NewCell()
		cell.SetMine()
		cell.Reveal()
		cell.SetAdjacent(3)

		// この呼び出しはパスするはず
		AssertCellState(t, cell, true, true, false, 3)
	})

	t.Run("AssertBoardState", func(t *testing.T) {
		board := NewBoardBuilder(3, 3, 2).
			WithMineAt(0, 0).
			WithMineAt(1, 1).
			WithRevealedAt(0, 1).
			WithRevealedAt(0, 2).
			WithFlagAt(2, 2).
			Build()

		AssertBoardState(t, board, 2, 2, 1)
	})

	t.Run("AssertPositionsEqual", func(t *testing.T) {
		actual := []game.Position{
			{Row: 1, Col: 2},
			{Row: 0, Col: 0},
			{Row: 2, Col: 1},
		}
		expected := []game.Position{
			{Row: 0, Col: 0},
			{Row: 2, Col: 1},
			{Row: 1, Col: 2},
		}

		// 順序は異なるが同じ要素
		AssertPositionsEqual(t, actual, expected)
	})
}

func TestDisplay(t *testing.T) {
	t.Run("DisplayBoard", func(t *testing.T) {
		board := NewBoardBuilder(3, 3, 1).
			WithPattern([]string{
				"12*",
				"2*F",
				"?2?",
			}).
			Build()

		display := DisplayBoard(board)

		// ヘッダーがあることを確認
		if !strings.Contains(display, "0 1 2") {
			t.Error("Column headers missing")
		}

		// 各セルの表示を確認
		// 地雷は開かれていないので'?'として表示される
		if !strings.Contains(display, "1 2 ?") {
			t.Error("First row not displayed correctly")
		}
		if !strings.Contains(display, "F") {
			t.Error("Flag not displayed")
		}
	})

	t.Run("DisplayBoardCompact", func(t *testing.T) {
		board := NewBoardBuilder(3, 3, 1).
			WithPattern([]string{
				"12*",
				"2*F",
				"?2?",
			}).
			Build()

		display := DisplayBoardCompact(board)
		// 地雷は開かれていないので'?'として表示される
		expected := "12?\n2?F\n?2?"

		if display != expected {
			t.Errorf("Compact display = %q, want %q", display, expected)
		}
	})

	t.Run("CompareBoardStates", func(t *testing.T) {
		board1 := NewBoardBuilder(2, 2, 0).
			WithRevealedAt(0, 0).
			Build()

		board2 := NewBoardBuilder(2, 2, 0).
			WithRevealedAt(0, 0).
			WithRevealedAt(1, 1).
			Build()

		diff := CompareBoardStates(board1, board2)

		if !strings.Contains(diff, "[1,1]") {
			t.Error("Difference at [1,1] not detected")
		}
	})
}

func TestScenarios(t *testing.T) {
	scenarios := []struct {
		name         string
		scenarioType ScenarioType
		checkFunc    func(*testing.T, *game.Board)
	}{
		{
			name:         "SimpleNumberPattern",
			scenarioType: SimpleNumberPattern,
			checkFunc: func(t *testing.T, board *game.Board) {
				// 1-2-1パターンの確認
				if board.Cells[1][0].Adjacent != 1 {
					t.Error("Expected 1 at [1,0]")
				}
				if board.Cells[1][1].Adjacent != 2 {
					t.Error("Expected 2 at [1,1]")
				}
			},
		},
		{
			name:         "CornerMinePattern",
			scenarioType: CornerMinePattern,
			checkFunc: func(t *testing.T, board *game.Board) {
				// コーナーに地雷があることを確認
				if !board.Cells[0][0].IsMine {
					t.Error("Expected mine at corner [0,0]")
				}
				if !board.Cells[1][0].IsMine {
					t.Error("Expected mine at [1,0]")
				}
			},
		},
		{
			name:         "ChainReactionPattern",
			scenarioType: ChainReactionPattern,
			checkFunc: func(t *testing.T, board *game.Board) {
				// 空のエリアがあることを確認
				emptyCount := 0
				for i := 0; i < board.Height; i++ {
					for j := 0; j < board.Width; j++ {
						if board.Cells[i][j].IsRevealed && board.Cells[i][j].Adjacent == 0 {
							emptyCount++
						}
					}
				}
				if emptyCount == 0 {
					t.Error("Expected empty cells for chain reaction")
				}
			},
		},
	}

	for _, tc := range scenarios {
		t.Run(tc.name, func(t *testing.T) {
			board := CreateScenario(tc.scenarioType)
			tc.checkFunc(t, board)
		})
	}
}

func TestGameScenarios(t *testing.T) {
	scenarios := CreateGameScenarios()

	// 最低限のシナリオがあることを確認
	if len(scenarios) < 3 {
		t.Error("Expected at least 3 game scenarios")
	}

	// 各シナリオが有効であることを確認
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			if scenario.Setup == nil {
				t.Error("Setup function is nil")
			}
			if len(scenario.Actions) == 0 {
				t.Error("No actions defined")
			}

			// セットアップが動作することを確認
			g := scenario.Setup()
			if g == nil {
				t.Error("Setup returned nil game")
			}
		})
	}
}
