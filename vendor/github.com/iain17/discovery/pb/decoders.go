package pb

import (
	"errors"
	"fmt"
	"github.com/iain17/discovery/env"
	"github.com/golang/protobuf/proto"
	"io"
	"github.com/iain17/framed"
)

func Decode(r io.Reader) (*Message, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				err = fmt.Errorf("panic while decoding message: %v", err)
			}
			return
		}
	}()
	var data []byte
	data, err = framed.Read(r)
	if err != nil {
		return nil, err
	}
	var msg Message
	if err := proto.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	if msg.Version != env.VERSION {
		return nil, errors.New(fmt.Sprintf("Invalid version. Received %d, expected %d", msg.Version, env.VERSION))
	}
	return &msg, err
}

func DecodeHeartBeat(r io.Reader) error {
	message, err := Decode(r)
	if err != nil {
		return err
	}
	result := message.GetHeartbeat()
	if result == nil {
		return errors.New(fmt.Sprintf("Did not receive a HeartBeat message"))
	}
	return nil
}

func DecodePeerInfo(r io.Reader, network string) (*DPeerInfo, error) {
	message, err := Decode(r)
	if err != nil {
		return nil, err
	}
	result := message.GetPeerInfo()
	if result == nil {
		return nil, errors.New(fmt.Sprintf("Did not receive a PeerInfo message"))
	}
	//Check network
	if result.Network != network {
		return nil, errors.New(fmt.Sprintf("Peer not from the same network. Received %s got %s", result.Network, network))
	}
	return result, nil
}
