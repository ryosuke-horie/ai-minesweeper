# AIマインスイーパー

AIにネタバレされるマインスイーパー。論理的に確定できるマスはAIが自動で開き、運が絡む部分だけプレイヤーが選択します。

## インストール

```bash
git clone https://github.com/ryosuke-horie/ai-minesweeper.git
cd ai-minesweeper
go build -o ai-minesweeper
```

## 実行

```bash
./ai-minesweeper
```

## 操作方法

- **↑↓←→** / **hjkl**: カーソル移動
- **スペース** / **Enter**: マスを開く
- **f**: 旗を立てる/外す
- **r**: 新しいゲーム
- **1/2/3**: 難易度変更（初級/中級/上級）
- **q** / **Ctrl+C** / **Ctrl+Q**: 終了

## 必要環境

- Go 1.21以上
