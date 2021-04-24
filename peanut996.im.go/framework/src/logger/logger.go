package logger

import (
	"fmt"
	"framework/cfgargs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	LEVEL_FATAL = iota
	LEVEL_ERROR
	LEVEL_WARN
	LEVEL_INFO
	LEVEL_DEBUG
)

const logTimeFMT string = "2006-01-02 15:04:05.999"
const logFullFMTConst string = "##%v##%-23v##%v##%v##%v"

var logFullFMT string = logFullFMTConst // TIME LEVEL CODE MSG
var logLevelFMT = []string{"FATAL", "ERROR", "WARN", "INFO", "DEBUG"}

var once sync.Once

var logger *Logger

type Logger struct {
	level      int
	logChan    chan string
	logWriters []LogWriter
	binName    string
	sync       bool
}

// InitLogger 初始化Logger单实例对象
func InitLogger(cfg *cfgargs.SrvConfig) {
	once.Do(func() {
		l := LEVEL_INFO
		switch strings.ToLower(cfg.Log.Level) {
		case "fatal":
			l = LEVEL_FATAL
		case "error":
			l = LEVEL_ERROR
		case "warn":
			l = LEVEL_WARN
		case "info":
			l = LEVEL_INFO
		case "debug":
			l = LEVEL_DEBUG
		}
		logger = &Logger{
			level:   l,
			logChan: make(chan string, 2000),
		}
		logger.logWriters = make([]LogWriter, 0, 10)
		filePath, _ := exec.LookPath(os.Args[0])
		logger.binName = filepath.Base(filePath)
		if len(logger.binName) == 0 || "." == logger.binName {
			logger.binName = filepath.Base(os.Args[0])
		}

		logger.sync = cfg.Log.Sync
		if !cfg.Log.Sync {
			go func() {
				logger.handloop()
			}()
		}

		if len(cfg.Log.Path) > 0 {
			loggerWriterFile := NewLogWriterFile()
			err := loggerWriterFile.Open(cfg.Log.Path)
			if nil != err {
				log.Printf("open log file err: %v", err)
			}
			logger.logWriters = append(logger.logWriters, loggerWriterFile)
		}
		if cfg.Log.Console {
			logger.logWriters = append(logger.logWriters, NewLogWriterConsole())
		}
	})
}

// DefLogger 获取默认的Logger对象
func DefLogger() *Logger {
	return logger
}

func (l *Logger) GetBinname() string {
	return l.binName
}
func GetLogLevel() int {
	return DefLogger().level
}

func GetLogLevelToString(l int) string {
	switch l {
	case LEVEL_FATAL:
		return "FATAL"
	case LEVEL_ERROR:
		return "Error"
	case LEVEL_WARN:
		return "Warn"
	case LEVEL_INFO:
		return "Info"
	case LEVEL_DEBUG:
		return "Debug"
	default:
		return "Info"
	}
}

// Debug 输出Debug级别日志
func Debug(format string, args ...interface{}) {
	if LEVEL_DEBUG > DefLogger().level {
		return
	}
	DefLogger().WriteLog(LEVEL_DEBUG, format, args...)
}

// Info 输出Info级别日志
func Info(format string, args ...interface{}) {
	if LEVEL_INFO > DefLogger().level {
		return
	}
	DefLogger().WriteLog(LEVEL_INFO, format, args...)
}

// Warn 输出Warn级别日志
func Warn(format string, args ...interface{}) {
	if LEVEL_WARN > DefLogger().level {
		return
	}
	DefLogger().WriteLog(LEVEL_WARN, format, args...)
}

// Error 输出Error级别日志
func Error(format string, args ...interface{}) {
	if LEVEL_ERROR > DefLogger().level {
		return
	}
	DefLogger().WriteLog(LEVEL_ERROR, format, args...)
}

// Fatal 输出Fatal级别日志
func Fatal(format string, args ...interface{}) {
	if LEVEL_FATAL > DefLogger().level {
		return
	}
	DefLogger().WriteLog(LEVEL_FATAL, format, args...)

	// 准备退出
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)
	logMsg := fmt.Sprintf(
		logFullFMT,
		logLevelFMT[LEVEL_FATAL],
		time.Now().Format(logTimeFMT), // 2019-01-23 14:41:26.10269679
		DefLogger().GetBinname(),
		fmt.Sprintf("%v:%v", filename, strconv.FormatInt(int64(line), 10)),
		msg)

	<-time.After(1 * time.Second)
	panic(logMsg)
}

func (l *Logger) handloop() {
	counter := 0
	for {
		select {
		case logMsg, ok := <-l.logChan:
			if !ok {
				continue
			}
			for _, wtr := range logger.logWriters {
				_, err := wtr.WriteString(logMsg)
				if err != nil {
					fmt.Println(err)
				}
			}
			counter++
		case <-time.After(time.Second * 2):
			if counter > 0 {
				counter = 0
				for _, wtr := range logger.logWriters {
					wtr.Flush()
				}
			}
		}
	}
}

func ChangeLogFMT(sep string, trim bool, enableAppName bool) {
	logFullFMT = strings.ReplaceAll(logFullFMTConst, "##", sep)
	if trim {
		logFullFMT = strings.Trim(logFullFMT, sep)
	}
	if !enableAppName {
		DefLogger().binName = ""
	}
}

func (l *Logger) WriteLog(logLevel int, msg string, v ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	if len(v) > 0 {
		msg = fmt.Sprintf(msg, v...)
	}

	_, filename := path.Split(file)
	logMsg := fmt.Sprintf(
		logFullFMT,
		logLevelFMT[logLevel],
		time.Now().Format(logTimeFMT), // 2019-01-23 14:41:26.10269679
		l.binName,
		fmt.Sprintf("%v:%v", filename, strconv.FormatInt(int64(line), 10)),
		msg)

	if !l.sync {
		fmt.Printf("write to channel, %v\n", logMsg)
		l.logChan <- logMsg
		return
	}

	for _, wtr := range logger.logWriters {
		_, err := wtr.WriteString(logMsg)
		if err != nil {
			panic(err)
		}
	}
}
