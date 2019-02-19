package mux

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ochipin/router"
)

type Example struct {
	*Controller
}

func (c *Example) PrePostRegister(prepost *PrePost) {
	// 事前コール関数を登録
	prepost.AddBeginFunc(func() Result {
		return nil
	})
	// 事後コール関数を登録
	prepost.AddCommitFunc(func() Result {
		return nil
	})
}

func (c *Example) Hello() Result {
	return c.Render().HTML()
}

func (c *Example) World() Result {
	fmt.Println(c.I18n("sample"))
	return c.Render().Template("example/%s", "world")
}

func (c *Example) InternalError() Result {
	return c.Controller.InternalError("Internal Server Error")
}

func (c *Example) NotFound() Result {
	return c.Controller.NotFound("Not Found")
}

// 正常系のテストを行う
func NewMux() (http.Handler, error) {
	// ルーティングテーブルを作成
	r := router.New()

	// コントローラを登録
	r.AddClass(Example{})
	r.AddClass(Serve{})

	// 正規表現を登録
	r.AddRegexp("static", `(.+\.[a-zA-Z0-9_]+)`)

	// ルートパスを登録
	r.Register("GET", "/", "Example.Hello")
	r.Register("GET", "/world", "Example.World")
	r.Register("GET", "/assets/:static", "Serve.File")

	// ルーティングテーブルを生成
	routes, err := r.Create()
	if err != nil {
		return nil, err
	}

	// Mux を生成
	var mux = &Mux{
		Router: routes,
		Tables: r.TableList(),
	}

	handler, err := mux.New()
	if err != nil {
		return nil, err
	}

	return handler, nil
}

func Test_Mux1(t *testing.T) {
	// ハンドラを受け取る
	handler, err := NewMux()
	if err != nil {
		t.Fatal(err)
	}
	// サーバを起動
	ts := httptest.NewServer(handler)
	defer ts.Close()

	client := ts.Client()

	// Example.Hello アクションを実行
	res, err := client.Get(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}

	var buf = make([]byte, res.ContentLength)
	res.Body.Read(buf)
	res.Body.Close()

	// Example.World アクションを実行
	res, err = client.Get(ts.URL + "/world")
	if err != nil {
		t.Fatal(err)
	}

	buf = make([]byte, res.ContentLength)
	res.Body.Read(buf)
	res.Body.Close()
}

func Test_Mux2(t *testing.T) {
	// ハンドラを受け取る
	handler, err := NewMux()
	if err != nil {
		t.Fatal(err)
	}
	// サーバを起動
	ts := httptest.NewServer(handler)
	defer ts.Close()

	client := ts.Client()

	// NotFoundError
	res, err := client.Get(ts.URL + "/notfound")
	if err != nil {
		t.Fatal(err)
	}

	var buf = make([]byte, res.ContentLength)
	res.Body.Read(buf)
	res.Body.Close()
}
