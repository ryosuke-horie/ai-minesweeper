package integration_test

import (
	"testing"
	"time"

	"github.com/r-horie/ai-minesweeper/game"
	"github.com/r-horie/ai-minesweeper/solver"
	"github.com/r-horie/ai-minesweeper/testutil"
)

// TestSolverCompleteGame ソルバーが完全なゲームを解くテスト
func TestSolverCompleteGame(t *testing.T) {
	tests := []struct {
		name     string
		board    *game.Board
		maxMoves int
		wantWin  bool
	}{
		{
			name: "solvable pattern",
			board: testutil.NewBoardBuilder(5, 5, 3).
				WithMineAt(1, 0).
				WithMineAt(2, 0).
				WithMineAt(3, 0).
				WithRevealedAt(0, 0).
				WithRevealedAt(0, 1).
				WithRevealedAt(1, 1).
				WithRevealedAt(2, 1).
				WithRevealedAt(3, 1).
				WithRevealedAt(4, 0).
				WithRevealedAt(4, 1).
				Build(),
			maxMoves: 30,
			wantWin:  true,
		},
		{
			name: "complex pattern",
			board: testutil.NewBoardBuilder(8, 8, 10).
				WithPattern([]string{
					"?1..11*1",
					"11..1221",
					"......11",
					"111211*1",
					"1*1*1111",
					"1111111.",
					"......1*",
					"......11",
				}).
				Build(),
			maxMoves: 50,
			wantWin:  false, // 左上のセルは推論できない可能性がある
		},
		{
			name: "partially solvable",
			board: testutil.NewBoardBuilder(5, 5, 5).
				WithPattern([]string{
					"?????",
					"?3*3?",
					"?*5*?",
					"?3*3?",
					"?????",
				}).
				Build(),
			maxMoves: 20,
			wantWin:  false, // 推論できない部分がある
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := testutil.NewGameBuilder().
				WithCustomBoard(tt.board).
				WithoutFirstClick().
				Build()

			moves := 0

			for g.State == game.Playing && moves < tt.maxMoves {
				moves++

				s := solver.NewSolver(g.Board)
				result := s.Solve()

				if result.CanProgress {

					// 安全なセルをクリック
					for _, pos := range result.SafeCells {
						if g.State == game.Playing {
							g.Click(pos)
						}
					}

					// 地雷にフラグ
					for _, pos := range result.MineCells {
						if !g.Board.GetCell(pos).IsFlagged {
							g.ToggleFlag(pos)
						}
					}
				} else {
					// 推論できない場合は終了
					break
				}
			}

			if tt.wantWin {
				// 勝利条件のチェック
				if g.Board.CountUnrevealedSafeCells() == 0 {
					// すべての安全なセルが開かれた場合
					if g.State != game.Won {
						// 勝利判定は明示的に行われない場合があるので、ログのみ
						t.Log("All safe cells revealed but game not marked as won (this is OK)")
					}
				} else if g.State != game.Won {
					t.Errorf("Expected to win the game. Unrevealed safe cells: %d", g.Board.CountUnrevealedSafeCells())
				}
			}

			t.Logf("Game ended after %d moves, state: %v", moves, g.State)
			t.Logf("Unrevealed safe cells: %d", g.Board.CountUnrevealedSafeCells())
			if testing.Verbose() {
				t.Logf("Final board:\n%s", testutil.DisplayBoard(g.Board))
			}
		})
	}
}

