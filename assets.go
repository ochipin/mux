package mux

// AssetsTemplate : 静的ファイルを取り扱う構造体
type AssetsTemplate struct {
	path    string // クエリパス (ex: /assets/stylesheets/controllers/index.css)
	content string // 拡張子に応じた Content-Type
}

func (r *AssetsTemplate) pointer() {}
