package helpers

import "fmt"

// Parameters : テンプレートファイル内で扱うマップ型
type Parameters map[string]interface{}

// Set : マップにデータを追加、または書き換える
func (params Parameters) Set(key string, val interface{}) string {
	params[key] = StringType(fmt.Sprint(val))
	return ""
}

// Delete : マップ要素にあるデータを削除する
func (params Parameters) Delete(key string) string {
	if _, ok := params[key]; ok {
		delete(params, key)
	}
	return ""
}

// Clear : マップ要素内のデータをすべて削除する
func (params Parameters) Clear() string {
	for name, _ := range params {
		delete(params, name)
	}
	return ""
}

// Copy : マップ要素内のデータをすべて別のマップへコピーする
func (params Parameters) Copy() Parameters {
	var result = make(Parameters)
	for key, value := range params {
		result[key] = value
	}
	return result
}

// HasItem : 指定した名前でデータが登録されているか否かを検出する
func (params Parameters) HasItem(name string) bool {
	_, ok := params[name]
	return ok
}

// T : 指定した名前をキーに、データを取得する
func (params Parameters) T(name string) StringType {
	if v, ok := params[name]; ok {
		return StringType(fmt.Sprint(v))
	}
	return ""
}
