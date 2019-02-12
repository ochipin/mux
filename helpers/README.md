mux.Helper ヘルパ構造体
===
テンプレートファイル内で使用できる関数について説明します。  
muxライブラリでは、テンプレートファイル内で次の関数をコールできます。

## add (a, b int) int
`add` 関数は、2つの`int`型の引数を受け取ります。引数に渡された `a` と `b` を足した値を返却します。
```go
{{add 1 2}} {{/* 1 + 2 を実行し、 3 を出力する */}}
```

## min (a, b int) int
`min` 関数は、2つの`int`型の引数を受け取ります。引数に渡された `a` から `b` を引いた値を返却します。
```go
{{min 1 2}} {{/* 1 - 2 を実行し、 -1 を出力する */}}
{{min 3 2}} {{/* 3 - 2 を実行し、 1 を出力する */}}
```

## mul (a, b int) int
`mul` 関数は、2つの`int`型の引数を受け取ります。引数に渡された `a` と `b` を掛けた値を返却します。
```go
{{mul 1 2}} {{/* 1 * 2 を実行し、 2 を出力する */}}
{{mul 3 2}} {{/* 3 * 2 を実行し、 6 を出力する */}}
```

## dev (a, b int) int
`dev` 関数は、2つの`int`型の引数を受け取ります。引数に渡された `a` と `b` を割った値を返却します。
```go
{{dev 2 2}} {{/* 2 / 2 を実行し、 1 を出力する */}}
{{dev 6 2}} {{/* 6 / 2 を実行し、 3 を出力する */}}
{{dev 3 2}} {{/* 3 / 2 を実行し、 1 を出力する */}}
```

## mod (a, b int) int
`mod` 関数は、2つの`int`型の引数を受け取ります。引数に渡された `a` と `b` を割った余りを返却します。
```go
{{mod 2 2}} {{/* 2 % 2 を実行し、 0 を出力する */}}
{{mod 6 2}} {{/* 6 % 2 を実行し、 0 を出力する */}}
{{mod 3 2}} {{/* 3 % 2 を実行し、 1 を出力する */}}
{{mod 5 3}} {{/* 5 % 3 を実行し、 2 を出力する */}}
{{mod 5 2}} {{/* 5 % 2 を実行し、 1 を出力する */}}
```

## sprintf (i ...interface{}) StringType
与えられた引数を`StringType`型に変換します。

```go
{{$v := sprintf "%d/%s" 200 "hello"}}
{{$v}} {{/* '200/hello' を表示する */}}
```

第1引数には、フォーマットを指定することができます。上記例の場合、第2引数以降が文字列へ変換される値となります。  
引数が1つの場合、その引数を`StringType`型へ変換した値を返却します。

```go
{{$v := sprintf "HELLO WORLD"}}
{{$v}} {{/* 'HELLO WORLD' を表示 */}}
```

`StringType`型は、次のメソッドを所持しています。

### Index
指定した文字が出現する箇所のインデックスを返却する。一致しない場合は、-1を返却する。

```go
{{$v := sprintf "HELLO WORLD"}}

{{$v.Index "W"}} {{/* 6 */}}
```

### Count
指定した文字が含まれている数を返却する。

```go
{{$v := sprintf "HELLO WORLD"}}

{{$v.Count "L"}} {{/* 3 */}}
```

### Len
`[]rune`での文字列の長さを取得する

```go
{{$v := sprintf "ハロー WORLD"}}

{{$v.Len}} {{/* 9 */}}
{{/* 'ハ', 'ロ', 'ー', ' ', 'W', 'O', 'R', 'L', 'D' */}}
{{/*  1     2     3    4    5    6    7    8    9 */}}
```

### Match (i interface{}) (bool, error)
正規表現文字列 `i` が文字列と一致するか調べます。
```go
{{$v := sprintf "HELLO WORLD"}}

{{$v.Match "HELLO"}}         {{/* true */}}
{{$v.Match "^HELLO$"}}       {{/* false */}}
{{$v.Match "HELLO\s*WORLD"}} {{/* true */}}
```

### Lower () StringType
文字列をUnicodeの小文字にマッピングした値を返却します。
```go
{{$v := sprintf "HELLO WORLD"}}

{{$v.Lower}} {{/* HELLO WORLD --> hello world */}}
```

