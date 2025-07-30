.PHONY: help test coverage lint build clean run install-tools integration-test

# デフォルトターゲット
.DEFAULT_GOAL := help

# ヘルプメッセージ
help: ## ヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# セットアップ
setup: ## 開発環境のセットアップ
	@echo "Setting up git hooks..."
	git config core.hooksPath .githooks
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "Setup complete!"

# ビルド関連
build: ## アプリケーションをビルド
	go build -v -o bin/minesweeper ./main.go

build-all: ## 全プラットフォーム向けにビルド
	GOOS=darwin GOARCH=amd64 go build -o bin/minesweeper-darwin-amd64 ./main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/minesweeper-darwin-arm64 ./main.go
	GOOS=linux GOARCH=amd64 go build -o bin/minesweeper-linux-amd64 ./main.go
	GOOS=windows GOARCH=amd64 go build -o bin/minesweeper-windows-amd64.exe ./main.go

# テスト関連
test: ## ユニットテストを実行
	go test -v -race ./...

test-short: ## 短時間のテストのみ実行
	go test -v -short ./...

integration-test: ## 統合テストを実行
	go test -v ./integration_test/...

test-all: test integration-test ## すべてのテストを実行

coverage: ## カバレッジレポートを生成
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

coverage-report: coverage ## カバレッジレポートをブラウザで開く
	open coverage.html

benchmark: ## ベンチマークテストを実行
	go test -bench=. -benchmem ./...

# 静的解析
lint: ## golangci-lintを実行
	golangci-lint run --timeout=5m

fmt: ## コードフォーマット
	go fmt ./...
	goimports -w .

vet: ## go vetを実行
	go vet ./...

# 実行
run: ## アプリケーションを実行
	go run main.go

# クリーンアップ
clean: ## ビルド成果物を削除
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f ai-minesweeper
	go clean -cache

# CI/CD関連
ci-test: ## CI環境でのテスト実行
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# 開発前チェック（fmt + lint + test）
check: fmt lint test ## 開発前の品質チェック

# プロジェクト情報
info: ## プロジェクト情報を表示
	@echo "AI Minesweeper Project"
	@echo "====================="
	@echo "Go version: $$(go version)"
	@echo "Module: $$(go list -m)"