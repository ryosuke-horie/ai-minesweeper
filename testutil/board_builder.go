package testutil

import (
	"github.com/r-horie/ai-minesweeper/game"
)

// BoardBuilder はテスト用のボードを構築するためのビルダー.
type BoardBuilder struct {
	board *game.Board
}

// NewBoardBuilder は新しいBoardBuilderを作成.
func NewBoardBuilder(width, height, mines int) *BoardBuilder {
	return &BoardBuilder{
		board: game.NewBoard(width, height, mines),
	}
}

// WithMineAt 指定位置に地雷を配置.
func (b *BoardBuilder) WithMineAt(row, col int) *BoardBuilder {
	if b.isValidPosition(row, col) {
		b.board.Cells[row][col].SetMine()
	}
	return b
}

// WithRevealedAt 指定位置のセルを開く.
func (b *BoardBuilder) WithRevealedAt(row, col int) *BoardBuilder {
	if b.isValidPosition(row, col) {
		b.board.Cells[row][col].IsRevealed = true
	}
	return b
}

// WithFlagAt 指定位置にフラグを配置.
func (b *BoardBuilder) WithFlagAt(row, col int) *BoardBuilder {
	if b.isValidPosition(row, col) {
		b.board.Cells[row][col].IsFlagged = true
	}
	return b
}

// WithAdjacentAt 指定位置の隣接数を設定.
func (b *BoardBuilder) WithAdjacentAt(row, col int, adjacent int) *BoardBuilder {
	if b.isValidPosition(row, col) {
		b.board.Cells[row][col].SetAdjacent(adjacent)
	}
	return b
}

// WithPattern 文字列パターンからボードを構築.
// '.' = 空のセル
// '*' = 地雷
// '1'-'8' = 隣接数を持つ開かれたセル
// 'F' = フラグ
// '?' = 未開放.
func (b *BoardBuilder) WithPattern(pattern []string) *BoardBuilder {
	for i, row := range pattern {
		if i >= b.board.Height {
			break
		}
		for j, cell := range row {
			if j >= b.board.Width {
				break
			}
			switch cell {
			case '*':
				b.board.Cells[i][j].SetMine()
			case 'F':
				b.board.Cells[i][j].IsFlagged = true
			case '?':
				// 未開放（デフォルト）
			case '.':
				b.board.Cells[i][j].IsRevealed = true
			case '1', '2', '3', '4', '5', '6', '7', '8':
				b.board.Cells[i][j].IsRevealed = true
				b.board.Cells[i][j].SetAdjacent(int(cell - '0'))
			}
		}
	}
	return b
}

// Build 構築したボードを返す.
func (b *BoardBuilder) Build() *game.Board {
	// 隣接数を自動計算（地雷が配置されている場合）
	b.calculateAdjacents()
	return b.board
}

// isValidPosition 位置が有効かチェック.
func (b *BoardBuilder) isValidPosition(row, col int) bool {
	return row >= 0 && row < b.board.Height && col >= 0 && col < b.board.Width
}

// calculateAdjacents 隣接地雷数を計算.
func (b *BoardBuilder) calculateAdjacents() {
	for i := 0; i < b.board.Height; i++ {
		for j := 0; j < b.board.Width; j++ {
			if !b.board.Cells[i][j].IsMine && b.board.Cells[i][j].Adjacent == 0 {
				count := 0
				for _, pos := range b.board.GetAdjacentPositions(game.Position{Row: i, Col: j}) {
					if b.board.Cells[pos.Row][pos.Col].IsMine {
						count++
					}
				}
				// count が 0 でも設定する（空のセルも重要）
				b.board.Cells[i][j].SetAdjacent(count)
			}
		}
	}
}
