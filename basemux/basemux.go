package basemux

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type muxHandler struct {
	mux *Mux
}

func (muxHandler *muxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 独自ResponseWriterを作成
	response := &ResponseWriter{
		w, muxHandler.mux.referer.Create(),
	}
	// リクエストを処理
	muxHandler.mux.active(response, r)
}

// View : Main関数の処理後、エラーがなければMainの復帰値として返却される
type View struct {
	Buffer      []byte
	StatusCode  int
	ContentType string
	Path        string
}

// Render : HTML/TEXT 描画関数
func (v *View) Render(w http.ResponseWriter, r *http.Request) {
	if v.StatusCode == 301 || v.StatusCode == 302 {
		http.Redirect(w, r, v.Path, v.StatusCode)
	} else {
		w.Header().Set("Content-Type", v.ContentType)
		w.WriteHeader(v.StatusCode)
		w.Write(v.Buffer)
	}
}

// Render : HTML/TEXT の描画に使用するインタフェース
type Render interface {
	Render(http.ResponseWriter, *http.Request)
}

// Handler : Main関数とError関数を実装したハンドラ
type Handler interface {
	// ルータや、アクションの実行、フィルタリングなどのメイン制御を行う関数
	Main(http.ResponseWriter, *http.Request, Referer, *Values) (Render, error)
	// PANICの発生や、何らかのエラーが発生した場合のエラー画面の制御を行う関数
	Error(error, http.ResponseWriter, *http.Request)
}

// Mux : 受け付けるリクエストの最大数を管理する構造体
type Mux struct {
	MaxClients int           // 最大リクエスト同時接続数
	Timeout    int           // タイムアウト時間(秒単位)
	MaxMemory  int64         // アップロードファイルを処理する際に使用する最大メモリ量(MB単位)
	MethodName string        // オリジナルメソッドキー名
	Handler    Handler       // リクエストを処理するハンドラ
	sem        chan struct{} // リクエスト受付管理チャネル
	referer    *referer      // 1つ前のページ情報を保持している独自リファラ
}

// GenerateHandler : 設定したMux構造体のパラメータから、ハンドラを作成する
func (mux *Mux) GenerateHandler() (http.Handler, error) {
	if mux.Handler == nil {
		return nil, fmt.Errorf("Handler is nil")
	}
	// 最大リクエスト同時接続数が未設定の場合、最大リクエスト同時接続数を100とする
	if mux.MaxClients <= 0 {
		mux.MaxClients = 100
	}
	// リクエスト受付管理チャネルを設定する
	mux.sem = make(chan struct{}, mux.MaxClients)
	for i := 0; i < mux.MaxClients; i++ {
		mux.sem <- struct{}{}
	}
	// タイムアウト時間が未設定の場合、60秒をデフォルトの時間とする
	if mux.Timeout <= 0 {
		mux.Timeout = 60
	}
	// アップロードファイルを処理する際に使用する最大メモリ量を設定(MB単位)
	if mux.MaxMemory <= 0 {
		mux.MaxMemory = 32
	}
	mux.MaxMemory = mux.MaxMemory << 20
	// オリジナルメソッドのキー名が未設定の場合、_methodをキー名とする
	if mux.MethodName == "" {
		mux.MethodName = "_method"
	}
	// リファラ管理構造体を作成する
	mux.referer = newRefer(60)

	return &muxHandler{mux}, nil
}

// リクエストを処理する
func (mux *Mux) active(w http.ResponseWriter, r *http.Request) {
	// 使用するプロトコルを選定
	var proto = "http"
	if r.TLS != nil {
		proto = "https"
	}
	// URLを一旦パースする
	if URL, err := r.URL.Parse(proto + "://" + r.Host + r.RequestURI); err == nil {
		r.URL = URL
	}
	// パース後、クエリパスを整形 (ex: /path/to/url/ => /path/to/url)
	r.URL.Path = strings.TrimRight(r.URL.Path, "/")
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}

	// ステータスコードを管理する変数
	var status int
	// リクエストタイムアウトを検知するコンテキストを生成
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(mux.Timeout)*time.Second)
	// タイムアウトせずに、リクエストを処理したか検知するフラグ
	isfinish := make(chan interface{}, 1)

	defer func() {
		cancel()
		// panic 発生時は recover を実施する
		if err := recover(); err != nil {
			status := PanicDump(0, err)
			if p, ok := status.(*PanicError); ok {
				log.Println(p.Error())
				for i := 0; i < len(p.StackTrace); i++ {
					log.Println(p.StackTrace[i])
				}
			}
		}
		mux.sem <- struct{}{}
	}()

	go func() {
		select {
		// 最大リクエスト同時接続数に到達していない場合、リクエストを受け付ける
		case <-mux.sem:
			status = 408
			mux.main(isfinish, w, r)
		// 最大リクエスト同時接続数を超過している場合、処理を待つ
		default:
			status = 503
			<-mux.sem
			mux.main(isfinish, w, r)
		}
	}()

	// リクエストタイムアウトを検知
	select {
	// タイムアウトの場合、408 or 503 エラーとする
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			mux.Handler.Error(accessError(status), w, r)
		}
	// タイムアウトせずにリクエストを処理
	case v := <-isfinish:
		switch types := v.(type) {
		// リクエスト処理中にエラー発生
		case error:
			mux.Handler.Error(types, w, r)
		// 正常にリクエストを処理
		case Render:
			types.Render(w, r)
		}
	}
}

func (mux *Mux) main(isfinish chan interface{}, w http.ResponseWriter, r *http.Request) {
	var render Render
	var err error

	// main 終了後、実行する
	defer func() {
		if e := recover(); e != nil {
			// PANIC
			isfinish <- PanicDump(0, e)
		} else if err != nil {
			// エントリポイントでエラー発生
			isfinish <- err
		} else {
			// 正常終了
			isfinish <- render
		}
	}()

	if strings.ToUpper(r.Method) != "GET" {
		// ファイルをアップロードされているか検出する
		for _, v := range r.Header["Content-Type"] {
			if strings.Index(v, "multipart/form-data") != -1 {
				r.ParseMultipartForm(mux.MaxMemory)
				break
			}
		}
		r.ParseForm()
		// GET/POST以外のリクエストメソッドを指定している場合、r.Methodに指定されたメソッド名を格納する
		if len(r.PostForm[mux.MethodName]) != 0 && r.PostForm[mux.MethodName][0] != "" {
			r.Method = r.PostForm[mux.MethodName][0]
		}
	}

	v := w.(*ResponseWriter)
	// メイン処理実行
	render, err = mux.Handler.Main(w, r, mux.referer, v.Values)
}
