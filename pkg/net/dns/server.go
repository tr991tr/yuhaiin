package dns

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/Asutorufa/yuhaiin/pkg/log"
	"github.com/Asutorufa/yuhaiin/pkg/net/dialer"
	"github.com/Asutorufa/yuhaiin/pkg/net/interfaces/dns"
	"github.com/Asutorufa/yuhaiin/pkg/net/interfaces/proxy"
	"github.com/Asutorufa/yuhaiin/pkg/net/interfaces/server"
	"github.com/Asutorufa/yuhaiin/pkg/net/nat"
	"github.com/Asutorufa/yuhaiin/pkg/utils/pool"
	"golang.org/x/exp/slog"
	"golang.org/x/net/dns/dnsmessage"
)

type dnsServer struct {
	server      string
	resolver    dns.DNS
	listener    net.PacketConn
	tcpListener net.Listener
}

func NewDnsServer(server string, process dns.DNS) server.DNSServer {
	d := &dnsServer{server: server, resolver: process}

	if server == "" {
		log.Warn("dns server is empty, skip to listen tcp and udp")
		return d
	}

	go func() {
		if err := d.start(); err != nil {
			log.Error("start udp dns server failed", slog.Any("err", err))
		}
	}()

	go func() {
		if err := d.startTCP(); err != nil {
			log.Error("start tcp dns server failed", slog.Any("err", err))
		}
	}()

	return d
}

func (d *dnsServer) Close() error {
	if d.listener != nil {
		d.listener.Close()
	}
	if d.tcpListener != nil {
		d.tcpListener.Close()
	}

	return nil
}

func (d *dnsServer) start() (err error) {
	d.listener, err = dialer.ListenPacket("udp", d.server)
	if err != nil {
		return fmt.Errorf("dns udp server listen failed: %w", err)
	}
	defer d.listener.Close()
	log.Info("new udp dns server", "host", d.server)

	for {
		buf := pool.GetBytes(nat.MaxSegmentSize)

		n, addr, err := d.listener.ReadFrom(buf)
		if err != nil {
			if e, ok := err.(net.Error); ok {
				if e.Temporary() {
					continue
				}
			}
			return fmt.Errorf("dns udp server handle failed: %w", err)
		}

		go func() {
			defer pool.PutBytes(buf)

			ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
			defer cancel()

			data, err := d.handle(ctx, buf[:n])
			if err != nil {
				log.Error("dns server handle data failed", slog.Any("err", err))
				return
			}

			if _, err = d.listener.WriteTo(data, addr); err != nil {
				log.Error("write dns response to client failed", slog.Any("err", err))
			}
		}()
	}
}

func (d *dnsServer) startTCP() (err error) {
	d.tcpListener, err = net.Listen("tcp", d.server)
	if err != nil {
		return fmt.Errorf("dns tcp server listen failed: %w", err)
	}
	defer d.tcpListener.Close()
	log.Error("new tcp dns server", "host", d.server)
	for {
		conn, err := d.tcpListener.Accept()
		if err != nil {
			if e, ok := err.(net.Error); ok {
				if e.Temporary() {
					continue
				}
			}
			return fmt.Errorf("dns server accept failed: %w", err)
		}

		go func() {
			defer conn.Close()

			ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
			defer cancel()

			if err := d.HandleTCP(ctx, conn); err != nil {
				log.Error("handle dns tcp failed", "err", err)
			}
		}()
	}
}

func (d *dnsServer) HandleTCP(ctx context.Context, c net.Conn) error {
	var length uint16
	if err := binary.Read(c, binary.BigEndian, &length); err != nil {
		return fmt.Errorf("read dns length failed: %w", err)
	}

	data := pool.GetBytes(int(length))
	defer pool.PutBytes(data)

	n, err := io.ReadFull(c, data[:length])
	if err != nil {
		return fmt.Errorf("dns server read data failed: %w", err)
	}

	data, err = d.handle(ctx, data[:n])
	if err != nil {
		return fmt.Errorf("dns server handle failed: %w", err)
	}

	if err = binary.Write(c, binary.BigEndian, uint16(len(data))); err != nil {
		return fmt.Errorf("dns server write length failed: %w", err)
	}
	_, err = c.Write(data)
	return err
}

