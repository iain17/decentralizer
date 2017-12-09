package app

import (
	"io"
	"github.com/getlantern/framed"
	"errors"
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