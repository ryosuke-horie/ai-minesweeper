# AI Minesweeper テスティングガイド

## 概要

このドキュメントは、AI Minesweeperプロジェクトのテスト戦略とベストプラクティスについて説明します。

## テスト構造

```
.
├── game/              # ゲームコアのユニットテスト
│   ├── board_test.go
│   ├── cell_test.go
│   └── game_test.go
├── solver/            # AIソルバーのユニットテスト
│   ├── solver_test.go
│   └── solver_bench_test.go
├── testutil/          # テストユーティリティ
│   ├── assertions.go
│   ├── board_builder.go
│   ├── display.go
│   ├── game_builder.go
│   └── scenarios.go
└── integration_test/  # 統合テスト
    ├── game_flow_test.go
    └── solver_integration_test.go
```

## テストの実行

### すべてのテストを実行
```bash
make test-all
```

### ユニットテストのみ
```bash
make test
```

### 統合テストのみ
```bash
make integration-test
```

### カバレッジレポート生成
```bash
make coverage
```

### ベンチマークテスト
```bash
make benchmark
```

## テストヘルパーの使用

### BoardBuilder

複雑なボード状態を簡単に作成：

```go
board := testutil.NewBoardBuilder(5, 5, 3).
    WithPattern([]string{
        "12*21",
        "2*321",
        "2221*",
        "1110?",
        "?10??",
    }).
    Build()
```

パターン記法：
- `*` - 地雷
- `0-8` - 数字（自動的に開かれる）
- `?` - 未開放セル
- `F` - フラグ
- `.` - 開かれた空セル

### GameBuilder

特定の状態のゲームを作成：

```go
game := testutil.NewGameBuilder().
    WithDifficulty(game.Expert).
    WithCustomBoard(board).
    WithState(game.Playing).
    WithoutFirstClick().
    Build()
```

### アサーション

一貫性のある検証：

```go
// セルの状態を検証
testutil.AssertCellState(t, cell, isMine, isRevealed, isFlagged, adjacent)

// ボードの統計を検証
testutil.AssertBoardState(t, board, mines, revealed, flagged)

// ゲーム状態を検証
testutil.AssertGameState(t, game, expectedState)
```

### デバッグ表示

テスト失敗時のデバッグに便利：

```go
if testing.Verbose() {
    t.Logf("Board state:\n%s", testutil.DisplayBoard(board))
}

// コンパクト表示
compact := testutil.DisplayBoardCompact(board)
```

## ベストプラクティス

### 1. テーブル駆動テスト

```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
    }{
        // テストケース
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テスト実装
        })
    }
}
```

### 2. ヘルパー関数の活用

```go
func TestComplexScenario(t *testing.T) {
    // シナリオヘルパーを使用
    board := testutil.CreateScenario(testutil.ComplexLogicPattern)
    
    // ゲームをセットアップ
    g := testutil.NewGameBuilder().
        WithCustomBoard(board).
        Build()
    
    // テスト実行
}
```

### 3. エラーメッセージの改善

```go
// 悪い例
if got != want {
    t.Error("Test failed")
}

// 良い例
if got != want {
    t.Errorf("Cell adjacent = %d, want %d", got, want)
}
```

### 4. 並行テストの活用

```go
func TestConcurrent(t *testing.T) {
    t.Parallel() // 他のテストと並行実行
    
    // テスト実装
}
```

## CI/CD統合

### GitHub Actions

プロジェクトには以下のCIジョブが設定されています：

1. **test** - ユニットテストとカバレッジ
2. **lint** - 静的解析
3. **build** - クロスプラットフォームビルド
4. **integration-test** - 統合テスト

### カバレッジ目標

- 全体: 70%以上
- 新規コード: 80%以上
- 重要パッケージ（game, solver）: 90%以上

## トラブルシューティング

### テストが遅い場合

```bash
# 短時間テストのみ実行
make test-short

# 特定パッケージのみテスト
go test -v ./game/...
```

### カバレッジが低い場合

```bash
# 関数ごとのカバレッジを確認
go tool cover -func=coverage.out

# HTMLレポートで詳細確認
make coverage-report
```

### 統合テストが失敗する場合

```bash
# 詳細ログ付きで実行
go test -v ./integration_test/... -run TestName
```

## 新しいテストの追加

1. 適切なディレクトリにテストファイルを作成
2. テストヘルパーを活用して簡潔に記述
3. カバレッジを確認
4. CIが通ることを確認

## 参考資料

- [Go Testing Best Practices](https://golang.org/doc/tutorial/add-a-test)
- [Table Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Test Fixtures in Go](https://medium.com/@benbjohnson/structuring-tests-in-go-46ddee7a25c)