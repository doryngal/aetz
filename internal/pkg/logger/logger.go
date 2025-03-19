package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

type Logger struct {
	*log.Logger
}

func Setup(logFile string) (*Logger, error) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return nil, fmt.Errorf("open logfile - %v", err)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	l := &Logger{
		Logger: log.New(multiWriter, "", log.Ldate|log.Ltime|log.Lshortfile),
	}

	return l, nil
}

func (l *Logger) logWithCallerInfo(prefix string, v ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		l.SetPrefix(fmt.Sprintf("%s %s:%d ", prefix, file, line))
	} else {
		l.SetPrefix(prefix)
	}
	l.Logger.Println(v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.logWithCallerInfo("FATAL:", v...)
	os.Exit(1)
}

func (l *Logger) Info(v ...interface{}) {
	l.logWithCallerInfo("INFO:", v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.logWithCallerInfo("ERROR:", v...)
}
