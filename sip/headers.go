package sip

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type URI struct {
	Login     string `json:"login"`
	Host      string `json:"host"`
	Transport string `json:"target"`
	LR        bool   `json:"lr"`
}

func (uri URI) String() string {
	var builder strings.Builder
	builder.WriteString("sip:")
	if uri.Login != "" {
		builder.WriteString(uri.Login)
		builder.WriteString("@")
	}
	builder.WriteString(uri.Host)
	if uri.LR {
		builder.WriteString(";lr")
	}
	if uri.Transport != "" {
		builder.WriteString(fmt.Sprintf(";transport=%s", uri.Transport))
	}
	return builder.String()
}

func (uri URI) GetAddr() net.Addr {
	return &net.UDPAddr{}
}

func ParseURI(v string) (URI, error) {
	uri := URI{
		LR: false,
	}

	raw := strings.Trim(v, "<>")
	raw = strings.TrimLeft(raw, "sip:")
	parts := strings.Split(raw, ";")
	rawURI := parts[0]
	if strings.Index(rawURI, "@") != -1 {
		rawURIParts := strings.Split(rawURI, "@")
		uri.Login = rawURIParts[0]
		uri.Host = rawURIParts[1]
	} else {
		uri.Host = rawURI
	}

	for _, p := range parts[1:] {
		if strings.Index(p, "=") != -1 {
			kv := strings.Split(p, "=")
			if kv[0] == "transport" {
				uri.Transport = kv[1]
			}
		} else if p == "lr" {
			uri.LR = true
		}
	}

	return uri, nil
}

type Address struct {
	Name string `json"name"`
	URI  URI    `json:"uri"`
}

func (t Address) String() string {
	var builder strings.Builder
	if t.Name != "" {
		v := fmt.Sprintf("\"%s\"", t.Name)
		builder.WriteString(fmt.Sprintf("%s ", v))
	}
	builder.WriteString(fmt.Sprintf("<%s>", t.URI))
	return builder.String()
}

// "A. G. Bell" <sip:agb@bell-telephone.com>
// The Operator <sip:operator@cs.columbia.edu>
// "Bob" <sips:bob@biloxi.com>
// sip:+12125551212@phone2net.com
// Anonymous <sip:c8oqz84zk7z@privacy.org>
func ParseTarget(v string) (Address, error) {
	var address Address
	i := strings.Index(v, "<sip:")
	if i == -1 {
		i = strings.Index(v, "sip:")
	}
	var rawURI string

	if i == 0 {
		rawURI = v
	} else if i != -1 {
		address.Name = strings.Trim(v[:i-1], " \"")
		rawURI = v[i:]
	}

	if uri, err := ParseURI(rawURI); err != nil {
		return address, err
	} else {
		address.URI = uri
	}

	return address, nil
}

type PlainHeader struct {
	Value      string
	Properties map[string]string
}

func (ph PlainHeader) String() string {
	return string(ph.Value)
}

type IntegerHeader struct {
	Value      int
	Properties map[string]string
}

func (ih IntegerHeader) String() string {
	return strconv.Itoa(int(ih.Value))
}

type Destination struct {
	Address Address `json:"addres"`
	Tag     string
}

func (d Destination) String() string {
	var builder strings.Builder
	builder.WriteString(d.Address.String())
	if d.Tag != "" {
		builder.WriteString(fmt.Sprintf(";tag=%s", d.Tag))
	}
	return builder.String()
}

type Headers struct {
	Vias            []Via
	From            *Destination
	To              *Destination
	CallID          *PlainHeader
	Contacts        []Contact
	CSeq            *CSeq
	Allows          []Allow
	MaxForwards     *IntegerHeader
	WWWAuthenticate *WWWAuthenticate
	Authorization   *Authorization
	ContentLength   *IntegerHeader
	Lines           map[string][]Line
	p               *Parser
}

var ErrHeaderFieldNotExists error = errors.New("header field not exists")

func (hs *Headers) GetRawFields(key string) ([]Line, error) {
	if hs.Lines == nil && hs.p != nil {
		hs.Lines = hs.p.PrepareFields()
	} else if hs.p == nil {
		hs.Lines = make(map[string][]Line)
	}

	if fields, ok := hs.Lines[key]; ok {
		return fields, nil
	}
	return nil, ErrHeaderFieldNotExists
}

