package sdp

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// v=0
// o=- 1815849 0 IN IP4 194.167.15.181
// s=Cisco SDP 0
// c=IN IP4 194.167.15.181
// t=0 0
// m=audio 20062 RTP/AVP 99 18 101 100
// a=rtpmap:99 G.729b/8000
// a=rtpmap:101 telephone-event/8000
// a=fmtp:101 0-15
// a=rtpmap:100 X-NSE/8000
// a=fmtp:100 200-202

var ErrSDPFieldNotExists = errors.New("SDP field not exists")

func NewSDPParseError(f, v string) error {
	msg := fmt.Sprintf("SDP field \"%s\" parse error. field value is: %s", f, v)
	return errors.New(msg)
}

type Nettype string

const (
	IN Nettype = "IN"
)

type Addrtype string

const (
	IP4 Addrtype = "IP4"
	IP6 Addrtype = "IP6"
)

type Origin struct {
	Username       string
	SessID         string
	SessVersion    string
	Nettype        Nettype
	Addrtype       Addrtype
	UnicastAddress net.IP
}

func (o *Origin) String() string {
	// o=<username> <sess-id> <sess-version> <nettype> <addrtype> <unicast-address>
	return fmt.Sprintf("%s %s %s %s %s %s", o.Username, o.SessID, o.SessVersion, string(o.Nettype), string(o.Addrtype), o.UnicastAddress.String())
}

type ConnectionAddress struct {
	IP                net.IP
	TTL               *int
	NumberOfAddresses *int
}

func (ca *ConnectionAddress) String() string {
	// <base multicast address>[/<ttl>]/<number of addresses>
	var builder strings.Builder
	builder.WriteString(ca.IP.String())
	if ca.TTL != nil {
		builder.WriteString(fmt.Sprintf("/%d", *ca.TTL))
	}
	if ca.NumberOfAddresses != nil {
		builder.WriteString(fmt.Sprintf("/%d", *ca.NumberOfAddresses))
	}
	return builder.String()
}

type ConnectionData struct {
	Nettype           Nettype
	Addrtype          Addrtype
	ConnectionAddress ConnectionAddress
}

func (c *ConnectionData) String() string {
	// c=<nettype> <addrtype> <connection-address>
	return fmt.Sprintf("%s %s %s", string(c.Nettype), string(c.Addrtype))
}

type Bandwidth struct {
	Bwtype    string
	Bandwidth int
}

type Timing struct {
	StartTime int
	StopTime  int
}

func (t *Timing) String() string {
	// t=<start-time> <stop-time>
	return fmt.Sprintf("%d %d", t.StartTime, t.StopTime)
}

type RepeatTimes struct {
	RepeatInterval       int
	ActiveDuration       int
	OffsetsFromStartTime int
}

func (rt *RepeatTimes) String() string {
	// r=<repeat interval> <active duration> <offsets from start-time>
	return fmt.Sprintf("%d %d %d", rt.RepeatInterval, rt.ActiveDuration, rt.OffsetsFromStartTime)
}

type TimeZone struct {
	AdjustmentTime int
	Offset         string
}

type EncryptionKey struct {
	Method        string
	EncryptionKey string
}

type Attribute struct {
	Name  string
	Value string
}

func (a *Attribute) String() string {
	if a.Value != "" {
		return a.Name
	} else {
		return fmt.Sprintf("%s:%s", a.Name, a.Value)
	}
}

// a=rtpmap:<payload type> <encoding name>/<clock rate> [/<encoding parameters>]
type RTPmap struct {
	EncodingName string // <encoding name>
	ClockRate    int    // <clock rate>
	Parameters   []int  // <encoding parameters>
}

type MediaDescriptionMode string

const (
	Recvonly MediaDescriptionMode = "recvonly" // a=recvonly
	Sendrecv MediaDescriptionMode = "sendrecv" // a=sendrecv
	Sendonly MediaDescriptionMode = "sendonly" // a=sendonly
	Inactive MediaDescriptionMode = "inactive" // a=inactive
)

type MediaDescriptionOrient string

const (
	Portrait  MediaDescriptionOrient = "portrait"
	Landscape MediaDescriptionOrient = "landscape"
	Seascape  MediaDescriptionOrient = "seascape"
)

