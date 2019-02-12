package mux

import (
	"fmt"
	"net/http"
)

// Render : basemux.Render インタフェースに対応したRender構造体
type Render struct {
	Buffer      []byte // 表示する内容
	StatusCode  int    // ステータスコード
	ContentType string // Content-Type
	Path        string // リダイレクト先のパス
}

// Render : HTML/TEXT/JSON 等を表示する関数
func (render *Render) Render(w http.ResponseWriter, r *http.Request) {
	if render.StatusCode == 301 || render.StatusCode == 302 {
		http.Redirect(w, r, render.Path, render.StatusCode)
	} else {
		w.Header().Set("Content-Type", render.ContentType)
		w.WriteHeader(render.StatusCode)
		w.Write(render.Buffer)
	}
}

// RenderTemplate : ビュー情報を管理する構造体
type RenderTemplate struct {
	ctlname    string      // コントローラ名 (Base)
	actname    string      // アクション名 (Index)
	path       string      // テンプレート名(ex: base/index)
	content    string      // Content-Type(ex: text/html)
	ext        string      // 拡張子(ex: .html)
	statuscode int         // 2xx, 4xx, 5xx などのエラー値
	data       interface{} // ビュー内で使用するデータ
	helper     interface{} // ビュー内で使用する関数
}

// Data : ビュー内で使用するデータ、もしくはJSON/XMLデータを登録する
func (r *RenderTemplate) Data(i interface{}) *RenderTemplate {
	r.data = i
	return r
}

// Code : ステータスコードを指定する(デフォルトは200)
func (r *RenderTemplate) Code(status int) *RenderTemplate {
	r.statuscode = status
	return r
}

// Helper : ビュー内で使用するヘルパを登録する
func (r *RenderTemplate) Helper(i interface{}) *RenderTemplate {
	r.helper = i
	return r
}

// Template : 表示するテンプレートファイルを指定する
func (r *RenderTemplate) Template(path string, i ...interface{}) *RenderTemplate {
	// 引数で渡されたパスが空文字列の場合、何もせず復帰する
	if path == "" {
		return r
	}
	r.path = path
	if len(i) > 0 {
		r.path = fmt.Sprintf(path, i...)
	}
	return r
}

// HTML : HTMLを出力する
func (r *RenderTemplate) HTML() *RenderTemplate {
	r.content = "text/html"
	r.ext = "html"
	return r
}

// TEXT : text/plain でテキストを出力する
func (r *RenderTemplate) TEXT() *RenderTemplate {
	r.content = "text/plain"
	r.ext = "text"
	return r
}

// JSON : Data関数で登録されたデータをJSON形式に変換して出力する
func (r *RenderTemplate) JSON() *RenderTemplate {
	r.content = "application/json"
	r.ext = "json"
	return r
}

// XML : Data関数で登録されたデータをXML形式に変換して出力する
func (r *RenderTemplate) XML() *RenderTemplate {
	r.content = "application/xml"
	r.ext = "xml"
	return r
}

// mux.Result に対応するための、空メソッド
func (r *RenderTemplate) pointer() {}
