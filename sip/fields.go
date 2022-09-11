package sip

import (
	"fmt"
	"strconv"
	"strings"
)

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

func (hs *Headers) parseVias() ([]Via, error) {
	vias := make([]Via, 0)

	if lines, err := hs.GetRawFields("Via"); err != nil {
		return vias, err
	} else {
		for _, line := range lines {
			host := strings.TrimLeft(line.Value, "Via: SIP/2.0/UDP ")
			v := Via{
				Host: host,
			}

			if branch, ok := line.Properties["branch"]; ok {
				v.Branch = branch
			}

			if received, ok := line.Properties["received"]; ok {
				v.Received = received
			}

			if _, ok := line.Properties["rport"]; ok {
				v.Rport = true
			}

			vias = append(vias, v)
		}
	}

	return vias, nil
}

func (hs *Headers) GetVias() ([]Via, error) {
	if hs.Vias == nil {
		if vias, err := hs.parseVias(); err != nil {
			return nil, err
		} else {
			hs.Vias = vias
		}
	}

	return hs.Vias, nil
}

func (hs *Headers) PopVia() (Via, error) {
	if vias, err := hs.GetVias(); err != nil {
		return Via{}, err
	} else {
		lastVia := vias[len(vias)-1]

		return lastVia, nil
	}
}

func (hs *Headers) PushVia(via Via) {
	if hs.Vias == nil {
		hs.Vias = make([]Via, 0)
	}
	hs.Vias = append(hs.Vias, via)
}

func (hs *Headers) GetFrom() (Destination, error) {
	if hs.From == nil {
		if from, err := hs.parseDestination("From"); err != nil {
			return Destination{}, nil
		} else {
			hs.From = &from
		}
	}
	return *hs.From, nil
}

func (hs Headers) GetHostLoginByFrom() (string, string, error) {
	if from, err := hs.GetFrom(); err != nil {
		return "", "", err
	} else {
		host := from.Address.URI.Host
		login := from.Address.URI.Login
		return host, login, nil
	}
}

func (hs *Headers) GetTo() (Destination, error) {
	if hs.To == nil {
		if to, err := hs.parseDestination("To"); err != nil {
			return Destination{}, nil
		} else {
			hs.To = &to
		}
	}
	return *hs.To, nil
}

func (hs Headers) GetHostLoginByTo() (string, string, error) {
	if to, err := hs.GetTo(); err != nil {
		return "", "", err
	} else {
		host := to.Address.URI.Host
		login := to.Address.URI.Login
		return host, login, nil
	}
}

func (hs *Headers) GetMaxForwards() (IntegerHeader, error) {
	if hs.MaxForwards == nil {
		if f, err := hs.parseIntegerHeader("Max-Forwards"); err != nil {
			return IntegerHeader{}, err
		} else {
			hs.MaxForwards = &f[0]
		}
	}

	return *hs.MaxForwards, nil
}

func (hs *Headers) GetCallID() (string, error) {
	if hs.CallID == nil {
		if f, err := hs.parsePlainHeader("Call-ID"); err != nil {
			return "", err
		} else {
			hs.CallID = &f[0]
		}
	}

	return hs.CallID.Value, nil
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
func (hs *Headers) parseContacts() ([]Contact, error) {
	if lines, err := hs.GetRawFields("Contact"); err != nil {
		return nil, err
	} else {
		contacts := make([]Contact, 0)
		for _, line := range lines {
			var contact Contact
			if address, err := ParseTarget(line.Value); err != nil {
				return nil, err
			} else {
				contact.Address = address
			}
			if rawQ, ok := line.Properties["q"]; ok {
				if q, err := strconv.ParseFloat(rawQ, 32); err != nil {
					return nil, err
				} else {
					contact.Q = q
				}
			}

			contacts = append(contacts, contact)
		}
		return contacts, nil
	}
}

func (hs *Headers) GetContacts() ([]Contact, error) {
	if hs.Contacts == nil {
		if contacts, err := hs.parseContacts(); err != nil {
			return nil, err
		} else {
			hs.Contacts = contacts
		}
	}

	return hs.Contacts, nil
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

func (hs *Headers) parseCSeq() (CSeq, error) {
	if lines, err := hs.GetRawFields("CSeq"); err != nil {
		return CSeq{}, err
	} else {
		parts := strings.Split(lines[0].Value, " ")
		if v, err := strconv.Atoi(parts[0]); err != nil {
			return CSeq{}, err
		} else {
			return CSeq{
				Value:  v,
				Method: MethodType(parts[1]),
			}, nil
		}
	}
}

func (hs *Headers) GetCSeq() (CSeq, error) {
	if hs.CSeq == nil {
		if cseq, err := hs.parseCSeq(); err != nil {
			return CSeq{}, err
		} else {
			hs.CSeq = &cseq
		}
	}

	return *hs.CSeq, nil
}

type Allow MethodType

func (a Allow) String() string {
	return string(a)
}

// Allow: INVITE, ACK, OPTIONS, CANCEL, BYE
func (hs *Headers) parseAllows() ([]Allow, error) {
	if lines, err := hs.GetRawFields("Allow"); err != nil {
		return nil, err
	} else {
		allows := make([]Allow, 0)
		for _, line := range lines {
			v := strings.Trim(line.Value, " ")
			allows = append(allows, Allow(v))
		}
		return allows, nil
	}
}

func (hs *Headers) GetAllows() ([]Allow, error) {
	if hs.Allows == nil {
		if allows, err := hs.parseAllows(); err != nil {
			return nil, err
		} else {
			hs.Allows = allows
		}
	}

	return hs.Allows, nil
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
func (hs *Headers) parseAuthorization() (Authorization, error) {
	if lines, err := hs.GetRawFields("Authorization"); err != nil {
		return Authorization{}, err
	} else {
		l := lines[0]
		return Authorization{
			Username: l.Properties["username"],
			Realm:    l.Properties["realm"],
			Nonce:    l.Properties["nonce"],
			Response: l.Properties["response"],
		}, nil
	}
}

func (hs *Headers) GetAuthorization() (Authorization, error) {
	if hs.Authorization == nil {
		if authorization, err := hs.parseAuthorization(); err != nil {
			return Authorization{}, err
		} else {
			hs.Authorization = &authorization
		}
	}

	return *hs.Authorization, nil
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
