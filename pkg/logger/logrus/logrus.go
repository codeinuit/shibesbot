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

func (l LogrusLogger) Error(v ...any) {
	l.logger.Error(v...)
}

func (l LogrusLogger) Warn(v ...any) {
	l.logger.Warn(v...)
}
