package transport

import (
	"errors"
	"net"
	"signal/sip"
)

type TransportType string

var ErrConnectionDoesNotExists error = errors.New("connection does not exists")

type Transport interface {
	Run(chan sip.Message)
	Send(string, []byte) error
	SendSIP(sip.Message) error
	// Close(net.Addr)
}

type Connection interface {
	Listen(chan sip.Message)
	Send(net.Addr, []byte) error
	Close()
}
