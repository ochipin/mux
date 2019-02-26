package mux

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/ochipin/locale"
	"github.com/ochipin/logger/errorlog"
	"github.com/ochipin/mux/basemux"
	"github.com/ochipin/mux/helpers"
	"github.com/ochipin/render"
	"github.com/ochipin/render/core"
	"github.com/ochipin/report"
	"github.com/ochipin/router"
	"github.com/ochipin/uploadfile"
)

// Values : basemux.Values のエイリアス
type Values = basemux.Values

// Mux : basemux.Mux を継承したマルチプレクサ
type Mux struct {
	basemux.Mux
	Router      router.Router           // ルーティングテーブル
	Tables      map[string][][]string   // ルーティングテーブル情報
	RenderFiles render.Render           // クライアント画面管理用構造体
	ErrorsFiles render.Render           // エラー画面管理用構造体
	StaticFiles render.Render           // 静的ファイル管理用構造体
	Log         errorlog.Logger         // ログ管理インタフェース
	Locale      locale.Parse            // 多言語設定
	UploadFiles *uploadfile.UploadFiles // ファイルアップロードの詳細
	Charset     string                  // charset
	BaseURL     string                  // ベースとなるURLデフォルトは'/'
	ContentList map[string]string       // コンテンツリスト
	Helpers     interface{}             // ヘルパ
	RestrictIP  RestrictIP              // IP制限
	Trigger     Trigger                 // トリガ
}

// New : Mux を初期化し、http.Handler を生成する関数
func (mux *Mux) New() (http.Handler, error) {
	// ルーティングテーブル未設定の場合はエラーとする
	if mux.Router == nil {
		return nil, fmt.Errorf("no route table")
	}

	// レンダーオブジェクトが未設定の場合は、デフォルト設定を適用する
	if mux.RenderFiles == nil {
		c := &render.Config{
			Directory: "app/views/contents",
			Targets:   []string{".html", ".text"},
			Cache:     true,
			Binary:    false,
		}
		v, err := c.New()
		if err != nil {
			return nil, err
		}
		mux.RenderFiles = v
	}

	// エラーレンダーオブジェクトが未設定の場合は、デフォルト設定を適用する
	if mux.ErrorsFiles == nil {
		c := &render.Config{
			Directory: "app/views/errors",
			Targets:   []string{".html", ".text"},
			Cache:     true,
			Binary:    false,
		}
		v, err := c.New()
		if err != nil {
			return nil, err
		}
		mux.ErrorsFiles = v
	}

	// 静的ファイル管理用構造体が未設定の場合、デフォルト設定を適用する
	if mux.StaticFiles == nil {
		c := &render.Config{
			Directory: "app/assets",
			Exclude:   regexp.MustCompile(`(^|[|\n])//=\s*(.+?)\s*$|(^|[|\n])/\*=\s*([\s\S]+?)\s*\*/`),
			Cache:     true,
			Binary:    true,
			MaxSize:   5 << 20,
		}
		v, err := c.New()
		if err != nil {
			return nil, err
		}
		mux.StaticFiles = v
	}

	// ロギングインタフェースが未設定の場合、デフォルト値を設定する
	if mux.Log == nil {
		logger := &errorlog.Log{
			Depth:  3,
			Format: "%D %T %b[%p]: %f(%m:%l) %L: %M",
			Level:  7,
		}
		mux.Log, _ = logger.MakeLog(os.Stderr)
	}

	// 多言語判定が未設定の場合、デフォルト値を適用する
	if mux.Locale == nil {
		l := &locale.Locale{
			Default: "ja",
			Langs: map[string][]string{
				"ja": []string{"ja"},
				"en": []string{"en"},
			},
			LocaleDir: "config/locales",
		}
		mux.Locale, _ = l.CreateLocale()
	}

	// ファイルアップロード時の設定が未設定の場合、デフォルト値を適用する
	if mux.UploadFiles == nil {
		mux.UploadFiles = &uploadfile.UploadFiles{
			SaveFile:  "files/%Y/%m/%y%m%d%H%M%S_%g_%f",
			Overwrite: true,
			Perm:      0644,
			MaxSize:   10 << 20,
		}
	}

	// 文字コード設定が未設定の場合、デフォルト値であるUTF-8を適用する
	if mux.Charset == "" {
		mux.Charset = "UTF-8"
	}

	// Content-Type リストが未登録の場合、デフォルト値を適用する
	if mux.ContentList == nil {
		mux.ContentList = basemux.ContentList()
	}

	// BaseURL が未登録の場合は、'/' が BaseURL とする
	mux.BaseURL = strings.Trim(mux.BaseURL, " ")
	if mux.BaseURL == "" {
		mux.BaseURL = "/"
	}
	// 先頭に'/'が付与されていない場合は付与する
	if mux.BaseURL[0] != '/' {
		mux.BaseURL = "/" + mux.BaseURL
	}
	// 最後尾の'/'は除外する
	if mux.BaseURL != "/" {
		mux.BaseURL = strings.TrimRight(mux.BaseURL, "/")
	}

	// ヘルパが未登録の場合、デフォルトのヘルパ設定を登録する
	if mux.Helpers == nil {
		mux.Helpers = helpers.Helpers{}
	}
	// 登録されているヘルパが構造体型ではない場合、エラーを返却する
	types := reflect.TypeOf(mux.Helpers)
	if types.Kind() != reflect.Struct {
		return nil, fmt.Errorf("'Helpers' parameter type not struct")
	}

	// IP制限設定をされている場合、IPNetを初期化する
	if mux.RestrictIP == nil {
		mux.RestrictIP = RestrictIP{}
	}
	if mux.RestrictIP != nil {
		if err := mux.RestrictIP.MakeIPNet(); err != nil {
			return nil, err
		}
	}

	// トリガ未設定の場合は、空トリガを記憶させる
	if mux.Trigger == nil {
		mux.Trigger = &BaseTrigger{}
	}
	mux.Handler = mux
	return mux.GenerateHandler()
}

