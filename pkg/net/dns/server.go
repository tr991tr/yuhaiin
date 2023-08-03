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
	proxy "github.com/Asutorufa/yuhaiin/pkg/net/interfaces"
	"github.com/Asutorufa/yuhaiin/pkg/net/nat"
	"github.com/Asutorufa/yuhaiin/pkg/utils/pool"
	"github.com/Asutorufa/yuhaiin/pkg/utils/system"
	"golang.org/x/exp/slog"
	"golang.org/x/net/dns/dnsmessage"
)

type dnsServer struct {
	server      string
	resolver    proxy.Resolver
	listener    net.PacketConn
	tcpListener net.Listener

	reqChan   chan dnsRequest
	doneCtx   context.Context
	cancelCtx func()
}

type dnsRequest struct {
	payload   []byte
	writeBack func([]byte) error
}

func NewDnsServer(server string, process proxy.Resolver) proxy.DNSHandler {
	ctx, cancel := context.WithCancel(context.TODO())

	d := &dnsServer{
		server:    server,
		resolver:  process,
		reqChan:   make(chan dnsRequest, system.Procs),
		doneCtx:   ctx,
		cancelCtx: cancel,
	}

	do := func(req dnsRequest) error {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*7)
		defer cancel()

		data, err := d.handle(ctx, req.payload)
		if err != nil {
			return err
		}
		return req.writeBack(data)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case req := <-d.reqChan:
				go func() {
					if err := do(req); err != nil {
						log.Error("dns server handle failed", "err", err)
					}
				}()
			}

		}
	}()

	if server == "" {
		log.Warn("dns server is empty, skip to listen tcp and udp")
		return d
	}

	if err := d.startUDP(); err != nil {
		log.Error("start udp dns server failed", slog.Any("err", err))
	}

	go func() {
		if err := d.startTCP(); err != nil {
			log.Error("start tcp dns server failed", slog.Any("err", err))
		}
	}()

	return d
}

func (d *dnsServer) Close() error {
	d.cancelCtx()

	if d.listener != nil {
		d.listener.Close()
	}
	if d.tcpListener != nil {
		d.tcpListener.Close()
	}

	return nil
}

func (d *dnsServer) startUDP() (err error) {
	d.listener, err = dialer.ListenPacket("udp", d.server)
	if err != nil {
		return fmt.Errorf("dns udp server listen failed: %w", err)
	}

	log.Info("new udp dns server", "host", d.server)

	go func() {
		defer d.Close()
		for {
			buf := pool.GetBytes(nat.MaxSegmentSize)
			n, addr, err := d.listener.ReadFrom(buf)
			if err != nil {
				pool.PutBytes(buf)
				if e, ok := err.(net.Error); ok && e.Temporary() {
					continue
				}
				log.Error("dns udp server handle failed", "err", err)
				return
			}

			writeBack := func(b []byte) error {
				defer pool.PutBytes(buf)
				if _, err = d.listener.WriteTo(b, addr); err != nil {
					return fmt.Errorf("write dns response to client failed: %w", err)
				}
				return nil
			}

			if err = d.Do(context.TODO(), buf[:n], writeBack); err != nil {
				log.Error("dns server handle data failed", slog.Any("err", err))
			}

		}
	}()

	return nil
}

