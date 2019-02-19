package mux

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/ochipin/locale"
	"github.com/ochipin/logger/errorlog"
	"github.com/ochipin/uploadfile"
)

// PrePost : アクション関数の事前、事後にコールする関数を管理する構造体
type PrePost struct {
	begins  []func() Result
	commits []func() Result
}

// AddBeginFunc : 事前関数の登録
func (prepost *PrePost) AddBeginFunc(fn func() Result) {
	prepost.begins = append(prepost.begins, fn)
}

// AddCommitFunc : 事後関数の登録
func (prepost *PrePost) AddCommitFunc(fn func() Result) {
	prepost.commits = append(prepost.commits, fn)
}

// Result : アクションの復帰値
type Result interface {
	pointer()
}

// Controller : コントローラ
type Controller struct {
	r           *http.Request
	w           http.ResponseWriter
	path        string
	controller  string
	action      string
	Log         errorlog.Logger
	files       *uploadfile.File
	contentlist map[string]string
	locale      locale.Data
}

// PrePostRegister : アクション実行前の事前、事後実行関数を登録する初期化関数
func (c *Controller) PrePostRegister(prepost *PrePost) {}

// Form : 入力フォームから得た情報を返却する
func (c *Controller) Form() url.Values {
	return c.r.PostForm
}

// Query : クエリパラメータの情報を返却する
func (c *Controller) Query() url.Values {
	return c.r.URL.Query()
}

// I18n : 言語設定パラメータの値を取得する
func (c *Controller) I18n(name string) string {
	return fmt.Sprint(c.locale.T(name))
}

// Render : HTML/TEXT/JSON のいずれかを表示する
func (c *Controller) Render() *RenderTemplate {
	return &RenderTemplate{
		ctlname:    c.controller,
		actname:    c.action,
		path:       strings.ToLower(c.controller) + "/" + strings.ToLower(c.action),
		statuscode: 200,
		ext:        "html",
		content:    "text/html",
	}
}

// Redirect : 301, 302 リダイレクトを行う
func (c *Controller) Redirect(i ...interface{}) *Redirect {
	var path = c.path
	if len(i) == 1 {
		path = fmt.Sprint(i[0])
	} else if len(i) > 1 {
		path = fmt.Sprintf(fmt.Sprint(i[0]), i[1:]...)
	}

	return &Redirect{
		path:       path,
		statuscode: 302,
	}
}

// InternalError : 500 InternalServerError を発生させる
func (c *Controller) InternalError(message string, i ...interface{}) *InternalError {
	// トレース情報を取得
	var trace []string
	for i := 0; ; i++ {
		pc, filename, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		funcname := runtime.FuncForPC(pc).Name()
		trace = append(trace, fmt.Sprintf("=======>> %d: %s: %s(%d)", i+1, funcname, filename, line))
	}
	// フォーマット指定子の場合は、message を Sprintf で整形
	if len(i) > 1 {
		message = fmt.Sprintf(message, i...)
	}
	// InternalError を返却する
	return &InternalError{
		ErrorReturn: &ErrorReturn{
			trace:      trace,
			message:    message,
			statuscode: 500,
		},
	}
}

// Maintenance : 503 Service Temporarily Unavailable を発生させる
func (c *Controller) Maintenance(message string, i ...interface{}) *Maintenance {
	// フォーマット指定子の場合は、message を Sprintf で整形
	if len(i) > 1 {
		message = fmt.Sprintf(message, i...)
	}
	// Maintenance を返却する
	return &Maintenance{
		ErrorReturn: &ErrorReturn{
			message:    message,
			statuscode: 503,
		},
	}
}

// Forbidden : 403 Forbidden を発生させる
func (c *Controller) Forbidden(message string, i ...interface{}) *Forbidden {
	// フォーマット指定子の場合は、message を Sprintf で整形
	if len(i) > 1 {
		message = fmt.Sprintf(message, i...)
	}
	// Forbidden を返却する
	return &Forbidden{
		ErrorReturn: &ErrorReturn{
			message:    message,
			statuscode: 403,
		},
	}
}

// NotFound : 404 Not Found を発生させる
func (c *Controller) NotFound(message string, i ...interface{}) *NotFound {
	// フォーマット指定子の場合は、message を Sprintf で整形
	if len(i) > 1 {
		message = fmt.Sprintf(message, i...)
	}
	// Forbidden を返却する
	return &NotFound{
		ErrorReturn: &ErrorReturn{
			message:    message,
			statuscode: 404,
		},
	}
}

// Serve : 静的ファイルを取り扱うコントローラ
type Serve struct {
	*Controller
}

// File : 静的ファイルを表示するアクション
func (c *Serve) File(name string) Result {
	// ex) name => /assets/stylesheets/controllers/index.css
	var path = name
	var content string

	// 拡張子が、Content-Typeリストに登録されている拡張子の場合
	if idx := strings.LastIndex(name, "."); idx != -1 {
		ext := name[idx:]
		// Content-Typeリストに拡張子が登録されているか確認する
		if v, ok := c.contentlist[ext]; ok {
			content = v
		}
	}
	// Content-Typeリストに存在しない場合は、バイナリとして扱う
	if content == "" {
		content = "application/octet-stream"
	}

	// 静的ファイル情報を返却する
	return &AssetsTemplate{
		content: content,
		path:    path,
	}
}

// Maintenance : メンテナンス画面を表示する
func (c *Serve) Maintenance() Result {
	return c.Controller.Maintenance("Maintenance")
}
