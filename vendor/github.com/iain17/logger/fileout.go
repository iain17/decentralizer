package logger

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

type Fileout struct {
	l *lumberjack.Logger
	MinLevel int
}

func NewFileOut(path string, minLevel int) (*Fileout, error) {
	return &Fileout{
		l: &lumberjack.Logger{
			Filename:   path,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
			Compress:   true, // disabled by default
		},
		MinLevel: minLevel,
	}, nil
}

func (s *Fileout) Close() {
	s.l.Close()
}

func (s *Fileout) Print(level int, message string) error {
	//Check if we want to log this
	if level < s.MinLevel {
		return nil
	}

	data := prefixes[level] + message + "\n"
	_, err := s.l.Write([]byte(data))
	return err
}

func (s *Fileout) Write(p []byte) (n int, err error) {
	return s.l.Write(p)
}