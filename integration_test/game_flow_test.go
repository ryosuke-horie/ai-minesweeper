package integration_test

import (
	"testing"

	"github.com/r-horie/ai-minesweeper/game"
	"github.com/r-horie/ai-minesweeper/solver"
	"github.com/r-horie/ai-minesweeper/testutil"
)

// TestGameFlowComplete ゲーム開始から終了までの完全なフローをテスト
func TestGameFlowComplete(t *testing.T) {
	tests := []struct {
		name       string
		difficulty game.Difficulty
		scenario   func(*testing.T, *game.Game)
	}{
		{
			name:       "beginner game flow",
			difficulty: game.Beginner,
			scenario:   playBeginnerGame,
		},
		{
			name:       "intermediate game flow",
			difficulty: game.Intermediate,
			scenario:   playIntermediateGame,
		},
		{
			name:       "expert game flow",
			difficulty: game.Expert,
			scenario:   playExpertGame,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := game.NewGame(tt.difficulty)
			tt.scenario(t, g)
		})
	}
}

// playBeginnerGame 初級ゲームのシナリオ
func playBeginnerGame(t *testing.T, g *game.Game) {
	// 初期状態の確認
	if g.State != game.Playing {
		t.Error("Game should start in Playing state")
	}
	if !g.FirstClick {
		t.Error("FirstClick should be true initially")
	}

	// 初回クリック（安全な位置）
	firstClickPos := game.Position{Row: 4, Col: 4}
	g.Click(firstClickPos)

	// 初回クリック後の確認
	if g.FirstClick {
		t.Error("FirstClick should be false after first click")
	}
	if g.Board.GetCell(firstClickPos).IsMine {
		t.Error("First click should never hit a mine")
	}

	// AIソルバーを使用して推論
	s := solver.NewSolver(g.Board)
	result := s.Solve()

	// 安全なセルがある場合はクリック
	for _, pos := range result.SafeCells {
		if g.State == game.Playing {
			g.Click(pos)
		}
	}

	// 地雷セルにフラグを立てる
	for _, pos := range result.MineCells {
		g.ToggleFlag(pos)
	}

	// ゲーム状態の確認
	if g.State == game.Lost {
		// 負けた場合、すべての地雷が表示されているか確認
		mineCount := 0
		for i := 0; i < g.Board.Height; i++ {
			for j := 0; j < g.Board.Width; j++ {
				cell := g.Board.Cells[i][j]
				if cell.IsMine && cell.IsRevealed {
					mineCount++
				}
			}
		}
		if mineCount == 0 {
			t.Error("All mines should be revealed when game is lost")
		}
	}
}

// playIntermediateGame 中級ゲームのシナリオ
func playIntermediateGame(t *testing.T, g *game.Game) {
	// より大きな盤面での初回クリック
	firstClickPos := game.Position{Row: 8, Col: 8}
	g.Click(firstClickPos)

	// 複数回のソルバー実行
	maxIterations := 10
	for i := 0; i < maxIterations && g.State == game.Playing; i++ {
		s := solver.NewSolver(g.Board)
		result := s.Solve()

		if !result.CanProgress {
			// 推論できない場合は、未開放セルからランダムに選択（実際のゲームプレイをシミュレート）
			unrevealedPositions := g.Board.GetAllUnrevealedPositions()
			if len(unrevealedPositions) > 0 {
				// 最初の未開放セルをクリック（テストなので決定的に）
				for _, pos := range unrevealedPositions {
					if !g.Board.GetCell(pos).IsFlagged {
						g.Click(pos)
						break
					}
				}
			}
		} else {
			// 推論結果を適用
			for _, pos := range result.SafeCells {
				if g.State == game.Playing {
					g.Click(pos)
				}
			}
			for _, pos := range result.MineCells {
				g.ToggleFlag(pos)
			}
		}
	}
}

// playExpertGame 上級ゲームのシナリオ
func playExpertGame(t *testing.T, g *game.Game) {
	// 大きな盤面でのテスト
	if g.Board.Width != 30 || g.Board.Height != 16 {
		t.Error("Expert board should be 30x16")
	}
	if g.Board.Mines != 99 {
		t.Error("Expert board should have 99 mines")
	}

	// パフォーマンステスト: 初回クリックと初期展開
	firstClickPos := game.Position{Row: 8, Col: 15}
	g.Click(firstClickPos)

	// 大量のセルが開かれることを確認
	revealedCount := testutil.CountRevealed(g.Board)
	if revealedCount < 10 {
		t.Error("Expert first click should reveal many cells")
	}
}

