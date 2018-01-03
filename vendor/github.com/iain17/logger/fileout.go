package logger

import (
	"os"
)

type Fileout struct {
	file *os.File
}

func NewFileOut(path string) (*Fileout, error) {
	f, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &Fileout{
		file: f,
	}, nil
}

func (s *Fileout) Close() {
	s.file.Close()
}

func (s *Fileout) Print(level int, message string) error {
	data := prefixes[level] + message + "\n"
	_, err := s.file.WriteString(data)
	return err
}
