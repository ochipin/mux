package helpers

import (
	"fmt"
	"sort"
	"strings"
)

// AutoComplete : autocomplete='off' を生成する型
type AutoComplete string

// NoValidate : novalidate='novalidate' を生成する型
type NoValidate string

// Multipart : enctype='multipart/form-data' を生成する型
type Multipart string

// Form : <form>タグを生成する構造体
type Form struct {
	data    map[string]interface{}
	method  string
	keyname string
	end     bool
}

func (form *Form) String() string {
	var result string

	if form.end == false {
		// 開始タグ生成
		if strings.ToUpper(form.method) == "GET" {
			result += "<form method='GET'"
		} else {
			result += "<form method='POST'"
		}
		keys := make([]string, 0, len(form.data))
		for k := range form.data {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			result += fmt.Sprintf(" %s='%v'", k, form.data[k])
		}
	} else {
		// 閉じタグ生成
		if strings.ToUpper(form.method) != "GET" {
			result += fmt.Sprintf("<input type='hidden' name='%s' value='%s' />", form.keyname, form.method)
		}
		result += "</form"
	}

	return strings.TrimSpace(result) + ">"
}

// Attr : 属性値を設定する
func (form *Form) Attr(name string, value interface{}) string {
	form.data[name] = value
	return ""
}
