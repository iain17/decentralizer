package pb

import (
	"bufio"
	"fmt"
	"io"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/iain17/logger"
	"bytes"
	"time"
)

var delimiter = []byte{'\n', '\r', '\n'}
const VERSION = 1
const MAXSIZEMESSAGESIZE = 32768

func debug(data []byte) {
	logger.Info("writing")
	for _, d := range data {
		logger.Infof("%d", d)
	}
	logger.Info("end")
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-2]
	}
	return data
}

func ScanDelimiter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	//If we are at EOF or there is no further data to be processed. return it.
	if atEOF && len(data) == 0 {
		time.Sleep(1 * time.Second)
		return 0, nil, nil
	}
	if i := bytes.Index(data, delimiter); i >= 0 {
		// We have a full newline-terminated line.
		return i + 3, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
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
	s, err = w.Write(delimiter)
	//debug([]byte{delimiter})
	if err != nil {
		return err
	}
	if s != len(delimiter) {
		return errors.New("Didn't write the delimiter")
	}
	return err
}

func Decode(r io.Reader, resultChan chan *RPCMessage) (error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(ScanDelimiter)
	var data []byte
	scanner.Buffer(data, MAXSIZEMESSAGESIZE)
	for scanner.Scan() {
		data := scanner.Bytes()
		var result RPCMessage
		if err := proto.Unmarshal(data, &result); err != nil {
			return err
		}
		if result.Version != VERSION {
			return errors.New(fmt.Sprintf("Invalid version. Received %d, expected %d", result.Version, VERSION))
		}
		resultChan <- &result
	}
	close(resultChan)
	return nil
}

func Write(w io.Writer, msg *RPCMessage) error {
	msg.Version = VERSION
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return write(w, data)
}