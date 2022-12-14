package sdp_test

import (
	"signal/sdp"
	"testing"
)

var SDP = `v=0
o=jdoe 2890844526 2890842807 IN IP4 10.47.16.5
s=SDP Seminar
i=A Seminar on the session description protocol
u=http://www.example.com/seminars/sdp.pdf
e=j.doe@example.com (Jane Doe)
c=IN IP4 224.2.17.12/127
t=2873397496 2873404696
a=recvonly
m=audio 49170 RTP/AVP 0
m=video 51372 RTP/AVP 99
a=rtpmap:99 h263-1998/90000`

func TestDecodeResponse(t *testing.T) {
	if s, err := sdp.DecodeSDP(SDP); err != nil {
		t.Error(err)
	} else {
		t.Log(s.Origin)
	}
}