### Upper () StringType
文字列をUnicodeの大文字にマッピングした値を返却します。
```go
{{$v := sprintf "hello world"}}

{{$v.Lower}} {{/* hello world --> HELLO WORLD */}}
```

### Title () StringType
文字列の先頭文字をUnicodeの大文字にマッピングした値を返却します。
```go
{{$v := sprintf "hello"}}

{{$v.Lower}} {{/* hello --> Hello */}}
```

### Strip () StringType
文字列の前後にある空白、改行コードを除去します。
```go
{{$v := sprintf "   HELLO WORLD   "}}
{{$v.Strip}} {{/* 'HELLO WOLRD' を表示 */}}
```

### Trim (cutset string) StringType
`cutset`に含まれるUnicodeコードポイントを文字列の先頭と末尾からすべて削除します。
```go
{{$v := sprintf "***HELLO WORLD***"}}
{{$v.Trim "*"}} {{/* 'HELLO WOLRD' を表示 */}}
```

### TrimLeft (cutset string) StringType
`cutset`に含まれるUnicodeコードポイントを文字列の先頭からすべて削除します。
```go
{{$v := sprintf "***HELLO WORLD***"}}
{{$v.TrimLeft "*"}} {{/* 'HELLO WOLRD***' を表示 */}}
```

### TrimRight (cutset string) StringType
`cutset`に含まれるUnicodeコードポイントを文字列の末尾からすべて削除します。
```go
{{$v := sprintf "***HELLO WORLD***"}}
{{$v.TrimRight "*"}} {{/* '***HELLO WOLRD' を表示 */}}
```

### Template (i interface{}) (StringType, error)
文字列テンプレートを解析した結果を返却します。
```go
{{$v := sprintf "hello {{.value}}"}}

{{$p := makemap}}
{{$p.Set "value" "world"}}

{{$v.Template $p}} {{/* hello {{.value}} --> hello world */}}
```

### Slice(i ...int) (StringType, error)
文字列を、指定した範囲の部分のみ切り出す。
引数`i`に「開始位置」と「終了位置」を指定。終了位置は省略可能。省略すると最後まで切り出す。
開始位置をマイナス値にすると、後からの桁数になる(右端のみ切り出せる)。

```go
{{$v := sprintf "0123456789"}}

{{$v.Slice 1}}     {{/* 123456789 */}}
{{$v.Slice -1}}    {{/* 9 */}}
{{$v.Slice 4}}     {{/* 456789 */}}
{{$v.Slice 4 1}}   {{/* 4 */}}
{{$v.Slice 5 100}} {{/* 56789 */}}
{{$v.Slice 5 3}}   {{/* 567 */}}
```

### Replace (old, new string) StringType
`old`の部分を`new`に置換します。
```go
{{$v := sprintf "hello <?> hello <?>"}}

{{$v.Replace "<?>" "world"}} {{/* hello <?> hello <?> --> hello world hello world */}}
```

### Split (sep string) Strings
引数`sep`に指定した区切り文字で、文字列配列にします。

```go
{{$v := sprintf "name1, name2, name3"}}

{{/* $v.Split "," --> Strings{"name1", "name2", "name3"} */}}
{{range $v.Split ","}}
  {{.}}
{{end}}
```
`Split`が返却する文字列配列は、`Strings`型となっており、次の4つの関数をコール可能です。

| Method                        | Description |
|:--                            |:-- |
| `Sort() Strings`              | 文字列配列を昇順でソートする |
| `Reverse() Strings`           | 文字列配列を逆順にする |
| `Uniq() Strings`              | 重複した文字列を取り除く |
| `Join(sep string) StringType` | 文字列配列を連結する |

次のように使用可能です。

