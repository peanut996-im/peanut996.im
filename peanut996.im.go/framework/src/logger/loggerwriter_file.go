package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type LogWriterFile struct {
	file               *os.File
	logPath            string
	logPerm            string
	checkRotateCounter int
	timeStamp          time.Time
}

func NewLogWriterFile() *LogWriterFile {
	return &LogWriterFile{
		logPerm: "0666",
	}
}

func (l *LogWriterFile) Open(logPath string) error {
	dir, _ := path.Split(logPath)
	err := os.MkdirAll(dir, 0777)
	if nil != err {
		return err
	}

	l.logPath = logPath
	ret := l.doOpenFile(l.logPath, l.logPerm)
	if nil != ret {
		panic(ret)
	}
	l.checkRotate()
	return ret
}

func (l *LogWriterFile) doOpenFile(logPath string, logPerm string) error {
	logPath = strings.Replace(logPath, "{timestamp}", time.Now().Format("20060102_150405"), -1)

	// Open the log file
	perm, err := strconv.ParseInt(logPerm, 8, 64)
	if err != nil {
		return err
	}
	fd, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	if err == nil {
		// Make sure file perm is user set perm cause of `os.OpenFile` will obey umask
		err = os.Chmod(logPath, os.FileMode(perm))
		if err != nil {
			log.Fatal(err)
		}
	}
	l.file = fd
	l.timeStamp = time.Now()
	return err
}

func (l *LogWriterFile) Close() {
	l.file.Close()
}

func (l *LogWriterFile) Flush() {
	err := l.file.Sync()
	if err != nil {
		fmt.Println(err)
	}
}

func (l *LogWriterFile) WriteString(s string) (n int, err error) {
	l.checkRotateCounter += 1
	if l.checkRotateCounter > 1000 {
		l.checkRotate()
		l.checkRotateCounter = 0
	}
	return fmt.Fprintln(l.file, s)
}

func (l *LogWriterFile) checkRotate() {
	rotate := false
	if l.timeStamp.Day() != time.Now().Day() {
		rotate = true
	}

	fi, _ := l.file.Stat()
	if fi.Size() > 0x1F400000 {
		rotate = true
	}

	if !rotate {
		return
	}

	if -1 == strings.Index(l.logPath, "{timestamp}") {
		if -1 != strings.Index(l.logPath, ".log") {
			l.logPath = strings.Replace(l.logPath, ".log", "-{timestamp}.log", -1)
		} else {
			l.logPath = fmt.Sprintf("%v-{timestamp}.log", l.logPath)
		}
	}

	l.Flush()
	l.Close()
	err := l.doOpenFile(l.logPath, l.logPerm)
	if err != nil {
		log.Fatal(err)
	}
}
