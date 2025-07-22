# 設計書

## 概要

ToDoリストアプリケーションは、GolangとEbitengineを使用したデスクトップアプリケーションとして実装します。JSONファイルを使用してデータを永続化し、Ebitengineの2Dグラフィック機能を活用してユーザーフレンドリーなインターフェースを提供します。

## アーキテクチャ

### 技術スタック
- **言語**: Go 1.24
- **グラフィックエンジン**: Ebitengine v2
- **データストレージ**: JSONファイル（ローカルファイルシステム）
- **UI**: Ebitengineの描画機能を使用したカスタムUI

### アプリケーション構造
```
todo-app/
├── main.go            # エントリーポイント
├── internal/
│   ├── game/
│   │   └── game.go    # メインゲームループ
│   ├── ui/
│   │   ├── button.go  # ボタンコンポーネント
│   │   ├── textbox.go # テキスト入力コンポーネント
│   │   └── todoitem.go # ToDoアイテムコンポーネント
│   ├── models/
│   │   └── todo.go    # Todoデータモデル
│   └── storage/
│       └── storage.go # ファイルストレージ管理
├── assets/
│   └── fonts/         # フォントファイル
├── data/
│   └── todos.json     # ToDoデータファイル
├── go.mod
├── go.sum
└── README.md
```

## コンポーネントとインターフェース

### 1. Todo データモデル
```go
type Todo struct {
    ID        string    `json:"id"`
    Text      string    `json:"text"`
    Completed bool      `json:"completed"`
    CreatedAt time.Time `json:"created_at"`
}

type TodoList struct {
    Todos []Todo `json:"todos"`
}
```

### 2. Game メインストラクト
```go
type Game struct {
    todos       []Todo
    currentFilter FilterType
    inputText   string
    editingID   string
    ui          *UIManager
    storage     *Storage
}

// 主要メソッド
func (g *Game) AddTodo(text string) error
func (g *Game) DeleteTodo(id string) error
func (g *Game) ToggleTodo(id string) error
func (g *Game) EditTodo(id, newText string) error
func (g *Game) SetFilter(filter FilterType)
func (g *Game) Update() error
func (g *Game) Draw(screen *ebiten.Image)
```

### 3. Storage インターフェース
```go
type Storage interface {
    SaveTodos(todos []Todo) error
    LoadTodos() ([]Todo, error)
    ClearTodos() error
}

type FileStorage struct {
    filepath string
}
```

### 4. UIコンポーネント

#### UIManager
```go
type UIManager struct {
    inputBox    *TextBox
    addButton   *Button
    filterButtons map[FilterType]*Button
    todoItems   []*TodoItem
    scrollOffset int
}
```

#### Button コンポーネント
```go
type Button struct {
    X, Y, Width, Height int
    Text               string
    OnClick            func()
    Hovered            bool
    Pressed            bool
}
```

#### TextBox コンポーネント
```go
type TextBox struct {
    X, Y, Width, Height int
    Text               string
    Focused            bool
    CursorPos          int
}
```

#### TodoItem コンポーネント
```go
type TodoItem struct {
    Todo       *Todo
    X, Y       int
    Width      int
    Height     int
    Editing    bool
    EditText   string
    Checkbox   *Button
    DeleteBtn  *Button
}
```

## データモデル

### Todo 構造体
```go
type Todo struct {
    ID        string    `json:"id"`        // 一意識別子（UUID）
    Text      string    `json:"text"`      // タスクテキスト
    Completed bool      `json:"completed"` // 完了状態
    CreatedAt time.Time `json:"created_at"` // 作成日時
}
```

### フィルタータイプ
```go
type FilterType int

const (
    FilterAll FilterType = iota
    FilterActive
    FilterCompleted
)
```

### アプリケーション状態
```go
type AppState struct {
    Todos         []Todo     // タスク配列
    CurrentFilter FilterType // 現在のフィルター
    EditingID     string     // 編集中のタスクID
    InputText     string     // 入力中のテキスト
}
```

## エラーハンドリング

### 入力検証
- 空文字列のタスク追加を防止
- 編集時の空文字列を防止
- 不正なIDでの操作を防止
- 文字列長制限の実装

### ファイルストレージエラー
- ファイル読み書きエラーの処理
- JSON解析エラーの処理
- ディスク容量不足への対応
- ファイル権限エラーの処理

### ユーザーフィードバック
- エラーメッセージの画面表示
- 操作成功時の視覚的フィードバック
- ローディング状態の表示（ファイル操作時）

### エラータイプ定義
```go
type AppError struct {
    Type    ErrorType
    Message string
    Err     error
}

type ErrorType int

const (
    ErrorValidation ErrorType = iota
    ErrorStorage
    ErrorUI
)
```

## テスト戦略

### 単体テスト
- Todo 構造体のメソッドテスト
- Storage インターフェースの実装テスト
- Game の各メソッドテスト
- UIコンポーネントの個別テスト

### 統合テスト
- UI操作とデータ更新の連携テスト
- ファイルストレージとの連携テスト
- フィルタリング機能のテスト
- キーボード・マウス入力のテスト

### テストツール
- Go標準のtestingパッケージ
- testifyライブラリ（アサーション）
- モックファイルシステム（テスト用）

## UI/UXデザイン

### デザイン原則
- シンプルで直感的なインターフェース
- キーボード操作対応（Tab、Enter、Escapeキー）
- ウィンドウサイズ変更対応

### カラーパレット
```go
var (
    ColorPrimary   = color.RGBA{0, 123, 255, 255}  // #007bff（青）
    ColorSecondary = color.RGBA{108, 117, 125, 255} // #6c757d（グレー）
    ColorSuccess   = color.RGBA{40, 167, 69, 255}   // #28a745（緑）
    ColorDanger    = color.RGBA{220, 53, 69, 255}   // #dc3545（赤）
    ColorBackground = color.RGBA{248, 249, 250, 255} // #f8f9fa（ライトグレー）
    ColorText      = color.RGBA{33, 37, 41, 255}    // #212529（ダークグレー）
)
```

### フォント設定
- システムフォントを使用
- ベースフォントサイズ: 16px
- 行間: 1.5倍

### ウィンドウ設定
- 初期サイズ: 800x600px
- 最小サイズ: 400x300px
- リサイズ可能

### インタラクション
- マウスホバー効果: 色の変化
- クリック時の視覚的フィードバック
- キーボードフォーカス表示
- スムーズなアニメーション（Ebitengineの描画機能を使用）

## パフォーマンス考慮事項

### 最適化戦略
- 描画の最小化（変更があった場合のみ再描画）
- 効率的なスライス操作
- ファイルI/Oの最小化
- メモリプールの活用

### メモリ管理
- 不要なオブジェクトの適切な解放
- 大量データ時の仮想スクロール実装検討
- Goのガベージコレクション最適化

### Ebitengine固有の最適化
- 画像の事前読み込み
- テキスト描画の最適化
- フレームレート制御（60FPS）
- バッチ描画の活用