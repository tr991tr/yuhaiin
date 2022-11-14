package dns

import (
	"fmt"
	"math"
	"math/big"
	"net"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/Asutorufa/yuhaiin/pkg/net/interfaces/dns"
	"github.com/Asutorufa/yuhaiin/pkg/utils/lru"
	"github.com/Asutorufa/yuhaiin/pkg/utils/pool"
	"golang.org/x/net/dns/dnsmessage"
)

type Fake struct {
	domainToIP *lru.LRU[string, string]
	ipRange    *net.IPNet

	mu sync.Mutex
}

func NewFake(ipRange *net.IPNet) *Fake {
	ones, bits := ipRange.Mask.Size()
	lruSize := int(math.Pow(2, float64(bits-ones)) - 1)
	// if lruSize > 250 {
	// 	lruSize = 250
	// }
	return &Fake{
		ipRange:    ipRange,
		domainToIP: lru.NewLru[string, string](uint(lruSize), 0*time.Minute),
	}
}

// GetFakeIPForDomain checks and generates a fake IP for a domain name
func (fkdns *Fake) GetFakeIPForDomain(domain string) string {
	fkdns.mu.Lock()
	defer fkdns.mu.Unlock()

	if v, ok := fkdns.domainToIP.Load(domain); ok {
		return v
	}
	currentTimeMillis := uint64(time.Now().UnixNano() / 1e6)
	ones, bits := fkdns.ipRange.Mask.Size()
	rooms := bits - ones
	if rooms < 64 {
		currentTimeMillis %= (uint64(1) << rooms)
	}

	bigIntIP := big.NewInt(0).SetBytes(fkdns.ipRange.IP)
	bigIntIP = bigIntIP.Add(bigIntIP, new(big.Int).SetUint64(currentTimeMillis))

	var bytesLen, fillIndex int
	if fkdns.ipRange.IP.To4() == nil { // ipv6
		bytesLen = net.IPv6len
		if len(bigIntIP.Bytes()) != net.IPv6len {
			fillIndex = 1
		}
	} else {
		bytesLen = net.IPv4len
	}

	bytes := pool.GetBytes(bytesLen)
	defer pool.PutBytes(bytes)

	var ip net.IP
	for {
		bigIntIP.FillBytes(bytes[fillIndex:])
		ip = net.IP(bytes)

		// if we run for a long time, we may go back to beginning and start seeing the IP in use
		if ok := fkdns.domainToIP.ValueExist(ip.String()); !ok {
			break
		}

		bigIntIP = bigIntIP.Add(bigIntIP, big.NewInt(1))

		bigIntIP.FillBytes(bytes[fillIndex:])
		if !fkdns.ipRange.Contains(bytes) {
			bigIntIP = big.NewInt(0).SetBytes(fkdns.ipRange.IP)
		}
	}
	fkdns.domainToIP.Add(domain, ip.String())
	return ip.String()
}

func (fkdns *Fake) GetDomainFromIP(ip string) (string, bool) {
	fkdns.mu.Lock()
	defer fkdns.mu.Unlock()
	return fkdns.domainToIP.ReverseLoad(ip)
}

var _ dns.DNS = (*FakeDNS)(nil)

type FakeDNS struct {
	upStreamDo func(b []byte) ([]byte, error)
	pool       *Fake
}

func WrapFakeDNS(upStreamDo func(b []byte) ([]byte, error), pool *Fake) *FakeDNS {
	return &FakeDNS{upStreamDo: upStreamDo, pool: pool}
}
func (f *FakeDNS) LookupIP(domain string) ([]net.IP, error) {
	ip := f.pool.GetFakeIPForDomain(domain)
	// log.Println("map", ip, "to", domain)

	return []net.IP{net.ParseIP(ip).To4()}, nil
}

func (f *FakeDNS) Record(domain string, t dnsmessage.Type) (dns.IPResponse, error) {
	ip := f.pool.GetFakeIPForDomain(domain)

	if t == dnsmessage.TypeAAAA {
		return dns.NewIPResponse([]net.IP{net.ParseIP(ip).To16()}, 600), nil
	}
	return dns.NewIPResponse([]net.IP{net.ParseIP(ip).To4()}, 600), nil
}

func (f *FakeDNS) LookupPtr(name string) (string, error) {
	ip := pool.GetBuffer()
	defer pool.PutBuffer(ip)

	i := strings.Index(name, ".in-addr.arpa.")
	if i == -1 {
		i = strings.Index(name, ".ip6.arpa.")
	}

	if i == -1 {
		return "", fmt.Errorf("ptr format error: %s", name)
	}

	p := strings.Split(name[:i], ".")
	for i, v4 := len(p)-1, len(p) == 4; i >= 0; i-- {
		ip.WriteString(p[i])
		if i != 0 {
			if v4 {
				ip.WriteByte('.')
			} else if i%4 == 0 {
				ip.WriteByte(':')
			}
		}
	}

	b := ip.Bytes()
	r, ok := f.pool.GetDomainFromIP(*(*string)(unsafe.Pointer(&b)))
	if !ok {
		return "", fmt.Errorf("not found %s[%s] ptr", ip, name)
	}

	return r, nil
}

func (f *FakeDNS) Do(b []byte) ([]byte, error) { return f.upStreamDo(b) }

func (f *FakeDNS) Close() error { return nil }
