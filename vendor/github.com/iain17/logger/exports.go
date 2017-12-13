package logger

import (
	"fmt"
)

func Debug(v ...interface{}) {
	for _, output := range outputs {
		output.Print(DEBUG, fmt.Sprint(v...))
	}
}

func Debugf(format string, v ...interface{}) {
	for _, output := range outputs {
		output.Print(DEBUG, fmt.Sprintf(format, v...))
	}
}

func Info(v ...interface{}) {
	for _, output := range outputs {
		output.Print(INFO, fmt.Sprint(v...))
	}
}

func Infof(format string, v ...interface{}) {
	for _, output := range outputs {
		output.Print(INFO, fmt.Sprintf(format, v...))
	}
}

func Warning(v ...interface{}) {
	for _, output := range outputs {
		output.Print(WARNING, fmt.Sprint(v...))
	}
}

func Warningf(format string, v ...interface{}) {
	for _, output := range outputs {
		output.Print(WARNING, fmt.Sprintf(format, v...))
	}
}

func Error(v ...interface{}) {
	for _, output := range outputs {
		output.Print(ERROR, fmt.Sprint(v...))
	}
}

func Errorf(format string, v ...interface{}) {
	for _, output := range outputs {
		output.Print(ERROR, fmt.Sprintf(format, v...))
	}
}

func Fatal(v ...interface{}) {
	s := fmt.Sprint(v...)

	for _, output := range outputs {
		output.Print(FATAL, s)
	}

	panic(s)
}

func Fatalf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)

	for _, output := range outputs {
		output.Print(FATAL, s)
	}

	panic(s)
}
