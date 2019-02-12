package basemux

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

// TestHandler : basemux.Mux を使用するためのハンドラ
type TestHandler struct{}
type PanicType struct {
	Message string
}

func (err *PanicType) Error() string {
	return err.Message
}

// エラー処理
func (h *TestHandler) Error(err error, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	switch t := err.(type) {
	case *PanicType:
		panic("PANIC")
	case *TimeoutError:
		w.WriteHeader(t.StatusCode)
	case *PanicError:
		w.WriteHeader(t.StatusCode)
	case *MaxClientsError:
		w.WriteHeader(t.StatusCode)
	default:
		w.WriteHeader(500)
	}
	w.Write([]byte(err.Error()))
}

// Main : リクエストを処理するエントリポイント
func (h *TestHandler) Main(w http.ResponseWriter, r *http.Request, refer Referer, v *Values) (Render, error) {
	switch r.URL.Path {
	// 正常処理のチェック
	case "/":
		return &View{
			Buffer:      []byte("HELLO WORLD"),
			ContentType: "text/html",
			StatusCode:  200,
		}, nil
	// リダイレクトチェック
	case "/redirect":
		return &View{
			StatusCode: 301,
			Path:       "/",
		}, nil
	// error 発生
	case "/error":
		return nil, fmt.Errorf("ERROR")
	// Error 関数内でのパニック
	case "/panic":
		return nil, &PanicType{"PANIC"}
	// PANIC発生時のチェック
	case "/error/panic":
		panic("panic")
	// リファラに登録された内容をチェック
	case "/value":
		v.Set("key", "value")
		if v.Get("key") != "value" {
			panic("/value")
		}
		// 存在しないキーを指定された場合は、空文字列を返却する
		if v.Get("undefined") != "" {
			panic("/value")
		}
		if fmt.Sprint(v.Val("key")) != "value" {
			panic("/value")
		}
		// 現在のリファラIDを返却する
		if len(v.ID()) == 0 {
			panic("/value")
		}
		// 過去に登録されていたリファラIDを取得する
		if v.Old(v.ID()) == nil {
			panic("/value")
		}
		// 存在しないIDが指定された場合は、エラーとする
		if v.Old("undefined") != nil {
			panic("/value")
		}
		// ここまでチェックが完了すれば処理OK
		return &View{
			Buffer:      []byte("SUCCESS"),
			ContentType: "text/html",
			StatusCode:  200,
		}, nil
	case "/timeout":
		time.Sleep(5 * time.Second)
		return &View{
			Buffer:      []byte("TIMEOUT"),
			ContentType: "text/html",
			StatusCode:  200,
		}, nil
	case "/upload":
		return &View{
			Buffer:      []byte("UPLOAD"),
			ContentType: "text/html",
			StatusCode:  200,
		}, nil
	}
	return nil, nil
}

