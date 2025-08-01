# テスト実装計画

## 概要
ADR-001で定義したテスト戦略に基づき、段階的にテストを実装していく。

## フェーズ分け

### Phase 1: ゲームコアのユニットテスト（基礎）
**目的**: ゲームの基本機能が正しく動作することを保証

#### タスク
- [ ] `game/cell_test.go`の実装
  - [ ] NewCell()のテスト
  - [ ] Reveal()のテスト
  - [ ] ToggleFlag()のテスト
  - [ ] SetMine()とSetAdjacent()のテスト

- [ ] `game/board_test.go`の実装（基本機能）
  - [ ] NewBoard()のテスト
  - [ ] GetCell()とIsValidPosition()のテスト
  - [ ] GetAdjacentPositions()のテスト

- [ ] `game/game_test.go`の実装（基本機能）
  - [ ] NewGame()のテスト
  - [ ] Reset()のテスト
  - [ ] GetRemainingMines()のテスト

### Phase 2: ゲームコアのユニットテスト（高度な機能）
**目的**: 複雑なゲームロジックの正確性を保証

#### タスク
- [ ] `game/board_test.go`の拡張
  - [ ] Initialize()のテスト（地雷配置の検証）
  - [ ] RevealCell()のテスト（連鎖開放を含む）
  - [ ] CountUnrevealedSafeCells()のテスト

- [ ] `game/game_test.go`の拡張
  - [ ] Click()のテスト（初回クリック保証を含む）
  - [ ] ゲーム状態遷移のテスト（Playing→Won/Lost）
  - [ ] ToggleFlag()のゲームレベルでのテスト

### Phase 3: AIソルバーのユニットテスト
**目的**: AIの推論ロジックが正しく動作することを保証

#### タスク
- [ ] `solver/solver_test.go`の実装
  - [ ] NewSolver()のテスト
  - [ ] findDefiniteMines()のテスト（確実な地雷検出）
  - [ ] findDefiniteSafeCells()のテスト（確実な安全マス検出）
  - [ ] Solve()の統合テスト
  - [ ] エッジケースのテスト（角、端、パターン）

### Phase 4: テストヘルパーとユーティリティ
**目的**: テストの記述を容易にし、保守性を向上

#### タスク
- [ ] テストヘルパーの作成
  - [ ] 盤面生成ヘルパー（特定のパターンを作成）
  - [ ] 盤面検証ヘルパー（状態の確認）
  - [ ] ゲーム状態のダンプ機能

- [ ] テストデータの整備
  - [ ] 典型的なゲームパターンの定義
  - [ ] エッジケースのパターン定義

### Phase 5: 統合テスト
**目的**: コンポーネント間の連携が正しく動作することを保証

#### タスク
- [ ] `integration_test.go`の実装
  - [ ] ゲーム開始から終了までのシナリオテスト
  - [ ] AIソルバーとゲームロジックの連携テスト
  - [ ] 難易度別のゲームプレイテスト
  - [ ] エラーケースの統合テスト

### Phase 6: CI/CD統合とカバレッジ
**目的**: 継続的な品質保証体制の確立

#### タスク
- [ ] GitHub Actionsの設定
  - [ ] テスト自動実行の設定
  - [ ] カバレッジレポートの生成
  - [ ] PR時の必須チェック設定

- [ ] カバレッジ目標の達成
  - [ ] 各パッケージ80%以上のカバレッジ
  - [ ] カバレッジバッジの追加

## 成功基準

### 定量的な基準
- ユニットテストカバレッジ: 80%以上
- 全テストの実行時間: 10秒以内
- テストの成功率: 100%（フレーキーなテストなし）

### 定性的な基準
- 新機能追加時にテストが書きやすい
- バグ発生時にテストで再現可能
- リファクタリング時の安心感

## 見積もり

### 工数見積もり
- Phase 1: 2-3時間
- Phase 2: 3-4時間
- Phase 3: 4-5時間
- Phase 4: 2-3時間
- Phase 5: 3-4時間
- Phase 6: 2-3時間

**合計: 16-22時間**

## リスクと対策

### リスク
1. **ランダム性のテスト**
   - 対策: シード値固定、モック使用

2. **非同期処理のテスト**
   - 対策: 適切なwait処理、タイムアウト設定

3. **TUI部分のテスト困難性**
   - 対策: ロジックとUIの分離を徹底

## 次のステップ
1. この計画のレビューと承認
2. 各フェーズのIssue作成
3. Phase 1から順次実装開始