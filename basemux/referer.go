package basemux

import (
	"container/list"
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// referer 構造体を生成する
func newRefer(latency time.Duration) *referer {
	r := &referer{
		list: list.New(),
		data: make(map[string]*list.Element),
	}

	t := time.NewTicker(latency * time.Second)
	go func() {
		for {
			select {
			// latency に設定した値を元に、データの期限切れをチェックして削除する
			case <-t.C:
				r.inspection(latency)
			}
		}
	}()
	return r
}

type referer struct {
	mu   sync.Mutex
	list *list.List
	data map[string]*list.Element
}

// 新規リンクIDを生成する
func (r *referer) generateID() string {
	// ランダムIDを生成するための、元データ
	const randomid = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// シード値を設定
	random.Seed(time.Now().UnixNano())
	// ランダム文字列を生成
	buf := make([]byte, 64)
	for i := range buf {
		buf[i] = randomid[random.Intn(len(randomid))]
	}
	// SHA256オブジェクト形式のランダムIDを生成する
	hash := md5.New()
	io.WriteString(hash, string(buf)+fmt.Sprint(time.Now().UnixNano()))
	return fmt.Sprintf("%X", hash.Sum(nil))
}

// リンク切れの情報が存在しないかチェックする
func (r *referer) inspection(latency time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for {
		// 古いアクセス履歴が前に存在しているため、前から探す
		elem := r.list.Front()
		// アクセス履歴がみつからない場合はループを抜ける
		if elem == nil {
			break
		}
		// アクセス履歴が見つかった場合、期限切れチェックを行う
		data := elem.Value.(*Values)
		access := data.access + int64(latency*2)
		if access < time.Now().Unix() {
			// 期限切れデータを破棄
			delete(r.data, data.id)
			r.list.Remove(elem)
		} else {
			// 期限切れのデータ存在しない場合処理を終了
			break
		}
	}
}

// リファラのIDをキーに、登録されているデータを取り出す
func (r *referer) Get(id string) *Values {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 登録データが存在しない場合 nil を返却する
	v, ok := r.data[id]
	if !ok {
		return nil
	}

	// 登録データのaccessを更新
	values := v.Value.(*Values)
	values.access = time.Now().Unix()

	// 登録データを返却する
	return values
}

// リファラに登録するデータを作成する
func (r *referer) Create() *Values {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 後方に、生成するデータを追加する
	v := &Values{
		access: time.Now().Unix(),
		id:     r.generateID(),
		data:   make(map[string]interface{}),
		old:    r,
	}
	r.data[v.id] = r.list.PushBack(v)

	return v
}

// Referer : リファラを取り扱うインタフェース
type Referer interface {
	Get(string) *Values
}
