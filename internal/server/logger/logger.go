package logging

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/Sofja96/GophKeeper.git/internal/server/settings"
)

// ILogger интерфейс для логирования с использованием библиотеки logrus.
// Обеспечивает методы для записи логов разных уровней.
type ILogger interface {
	Log() *logrus.Logger
	Error(format string, args ...interface{})
	Info(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Fatal(format string, args ...interface{})
}

type logger struct {
	*logrus.Logger
}

// New создает новый экземпляр ILogger на основе конфигурации.
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

// Log возвращает объект logrus.Logger для более гибкого контроля над логированием.
func (l *logger) Log() *logrus.Logger {
	return l.Logger
}

// Error записывает сообщение об ошибке с указанным форматом и параметрами.
func (l *logger) Error(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

// Info записывает информационное сообщение с указанным форматом и параметрами.
func (l *logger) Info(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

// Debug записывает отладочное сообщение с указанным форматом и параметрами.
func (l *logger) Debug(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}

// Fatal записывает сообщение об ошибке и завершает программу с кодом ошибки.
func (l *logger) Fatal(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
}
