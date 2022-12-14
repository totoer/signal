package sip

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type RawHeader struct {
	Value      string
	Properties map[string]string
}

func prepareHeader(line string) (string, []RawHeader) {
	keyBuffer := make([]string, 0)
	i := 0

	for ii, r := range line {
		s := string(r)
		if s != ":" {
			keyBuffer = append(keyBuffer, s)
		} else {
			i = ii + 1
			break
		}
	}

	key := strings.TrimSpace(strings.Join(keyBuffer, ""))
	value := strings.TrimSpace(line[i:])

	rhs := make([]RawHeader, 0)

	// Authorization: Digest username="Alice", realm="atlanta.com", nonce="84a4cc6f3082121f32b42a2187831a9e", response="7587245234b3434cc3412213e5f113a5432"
	// WWW-Authenticate: Digest realm="atlanta.com", nonce="f84f1cec41e6cbe5aea9c8e88d359", algorithm=MD5
	if key == "Authorization" || key == "WWW-Authenticate" {
		rawValues := strings.Split(value, "Digest ")
		parts := strings.Split(rawValues[0], ",")
		props := make(map[string]string)
		for _, part := range parts {
			propParts := strings.Split(part, "=")
			propKey := propParts[0]
			propValue := propParts[1]
			props[propKey] = strings.Trim(propValue, "\" ")
		}
		rhs = append(rhs, RawHeader{
			Properties: props,
		})
		return key, rhs
	} else {
		parts := strings.Split(value, ",")
		for _, part := range parts {
			rawl := strings.Split(part, ";")
			rh := RawHeader{
				Value: rawl[0],
			}
			if len(rawl[1:]) > 0 {
				props := make(map[string]string)
				for _, prop := range rawl[1:] {
					if strings.Count(prop, "=") > 0 {
						propParts := strings.Split(prop, "=")
						propKey := propParts[0]
						propValue := strings.Trim(propParts[1], " \"")
						props[propKey] = propValue
					} else {
						props[prop] = ""
					}
				}
				rh.Properties = props
			}
			rhs = append(rhs, rh)
		}

		return key, rhs
	}
}

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