func (d *dnsServer) startTCP() (err error) {
	defer d.Close()

	d.tcpListener, err = dialer.ListenContext(context.TODO(), "tcp", d.server)
	if err != nil {
		return fmt.Errorf("dns tcp server listen failed: %w", err)
	}

	log.Info("new tcp dns server", "host", d.server)

	for {
		conn, err := d.tcpListener.Accept()
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Temporary() {
				continue
			}
			return fmt.Errorf("dns server accept failed: %w", err)
		}

		go func() {
			defer conn.Close()

			if err := d.HandleTCP(context.TODO(), conn); err != nil {
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

	n, err := io.ReadFull(c, data[:length])
	if err != nil {
		pool.PutBytes(data)
		return fmt.Errorf("dns server read data failed: %w", err)
	}

	return d.Do(ctx, data[:n], func(b []byte) error {
		defer pool.PutBytes(data)
		if err = binary.Write(c, binary.BigEndian, uint16(len(b))); err != nil {
			return fmt.Errorf("dns server write length failed: %w", err)
		}
		_, err = c.Write(b)
		return err
	})
}

func (d *dnsServer) HandleUDP(ctx context.Context, l net.PacketConn) error {
	buf := pool.GetBytes(nat.MaxSegmentSize)

	n, addr, err := l.ReadFrom(buf)
	if err != nil {
		pool.PutBytes(buf)
		return err
	}

	writeBack := func(b []byte) error {
		defer pool.PutBytes(buf)
		_, err = l.WriteTo(b, addr)
		return err
	}

	return d.Do(context.TODO(), buf[:n], writeBack)
}

func (d *dnsServer) Do(ctx context.Context, b []byte, writeBack func([]byte) error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-d.doneCtx.Done():
		return io.EOF
	case d.reqChan <- dnsRequest{payload: b, writeBack: writeBack}:
		return nil
	}
}

func (d *dnsServer) handle(ctx context.Context, raw []byte) ([]byte, error) {
	var parse dnsmessage.Parser
	header, err := parse.Start(raw)
	if err != nil {
		return nil, fmt.Errorf("dns server parse failed: %w", err)
	}

	question, err := parse.Question()
	if err != nil {
		return nil, fmt.Errorf("dns server parse failed: %w", err)
	}

	reqMsg := &reqMsg{header, question, raw}

	// PTR
	if question.Type == dnsmessage.TypePTR {
		return d.handlePtr(ctx, reqMsg)
	}

	// A or AAAA
	if question.Type == dnsmessage.TypeA || question.Type == dnsmessage.TypeAAAA {
		return d.handleAOrAAAA(ctx, reqMsg)
	}

	// other question Type
	log.Debug("other dns question Type", "type", question.Type)
	return d.resolver.Do(ctx, reqMsg.Addr(), raw)
}

type reqMsg struct {
	header   dnsmessage.Header
	question dnsmessage.Question
	raw      []byte
}

func (r *reqMsg) Addr() string { return strings.TrimSuffix(r.question.Name.String(), ".") }

func (r *reqMsg) newResponse(f ...func(*dnsmessage.Message)) *dnsmessage.Message {
	msg := &dnsmessage.Message{
		Header: dnsmessage.Header{
			ID:                 r.header.ID,
			Response:           true,
			Authoritative:      false,
			RecursionDesired:   false,
			RCode:              dnsmessage.RCodeSuccess,
			RecursionAvailable: false,
		},
		Questions: []dnsmessage.Question{
			{
				Name:  r.question.Name,
				Type:  r.question.Type,
				Class: dnsmessage.ClassINET,
			},
		},
	}

	for _, f := range f {
		f(msg)
	}

	return msg
}
func (d *dnsServer) handleAOrAAAA(ctx context.Context, reqMsg *reqMsg) ([]byte, error) {
	records, ttl, err := d.resolver.Record(ctx, reqMsg.Addr(), reqMsg.question.Type)
	if err != nil {
		noIPFound := errors.Is(err, ErrNoIPFound)

		if !noIPFound && !errors.Is(err, ErrCondEmptyResponse) {
			if errors.Is(err, proxy.ErrBlocked) {
				log.Debug(err.Error())
			} else {
				log.Error("lookup domain failed", slog.String("domain", reqMsg.question.Name.String()), slog.Any("err", err))
			}
		}

		if noIPFound {
			return reqMsg.newResponse().Pack()
		}

		return reqMsg.newResponse(func(m *dnsmessage.Message) { m.RCode = dnsmessage.RCodeNameError }).Pack()

	}

	msg := reqMsg.newResponse(func(m *dnsmessage.Message) {
		m.Answers = make([]dnsmessage.Resource, 0, len(records))

		for _, ip := range records {
			answer := dnsmessage.Resource{
				Header: dnsmessage.ResourceHeader{
					Name:  reqMsg.question.Name,
					Type:  reqMsg.question.Type,
					Class: dnsmessage.ClassINET,
					TTL:   ttl,
				},
			}

			if reqMsg.question.Type == dnsmessage.TypeA {
				answer.Body = &dnsmessage.AResource{A: [4]byte(ip.To4())}
			} else {
				answer.Body = &dnsmessage.AAAAResource{AAAA: [16]byte(ip.To16())}
			}

			m.Answers = append(m.Answers, answer)
		}
	})

	return msg.Pack()
}

func (d *dnsServer) handlePtr(ctx context.Context, req *reqMsg) ([]byte, error) {

	ff, ok := d.resolver.(interface{ LookupPtr(string) (string, error) })
	if ok {
		r, err := ff.LookupPtr(req.question.Name.String())
		if err == nil {
			msg := req.newResponse(func(m *dnsmessage.Message) {
				m.Answers = []dnsmessage.Resource{
					{
						Header: dnsmessage.ResourceHeader{
							Name:  req.question.Name,
							Class: dnsmessage.ClassINET,
							TTL:   600,
						},
						Body: &dnsmessage.PTRResource{
							PTR: dnsmessage.MustNewName(r + "."),
						},
					},
				}
			})

			return msg.Pack()
		}
	}

	return d.resolver.Do(ctx, req.Addr(), req.raw)
}
