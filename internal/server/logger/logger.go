package logging

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/Sofja96/GophKeeper.git/internal/server/settings"
)

type ILogger interface {
	Log() *logrus.Logger
	Error(format string, args ...interface{})
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
}

type logger struct {
	*logrus.Logger
}

func New(conf *settings.Settings) ILogger {
	l := new(logger)
	l.Logger = logrus.New()
	l.Out = os.Stderr

	l.Level = logrus.ErrorLevel
	if conf.Debug {
		l.Level = logrus.DebugLevel
	}

	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	l.SetFormatter(formatter)

	return l

}

func (l *logger) Log() *logrus.Logger {
	return l.Logger
}

func (l *logger) Error(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

func (l *logger) Info(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

func (l *logger) Debug(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}
