package quic

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/Asutorufa/yuhaiin/pkg/log"
	"github.com/Asutorufa/yuhaiin/pkg/net/dialer"
	"github.com/Asutorufa/yuhaiin/pkg/net/netapi"
	"github.com/Asutorufa/yuhaiin/pkg/protos/node/protocol"
	"github.com/Asutorufa/yuhaiin/pkg/utils/id"
	"github.com/Asutorufa/yuhaiin/pkg/utils/syncmap"
	"github.com/quic-go/quic-go"
)

type Client struct {
	netapi.EmptyDispatch

	tlsConfig  *tls.Config
	quicConfig *quic.Config
	dialer     netapi.Proxy

	session     quic.Connection
	sessionMu   sync.Mutex
	sessionUnix int64

	packetConn *ConnectionPacketConn
	natMap     syncmap.SyncMap[uint64, chan []byte]

	idg id.IDGenerator

	host      *net.UDPAddr
	asNetwork bool
}

func New(config *protocol.Protocol_Quic) protocol.WrapProxy {
	return func(dialer netapi.Proxy) (netapi.Proxy, error) {

		var host *net.UDPAddr = &net.UDPAddr{IP: net.IPv4zero}

		if config.Quic.AsNetwork {
			addr, err := net.ResolveUDPAddr("udp", config.Quic.Host)
			if err != nil {
				return nil, err
			}

			host = addr
		}

		tlsConfig := protocol.ParseTLSConfig(config.Quic.Tls)
		if tlsConfig == nil {
			tlsConfig = &tls.Config{}
		}

		c := &Client{
			dialer:    dialer,
			tlsConfig: tlsConfig,
			quicConfig: &quic.Config{
				MaxIncomingStreams: 2048,
				MaxIdleTimeout:     60 * time.Second,
				KeepAlivePeriod:    45 * time.Second,
				EnableDatagrams:    true,
			},
			asNetwork: config.Quic.AsNetwork,
			host:      host,
		}

		return c, nil
	}
}

func (c *Client) initSession(ctx context.Context) error {
	if c.session != nil {
		select {
		case <-c.session.Context().Done():
		default:
			return nil
		}
	}

	c.sessionMu.Lock()
	defer c.sessionMu.Unlock()

	if c.session != nil {
		select {
		case <-c.session.Context().Done():
		default:
			return nil
		}
	}

	var conn net.PacketConn
	var err error

	if c.asNetwork {
		conn, err = dialer.ListenPacket("udp", "")
	} else {
		conn, err = c.dialer.PacketConn(ctx, netapi.EmptyAddr)
	}
	if err != nil {
		return err
	}

	tr := quic.Transport{
		Conn:               conn,
		ConnectionIDLength: 12,
	}

	session, err := tr.DialEarly(ctx, c.host, c.tlsConfig, c.quicConfig)
	if err != nil {
		return err
	}

	pconn := NewConnectionPacketConn(session)

	c.session = session
	c.packetConn = pconn
	c.sessionUnix = time.Now().Unix()

	go func() {
		for {
			id, data, err := pconn.Receive(context.TODO())
			if err != nil {
				return
			}

			cchan, ok := c.natMap.Load(id)
			if !ok {
				continue
			}

			cchan <- data
		}
	}()
	return nil
}

func (c *Client) Conn(ctx context.Context, s netapi.Address) (net.Conn, error) {
	if err := c.initSession(ctx); err != nil {
		log.Error("init session failed:", "err", err)
		return nil, err
	}

	stream, err := c.session.OpenStream()
	if err != nil {
		return nil, err
	}

	return &interConn{
		Stream:  stream,
		session: c.session,
		time:    c.sessionUnix,
	}, nil
}

func (c *Client) PacketConn(ctx context.Context, host netapi.Address) (net.PacketConn, error) {
	if err := c.initSession(ctx); err != nil {
		return nil, err
	}

	id := c.idg.Generate()
	cchan := make(chan []byte, 10)
	c.natMap.Store(id, cchan)

	return &clientPacketConn{
		c:       c,
		session: c.packetConn,
		id:      id,
		msg:     cchan,
	}, nil
}

var _ net.Conn = (*interConn)(nil)

type interConn struct {
	quic.Stream
	session quic.Connection
	time    int64
}

func (c *interConn) Read(p []byte) (n int, err error) {
	n, err = c.Stream.Read(p)

	if err != nil && err != io.EOF {
		qe, ok := err.(*quic.StreamError)
		if ok && qe.ErrorCode == quic.StreamErrorCode(quic.NoError) {
			err = io.EOF
		}
	}
	return
}

func (c *interConn) Write(p []byte) (n int, err error) {
	n, err = c.Stream.Write(p)
	if err != nil && err != io.EOF {
		qe, ok := err.(*quic.StreamError)
		if ok && qe.ErrorCode == quic.StreamErrorCode(quic.NoError) {
			err = io.EOF
		}
	}
	return
}

func (c *interConn) Close() error {
	return c.Stream.Close()
}

func (c *interConn) LocalAddr() net.Addr {
	return &QuicAddr{
		Addr: c.session.LocalAddr(),
		ID:   c.Stream.StreamID(),
		time: c.time,
	}
}

func (c *interConn) RemoteAddr() net.Addr {
	return &QuicAddr{
		Addr: c.session.RemoteAddr(),
		ID:   c.Stream.StreamID(),
		time: c.time,
	}
}

type QuicAddr struct {
	Addr net.Addr
	ID   quic.StreamID
	time int64
}

func (q *QuicAddr) String() string {
	if q.time == 0 {
		return fmt.Sprint(q.Addr, q.ID)
	}
	return fmt.Sprint(q.Addr, q.time, q.ID)
}

func (q *QuicAddr) Network() string { return "quic" }

type clientPacketConn struct {
	c       *Client
	session *ConnectionPacketConn
	id      uint64

	msg      chan []byte
	deadline *time.Timer
}

func (x *clientPacketConn) ReadFrom(p []byte) (n int, _ net.Addr, err error) {
	msg, ok := <-x.msg
	if !ok {
		return 0, nil, io.EOF
	}

	n = copy(p, msg)
	return n, x.session.conn.RemoteAddr(), nil
}

func (x *clientPacketConn) WriteTo(p []byte, _ net.Addr) (n int, err error) {
	err = x.session.Write(p, x.id)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (x *clientPacketConn) Close() error {
	x.c.natMap.Delete(x.id)
	return nil
}

func (x *clientPacketConn) LocalAddr() net.Addr {
	return &QuicAddr{
		Addr: x.session.conn.LocalAddr(),
		ID:   quic.StreamID(x.id),
	}
}

func (x *clientPacketConn) SetDeadline(t time.Time) error {
	if x.deadline == nil {
		if !t.IsZero() {
			x.deadline = time.AfterFunc(time.Until(t), func() { x.Close() })
		}
		return nil
	}

	if t.IsZero() {
		x.deadline.Stop()
	} else {
		x.deadline.Reset(time.Until(t))
	}
	return nil
}
func (x *clientPacketConn) SetReadDeadline(t time.Time) error  { return x.SetDeadline(t) }
func (x *clientPacketConn) SetWriteDeadline(t time.Time) error { return x.SetDeadline(t) }
