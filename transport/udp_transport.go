package transport

import (
	"net"
	"strconv"

	"signal/sip"

	"github.com/rs/zerolog/log"
)

type UDPTransport struct {
	conn    *net.UDPConn
	clients []*net.UDPConn
	laddr   *net.UDPAddr
}

func (t *UDPTransport) Run(mq chan sip.Message) {
	if conn, err := net.ListenUDP("udp", t.laddr); err != nil {
		log.Err(err)
	} else {
		t.conn = conn
		defer t.conn.Close()

		for {
			buffer := make([]byte, 1024)
			if l, addr, err := t.conn.ReadFrom(buffer); err != nil {
				continue
			} else {
				body := string(buffer[:l])
				// log.Debug().Str("body", body).Msg("Receivd message")
				p := sip.NewParser(body)
				if m, err := sip.NewMessage(p, addr); err != nil {
					log.Error().Err(err).Str("transport", "recived").Err(err).Msg(body)
				} else {
					mq <- m
				}
			}
		}
	}
}

func (t *UDPTransport) Send(rawAddr string, body []byte) error {
	if rawHost, rawPort, err := net.SplitHostPort(rawAddr); err != nil {
		return err
	} else if port, err := strconv.Atoi(rawPort); err != nil {
		return err
	} else {
		addr := net.UDPAddr{
			IP:   net.ParseIP(rawHost),
			Port: port,
		}
		// log.Debug().Bytes("body", body).Msg("Send message")
		_, err := t.conn.WriteToUDP(body, &addr)
		return err
	}
}

func (t *UDPTransport) SendSIP(m sip.Message) error {
	// body := m.Data()
	// addrs := m.GetRecipientAddrs()
	// for _, addr := range addrs {
	// 	return t.Send(addr, body)
	// }
	// return nil
	return t.Send(m.GetSourceAddres().String(), m.Data())
}

func NewUDPTransport(ip string, port int) *UDPTransport {
	return &UDPTransport{
		clients: make([]*net.UDPConn, 0),
		laddr: &net.UDPAddr{
			IP:   net.ParseIP(ip),
			Port: port,
		},
	}
}
