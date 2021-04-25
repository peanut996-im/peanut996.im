package logger

import (
	"fmt"
)

type LogWriterConsole struct {
}

func NewLogWriterConsole() *LogWriterConsole {
	return &LogWriterConsole{}
}

func (l *LogWriterConsole) Open(logPath string) error {
	return nil
}

func (l *LogWriterConsole) Close() {
}

func (l *LogWriterConsole) Flush() {
}

func (l *LogWriterConsole) WriteString(s string) (n int, err error) {
	return fmt.Println(s)
}
