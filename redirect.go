package mux

// Redirect : 301, 302 リダイレクトを実施する構造体
type Redirect struct {
	path       string // リダイレクト先のパス
	statuscode int    // 301, 302 などのステータスコード
}

// Perm : 301 リダイレクト
func (r *Redirect) Perm() *Redirect {
	r.statuscode = 301
	return r
}

// Temp : 302 リダイレクト
func (r *Redirect) Temp() *Redirect {
	r.statuscode = 302
	return r
}

func (r *Redirect) pointer() {}
