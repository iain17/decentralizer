package pb

import (
	"errors"
	"fmt"
	"github.com/iain17/decentralizer/discovery/env"
	"github.com/golang/protobuf/proto"
	"io"
	"github.com/getlantern/framed"
)

//TODO: Replace framed for something better.
func Write(w io.Writer, data []byte) error {
	fw := framed.NewWriter(w)
	s, err := fw.Write(data)
	if err != nil {
		return err
	}
	if len(data) != s {
		return errors.New("Didn't write all of the data")
	}
	return err
}

func Read(r io.Reader) ([]byte, error) {
	fr := framed.NewReader(r)
	return fr.ReadFrame()
}

func Decode(r io.Reader) (*Message, error) {
	data, err := Read(r)
	if err != nil {
		return nil, err
	}
	var result Message
	if err := proto.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	if result.Version != env.VERSION {
		return nil, errors.New(fmt.Sprintf("Invalid version. Received %d, expected %d", result.Version, env.VERSION))
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