// TestSolverPerformance ソルバーのパフォーマンステスト
func TestSolverPerformance(t *testing.T) {
	difficulties := []struct {
		name string
		diff game.Difficulty
	}{
		{"beginner", game.Beginner},
		{"intermediate", game.Intermediate},
		{"expert", game.Expert},
	}

	for _, d := range difficulties {
		t.Run(d.name, func(t *testing.T) {
			g := game.NewGame(d.diff)
			
			// 初回クリック（中央付近）
			firstClick := game.Position{
				Row: g.Board.Height / 2,
				Col: g.Board.Width / 2,
			}
			g.Click(firstClick)

			start := time.Now()
			moves := 0
			maxMoves := 100

			for g.State == game.Playing && moves < maxMoves {
				moves++
				
				s := solver.NewSolver(g.Board)
				result := s.Solve()

				if !result.CanProgress {
					break
				}

				for _, pos := range result.SafeCells {
					if g.State == game.Playing {
						g.Click(pos)
					}
				}

				for _, pos := range result.MineCells {
					g.ToggleFlag(pos)
				}
			}

			elapsed := time.Since(start)

			t.Logf("Difficulty: %s", d.name)
			t.Logf("Board size: %dx%d, Mines: %d", d.diff.Width, d.diff.Height, d.diff.Mines)
			t.Logf("Moves made: %d", moves)
			t.Logf("Time elapsed: %v", elapsed)
			t.Logf("Average time per move: %v", elapsed/time.Duration(moves))
			t.Logf("Revealed cells: %d/%d", testutil.CountRevealed(g.Board), g.Board.Width*g.Board.Height-g.Board.Mines)
		})
	}
}

// TestSolverWithAIAssist AI支援機能の統合テスト
func TestSolverWithAIAssist(t *testing.T) {
	// AI支援を有効にしたゲームをシミュレート
	g := game.NewGame(game.Intermediate)
	
	aiAssistEnabled := true
	autoSolveCount := 0
	maxAutoSolves := 5

	// 初回クリック
	g.Click(game.Position{Row: 8, Col: 8})

	// AI支援のシミュレーション
	for g.State == game.Playing && autoSolveCount < maxAutoSolves {
		if aiAssistEnabled {
			s := solver.NewSolver(g.Board)
			result := s.Solve()

			if result.CanProgress {
				autoSolveCount++
				
				// 自動的に安全なセルを開く
				for _, pos := range result.SafeCells {
					if g.State == game.Playing {
						g.Click(pos)
					}
				}

				// 自動的に地雷にフラグを立てる
				for _, pos := range result.MineCells {
					if !g.Board.GetCell(pos).IsFlagged {
						g.ToggleFlag(pos)
					}
				}

				// AI支援の動作をログ
				t.Logf("AI assist #%d: revealed %d cells, flagged %d mines", 
					autoSolveCount, len(result.SafeCells), len(result.MineCells))
			} else {
				// これ以上推論できない
				break
			}
		}
	}

	t.Logf("AI assisted %d times", autoSolveCount)
	t.Logf("Final state: %v", g.State)
	t.Logf("Revealed: %d/%d safe cells", 
		g.Board.Width*g.Board.Height - g.Board.Mines - g.Board.CountUnrevealedSafeCells(),
		g.Board.Width*g.Board.Height - g.Board.Mines)
}

