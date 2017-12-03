package pb

import (
	"io"
	"errors"
	"bufio"
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

//TODO: Read MAX size!!
func Read(r io.Reader) ([]byte, error) {
	data, err := bufio.NewReader(r).ReadBytes(delimiter)
	if err != nil {
		return nil, err
	}
	return data[:len(data)-1], err
}