type MediaDescriptionType string

const (
	Broadcast MediaDescriptionType = "broadcast"
	Meeting   MediaDescriptionType = "meeting"
	Moderated MediaDescriptionType = "moderated"
	Test      MediaDescriptionType = "test"
	H332      MediaDescriptionType = "H332"
)

type MediaDescriptionFmtp struct {
	Format     int    // <format>
	Parameters string // <format specific parameters>
}

// m=<media> <port> <proto> <fmt> ...
// m=<media> <port>/<number of ports> <proto> <fmt> ...
// m=video 49170/2 RTP/AVP 31
type MediaDescription struct {
	Media         string
	Port          int
	NumberOfPorts int
	Proto         string
	Fmt           int

	Ptime     *int                    // a=ptime:<packet time>
	Maxptime  *int                    // a=maxptime:<maximum packet time>
	RTPmaps   map[int]RTPmap          // a=rtpmap:<payload type> <encoding name>/<clock rate> [/<encoding parameters>]
	Sendrecv  *bool                   // a=sendrecv
	Sendonly  *bool                   // a=sendonly
	Inactive  *bool                   // a=inactive
	Orient    *MediaDescriptionOrient // a=orient:<orientation>
	SDPlang   *string                 // a=sdplang:<language tag>
	Lang      *string                 // a=lang:<language tag>
	Framerate *float64                // a=framerate:<frame rate>
	Quality   *int                    // a=quality:<quality>
	Fmtp      map[int]string          // a=fmtp:<format> <format specific parameters>
}

func (md *MediaDescription) String() string {
	if md.NumberOfPorts != 0 {
		// m=<media> <port>/<number of ports> <proto> <fmt>
		return fmt.Sprintf("%s %d/%d %s %d", md.Media, md.Port, md.NumberOfPorts, md.Proto, md.Fmt)
	} else {
		// m=<media> <port> <proto> <fmt>
		return fmt.Sprintf("%s %d %s %d", md.Media, md.Port, md.Proto, md.Fmt)
	}
}

func (md *MediaDescription) PtimeString() string {
	return fmt.Sprintf("ptime:%d", md.Ptime)
}

// func (md *MediaDescription) MaxptimeString() string {}

// func (md *MediaDescription) RTPmapsString() string {}

// func (md *MediaDescription) SendrecvString() string {}

// func (md *MediaDescription) SendonlyString() string {}

// func (md *MediaDescription) InactiveString() string {}

// func (md *MediaDescription) OrientString() string {}

// func (md *MediaDescription) SDPlangString() string {}

// func (md *MediaDescription) LangString() string {}

// func (md *MediaDescription) FramerateString() string {}

// func (md *MediaDescription) QualityString() string {}

// func (md *MediaDescription) FmtpString() string {}

type SDP struct {
	// Session description
	Version            int             //v=
	Origin             *Origin         //o=
	SessionName        *string         //s=
	SessionInformation *string         //i=
	URI                *string         //u=
	EmailAddress       *string         //e=
	PhoneNumber        *string         //p=
	ConnectionData     *ConnectionData //c=
	Bandwidth          *Bandwidth      //b=
	TimeZones          []TimeZone      //z=
	EncryptionKeys     []EncryptionKey //k=
	Attributes         []Attribute     //a=

	// Time description
	Timing      *Timing      //t=
	RepeatTimes *RepeatTimes //r=

	// Media description
	MediaDescriptions []MediaDescription // m=
}