// Main : Muxエントリポイント
func (mux *Mux) Main(w http.ResponseWriter, r *http.Request, refer basemux.Referer, v *basemux.Values) (basemux.Render, error) {
	mux.Log.Debug("BEGIN")
	var err error

	// アクセスが '/' の場合のみ、BaseURLの確認を行う
	if r.URL.Path == "/" {
		// BaseURL が '/' のみではない場合、BaseURLへリダイレクトする
		if mux.BaseURL != "/" {
			return &Render{
				Path:       mux.BaseURL,
				StatusCode: 302,
			}, nil
		}
	}

	// 最後に、Commitをコール
	defer func() {
		if e := mux.Trigger.Commit(w, r, v); e != nil {
			if err == nil {
				err = e
			}
		}
	}()

	// アクション実行結果を受け取る
	object := mux.CallAction(w, r, v)

	switch result := object.(type) {
	// エラー
	case error:
		err = result
	// レンダリング
	case basemux.Render:
		return result, nil
	}

	mux.Log.Debug("END")
	// 上記以外
	return nil, err
}

// CallAction : アクション実行結果を判定し、返却する
func (mux *Mux) CallAction(w http.ResponseWriter, r *http.Request, v *Values) interface{} {
	mux.Log.Debug("BEGIN")
	// ヘルパの雛形を作成し、バッファにデフォルトヘルパのアドレスを登録する
	helper := &helpers.Helpers{
		Params: helpers.Parameters{
			"charset": mux.Charset,
			"lang":    mux.Locale.Lookup(r.Header.Get("Accept-Language")),
		},
		MethodName:   mux.MethodName,
		LinkID:       v.ID(),
		RemoteURI:    &helpers.URL{URL: r.URL},
		Locale:       mux.Locale,
		BaseURL:      mux.BaseURL,
		LangData:     mux.I18n(r, "", ""),
		SubmitMethod: r.Method,
	}
	v.Set("defaultHelper", mux.Trigger.SetHelper(helper))
	mux.Log.Debug("default helper created.")

	// アクションを実行し、実行結果を判定する
	var result interface{}
	switch object := mux.ExecAction(w, r, v, helper).(type) {
	// エラーが発生している場合は、どのタイプのエラーか判定する
	case error:
		result = object
	// エラー発生ではない場合,どの型がコントローラから返却されたのか判定する
	default:
		switch types := object.(type) {
		// HTML/TEXT/JSON/XML のいずれかの表示
		case *RenderTemplate:
			object, err := mux.Render(types, v)
			if err != nil {
				mux.Log.Error(err)
				return err
			}
			// RenderTemplate で登録されたオリジナルヘルパ、データを記憶させる
			v.Set("helper", types.helper)
			v.Set("data", types.data)
			result = object
		// 静的ファイルを処理の表示
		case *AssetsTemplate:
			var id string
			// クエリパラメータからリンクIDを取得
			if query := r.URL.Query(); query != nil {
				id = query.Get("id")
			}
			// 呼び出し元のコントローラ名、アクション名、リンクID、言語名に置き換える
			source := v.Old(id)
			if source != nil {
				// バッファ内容を呼び出し元へ書き換える
				v.Set("ctlname", source.Get("ctlname"))
				v.Set("actname", source.Get("actname"))
				v.Set("execname", source.Get("ctlname")+"."+source.Get("actname"))
				v.Set("linkid", source.Get("linkid"))
				v.Set("helper", source.Val("helper"))
				v.Set("data", v.Val("data"))
				// 多言語情報を取得する
				langdata := mux.I18n(r, v.Get("ctlname"), v.Get("actname"))
				// ヘルパのデータを書き換える
				helper.Params["controller"] = source.Get("ctlname")
				helper.Params["action"] = v.Get("actname")
				helper.LangData = langdata
				helper.LinkID = v.Get("linkid")
			}
			// 静的ファイルを表示する
			object, err := mux.Static(types, v)
			if err != nil {
				return err
			}
			result = object
		// リダイレクトの場合
		case *Redirect:
			// https, http から始まる文字列ではない場合、BaseURLを考慮したパスに変換する
			if !regexp.MustCompile(`^http[s]*://`).MatchString(types.path) {
				// /baseurl/path/to/url/ => baseurl/path/to/url
				trimpath := strings.Trim(mux.BaseURL+"/"+types.path, "/")
				// baseurl/path/to/url => [baseurl path to url]
				paths := strings.Fields(strings.Replace(trimpath, "/", " ", -1))
				// [baseurl path to url] => /baseurl/path/to/url
				types.path = "/" + strings.Join(paths, "/")
			}
			// リダイレクト用オブジェクトを作成
			result = &Render{
				Path:       types.path,
				StatusCode: types.statuscode,
			}
		// 上記以外の場合、復帰値エラーとして扱う
		default:
			result = &InvalidReturn{
				Message: "invalid return value. value is nil",
			}
		}
	}

	mux.Log.Debug("END")
	return result
}

