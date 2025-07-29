package examples

import (
	"testing"

	"github.com/r-horie/ai-minesweeper/game"
	"github.com/r-horie/ai-minesweeper/testutil"
)

// TestGameWithHelpers ヘルパーを使用したテストの例
func TestGameWithHelpers(t *testing.T) {
	t.Run("using BoardBuilder", func(t *testing.T) {
		// 複雑なボードパターンを簡単に作成
		board := testutil.NewBoardBuilder(5, 5, 3).
			WithPattern([]string{
				"12*21",
				"2*321",
				"2221*",
				"1110?",
				"?10??",
			}).
			Build()

		// ゲームに適用
		g := testutil.NewGameBuilder().
			WithCustomBoard(board).
			WithoutFirstClick().
			Build()

		// クリックして連鎖反応をテスト
		g.Click(game.Position{Row: 3, Col: 2})

		// アサーションヘルパーを使用
		testutil.AssertGameState(t, g, game.Playing)
		
		// ボードの状態を視覚的に確認（デバッグ時に便利）
		if testing.Verbose() {
			t.Logf("Board after click:\n%s", testutil.DisplayBoard(board))
		}
	})

	t.Run("using scenarios", func(t *testing.T) {
		// 事前定義されたシナリオを使用
		board := testutil.CreateScenario(testutil.ComplexLogicPattern)
		
		g := testutil.NewGameBuilder().
			WithCustomBoard(board).
			WithoutFirstClick().
			Build()

		// AIソルバーのテストに適したパターンが自動生成される
		revealed := testutil.CountRevealed(board)
		if revealed == 0 {
			t.Error("Complex pattern should have some revealed cells")
		}
	})

	t.Run("game state transitions", func(t *testing.T) {
		// ほぼ完成したゲームのシナリオ
		board := testutil.CreateScenario(testutil.AlmostCompletePattern)
		
		g := testutil.NewGameBuilder().
			WithCustomBoard(board).
			WithoutFirstClick().
			Build()

		// 最後のセルをクリックして勝利
		g.Click(game.Position{Row: 4, Col: 0})
		
		// 複数の条件を一度にチェック
		testutil.AssertGameState(t, g, game.Won)
		testutil.AssertBoardState(t, board, 4, board.Width*board.Height-4, 0)
	})
}

// TestSolverWithHelpers ソルバーのテストでヘルパーを使用する例
func TestSolverWithHelpers(t *testing.T) {
	t.Run("definite mine detection", func(t *testing.T) {
		// 1-2-1パターンで確実な地雷を検出
		board := testutil.CreateScenario(testutil.SimpleNumberPattern)
		
		// ボードの状態を確認
		if testing.Verbose() {
			t.Logf("Test pattern:\n%s", testutil.DisplayBoardCompact(board))
		}
		
		// ここでソルバーをテスト（実際のソルバーテストで使用）
	})
}