// TestGameWithAISolver ゲームとAIソルバーの連携テスト
func TestGameWithAISolver(t *testing.T) {
	// カスタムボードで確実に解けるパターンを作成
	board := testutil.NewBoardBuilder(5, 5, 5).
		WithPattern([]string{
			"11100",
			"*2210",
			"*3*10",
			"*3220",
			"11*10",
		}).
		Build()

	g := testutil.NewGameBuilder().
		WithCustomBoard(board).
		WithoutFirstClick().
		Build()

	// AIソルバーで解く
	iterations := 0
	maxIterations := 20

	for g.State == game.Playing && iterations < maxIterations {
		iterations++
		
		s := solver.NewSolver(g.Board)
		result := s.Solve()

		if !result.CanProgress {
			// これ以上推論できない
			break
		}

		// 安全なセルをクリック
		for _, pos := range result.SafeCells {
			if g.State == game.Playing {
				oldRevealed := testutil.CountRevealed(g.Board)
				g.Click(pos)
				newRevealed := testutil.CountRevealed(g.Board)
				
				if newRevealed <= oldRevealed {
					t.Errorf("Click on safe cell should reveal at least one cell")
				}
			}
		}

		// 地雷にフラグ
		for _, pos := range result.MineCells {
			g.ToggleFlag(pos)
		}

		// 勝利条件のチェック
		if g.Board.CountUnrevealedSafeCells() == 0 && g.State == game.Playing {
			t.Error("Game should be won when all safe cells are revealed")
		}
	}

	if iterations >= maxIterations {
		t.Error("Solver took too many iterations")
	}

	// 最終状態の確認
	t.Logf("Game ended with state: %v, iterations: %d", g.State, iterations)
	t.Logf("Final board:\n%s", testutil.DisplayBoard(g.Board))
}

// TestDifficultyTransitions 難易度切り替えのテスト
func TestDifficultyTransitions(t *testing.T) {
	difficulties := []game.Difficulty{
		game.Beginner,
		game.Intermediate,
		game.Expert,
	}

	for i, fromDiff := range difficulties {
		for j, toDiff := range difficulties {
			if i == j {
				continue
			}

			t.Run(fromDiff.Name+"_to_"+toDiff.Name, func(t *testing.T) {
				// 最初の難易度でゲーム開始
				g := game.NewGame(fromDiff)
				
				// いくつかクリック
				g.Click(game.Position{Row: 0, Col: 0})
				
				// 難易度を変更（新しいゲームを作成）
				g = game.NewGame(toDiff)
				
				// 新しい難易度が適用されているか確認
				if g.Board.Width != toDiff.Width {
					t.Errorf("Width = %d, want %d", g.Board.Width, toDiff.Width)
				}
				if g.Board.Height != toDiff.Height {
					t.Errorf("Height = %d, want %d", g.Board.Height, toDiff.Height)
				}
				if g.Board.Mines != toDiff.Mines {
					t.Errorf("Mines = %d, want %d", g.Board.Mines, toDiff.Mines)
				}
				
				// ゲームがリセットされているか確認
				if g.State != game.Playing {
					t.Error("New game should be in Playing state")
				}
				if !g.FirstClick {
					t.Error("New game should have FirstClick = true")
				}
			})
		}
	}
}

// TestResetBehavior リセット機能の詳細なテスト
func TestResetBehavior(t *testing.T) {
	g := game.NewGame(game.Beginner)
	
	// ゲームを進める
	g.Click(game.Position{Row: 4, Col: 4})
	g.ToggleFlag(game.Position{Row: 0, Col: 0})
	g.ToggleFlag(game.Position{Row: 1, Col: 1})
	
	// 状態を記録
	oldBoard := g.Board
	flagCount := testutil.CountFlagged(g.Board)
	
	if flagCount != 2 {
		t.Error("Should have 2 flags before reset")
	}
	
	// リセット
	g.Reset()
	
	// リセット後の確認
	if g.Board == oldBoard {
		t.Error("Reset should create a new board instance")
	}
	
	if g.State != game.Playing {
		t.Error("Reset should set state to Playing")
	}
	
	if !g.FirstClick {
		t.Error("Reset should set FirstClick to true")
	}
	
	if testutil.CountFlagged(g.Board) != 0 {
		t.Error("Reset should clear all flags")
	}
	
	if testutil.CountRevealed(g.Board) != 0 {
		t.Error("Reset should clear all revealed cells")
	}
	
	if g.ElapsedTime != 0 {
		t.Error("Reset should clear elapsed time")
	}
}

// TestConcurrentGames 複数ゲームの同時実行テスト
func TestConcurrentGames(t *testing.T) {
	// 異なる難易度で複数のゲームを作成
	games := []*game.Game{
		game.NewGame(game.Beginner),
		game.NewGame(game.Intermediate),
		game.NewGame(game.Expert),
	}
	
	// 各ゲームで異なる操作を実行
	for i, g := range games {
		// 異なる位置をクリック
		clickPos := game.Position{Row: i, Col: i}
		if g.Board.IsValidPosition(clickPos) {
			g.Click(clickPos)
		}
		
		// 各ゲームが独立していることを確認
		for j, other := range games {
			if i != j {
				if g.Board == other.Board {
					t.Error("Games should have independent boards")
				}
			}
		}
	}
}