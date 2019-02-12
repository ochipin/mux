package helpers

import (
	"fmt"
	"net/url"
	"testing"
)

// Helpers 構造体をテストするための設定値を付与したHelpers構造体を生成する
func CreateHelper() *Helpers {
	url, _ := url.Parse("https://localhost:8080/path/to/url?query=sample")

	return &Helpers{
		Params: Parameters{
			"controller": "Test",
			"action":     "Action",
			"charset":    "UTF-8",
			"lang":       "ja",
		},
		LangData: map[string]interface{}{
			"name": "ok",
		},
		BaseURL:      "baseurl",
		LinkID:       "0123456789",
		MethodName:   "_method",
		RemoteURI:    &URL{url},
		SubmitMethod: "POST",
	}
}

// {{add}} のチェック
func Test_add(t *testing.T) {
	helper := CreateHelper()
	if helper.Add(1, 1) != 2 {
		t.Fatal("ERROR")
	}
}

// {{min}} のチェック
func Test_min(t *testing.T) {
	helper := CreateHelper()
	if helper.Min(10, 2) != 8 {
		t.Fatal("ERROR")
	}
}

// {{mul}} のチェック
func Test_mul(t *testing.T) {
	helper := CreateHelper()
	if helper.Mul(2, 3) != 6 {
		t.Fatal("ERROR")
	}
}

// {{div}} のチェック
func Test_div(t *testing.T) {
	helper := CreateHelper()
	num, _ := helper.Div(10, 2)
	if num != 5 {
		t.Fatal("ERROR")
	}

	_, err := helper.Div(10, 0)
	if err == nil {
		t.Fatal("ERROR")
	}
}

// {{mod}} のチェック
func Test_mod(t *testing.T) {
	helper := CreateHelper()
	num, _ := helper.Mod(10, 2)
	if num != 0 {
		t.Fatal("ERROR")
	}

	if _, err := helper.Mod(10, 0); err == nil {
		t.Fatal("ERROR")
	}
}

// {{sprintf}} テスト
func Test__Sprintf(t *testing.T) {
	helper := CreateHelper()

	// 空文字列を返却する
	if helper.Sprintf() != "" {
		t.Fatal("ERROR")
	}
	// test_string を返却する
	if helper.Sprintf("test_string") != "test_string" {
		t.Fatal("ERROR")
	}
	// test_string を返却する
	if helper.Sprintf("%s_%s", "test", "string") != "test_string" {
		t.Fatal("ERROR")
	}
}

// StringTypeのMatchテスト
func Test_StringType_Match(t *testing.T) {
	testString := CreateHelper().Sprintf("test_string")

	// string文字列が含まれているかチェック
	ok, err := testString.Match("string")
	if err != nil {
		t.Fatal("ERROR")
	}
	if !ok {
		t.Fatal("ERROR")
	}
}

func Test_StringType(t *testing.T) {
	testString := CreateHelper().Sprintf("test_string")
	// 大文字変換
	if testString.Upper() != "TEST_STRING" {
		t.Fatal("ERROR")
	}
	// タイトル変換
	if testString.Title() != "Test_string" {
		t.Fatal("ERROR", testString, testString.Title())
	}
}

// StringTypeのTemplateテスト
func Test_StringType_Template(t *testing.T) {
	testString := CreateHelper().Sprintf("{{.value}}")
	// {{.value}} に対応するパラメータを付与
	v, err := testString.Template(map[string]interface{}{
		"value": "test_string",
	})
	// err が帰ってきた場合は、テスト失敗
	if err != nil {
		t.Fatal(err)
	}
	// test_string 文字列ではない場合はテスト失敗
	if v != "test_string" {
		t.Fatal(v)
	}
}

// MakeMap テスト
func Test_MakeMap(t *testing.T) {
	helper := CreateHelper()

	p1 := helper.MakeMap()
	p1.Set("t1", "v1")
	if p1.HasItem("t1") != true {
		t.Fatal("ERROR")
	}
	if p1.HasItem("t2") == true {
		t.Fatal("ERROR")
	}

	p2 := p1.Copy()
	p2.Set("t2", "v2")
	if p2.HasItem("t1") != true {
		t.Fatal("ERROR")
	}
	if p2.HasItem("t2") != true {
		t.Fatal("ERROR")
	}
	p2.Delete("t2")
	if p2.HasItem("t2") == true {
		t.Fatal("ERROR")
	}
	p2.Clear()
	if p1.HasItem("t1") != true {
		t.Fatal("ERROR")
	}
}

// Parameters に設定されている値をチェックする
func Test_Params(t *testing.T) {
	helper := CreateHelper()
	if helper.Controller() != "Test" {
		t.Fatal("ERROR", helper.Controller())
	}
	if helper.Action() != "Action" {
		t.Fatal("ERROR")
	}
	if helper.ID() != "Test_Action" {
		t.Fatal("ERROR")
	}
	if helper.Class() != "Test_Action" {
		t.Fatal("ERROR")
	}
	if helper.Title() != "Test.Action" {
		t.Fatal("ERROR")
	}
	helper.Set("id", "id_name")
	helper.Set("class", "class_name")
	helper.Set("title", helper.Sprintf("title_name"))
	if helper.ID() != "id_name" {
		t.Fatal("ERROR")
	}
	if helper.Class() != "class_name" {
		t.Fatal("ERROR")
	}
	if helper.Title() != "title_name" {
		t.Fatal("ERROR")
	}
}

func Test_Date(t *testing.T) {
	helper := CreateHelper()
	fmt.Println(helper.Date(""))
	fmt.Println(helper.Date("%Y-%m-%d", "1 days"))
}
func Test_Form(t *testing.T) {
	helper := CreateHelper()
	fmt.Println(helper.Form(helper.NoComplete(), helper.NoValidate(), helper.Multipart()))
	fmt.Println(helper.Form())
}

// {{href}} テスト
func Test_Href(t *testing.T) {
	helper := CreateHelper()
	// BaseURL の文字列が付与されたクエリパスが返却される
	if helper.Href("/path/to/url") != "/baseurl/path/to/url" {
		t.Fatal("ERROR")
	}
	// https, http から始まる文字列の場合、BaseURLは付与されない
	if helper.Href("https://localhost:8080/path/to/url") != "https://localhost:8080/path/to/url" {
		t.Fatal("ERROR")
	}
}

// {{url}} テスト
func Test_URL(t *testing.T) {
	helper := CreateHelper()
	if helper.URL().Hostname() != "localhost" {
		t.Fatal("ERROR")
	}
	if helper.URL().Port() != "8080" {
		t.Fatal("ERROR")
	}
	if helper.URL().Host() != "localhost:8080" {
		t.Fatal("ERROR")
	}
	if helper.URL().Path() != "/path/to/url" {
		t.Fatal("ERROR")
	}
	if helper.URL().Proto() != "https" {
		t.Fatal("ERROR")
	}
	if helper.URL().Search() != "query=sample" {
		t.Fatal("ERROR")
	}
	if fmt.Sprint(helper.URL().Get("query")) != "sample" {
		t.Fatal("ERROR")
	}
	if fmt.Sprint(helper.URL().Get("query2")) == "sample" {
		t.Fatal("ERROR")
	}
	if helper.URL().Origin() != "https://localhost:8080" {
		t.Fatal("ERROR")
	}
}
