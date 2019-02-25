package mux

import (
	"fmt"
	"net"
	"strings"
)

// IP : パス毎にIP制限を設ける構造体
type IP struct {
	IsAllow bool
	Path    string
	Addr    []string
	ipnet   []*net.IPNet
}

// CreateIPNet : IP 構造体に登録されている情報から net.IPNet を生成する
func (ip *IP) CreateIPNet() error {
	// ip リストが未作成の場合はエラーを返却する
	if len(ip.Addr) == 0 {
		return fmt.Errorf("not setting ip tables")
	}
	// /path/to/url => /path/to/url/ へ変換する
	ip.Path = strings.TrimRight(ip.Path, "/") + "/"

	// 指定されたIPリスト順に、net.IPNet を生成する
	for _, addr := range ip.Addr {
		// all が指定されている場合は、 0.0.0.0/0 に変換する
		if strings.ToLower(addr) == "all" {
			addr = "0.0.0.0/0"
		}

		// IPアドレスに '/' が含まれていない場合は、 /32 を付け加える
		if strings.Index(addr, "/") == -1 {
			addr += "/32"
		}

		// IPアドレスをCIDR形式で解析
		_, ipnet, err := net.ParseCIDR(addr)
		if err != nil {
			return err
		}
		ip.ipnet = append(ip.ipnet, ipnet)
	}

	// エラーがなければ nil を返却して復帰する
	return nil
}

// Contains : 指定されたIPアドレスにマッチした場合は true を返却する
func (ip *IP) Contains(addr string) bool {
	// 127.0.0.1 などのIPを net.IP 型経パースする
	parseIP := net.ParseIP(addr)
	for _, v := range ip.ipnet {
		// IPがマッチした場合
		if v.Contains(parseIP) {
			// Allow の場合は true を、Denyの場合はfalseを返却する
			return ip.IsAllow
		}
	}
	// マッチしない場合、Allowの場合は false を、Denyの場合は false を返却する
	return !ip.IsAllow
}

// RestrictIP : IP構造体を一括管理する配列
type RestrictIP []*IP

// MakeIPNet : IP構造体に登録されているデータを用いて、IPNet構造体を生成する
func (iplist RestrictIP) MakeIPNet() error {
	for _, ip := range iplist {
		if err := ip.CreateIPNet(); err != nil {
			return err
		}
	}
	return nil
}

// Contains : 指定されたパスとIPに制限がかかっていないか確認する関数
func (iplist RestrictIP) Contains(path, addr string) bool {
	path = strings.TrimRight(path, "/") + "/"
	// アクセス制限用のリストが存在していない場合は true を返却する
	if len(iplist) == 0 {
		return true
	}
	var result bool
	// アクセス元のクエリパスから、該当するパスでIP制限がかかっているか確認する
	for _, v := range iplist {
		// アクセス元のクエリパスが、IP制限パスと不一致の場合は次の設定へ
		if strings.Index(path, v.Path) != 0 {
			continue
		}
		result = v.Contains(addr)
		if v.IsAllow {
			// Allow 設定でかつ、IP制限でアクセスできない場合は、次のAllow設定を参照する
			if result == false {
				continue
			}
		} else {
			// Deny 設定でかつ、IP制限にかかっていない場合は、次のDeny設定を参照する
			if result == true {
				continue
			}
		}
		return result
	}
	// 該当パスが存在しない場合は true を返却する
	return true
}
