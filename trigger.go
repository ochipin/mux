package mux

import (
	"net/http"
	"reflect"

	"github.com/ochipin/mux/helpers"
)

// Trigger : Mux構造体に登録するトリガ
type Trigger interface {
	Begin(http.ResponseWriter, *http.Request, *Values) error
	Commit(http.ResponseWriter, *http.Request, *Values) error
	SetController(reflect.Value, *Values)
	SetHelper(*helpers.Helpers) interface{}
	ErrorReport(err error, code int, status string)
}

// BaseTrigger : Mux構造体に登録するベースとなるトリガ。実装は空となっており、何もしない。
type BaseTrigger struct{}

// Begin : アクション実行前に何らかの処理を実行するトリガ
func (t BaseTrigger) Begin(w http.ResponseWriter, r *http.Request, v *Values) error { return nil }

// Commit : アクション実行後に何らかの処理を実行するトリガ
func (t BaseTrigger) Commit(w http.ResponseWriter, r *http.Request, v *Values) error { return nil }

// SetController : 独自コントローラを登録するトリガ
func (t BaseTrigger) SetController(elem reflect.Value, v *Values) {}

// SetHelper : 独自ヘルパを登録するトリガ
func (t BaseTrigger) SetHelper(helpers *helpers.Helpers) interface{} { return helpers }

// ErrorReport : 2xx 以外のエラー発生時にコールされるトリガ
func (t BaseTrigger) ErrorReport(err error, code int, status string) {}