func (sdp *SDP) Encode() string {
	var builder strings.Builder
	builder.WriteString("v=0\n")

	// o=<username> <sess-id> <sess-version> <nettype> <addrtype> <unicast-address>
	if sdp.Origin != nil {
		o := fmt.Sprintf("o=%s\n", sdp.Origin.String())
		builder.WriteString(o)
	}

	// s=<session name>
	if sdp.SessionName != nil {
		s := fmt.Sprintf("s=%s\n", *sdp.SessionName)
		builder.WriteString(s)
	}

	// i=<session description>
	if sdp.SessionInformation != nil {
		i := fmt.Sprintf("i=%s\n", *sdp.SessionInformation)
		builder.WriteString(i)
	}

	// u=<uri>
	if sdp.URI != nil {
		u := fmt.Sprintf("u=%s\n", *sdp.URI)
		builder.WriteString(u)
	}

	// e=<email-address>
	if sdp.EmailAddress != nil {
		e := fmt.Sprintf("e=%s\n", *sdp.EmailAddress)
		builder.WriteString(e)
	}

	// p=<phone-number>
	if sdp.PhoneNumber != nil {
		p := fmt.Sprintf("p=%s\n", *sdp.PhoneNumber)
		builder.WriteString(p)
	}

	// c=<nettype> <addrtype> <connection-address>
	if sdp.ConnectionData != nil {
		c := fmt.Sprintf("c=%s\n", sdp.ConnectionData.String())
		builder.WriteString(c)
	}

	// b=<bwtype>:<bandwidth>
	if sdp.Bandwidth != nil {
		b := fmt.Sprintf("b=%d:%d\n", sdp.Bandwidth.Bwtype, sdp.Bandwidth.Bandwidth)
		builder.WriteString(b)
	}

	// t=<start-time> <stop-time>
	if sdp.Timing != nil {
		t := fmt.Sprintf("t=%s\n", sdp.Timing.String())
		builder.WriteString(t)
	}

	// r=<repeat interval> <active duration> <offsets from start-time>
	if sdp.RepeatTimes != nil {
		r := fmt.Sprintf("r=%s\n", sdp.RepeatTimes.String())
		builder.WriteString(r)
	}

	// z=<adjustment time> <offset> <adjustment time> <offset> ....
	// for _, tz := range sdp.TimeZones {

	// }

	// k=<method>
	// k=<method>:<encryption key>
	// for _, ek := range sdp.EncryptionKeys {

	// }

	// a=<attribute>
	// a=<attribute>:<value>
	for _, item := range sdp.Attributes {
		a := fmt.Sprintf("r=%s\n", item.String())
		builder.WriteString(a)
	}

	// m=<media> <port> <proto> <fmt>
	for _, item := range sdp.MediaDescriptions {
		m := fmt.Sprintf("m=%s\n", item.String())
		builder.WriteString(m)
	}

	return builder.String()
}