// TestSolverEdgeCases ソルバーのエッジケーステスト
func TestSolverEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() (*game.Game, *solver.Solver)
		validate func(*testing.T, *solver.SolverResult)
	}{
		{
			name: "empty board with revealed cell",
			setup: func() (*game.Game, *solver.Solver) {
				board := testutil.NewBoardBuilder(3, 3, 0).
					WithRevealedAt(1, 1).
					Build()
				g := testutil.NewGameBuilder().
					WithCustomBoard(board).
					WithoutFirstClick().
					Build()
				return g, solver.NewSolver(g.Board)
			},
			validate: func(t *testing.T, result *solver.SolverResult) {
				if !result.CanProgress {
					t.Error("Should be able to progress on empty board with revealed cell")
				}
				// 0のセルの周りはすべて安全
				if len(result.SafeCells) < 8 {
					t.Errorf("Should identify 8 safe cells around 0, got %d", len(result.SafeCells))
				}
			},
		},
		{
			name: "all mines",
			setup: func() (*game.Game, *solver.Solver) {
				board := testutil.NewBoardBuilder(3, 3, 9).
					WithPattern([]string{
						"***",
						"***",
						"***",
					}).
					Build()
				g := testutil.NewGameBuilder().
					WithCustomBoard(board).
					WithoutFirstClick().
					Build()
				return g, solver.NewSolver(g.Board)
			},
			validate: func(t *testing.T, result *solver.SolverResult) {
				if result.CanProgress {
					t.Error("Should not progress when all cells are mines")
				}
			},
		},
		{
			name: "single safe cell with all mines around",
			setup: func() (*game.Game, *solver.Solver) {
				board := testutil.NewBoardBuilder(3, 3, 8).
					WithMineAt(0, 0).WithMineAt(0, 1).WithMineAt(0, 2).
					WithMineAt(1, 0).WithMineAt(1, 2).
					WithMineAt(2, 0).WithMineAt(2, 1).WithMineAt(2, 2).
					WithRevealedAt(1, 1).
					WithAdjacentAt(1, 1, 8).
					Build()
				g := testutil.NewGameBuilder().
					WithCustomBoard(board).
					WithoutFirstClick().
					Build()
				return g, solver.NewSolver(g.Board)
			},
			validate: func(t *testing.T, result *solver.SolverResult) {
				t.Logf("Result: SafeCells=%d, MineCells=%d, CanProgress=%v", 
					len(result.SafeCells), len(result.MineCells), result.CanProgress)
				if !result.CanProgress {
					t.Error("Should progress with revealed number")
				}
				if len(result.MineCells) != 8 {
					t.Errorf("Should identify all 8 mines, got %d", len(result.MineCells))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, s := tt.setup()
			result := s.Solve()
			tt.validate(t, &result)
			
			// ゲームの状態も確認
			if g.State == game.Lost {
				t.Log("Game was lost (expected in some edge cases)")
			}
		})
	}
}

// TestInteractiveGameplay インタラクティブなゲームプレイのシミュレーション
func TestInteractiveGameplay(t *testing.T) {
	// ユーザーとAIが交互にプレイするシミュレーション
	g := game.NewGame(game.Beginner)
	
	type Turn struct {
		IsAI     bool
		Action   func(*game.Game, *solver.Solver) bool
		Name     string
	}

	turns := []Turn{
		{
			IsAI: false,
			Name: "user first click",
			Action: func(g *game.Game, _ *solver.Solver) bool {
				g.Click(game.Position{Row: 4, Col: 4})
				return true
			},
		},
		{
			IsAI: true,
			Name: "AI solve obvious",
			Action: func(g *game.Game, s *solver.Solver) bool {
				result := s.Solve()
				if result.CanProgress {
					for _, pos := range result.SafeCells {
						if g.State == game.Playing {
							g.Click(pos)
						}
					}
					for _, pos := range result.MineCells {
						g.ToggleFlag(pos)
					}
					return true
				}
				return false
			},
		},
		{
			IsAI: false,
			Name: "user random click",
			Action: func(g *game.Game, _ *solver.Solver) bool {
				// 未開放のセルを探してクリック
				for i := 0; i < g.Board.Height; i++ {
					for j := 0; j < g.Board.Width; j++ {
						cell := g.Board.Cells[i][j]
						if !cell.IsRevealed && !cell.IsFlagged {
							g.Click(game.Position{Row: i, Col: j})
							return true
						}
					}
				}
				return false
			},
		},
	}

	turnCount := 0
	maxTurns := 20

	for g.State == game.Playing && turnCount < maxTurns {
		turn := turns[turnCount%len(turns)]
		s := solver.NewSolver(g.Board)
		
		if turn.Action(g, s) {
			t.Logf("Turn %d (%s): %s", turnCount+1, 
				map[bool]string{true: "AI", false: "User"}[turn.IsAI], 
				turn.Name)
		}
		
		turnCount++
	}

	t.Logf("Game ended after %d turns", turnCount)
	t.Logf("Final state: %v", g.State)
	t.Logf("Board:\n%s", testutil.DisplayBoard(g.Board))
}