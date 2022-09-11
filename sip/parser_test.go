package sip_test

import (
	"fmt"
	"signal/sip"
	"testing"
)

var SIP_REQUEST = `INVITE sip:test@foo.bar.com;transport=UDP SIP/2.0
Via: SIP/2.0/UDP 10.10.10.10:44444;branch=z9hG4bK-524287-1---3c38414a643cc244;rport
Max-Forwards: 70
Contact: <sip:user@bar.foo.com:44444;transport=UDP>
To: <sip:test@foo.bar.com>
From: <sip:user@foo.bar.com;transport=UDP>;tag=902cba13
Call-ID: gwQlUuwZxsFHSoh5XE8AOA
CSeq: 2 INVITE
Allow: INVITE, ACK, CANCEL, BYE, NOTIFY, REFER, MESSAGE, OPTIONS, INFO, SUBSCRIBE
Content-Type: application/sdp
User-Agent: Z 5.5.10 v2.10.17.3
Allow-Events: presence, kpml, talk
Content-Length: 333

v=0
o=Z 697308982 1 IN IP4 0.0.0.0
s=Z
c=IN IP4 0.0.0.0
t=0 0
m=audio 60417 RTP/AVP 106 9 98 101 0 8 3
a=rtpmap:106 opus/48000/2
a=fmtp:106 sprop-maxcapturerate=16000; minptime=20; useinbandfec=1
a=rtpmap:98 telephone-event/48000
a=fmtp:98 0-16
a=rtpmap:101 telephone-event/8000
a=fmtp:101 0-16
a=sendrecv
`

var SIP_RESPONSE = `SIP/2.0 100 Trying
Via: SIP/2.0/UDP 127.0.0.1:5080;branch=19048c29-9263-4d8f-b2b0-b53ee05331c2
To: \"foo\" <sip:foo@127.0.0.1:5080>
From: <sip:test@127.0.0.1:5080>;tag=d1712d60
Call-ID: 577cf2e4-fed5-48d3-a25b-b56a5a5f24d4
CSeq: 0 INVITE
Server: Twinkle/1.10.1
Content-Length: 0

`

func TestParseResponse(t *testing.T) {
	p := sip.NewParser(SIP_RESPONSE)
	if m, err := p.Parse(); err != nil {
		t.Error(err)
	} else {
		resp := m.(sip.Response)
		t.Logf("Response code: %s", fmt.Sprint(resp.Code))
		if cid, err := resp.GetHeaders().GetCallID(); err != nil {
			t.Error(err)
		} else {
			t.Log(cid)
		}
	}
}

func RequestTest(t testing.TB, d string) {
	p := sip.NewParser(d)
	if p.IsRequest() {
		if r, err := p.ParseRequest(); err != nil {
			t.Error(err)
		} else {
			// Request line
			if t.Logf("Message INVITE login is %s", r.URI.Login); r.URI.Login != "test" {
				t.Error("Login != test")
			} else if t.Logf("Message INVITE host is %s", r.URI.Host); r.URI.Host != "foo.bar.com" {
				t.Error("Host != foo.bar.com")
			}

			// Max-Forwards
			if maxForwarders, err := r.Headers.GetMaxForwards(); err != nil {
				t.Error(err)
			} else {
				if t.Logf("Max-Forwards is %d", maxForwarders.Value); maxForwarders.Value != 70 {
					t.Error("Max-Forwards != 70")
				}
			}

			// Via
			if vias, err := r.Headers.GetVias(); err != nil {
				t.Error(err)
			} else {
				via := vias[0]
				if t.Logf("Via host is %s", via.Host); via.Host != "10.10.10.10:44444" {
					t.Error("Via host != 10.10.10.10:44444")
				} else if t.Logf("Via branch is %s", via.Branch); via.Branch != "z9hG4bK-524287-1---3c38414a643cc244" {
					t.Error("Via branch != z9hG4bK-524287-1---3c38414a643cc244")
				}
			}

			// From
			// if f, err := r.Headers.GetFirst("From"); err != nil {
			// 	t.Error(err)
			// } else if from, ok := f.(sip.Destination); !ok {
			// 	t.Error("From cast error")
			// } else if t.Logf("From login is %s", from.Address.URI.Login); from.Address.URI.Login != "user" {
			// 	t.Error("From login != user")
			// } else if t.Logf("From host is %s", from.Address.URI.Host); from.Address.URI.Host != "foo.bar.com" {
			// 	t.Error("From host != foo.bar.com")
			// } else if t.Logf("From tag is %s", from.Tag); from.Tag != "902cba13" {
			// 	t.Error("From tag != 902cba13")
			// }

			// Call-Id
			// if f, err := r.Headers.GetFirst("Call-ID"); err != nil {
			// 	t.Error(err)
			// } else if cid, ok := f.(sip.PlainHeader); !ok {
			// 	t.Error("Call-ID cast error")
			// } else if t.Logf("Call-id is %s", cid.Value); cid.Value != "gwQlUuwZxsFHSoh5XE8AOA" {
			// 	t.Error("Call-ID != gwQlUuwZxsFHSoh5XE8AOA")
			// }

			// CSeq
			// if f, err := r.Headers.GetFirst("CSeq"); err != nil {
			// 	t.Error(err)
			// } else if cseq, ok := f.(sip.CSeq); !ok {
			// 	t.Error("CSeq cast error")
			// } else if t.Logf("CSeq value is %d", cseq.Value); cseq.Value != 2 {
			// 	t.Errorf("CSeq value != 2")
			// }

			// if fields, err := r.Headers.Get("Allow"); err != nil {
			// 	t.Error(err)
			// } else {
			// 	for _, f := range fields {
			// 		if a, ok := f.(sip.Allow); !ok {
			// 			t.Error("Allow cast error")
			// 		} else {
			// 			t.Logf("Allow: %s", a.String())
			// 		}
			// 	}
			// }
		}
	}
}

// go clean -testcache && /usr/local/go/bin/go test -parallel 1 -timeout 30s -run ^TestPerformanceParser$ signal/sip
// go clean -testcache && go test -bench=. -parallel 1 -timeout 30s
func BenchmarkPerformanceParser(b *testing.B) {
	sip.PROTOCOL = "UDP"
	for i := 0; i < 100000; i++ {
		RequestTest(b, SIP_REQUEST)
	}
}

// go clean -testcache && go test -bench=. -parallel 1 -timeout 30s -run ^TestParser$ signal/sip
func TestParser(t *testing.T) {
	sip.PROTOCOL = "UDP"
	RequestTest(t, SIP_REQUEST)
}
