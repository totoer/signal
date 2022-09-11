package sip

import "net"

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

type ConnectionAddress struct {
	IP                net.IP
	TTL               int
	NumberOfAddresses int
}

type ConnectionData struct {
	Nettype           Nettype
	Addrtype          Addrtype
	ConnectionAddress ConnectionAddress
}

type Bandwidth struct {
	Bwtype    string
	Bandwidth int
}

type Timing struct {
	StartTime int
	StopTime  int
}

type RepeatTimes struct {
	RepeatInterval       int
	ActiveDuration       int
	OffsetsFromStartTime int
}

type TimeZones struct {
	AdjustmentTime int
	Offset         string
}

type EncryptionKey struct {
	Method        string
	EncryptionKey string
}

type Attributes struct {
	Attribute string
	Value     string
}

type RTPmap struct {
	EncodingName string // <encoding name
	ClockRate    int    // <clock rate>
	Parameters   []int  //<encoding parameters>
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

type MediaDescription struct {
	Media         string
	Port          int
	NumberOfPorts int
	Proto         string
	Fmt           []int

	Cat       []string               // a=cat:<category>
	Keywds    []string               // a=keywds:<keywords>
	Tool      string                 // a=tool:<name and version of tool>
	Ptime     int                    // a=ptime:<packet time>
	Maxptime  int                    // a=maxptime:<maximum packet time>
	RTPmaps   map[int]RTPmap         // a=rtpmap:<payload type> <encoding name>/<clock rate> [/<encoding parameters>]
	Mode      MediaDescriptionMode   // a=recvonly/sendrecv/sendonly/inactive
	Orient    MediaDescriptionOrient // a=orient:<orientation>
	Type      MediaDescriptionType   // a=type:<conference type>
	Charset   string                 // a=charset:<character set>
	SDPlang   string                 // a=sdplang:<language tag>
	Lang      string                 // a=lang:<language tag>
	Framerate float32                // a=framerate:<frame rate>
	Quality   int                    // a=quality:<quality>
	Fmtp      map[int]string         // a=fmtp:<format> <format specific parameters>
}

type SDP struct {
	// Session description
	Version            int             //v=
	Origin             Origin          //o=
	SessionName        string          //s=
	SessionInformation string          //i=
	URI                string          //u=
	EmailAddress       string          //e=
	PhoneNumber        string          //p=
	ConnectionData     ConnectionData  //c=
	Bandwidth          Bandwidth       //b=
	TimeZones          []TimeZones     //z=
	EncryptionKeys     []EncryptionKey //k=
	Attributes         []Attributes    //a=

	// Time description
	Timing      Timing      //t=
	RepeatTimes RepeatTimes //r=

	// Media description
	MediaDescriptions []MediaDescription // m=
}

func (sdp *SDP) Data() []byte {
	return []byte("")
}

func ParseSDP() SDP {
	return SDP{}
}
