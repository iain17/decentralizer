package pb

import (
	"io"
	"errors"
	"fmt"
	"github.com/iain17/decentralizer/discovery/env"
	"bufio"
	"github.com/golang/protobuf/proto"
)

var delimiter = byte(255)

func Write(w io.Writer, data []byte) error {
	s, err := w.Write(data)
	if err != nil {
		return err
	}
	if len(data) != s {
		return errors.New("Didn't write all of the data")
	}
	s, err = w.Write([]byte{delimiter})
	if err != nil {
		return err
	}
	if s != 1 {
		return errors.New("Didn't write the delimiter")
	}
	return err
}

func Decode(r io.Reader) (*Message, error) {
	data, err := bufio.NewReader(r).ReadBytes(delimiter)
	if err != nil {
		return nil, err
	}
	var result Message
	if err := proto.Unmarshal(data[:len(data)-1], &result); err != nil {
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

func DecodePeerInfo(r io.Reader, network string) (*PeerInfo, error) {
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