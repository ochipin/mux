package helpers

import "net/url"

// URL : URL構造体情報を取り扱う構造体
type URL struct {
	*url.URL
}

// Host : ホスト:ポート番号を返却する
func (u *URL) Host() StringType {
	return StringType(u.URL.Host)
}

// Hostname : ホスト名を返却する
func (u *URL) Hostname() StringType {
	return StringType(u.URL.Hostname())
}

// Port : ポート番号を返却する
func (u *URL) Port() StringType {
	return StringType(u.URL.Port())
}

// Path : クエリパスを返却する
func (u *URL) Path() StringType {
	return StringType(u.URL.Path)
}

// Proto : http[s] を返却する
func (u *URL) Proto() StringType {
	return StringType(u.URL.Scheme)
}

// Search : クエリパラメータのみを返却する
func (u *URL) Search() StringType {
	return StringType(u.URL.RawQuery)
}

// Query : クエリパラメータを取得する
func (u *URL) Query(i ...string) interface{} {
	// クエリパラメータがない場合、関数を抜ける
	if u.URL.Query() == nil {
		return StringType("")
	}

	var result interface{}

	if len(i) == 0 {
		// 引数がない場合、map[string]StringType型にクエリパラメータを格納する
		var m = make(map[string]StringType)
		for k, v := range u.URL.Query() {
			m[k] = StringType(v[0])
		}
		result = m
	} else {
		// 引数がある場合、指定されたクエリパラメータの情報のみを格納する
		var m = make(map[string]StringType)
		for _, v := range i {
			m[v] = StringType(u.URL.Query().Get(v))
		}
		result = m
	}

	return result
}

// Get : 単体のクエリパラメータの情報のみを取得する
func (u *URL) Get(name string) interface{} {
	return StringType(u.URL.Query().Get(name))
}

// Origin : プロトコル名,ホスト名,ポート番号を付与したURLを返却する
func (u *URL) Origin() StringType {
	return StringType(u.URL.Scheme + "://" + u.URL.Host)
}
