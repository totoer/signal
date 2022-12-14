package media

import (
	"bytes"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type MediaChanal struct {
	conn        *net.UDPConn
	inputBuffer bytes.Buffer
}

func (mc *MediaChanal) write() {
	// rtpp := RTPPacket{}
	// rtpp.flags |= (1 << pos)
	// enc := gob.NewEncoder(&mc.inputBuffer)
}

func (mc *MediaChanal) Start() {}

func (mc *MediaChanal) IsStarted() bool {
	return false
}

func (mc *MediaChanal) End() {}

func (mc *MediaChanal) Listen() {
	defer mc.conn.Close()

	for {
		buffer := make([]byte, 1024)
		if l, _, err := mc.conn.ReadFrom(buffer); err != nil {
			continue
		} else {
			rtpp := NewRTPPacket()
			rtpp.Decode(buffer[:l])
			fmt.Println("Padding:", rtpp.Padding)
			fmt.Println("Extension:", rtpp.Extension)
			fmt.Println("CSRCCount:", rtpp.CSRCCount)
			fmt.Println("Marker:", rtpp.Marker)
			fmt.Println("PayloadType:", rtpp.PayloadType)
			fmt.Println("SequenceNumber:", rtpp.SequenceNumber)
			fmt.Println("Timestamp:", rtpp.Timestamp)
			fmt.Println("SSRC:", rtpp.SSRC)
			fmt.Println("CSRCs:", rtpp.CSRCs)
			fmt.Println("ExtensionProfile:", rtpp.ExtensionProfile)
			fmt.Println("ExtensionLength:", rtpp.ExtensionLength)
		}
	}
}

func (mc *MediaChanal) Connect(am *MediaChanal) {}

func (mc *MediaChanal) Play(f ...string) {
	// if f, err := os.Open(fmt.Sprintf("audio/%s", f)); err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	defer f.Close()
	// 	decoder := wav.NewDecoder(f)
	// 	if decoder.IsValidFile() {
	// 		if buffer, err := decoder.FullPCMBuffer(); err != nil {

	// 		} else {
	// 			for _, lpcm := range buffer.Data {
	// 				chunk := g711.EncodeAlawFrame(int16(lpcm))
	// 			}
	// 		}
	// 	}
	// }
}

func (ms *MediaChanal) Beeps() {}

func (mc *MediaChanal) Stop() {}

func NewMediaChanal(mid uuid.UUID, cid string) (*MediaChanal, error) {
	host := viper.GetString("media.host")
	port := viper.GetInt("media.port")

	addr := &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: port,
	}

	if conn, err := net.ListenUDP("udp", addr); err != nil {
		return nil, err
	} else {
		return &MediaChanal{
			conn: conn,
		}, nil
	}
}

type MediaMixer struct{}

func (m MediaMixer) Join(mc *MediaChanal) {}

func (m *MediaMixer) Play(f string) {}

func (m *MediaMixer) Stop() {}

func NewMediaMixer() (*MediaMixer, error) {
	return &MediaMixer{}, nil
}
