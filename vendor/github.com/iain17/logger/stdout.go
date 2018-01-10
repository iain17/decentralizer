package logger

import (
	"log"
	"os"
	"github.com/aybabtme/color/brush"
)

var stdout_std_logger *log.Logger
var stderr_std_logger *log.Logger

var prefixes = map[int]string{
	DEBUG:   "[debug] ",
	INFO:    "[info] ",
	WARNING: "[warning] ",
	ERROR:   "[error] ",
	FATAL:   "[fatal] ",
}

var coloredPrefixes = map[int]string{
	DEBUG:   brush.Cyan("[debug] ").String(),
	INFO:    brush.Green("[info] ").String(),
	WARNING: brush.Yellow("[warning] ").String(),
	ERROR:   brush.Red("[error] ").String(),
	FATAL:   brush.DarkRed("[fatal] ").String(),
}

type Stdout struct {
	MinLevel int
	Colored  bool
}

func init() {
	stdout_std_logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	stderr_std_logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func (s Stdout) std_format_message(level int, message string) (int, string) {
	if s.Colored {
		return 4, coloredPrefixes[level] + message
	}

	return 3, prefixes[level] + message
}

func (s Stdout) Print(level int, message string) error {
	//Check if we want to log this
	if level < s.MinLevel {
		return nil
	}

	d, msg := s.std_format_message(level, message)
	return stdout_std_logger.Output(d, msg)

	return nil
}
