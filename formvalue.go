package mux

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// FormValues 構造体は、GET/POST の入力データを管理する
type FormValues map[string][]string

// Get gets the first value associated with the given key. If there are no values associated with the key, Get returns the empty string. To access multiple values, use the map directly.
func (v FormValues) Get(key string) string {
	return url.Values(v).Get(key)
}

// Del deletes the values associated with key.
func (v FormValues) Del(key string) {
	url.Values(v).Del(key)
}

// Add adds the value to key. It appends to any existing values associated with key.
func (v FormValues) Add(key, value string) {
	url.Values(v).Add(key, value)
}

// Encode encodes the values into “URL encoded” form ("bar=baz&foo=quux") sorted by key.
func (v FormValues) Encode() string {
	return url.Values(v).Encode()
}

// Set sets the key to value. It replaces any existing values.
func (v FormValues) Set(key, value string) {
	url.Values(v).Set(key, value)
}

// Values : FormValues => url.Values へ変換する
func (v FormValues) Values() url.Values {
	return url.Values(v)
}

// Copy FormValues の値を構造体へ格納する
func (v FormValues) Copy(i interface{}) error {
	val := reflect.ValueOf(i)

	// 無効なデータが渡された場合は、エラーを返却する
	if val.IsValid() == false {
		return fmt.Errorf("Invalid data")
	}
	// ポインタ型ではない場合、エラーを返却する
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("%s is not a pointer", val.Type())
	}
	// ポインタが nil の場合、エラーを返却する
	if val.IsNil() {
		return fmt.Errorf("%s is nil pointer", val.Type())
	}
	// ポインタの中身が構造体ではない場合、エラーを返却する
	if val = val.Elem(); val.Kind() != reflect.Struct {
		return fmt.Errorf("%s is not a struct", val.Type())
	}

	// 該当するデータを構造体が所持する全フィールドに適用する
	for i := 0; i < val.NumField(); i++ {
		// フィールドを取得
		f := val.Type().Field(i)
		// フィールド名、またはフィールドに付与しているタグ名から該当するデータを取得する
		value, ok := v[f.Name]
		if !ok {
			value, ok = v[strings.ToLower(f.Name)]
			if !ok {
				if tag := f.Tag.Get("json"); tag != "" {
					value, ok = v[tag]
				}
				if !ok {
					continue
				}
			}
		}

		if f.Type.Kind() == reflect.Slice {
			// フィールドの型情報がスライスの場合、要素の型情報と、値を定義する
			rt := reflect.MakeSlice(f.Type, 1, 1).Index(0).Type()
			rv := reflect.MakeSlice(f.Type, 0, 0)
			for _, str := range value {
				v, err := setvalue(rt, str)
				if err != nil {
					return err
				}
				rv = reflect.Append(rv, v)
			}
			val.Field(i).Set(rv)
		} else {
			// フィールドの型情報がスライス以外の場合
			v, err := setvalue(f.Type, value[0])
			if err != nil {
				return err
			}
			val.Field(i).Set(v)
		}
	}
	return nil
}

func setvalue(rt reflect.Type, value string) (rv reflect.Value, err error) {
	var n interface{}
	rv = reflect.New(rt).Elem()
	switch rt.Kind() {
	case reflect.Int:
		if n, err = strconv.ParseInt(value, 10, 32); err != nil {
			return
		}
		rv.SetInt(n.(int64))
	case reflect.Int8:
		if n, err = strconv.ParseInt(value, 10, 8); err != nil {
			return
		}
		rv.SetInt(n.(int64))
	case reflect.Int16:
		if n, err = strconv.ParseInt(value, 10, 16); err != nil {
			return
		}
		rv.SetInt(n.(int64))
	case reflect.Int32:
		if n, err = strconv.ParseInt(value, 10, 32); err != nil {
			return
		}
		rv.SetInt(n.(int64))
	case reflect.Int64:
		if n, err = strconv.ParseInt(value, 10, 64); err != nil {
			return
		}
		rv.SetInt(n.(int64))
	case reflect.Uint:
		if n, err = strconv.ParseUint(value, 10, 32); err != nil {
			return
		}
		rv.SetUint(n.(uint64))
	case reflect.Uint8:
		if n, err = strconv.ParseUint(value, 10, 8); err != nil {
			return
		}
		rv.SetUint(n.(uint64))
	case reflect.Uint16:
		if n, err = strconv.ParseUint(value, 10, 16); err != nil {
			return
		}
		rv.SetUint(n.(uint64))
	case reflect.Uint32:
		if n, err = strconv.ParseUint(value, 10, 32); err != nil {
			return
		}
		rv.SetUint(n.(uint64))
	case reflect.Uint64:
		if n, err = strconv.ParseUint(value, 10, 64); err != nil {
			return
		}
		rv.SetUint(n.(uint64))
	case reflect.Float32:
		if n, err = strconv.ParseFloat(value, 32); err != nil {
			return
		}
		rv.SetFloat(n.(float64))
	case reflect.Float64:
		if n, err = strconv.ParseFloat(value, 32); err != nil {
			return
		}
		rv.SetFloat(n.(float64))
	case reflect.Bool:
		if n, err = strconv.ParseBool(value); err != nil {
			return
		}
		rv.SetBool(n.(bool))
	case reflect.String:
		rv.SetString(value)
	case reflect.Struct:
		if rt.Name() == "time.Time" {
			datetime, err := time.Parse("2006-01-02 15:04:05", value)
			if err != nil {
				return rv, err
			}
			rv.Set(reflect.ValueOf(datetime))
		} else {
			err = fmt.Errorf("Unusable type")
		}
	default:
		err = fmt.Errorf("Unusable type")
	}
	return
}

// Map : map[string][]string => map[string]interface へ変換する
func (v FormValues) Map() map[string]interface{} {
	var result = make(map[string]interface{})
	for key, value := range v {
		if len(value) == 1 {
			result[key] = value[0]
		} else {
			result[key] = value
		}
	}
	return result
}
