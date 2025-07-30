package testutil

import (
	"fmt"
	"strings"

	"github.com/r-horie/ai-minesweeper/game"
)

// DisplayBoard ボードの状態を文字列で表示（デバッグ用）.
func DisplayBoard(board *game.Board) string {
	var sb strings.Builder

	// ヘッダー
	sb.WriteString("  ")
	for j := 0; j < board.Width; j++ {
		sb.WriteString(fmt.Sprintf("%d ", j))
	}
	sb.WriteString("\n")

	// ボード
	for i := 0; i < board.Height; i++ {
		sb.WriteString(fmt.Sprintf("%d ", i))
		for j := 0; j < board.Width; j++ {
			sb.WriteString(getCellChar(board.Cells[i][j]))
			sb.WriteString(" ")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// DisplayBoardCompact コンパクトなボード表示.
func DisplayBoardCompact(board *game.Board) string {
	var sb strings.Builder

	for i := 0; i < board.Height; i++ {
		for j := 0; j < board.Width; j++ {
			sb.WriteString(getCellChar(board.Cells[i][j]))
		}
		if i < board.Height-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// getCellChar セルを文字で表現.
func getCellChar(cell *game.Cell) string {
	if cell.IsFlagged {
		return "F"
	}
	if !cell.IsRevealed {
		return "?"
	}
	if cell.IsMine {
		return "*"
	}
	if cell.Adjacent == 0 {
		return "."
	}
	return fmt.Sprintf("%d", cell.Adjacent)
}

// CompareBoardStates 2つのボードの差分を表示.
func CompareBoardStates(board1, board2 *game.Board) string {
	if board1.Width != board2.Width || board1.Height != board2.Height {
		return "Boards have different dimensions"
	}

	var sb strings.Builder
	sb.WriteString("Differences (row, col):\n")

	differences := 0
	for i := 0; i < board1.Height; i++ {
		for j := 0; j < board1.Width; j++ {
			cell1 := board1.Cells[i][j]
			cell2 := board2.Cells[i][j]

			if !cellsEqual(cell1, cell2) {
				differences++
				sb.WriteString(fmt.Sprintf("[%d,%d]: %s -> %s\n",
					i, j, getCellChar(cell1), getCellChar(cell2)))
			}
		}
	}

	if differences == 0 {
		return "No differences found"
	}

	return sb.String()
}

// cellsEqual 2つのセルが等しいか判定.
func cellsEqual(c1, c2 *game.Cell) bool {
	return c1.IsMine == c2.IsMine &&
		c1.IsRevealed == c2.IsRevealed &&
		c1.IsFlagged == c2.IsFlagged &&
		c1.Adjacent == c2.Adjacent
}
