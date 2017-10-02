package pb

import (
	"io"
	"errors"
	"github.com/gogo/protobuf/proto"
	"fmt"
	"io/ioutil"
	"github.com/iain17/decentralizer/discovery/env"
)

func Decode(r io.Reader) (*Message, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var result Message
	if err := proto.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	if result.Version != env.VERSION {
		return nil, errors.New(fmt.Sprintf("Invalid version"))
	}
	return &result, nil
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

func DecodePeerInfo(r io.Reader) (*PeerInfo, error) {
	message, err := Decode(r)
	if err != nil {
		return nil, err
	}
	result := message.GetPeerInfo()
	if result == nil {
		return nil, errors.New(fmt.Sprintf("Did not receive a PeerInfo message"))
	}
	return result, nil
}