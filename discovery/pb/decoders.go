package pb

import (
	"io"
	"errors"
	"github.com/gogo/protobuf/proto"
)

//This file is to remove code dup

//We don't actually call protobuf here because a heartbeat doesn't contain anything.
func DecodeHeartBeat(r io.Reader) error {
	packet, err := Decode(r)
	if err != nil {
		return err
	}
	if packet.Body.Type != HearBeatMessage {
		return errors.New("message type was incorrect")
	}
	return nil
}

func DecodePeerInfo(r io.Reader) (*PeerInfo, error) {
	packet, err := Decode(r)
	if err != nil {
		return nil, err
	}
	if packet.Body.Type != PeerInfoMessage {
		return nil, errors.New("message type was incorrect")
	}
	var result PeerInfo
	if err := proto.Unmarshal(packet.Body.Data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}