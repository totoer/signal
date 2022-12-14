package sip

import (
	"errors"
	"strconv"
	"strings"
)

var ErrCantParseMessage = errors.New("cant parse message")

var HEADERS_ORDER = []string{
	"Via",
}

type Line struct {
	Value      string
	Properties map[string]string
}

func ParseLine(v string) []Line {
	lines := make([]Line, 0)

	parts := strings.Split(v, ",")
	for _, p := range parts {
		rawl := strings.Split(p, ";")
		l := Line{
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
			l.Properties = props
		}
		lines = append(lines, l)
	}
	return lines
}

func ParseAuthorizationLine(v string) []Line {
	value := strings.Split(v, "Digest")
	parts := strings.Split(value[1], ",")
	props := make(map[string]string)
	l := Line{
		Value: "",
	}
	for _, p := range parts {
		if strings.Count(p, "=") > 0 {
			propParts := strings.Split(p, "=")
			propKey := propParts[0]
			propValue := strings.Trim(propParts[1], " \"")
			props[propKey] = propValue
		}
		l.Properties = props

	}
	return []Line{l}
}

type Parser struct {
	rawBody  string
	rawLines []string
}

func (p *Parser) getFirstLine() string {
	return p.rawLines[0]
}

func (p *Parser) IsRequest() bool {
	parts := strings.Split(p.getFirstLine(), " ")
	return parts[0] != "SIP/2.0"
}

func (p *Parser) parseRequestLine() (MethodType, URI, error) {
	fl := p.getFirstLine()
	parts := strings.Split(fl, " ")
	if len(parts) != 3 {
		return INVITE, URI{}, ErrCantParseMessage
	}
	m := strings.TrimSpace(parts[0])
	rawURI := strings.TrimSpace(parts[1])
	if uri, err := ParseURI(rawURI); err != nil {
		return INVITE, URI{}, err
	} else {
		return MethodType(m), uri, nil
	}
}

func (p *Parser) PrepareFields() map[string][]Line {
	fields := make(map[string][]Line)
	for _, rawLine := range p.rawLines[1:] {
		if rawLine != "" && strings.Index(rawLine, "\n") == -1 {
			lineKey, value := separateHeaderLine(rawLine)

			var lines = make([]Line, 0)
			if lineKey != "Authorization" && lineKey != "WWW-Authenticate" {
				lines = ParseLine(strings.TrimSpace(value))
			} else {
				lines = ParseAuthorizationLine(strings.TrimSpace(value))
			}

			if len(lines) > 0 {
				if fields[lineKey] == nil {
					fields[lineKey] = lines
				} else {
					fields[lineKey] = append(fields[lineKey], lines...)
				}
			}
		}
	}

	return fields
}

var ErrIsNotRequest = errors.New("is not request")

func (p *Parser) ParseRequest() (Request, error) {
	if !p.IsRequest() {
		return Request{}, ErrIsNotRequest
	}

	if m, uri, err := p.parseRequestLine(); err != nil {
		return Request{}, err
	} else {
		r := NewRequest(m, p.rawBody, uri, NewHeaders(p))
		return r, nil
	}
}

var ErrIsNotResponse = errors.New("is not response")
var ErrWrongResponseCode = errors.New("wrong response code")

func (p *Parser) ParseResponse() (Response, error) {
	if p.IsRequest() {
		return Response{}, ErrIsNotResponse
	}
	pfl := strings.Split(p.getFirstLine(), " ")
	if c, err := strconv.Atoi(pfl[1]); err != nil {
		return Response{}, ErrWrongResponseCode
	} else if _, ok := ResponseCodes[c]; !ok {
		return Response{}, ErrWrongResponseCode
	} else {
		r := NewResponse(ResponseCode(c), p.rawBody, NewHeaders(p))
		return r, nil
	}
}

func (p *Parser) Parse() (Message, error) {
	if p.IsRequest() {
		return p.ParseRequest()
	} else {
		return p.ParseResponse()
	}
}

func NewParser(b string) *Parser {
	return &Parser{
		rawBody:  b,
		rawLines: strings.Split(strings.ReplaceAll(b, "\r\n", "\n"), "\n"),
	}
}

func decodeFirseLine() {

}

func Decode(b string) {
	b = strings.ReplaceAll(b, "\r\n", "\n")
	parts := strings.Split(b, "\n\n")
	lines := strings.Split(parts[0], "\n")

	// lines := strings.Split(strings.ReplaceAll(b, "\r\n", "\n"), "\n")
}
