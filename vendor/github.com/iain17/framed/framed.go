package framed

import (
	"io"
	"errors"
	"encoding/binary"
	"bytes"
	"math"
)

//Simple framing.
func Write(w io.Writer, data []byte) error {
	buflen := uint32(len(data))
	if buflen > uint32(MAX_SIZE) || buflen > math.MaxUint32 {
		return errors.New("writes exceeded max size")
	}
	//Send length
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, buflen)
	if err != nil {
		return err
	}
	//Write data
	_, err = buffer.Write(data)
	s, err := buffer.WriteTo(w)
	if buflen + 4 != uint32(s) {
		return errors.New("incomplete write")
	}
	return err
}

func Read(r io.Reader) ([]byte, error) {
	headerBytes := make([]byte, 4)
	if _, err := io.ReadFull(r, headerBytes); err != nil {
		return nil, err
	}
	buflen := binary.LittleEndian.Uint32(headerBytes)
	if buflen > uint32(MAX_SIZE) || buflen > math.MaxUint32 {
		return nil, errors.New("read exceeded max size")
	}
	data := make([]byte, buflen)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}
	return data, nil
}