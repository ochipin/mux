package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ochipin/locale"
)

// Helpers : ビュー内で使用する関数群を管理する構造体
type Helpers struct {
	MethodName   string       // <form>タグ生成時に付与されるメソッド名を取り出すキー名
	FormData     *Form        // <form>タグを生成するマップ
	Locale       locale.Parse // 言語パース
	LangData     locale.Data  // 言語設定情報
	Params       Parameters   // controller, action, language, charset,
	LinkID       string       // <link rel=... 時に同時に付与されるリンクID
	BaseURL      string       // ベースURL
	RemoteURI    *URL         // URL情報
	SubmitMethod string       // リクエスト情報に付与されるメソッド名
}

// Add : 足し算コマンド
func (cmd *Helpers) Add(a, b int) int {
	return a + b
}

// Min : 引き算コマンド
func (cmd *Helpers) Min(a, b int) int {
	return a - b
}

// Mul : 掛け算コマンド
func (cmd *Helpers) Mul(a, b int) int {
	return a * b
}

// Div : 割り算コマンド
func (cmd *Helpers) Div(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("error: integer divide by zero")
	}
	return a / b, nil
}

// Mod : 剰余演算コマンド
func (cmd *Helpers) Mod(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("error: integer divide by zero")
	}
	return a % b, nil
}

// Sprintf : 与えられた引数をStringType型へ変換する
func (cmd *Helpers) Sprintf(i ...interface{}) StringType {
	switch len(i) {
	case 0:
		return ""
	case 1:
		return StringType(fmt.Sprint(i[0]))
	}
	return StringType(fmt.Sprintf(fmt.Sprint(i[0]), i[1:]...))
}

// MakeMap : テンプレート内で扱うマップ型を返却する
func (cmd *Helpers) MakeMap() Parameters {
	return make(Parameters)
}

// Date : 日付を返却する
func (cmd *Helpers) Date(strs ...string) (*DateTime, error) {
	// 引数が2つ以上指定されている場合は、エラーとする
	if len(strs) > 2 {
		return nil, fmt.Errorf("date: invalid date arguments")
	}

	// 引数がない場合、デフォルト時刻表記を返却する
	if len(strs) == 0 {
		return &DateTime{
			datetime: time.Now(),
			format:   "%a %b %d %H:%M:%S %Z %Y",
		}, nil
	}

	// フォーマットが空文字列の場合、デフォルトフォーマットを指定する
	if strs[0] == "" {
		strs[0] = "%a %b %d %H:%M:%S %Z %Y"
	}

	var now = &DateTime{format: strs[0]}
	// 1つしか引数がない場合は、現在時刻を取得して関数を復帰する
	if len(strs) == 1 {
		now.datetime = time.Now()
		return now, nil
	}

	// 引数が2つ指定されている場合、2つ目の引数を解析する
	args := strings.Fields(strings.ToLower(strs[1]))
	// 引数が不正な場合、エラーとみなす
	if len(args) < 1 || len(args) > 3 {
		return nil, fmt.Errorf("date: invalid date '%s'", strs[1])
	}

	// 1つめの引数は、整数でなければならない。整数でない場合は、エラーとみなす
	n, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, fmt.Errorf("date: invalid date '%s'", strs[1])
	}
	duration := time.Duration(n)

	switch len(args) {
	// 引数の数が1個しかない場合は、分、秒を0とする
	case 1:
		// 0-23 の間でなければエラーとみなす
		if n < 0 || n > 23 {
			return nil, fmt.Errorf("date: invalid date '%s'", strs[1])
		}
		d := time.Now()
		now.datetime = time.Date(d.Year(), d.Month(), d.Day(), n, 0, 0, 0, time.Local)
	// 引数の数が2つの場合は、指定した単位で日付をずらす
	case 2:
		switch args[1] {
		case "second", "seconds":
			now.datetime = time.Now().Add(duration * time.Second)
		case "minute", "minutes":
			now.datetime = time.Now().Add(duration * time.Minute)
		case "hour", "hours":
			now.datetime = time.Now().Add(duration * time.Hour)
		case "day", "days":
			now.datetime = time.Now().AddDate(0, 0, n)
		case "month", "months":
			now.datetime = time.Now().AddDate(0, n, 0)
		case "year", "years":
			now.datetime = time.Now().AddDate(n, 0, 0)
		default:
			return nil, fmt.Errorf("date: invalid date '%s'", strs[1])
		}
	}

	return now, nil
}

// Multipart : <form enctype='multipart/form-data' を実現する
func (cmd *Helpers) Multipart() Multipart { return "" }

// NoComplete : <form autocomplete='off' を実現する
func (cmd *Helpers) NoComplete() AutoComplete { return "" }

