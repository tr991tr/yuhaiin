package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"unsafe"

	"github.com/Asutorufa/yuhaiin/pkg/log"
	"github.com/Asutorufa/yuhaiin/pkg/net/netapi"
	"github.com/Asutorufa/yuhaiin/pkg/net/proxy/socks5/tools"
	"github.com/Asutorufa/yuhaiin/pkg/protos/config/listener"
	"github.com/Asutorufa/yuhaiin/pkg/protos/statistic"
	"github.com/Asutorufa/yuhaiin/pkg/utils/pool"
	"github.com/Asutorufa/yuhaiin/pkg/utils/relay"
)

func (s *Socks5) startTCPServer() error {
	lis, err := s.lis.Stream(s.ctx)
	if err != nil {
		return err
	}
	go func() {
		defer s.Close()
		defer lis.Close()

		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Error("socks5 accept failed", "err", err)

				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					continue
				}
				return
			}

			go func() {
				if err := s.Handle(conn); err != nil {
					if errors.Is(err, netapi.ErrBlocked) {
						log.Debug(err.Error())
					} else {
						log.Error("socks5 server handle failed", "err", err)
					}
				}
			}()

		}
	}()

	return nil
}

func (s *Socks5) Handle(client net.Conn) (err error) {
	b := pool.GetBytes(pool.DefaultSize)
	defer pool.PutBytes(b)

	err = s.handshake1(client, b)
	if err != nil {
		return fmt.Errorf("first hand failed: %w", err)
	}

	if err = s.handshake2(client, b); err != nil {
		return fmt.Errorf("second hand failed: %w", err)
	}

	return
}

func (s *Socks5) handshake1(client net.Conn, buf []byte) error {
	//socks5 first handshake
	if _, err := io.ReadFull(client, buf[:2]); err != nil {
		return fmt.Errorf("read first handshake failed: %w", err)
	}

	if buf[0] != 0x05 { // ver
		err := writeHandshake1(client, tools.NoAcceptableMethods)
		return fmt.Errorf("no acceptable method: %d, resp err: %w", buf[0], err)
	}

	nMethods := int(buf[1])

	if nMethods > len(buf) {
		err := writeHandshake1(client, tools.NoAcceptableMethods)
		return fmt.Errorf("nMethods length of methods out of buf, resp err: %w", err)
	}

	if _, err := io.ReadFull(client, buf[:nMethods]); err != nil {
		return fmt.Errorf("read methods failed: %w", err)
	}

	noNeedVerify := s.username == "" && s.password == ""
	userAndPasswordSupport := false

	for _, v := range buf[:nMethods] { // range all supported methods
		if v == tools.NoAuthenticationRequired && noNeedVerify {
			return writeHandshake1(client, tools.NoAuthenticationRequired)
		}

		if v == tools.UserAndPassword {
			userAndPasswordSupport = true
		}
	}

	if userAndPasswordSupport {
		return verifyUserPass(client, s.username, s.password)
	}

	err := writeHandshake1(client, tools.NoAcceptableMethods)

	return fmt.Errorf("no acceptable authentication methods: [length: %d, method:%v], response err: %w", nMethods, buf[:nMethods], err)
}

func verifyUserPass(client net.Conn, user, key string) error {
	if err := writeHandshake1(client, tools.UserAndPassword); err != nil {
		return err
	}

	b := pool.GetBytes(pool.DefaultSize)
	defer pool.PutBytes(b)

	if _, err := io.ReadFull(client, b[:2]); err != nil {
		return fmt.Errorf("read ver and user name length failed: %w", err)
	}

	// if b[0] != 0x01 {
	// 	return fmt.Errorf("unknown ver: %d", b[0])
	// }

	usernameLength := int(b[1])

	if _, err := io.ReadFull(client, b[2:2+usernameLength]); err != nil {
		return fmt.Errorf("read username failed: %w", err)
	}

	username := b[2 : 2+usernameLength]

	if _, err := io.ReadFull(client, b[2+usernameLength:2+usernameLength+1]); err != nil {
		return fmt.Errorf("read password length failed: %w", err)
	}

	passwordLength := int(b[2+usernameLength])

	if _, err := io.ReadFull(client, b[2+usernameLength+1:2+usernameLength+1+passwordLength]); err != nil {
		return fmt.Errorf("read password failed: %w", err)
	}

	password := b[2+usernameLength+1 : 2+usernameLength+1+passwordLength]

	if (len(user) > 0 && (usernameLength <= 0 || user != unsafe.String(&username[0], usernameLength))) ||
		(len(key) > 0 && (passwordLength <= 0 || key != unsafe.String(&password[0], passwordLength))) {
		_, err := client.Write([]byte{1, 1})
		return fmt.Errorf("verify username and password failed, resp err: %w", err)
	}

	_, err := client.Write([]byte{1, 0})
	return err
}

