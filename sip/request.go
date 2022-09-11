package sip

import (
	"bytes"
	"fmt"
	"net"
)

type MethodType string

const (
	INVITE   MethodType = "INVITE"
	ACK      MethodType = "ACK"
	BYE      MethodType = "BYE"
	CANCEL   MethodType = "CANCEL"
	REGISTER MethodType = "REGISTER"
	OPTIONS  MethodType = "OPTIONS"
	INFO     MethodType = "INFO"
)

func (mt *MethodType) IncludeIn(mts ...MethodType) bool {
	for _, i := range mts {
		if i == *mt {
			return true
		}
	}
	return false
}

type Request struct {
	Method       MethodType
	URI          URI
	Headers      Headers
	SDP          SDP
	SourceAddres net.Addr
	rawBody      string
}

func (req Request) GetHeaders() *Headers {
	return &req.Headers
}

func (req Request) Data() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s %s SIP/2.0", req.Method, req.URI.String()))
	buffer.WriteString("\r\n")
	buffer.Write(req.Headers.Data())
	buffer.WriteString("\r\n")

	return buffer.Bytes()
}

func (req Request) GetRawBody() string {
	return req.rawBody
}

func (req Request) GetSourceAddres() net.Addr {
	return req.SourceAddres
}

func (req Request) GetRecipientAddrs() []string {
	if contacts, err := req.GetHeaders().GetContacts(); err == nil {
		addrs := make([]string, 0)
		for _, contact := range contacts {
			addrs = append(addrs, contact.Address.URI.Host)
		}
		return addrs
	} else {
		return []string{req.GetSourceAddres().String()}
	}
}

func (req Request) MakeResponse(c ResponseCode) (Response, error) {
	if vias, err := req.GetHeaders().GetVias(); err != nil {
		return Response{}, err
	} else if cid, err := req.GetHeaders().GetCallID(); err != nil {
		return Response{}, err
	} else if cseq, err := req.GetHeaders().GetCSeq(); err != nil {
		return Response{}, err
	} else if from, err := req.GetHeaders().GetFrom(); err != nil {
		return Response{}, err
	} else if to, err := req.GetHeaders().GetTo(); err != nil {
		return Response{}, err
	} else {
		r := NewResponse(c, "", Headers{})
		r.Headers.Vias = vias
		r.Headers.CallID = &PlainHeader{
			Value: cid,
		}
		r.Headers.From = &from
		r.Headers.To = &to
		r.Headers.CSeq = &CSeq{
			Value:  cseq.Value + 1,
			Method: req.Method,
		}
		r.SourceAddres = req.GetSourceAddres()
		return r, nil
	}
}

func NewRequest(m MethodType, b string, uri URI, h Headers) Request {
	return Request{
		Method:  m,
		URI:     uri,
		Headers: h,
		rawBody: b,
	}
}
