package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/Asutorufa/yuhaiin/pkg/log"
	"github.com/Asutorufa/yuhaiin/pkg/net/dialer"
	"github.com/Asutorufa/yuhaiin/pkg/net/interfaces/proxy"
	"github.com/Asutorufa/yuhaiin/pkg/net/nat"
	s5c "github.com/Asutorufa/yuhaiin/pkg/net/proxy/socks5/client"
	"github.com/Asutorufa/yuhaiin/pkg/protos/statistic"
	"github.com/Asutorufa/yuhaiin/pkg/utils/pool"
)

type udpServer struct {
	net.PacketConn
	natTable *nat.Table
}

func (s *Socks5) newUDPServer() error {
	l, err := dialer.ListenPacket("udp", s.addr)
	if err != nil {
		return fmt.Errorf("listen udp failed: %w", err)
	}

	u := &udpServer{PacketConn: l, natTable: nat.NewTable(s.dialer)}
	s.udpServer = u

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go func() {
			defer s.Close()
			if err := u.handle(); err != nil && !errors.Is(err, net.ErrClosed) {
				log.Errorln("handle udp request failed:", err)
			}

		}()
	}

	return nil
}

func (u *udpServer) handle() error {
	buf := pool.GetBytes(nat.MaxSegmentSize)
	defer pool.PutBytes(buf)

	for {
		n, src, err := u.PacketConn.ReadFrom(buf)
		if err != nil {
			return err
		}

		addr, err := s5c.ResolveAddr(bytes.NewReader(buf[3:n]))
		if err != nil {
			return fmt.Errorf("resolve addr failed: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*15)
		defer cancel()

		dst := addr.Address(statistic.Type_udp)
		dst.WithContext(ctx)

		return u.natTable.Write(
			&nat.Packet{
				Src:     src,
				Dst:     dst,
				Payload: buf[3+len(addr) : n],
				WriteBack: func(b []byte, source net.Addr) (int, error) {
					sourceAddr, err := proxy.ParseSysAddr(source)
					if err != nil {
						return 0, err
					}
					b = bytes.Join([][]byte{{0, 0, 0}, s5c.ParseAddr(sourceAddr), b}, nil)

					return u.PacketConn.WriteTo(b, src)
				},
			},
		)
	}
}

func (u *udpServer) Close() error {
	u.natTable.Close()
	return u.PacketConn.Close()
}