func (s *Socks5) handshake2(client net.Conn, buf []byte) error {
	// socks5 second handshake
	if _, err := io.ReadFull(client, buf[:3]); err != nil {
		return fmt.Errorf("read second handshake failed: %w", err)
	}

	if buf[0] != 0x05 { // ver
		err := writeHandshake2(client, tools.NoAcceptableMethods, netapi.EmptyAddr)
		return fmt.Errorf("no acceptable method: %d, resp err: %w", buf[0], err)
	}

	var err error

	switch tools.CMD(buf[1]) { // mode
	case tools.Connect:
		var adr tools.ADDR
		adr, err = tools.ResolveAddr(client)
		if err != nil {
			return fmt.Errorf("resolve addr failed: %w", err)
		}

		addr := adr.Address(statistic.Type_tcp)

		caddr, err := netapi.ParseSysAddr(client.LocalAddr())
		if err != nil {
			return fmt.Errorf("parse local addr failed: %w", err)
		}
		err = writeHandshake2(client, tools.Succeeded, caddr) // response to connect successful
		if err != nil {
			return err
		}

		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		case s.tcpChannel <- &netapi.StreamMeta{
			Source:      client.RemoteAddr(),
			Destination: addr,
			Inbound:     client.LocalAddr(),
			Src:         client,
			Address:     addr,
		}:
		}

	case tools.Udp: // udp
		if s.udp {
			err = handleUDP(client)
			break
		}
		fallthrough

	case tools.Bind: // bind request
		fallthrough

	default:
		err := writeHandshake2(client, tools.CommandNotSupport, netapi.EmptyAddr)
		return fmt.Errorf("not Support Method %d, resp err: %w", buf[1], err)
	}

	if err != nil {
		_ = writeHandshake2(client, tools.HostUnreachable, netapi.EmptyAddr)
	}
	return err
}

func handleUDP(client net.Conn) error {
	laddr, err := netapi.ParseSysAddr(client.LocalAddr())
	if err != nil {
		return fmt.Errorf("parse sys addr failed: %w", err)
	}
	err = writeHandshake2(client, tools.Succeeded, netapi.ParseAddressPort(statistic.Type_tcp, "0.0.0.0", laddr.Port()))
	if err != nil {
		return err
	}
	_, _ = relay.Copy(io.Discard, client)
	return nil
}

func writeHandshake1(conn net.Conn, errREP byte) error {
	_, err := conn.Write([]byte{0x05, errREP})
	return err
}

func writeHandshake2(conn net.Conn, errREP byte, addr netapi.Address) error {
	_, err := conn.Write(append([]byte{0x05, errREP, 0x00}, tools.ParseAddr(addr)...))
	return err
}

type Socks5 struct {
	udp      bool
	lis      netapi.Listener
	username string
	password string

	ctx   context.Context
	close context.CancelFunc

	tcpChannel chan *netapi.StreamMeta
	udpChannel chan *netapi.Packet
}

func (s *Socks5) Close() error {
	s.close()
	return s.lis.Close()
}

func (s *Socks5) AcceptStream() (*netapi.StreamMeta, error) {
	select {
	case <-s.ctx.Done():
		return nil, s.ctx.Err()
	case meta := <-s.tcpChannel:
		return meta, nil
	}
}

func (s *Socks5) AcceptPacket() (*netapi.Packet, error) {
	select {
	case <-s.ctx.Done():
		return nil, s.ctx.Err()
	case packet := <-s.udpChannel:
		return packet, nil
	}
}

func init() {
	listener.RegisterProtocol2(NewServer)
}

func NewServer(o *listener.Inbound_Socks5) func(netapi.Listener) (netapi.ProtocolServer, error) {
	return func(ii netapi.Listener) (netapi.ProtocolServer, error) {
		ctx, cancel := context.WithCancel(context.TODO())
		s := &Socks5{
			udp:        o.Socks5.Udp,
			username:   o.Socks5.Username,
			password:   o.Socks5.Password,
			lis:        ii,
			ctx:        ctx,
			close:      cancel,
			tcpChannel: make(chan *netapi.StreamMeta, 100),
			udpChannel: make(chan *netapi.Packet, 100),
		}

		if s.udp {
			if err := s.startUDPServer(); err != nil {
				return nil, err
			}
		}

		if err := s.startTCPServer(); err != nil {
			s.Close()
			return nil, err
		}

		return s, nil
	}
}
