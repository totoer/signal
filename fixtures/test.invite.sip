INVITE sip:test@foo.bar.com;transport=UDP SIP/2.0
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