// ExecAction : アクションを実行する
func (mux *Mux) ExecAction(w http.ResponseWriter, r *http.Request, v *Values, helper *helpers.Helpers) interface{} {
	mux.Log.Debug("BEGIN")
	// 事前共通処理を実施
	if err := mux.Trigger.Begin(w, r, v); err != nil {
		mux.Log.Error(err)
		return err
	}

	// ルーティングテーブルから、アクセスされたクエリパスに該当するアクション情報を取得する
	res, args, err := mux.RoutePath(r)
	if err != nil {
		mux.Log.Error(err)
		return err
	}

	// アクション情報を取得
	action, _ := res.Get()
	// アクション情報が所持するコントローラ名とアクション名を取得
	ctlname, actname := res.Name()
	// リンクIDを記憶させる
	linkid := v.ID()
	// 多言語情報を取得する
	langdata := mux.I18n(r, ctlname, actname)

	// 一時バッファに格納する
	v.Set("ctlname", ctlname)
	v.Set("actname", actname)
	v.Set("execname", ctlname+"."+actname)
	v.Set("linkid", linkid)
	mux.Log.Debug("buffer setting complete")

	// ヘルパのデータを完成させる
	helper.Params["controller"] = ctlname
	helper.Params["action"] = actname
	helper.LangData = langdata
	helper.LinkID = linkid
	// 独自ヘルパを設定
	// v.Set("defaultHelper", mux.Trigger.SetHelper(helper))
	mux.Log.Debug("helpers.Helper parameters set complete")

	// 基本コントローラを生成
	controller := &Controller{
		w:           w,                                  // http.ResponseWriter
		r:           r,                                  // *http.Request
		controller:  ctlname,                            // コントローラ名
		action:      actname,                            // アクション名
		locale:      langdata,                           // 多言語設定
		files:       uploadfile.New(r, mux.UploadFiles), // アップロードファイルの成約
		path:        r.URL.Path,                         // クエリパス
		Log:         mux.Log,                            // ロギング
		contentlist: mux.ContentList,                    // Content-Type 一覧
	}

	// アクション情報に、コントローラをセット
	mux.Trigger.SetController(action, v)
	if err := router.SetStruct(action, controller); err != nil {
		mux.Log.Error(err)
		return err
	}

	// 実行するアクションが正しい情報で構築されているか確認
	fn, err := res.Valid(action, args, "mux.Result")
	if err != nil {
		mux.Log.Error(err)
		return err
	}

	// 型情報に問題がなければコントローラが所持するPrePostRegisterをコールする
	var prepost = &PrePost{}
	_, err = res.Callname(action, "PrePostRegister", []reflect.Value{
		reflect.ValueOf(prepost),
	})
	if err != nil {
		mux.Log.Error(err)
		return err
	}

	// アクション実行前に、事前関数を実行する
	for _, v := range prepost.begins {
		// 事前関数の復帰値が、nil以外の場合は、処理を中断し関数を復帰する
		if result := v(); result != nil {
			return result
		}
	}
	// アクションを実行し、実行結果を返却する
	out := fn.Call(args)
	// アクション実行後、事後関数を実行する
	for _, v := range prepost.commits {
		// 事後関数の復帰値が、nil以外の場合は処理を中断し関数を復帰する
		if result := v(); result != nil {
			return result
		}
	}

	// 事前、事後関数共に復帰値がnilの場合、アクション実行結果を検証する
	if out[0].Interface() == nil {
		return nil
	}

	mux.Log.Debug("END")
	return out[0].Interface()
}

