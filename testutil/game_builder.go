package testutil

import (
	"github.com/r-horie/ai-minesweeper/game"
)

// GameBuilder はテスト用のゲームを構築するためのビルダー.
type GameBuilder struct {
	game *game.Game
}

// NewGameBuilder は新しいGameBuilderを作成.
func NewGameBuilder() *GameBuilder {
	return &GameBuilder{
		game: game.NewGame(game.Beginner),
	}
}

// WithDifficulty 難易度を設定.
func (g *GameBuilder) WithDifficulty(difficulty game.Difficulty) *GameBuilder {
	g.game = game.NewGame(difficulty)
	return g
}

// WithCustomBoard カスタムボードを設定.
func (g *GameBuilder) WithCustomBoard(board *game.Board) *GameBuilder {
	g.game.Board = board
	return g
}

// WithState ゲーム状態を設定.
func (g *GameBuilder) WithState(state game.GameState) *GameBuilder {
	g.game.State = state
	return g
}

// WithoutFirstClick 初回クリックを済ませた状態にする.
func (g *GameBuilder) WithoutFirstClick() *GameBuilder {
	g.game.FirstClick = false
	return g
}

// WithElapsedTime 経過時間を設定.
func (g *GameBuilder) WithElapsedTime(seconds int64) *GameBuilder {
	g.game.ElapsedTime = seconds
	return g
}

// Build 構築したゲームを返す.
func (g *GameBuilder) Build() *game.Game {
	return g.game
}