// 正常に終了したかチェック
func Test_Basemux(t *testing.T) {
	// Mux を設定
	mux := &Mux{
		MaxClients: 2,              // 同時受付リクエスト数(2リクエストまで)
		Timeout:    10,             // リクエストタイムアウト(10秒間)
		MaxMemory:  32,             // アップロードファイルの処理に使用する最大使用メモリ量(32MB)
		Handler:    &TestHandler{}, // ハンドラを登録
	}

	// ハンドラを生成
	handler, err := mux.GenerateHandler()
	if err != nil {
		t.Fatal(err)
	}

	// http サーバを立てる
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// リクエストを投入
	client := ts.Client()
	res, err := client.Get(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var buf = make([]byte, res.ContentLength)
	res.Body.Read(buf)
	if string(buf) != "HELLO WORLD" {
		t.Fatal("Request Error")
	}
}

// ハンドラの登録忘れ
func Test_BasemuxHandlerError(t *testing.T) {
	// Mux を設定
	mux := &Mux{
		MaxClients: 2,   // 同時受付リクエスト数(2リクエストまで)
		Timeout:    10,  // リクエストタイムアウト(10秒間)
		MaxMemory:  32,  // アップロードファイルの処理に使用する最大使用メモリ量(32MB)
		Handler:    nil, // ハンドラを登録
	}

	// ハンドラを生成
	_, err := mux.GenerateHandler()
	if err == nil {
		t.Fatal(err)
	}
}

// リダイレクト処理
func Test_BasemuxRedirect(t *testing.T) {
	// Mux を設定
	mux := &Mux{
		MaxClients: 2,              // 同時受付リクエスト数(2リクエストまで)
		Timeout:    10,             // リクエストタイムアウト(10秒間)
		MaxMemory:  32,             // アップロードファイルの処理に使用する最大使用メモリ量(32MB)
		Handler:    &TestHandler{}, // ハンドラを登録
	}

	// ハンドラを生成
	handler, err := mux.GenerateHandler()
	if err != nil {
		t.Fatal(err)
	}

	// http サーバを立てる
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// リクエストを投入
	client := ts.Client()
	res, err := client.Get(ts.URL + "/redirect")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var buf = make([]byte, res.ContentLength)
	res.Body.Read(buf)
	// リダイレクト先の出力結果を得る
	if string(buf) != "HELLO WORLD" {
		t.Fatal("Request Error")
	}
}

// PANIC時の精査
func Test_BasemuxPanic(t *testing.T) {
	// Mux を設定
	mux := &Mux{
		MaxClients: 2,              // 同時受付リクエスト数(2リクエストまで)
		Timeout:    10,             // リクエストタイムアウト(10秒間)
		MaxMemory:  32,             // アップロードファイルの処理に使用する最大使用メモリ量(32MB)
		Handler:    &TestHandler{}, // ハンドラを登録
	}

	// ハンドラを生成
	handler, err := mux.GenerateHandler()
	if err != nil {
		t.Fatal(err)
	}

	// http サーバを立てる
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// リクエストを投入
	client := ts.Client()
	res, err := client.Get(ts.URL + "/error/panic")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var buf = make([]byte, res.ContentLength)
	res.Body.Read(buf)
	if string(buf) != "runtime error. panic" {
		t.Error(string(buf))
	}
}

// リファラチェック
func Test_BasemuxValues(t *testing.T) {
	// Mux を設定
	mux := &Mux{
		MaxClients: 0,              // 同時受付リクエスト数(2リクエストまで)
		Timeout:    0,              // リクエストタイムアウト(10秒間)
		MaxMemory:  0,              // アップロードファイルの処理に使用する最大使用メモリ量(32MB)
		Handler:    &TestHandler{}, // ハンドラを登録
	}
	// ハンドラを生成
	handler, err := mux.GenerateHandler()
	if err != nil {
		t.Fatal(err)
	}
	// http サーバを立てる
	ts := httptest.NewServer(handler)
	defer ts.Close()
	// リクエストを投入
	client := ts.Client()
	res, err := client.Get(ts.URL + "/value")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var buf = make([]byte, res.ContentLength)
	res.Body.Read(buf)
	if string(buf) != "SUCCESS" {
		t.Fatal(string(buf))
	}
}

// リファラ管理チェック
func Test_BasemuxRefer(t *testing.T) {
	// リファラ生成
	ref := newRefer(1)
	// 要素が1つもない状態でのinspectionは何もしない
	ref.inspection(1)

	// リファラ登録データを生成
	values := ref.Create()
	time.Sleep(4 * time.Second)
	// 登録データ削除後、再度データ取得可能かチェックする
	if ref.Get(values.ID()) != nil {
		t.Fatal("ERROR")
	}
}

// Content-List チェック
func Test_ContentList(t *testing.T) {
	list := ContentList()
	if list == nil {
		t.Fatal("ERROR")
	}
}

// タイムアウトチェック
func Test_TimeoutMaxClients(t *testing.T) {
	// Mux を設定
	mux := &Mux{
		MaxClients: 2,              // 同時受付リクエスト数(2リクエストまで)
		Timeout:    2,              // リクエストタイムアウト(10秒間)
		MaxMemory:  32,             // アップロードファイルの処理に使用する最大使用メモリ量(32MB)
		Handler:    &TestHandler{}, // ハンドラを登録
	}

	// ハンドラを生成
	handler, err := mux.GenerateHandler()
	if err != nil {
		t.Fatal(err)
	}

	// http サーバを立てる
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// リクエストを投入
	client := ts.Client()
	res, err := client.Get(ts.URL + "/timeout")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var buf = make([]byte, res.ContentLength)
	res.Body.Read(buf)
	if string(buf) != "request time-out" {
		t.Fatal("Request Error")
	}
	// MaxClients エラー
	go client.Get(ts.URL + "/timeout")
	go client.Get(ts.URL + "/timeout")
	go client.Get(ts.URL + "/timeout")
	time.Sleep(3 * time.Second)
}

// ファイルアップロード時の処理
func Test_Main(t *testing.T) {
	// Mux を設定
	mux := &Mux{
		MaxClients: 2,              // 同時受付リクエスト数(2リクエストまで)
		Timeout:    2,              // リクエストタイムアウト(10秒間)
		MaxMemory:  32,             // アップロードファイルの処理に使用する最大使用メモリ量(32MB)
		Handler:    &TestHandler{}, // ハンドラを登録
	}

	// ハンドラを生成
	handler, err := mux.GenerateHandler()
	if err != nil {
		t.Fatal(err)
	}

	// http サーバを立てる
	ts := httptest.NewServer(handler)
	defer ts.Close()

	client := ts.Client()
	// 対象となるファイルを読み込む
	var buf bytes.Buffer
	var writer = multipart.NewWriter(&buf)
	r, _ := os.Open("Makefile")
	fw, err := writer.CreateFormFile("file1", r.Name())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(fw, r); err != nil {
		t.Fatal(err)
	}
	defer writer.Close()

	req, err := http.NewRequest("POST", ts.URL+"/upload", &buf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	var b = make([]byte, res.ContentLength)
	res.Body.Read(b)
	if string(b) != "UPLOAD" {
		t.Fatal(string(b))
	}
}

// Error 発生テスト
func Test_BasemuxError(t *testing.T) {
	// Mux を設定
	mux := &Mux{
		MaxClients: 2,              // 同時受付リクエスト数(2リクエストまで)
		Timeout:    2,              // リクエストタイムアウト(10秒間)
		MaxMemory:  32,             // アップロードファイルの処理に使用する最大使用メモリ量(32MB)
		Handler:    &TestHandler{}, // ハンドラを登録
	}

	// ハンドラを生成
	handler, err := mux.GenerateHandler()
	if err != nil {
		t.Fatal(err)
	}

	// http サーバを立てる
	ts := httptest.NewTLSServer(handler)
	defer ts.Close()
	// リクエストを投入
	client := ts.Client()
	res, err := client.Get(ts.URL + "/error")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var buf = make([]byte, res.ContentLength)
	res.Body.Read(buf)
	if string(buf) != "ERROR" {
		t.Fatal(string(buf))
	}
}

// PostForm テスト
func Test_PostForm(t *testing.T) {
	// Mux を設定
	mux := &Mux{
		MaxClients: 2,              // 同時受付リクエスト数(2リクエストまで)
		Timeout:    2,              // リクエストタイムアウト(10秒間)
		MaxMemory:  32,             // アップロードファイルの処理に使用する最大使用メモリ量(32MB)
		Handler:    &TestHandler{}, // ハンドラを登録
	}

	// ハンドラを生成
	handler, err := mux.GenerateHandler()
	if err != nil {
		t.Fatal(err)
	}

	// http サーバを立てる
	ts := httptest.NewTLSServer(handler)
	defer ts.Close()

	v := url.Values{}
	v.Set("_method", "DELETE")
	ts.Client().PostForm(ts.URL+"/postform", v)
}

func Test_BasemuxErrorPanic(t *testing.T) {
	// Mux を設定
	mux := &Mux{
		MaxClients: 2,              // 同時受付リクエスト数(2リクエストまで)
		Timeout:    2,              // リクエストタイムアウト(10秒間)
		MaxMemory:  32,             // アップロードファイルの処理に使用する最大使用メモリ量(32MB)
		Handler:    &TestHandler{}, // ハンドラを登録
	}

	// ハンドラを生成
	handler, err := mux.GenerateHandler()
	if err != nil {
		t.Fatal(err)
	}

	// http サーバを立てる
	ts := httptest.NewTLSServer(handler)
	defer ts.Close()
	// リクエストを投入
	client := ts.Client()
	res, err := client.Get(ts.URL + "/panic")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
}