func DecodeSDP(raw string) (*SDP, error) {
	sdp := &SDP{
		TimeZones:         make([]TimeZone, 0),
		EncryptionKeys:    make([]EncryptionKey, 0),
		Attributes:        make([]Attribute, 0),
		MediaDescriptions: make([]MediaDescription, 0),
	}
	lines := strings.Split(raw, "\n")
	currentMediaDescription := -1
	for _, line := range lines {
		parts := strings.Split(line, "=")
		key := parts[0]
		value := parts[0]
		if value == "" {
			return nil, NewSDPParseError(key, "EMPTY")
		} else {
			switch key {
			case "v":
				if value != "0" {
					return nil, NewSDPParseError(key, value)
				} else {
					sdp.Version = 0
				}
			case "o":
				// o=<username> <sess-id> <sess-version> <nettype> <addrtype> <unicast-address>
				if valueParts := strings.Split(value, " "); len(valueParts) != 6 {
					return nil, NewSDPParseError(key, value)
				} else {
					*sdp.Origin = Origin{
						Username:       valueParts[0],
						SessID:         valueParts[1],
						SessVersion:    valueParts[2],
						Nettype:        Nettype(valueParts[3]),
						Addrtype:       Addrtype(valueParts[4]),
						UnicastAddress: net.ParseIP(valueParts[5]),
					}
				}
			case "s":
				// s=<session name>
				*sdp.SessionName = value
			case "i":
				// i=<session description>
				*sdp.SessionInformation = value
			case "u":
				// u=<uri>
				*sdp.URI = value
			case "e":
				// e=<email-address>
				*sdp.EmailAddress = value
			case "p":
				// p=<phone-number>
				*sdp.PhoneNumber = value
			case "c":
				// c=<nettype> <addrtype> <connection-address>
				// c=IN IP4 224.2.1.1/127/3
				if valueParts := strings.Split(value, " "); len(valueParts) < 3 {
					return nil, NewSDPParseError(key, value)
				} else {
					connAddrParts := strings.Split(valueParts[2], "/")
					condAddr := ConnectionAddress{
						IP: net.ParseIP(connAddrParts[0]),
					}

					if len(connAddrParts) > 1 {
						if ttl, err := strconv.Atoi(connAddrParts[1]); err != nil {
							condAddr.TTL = &ttl
						} else if len(connAddrParts) > 2 {
							if noa, err := strconv.Atoi(connAddrParts[3]); err != nil {
								condAddr.NumberOfAddresses = &noa
							}
						}
					}

					*sdp.ConnectionData = ConnectionData{
						Nettype:           Nettype(valueParts[0]),
						Addrtype:          Addrtype(valueParts[1]),
						ConnectionAddress: condAddr,
					}
				}
			case "b":
				// b=<bwtype>:<bandwidth>
				// b=X-YZ:128
				if valueParts := strings.Split(value, ":"); len(valueParts) != 2 {
					return nil, NewSDPParseError(key, value)
				} else if bandwidth, err := strconv.Atoi(valueParts[1]); err != nil {
					return nil, NewSDPParseError(key, value)
				} else {
					*sdp.Bandwidth = Bandwidth{
						Bwtype:    valueParts[0],
						Bandwidth: bandwidth,
					}
				}
			case "z":
				// z=<adjustment time> <offset> <adjustment time> <offset> ....
			case "k":
				// k=<method>
				// k=<method>:<encryption key>
				// k=clear:<encryption key>
				// k=base64:<encoded encryption key>
				// k=uri:<URI to obtain key>
				// k=prompt
			case "t":
				// t=<start-time> <stop-time>
				if valueParts := strings.Split(value, " "); len(valueParts) != 2 {
					return nil, NewSDPParseError(key, value)
				} else if start, err := strconv.Atoi(valueParts[0]); err != nil {
					return nil, NewSDPParseError(key, value)
				} else if stop, err := strconv.Atoi(valueParts[1]); err != nil {
					return nil, NewSDPParseError(key, value)
				} else {
					*sdp.Timing = Timing{
						StartTime: start,
						StopTime:  stop,
					}
				}
			case "r":
				// r=<repeat interval> <active duration> <offsets from start-time>
				// r=604800 3600 0 90000
				// sdp.RepeatTimes = RepeatTimes{
				// 	RepeatInterval       int
				// 	ActiveDuration       int
				// 	OffsetsFromStartTime int
				// }
			case "m":
				// m=<media> <port> <proto> <fmt> ...
				// m=<media> <port>/<number of ports> <proto> <fmt> ...
				// m=video 49170/2 RTP/AVP 31
				if valueParts := strings.Split(value, " "); len(valueParts) != 4 {
					return nil, NewSDPParseError(key, value)
				} else {
					mediaDescription := MediaDescription{
						Media:   valueParts[0],
						Proto:   valueParts[2],
						RTPmaps: make(map[int]RTPmap),
						Fmtp:    make(map[int]string),
					}

					if fmt, err := strconv.Atoi(valueParts[3]); err != nil {
						return nil, NewSDPParseError(key, value)
					} else {
						mediaDescription.Fmt = fmt
					}

					if portParts := strings.Split(valueParts[1], "/"); len(portParts) > 1 {
						if port, err := strconv.Atoi(portParts[0]); err != nil {
							return nil, NewSDPParseError(key, value)
						} else if numberOfPorts, err := strconv.Atoi(portParts[1]); err != nil {
							return nil, NewSDPParseError(key, value)
						} else {
							mediaDescription.Port = port
							mediaDescription.NumberOfPorts = numberOfPorts
						}
					} else if port, err := strconv.Atoi(valueParts[1]); err != nil {
						return nil, NewSDPParseError(key, value)
					} else {
						mediaDescription.Port = port
					}

					sdp.MediaDescriptions = append(sdp.MediaDescriptions, mediaDescription)
					currentMediaDescription = len(sdp.MediaDescriptions) - 1
				}
			case "a":
				// a=<attribute>
				// a=<attribute>:<value>
				valueParts := strings.Split(value, ":")
				attribute := valueParts[0]
				var valueAttribute string
				if len(valueParts) == 2 {
					valueAttribute = valueParts[1]
				}
				switch attribute {
				case "cat", "keywds", "tool", "type", "charset":
					// a=cat:<category>
					// a=keywds:<keywords>
					// a=tool:<name and version of tool>
					// a=type:<conference type>
					// a=charset:<character set>
					sdp.Attributes = append(sdp.Attributes, Attribute{
						Name:  attribute,
						Value: valueAttribute,
					})
				case "sdplang":
					// a=sdplang:<language tag>
					if currentMediaDescription != -1 {

					} else {
						sdp.Attributes = append(sdp.Attributes, Attribute{
							Name: attribute,
						})
					}
				case "lang":
					// a=lang:<language tag>
					if currentMediaDescription != -1 {

					} else {
						sdp.Attributes = append(sdp.Attributes, Attribute{
							Name: attribute,
						})
					}
				default:
					if currentMediaDescription != -1 {
						switch attribute {
						case "ptime":
							// a=ptime:<packet time>
							if ptimeValue, err := strconv.Atoi(valueAttribute); err != nil {
								return nil, NewSDPParseError(key, value)
							} else {
								*sdp.MediaDescriptions[currentMediaDescription].Ptime = ptimeValue
							}
						case "maxptime":
							// a=maxptime:<maximum packet time>
							if maxptimeValue, err := strconv.Atoi(valueAttribute); err != nil {
								return nil, NewSDPParseError(key, value)
							} else {
								*sdp.MediaDescriptions[currentMediaDescription].Maxptime = maxptimeValue
							}
						case "rtpmap":
							// a=rtpmap:<payload type> <encoding name>/<clock rate> [/<encoding parameters>]
							if attributeValueParts := strings.Split(valueAttribute, " "); len(attributeValueParts) != 2 {
								return nil, NewSDPParseError(key, value)
							} else if payloadType, err := strconv.Atoi(attributeValueParts[0]); err != nil {
								return nil, NewSDPParseError(key, value)
							} else if encoding := strings.Split(attributeValueParts[1], "/"); len(encoding) < 2 {
								return nil, NewSDPParseError(key, value)
							} else if clockRate, err := strconv.Atoi(encoding[1]); err != nil {
								return nil, NewSDPParseError(key, value)
							} else {
								parameters := make([]int, 0)
								if len(encoding) > 2 {
									for _, param := range encoding[2:] {
										if p, err := strconv.Atoi(param); err != nil {
											return nil, NewSDPParseError(key, value)
										} else {
											parameters = append(parameters, p)
										}
									}
								}
								sdp.MediaDescriptions[currentMediaDescription].RTPmaps[payloadType] = RTPmap{
									EncodingName: encoding[0],
									ClockRate:    clockRate,
									Parameters:   parameters,
								}
							}
						case "sendrecv":
							// a=sendrecv
							*sdp.MediaDescriptions[currentMediaDescription].Sendrecv = true
						case "sendonly":
							// a=sendonly
							*sdp.MediaDescriptions[currentMediaDescription].Sendonly = true
						case "inactive":
							// a=inactive
							*sdp.MediaDescriptions[currentMediaDescription].Inactive = true
						case "orient":
							// a=orient:<orientation>
							*sdp.MediaDescriptions[currentMediaDescription].Orient = MediaDescriptionOrient(valueAttribute)
						case "framerate":
							// a=framerate:<frame rate>
							if framerate, err := strconv.ParseFloat(valueAttribute, 64); err != nil {
								return nil, NewSDPParseError(key, value)
							} else {
								*sdp.MediaDescriptions[currentMediaDescription].Framerate = framerate
							}
						case "quality":
							// a=quality:<quality>
							if quality, err := strconv.Atoi(valueAttribute); err != nil {
								return nil, NewSDPParseError(key, value)
							} else {
								*sdp.MediaDescriptions[currentMediaDescription].Quality = quality
							}
						case "fmtp":
							// a=fmtp:<format> <format specific parameters>
							if attributeValueParts := strings.Split(valueAttribute, " "); len(attributeValueParts) != 2 {
								return nil, NewSDPParseError(key, value)
							} else if format, err := strconv.Atoi(attributeValueParts[0]); err != nil {
								return nil, NewSDPParseError(key, value)
							} else {
								sdp.MediaDescriptions[currentMediaDescription].Fmtp[format] = attributeValueParts[1]
							}
						}
					}
				}
			}
		}
	}
	return sdp, nil
}
