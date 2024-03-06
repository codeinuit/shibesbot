package logrus

import "github.com/sirupsen/logrus"

type LogrusLogger struct {
	logger *logrus.Logger
}

func NewLogrusLogger() *LogrusLogger {
	logger := logrus.New()
	return &LogrusLogger{
		logger: logger,
	}
}

func (l LogrusLogger) Info(v ...any) {
	l.logger.Info(v...)
}

func (l LogrusLogger) Infof(format string, v ...any) {
	l.logger.Infof(format, v...)
}

func (l LogrusLogger) Error(v ...any) {
	l.logger.Error(v...)
}

func (l LogrusLogger) Errorf(format string, v ...any) {
	l.logger.Errorf(format, v...)
}

func (l LogrusLogger) Warn(v ...any) {
	l.logger.Warn(v...)
}

func (l LogrusLogger) Warnf(format string, v ...any) {
	l.logger.Warnf(format, v...)
}
