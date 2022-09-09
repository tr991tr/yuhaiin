package yuhaiin

import "github.com/Asutorufa/yuhaiin/pkg/utils/syncmap"

type Opts struct {
	Host     string      `json:"host"`
	Savepath string      `json:"savepath"`
	Socks5   string      `json:"socks5"`
	Http     string      `json:"http"`
	IPv6     bool        `json:"ipv6"`
	Bypass   *Bypass     `json:"bypass"`
	DNS      *DNSSetting `json:"dns"`
	TUN      *TUN        `json:"tun"`
	Log      *Log        `json:"log"`
}

type Log struct {
	SaveLogcat bool `json:"save_logcat"`
	// 0:verbose, 1:debug, 2:info, 3:warning, 4:error, 5: fatal
	LogLevel int32 `json:"log_level"`
}

type Bypass struct {
	// 0: bypass, 1: proxy, 2: direct, 3: block
	TCP int32 `json:"tcp"`
	// 0: bypass, 1: proxy, 2: direct, 3: block
	UDP int32 `json:"udp"`

	Block  string `json:"block"`
	Proxy  string `json:"proxy"`
	Direct string `json:"direct"`
}

type DNSSetting struct {
	Server         string `json:"server"`
	Fakedns        bool   `json:"fakedns"`
	FakednsIpRange string `json:"fakedns_ip_range"`
	Remote         *DNS   `json:"remote"`
	Local          *DNS   `json:"local"`
	Bootstrap      *DNS   `json:"bootstrap"`
}

type DNS struct {
	Host string `json:"host"`
	// Type
	// 0: reserve
	// 1: udp
	// 2: tcp
	// 3: doh
	// 4: dot
	// 5: doq
	// 6: doh3
	Type          int32  `json:"type"`
	Proxy         bool   `json:"proxy"`
	Subnet        string `json:"subnet"`
	TlsServername string `json:"tls_servername"`
}

type TUN struct {
	FD           int32  `json:"fd"`
	MTU          int32  `json:"mtu"`
	Gateway      string `json:"gateway"`
	DNSHijacking bool   `json:"dns_hijacking"`
	// Driver
	// 0: fdbased
	// 1: channel
	Driver    int32 `json:"driver"`
	UidDumper UidDumper
}

type UidDumper interface {
	DumpUid(ipProto int32, srcIp string, srcPort int32, destIp string, destPort int32) (int32, error)
	GetUidInfo(uid int32) (string, error)
}

type uidDumper struct {
	UidDumper
	cache syncmap.SyncMap[int32, string]
}

func NewUidDumper(ud UidDumper) UidDumper {
	if ud == nil {
		return nil
	}
	return &uidDumper{UidDumper: ud}
}

func (u *uidDumper) GetUidInfo(uid int32) (string, error) {
	if r, ok := u.cache.Load(uid); ok {
		return r, nil
	}

	r, err := u.UidDumper.GetUidInfo(uid)
	if err != nil {
		return "", err
	}

	u.cache.Store(uid, r)
	return r, nil
}