```go
{{$str := sprintf "name3, name1, name2, name1"}}

{{/* $str.Split "," --> Strings{"name3", "name1", "name2", "name1"} */}}
{{$v := $str.Split ","}}

{{/* sort: ソートする --> Strings{"name1", "name1", "name2", "name3"} */}}
{{$v.Sort}}

{{/* reverse: 逆順にする --> Strings{"name3", "name2", "name1", "name1"} */}}
{{$v.Sort.Reverse}}

{{/* uniq: 重複を取り除く --> Strings{"name3", "name1", "name2"} */}}
{{$v.Uniq}}
{{/* uniq: 重複を取り除き、逆順ソートする --> Strings{"name3", "name2", "name1"} */}}
{{$v.Uniq.Sort.Reverse}}

{{/* join: 文字列を連結する --> "name3_name2_name1" */}}
{{$v.Uniq.Sort.Reverse.Join "_"}}
```

## makemap() Parameters
テンプレートファイル内で使用するマップを生成します。

```go
{{$p := makemap}} {{/* map[string]StringType */}}
```
生成された、`map`は、次のメソッドを所持しています。

|Method   | Description |
|:--      |:--     |
|`Set(key string, val interface{})`    | キーと値でデータを登録する |
|`Delete(key string)` | 不要なデータを、キー名を指定することで削除する |
|`Clear()`  | 全データを削除する |
|`Copy() Parameters`   | 全データをコピーする |
|`HasItem(name string) bool`| データが存在するかチェックする |

### 使用例
```go
{{/* Map を生成する */}}
{{$p := makemap}}

{{/* Set は、キーと値でデータを登録する */}}
{{$p.Set "key" "value"}} {{/* {"key": "value"} */}}
{{$p.Set "name" "app"}}  {{/* {"key": "value", "name": "app"} */}}

{{/* Set は、同名のキー名を指定すると、値を上書きする */}}
{{$p.Set "key" "val"}}   {{/* {"key": "val", "name": "app"} */}}

{{/* キー名を指定することで、指定したキー名で値を取得可能 */}}
{{$p.key}} {{/* val を出力 */}}

{{/* Copy でキーと値を複製 */}}
{{$p2 := $p.Copy}}

{{/* Delete にキー名を指定することで、登録されたデータを削除する */}}
{{$p.Delete "name"}}     {{/* {"key": "val"} */}}
{{/* Clear は全データを削除する */}}
{{$p.Clear}}             {{/* {} */}}

{{/* $p から複製されたデータ $p2 は、別データとして扱われるため、 */}}
{{/* 元データである $p のデータを削除、または追加を行っても $p2 は影響を受けない */}}
{{range $k, $v := $p2}}
  {{printf "%s => %s" $k $v}} {{/* key => val, name => app を出力 */}}
{{end}}

{{/* 指定したキー名でデータが登録されているかチェックする */}}
{{$p2.HasItem "name"}}    {{/* true */}}
{{$p2.HasItem "unknown"}} {{/* false */}}
```

## date () (*DateTime, error)
`date`関数は、現在時刻を返却します。`date`関数の第1引数には、次のフォーマット指定子を指定して時刻の出力方式を変更することができます。

| Format | Usage           | Result     | Description |
|:--     |:--              |:--         |:--     |
| `%A`   | `{{date "%A"}}` | `Thursday` | 曜日の名称(Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday)|
| `%a`   | `{{date "%a"}}` | `Thu`      | 曜日の省略名(Sun, Mon, Tue, Wed, Thu, Fri, Sat) |
| `%B`   | `{{date "%B"}}` | `December` | 月の名称(January, February, March, April, May, June, July, August, September, October, November, December) |
| `%b`   | `{{date "%b"}}` | `Dec`      | 月の省略名(Jan, Feb, Mar, Aprm May, Jun, Jul, Aug, Sep, Oct, Nov, Dec) |
| `%c`   | `{{date "%c"}}` | `Thu Dec 20 18:27:49 2018` | 日付と時刻 |
| `%d`   | `{{date "%d"}}` | `20`       | 日(01-31) |
| `%H`   | `{{date "%H"}}` | `18`       | 24時間制の時(00-23) |
| `%I`   | `{{date "%I"}}` | `06`       | 12時間制の時(01-12) |
| `%j`   | `{{date "%j"}}` | `354`      | 年中の通算日(001-366) |
| `%M`   | `{{date "%M"}}` | `27`       | 分(00-59) |
| `%m`   | `{{date "%m"}}` | `12`       | 月を表す数字(01-12) |
| `%P`   | `{{date "%P"}}` | `pm`       | 午前または午後(am,pm) |
| `%p`   | `{{date "%p"}}` | `PM`       | 午前または午後(AM,PM) |
| `%S`   | `{{date "%S"}}` | `49`       | 秒(00-59) |
| `%U`   | `{{date "%U"}}` | `51`       | 週を表す数。最初の日曜日が第1週の始まり(00-53) |
| `%W`   | `{{date "%W"}}` | `50`       | 週を表す数。最初の月曜日が第1週の始まり(00-53) |
| `%w`   | `{{date "%w"}}` | `4`        | 曜日を表す数。日曜日が0(0-6) |
| `%X`   | `{{date "%X"}}` | `18:27:49` | `HH:MM:SS`形式の時刻 |
| `%x`   | `{{date "%x"}}` | `12/20/18` | `mm/dd/yy`形式の日付 |
| `%Y`   | `{{date "%Y"}}` | `2018`     | `YYYY`形式の西暦     |
| `%y`   | `{{date "%y"}}` | `18`       | `yy`形式の西暦       |
| `%Z`   | `{{date "%Z"}}` | `JST`      | タイムゾーン         |

