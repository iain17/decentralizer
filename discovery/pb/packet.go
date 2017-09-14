package pb

import (
	"io"
	"encoding/gob"
	"errors"
)
//TODO: If this is slow, then look into writing the bytes manually.
type (
	Header struct {
		Version uint8
		//MessageNum uint64
	}
	Body struct {
		Type   MessageType
		Data []byte
	}

	Packet struct {
		Head Header
		Body Body
	}
)

const VERSION = 1

func NewPacket(messageType MessageType, msg []byte) *Packet{
	return &Packet{
		Head: Header{
			Version: VERSION,
			//MessageNum: messageNum,
		},
		Body: Body{
			Type: messageType,
			Data: msg,
		},
	}
}

func (p *Packet) Write(w io.Writer) (error) {
	enc := gob.NewEncoder(w)
	return enc.Encode(p)
}

func Decode(r io.Reader) (*Packet, error) {
	dec := gob.NewDecoder(r)
	var packet Packet
	err := dec.Decode(&packet)
	if packet.Head.Version != VERSION {
		return nil, errors.New("incorrect version")
	}
	return &packet, err
}