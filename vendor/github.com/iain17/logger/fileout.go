package logger

import (
	"os"
)

type Fileout struct {
	file *os.File
	MinLevel int
}

func NewFileOut(path string, minLevel int) (*Fileout, error) {
	f, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &Fileout{
		file: f,
		MinLevel: minLevel,
	}, nil
}

func (s *Fileout) Close() {
	s.file.Close()
}

func (s *Fileout) Print(level int, message string) error {
	//Check if we want to log this
	if level < s.MinLevel {
		return nil
	}

	data := prefixes[level] + message + "\n"
	_, err := s.file.WriteString(data)
	return err
}

func (s *Fileout) Write(p []byte) (n int, err error) {
	return s.file.Write(p)
}