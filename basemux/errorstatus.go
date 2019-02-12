package basemux

import (
	"fmt"
	"runtime"
)

// PanicError : 発生時のエラー型
type PanicError struct {
	StackTrace []string // スタックトレース
	Title      string   // エラータイトル
	Message    string   // エラーメッセージ
	StatusCode int      // ステータスコード
}

func (p *PanicError) Error() string {
	return p.Message
}

// PanicDump : PANIC発生時にコールし、トレース情報を収集した
func PanicDump(point int, err interface{}) error {
	var p = &PanicError{}
	// スタックトレースの取得
	for i := point; ; i++ {
		pc, filename, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		funcname := runtime.FuncForPC(pc).Name()
		p.StackTrace = append(p.StackTrace, fmt.Sprintf("=======>> %d: %s: %s(%d)", i-point, funcname, filename, line))
	}
	// エラータイトルを格納
	p.Title = "500 Internal Server Error"
	// エラーメッセージを格納
	p.Message = fmt.Sprintf("runtime error. %v", err)
	// ステータスコードを格納
	p.StatusCode = 500

	return p
}

// TimeoutError : タイムアウト発生時のエラー型
type TimeoutError struct {
	Title      string // エラータイトル
	Message    string // エラーメッセージ
	StatusCode int    // ステータスコード
}

func (p *TimeoutError) Error() string {
	return p.Message
}

// MaxClientsError : 最大同時リクエスト数を超過した際のエラー型
type MaxClientsError struct {
	Title      string // エラータイトル
	Message    string // エラーメッセージ
	StatusCode int    // ステータスコード
}

func (p *MaxClientsError) Error() string {
	return p.Message
}

// タイムアウトエラー、または最大同時リクエスト数が超過した場合コールする
func accessError(status int) error {
	var err error
	if status == 408 {
		err = &TimeoutError{
			Title:      "408 Request Time-out",
			Message:    fmt.Sprintf("request time-out"),
			StatusCode: status,
		}
	} else if status == 503 {
		err = &MaxClientsError{
			Title:      "503 Service Temporarily Unavailable",
			Message:    fmt.Sprintf("max clients number of limit exceeded"),
			StatusCode: status,
		}
	}
	return err
}
