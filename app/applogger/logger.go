package applogger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

var Log Logger

func NewLogger() (*LoggerWrapper, *os.File) {
	var fileName = "app.log"

	lw := logrus.New()
	lw.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	multiWriter := io.MultiWriter(os.Stdout)
	lw.SetOutput(multiWriter)
	loggerWrapper := &LoggerWrapper{
		lw: lw,
	}

	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return loggerWrapper, nil
	}

	multiWriter = io.MultiWriter(os.Stdout, logFile)
	lw.SetOutput(multiWriter)
	loggerWrapper = &LoggerWrapper{
		lw: lw,
	}

	return loggerWrapper, logFile
}

func SetLogger(log Logger) {
	Log = log
}

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	WithFields(args map[string]interface{}) Logger
}

type LoggerWrapper struct {
	lw    *logrus.Logger
	entry *logrus.Entry
}

func (logger *LoggerWrapper) Info(args ...interface{}) {
	if logger.entry != nil {
		logger.entry.Info(args...)
		logger.entry = nil
		return
	}
	logger.lw.Info(args...)
}

func (logger *LoggerWrapper) Infof(format string, args ...interface{}) {
	if logger.entry != nil {
		logger.entry.Infof(format, args...)
		logger.entry = nil
		return
	}
	logger.lw.Infof(format, args...)
}

func (logger *LoggerWrapper) Warn(args ...interface{}) {
	if logger.entry != nil {
		logger.entry.Warn(args...)
		logger.entry = nil
		return
	}
	logger.lw.Warn(args...)
}

func (logger *LoggerWrapper) Warnf(format string, args ...interface{}) {
	if logger.entry != nil {
		logger.entry.Warnf(format, args...)
		logger.entry = nil
		return
	}
	logger.lw.Warnf(format, args...)
}

func (logger *LoggerWrapper) Debug(args ...interface{}) {
	if logger.entry != nil {
		logger.entry.Debug(args...)
		logger.entry = nil
		return
	}
	logger.lw.Debug(args...)
}

func (logger *LoggerWrapper) Debugf(format string, args ...interface{}) {
	if logger.entry != nil {
		logger.entry.Debugf(format, args...)
		logger.entry = nil
		return
	}
	logger.lw.Debugf(format, args...)
}

func (logger *LoggerWrapper) Error(args ...interface{}) {
	if logger.entry != nil {
		logger.entry.Error(args...)
		logger.entry = nil
		return
	}
	logger.lw.Error(args...)
}

func (logger *LoggerWrapper) Errorf(format string, args ...interface{}) {
	if logger.entry != nil {
		logger.entry.Errorf(format, args...)
		logger.entry = nil
		return
	}
	logger.lw.Errorf(format, args...)
}

func (logger *LoggerWrapper) Fatal(args ...interface{}) {
	if logger.entry != nil {
		logger.entry.Fatal(args...)
		logger.entry = nil
		return
	}
	logger.lw.Fatal(args...)
}

func (logger *LoggerWrapper) Fatalf(format string, args ...interface{}) {
	if logger.entry != nil {
		logger.entry.Fatalf(format, args...)
		logger.entry = nil
		return
	}
	logger.lw.Fatalf(format, args...)
}

func (logger *LoggerWrapper) WithFields(args map[string]interface{}) Logger {
	entry := logger.lw.WithFields(args)
	return &LoggerWrapper{
		lw:    logger.lw,
		entry: entry,
	}
}
