package media

import (
	"encoding/binary"
)

const RTP_VERSION uint8 = 2

//  0                   1                   2                   3
//  0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |V=2|P|X|  CC   |M|     PT      |       sequence number         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                           timestamp                           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |           synchronization source (SSRC) identifier            |
// +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
// |            contributing source (CSRC) identifiers             |
// |                             ....                              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

//  0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |      defined by profile       |           length              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                        header extension                       |
// |                             ....                              |

type RTPPacket struct {
	Version          uint8
	Padding          uint8
	Extension        uint8
	CSRCCount        uint8
	Marker           uint8
	PayloadType      uint8
	SequenceNumber   uint16
	Timestamp        uint32
	SSRC             uint32
	CSRCs            []uint32
	ExtensionProfile uint16
	ExtensionLength  uint16
}

func (rtpp *RTPPacket) Encode() []byte {
	packet := make([]byte, 16)
	var headers uint64 = 0

	// Set version to 2, 1-2 bits
	headers |= uint64(RTP_VERSION) << (64 - 2)

	// Set padding, 3 bit
	headers |= uint64(rtpp.Padding) << (64 - 3)

	// Set extension, 4 bit
	headers |= uint64(rtpp.Extension) << (64 - 4)

	// Set CSRC count, 5-8 bits
	headers |= uint64(rtpp.CSRCCount) << (64 - 8)

	// Set marker, 9 bit
	headers |= uint64(rtpp.Marker) << (64 - 9)

	// Set payload type, 10-16 bits
	headers |= uint64(rtpp.PayloadType) << (64 - 16)

	// Set sequence number, 17-32 bits
	headers |= uint64(rtpp.SequenceNumber) << (64 - 32)

	// Set timestamp, 33-64 bits
	headers |= uint64(rtpp.Timestamp) << (64 - 64)

	binary.BigEndian.PutUint64(packet, headers)
	binary.BigEndian.PutUint32(packet[8:12], rtpp.SSRC)

	n := 12
	for i := 0; i < int(rtpp.CSRCCount); i++ {
		binary.BigEndian.PutUint32(packet[n:n+4], rtpp.CSRCs[i])
		n += 4
	}

	if rtpp.Extension == 1 {
		binary.BigEndian.PutUint16(packet[n:n+2], rtpp.ExtensionProfile)
		binary.BigEndian.PutUint16(packet[n+2:n+4], rtpp.ExtensionLength)
	}

	return packet
}

func Decode(packet []byte) RTPPacket {
	rtpp := RTPPacket{}

	rtpp.Version = packet[0] >> 6
	rtpp.Padding = (packet[0] >> 5) & 0x1
	rtpp.Extension = (packet[0] >> 4) & 0x1
	rtpp.CSRCCount = packet[0] & 0xF

	rtpp.Marker = (packet[1] >> 7) & 0x1
	rtpp.PayloadType = packet[1] & 0x7F

	rtpp.SequenceNumber = binary.BigEndian.Uint16(packet[2:4])

	rtpp.Timestamp = binary.BigEndian.Uint32(packet[4:8])

	rtpp.SSRC = binary.BigEndian.Uint32(packet[8:12])
	rtpp.CSRCs = make([]uint32, 0)
	n := 12
	for i := 0; i < int(rtpp.CSRCCount); i++ {
		CSRC := binary.BigEndian.Uint32(packet[n : n+4])
		rtpp.CSRCs = append(rtpp.CSRCs, CSRC)
		n += 4
	}

	if rtpp.Extension == 1 {
		rtpp.ExtensionProfile = binary.BigEndian.Uint16(packet[n : n+2])
		rtpp.ExtensionLength = binary.BigEndian.Uint16(packet[n+2 : n+4])
	}

	return rtpp
}