func (hs *Headers) parsePlainHeader(key string) ([]PlainHeader, error) {
	if lines, err := hs.GetRawFields(key); err != nil {
		return nil, err
	} else {
		headers := make([]PlainHeader, 0)
		for _, line := range lines {
			headers = append(headers, PlainHeader{line.Value, line.Properties})
		}
		return headers, nil
	}
}

func (hs *Headers) GetPlainHeaders(key string) ([]PlainHeader, error) {
	return hs.parsePlainHeader(key)
}

func (hs *Headers) parseIntegerHeader(key string) ([]IntegerHeader, error) {
	if lines, err := hs.GetRawFields(key); err != nil {
		return nil, err
	} else {
		headers := make([]IntegerHeader, 0)
		for _, line := range lines {
			if v, err := strconv.Atoi(line.Value); err != nil {
				return nil, err
			} else {
				headers = append(headers, IntegerHeader{v, line.Properties})
			}
		}
		return headers, nil
	}
}

func (hs *Headers) GetIntegerHeaders(key string) ([]IntegerHeader, error) {
	return hs.parseIntegerHeader(key)
}

func (hs *Headers) parseDestination(key string) (Destination, error) {
	if lines, err := hs.GetRawFields(key); err != nil {
		return Destination{}, err
	} else {
		if t, err := ParseTarget(lines[0].Value); err != nil {
			return Destination{}, err
		} else {
			return Destination{
				Address: t,
				Tag:     lines[0].Properties["tag"],
			}, nil
		}
	}
}

func (hs *Headers) Data() []byte {
	var buffer bytes.Buffer

	if len(hs.Vias) != 0 {
		for _, via := range hs.Vias {
			buffer.WriteString("Via: ")
			buffer.WriteString(via.String())
			buffer.WriteString("\r\n")
		}
	}

	if hs.From != nil {
		buffer.WriteString("From: ")
		buffer.WriteString(hs.From.String())
		buffer.WriteString("\r\n")
	}

	if hs.To != nil {
		buffer.WriteString("To: ")
		buffer.WriteString(hs.To.String())
		buffer.WriteString("\r\n")
	}

	if hs.CallID != nil {
		buffer.WriteString("Call-ID: ")
		buffer.WriteString(hs.CallID.String())
		buffer.WriteString("\r\n")
	}

	if hs.Contacts != nil {
		for _, contact := range hs.Contacts {
			buffer.WriteString("Contact: ")
			buffer.WriteString(contact.String())
			buffer.WriteString("\r\n")
		}
	}

	if hs.CSeq != nil {
		buffer.WriteString("CSeq: ")
		buffer.WriteString(hs.CSeq.String())
		buffer.WriteString("\r\n")
	}

	if hs.Allows != nil {
		rawAllows := make([]string, 0)
		for _, allow := range hs.Allows {
			rawAllows = append(rawAllows, allow.String())
		}
		buffer.WriteString("Allow: ")
		buffer.WriteString(strings.Join(rawAllows, ","))
		buffer.WriteString("\r\n")
	}

	if hs.MaxForwards != nil {
		buffer.WriteString("Max-Forwards: ")
		buffer.WriteString(hs.MaxForwards.String())
		buffer.WriteString("\r\n")
	}

	if hs.WWWAuthenticate != nil {
		buffer.WriteString("WWW-Authenticate: ")
		buffer.WriteString(hs.WWWAuthenticate.String())
		buffer.WriteString("\r\n")
	}

	if hs.Authorization != nil {
		buffer.WriteString("Authorization: ")
		buffer.WriteString(hs.Authorization.String())
		buffer.WriteString("\r\n")
	}

	hs.ContentLength = &IntegerHeader{
		Value: 0,
	}
	buffer.WriteString("Content-Length: ")
	buffer.WriteString(hs.ContentLength.String())
	buffer.WriteString("\r\n")

	return buffer.Bytes()
}

func NewHeaders(p *Parser) Headers {
	return Headers{
		Lines: nil,
		p:     p,
	}
}