// NoValidate : <form novalidate='novalidate' を実現する
func (cmd *Helpers) NoValidate() NoValidate { return "" }

// Form : <form> を生成する
func (cmd *Helpers) Form(params ...interface{}) (*Form, error) {
	// form タグが閉じていない状態で、formがコールされた場合、閉じタグを生成する
	if cmd.FormData != nil {
		cmd.FormData.end = true
		form := cmd.FormData
		cmd.FormData = nil
		return form, nil
	}

	var form = &Form{
		data: map[string]interface{}{
			"action": cmd.RemoteURI.URL.Path,
		},
		keyname: cmd.MethodName,
		end:     false,
		method:  "POST",
	}
	cmd.FormData = form

	// 引数がない場合は、デフォルトフォームタグを返却する
	if len(params) == 0 {
		return form, nil
	}

	// 引数で指定されたパラメータを解析
	for _, v := range params {
		switch types := v.(type) {
		// enctype='multipart/form-data' 指定
		case Multipart:
			form.Attr("enctype", "multipart/form-data")
		// autocomplete='off' 指定
		case AutoComplete:
			form.Attr("autocomplete", "off")
		// novalidate='novalidate' 指定
		case NoValidate:
			form.Attr("novalidate", "novalidate")
		// 文字列の場合、先頭文字から
		case string, StringType:
			t := strings.TrimSpace(fmt.Sprint(types))
			if len(t) == 0 {
				continue
			}
			switch t[0] {
			// 先頭文字が'/'の場合、クエリパス指定とみなす
			case '/':
				form.Attr("action", t)
			// 先頭文字が'#'の場合、id指定とみなす
			case '#':
				form.Attr("id", t[1:])
			// 先頭文字が'.'の場合、class指定とみなす
			case '.':
				form.Attr("class", t[1:])
			// 先頭文字が':'の場合、target指定とみなす
			case ':':
				form.Attr("target", t[1:])
			// 先頭文字が'@'の場合、accept-charset指定とみなす
			case '@':
				form.Attr("charset", t[1:])
			// 先頭文字が'$'の場合、method指定とみなす
			case '$':
				form.method = t[1:]
			// 特殊記号文字がない場合、name指定とみなす
			default:
				form.Attr("name", t)
			}
		// 上記以外の型が指定された場合、エラーを返却する
		default:
			return nil, fmt.Errorf("form: invalid form arguments")
		}
	}
	return form, nil
}

// Controller : コントローラ名を返却する
func (cmd *Helpers) Controller() StringType {
	return cmd.Params.T("controller")
}

// Action : アクション名を返却する
func (cmd *Helpers) Action() StringType {
	return cmd.Params.T("action")
}

// ID : id 属性値を返却する
func (cmd *Helpers) ID() StringType {
	if id := cmd.Params.T("id"); id != "" {
		return id
	}
	return cmd.Controller() + "_" + cmd.Action()
}

// Class : class 属性値を返却する
func (cmd *Helpers) Class() StringType {
	if class := cmd.Params.T("class"); class != "" {
		return class
	}
	return cmd.Controller() + "_" + cmd.Action()
}

// Charset : charset に設定した値(UTF-8)を返却する
func (cmd *Helpers) Charset() StringType {
	return cmd.Params.T("charset")
}

// Lang : 適用されている自然言語名を返却する
func (cmd *Helpers) Lang() StringType {
	return cmd.Params.T("lang")
}

// Title : <title>タグに設定するタイトル名を返却する
func (cmd *Helpers) Title() StringType {
	if title := cmd.Params.T("title"); title != "" {
		return title
	}
	return cmd.Controller() + "." + cmd.Action()
}

// Set : キーと値で、変数を取り扱う
func (cmd *Helpers) Set(key string, value interface{}) string {
	var v interface{}
	switch types := value.(type) {
	case string:
		v = StringType(types)
	default:
		v = value
	}
	cmd.Params[key] = v
	return ""
}

// Parameter : 設定されている変数を取り出す
func (cmd *Helpers) Parameter(name string) interface{} {
	return cmd.Params.T(name)
}

// T : 言語情報を取得する
func (cmd *Helpers) T(name string) interface{} {
	if cmd.LangData == nil {
		return ""
	}
	return cmd.LangData.T(name)
}

// Global : i18n 関数の第2引数に指定する型
type Global struct{}

// Global : i18n 関数の第2引数に指定することで、全体言語設定を変更可能
func (cmd *Helpers) Global() Global { return Global{} }

