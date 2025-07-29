.PHONY: setup lint fmt test build clean

# セットアップ
setup:
	@echo "Setting up git hooks..."
	git config core.hooksPath .githooks
	@echo "Installing golangci-lint..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Setup complete!"

# lint実行
lint:
	golangci-lint run

# フォーマット
fmt:
	go fmt ./...

# テスト実行
test:
	go test ./...

# ビルド
build:
	go build -o ai-minesweeper

# クリーン
clean:
	rm -f ai-minesweeper

# 開発前チェック（fmt + lint + test）
check: fmt lint test