func DecodeURI(v string) (URI, error) {
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
func DecodeTarget(v string) (Address, error) {
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

	if uri, err := DecodeURI(rawURI); err != nil {
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

func decodePlainHeader(rh RawHeader) (PlainHeader, error) {
	return PlainHeader{rh.Value, rh.Properties}, nil
}

type IntegerHeader struct {
	Value      int
	Properties map[string]string
}

func (ih IntegerHeader) String() string {
	return strconv.Itoa(int(ih.Value))
}

func decodeIntegerHeader(rh RawHeader) (IntegerHeader, error) {
	if v, err := strconv.Atoi(rh.Value); err != nil {
		return IntegerHeader{}, err
	} else {
		return IntegerHeader{v, rh.Properties}, nil
	}
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

func decodeDestinations(rh RawHeader) (Destination, error) {
	if t, err := DecodeTarget(rh.Value); err != nil {
		return Destination{}, err
	} else {
		return Destination{
			Address: t,
			Tag:     rh.Properties["tag"],
		}, nil
	}
}

type Via struct {
	Host     string
	Branch   string
	Received string
	Rport    bool
}

func (via Via) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("SIP/2.0/UDP %s", via.Host))
	if via.Branch != "" {
		builder.WriteString(fmt.Sprintf(";branch=%s", via.Branch))
	}
	if via.Received != "" {
		builder.WriteString(fmt.Sprintf(";received=%s", via.Received))
	}
	if via.Rport {
		builder.WriteString(";rpotr")
	}
	return builder.String()
}

func decodeVia(rh RawHeader) (Via, error) {
	via := Via{
		Host: strings.TrimLeft(rh.Value, "Via: SIP/2.0/UDP "),
	}

	if branch, ok := rh.Properties["branch"]; ok {
		via.Branch = branch
	}

	if received, ok := rh.Properties["received"]; ok {
		via.Received = received
	}

	if _, ok := rh.Properties["rport"]; ok {
		via.Rport = true
	}

	return via, nil
}

type Contact struct {
	Address Address
	Q       float64
	Expires int
}

func (c Contact) String() string {
	var builder strings.Builder
	builder.WriteString(c.Address.String())
	if c.Q != 0 {
		builder.WriteString(fmt.Sprintf("%f", c.Q))
	}
	if c.Expires != 0 {
		builder.WriteString(strconv.Itoa(c.Expires))
	}
	return builder.String()
}

// Contact: "Mr. Watson" <sip:watson@worcester.bell-telephone.com>;q=0.7; expires=3600
func decodeContact(rh RawHeader) (Contact, error) {
	var contact Contact
	if address, err := DecodeTarget(rh.Value); err != nil {
		return Contact{}, err
	} else {
		contact.Address = address
	}
	if rawQ, ok := rh.Properties["q"]; ok {
		if q, err := strconv.ParseFloat(rawQ, 32); err != nil {
			return Contact{}, err
		} else {
			contact.Q = q
		}
	}
	return contact, nil
}

type CSeq struct {
	Value  int
	Method MethodType
}

func (cseq CSeq) String() string {
	var builder strings.Builder
	builder.WriteString(strconv.Itoa(cseq.Value))
	builder.WriteString(" ")
	builder.WriteString(string(cseq.Method))
	return builder.String()
}

func decodeCSeq(rh RawHeader) (CSeq, error) {
	parts := strings.Split(rh.Value, " ")
	if v, err := strconv.Atoi(parts[0]); err != nil {
		return CSeq{}, err
	} else {
		return CSeq{
			Value:  v,
			Method: MethodType(parts[1]),
		}, nil
	}
}

type Allow MethodType

func (a Allow) String() string {
	return string(a)
}

// Allow: INVITE, ACK, OPTIONS, CANCEL, BYE
func decodeAllow(rh RawHeader) (Allow, error) {
	v := strings.Trim(rh.Value, " ")
	return Allow(v), nil
}

type Authorization struct {
	Username string
	Realm    string
	Nonce    string
	Response string
}

func (a *Authorization) String() string {
	var builder strings.Builder
	builder.WriteString("Digest ")
	builder.WriteString(fmt.Sprintf("username=%s", a.Username))
	builder.WriteString(fmt.Sprintf("realm=%s", a.Realm))
	builder.WriteString(fmt.Sprintf("nonce=%s", a.Nonce))
	builder.WriteString(fmt.Sprintf("response=%s", a.Response))
	return builder.String()
}

// Authorization: Digest username="Alice", realm="atlanta.com", nonce="84a4cc6f3082121f32b42a2187831a9e", response="7587245234b3434cc3412213e5f113a5432"
func decodeAuthorization(rh RawHeader) (Authorization, error) {
	return Authorization{
		Username: rh.Properties["username"],
		Realm:    rh.Properties["realm"],
		Nonce:    rh.Properties["nonce"],
		Response: rh.Properties["response"],
	}, nil
}

type WWWAuthenticate struct {
	Realm     string
	Nonce     string
	Algorithm string
}

func (a *WWWAuthenticate) String() string {
	return fmt.Sprintf("Digest realm=\"%s\", nonce=\"%s\", algorithm=\"%s\"", a.Realm, a.Nonce, a.Algorithm)
}

// WWW-Authenticate: Digest realm="atlanta.com", nonce="f84f1cec41e6cbe5aea9c8e88d359", algorithm=MD5
func decodeWWWAuthenticate(rh RawHeader) (WWWAuthenticate, error) {
	return WWWAuthenticate{
		Realm:     rh.Properties["realm"],
		Nonce:     rh.Properties["nonce"],
		Algorithm: rh.Properties["algorithm"],
	}, nil
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
}

func (hs *Headers) Encode() []byte {
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

func DecodeHeaders(lines []string) (*Headers, error) {
	hs := &Headers{
		Vias: make([]Via, 0),
	}
	for _, line := range lines {
		key, rhs := prepareHeader(line)
		for _, rh := range rhs {
			switch key {
			case "Via":
				if via, err := decodeVia(rh); err != nil {
					hs.Vias = append(hs.Vias, via)
				}
			case "From", "To":
				if dist, err := decodeDestinations(rh); err != nil {
					return nil, err
				} else {
					switch key {
					case "From":
						hs.From = &dist
					case "To":
						hs.To = &dist
					}
				}
			case "Max-Forwards", "Content-Length":
				if h, err := decodeIntegerHeader(rh); err != nil {
					return nil, err
				} else {
					switch key {
					case "Max-Forwards":
						hs.MaxForwards = &h
					case "Content-Length":
						hs.ContentLength = &h
					}
				}
			case "Call-ID":
				if h, err := decodePlainHeader(rh); err != nil {
					return nil, err
				} else {
					hs.CallID = &h
				}
			case "Contact":
				if contact, err := decodeContact(rh); err != nil {
					return nil, err
				} else {
					hs.Contacts = append(hs.Contacts, contact)
				}
			case "CSeq":
				if cseq, err := decodeCSeq(rh); err != nil {
					return nil, err
				} else {
					hs.CSeq = &cseq
				}
			case "Allow":
				if allow, err := decodeAllow(rh); err != nil {
					return nil, err
				} else {
					hs.Allows = append(hs.Allows, allow)
				}
			case "Authorization":
				if auth, err := decodeAuthorization(rh); err != nil {
					return nil, err
				} else {
					hs.Authorization = &auth
				}
			case "WWW-Authenticate":
				if wwwauth, err := decodeWWWAuthenticate(rh); err != nil {
					return nil, err
				} else {
					hs.WWWAuthenticate = &wwwauth
				}
			}
		}
	}
	return hs, nil
}
