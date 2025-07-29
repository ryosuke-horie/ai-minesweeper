package testutil

import (
	"github.com/r-horie/ai-minesweeper/game"
)

// ScenarioType テストシナリオの種類
type ScenarioType int

const (
	// SimpleNumberPattern 単純な数字パターン
	SimpleNumberPattern ScenarioType = iota
	// CornerMinePattern コーナーの地雷パターン
	CornerMinePattern
	// ChainReactionPattern 連鎖反応パターン
	ChainReactionPattern
	// ComplexLogicPattern 複雑な論理パターン
	ComplexLogicPattern
	// AlmostCompletePattern ほぼ完成したゲーム
	AlmostCompletePattern
)

// CreateScenario 事前定義されたシナリオを作成
func CreateScenario(scenarioType ScenarioType) *game.Board {
	switch scenarioType {
	case SimpleNumberPattern:
		return createSimpleNumberPattern()
	case CornerMinePattern:
		return createCornerMinePattern()
	case ChainReactionPattern:
		return createChainReactionPattern()
	case ComplexLogicPattern:
		return createComplexLogicPattern()
	case AlmostCompletePattern:
		return createAlmostCompletePattern()
	default:
		return game.NewBoard(9, 9, 10)
	}
}

// createSimpleNumberPattern 1-2-1パターン
func createSimpleNumberPattern() *game.Board {
	pattern := []string{
		"...",
		"121",
		"?*?",
	}
	return NewBoardBuilder(3, 3, 2).
		WithPattern(pattern).
		Build()
}

// createCornerMinePattern コーナーに地雷
func createCornerMinePattern() *game.Board {
	pattern := []string{
		"*21.",
		"*31.",
		"221.",
		"....",
	}
	return NewBoardBuilder(4, 4, 3).
		WithPattern(pattern).
		Build()
}

// createChainReactionPattern 連鎖反応をテスト
func createChainReactionPattern() *game.Board {
	pattern := []string{
		".....",
		".....",
		"..221",
		"..2*1",
		"..221",
	}
	return NewBoardBuilder(5, 5, 1).
		WithPattern(pattern).
		Build()
}

// createComplexLogicPattern AIの推論が必要なパターン
func createComplexLogicPattern() *game.Board {
	pattern := []string{
		"?????",
		"?212?",
		"?*3*?",
		"?212?",
		"?????",
	}
	return NewBoardBuilder(5, 5, 4).
		WithPattern(pattern).
		Build()
}

// createAlmostCompletePattern ほぼクリア状態
func createAlmostCompletePattern() *game.Board {
	pattern := []string{
		"11111",
		"1*22*",
		"11*21",
		"112*1",
		"..111",
	}
	// 最後の1マスだけ未開放
	board := NewBoardBuilder(5, 5, 4).
		WithPattern(pattern).
		Build()
	board.Cells[4][0].IsRevealed = false // 左下を未開放に
	return board
}

// GameScenario ゲーム全体のシナリオ
type GameScenario struct {
	Name        string
	Description string
	Setup       func() *game.Game
	Actions     []GameAction
	Expected    GameExpectation
}

// GameAction ゲームに対するアクション
type GameAction struct {
	Type     ActionType
	Position game.Position
}

// ActionType アクションの種類
type ActionType int

const (
	Click ActionType = iota
	Flag
	Reset
)

// GameExpectation 期待される結果
type GameExpectation struct {
	State          game.GameState
	RevealedCount  int
	FlaggedCount   int
	RemainingMines int
}

// CreateGameScenarios よく使うゲームシナリオのセット
func CreateGameScenarios() []GameScenario {
	return []GameScenario{
		{
			Name:        "First click win",
			Description: "初回クリックで大量に開く",
			Setup: func() *game.Game {
				return NewGameBuilder().
					WithDifficulty(game.Beginner).
					Build()
			},
			Actions: []GameAction{
				{Type: Click, Position: game.Position{Row: 0, Col: 0}},
			},
			Expected: GameExpectation{
				State:         game.Playing,
				RevealedCount: -1, // 不定
			},
		},
		{
			Name:        "Quick loss",
			Description: "地雷を踏んで即負け",
			Setup: func() *game.Game {
				board := NewBoardBuilder(3, 3, 1).
					WithMineAt(1, 1).
					Build()
				return NewGameBuilder().
					WithCustomBoard(board).
					WithoutFirstClick().
					Build()
			},
			Actions: []GameAction{
				{Type: Click, Position: game.Position{Row: 1, Col: 1}},
			},
			Expected: GameExpectation{
				State: game.Lost,
			},
		},
		{
			Name:        "Flag and win",
			Description: "フラグを立てて勝利",
			Setup: func() *game.Game {
				board := createAlmostCompletePattern()
				return NewGameBuilder().
					WithCustomBoard(board).
					WithoutFirstClick().
					Build()
			},
			Actions: []GameAction{
				{Type: Click, Position: game.Position{Row: 4, Col: 0}},
			},
			Expected: GameExpectation{
				State: game.Won,
			},
		},
	}
}