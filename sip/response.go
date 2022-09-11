package sip

import (
	"bytes"
	"fmt"
	"net"
)

type Response struct {
	Code         ResponseCode
	Headers      Headers
	SDP          SDP
	SourceAddres net.Addr
	rawBody      string
}

func (resp Response) GetHeaders() *Headers {
	return &resp.Headers
}

func (resp Response) Data() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("SIP/2.0 %d %s", resp.Code, ResponseCodes[int(resp.Code)]))
	buffer.WriteString("\r\n")
	buffer.Write(resp.Headers.Data())
	buffer.WriteString("\r\n")

	return buffer.Bytes()
}

func (resp Response) GetRawBody() string {
	return resp.rawBody
}

func (resp Response) GetSourceAddres() net.Addr {
	return resp.SourceAddres
}

func (resp Response) GetRecipientAddrs() []string {
	if contacts, err := resp.GetHeaders().GetContacts(); err == nil {
		addrs := make([]string, 0)
		for _, contact := range contacts {
			addrs = append(addrs, contact.Address.URI.Host)
		}
		return addrs
	} else {
		return []string{resp.GetSourceAddres().String()}
	}
}

func NewResponse(c ResponseCode, b string, h Headers) Response {
	return Response{
		Code:    c,
		Headers: h,
		rawBody: b,
	}
}
