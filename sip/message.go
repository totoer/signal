package sip

import (
	"net"

	"github.com/spf13/viper"
)

var PROTOCOL string
var VERSION string

func init() {
	VERSION = "SIP/2.0"
	PROTOCOL = viper.GetString("server.transport")
}

type Message interface {
	GetRawBody() string
	Data() []byte
	GetHeaders() *Headers
	GetSourceAddres() net.Addr
	GetRecipientAddrs() []string
}

func NewMessage(p *Parser, addr net.Addr) (Message, error) {
	if m, err := p.Parse(); err != nil {
		return nil, err
	} else {
		switch m.(type) {
		case Request:
			r := m.(Request)
			r.SourceAddres = addr
			return r, nil
		case Response:
			r := m.(Response)
			r.SourceAddres = addr
			return r, nil
		default:
			return nil, nil
		}
	}
}
