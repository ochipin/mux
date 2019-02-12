package basemux

import (
	"fmt"
	"net/http"
	"sync"
)

// Values : キーと値でデータを管理する構造体
type Values struct {
	mu     sync.Mutex
	id     string
	access int64
	data   map[string]interface{}
	old    Referer
}

// ID : リンクIDを返却する
func (v *Values) ID() string {
	return v.id
}

// Set : 値をセットする
func (v *Values) Set(key string, val interface{}) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.data[key] = val
}

// Get : セットした値を文字列で取得する
func (v *Values) Get(key string) string {
	if val := v.Val(key); val != nil {
		return fmt.Sprint(val)
	}
	return ""
}

// Val : セットした値をinterface{}で取得する
func (v *Values) Val(key string) interface{} {
	v.mu.Lock()
	defer v.mu.Unlock()
	val, _ := v.data[key]
	return val
}

// Old : 過去にリファラに登録されていた情報を、指定したIDで取得する
func (v *Values) Old(id string) *Values {
	return v.old.Get(id)
}

// ResponseValues : http.ResponseWriter + Values を合わせた機能を持つインタフェース
type ResponseValues interface {
	Header() http.Header
	Write([]byte, int, error)
	WriteHeader(int)
	Set(string, interface{})
	Get(string) string
	Val(string) interface{}
}

// ResponseWriter : 独自ResponseWriter
type ResponseWriter struct {
	http.ResponseWriter
	*Values
}
