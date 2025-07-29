# Issues

## #1 Goのlintツール導入とGitフックの設定

### 概要
コード品質を保つため、Goのlintツール（golangci-lint）を導入し、コミット時に自動でチェックが走るようにGitフックを設定する。

### 実装内容
1. golangci-lintの導入
   - `.golangci.yml`の設定ファイル作成
   - 主要なlinterの有効化（gofmt, govet, errcheck, ineffassign, unused等）

2. pre-commitフックの設定
   - `.githooks/pre-commit`スクリプトの作成
   - `go fmt`の自動実行
   - `golangci-lint run`の実行
   - テストの実行

3. CLAUDE.mdへのルール追加
   - 変更時のコミットルール明文化
   - lintエラーがある場合はコミットできないルール

### 期待される効果
- コード品質の向上
- 一貫性のあるコードスタイル
- バグの早期発見
- チーム開発時の品質担保