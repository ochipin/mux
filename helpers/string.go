package helpers

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"sort"
	"strings"
)

// StringType : 文字列型
type StringType string

func (s StringType) String() string {
	return string(s)
}

// Index : 指定した文字が出現する箇所のインデックスを返却する。一致しない場合は、-1を返却する
func (s StringType) Index(search string) int {
	return strings.Index(s.String(), search)
}

// Count : 指定した文字が含まれている数を返却する
func (s StringType) Count(search string) int {
	return strings.Count(s.String(), search)
}

// Len : []runeでの文字列の長さを取得する
func (s StringType) Len() int {
	return len([]rune(s.String()))
}

// Match : 指定した文字列と一致した場合、trueを返却する
func (s StringType) Match(i interface{}) (bool, error) {
	var str string
	switch types := i.(type) {
	case string:
		str = types
	case StringType:
		str = types.String()
	default:
		return false, fmt.Errorf("Match invalid arguments")
	}
	r, err := regexp.Compile(str)
	if err != nil {
		return false, err
	}
	return r.MatchString(s.String()), nil
}

// Lower : アルファベットをすべて小文字にする
func (s StringType) Lower() StringType {
	return StringType(strings.ToLower(string(s)))
}

// Upper : アルファベットをすべて大文字にする
func (s StringType) Upper() StringType {
	return StringType(strings.ToUpper(string(s)))
}

// Title : アルファベットの先頭文字のみを大文字にする
func (s StringType) Title() StringType {
	return StringType(strings.Title(string(s)))
}

// Strip : 前後にある空白、改行コードを除去する
func (s StringType) Strip() StringType {
	return StringType(strings.TrimSpace(s.String()))
}

// Trim : 前後にある、指定した文字列を除去する
func (s StringType) Trim(cutset string) StringType {
	return StringType(strings.Trim(s.String(), cutset))
}

// TrimLeft : 左側にある、指定した文字列を除去する
func (s StringType) TrimLeft(cutset string) StringType {
	return StringType(strings.TrimLeft(s.String(), cutset))
}

// TrimRight : 左側にある、指定した文字列を除去する
func (s StringType) TrimRight(cutset string) StringType {
	return StringType(strings.TrimRight(s.String(), cutset))
}

// Template : 指定されたパラメータを元に、文字列内に埋め込まれた変数を展開する
func (s StringType) Template(i interface{}) (StringType, error) {
	tmpl, err := template.New("String.Template").Parse(s.String())
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, i); err != nil {
		return "", err
	}
	return StringType(buf.String()), nil
}

// Slice : 文字列の切り出しを行う
func (s StringType) Slice(i ...int) (StringType, error) {
	// 引数が1つ、または2つではない場合、エラーとする
	if len(i) != 1 && len(i) != 2 {
		return "", fmt.Errorf("slice: invalid arguments error")
	}

	// 第1引数の値を検証する
	start := i[0]
	if start < 0 {
		// 負の値を指定されている場合、後尾から文字列を切り出す
		start = len(s.String()) - start
		// 指定された値が、配列の要素数を超過してしまう場合、強制的に "0" 番目とする
		if start < 0 {
			start = 0
		}
	} else {
		// 正の値を指定されている場合、先頭から文字列を切り出す
		if start >= len(s.String()) {
			start = len(s.String())
		}
	}

	end := len(s.String())
	// 第2引数が指定されている場合
	if len(i) == 2 {
		// 与えられた引数の数が配列の要素数内であるかチェックする
		end = i[1]
		if end <= 0 {
			end = start
		} else if (end + start) >= len(s.String()) {
			end = len(s.String())
		} else {
			end += start
		}
	}

	return StringType(s.String()[start:end]), nil
}

// Replace : 文字列を置き換える
func (s StringType) Replace(old, new string) StringType {
	return StringType(strings.Replace(s.String(), old, new, -1))
}

// Split : 指定した区切り文字で、文字列を配列にする
func (s StringType) Split(sep string) Strings {
	return strings.Fields(s.Replace(sep, " ").String())
}

// Strings : []string 型
type Strings []string

// Sort : 配列を昇順でソートする
func (s Strings) Sort() Strings {
	sort.Strings(s)
	return s
}

// Reverse : 配列を逆順にする
func (s Strings) Reverse() Strings {
	if len(s) == 0 {
		return s
	}
	var result = make(Strings, len(s))
	for i := 0; i < len(s); i++ {
		result[i] = s[len(s)-(i+1)]
	}
	return result
}

// Uniq : 配列内の重複した文字列を削除する
func (s Strings) Uniq() Strings {
	if len(s) == 0 {
		return s
	}
	var result []string
	var dup = make(map[string]string)
	for _, v := range s {
		if _, ok := dup[v]; !ok {
			result = append(result, v)
		}
	}
	return result
}

// Join : 文字列を連結する
func (s Strings) Join(sep string) StringType {
	return StringType(strings.Join(s, sep))
}
