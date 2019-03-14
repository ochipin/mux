package mux

import (
	"fmt"
	"net/http"

	"github.com/ochipin/report"
)

// ErrorStatus : エラー画面レンダリング時に使用されるステータス構造体
type ErrorStatus struct {
	Trace      *report.Trace
	Message    string
	Title      string
	ErrorTitle string
	StatusCode int
	StatusName string
	Interface  interface{}
	r          *http.Request
}

func (err *ErrorStatus) Error() string {
	return fmt.Sprintf("%d \"%s\" [%s \"%s\"] %s",
		err.StatusCode, err.r.Method, err.r.RemoteAddr, err.r.URL.Path, err.Message)
}

// ErrorReturn : アクションの復帰値エラーを管理するベース構造体
type ErrorReturn struct {
	trace      []string
	message    string
	statuscode int
	data       interface{}
}

// Data : エラー画面内で使用するデータを登録する
func (err *ErrorReturn) Data(i interface{}) {
	err.data = i
}

// Code : ステータスコード変更する
func (err *ErrorReturn) Code(status int) {
	err.statuscode = status
}

func (err ErrorReturn) Error() string {
	return err.message
}

func (err *ErrorReturn) pointer() {}

// Maintenance : 503 Service Temporarily Unavailable
type Maintenance struct {
	*ErrorReturn
}

// InternalError : 500 Internal Server Error
type InternalError struct {
	*ErrorReturn
}

// NotFound : 404 Not Found
type NotFound struct {
	*ErrorReturn
}

// Forbidden : 403 Forbidden
type Forbidden struct {
	*ErrorReturn
}

// InvalidReturn : アクションの復帰値が nil
type InvalidReturn struct {
	Message string
}

func (err *InvalidReturn) Error() string {
	return err.Message
}

// MarshalError : json/xml 等の解析エラー
type MarshalError struct {
	Message string
}

func (err *MarshalError) Error() string {
	return err.Message
}

// BeginError : 事前共通処理エラー
type BeginError struct {
	Message string
}

func (err *BeginError) Error() string {
	return err.Message
}

// CommitError : 事後共通処理エラー
type CommitError struct {
	Message string
}

func (err *CommitError) Error() string {
	return err.Message
}

// AccessDenied : IP制限に引っかかった場合のエラー
type AccessDenied struct {
	Message    string
	StatusCode int
	IP         string
}

func (err *AccessDenied) Error() string {
	return err.Message
}

// Unauthorized : 401 Unauthorized エラー
type Unauthorized struct {
	Message    string
	Title      string
	StatusCode int
	Data       interface{}
}

func (err *Unauthorized) Error() string {
	return err.Message
}

func (err *Unauthorized) pointer() {}
