package logger

import (
	"io"
	"os"
)

type Logger struct {
	Output io.WriteCloser
}

func New(output io.WriteCloser) *Logger {
	return &Logger{
		Output: output,
	}
}

func NewFile(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return &Logger{
		Output: f,
	}, nil
}

func (l *Logger) Log() error {
	return nil
}

func (l *Logger) Close() error {
	return l.Output.Close()
}
