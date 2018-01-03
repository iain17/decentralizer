package logger

import "fmt"

type Logger struct {
	module string
}

func New(module string) *Logger {
	return &Logger{
		module: module,
	}
}

func (l *Logger) Debug(v ...interface{}) {
	for _, output := range outputs {
		output.Print(DEBUG, fmt.Sprint(l.module, v))
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	for _, output := range outputs {
		output.Print(DEBUG, fmt.Sprintf(l.module+"	"+format, v...))
	}
}

func (l *Logger) Info(v ...interface{}) {
	for _, output := range outputs {
		output.Print(INFO, fmt.Sprint(l.module, v))
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	for _, output := range outputs {
		output.Print(INFO, fmt.Sprintf(l.module+"	"+format, v...))
	}
}

func (l *Logger) Warning(v ...interface{}) {
	for _, output := range outputs {
		output.Print(WARNING, fmt.Sprint(l.module, v))
	}
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	for _, output := range outputs {
		output.Print(WARNING, fmt.Sprintf(l.module+"	"+format, v...))
	}
}

func (l *Logger) Error(v ...interface{}) {
	for _, output := range outputs {
		output.Print(ERROR, fmt.Sprint(l.module, v))
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	for _, output := range outputs {
		output.Print(ERROR, fmt.Sprintf(l.module+"	"+format, v...))
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	s := fmt.Sprint(v...)

	for _, output := range outputs {
		output.Print(FATAL, s)
	}

	panic(s)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)

	for _, output := range outputs {
		output.Print(FATAL, s)
	}

	panic(s)
}
