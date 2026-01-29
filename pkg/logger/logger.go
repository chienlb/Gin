package logger

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	warnLogger  *log.Logger
	debugLogger *log.Logger
}

var instance *Logger

func Init() *Logger {
	if instance == nil {
		instance = &Logger{
			infoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
			errorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
			warnLogger:  log.New(os.Stdout, "[WARN] ", log.LstdFlags),
			debugLogger: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
		}
	}
	return instance
}

func Get() *Logger {
	if instance == nil {
		return Init()
	}
	return instance
}

func (l *Logger) Info(message string) {
	l.infoLogger.Println(message)
}

func (l *Logger) Error(message string, err error) {
	if err != nil {
		l.errorLogger.Printf("%s: %v\n", message, err)
	} else {
		l.errorLogger.Println(message)
	}
}

func (l *Logger) Warn(message string) {
	l.warnLogger.Println(message)
}

func (l *Logger) Debug(message string) {
	l.debugLogger.Println(message)
}