```go
{{date}}      {{/* Thu Dec 20 18:27:49 JST 2018 */}}
{{date "%c"}} {{/* Thu Dec 20 18:27:49 2018 */}}

{{date "%Y/%m/%d %H:%M:%S"}} {{/* 2018/12/20 18:27:49 */}}
```
また、第2引数に次の単位を渡すことにより、時刻を調整することができます。

```go
{{date "..." "<number> <unit>"}}
```

指定する数字の先頭に `-` を用いることで、日時を進めるだけでなく、戻ることもできます。

| Unit   | Usage            | Result | Description |
|:--     |:--              |:--         |:--     |
| 秒     | `{{date "..." "120 second"}}` | `18:27:49 --> 18:29:49` | 秒単位で時刻を調整する|
| 分     | `{{date "..." "3 minute"}}` | `18:27:49 --> 18:30:49` | 分単位で時刻を調整する |
| 時     | `{{date "..." "3 hour"}}` | `18:27:49 --> 21:27:49` | 時単位で時刻を調整する |
| 日     |`{{date "..." "-18 day"}}` | `2018-12-20 --> 2018-12-02` | 日単位で時刻を調整する |
| 月     |`{{date "..." "-3 month"}}` | `2018-12-20 --> 2018-09-20` | 月単位で時刻を調整する |
| 年     |`{{date "..." "3 year"}}` | `2018-12-20 --> 2021-12-20` | 年単位で時刻を調整する |

```go
{{date "%c" "18 day"}}   {{/* 現在時刻に18日足した日付を表示する */}}
{{date "%c" "-1 month"}} {{/* 1ヶ月前の日付を表示する */}}
```
日時をより細かく扱えるよう、次のメソッドを使用できます。

| Method    | Type           | Description |
|:--        |:--             |:--     |
| `Second`  | `int`          |「秒(0-59)」を返却する |
| `Minute`  | `int`          |「分(0-59)」を返却する |
| `Hour`    | `int`          |「時(0-23)」を返却する |
| `Day`     | `int`          |「日(0-31)」を返却する |
| `Month`   | `time.Month`   |「月(1-12)」を返却する |
| `Year`    | `int`          |「年(YYYY)」を返却する |
| `YearDay` | `int`          | 年中の通算日(001-366)を返却する |
| `Weekday` | `time.Weekday` | 曜日を表す数を返却する。日曜日が0(0-6) |
| `Time`    | `time.Time`    | `time.Time`型へ変換する |
| `Format`  | `string`       | フォーマットを指定して日時を表示する |

### 使用例

```go
{{/* 日時情報を取得 */}}
{{$v := date}}

{{$v.Second}}  {{/* 49 */}}
{{$v.Minute}}  {{/* 27 */}}
{{$v.Hour}}    {{/* 18 */}}
{{$v.Day}}     {{/* 20 */}}
{{$v.Month}}   {{/* December */}}
{{(sprintf "%s" $v.Month).Slice 0 3}} {{/* Dec */}}
{{printf "%d" $v.Month}} {{/* 12 */}}
{{$v.Year}}    {{/* 2018 */}}
{{$v.YearDay}} {{/* 354 */}}
{{$v.Weekday}} {{/* Thursday */}}
{{printf "%d" $v.Weekday}} {{/* 4 */}}
{{(sprintf "%s" $v.Weekday).Slice 0 3}} {{/* Thu */}}
```

