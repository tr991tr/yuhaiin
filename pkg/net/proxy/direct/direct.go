package direct

import (
	"context"
	"fmt"
	"net"

	"github.com/Asutorufa/yuhaiin/pkg/net/dialer"
	"github.com/Asutorufa/yuhaiin/pkg/net/netapi"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node/point"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node/protocol"
)

type direct struct{ netapi.EmptyDispatch }

func init() {
	point.RegisterProtocol(func(p *protocol.Protocol_Direct) point.WrapProxy {
		return func(netapi.Proxy) (netapi.Proxy, error) {
			return Default, nil
		}
	})

	point.SetBootstrap(Default)
}

var Default netapi.Proxy = NewDirect()

func NewDirect() netapi.Proxy {
	return &direct{}
}

func (d *direct) Conn(ctx context.Context, s netapi.Address) (net.Conn, error) {
	return netapi.DialHappyEyeballs(ctx, s)
}

func (d *direct) PacketConn(context.Context, netapi.Address) (net.PacketConn, error) {
	p, err := dialer.ListenPacket("udp", "")
	if err != nil {
		return nil, fmt.Errorf("listen packet failed: %w", err)
	}

	return &PacketConn{p}, nil
}

type PacketConn struct{ net.PacketConn }

func (p *PacketConn) WriteTo(b []byte, addr net.Addr) (_ int, err error) {
	udpAddr, ok := addr.(*net.UDPAddr)
	if !ok {
		a, err := netapi.ParseSysAddr(addr)
		if err != nil {
			return 0, err
		}

		ur := a.UDPAddr(context.TODO())
		if ur.Err != nil {
			return 0, ur.Err
		}

		udpAddr = ur.V
	}

	return p.PacketConn.WriteTo(b, udpAddr)
}
