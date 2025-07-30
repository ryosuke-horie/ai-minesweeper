package testutil

import (
	"fmt"
	"testing"

	"github.com/r-horie/ai-minesweeper/game"
)

// AssertCellState セルの状態を検証
func AssertCellState(t *testing.T, cell *game.Cell, isMine, isRevealed, isFlagged bool, adjacent int) {
	t.Helper()
	
	if cell.IsMine != isMine {
		t.Errorf("Cell.IsMine = %v, want %v", cell.IsMine, isMine)
	}
	if cell.IsRevealed != isRevealed {
		t.Errorf("Cell.IsRevealed = %v, want %v", cell.IsRevealed, isRevealed)
	}
	if cell.IsFlagged != isFlagged {
		t.Errorf("Cell.IsFlagged = %v, want %v", cell.IsFlagged, isFlagged)
	}
	if cell.Adjacent != adjacent {
		t.Errorf("Cell.Adjacent = %d, want %d", cell.Adjacent, adjacent)
	}
}

// AssertBoardState ボードの状態を検証
func AssertBoardState(t *testing.T, board *game.Board, mines, revealed, flagged int) {
	t.Helper()
	
	actualMines := CountMines(board)
	actualRevealed := CountRevealed(board)
	actualFlagged := CountFlagged(board)
	
	if actualMines != mines {
		t.Errorf("Mine count = %d, want %d", actualMines, mines)
	}
	if actualRevealed != revealed {
		t.Errorf("Revealed count = %d, want %d", actualRevealed, revealed)
	}
	if actualFlagged != flagged {
		t.Errorf("Flagged count = %d, want %d", actualFlagged, flagged)
	}
}

// AssertGameState ゲームの状態を検証
func AssertGameState(t *testing.T, game *game.Game, state game.GameState) {
	t.Helper()
	
	if game.State != state {
		t.Errorf("Game.State = %v, want %v", game.State, state)
	}
}

// AssertPositionsEqual ポジションのスライスが等しいか検証（順序は問わない）
func AssertPositionsEqual(t *testing.T, actual, expected []game.Position) {
	t.Helper()
	
	if len(actual) != len(expected) {
		t.Errorf("Position count = %d, want %d", len(actual), len(expected))
		return
	}
	
	// マップを使って順序に依存しない比較
	expectedMap := make(map[string]bool)
	for _, pos := range expected {
		key := fmt.Sprintf("%d,%d", pos.Row, pos.Col)
		expectedMap[key] = true
	}
	
	for _, pos := range actual {
		key := fmt.Sprintf("%d,%d", pos.Row, pos.Col)
		if !expectedMap[key] {
			t.Errorf("Unexpected position: %v", pos)
		}
		delete(expectedMap, key)
	}
	
	for key := range expectedMap {
		t.Errorf("Missing expected position: %s", key)
	}
}

// CountMines ボード上の地雷数をカウント
func CountMines(board *game.Board) int {
	count := 0
	for i := 0; i < board.Height; i++ {
		for j := 0; j < board.Width; j++ {
			if board.Cells[i][j].IsMine {
				count++
			}
		}
	}
	return count
}

// CountRevealed 開かれたセルの数を数える
func CountRevealed(board *game.Board) int {
	count := 0
	for i := 0; i < board.Height; i++ {
		for j := 0; j < board.Width; j++ {
			if board.Cells[i][j].IsRevealed {
				count++
			}
		}
	}
	return count
}

// CountFlagged ボード上のフラグ数をカウント
func CountFlagged(board *game.Board) int {
	count := 0
	for i := 0; i < board.Height; i++ {
		for j := 0; j < board.Width; j++ {
			if board.Cells[i][j].IsFlagged {
				count++
			}
		}
	}
	return count
}

// CountUnrevealed ボード上の未開放セル数をカウント
func CountUnrevealed(board *game.Board) int {
	count := 0
	for i := 0; i < board.Height; i++ {
		for j := 0; j < board.Width; j++ {
			if !board.Cells[i][j].IsRevealed {
				count++
			}
		}
	}
	return count
}