## form (params ...interface{}) (Form, error)
formタグを生成します。
```go
{{/* <form method='POST' action='/path/to/url'>...</form> */}}
{{form}}
  ...
{{form}} {{/* <-- 閉じタグ(</form>)を生成 */}}
```
次のような、formタグがネストする構造を生成することはできません。
```go
{{form}}
  ...
  {{form}}
    ...
  {{form}}
  ...
{{form}}
```
form関数の引数に何も指定しない場合、method,action属性には次の値が付与されます。

| Attribute | Values |
|:--        |:--     |
| method | POST が設定される |
| action | アクセスしたクエリパスが設定される。`https://localhost/path/to/url` の場合、`/path/to/url` が設定される |

また、form関数の引数を使用することで、次の属性値を設定することができます。
| Attribute        | Description |
|:--               |:--     |
| `id`             | 先頭文字が`'#'`から始まる文字列の場合、`id`属性値としてみなす |
| `class`          | 先頭文字が`'.'`から始まる文字列の場合、`class`属性値としてみなす |
| `autocomplete`   | `nocomplete` を指定した場合、`autocomplete='off'`を生成する |
| `novalidate`     | `novalidate` を指定した場合、`novalidate='novalidate'`を生成する |
| `enctype`        | `multipart` を指定した場合、`enctype='multipart/form-data'`を生成する |
| `action`         | 先頭文字が`'/'`から始まる文字列の場合、`action`属性値としてみなす |
| `target`         | 先頭文字が`':'`から始まる文字列の場合、`target`属性値としてみなす |
| `accept-charset` | 先頭文字が`'@'`から始まる文字列の場合、`accept-charset`属性値としてみなす |
| `method`         | 先頭文字が`'$'`から始まる文字列の場合、`method`属性値としてみなす |
| `name`           | 特に何も指定しない場合、`name`属性値としてみなす |

### 使用例
```go
{{/* <form action='/path/to/change' name='formname' id='form_id' class='form_class' */}}
{{/*       autocomplete='off' novalidate='novalidate' enctype='multipart/form-data' */}}
{{/*       target='_blank' accept-charset='UTF-8' method='POST'> */}}
{{form "/path/to/change" "formname" "#form_id" ".form_class" nocomplete novalidate multipart ":_blank" "@UTF-8"}}
  ...
{{/* </form> */}}
{{form}}
```

引数では指定できない属性値を設定するには、`Attr`関数を使用することで設定可能です。

```go
{{/* <form action='/path/to/change' name='formname' ... data-text='data-value'> */}}
{{$f := form "/path/to/change" "formname" ...}}
{{$f.Attr "data-text" "data-value"}}{{$f}}
  ...
{{/*   <input type='hidden' name='_method' value='PUT' /> */}}
{{/* </form> */}}
{{form}}
```

method 属性を、GET/POST以外にする場合は、次のように指定します。

```go
{{/* <form action='...' ... method='GET'>...</form> */}}
{{form ... "$GET"}}...{{form}}

{{/* <form action='...' ... method='POST'> */}}
{{/*   ... */}}
{{/*   <input type='hidden' name='_method' value='DELETE' /> */}}
{{/* </form> */}}
{{form ... "$DELETE"}}...{{form}}
```

## controller() StringType
コントローラ名を返却します。

## action() StringType
アクション名を返却します。

## id() StringType
ID属性値を返却します。

## class() StringType
Class属性値を返却します。

## charset() StringType
charset設定値(UTF-8)を返却します。

## lang() StringType
適用されている自然言語名を返却します。返却される値は、`Accept-Language`の値によって異なります。