// RoutePath : アクセスされたクエリパスから該当するアクションを取得する
func (mux *Mux) RoutePath(r *http.Request) (router.Result, []reflect.Value, error) {
	var path = r.URL.Path
	// 先頭のクエリパスがBaseURLで設定したパスではない場合、エラーを返却する
	if strings.Index(path, mux.BaseURL) != 0 {
		return nil, nil, &router.NotRoutes{
			Message: fmt.Sprintf("'[%s]: %s' - not found", r.Method, r.URL.Path),
			Path:    r.URL.Path,
			Method:  r.Method,
		}
	}
	// /baseurl/path/to/url => /path/to/url へ変換する
	if mux.BaseURL != "/" {
		path = path[len(mux.BaseURL):]
		if path == "" {
			path = "/"
		}
		if path[0] != '/' {
			path = r.URL.Path
		}
	}

	// IP 制限がかかっていないか確認する
	var addr string
	idx := strings.Index(r.RemoteAddr, ":")
	if idx != -1 {
		addr = r.RemoteAddr[:idx]
	}
	if mux.RestrictIP.Contains(path, addr) == false {
		return nil, nil, &AccessDenied{
			Message:    "access forbidden by rule, client: " + addr,
			IP:         addr,
			StatusCode: 403,
		}
	}

	// アクセスされたクエリパスに該当するアクションを取得する
	res, args, err := mux.Router.Caller(r.Method, path)
	// アクション取得成功の場合、アクションの情報を返却する
	if err == nil {
		return res, args, nil
	}

	// 該当するアクションが見つからない場合、"*" で登録されたクエリパスがないか確認する
	if res, args, err := mux.Router.Caller("*", path); err == nil {
		return res, args, nil
	}

	// ルーティングテーブルから該当するアクションが見つからない場合は、エラーを返却する
	return nil, nil, err
}

// I18n : Accept-Languageを使用して多言語情報を取得する
func (mux *Mux) I18n(r *http.Request, ctlname, actname string) locale.Data {
	mux.Log.Debug("BEGIN")
	// コントローラ名とアクション名を小文字にする
	ctlname = strings.ToLower(ctlname)
	actname = strings.ToLower(actname)
	// Accept-Language から該当する言語名を取得する
	langname := mux.Locale.Lookup(r.Header.Get("Accept-Language"))
	// 該当する言語情報と、コントローラとアクションに該当する言語情報取得
	l1 := mux.Locale.Locale(langname)
	// コントローラ、またはアクション名が存在しない場合、1つの言語名のみを返却する
	if ctlname == "" || actname == "" {
		mux.Log.Debug("controller, action is empty. default language used")
		if l1 == nil {
			mux.Log.Notice("'%s' language file is nil", langname)
		}
		mux.Log.Debug("END")
		return l1
	}
	// コントローラとアクションに応じた言語情報を取得
	l2 := mux.Locale.Locale(fmt.Sprintf("%s/%s/%s", ctlname, actname, langname))
	mux.Log.Debugf("'%s.%s' is '%s' language used", ctlname, actname, langname)
	mux.Log.Debug("END")
	// l1 + l2 のマージした情報を返却する
	mergedata := locale.Merge(l1, l2)
	if mergedata != nil {
		mergedata = make(locale.Data)
	}
	return mergedata
}

