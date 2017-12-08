package pb

import (
	"bufio"
	"fmt"
	"io"
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/iain17/logger"
	"reflect"
)

var delimiter = byte(255)
const VERSION = 1

func debug(data []byte) {
	logger.Info("writing")
	for _, d := range data {
		logger.Infof("%d", d)
	}
	logger.Info("end")
}


func write(w io.Writer, data []byte) error {
	s, err := w.Write(data)
	//debug(data)
	if err != nil {
		return err
	}
	if len(data) != s {
		return errors.New("Didn't write all of the data")
	}
	s, err = w.Write([]byte{delimiter})
	//debug([]byte{delimiter})
	if err != nil {
		return err
	}
	if s != 1 {
		return errors.New("Didn't write the delimiter")
	}
	return err
}

func Decode(r io.Reader) (*RPCMessage, error) {
	data, err := bufio.NewReader(r).ReadBytes(delimiter)
	if err != nil {
		return nil, err
	}
	var result RPCMessage
	if err := proto.Unmarshal(data[:len(data)-1], &result); err != nil {
		return nil, err
	}
	if result.Version != VERSION {
		return nil, errors.New(fmt.Sprintf("Invalid version. Received %d, expected %d", result.Version, VERSION))
	}
	return &result, nil
}

func Write(w io.Writer, msg *RPCMessage) error {
	msg.Version = VERSION
	msgType := reflect.TypeOf(msg.GetMsg())
	msg.Type = int32(types[msgType])
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return write(w, data)
}