```go
{{/* Accept-Language が "ar-DZ,ar-JO;q=0.8,id;q=0.6,ug;q=0.4,ky;q=0.2" だった場合、 ja となる */}}
{{/* 理由は、Accept-Languageには、 ja, en 双方含まれていないため、Defaultで設定した ja が選定される */}}
{{lang}} {{/* ja */}}

{{/* Accept-Language が "ar-DZ,zh;q=0.8,ja;q=0.6,en-US;q=0.4,en;q=0.2" だった場合、 ja となる */}}
{{/* 理由は、Accept-Languageの解析は左から右に対して行うため、ja が en よりも先に出現するため */}}
{{lang}} {{/* ja */}}

{{/* Accept-Language が "ar-DZ,zh;q=0.8,en-US;q=0.4,en;q=0.2" だった場合、 en となる */}}
{{lang}} {{/* en */}}
```

返却される ja や en などの文字列は、`mux.Mux.Locale`で次のような設定をすることで追加、変更可能です。
```go
lang := &locale.Locale{
    /* Accept-Language の判定に使用する */
    Langs: map[string][]string{
        "ja": []string{"ja"},
        "en": []string{"en", "en-US", "en-*"},
    },
    /* Accept-Language 内に, ja, en 系の言語が存在しない場合、デフォルトとして ja を使用する */
    Default: "ja",
}
/* Mux 構造体に言語構造体を設定 */
l, _ := lang.CreateLocale()
&mux.Mux{
    Locale: l,
}
```

## title() StringType
`<title>`タグに設定するタイトル名を返却します。

## set(key string, value interface{}) string
キーと値で、変数を管理します。`set`関数で設定した値は、次の関数に影響します。

* controller
* action
* id
* class
* charset
* lang
* title

### 例

```go
{{set "title" "change_title"}}
{{title}} {{/* change_title を表示 */}}

{{/* 新規カラムを追加 */}}
{{set "name" "myname"}}
```

## parameter(name string) interface{}
`set`関数で設定された値を取り出します。

```go
{{set "name" "myname"}}

{{$p := parameter "name"}}
{{$p.name}} {{/* myname を表示 */}}
```

## t(name string) interface{}
適用されている自然言語情報を返却する

```go
{{/* 言語設定情報を出力する */}}
{{t "index.appname"}}
```

## i18n(name string, i ...Global) interface{}
```go
{{/* 言語設定情報を取得する */}}
{{$l := i18n parameter "lang"}}

{{/* index.appname を所持しているかチェック */}}
{{if $l.HasItem "index.appname"}}
  {{/* 所持している場合は、値を出力する */}}
  {{$l.index.appname}}
{{end}}

{{/* {"index":{"appname":"XXX"}} の XXX を出力 */}}
{{/* 設定情報が存在しない場合は、空文字列を返却する */}}
{{$l.T "index.appname"}}
 
{{/* 大域上に設定されている自然言語名を変更する */}}
{{i18n "en" global}}
{{t "index.appname"}}
```

## hostname () String
`hostname`関数は、`/bin/hostname`コマンド相当のホスト名を返却する。

```go
{{hostname}} {{/* localhost */}}
```

## env (name string) String
`env`関数は、環境変数の値を返却する。

```go
{{env "PATH"}} {{/* {{/home/user/bin:/usr/bin:/bin:...}} */}}
```

## stylesheet (path string) string
`<link rel='stylesheet' ...` タグを埋め込む。
```go
{{/* <link rel='stylesheet' type='text/css' href='style.css?id=3e167af...' /> */}}
{{stylesheet "style.css"}}
```

## script (path string) string
`<script src='...'` タグを埋め込む。
```go
{{/* <script src='application.js?id=3e167af...'></script>*/}}
{{script "application.js"}}
```

## url () *URL
アクセス先のURLを返却します。
```go
{{url}} {{/* https://localhost:8080/path/to/url?name=test&param=200 を返却 */}}
```
url が返却する *helper.URL は、次のメソッドを所持しています。

### Host () StringType
ホスト名:ポート番号を返却します。
```go
{{/* https://localhost:8080/path/to/url?name=test&param=200 の場合 */}}
{{url.Host}} {{/* localhost:8080 を返却する */}}
```

### Hostname () StringType
ホスト名を返却します。
```go
{{/* https://localhost:8080/path/to/url?name=test&param=200 の場合 */}}
{{url.Hostname}} {{/* localhost を返却する */}}
```

