package socks4a

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"unsafe"

	"github.com/Asutorufa/yuhaiin/pkg/log"
	"github.com/Asutorufa/yuhaiin/pkg/net/netapi"
	"github.com/Asutorufa/yuhaiin/pkg/protos/config/listener"
	"github.com/Asutorufa/yuhaiin/pkg/protos/statistic"
	"github.com/Asutorufa/yuhaiin/pkg/utils/pool"
)

const (
	CommandConnect byte = 0x01
	CommandBind    byte = 0x02
)

type Server struct {
	lis        net.Listener
	usernameID string

	*netapi.ChannelServer
}

func (s *Server) Handle(conn net.Conn) error {
	addr, err := s.Handshake(conn)
	if err != nil {
		_, _ = conn.Write([]byte{0, 91, 0, 0, 0, 0, 0, 0})
		return fmt.Errorf("handshake failed: %w", err)
	}

	return s.SendStream(&netapi.StreamMeta{
		Source:      conn.RemoteAddr(),
		Destination: addr,
		Inbound:     conn.LocalAddr(),
		Src:         conn,
		Address:     addr,
	})
}

func (s *Server) Handshake(conn net.Conn) (netapi.Address, error) {
	buf := pool.GetBytesBuffer(8)
	defer buf.Free()

	if _, err := io.ReadFull(conn, buf.Bytes()); err != nil {
		return nil, err
	}

	if buf.Bytes()[0] != 0x04 {
		return nil, fmt.Errorf("unknown socks version: %d", buf.Bytes()[0])
	}

	if buf.Bytes()[1] != CommandConnect {
		return nil, fmt.Errorf("unsupported command: %d", buf.Bytes()[1])
	}

	port := binary.BigEndian.Uint16(buf.Bytes()[2:4])
	dstAddr := buf.Bytes()[4:8]
	userId, err := readData(conn)
	if err != nil {
		return nil, err
	}

	if s.usernameID != "" && !bytes.Equal(userId, unsafe.Slice(unsafe.StringData(s.usernameID), len(s.usernameID))) {
		return nil, fmt.Errorf("username not match")
	}

	var target netapi.Address
	if dstAddr[0] == 0 && dstAddr[1] == 0 && dstAddr[2] == 0 && dstAddr[3] != 0 {
		host, err := readData(conn)
		if err != nil {
			return nil, err
		}
		target = netapi.ParseAddressPort(statistic.Type_tcp, string(host), netapi.ParsePort(port))
	} else {
		target = netapi.ParseIPAddrPort(statistic.Type_tcp, dstAddr, int(port))
	}

	_, _ = conn.Write([]byte{0, 90})
	_, _ = conn.Write(buf.Bytes()[2:8])
	return target, nil
}

func readData(conn net.Conn) ([]byte, error) {
	var data []byte

	buf := make([]byte, 1)

	for {
		if _, err := io.ReadFull(conn, buf); err != nil {
			return nil, err
		}

		if buf[0] == 0 {
			break
		}

		data = append(data, buf[0])
	}

	return data, nil
}

func (s *Server) Close() error {
	s.ChannelServer.Close()

	if s.lis != nil {
		return s.lis.Close()
	}

	return nil
}

func (s *Server) Server() {
	defer s.Close()
	for {
		conn, err := s.lis.Accept()
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
}

func (s *Server) AcceptPacket() (*netapi.Packet, error) {
	return nil, io.EOF
}

func init() {
	listener.RegisterProtocol(NewServer)
}

func NewServer(o *listener.Inbound_Socks4A) func(netapi.Listener) (netapi.Accepter, error) {
	return func(ii netapi.Listener) (netapi.Accepter, error) {
		lis, err := ii.Stream(context.TODO())
		if err != nil {
			return nil, err
		}

		s := &Server{
			usernameID:    o.Socks4A.Username,
			lis:           lis,
			ChannelServer: netapi.NewChannelServer(),
		}

		go s.Server()

		return s, nil
	}
}