// I18n : 言語情報を返却する
func (cmd *Helpers) I18n(name string, i ...Global) interface{} {
	var result locale.Data

	if cmd.Locale.LangList(name) {
		// ja, en などの locale.Locale で設定した言語名の場合、
		// jaファイルと、hello/world/jaなどのコントローラとアクションに該当する
		// 言語ファイルを結合した言語情報を対象とする
		l1 := cmd.Locale.Locale(name)
		ctlname := cmd.Params.T("controller").Lower().String()
		actname := cmd.Params.T("action").Lower().String()
		l2 := cmd.Locale.Locale(ctlname + "/" + actname + "/" + name)
		result = locale.Merge(l1, l2)
	} else {
		// hello/world/ja などパス指定の場合、指定したパスの言語ファイルのみを対象とする
		result = cmd.Locale.Locale(name)
	}

	if result == nil {
		result = make(locale.Data)
	}

	if len(i) > 0 {
		cmd.LangData = result
		return ""
	}

	return result
}

// Hostname : ホスト名表示コマンド
func (cmd *Helpers) Hostname() (StringType, error) {
	v, err := os.Hostname()
	return StringType(v), err
}

// Env : 環境変数を取得する
func (cmd *Helpers) Env(name string) StringType {
	return StringType(os.Getenv(name))
}

// Stylesheet : <link rel='stylesheet' ... タグを埋め込む
func (cmd *Helpers) Stylesheet(path string) string {
	path = filepath.Join(cmd.BaseURL, path)
	// リンクIDがない場合、idクエリパラメータなしの <link>を生成する
	if cmd.LinkID == "" {
		return fmt.Sprintf("<link rel='stylesheet' type='text/css' href='%s' />", path)
	}
	// リンクIDが存在する場合、idクエリパラメータにリンクIDを付与する
	return fmt.Sprintf("<link rel='stylesheet' type='text/css' href='%s?id=%s' />", path, cmd.LinkID)
}

// Script : <script src='...' タグを埋め込む
func (cmd *Helpers) Script(path string) string {
	path = filepath.Join(cmd.BaseURL, path)
	// リンクIDがない場合、idクエリパラメータなしの <link>を生成する
	if cmd.LinkID == "" {
		return fmt.Sprintf("<script src='%s'></script>", path)
	}
	// リンクIDが存在する場合、idクエリパラメータにリンクIDを付与する
	return fmt.Sprintf("<script src='%s?id=%s'></script>", path, cmd.LinkID)
}

// URL : アクセス先URLを返却する
func (cmd *Helpers) URL() *URL {
	return cmd.RemoteURI
}

// IsFile : ファイル有無を確認する
func (cmd *Helpers) IsFile(fname StringType) bool {
	_, err := os.Stat(fname.String())
	return err == nil
}

// IsDir : ディレクトリ有無を確認する
func (cmd *Helpers) IsDir(dir StringType) bool {
	f, err := os.Stat(dir.String())
	if err != nil || f.IsDir() == false {
		return false
	}
	return true
}

// Stat : ファイル情報を取得する
func (cmd *Helpers) Stat(fname StringType) (*FileInfo, error) {
	f, err := os.Stat(fname.String())
	if err != nil {
		return nil, err
	}
	var result = &FileInfo{f}
	return result, nil
}

// Item : 与えられた文字列を、HTMLの属性値にセットするよう加工する
func (cmd *Helpers) Item(value StringType) string {
	// ...value=''<name>'&"<key>"'
	// => value='&#39;&lt;name&gt;&#39;&amp;&quot;&lt;key&gt;&quot;'
	rep := strings.NewReplacer(
		"'", "&#39;",
		"\"", "&quot;",
		"<", "&lt;",
		">", "&gt;",
		"&", "&amp;",
	)
	return rep.Replace(value.String())
}

// Href : リンク先を生成する
func (cmd *Helpers) Href(path StringType) string {
	// 先頭、最後尾にある空白を除去
	url := strings.Trim(path.String(), " ")
	// https://, http:// から始まる場合は、別のURLへの遷移として扱う
	if strings.Index(url, "https://") == 0 || strings.Index(url, "http://") == 0 {
		return url
	}
	// /baseurl/path/to/url/ => baseurl/path/to/url => baseurl path to url => [baseurl path to url]
	paths := strings.Fields(strings.Replace(strings.Trim(cmd.BaseURL+"/"+url, "/"), "/", " ", -1))
	// [baseurl path to url] => baseurl/path/to/url => /baseurl/path/to/url
	return "/" + strings.Join(paths, "/")
}

// Method : メソッド名を取得する
func (cmd *Helpers) Method() string {
	return cmd.SubmitMethod
}

// FileInfo : Stat 関数が返却するファイル情報を取り扱う構造体
type FileInfo struct {
	os.FileInfo
}

// ModTime : ファイルの更新日、作成日を返却する
func (f *FileInfo) ModTime() *DateTime {
	return &DateTime{
		datetime: f.FileInfo.ModTime(),
		format:   "%w %b %d %H:%M:%S %Z %Y",
	}
}