### Port () StringType
ポート番号を返却します。
```go
{{/* https://localhost:8080/path/to/url?name=test&param=200 の場合 */}}
{{url.Port}} {{/* 8080 を返却する */}}
```

### Proto () StringType
プロトコル名を返却します。
```go
{{/* https://localhost:8080/path/to/url?name=test&param=200 の場合 */}}
{{url.Proto}} {{/* https を返却します */}}
```

### Path () StringType
クエリパスを返却します。
```go
{{/* https://localhost:8080/path/to/url?name=test&param=200 の場合 */}}
{{url.Path}} {{/* /path/to/url を返却します */}}
```

### Search () StringType
"?"以降のクエリパラメータを返却します。
```go
{{/* https://localhost:8080/path/to/url?name=test&param=200 の場合 */}}
{{url.Search}} {{/* name=test&param=200 を返却します */}}
```

### Query (i ...string) interface{}
`map`型に変換したクエリパラメータを返却します。

```go
{{/* https://localhost:8080/path/to/url?name=test&param=200 の場合 */}}
{{range $k, $v := url.Query}}
    {{printf "%s: %s" $k $v}} {{/* name: test, param: 200 を出力 */}}
{{end}}

{{/* 引数を指定することで、指定したクエリパラメータのみ取り出すことも可能 */}}
{{range $k, $v := url.Query "name"}}
    {{printf "%s: %s" $k $v}} {{/* name: test を出力 */}}
{{end}}
```

### Get (name string) interface{}
指定したクエリパラメータの値を返却します。

```go
{{/* https://localhost:8080/path/to/url?name=test&param=200 の場合 */}}
{{url.Get "name"}}  {{/* test */}}
{{url.Get "param"}} {{/* 200 */}}
```

### Origin () StringType
プロトコル名、ホスト名、ポート番号が付与されたURLを返却します。
```go
{{/* https://localhost:8080/path/to/url?name=test&param=200 の場合 */}}
{{url.Origin}} {{/* https://localhost:8080 を返却します */}}
```

## isfile (fname StringType) bool
ファイルの存在有無を確認する。

```go
{{/* sample.txt が存在する場合は true。それ以外は false。*/}}
{{isfile "sample.txt"}}
```

## isdir (dir StringType) bool
指定したファイル名がディレクトリか否かを確認する。
```go
{{/* sample.txt が存在し、かつディレクトリの場合は true。それ以外は false。*/}}
{{isfile "sample.txt"}}
```

## stat (fname StringType) (*FileInfo, error)
指定したファイル名の情報を取得する。取得した情報から、次の関数をコール可能。

| Method                | Description                              |
|:--                    |:--                                       |
| `IsDir() bool`        | ディレクトリの場合は true。それ以外は false |
| `ModTime() *DateTime` | ファイル更新、作成日                       |
| `Mode() os.FileMode`  | ファイルのパーミッション                    |
| `Name() string`       | ファイル名 |
| `Size() int64`        | ファイルサイズ |

```go
{{/* ファイル情報を取得 */}}
{{$f := stat "sample.txt"}}

{{/* ファイル情報の詳細を表示 */}}
{{$f.IsDir}}   {{/* false */}}
{{$f.ModTime}} {{/* Thu Dec 20 18:27:49 JST 2018 */}}
{{$f.Mode}}    {{/* -rw-rw-r-- */}}
{{$f.Name}}    {{/* sample.txt */}}
{{$f.Size}}    {{/* 165 */}}
```

## item(value StringType) string
属性値に、値を設定する際に使用する。&, <, >, ", ' をエンティティ文字に置き換える。

```html
<!--
  value には、 'name' が格納されている
  item を利用することにより &#39;name&#39; 文字列へ置き換える
-->
<input type='text' value='{{item .value}}' />
```

## href(path StringType) string
リンク先パス、URLを生成する。http://, https:// から始まるリンク先の場合、別URLへの遷移として扱う。

```html
<!-- http://... -->
<a href='{{href "http://"}}'>name</a>
<!-- https://... -->
<a href='{{href "https://..."}}'>name</a>
<!-- /path/to/url -->
<!-- BaseURLの値が設定されている場合、/baseurl/path/to/url へ変換される -->
<a href='/path/to/url'>name</a>
```