// Render : basemux.Render を生成する
func (mux *Mux) Render(r *RenderTemplate, v *Values) (basemux.Render, error) {
	mux.Log.Debug("BEGIN")
	render := mux.RenderFiles.Copy()

	// 独自ヘルパを設定されている場合、ヘルパを登録する
	if r.helper != nil {
		if err := render.Helper(r.helper); err != nil {
			return nil, err
		}
	}
	// デフォルトヘルパを登録
	if err := render.SmallHelper(v.Val("defaultHelper")); err != nil {
		mux.Log.Error(err)
		return nil, err
	}

	// 返却するRender型を生成
	var result = &Render{
		StatusCode:  r.statuscode,
		ContentType: fmt.Sprintf("%s; %s", r.content, mux.Charset),
	}

	switch r.ext {
	// JSON関数がコールされている場合、Data関数で登録されたデータをJSONとして処理する
	case "json":
		var data map[string]interface{}
		buf, _ := json.Marshal(r.data)
		if err := json.Unmarshal(buf, &data); err != nil {
			mux.Log.Error(err)
			return nil, &MarshalError{err.Error()}
		}
		result.Buffer = buf
	// XML関数がコールされている場合、Data関数で登録されたデータをXMLとして処理する
	case "xml":
		var data map[string]interface{}
		buf, _ := xml.Marshal(r.data)
		if err := xml.Unmarshal(buf, &data); err != nil {
			mux.Log.Error(err)
			return nil, &MarshalError{err.Error()}
		}
		result.Buffer = buf
	// HTML/TEXT関数がコールされている場合、HTML、またはTEXTとして処理する
	case "html", "text":
		mux.Log.Debugf("'%s.%s' html/text output", r.ctlname, r.actname)
		buf, err := render.Render(r.path+"."+r.ext, r.data)
		if err != nil {
			mux.Log.Error(err)
			return nil, err
		}
		result.Buffer = buf
	}

	mux.Log.Debug("END")
	return result, nil
}

// Static : 静的ファイルを処理する
func (mux *Mux) Static(r *AssetsTemplate, v *Values) (basemux.Render, error) {
	mux.Log.Debug("BEGIN")
	render := mux.StaticFiles.Copy()

	// 独自ヘルパを設定されている場合、ヘルパを登録する
	if helper := v.Val("helper"); helper != nil {
		if err := render.Helper(helper); err != nil {
			return nil, err
		}
	}

	// デフォルトヘルパを登録
	if err := render.SmallHelper(v.Val("defaultHelper")); err != nil {
		return nil, err
	}

	// 返却するRender型を生成
	var result = &Render{
		StatusCode:  200,
		ContentType: r.content,
	}

	// 静的ファイルを処理
	buf, err := render.Render(r.path, v.Val("data"))
	if err != nil {
		return nil, err
	}
	result.Buffer = buf

	mux.Log.Debug("END")
	return result, nil
}