func (d *dnsServer) HandleUDP(ctx context.Context, l net.PacketConn) error {
	p := pool.GetBytes(nat.MaxSegmentSize)
	defer pool.PutBytes(p)
	n, addr, err := l.ReadFrom(p)
	if err != nil {
		return err
	}

	data, err := d.handle(ctx, p[:n])
	if err != nil {
		return fmt.Errorf("dns server handle failed: %w", err)
	}
	_, err = l.WriteTo(data, addr)
	return err
}

func (d *dnsServer) Do(ctx context.Context, b []byte) ([]byte, error) { return d.handle(ctx, b) }

func (d *dnsServer) handle(ctx context.Context, b []byte) ([]byte, error) {
	var parse dnsmessage.Parser

	h, err := parse.Start(b)
	if err != nil {
		return nil, fmt.Errorf("dns server parse failed: %w", err)
	}

	q, err := parse.Question()
	if err != nil {
		return nil, fmt.Errorf("dns server parse failed: %w", err)
	}

	add := strings.TrimSuffix(q.Name.String(), ".")

	if q.Type != dnsmessage.TypeA && q.Type != dnsmessage.TypeAAAA &&
		q.Type != dnsmessage.TypePTR {
		log.Debug("not a, aaaa or ptr", "type", q.Type)
		return d.resolver.Do(ctx, add, b)
	}

	resp := dnsmessage.Message{
		Header: dnsmessage.Header{
			ID:                 h.ID,
			Response:           true,
			Authoritative:      false,
			RecursionDesired:   false,
			RCode:              dnsmessage.RCodeSuccess,
			RecursionAvailable: false,
		},
		Questions: []dnsmessage.Question{
			{
				Name:  q.Name,
				Type:  q.Type,
				Class: dnsmessage.ClassINET,
			},
		},
	}

	// PTR
	if q.Type == dnsmessage.TypePTR {
		return d.handlePtr(ctx, add, b, resp, d.resolver, q.Name)
	}

	// A or AAAA
	ips, ttl, err := d.resolver.Record(ctx, add, q.Type)
	if err != nil {
		if !errors.Is(err, ErrNoIPFound) && !errors.Is(err, ErrCondEmptyResponse) {
			if errors.Is(err, proxy.ErrBlocked) {
				log.Debug(err.Error())
			} else {
				log.Error("lookup domain failed", slog.String("domain", q.Name.String()), slog.Any("err", err))
			}
		}

		if !errors.Is(err, ErrNoIPFound) {
			resp.RCode = dnsmessage.RCodeNameError
		}
	}

	for _, a := range ips {
		var resource dnsmessage.ResourceBody
		if q.Type == dnsmessage.TypeA {
			rr := &dnsmessage.AResource{}
			copy(rr.A[:], a.To4())
			resource = rr
		} else {
			rr := &dnsmessage.AAAAResource{}
			copy(rr.AAAA[:], a.To16())
			resource = rr
		}
		resp.Answers = append(resp.Answers, dnsmessage.Resource{
			Header: dnsmessage.ResourceHeader{
				Name:  q.Name,
				Type:  q.Type,
				Class: dnsmessage.ClassINET,
				TTL:   ttl,
			},
			Body: resource,
		})
	}

	return resp.Pack()
}

func (d *dnsServer) handlePtr(ctx context.Context, address string, raw []byte, msg dnsmessage.Message,
	processor dns.DNS, name dnsmessage.Name) ([]byte, error) {
	if ff, ok := processor.(interface{ LookupPtr(string) (string, error) }); ok {
		r, err := ff.LookupPtr(name.String())
		if err == nil {
			msg.Answers = []dnsmessage.Resource{
				{
					Header: dnsmessage.ResourceHeader{
						Name:  name,
						Class: dnsmessage.ClassINET,
						TTL:   600,
					},
					Body: &dnsmessage.PTRResource{
						PTR: dnsmessage.MustNewName(r + "."),
					},
				},
			}

			return msg.Pack()
		}
	}

	return processor.Do(ctx, address, raw)
}
