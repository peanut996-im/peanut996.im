package logger

import (
	"fmt"
)

type LoggerWriterConsole struct {
}

func NewLoggerWriterConsole() *LoggerWriterConsole {
	return &LoggerWriterConsole{}
}

func (l *LoggerWriterConsole) Open(logPath string) error {
	return nil
}

func (l *LoggerWriterConsole) Close() {
}

func (l *LoggerWriterConsole) Flush() {
}

func (l *LoggerWriterConsole) WriteString(s string) (n int, err error) {
	return fmt.Println(s)
}

func (l *LoggerWriterConsole) Write(b []byte) (n int, err error) {
	return l.WriteString(string(b))
}