// Error : エントリポイントからエラーが返却されるとコールされるエラー画面出力関数
func (mux *Mux) Error(err error, res http.ResponseWriter, req *http.Request) {
	mux.Log.Debug("BEGIN")
	r := mux.ErrorsFiles.Copy()
	v := res.(*basemux.ResponseWriter)

	// デフォルトヘルパを登録
	if helper := v.Val("defaultHelper"); helper != nil {
		if err := r.SmallHelper(helper); err != nil {
			mux.Log.Error("default helper regsiter error")
			res.Header().Set("Content-Type", "text/html")
			res.WriteHeader(500)
			res.Write([]byte(err.Error()))
			return
		}
	}

	// 登録されているコントローラ名とアクション名を取得
	execname := v.Get("execname")
	if execname == "" {
		execname = "???.???"
	}

	// エラータイプを判定する
	var status = &ErrorStatus{
		Trace:     report.ServeTrace(0, res, req),
		Message:   err.Error(),
		Interface: err,
		r:         req,
	}

	switch types := err.(type) {
	// IP制限に引っかかった際のエラー
	case *AccessDenied:
		status.Title = "403 Forbidden"
		status.StatusCode = types.StatusCode
		status.ErrorTitle = "403 Access Denied"
		status.StatusName = "AccessDenied"
	// ルーティングテーブルから、クエリパスに該当するアクションが見つからない
	case *router.NotRoutes:
		status.Title = "404 Not Found"
		status.StatusCode = 404
		status.StatusName = "NotRoutes"
		status.ErrorTitle = "'" + types.Path + "' query path not found"
		status.Message = "action corresponding to '" + types.Path + "' query path not found in routing table"
		status.Interface = mux.Tables
	// 実行するアクションの引数の型情報に誤りがある
	case *router.IllegalArgs:
		status.Title = "500 Internal Server Error"
		status.StatusCode = 500
		status.StatusName = "IllegalArguments"
		status.ErrorTitle = "Illegal Arguments in '" + execname + "'"
	// 実行するアクションの引数の数に誤りがある
	case *router.NotEnoughArgs:
		status.Title = "500 Internal Server Error"
		status.StatusCode = 500
		status.StatusName = "NotEnoughArguments"
		status.ErrorTitle = "Not Enough Arguments in '" + execname + "'"
	// 実行するアクションの復帰値の型情報に誤りがある
	case *router.IllegalRets:
		status.Title = "500 Internal Server Error"
		status.StatusCode = 500
		status.StatusName = "IllegalReturn"
		status.ErrorTitle = "Illegal Return in '" + execname + "'"
	// 実行するアクションの復帰値の数に誤りがある
	case *router.NotEnoughRets:
		status.Title = "500 Internal Server Error"
		status.StatusCode = 500
		status.StatusName = "NotEnoughReturn"
		status.ErrorTitle = "Not Enough Return in '" + execname + "'"
	// ミックスイン構造体にセットする値が nil ポインタをセットされている
	case *router.InvalidError:
		status.Title = "500 Internal Server Error"
		status.StatusCode = 500
		status.StatusName = "MixInInvalid"
		status.ErrorTitle = "Mix-in Invalid arguments in '" + execname + "'"
	// ミックスイン構造体にセットする値が構造体型ではない場合のエラー
	case *router.NoStruct:
		status.Title = "500 Internal Server Error"
		status.StatusCode = 500
		status.StatusName = "MixInNoStructType"
		status.ErrorTitle = "Mix-in No Struct Type error in '" + execname + "'"
	// 実行するアクションに、必要な構造体がミックスインされていない
	case *router.NoMixin:
		status.Title = "500 Internal Server Error"
		status.StatusCode = 500
		status.StatusName = "NoMixIn"
		status.ErrorTitle = "No Mix-in '" + types.Basicname + "' in '" + execname + "'"
	// ヘルパ登録失敗
	case *core.HelperInvalid:
		status.Title = "500 Internal Server Error"
		status.StatusCode = 500
		status.StatusName = "HelperInvalid"
		status.ErrorTitle = "Helper Register Error in '" + execname + "'"
	// レンダーパースエラー
	case *core.TemplateError:
		status.Title = "500 Internal Server Error"
		status.StatusCode = 500
		status.StatusName = "TemplateError"
		status.ErrorTitle = "Template Error in '" + execname + "'"
	// JSON/XML の解析失敗時のエラー
	case *MarshalError:
		status.Title = "500 Internal Server Error"
		status.StatusCode = 500
		status.StatusName = "MarshalError"
		status.ErrorTitle = "Marshal Error in '" + execname + "'"
	// レンダーエラー
	case *core.RenderError:
		status.Title = "500 Internal Server Error"
		status.StatusCode = 500
		status.StatusName = "RenderError"
		status.ErrorTitle = "Render Error in '" + execname + "'"
		// エラー内容を配列化
		rep := strings.NewReplacer("<", "&lt;", ">", "&gt;")
		root := strings.Split(rep.Replace(types.Root), "\n")
		for i := 0; i < len(root); i++ {
			if types.Line-1 == i {
				root[i] = "<li class='alert'>" + root[i] + "</li>"
			} else {
				root[i] = "<li>" + root[i] + "</li>"
			}
		}
		// エラー情報をInterfaceメンバへ格納
		status.Interface = map[string]interface{}{
			"basename": types.Basename,
			"message":  err.Error(),
			"line":     types.Line,
			"root":     root,
		}
	// 実行時エラー
	case *basemux.PanicError:
		status.Title = "500 Internal Server Error"
		status.StatusCode = 500
		status.StatusName = "RuntimeError"
		status.ErrorTitle = "Runtime Error in '" + execname + "'"
		status.Trace.StackTrace = types.StackTrace
	// タイムアウトエラー
	case *basemux.TimeoutError:
		status.Title = "408 Request Time-out"
		status.StatusCode = 408
		status.StatusName = "Timeout"
		status.ErrorTitle = "Request Time-out in '" + execname + "'"
	// 最大同時アクセス数の超過エラー
	case *basemux.MaxClientsError:
		status.Title = "503 Service Temporarily Unavailable"
		status.StatusCode = 503
		status.StatusName = "MaxClientsOver"
		status.ErrorTitle = "Max Clients Over in '" + execname + "'"
	// 401 認証エラー
	case *Unauthorized:
		status.Title = "401 Unauthorized"
		status.StatusCode = types.StatusCode
		status.StatusName = "Unauthorized"
		status.ErrorTitle = "401 Unauthorized in '" + execname + "'"
		res.Header().Add("WWW-Authenticate", `Basic realm="`+types.Title+`"`)
	// コントローラから InternalError が返却された場合
	case *InternalError:
		status.Title = "500 Internal Server Error"
		status.StatusCode = types.statuscode
		status.StatusName = "InternalError"
		status.ErrorTitle = "Internal Server Error in '" + execname + "'"
		status.Interface = types.data
		status.Trace.StackTrace = types.trace
	// メンテナンス中である場合のエラー
	case *Maintenance:
		status.Title = "503 Service Temporarily Unavailable"
		status.StatusCode = types.statuscode
		status.StatusName = "Maintenance"
		status.ErrorTitle = "The Maintenance Mode in '" + execname + "'"
		status.Interface = types.data
	// コントローラからNotFoundが返却された場合
	case *NotFound:
		status.Title = "404 Not Found"
		status.StatusCode = types.statuscode
		status.StatusName = "NotFound"
		status.ErrorTitle = "Not Found in '" + execname + "'"
		status.Interface = types.data
	// コントローラからForbiddenが返却された場合
	case *Forbidden:
		status.Title = "403 Forbidden"
		status.StatusCode = types.statuscode
		status.StatusName = "Forbidden"
		status.ErrorTitle = "Forbidden in '" + execname + "'"
		status.Interface = types.data
	// コントローラの復帰値が nil の場合
	case *InvalidReturn:
		status.Title = "500 Internal Server Error"
		status.StatusCode = 500
		status.StatusName = "InvalidReturn"
		status.ErrorTitle = "Invalid Return Value in '" + execname + "'"
	// 事前共通関数のエラー
	case *BeginError:
		status.Title = "500 Internal Server Error"
		status.StatusCode = 500
		status.StatusName = "BeginError"
		status.ErrorTitle = "An error occurred in preprocessing '" + req.URL.Path + "'"
	// 事後共通関数のエラー
	case *CommitError:
		status.Title = "500 Internal Server Error"
		status.StatusCode = 500
		status.StatusName = "CommitError"
		status.ErrorTitle = "An error occurred in post processing '" + req.URL.Path + "'"
	// 上記以外のエラー
	default:
		res.Header().Set("Content-Type", "text/html")
		res.WriteHeader(500)
		res.Write([]byte(err.Error()))
		return
	}

	mux.Log.Debugf("status name is '%s'", status.StatusName)
	defer func() {
		// エラー発生時には、エラーレポートトリガ関数をコール
		mux.Trigger.ErrorReport(status, status.StatusCode, status.StatusName)
	}()

	// 静的ファイルを処理
	buf, err := r.Render("errors.html", status)
	if err != nil {
		mux.Log.Error(err)
		res.Header().Set("Content-Type", "text/html")
		res.WriteHeader(status.StatusCode)
		res.Write([]byte(status.Error()))
		return
	}

	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(status.StatusCode)
	res.Write(buf)

	mux.Log.Debug("END")
}
